//go:build windows

package backend

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

// DYTC DLL function declarations
var (
	modDYTCDll                 = windows.NewLazyDLL("LenovoDYTC.dll")
	procSetFanMode              = modDYTCDll.NewProc("Set_FanMode")
	procSetDYTCMode             = modDYTCDll.NewProc("Set_DYTCMode")
	procGetCapDCC              = modDYTCDll.NewProc("GET_Cap_DCC")
	procSetGEEKMode            = modDYTCDll.NewProc("Set_GEEKMode")
	procGetCapGEEK             = modDYTCDll.NewProc("GET_Cap_GEEK")
	procSetODVMode             = modDYTCDll.NewProc("Set_ODVMode")
	procGetDYTCCmdDispatcherFUNC = modDYTCDll.NewProc("Get_DYTC_CMD_MODE_DISPATCHERFUNCTION")
	procGetDYTCCmdFuncCap     = modDYTCDll.NewProc("Get_DYTC_CMD_FUNC_CAP")
	procGetDYTCCmdNITThreshold = modDYTCDll.NewProc("Get_DYTC_CMD_MODE_NIT_DISPATCHERTHRESHOLD")
	procGetDYTCCmdNITGet       = modDYTCDll.NewProc("Get_DYTC_CMD_MODE_NIT_GET")
)

// DYTC Mode constants
const (
	DYTC_MODE_BSM  uint32 = 1
	DYTC_MODE_IBSM uint32 = 2
	DYTC_MODE_AQM  uint32 = 3
	DYTC_MODE_STD  uint32 = 4
	DYTC_MODE_APM  uint32 = 5
	DYTC_MODE_IEPM uint32 = 6
	DYTC_MODE_EPM  uint32 = 7
	DYTC_MODE_DCC  uint32 = 13
)

// FUNC_CAP bitmap masks (bitmap = funcCapRet >> 16)
// Values from C++ FuncAbility definitions (decimal mask values, not bit positions)
var funcCapModeMasks = map[string]uint32{
	"STD":  1,     // STD_MASK
	"MYH":  8,     // MYH_MASK (Multi Yoga Hinge)
	"STP":  16,    // STP_MASK (Skin Temperature Protection)
	"APM":  32,    // APM_MASK (Auto Performance Mode)
	"AQM":  64,    // AQM_MASK (Auto Quiet Mode)
	"IEPM": 128,   // IEPM_MASK (Intelligent Extreme Performance Mode)
	"IBSM": 256,   // IBSM_MASK (Intelligent Battery Saving Mode)
	"CQL":  512,   // CQL_MASK (Cool Quiet on Laptop)
	"AAA":  1024,  // AAA_MASK (Ambient Acoustic Adaptive)
	"EPM":  2048,  // EPM_MASK
	"MSC":  4096,  // MSC_MASK (Modern Standby Control)
	"AIM":  8192,  // AIM_MASK
	"DCC":  32768, // DCC_MASK
}

// funcCapCachedBitmap caches the FUNC_CAP bitmap on first read
var (
	funcCapOnce    sync.Once
	funcCapBitmap  uint32
	funcCapValid   bool
)

// getFuncCapBitmap reads and caches the FUNC_CAP bitmap (bits 16-31 of Get_DYTC_CMD_FUNC_CAP)
func getFuncCapBitmap() (uint32, bool) {
	funcCapOnce.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[FUNC_CAP] DLL call failed (procedure not found or driver unavailable): %v", r)
				// funcCapValid stays false — callers will gracefully allow attempts
			}
		}()
		ret, _, _ := procGetDYTCCmdFuncCap.Call()
		rawVal := uint32(ret)
		if rawVal&1 != 0 { // bit 0 = valid
			funcCapBitmap = rawVal >> 16
			funcCapValid = true
		}
	})
	return funcCapBitmap, funcCapValid
}

// GetFuncCapSupportedModes returns a map of DYTC mode name → supported (bool)
// based on the FUNC_CAP bitmap from the DLL.
func GetFuncCapSupportedModes() map[string]bool {
	result := make(map[string]bool)
	bitmap, valid := getFuncCapBitmap()
	if !valid {
		// DLL call failed or returned invalid — report nothing as supported
		return result
	}
	for modeName, mask := range funcCapModeMasks {
		result[modeName] = (bitmap & mask) != 0
	}
	// BSM has mask=0 (receive only) — implicitly enabled when EPM/AIM/MSC/DCC is set
	// (per C++ source: EPM/AIM/MSC/DCC all set funcCap.bBSM = true)
	result["BSM"] = ((bitmap & 2048) != 0) || ((bitmap & 8192) != 0) || ((bitmap & 4096) != 0) || ((bitmap & 32768) != 0)
	// GEEK is not in FUNC_CAP bitmap — check via dedicated DLL call
	result["GEEK"] = GetCapGEEK() != 0
	return result
}

// IsModeSupportedByFuncCap checks if a DYTC mode is supported according to FUNC_CAP bitmap.
// Returns (supported bool, reason string). If FUNC_CAP is unavailable, returns (true, "")
// to allow the attempt (don't block when we can't verify).
func IsModeSupportedByFuncCap(modeName string) (bool, string) {
	bitmap, valid := getFuncCapBitmap()
	if !valid {
		// Cannot verify — allow attempt
		return true, ""
	}
	mask, exists := funcCapModeMasks[modeName]
	if !exists {
		// Unknown mode — allow attempt
		return true, ""
	}
	if bitmap&mask == 0 {
		return false, fmt.Sprintf("本设备不支持 %s 模式 (FUNC_CAP bitmap 未设置对应位)", modeName)
	}
	return true, ""
}

// DYTCInfo holds all DYTC related information
type DYTCInfo struct {
	CurrentMode           string               `json:"currentMode"`
	CurrentDispatcherMode string               `json:"currentDispatcherMode"`
	DCCCapability        uint32               `json:"dccCapability"`
	GEEKCapability       uint32               `json:"geekCapability"`
	AIEngineMode         string               `json:"aiEngineMode"`
	DispatcherFunction   uint32               `json:"dispatcherFunction"`
	DispatcherThreshold  uint32               `json:"dispatcherThreshold"`
	EnableFunc           uint32               `json:"enableFunc"`
	FuncCap              uint32               `json:"funcCap"`
	Nits                 uint32               `json:"nits"`
	DispatcherFeatures   []DispatcherFeature  `json:"dispatcherFeatures"`
}

// DispatcherFeature represents a single feature bit from DISPATCHER_FUNCTION
type DispatcherFeature struct {
	Bit     int    `json:"bit"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Enabled bool   `json:"enabled"`
	Group   string `json:"group"`
}

// ParseDispatcherFunction parses the DISPATCHER_FUNCTION value into feature list
func ParseDispatcherFunction(val uint32) []DispatcherFeature {
	defs := []struct {
		bit   int
		name  string
		desc  string
		group string
	}{
		{0, "Status", "Dispatcher status bit", "Core"},
		{4, "StopPPM", "Stop PPM / Use DTT EPO", "Core"},
		{5, "BIOSConfigCPUTH", "BIOS config CPU threshold", "Core"},
		{6, "BIOSConfigLatencyTH", "BIOS config latency threshold", "Core"},
		{7, "SendFPS2AI", "Send FPS to AI chip", "AI"},
		{8, "LowFreqforMM25", "Low frequency for MobileMark 2025", "Perf"},
		{9, "UsingODV", "Using ODV (OverDrive Voltage)", "Perf"},
		{10, "HiddenIEPMEnableAIGC", "Hidden IEPM enable AIGC", "AI"},
		{11, "BootUpDelay", "Boot-up delay enabled", "Core"},
		{12, "MouseMoveTurbo", "Mouse move turbo boost", "Turbo"},
		{13, "EPMDynamicEPP", "EPM dynamic EPP adjustment", "Perf"},
		{14, "APMCreat_BG_Turbo", "APM background process turbo", "Turbo"},
		{15, "AQMAPPtrubo2APM", "AQM app turbo to APM mode", "Turbo"},
		{16, "EnableSSD_b0", "SSD whitelist perf (bit 0)", "SSD"},
		{17, "EnableSSD_b1", "SSD whitelist perf (bit 1)", "SSD"},
		{18, "EnableSSD_b2", "SSD whitelist perf (bit 2)", "SSD"},
		{19, "EnableSSD_b3", "SSD whitelist perf (bit 3)", "SSD"},
		{20, "AutoMemoryClean", "Auto memory clean", "Memory"},
		{21, "FNQ2WaySync", "FN+Q 2-way sync", "UI"},
		{22, "EnableBSMTurbo", "Enable BSM turbo", "Turbo"},
		{24, "EnableSAGV_b0", "Enable SAGV (bit 0)", "Memory"},
		{25, "EnableSAGV_b1", "Enable SAGV (bit 1)", "Memory"},
		{26, "EnableDolbyCtrl", "Enable Dolby control", "Audio"},
		{27, "EnableDNPUTurbo", "Enable DNPU turbo", "AI"},
		{28, "SendOEMV_DTTReconnTest", "DTT reconnect test mode", "Debug"},
		{29, "SendOEMV_DTTTest", "DTT test mode (1s interval)", "Debug"},
		{30, "EnableDGPUPlug", "Enable dGPU plug/unplug control", "GPU"},
		{31, "NVGPUOC", "NVIDIA GPU overclocking", "GPU"},
	}
	features := make([]DispatcherFeature, 0, len(defs))
	for _, d := range defs {
		enabled := (val>>uint(d.bit))&1 == 1
		features = append(features, DispatcherFeature{
			Bit: d.bit, Name: d.name, Desc: d.desc, Enabled: enabled, Group: d.group,
		})
	}
	return features
}

// SetFanMode sets the fan mode
func SetFanMode(mode uint32) error {
	ret, _, err := procSetFanMode.Call(uintptr(mode))
	if ret == 0 && err != nil {
		return fmt.Errorf("SetFanMode failed: %v", err)
	}
	return nil
}

// SetDYTCMode sets the DYTC mode
func SetDYTCMode(mode uint32) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("SetDYTCMode DLL call failed: %v", r)
			log.Printf("[DYTC] SetDYTCMode panic: %v", r)
		}
	}()
	ret, _, err := procSetDYTCMode.Call(uintptr(mode))
	if ret == 0 && err != nil {
		return fmt.Errorf("SetDYTCMode failed: %v", err)
	}
	return nil
}

// GetCapDCC gets DCC capability
func GetCapDCC() (result uint32) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[DYTC] GetCapDCC panic: %v", r)
		}
	}()
	ret, _, _ := procGetCapDCC.Call()
	return uint32(ret)
}

// SetGEEKMode sets GEEK mode on/off
func SetGEEKMode(onoff bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("SetGEEKMode DLL call failed: %v", r)
			log.Printf("[DYTC] SetGEEKMode panic: %v", r)
		}
	}()
	var mode uint32
	if onoff {
		mode = 1
	}
	ret, _, err := procSetGEEKMode.Call(uintptr(mode))
	if ret == 0 && err != nil {
		return fmt.Errorf("SetGEEKMode failed: %v", err)
	}
	return nil
}

// GetCapGEEK gets GEEK capability
func GetCapGEEK() (result uint32) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[DYTC] GetCapGEEK panic: %v", r)
		}
	}()
	ret, _, _ := procGetCapGEEK.Call()
	return uint32(ret)
}

// SetODVMode sets ODV mode
// For ODV index 30-33 (ODV1 group), constructs DYTC_CMD_DISPATCHERODV1_UNION:
//   CMD_ID = 16 (bits 0-8)
//   MASK_ODV3x = 1 (bits 12-15 depending on index)
//   Value_ODV3x = value (4-bit fields at bits 16-31)
// For ODV index 34-41 (ODV2 group), constructs DYTC_CMD_DISPATCHERODV2_UNION:
//   CMD_ID = 17 (bits 0-8)
//   MASK_ODV3x = 1 (bits 12-15 depending on index)
//   Value_ODV3x = value (4-bit fields at bits 16-31)
func SetODVMode(index, value uint32) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("SetODVMode DLL call failed: %v", r)
			log.Printf("[DYTC] SetODVMode panic: %v", r)
		}
	}()
	var dwIoData uint32

	if index >= 30 && index <= 33 {
		// ODV1 group: CMD_ID=16
		// bit 0-8: CMD_ID = 16
		// bit 12: MASK_ODV30, bit 13: MASK_ODV31, bit 14: MASK_ODV32, bit 15: MASK_ODV33
		// bit 16-19: Value_ODV30, bit 20-23: Value_ODV31, bit 24-27: Value_ODV32, bit 28-31: Value_ODV33
		maskBit := uint32(12 + (index - 30)) // MASK bit position
		valueShift := uint32(16 + (index-30)*4) // Value field start bit
		dwIoData = 16 // CMD_ID = DYTC_CMD_DISPATCHERODV1
		dwIoData |= 1 << maskBit // Set MASK bit
		dwIoData |= (value & 0xF) << valueShift // Set Value field (4 bits)
	} else if index >= 34 && index <= 41 {
		// ODV2 group: CMD_ID=17
		// bit 0-8: CMD_ID = 17
		// bit 12: MASK_ODV34, bit 13: MASK_ODV35, bit 14: MASK_ODV36, bit 15: MASK_ODV37
		// bit 16-19: Value_ODV34, bit 20-23: Value_ODV35, bit 24-27: Value_ODV36, bit 28-31: Value_ODV37
		maskBit := uint32(12 + (index - 34))
		valueShift := uint32(16 + (index-34)*4)
		dwIoData = 17 // CMD_ID = DYTC_CMD_DISPATCHERODV2
		dwIoData |= 1 << maskBit
		dwIoData |= (value & 0xF) << valueShift
	} else {
		// Unknown index range — pass value as-is (legacy behavior)
		dwIoData = value
	}

	log.Printf("[ODV] SetODVMode(index=%d, value=%d) -> dwIoData=0x%08X", index, value, dwIoData)
	ret, _, err := procSetODVMode.Call(uintptr(index), uintptr(dwIoData))
	if ret == 0 && err != nil {
		return fmt.Errorf("SetODVMode failed: %v", err)
	}
	return nil
}

// GetDYTCCmdDispatcherFUNC gets dispatcher function
func GetDYTCCmdDispatcherFUNC() (result uint32) {
	defer func() { if r := recover(); r != nil { log.Printf("[DYTC] GetDYTCCmdDispatcherFUNC panic: %v", r) } }()
	ret, _, _ := procGetDYTCCmdDispatcherFUNC.Call()
	return uint32(ret)
}

// GetDYTCCmdFuncCap gets function capability
func GetDYTCCmdFuncCap() (result uint32) {
	defer func() { if r := recover(); r != nil { log.Printf("[DYTC] GetDYTCCmdFuncCap panic: %v", r) } }()
	ret, _, _ := procGetDYTCCmdFuncCap.Call()
	return uint32(ret)
}

// GetDYTCCmdNITThreshold gets NIT threshold
func GetDYTCCmdNITThreshold() (result uint32) {
	defer func() { if r := recover(); r != nil { log.Printf("[DYTC] GetDYTCCmdNITThreshold panic: %v", r) } }()
	ret, _, _ := procGetDYTCCmdNITThreshold.Call()
	return uint32(ret)
}

// GetDYTCCmdNITGet gets NIT value
func GetDYTCCmdNITGet() (result uint32) {
	defer func() { if r := recover(); r != nil { log.Printf("[DYTC] GetDYTCCmdNITGet panic: %v", r) } }()
	ret, _, _ := procGetDYTCCmdNITGet.Call()
	return uint32(ret)
}

// GetDYTCModeName returns human-readable name for DYTC mode
func GetDYTCModeName(mode uint32) string {
	names := map[uint32]string{
		1:  "BSM (Basic Performance)",
		2:  "IBSM (Intelligent Basic)",
		3:  "AQM (Active Performance)",
		4:  "STD (Standard Performance)",
		5:  "APM (Advanced Performance)",
		6:  "IEPM (Intelligent Performance)",
		7:  "EPM (Extreme Performance)",
		13: "DCC (Dynamic Performance Control)",
	}
	if name, ok := names[mode]; ok {
		return name
	}
	return fmt.Sprintf("Unknown Mode (%d)", mode)
}

// GetDispatcherModeName returns human-readable name for dispatcher mode
func GetDispatcherModeName(mode int) string {
	names := map[int]string{
		1:  "Battery Saving",
		2:  "Intelligent Battery Saving",
		3:  "Intelligent Auto Quiet",
		4:  "Intelligent Stand Mode",
		5:  "Intelligent Auto Performance",
		6:  "Intelligent Extreme",
		7:  "Extreme Performance",
		8:  "Yoga Tablet",
		9:  "Yoga Tent",
		10: "Yoga Flat",
		11: "Geek Mode",
	}
	if name, ok := names[mode]; ok {
		return name
	}
	return fmt.Sprintf("Unknown (%d)", mode)
}

// GetDYTCInfo returns all DYTC information
func GetDYTCInfo() (*DYTCInfo, error) {
	info := &DYTCInfo{}

	// Read all from registry — no DLL calls (avoids blocking on some machines)
	// Capabilities are cached by Dispatcher service from BIOS on startup
	info.DCCCapability = readDispatcherReg("ITS_DCCCapability", 0)
	info.GEEKCapability = readDispatcherReg("ITS_GEECKCapability", 0)

	// If not in registry, derive from DispatcherFunction
	if info.DCCCapability == 0 {
		dytcFunc := readDispatcherReg("Policy_DispatcherFunc", 0)
		// Bit 3 in DispatcherFunction indicates DCC support
		info.DCCCapability = (dytcFunc >> 3) & 1
	}
	if info.GEEKCapability == 0 {
		info.GEEKCapability = 1 // Assume supported if not set
	}

	info.DispatcherFunction = readDispatcherReg("Policy_DispatcherFunc", 0)
	info.DispatcherThreshold = readDispatcherReg("Policy_NITThreshold", 0)
	info.DispatcherFeatures = ParseDispatcherFunction(info.DispatcherFunction)

	enableFuncReg := readDispatcherReg("Policy_EnableFunc", 0xFFFF)
	if enableFuncReg == 0xFFFF {
		info.EnableFunc = info.DispatcherFunction
	} else {
		info.EnableFunc = enableFuncReg
	}

	currentMode := readDispatcherRegInt("ITS_AutomaticModeSetting", -1)
	info.CurrentDispatcherMode = GetDispatcherModeName(currentMode)

	aiEngine := readDispatcherReg("Policy_AIEngine", 0)
	if aiEngine == 0xc {
		info.AIEngineMode = "AI Engine"
	} else {
		info.AIEngineMode = "CPU Engine"
	}

	// Read FUNC_CAP and NITS via DLL calls
	info.FuncCap = GetDYTCCmdFuncCap()
	info.Nits = GetDYTCCmdNITGet()

	return info, nil
}

// SetDYTCModeByName sets DYTC mode by name
func SetDYTCModeByName(modeName string) (string, error) {
	modes := map[string]uint32{
		"BSM": DYTC_MODE_BSM, "IBSM": DYTC_MODE_IBSM, "AQM": DYTC_MODE_AQM,
		"STD": DYTC_MODE_STD, "APM": DYTC_MODE_APM, "IEPM": DYTC_MODE_IEPM,
		"EPM": DYTC_MODE_EPM, "DCC": DYTC_MODE_DCC,
	}
	if mode, ok := modes[modeName]; ok {
		// Check FUNC_CAP bitmap — block unsupported modes
		if supported, reason := IsModeSupportedByFuncCap(modeName); !supported {
			return "", fmt.Errorf("%s", reason)
		}
		err := SetDYTCMode(mode)
		if err != nil {
			return "", err
		}
		return GetDYTCModeName(mode), nil
	}
	if modeName == "GEEK" {
		geekCap := GetCapGEEK()
		if geekCap == 0 {
			return "", fmt.Errorf("GEEK mode not supported on this device")
		}
		err := SetGEEKMode(true)
		if err != nil {
			return "", err
		}
		return "GEEK Mode", nil
	}
	return "", fmt.Errorf("unknown mode: %s", modeName)
}

// SetODV sets ODV mode with index and value
func SetODV(index, value uint32) error {
	return SetODVMode(index, value)
}

// CheckBSM checks if device is BSM capable
func CheckBSM() bool {
	return GetCapDCC() != 0
}

// CheckGEEK checks if device is GEEK mode capable
func CheckGEEK() bool {
	return GetCapGEEK() != 0
}

// CheckDCC checks if device is DCC capable
func CheckDCC() bool {
	return GetCapDCC() != 0
}

// readDispatcherReg reads a DWORD from the Dispatcher registry path
func readDispatcherReg(name string, defaultVal uint32) uint32 {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.QUERY_VALUE)
	if err != nil {
		return defaultVal
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue(name)
	if err != nil {
		return defaultVal
	}
	return uint32(val)
}

// readDispatcherRegInt reads a DWORD from the Dispatcher registry path as int
func readDispatcherRegInt(name string, defaultVal int) int {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, registryPath, registry.QUERY_VALUE)
	if err != nil {
		return defaultVal
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue(name)
	if err != nil {
		return defaultVal
	}
	return int(val)
}
