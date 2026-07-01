//go:build windows

package backend

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sys/windows/registry"
)

// callWithTimeout runs f in a goroutine and returns after timeout if f hasn't finished.
// This prevents DLL calls from blocking forever if the driver is unresponsive.
func callWithTimeout(timeout time.Duration, f func() error) error {
	ch := make(chan error, 1)
	go func() {
		ch <- f()
	}()
	select {
	case err := <-ch:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("operation timed out after %v", timeout)
	}
}

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
//
// Entire function is capped at 8s via context timeout to prevent
// the Wails JS bridge from blocking indefinitely.
func PinDYTCMode(modeId string) (result string, err error) {
	// Top-level recover: any DLL panic is caught here and returned as error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("DLL call failed: %v", r)
			log.Printf("[PinMode] Panic recovered: %v", r)
			result = ""
		}
	}()

	start := time.Now()
	val, ok := modeToValue[modeId]
	if !ok {
		return "", fmt.Errorf("unknown mode: %s", modeId)
	}

	// Check FUNC_CAP bitmap - block unsupported modes
	// If FUNC_CAP DLL call fails (procedure not found), this gracefully allows the attempt
	if supported, reason := IsModeSupportedByFuncCap(modeId); !supported {
		return "", fmt.Errorf("%s", reason)
	}

	// 8s hard timeout for the entire function (frontend uses 10s)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	// Step 1: Stop Dispatcher service (best-effort, non-blocking)
	step1Start := time.Now()
	if isDispatcherRunning() {
		log.Printf("[PinMode] Stopping Dispatcher service before pinning...")
		stopDone := make(chan error, 1)
		go func() { stopDone <- stopDispatcherService() }()
		select {
		case err := <-stopDone:
			if err != nil {
				log.Printf("[PinMode] Warning: stop service failed: %v (%v)", err, time.Since(step1Start))
			} else {
				log.Printf("[PinMode] Step1 (stop service) took %v", time.Since(step1Start))
			}
		case <-ctx.Done():
			log.Printf("[PinMode] Step1 (stop service) skipped: context timeout (%v)", time.Since(step1Start))
		}
	}

	// Step 2: Write registry values (fast, no timeout needed)
	step2Start := time.Now()
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
	if err := k.SetDWordValue("Policy_Override", 3); err != nil {
		return "", fmt.Errorf("failed to write Policy_Override: %w", err)
	}
	log.Printf("[PinMode] Step2 (registry write) took %v", time.Since(step2Start))

	// Step 3: DLL call with context cancellation support
	step3Start := time.Now()
	dllDone := make(chan error, 1)
	if modeId == "GEEK" {
		go func() { dllDone <- SetGEEKMode(true) }()
	} else {
		go func() { dllDone <- SetDYTCMode(val) }()
	}
	select {
	case err := <-dllDone:
		if err != nil {
			log.Printf("[PinMode] DLL call failed: %v (%v)", err, time.Since(step3Start))
		} else {
			log.Printf("[PinMode] Step3 (DLL call) succeeded (%v)", time.Since(step3Start))
		}
	case <-ctx.Done():
		log.Printf("[PinMode] Step3 (DLL call) timed out after %v", time.Since(step3Start))
	}

	// Step 4: Send ODV33 (only if context not expired)
	step4Start := time.Now()
	odvMsg := ""
	select {
	case <-ctx.Done():
		log.Printf("[PinMode] Skipping ODV33: context already expired")
		odvMsg = " (ODV33 skipped: timeout)"
	default:
		odvDone := make(chan error, 1)
		go func() { odvDone <- SetODVMode(33, 15) }()
		select {
		case err := <-odvDone:
			if err != nil {
				log.Printf("[PinMode] SetODVMode(33,15) failed: %v (%v)", err, time.Since(step4Start))
				odvMsg = " (ODV33=UserScenario send failed)"
			} else {
				log.Printf("[PinMode] Step4 (ODV33) succeeded (%v)", time.Since(step4Start))
				odvMsg = " | ODV33 UserScenario=15 sent"
			}
		case <-ctx.Done():
			log.Printf("[PinMode] Step4 (ODV33) timed out after %v", time.Since(step4Start))
			odvMsg = " (ODV33 skipped: timeout)"
		}
	}

	log.Printf("[PinMode] Total PinDYTCMode took %v (ctx_err=%v)", time.Since(start), ctx.Err())
	return odvMsg, nil
}

// UnpinDYTCMode removes the pin by:
// 1. Restoring Policy_Override to 0 (auto mode)
// 2. Restarting the Dispatcher service so it resumes auto-switching
func UnpinDYTCMode() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("DLL call failed: %v", r)
			log.Printf("[PinMode] UnpinDYTCMode panic recovered: %v", r)
		}
	}()
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
		log.Printf("[PinMode] SetODVMode(33, 0) succeeded - ODV33 UserScenario reset to 0")
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
