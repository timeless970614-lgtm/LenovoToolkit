//go:build windows

package backend

import (
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

// TestBrightness tests screen brightness read/write capability.
// Returns current brightness, test result, and whether read/write are supported.
func TestBrightness() map[string]interface{} {
	result := make(map[string]interface{})

	// Read current brightness
	orig, err := getBrightness()
	if err != nil {
		result["readable"] = false
		result["writable"] = false
		result["error"] = "无法读取亮度: " + err.Error()
		result["currentBrightness"] = 0
		return result
	}
	result["readable"] = true
	result["currentBrightness"] = orig

	// Try to write a test value
	testLevel := orig
	if orig < 80 {
		testLevel = orig + 10
	} else if orig > 20 {
		testLevel = orig - 10
	} else {
		testLevel = 50
	}

	newLevel, err := setBrightness(testLevel)
	if err != nil {
		result["writable"] = false
		result["error"] = "无法设置亮度: " + err.Error()
		return result
	}
	result["writable"] = true
	result["testBrightness"] = newLevel

	// Restore original brightness after a short delay
	go func() {
		time.Sleep(1500 * time.Millisecond)
		setBrightness(orig)
	}()

	return result
}

// TestFNQ checks FN+Q hotkey support and availability via Lenovo registry.
func TestFNQ() map[string]interface{} {
	result := make(map[string]interface{})

	// Check current mode via DYTC
	info, err := GetDYTCInfo()
	if err == nil && info != nil {
		result["currentMode"] = info.CurrentMode
	} else {
		result["currentMode"] = "N/A"
	}

	// Check FN+Q support via registry and service status
	supported := false
	hotkeyAvailable := false
	var dispatcherFunc string
	var hotkeyStatus string
	dccCapability := false
	geekCapability := false

	script := `
$ErrorActionPreference = 'SilentlyContinue'
$results = @()

# Check Lenovo System Interface Foundation
$paths = @(
    'HKLM:\SOFTWARE\Lenovo\SystemInterfaceFoundation\Hotkey',
    'HKLM:\SOFTWARE\Lenovo\HOTKEY\SETTING',
    'HKLM:\SYSTEM\CurrentControlSet\Services\Lenovo Dispatcher',
    'HKCU:\SOFTWARE\Lenovo\AccessoryOne\Function'
)

foreach ($p in $paths) {
    if (Test-Path $p) {
        $results += "FOUND:$p"
        Get-ItemProperty -Path $p -ErrorAction SilentlyContinue | ForEach-Object {
            $_.PSObject.Properties | Where-Object { $_.Name -notlike 'PS*' } | ForEach-Object {
                $results += "$($_.Name)=$($_.Value)"
            }
        }
    } else {
        $results += "MISSING:$p"
    }
}

# Check Dispatcher service status
$svc = Get-Service -Name 'Lenovo Dispatcher' -ErrorAction SilentlyContinue
if ($svc) {
    $results += "SVC_STATUS:$($svc.Status)"
}

Write-Output ($results -join "|")
`
	out, err := runPowershellScript(script)
	output := strings.TrimSpace(out)

	if err == nil && output != "" {
		lines := strings.Split(output, "|")
		for _, line := range lines {
			kv := strings.SplitN(line, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key := strings.ToLower(kv[0])
			val := strings.TrimSpace(kv[1])

			switch {
			case strings.Contains(key, "fnq") || strings.Contains(key, "hotkey") ||
				strings.Contains(key, "fnf4") || strings.Contains(key, "fnf5") || strings.Contains(key, "fnf6"):
				if val != "" && val != "0" && val != "false" {
					supported = true
				}
			case key == "hkey" || key == "key":
				if val != "" && val != "0" {
					supported = true
					dispatcherFunc = val
				}
			case strings.HasPrefix(line, "SVC_STATUS:"):
				hotkeyStatus = val
				hotkeyAvailable = (val == "Running")
			case strings.HasPrefix(line, "FOUND:"):
				supported = true
			}
		}
	}

	// Check DCC and GEEK capabilities
	dccCapability = CheckDCC()
	geekCapability = CheckGEEK()

	// FN+Q is supported if we found any Lenovo hotkey registry entries
	if supported || dccCapability || geekCapability {
		supported = true
	}

	result["supported"] = supported
	result["hotkeyAvailable"] = hotkeyAvailable
	result["dispatcherFunction"] = dispatcherFunc
	result["hotkeyStatus"] = hotkeyStatus
	result["dccCapability"] = dccCapability
	result["geekCapability"] = geekCapability

	return result
}

// TestModeSwitch attempts to cycle through available DYTC modes and verify switching.
// Returns current mode, available modes, and per-mode attempt results.
func TestModeSwitch() map[string]interface{} {
	result := make(map[string]interface{})

	// Get current mode
	info, err := GetDYTCInfo()
	if err != nil {
		result["error"] = "无法获取 DYTC 信息: " + err.Error()
		return result
	}
	origModeName := info.CurrentMode
	result["currentMode"] = origModeName

	// Get available modes from FUNC_CAP
	supported := GetFuncCapSupportedModes()
	availableModes := []string{}
	for mode := range supported {
		availableModes = append(availableModes, mode)
	}
	result["availableModes"] = availableModes

	// Mode priority list (skip the current mode)
	modeOrder := []string{"STD", "IBSM", "AQM", "APM", "IEPM", "DCC"}
	attempts := []map[string]interface{}{}

	for _, modeName := range modeOrder {
		if modeName == origModeName {
			continue
		}
		// Skip unsupported modes
		if ok, _ := IsModeSupportedByFuncCap(modeName); !ok {
			continue
		}

		attempt := map[string]interface{}{
			"mode":     modeName,
			"success":  false,
			"verified": false,
		}

		// Try to set the mode
		_, err := SetDYTCModeByName(modeName)
		if err == nil {
			attempt["success"] = true

			// Brief pause then verify
			time.Sleep(800 * time.Millisecond)

			// Verify by re-reading current mode
			info2, err2 := GetDYTCInfo()
			if err2 == nil {
				verified := (info2.CurrentMode == modeName)
				attempt["verified"] = verified
				attempt["actualMode"] = info2.CurrentMode
			}
		} else {
			attempt["error"] = err.Error()
		}

		attempts = append(attempts, attempt)

		// Restore to original mode after testing
		time.Sleep(300 * time.Millisecond)
		SetDYTCModeByName(origModeName)
		time.Sleep(500 * time.Millisecond)
	}

	result["attempts"] = attempts
	return result
}

// OpenTestMode enables Windows test mode for driver testing.
// This allows loading unsigned/test-signed drivers.
func OpenTestMode() map[string]interface{} {
	// Execute bcdedit to enable test signing
	cmd := exec.Command("bcdedit", "/set", "testsigning", "on")
	cmd.SysProcAttr = &windows.SysProcAttr{
		CreationFlags: windows.CREATE_NO_WINDOW,
	}
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		return map[string]interface{}{
			"success": false,
			"message": "执行失败: " + outputStr,
		}
	}

	// Check if operation was successful
	if strings.Contains(outputStr, "操作成功完成") || strings.Contains(strings.ToLower(outputStr), "success") {
		return map[string]interface{}{
			"success": true,
			"message": "测试模式已启用！请重启电脑生效。\n\n桌面右下角将显示'测试模式'水印。",
		}
	}

	return map[string]interface{}{
		"success": false,
		"message": "未知结果: " + outputStr,
	}
}
