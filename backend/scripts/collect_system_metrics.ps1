<# collect_system_metrics.ps1
# Collects system metrics for ML Scenario log.
# Returns 53 pipe-delimited fields for CSV injection (cols 3-55, 0-indexed).
#
# Field order (matching mlscenario.go mlCollectRow):
#   cpuUsage|cpuFreq|cpuPerf|gpuTotal|igpuUsage|dgpuUsage|vpuUsage|npuUsage|
#   iGPUID|gGPUID|vpuID|memUsage|memFree|inputLat|lag64|lag100|lag200|
#   diskUsage|diskSpeed|diskReadLat|diskWriteLat|
#   sysPower|cpuPower|gpuPower|nvidiaPower|nvidiaTemp|
#   copy|gdiRender|legacyOverlay|security|3d|videoDec|videoEnc|videoProc|unknown|compute|
#   itsMode|battCap|fgApp|memSpeed|igpuVram|igpuShare|igpuTotal
#
# Note: AC_DC and PowerSlider are injected by Go, not collected here.
#>

$ErrorActionPreference = 'SilentlyContinue'

# ── Battery (for Battery_Capacity) ──────────────────────────────────────
$batt = Get-WmiObject Win32_Battery -ErrorAction SilentlyContinue

# ── CPU Usage ────────────────────────────────────────────────────────────
$cpu = Get-WmiObject Win32_PerfFormattedData_PerfOS_Processor |
       Where-Object { $_.Name -eq '_Total' }
if ($null -ne $cpu) {
    $cpuUsage = [math]::Round($cpu.PercentProcessorTime, 4)
} else {
    $cpuUsage = 0
}

# ── CPU Frequency ────────────────────────────────────────────────────────
$cpuProc = Get-WmiObject Win32_Processor -ErrorAction SilentlyContinue | Select-Object -First 1
if ($null -ne $cpuProc) {
    $cpuFreq = [math]::Round($cpuProc.CurrentClockSpeed, 0)
    $cpuPerf = [math]::Round($cpuUsage * $cpuFreq / 1000, 3)
} else {
    $cpuFreq = 0
    $cpuPerf = 0
}

# ── GPU Usage (Intel GPU Engine counters) ────────────────────────────────
$gpuCounters = Get-Counter '\GPU Engine(*)\Utilization Percentage' -ErrorAction SilentlyContinue
$igpuUsage = 0.0
$dgpuUsage = 0.0
$gpuTotal  = 0.0
$igpu3D    = 0.0
$igpuDec   = 0.0
$igpuEnc   = 0.0
$igpuProc  = 0.0

if ($null -ne $gpuCounters) {
    foreach ($c in $gpuCounters.CounterSamples) {
        $val = [math]::Round($c.CookedValue, 4)
        $gpuTotal += $val
        if ($c.InstanceName -match 'engtype_3D')              { $igpu3D   += $val }
        if ($c.InstanceName -match 'engtype_VideoDecode')       { $igpuDec  += $val }
        if ($c.InstanceName -match 'engtype_VideoEncode')       { $igpuEnc  += $val }
        if ($c.InstanceName -match 'engtype_VideoProcessing')   { $igpuProc += $val }
    }
    $igpuUsage = [math]::Round($igpu3D + $igpuDec + $igpuEnc + $igpuProc, 4)
}

# GPU IDs (placeholder — real IDs require NVML / ADL)
$iGPUID = 15
$gGPUID = 15
$vpuID  = 0
$npuUsage = 0.0

# ── Memory ────────────────────────────────────────────────────────────────
$mem = Get-WmiObject Win32_OperatingSystem -ErrorAction SilentlyContinue
if ($null -ne $mem) {
    $memTotal = [math]::Round($mem.TotalVisibleMemorySize / 1MB, 0)
    $memFree  = [math]::Round($mem.FreePhysicalMemory / 1MB, 0)
    $memUsage = [math]::Round(($memTotal - $memFree) / $memTotal * 100, 1)
} else {
    $memTotal = 0
    $memFree  = 0
    $memUsage = 0
}

# ── Disk I/O ─────────────────────────────────────────────────────────────
$diskCounters = Get-Counter `
    '\PhysicalDisk(_Total)\% Disk Time',`
    '\PhysicalDisk(_Total)\Disk Bytes/sec',`
    '\PhysicalDisk(_Total)\Avg. Disk sec/Read',`
    '\PhysicalDisk(_Total)\Avg. Disk sec/Write' `
    -ErrorAction SilentlyContinue

$diskUsage    = 0.0
$diskSpeed    = 0.0
$diskReadLat  = 0.0
$diskWriteLat = 0.0

if ($null -ne $diskCounters) {
    foreach ($s in $diskCounters.CounterSamples) {
        if ($s.Path -match '% Disk Time')              { $diskUsage    = [math]::Round($s.CookedValue, 4) }
        if ($s.Path -match 'Disk Bytes')               { $diskSpeed    = [math]::Round($s.CookedValue, 0) }
        if ($s.Path -match 'sec/Read')                 { $diskReadLat  = [math]::Round($s.CookedValue * 1000, 5) }
        if ($s.Path -match 'sec/Write')                { $diskWriteLat = [math]::Round($s.CookedValue * 1000, 5) }
    }
}

# ── Battery Capacity ──────────────────────────────────────────────────────
$battCap = 0
if ($null -ne $batt) {
    $battCap = $batt.EstimatedChargeRemaining
}

# ── Foreground App ───────────────────────────────────────────────────────
$fgApp = ''
try {
    $fgProc = Get-Process |
              Where-Object { $_.MainWindowTitle -ne '' } |
              Select-Object -First 1 -ExpandProperty ProcessName
    if ($fgProc) { $fgApp = $fgProc }
} catch { }

# ── ITS Mode (from Dispatcher registry) ─────────────────────────────────
$itsMode = 0
try {
    $regPath = 'HKLM:\SYSTEM\CurrentControlSet\Services\LenovoProcessManagement\Performance\PowerSlider'
    $itsVal = Get-ItemProperty -Path $regPath -Name 'ITS_AutomaticModeSetting' -ErrorAction SilentlyContinue
    if ($null -ne $itsVal) {
        $itsMode = [int]$itsVal.ITS_AutomaticModeSetting
    }
} catch { }

# ── Memory Speed ──────────────────────────────────────────────────────────
$memSpeed = 0
$memInfo = Get-WmiObject Win32_PhysicalMemory -ErrorAction SilentlyContinue | Select-Object -First 1
if ($null -ne $memInfo) {
    $memSpeed = [int]$memInfo.Speed
}

# ── GPU VRAM (estimated from system memory) ───────────────────────────────
$igpuVram  = [math]::Round($memTotal * 0.25, 0)
$igpuShare = [int]$memFree
$igpuTotal = [int]$memTotal

# ── Output 53 pipe-delimited fields ────────────────────────────────────
# Empty slots for HW counters not available via WMI (latency, GDI, copy, etc.)
$output = @(
    $cpuUsage, $cpuFreq, $cpuPerf,
    $gpuTotal, $igpuUsage, $dgpuUsage, $vpuUsage, $npuUsage,
    $iGPUID, $gGPUID, $vpuID,
    $memUsage, $memFree,
    0, 0, 0, 0,           # InputLatency, Lag_64ms, Lag_100ms, Lag_200ms
    $diskUsage, $diskSpeed, $diskReadLat, $diskWriteLat,
    0, 0, 0, 0, 0,       # SystemPower, CPUPower, GPU0Power, NvidiaPower, NvidiaTemp (injected by Go from IPF DLL)
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0,  # Copy, GDI_Render, Legacy_Overlay, Security, 3D, Video_Dec, Video_Enc, Video_Proc, Unknown, Compute
    $itsMode, $battCap, $fgApp,
    $memSpeed, $igpuVram, $igpuShare, $igpuTotal
) -join '|'

Write-Output $output
