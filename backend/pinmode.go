//go:build windows

package backend

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

// modeToValue maps DYTC mode name to ITS_AutomaticModeSetting value
var modeToValue = map[string]uint32{
	"BSM": 1, "IBSM": 2, "AQM": 3, "STD": 4,
	"APM": 5, "IEPM": 6, "EPM": 7, "GEEK": 11, "DCC": 13,
}

var valueToMode = map[uint32]string{
	1: "BSM", 2: "IBSM", 3: "AQM", 4: "STD",
	5: "APM", 6: "IEPM", 7: "EPM", 11: "GEEK", 13: "DCC",
}

// GetPinnedDYTCMode returns the currently pinned mode name, or "" if not pinned.
// Pinned = Policy_Override == 3
func GetPinnedDYTCMode() string {
	override := readModeCheckReg("Policy_Override", 0)
	if override != 3 {
		return ""
	}
	mode := readModeCheckReg("ITS_AutomaticModeSetting", 0)
	if name, ok := valueToMode[mode]; ok {
		return name
	}
	return fmt.Sprintf("Mode %d", mode)
}

// PinDYTCMode pins the given mode by:
// 1. Stopping the Dispatcher service (to prevent it from overriding our settings)
// 2. Writing ITS_AutomaticModeSetting + Policy_Override=3 to registry
// 3. Calling Set_DYTCMode DLL to actually switch the hardware mode
func PinDYTCMode(modeId string) (string, error) {
	val, ok := modeToValue[modeId]
	if !ok {
		return "", fmt.Errorf("unknown mode: %s", modeId)
	}

	// Check FUNC_CAP bitmap — block unsupported modes
	if supported, reason := IsModeSupportedByFuncCap(modeId); !supported {
		return "", fmt.Errorf("%s", reason)
	}

	// Step 1: Stop Dispatcher service to prevent it from overriding our settings
	if isDispatcherRunning() {
		log.Printf("[PinMode] Stopping Dispatcher service before pinning...")
		if err := stopDispatcherService(); err != nil {
			log.Printf("[PinMode] Warning: failed to stop Dispatcher service: %v", err)
			// Continue anyway - we'll still try to set the mode
		}
	}

	// Step 2: Write registry values
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Services\LenovoProcessManagement\Performance\PowerSlider`,
		registry.SET_VALUE,
	)
	if err != nil {
		return "", fmt.Errorf("cannot open registry for write: %w", err)
	}
	defer k.Close()

	if err := k.SetDWordValue("ITS_AutomaticModeSetting", val); err != nil {
		return "", fmt.Errorf("failed to write ITS_AutomaticModeSetting: %w", err)
	}
	// Policy_Override = 3 -> fixed mode, Dispatcher won't auto-switch on restart
	if err := k.SetDWordValue("Policy_Override", 3); err != nil {
		return "", fmt.Errorf("failed to write Policy_Override: %w", err)
	}
	log.Printf("[PinMode] Registry set: ITS_AutomaticModeSetting=%d, Policy_Override=3", val)

	// Step 3: Call the appropriate DLL to actually switch the hardware mode
	if modeId == "GEEK" {
		if err := SetGEEKMode(true); err != nil {
			log.Printf("[PinMode] Warning: SetGEEKMode DLL call failed: %v", err)
		} else {
			log.Printf("[PinMode] SetGEEKMode(true) DLL call succeeded")
		}
	} else {
		if err := SetDYTCMode(val); err != nil {
			log.Printf("[PinMode] Warning: SetDYTCMode DLL call failed: %v", err)
			// Don't return error - registry is already set, DLL may fail on some machines
			// The mode will still be applied when the service restarts
		} else {
			log.Printf("[PinMode] SetDYTCMode(%d) DLL call succeeded", val)
		}
	}

	// Step 4: Send ODV33 (UserScenario) to DTT via DYTC CMD DISPATCHERODV1 union
	odvMsg := ""
	if err := SetODVMode(33, 15); err != nil {
		log.Printf("[PinMode] Warning: SetODVMode(33, 15) failed: %v", err)
		odvMsg = " (ODV33=UserScenario send failed)"
	} else {
		log.Printf("[PinMode] SetODVMode(33, 15) succeeded — ODV33=UserScenario=15 sent via DISPATCHERODV1 union (0xF0008010)")
		odvMsg = " | ODV33 UserScenario=15 sent"
	}

	return odvMsg, nil
}

// UnpinDYTCMode removes the pin by:
// 1. Restoring Policy_Override to 0 (auto mode)
// 2. Restarting the Dispatcher service so it resumes auto-switching
func UnpinDYTCMode() error {
	// Step 1: Reset Policy_Override in registry
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Services\LenovoProcessManagement\Performance\PowerSlider`,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("cannot open registry for write: %w", err)
	}
	defer k.Close()

	if err := k.SetDWordValue("Policy_Override", 0); err != nil {
		return fmt.Errorf("failed to write Policy_Override: %w", err)
	}
	log.Printf("[PinMode] Registry set: Policy_Override=0 (auto mode)")

	// Step 2: Reset ODV33 (UserScenario) to default (0) when unpinning
	if err := SetODVMode(33, 0); err != nil {
		log.Printf("[PinMode] Warning: SetODVMode(33, 0) reset failed: %v", err)
	} else {
		log.Printf("[PinMode] SetODVMode(33, 0) succeeded — ODV33 UserScenario reset to 0 via DISPATCHERODV1 union (0x00008010)")
	}

	// Step 3: Start Dispatcher service so it resumes auto mode switching
	if !isDispatcherRunning() {
		log.Printf("[PinMode] Starting Dispatcher service for auto mode...")
		if err := startDispatcherService(); err != nil {
			log.Printf("[PinMode] Warning: failed to start Dispatcher service: %v", err)
		} else {
			log.Printf("[PinMode] Dispatcher service started, auto mode switching resumed")
		}
	}

	return nil
}

// isDispatcherRunning checks if the Lenovo Dispatcher service is currently running
func isDispatcherRunning() bool {
	status, err := GetServiceStatus()
	if err != nil {
		log.Printf("[PinMode] Warning: failed to query service status: %v", err)
		return false
	}
	return status == "Running"
}

// stopDispatcherService stops the Lenovo Dispatcher service
func stopDispatcherService() error {
	return StopService()
}

// startDispatcherService starts the Lenovo Dispatcher service
func startDispatcherService() error {
	return StartService()
}
