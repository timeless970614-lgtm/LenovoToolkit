<# get_acdc.ps1
# Returns AC/DC status: 1=AC (plugged in), 0=DC (on battery)
# Source: Win32_Battery.BatteryStatus
#   1=Discharging, 2=On AC (no battery), 3=Fully Charged, others=unknown
#>
$ErrorActionPreference = 'SilentlyContinue'
$b = Get-WmiObject Win32_Battery -ErrorAction SilentlyContinue
if ($null -eq $b) {
    # No battery at all -> treat as AC
    Write-Output "1"
} elseif ($b.BatteryStatus -eq 2 -or $b.BatteryStatus -eq 3) {
    Write-Output "1"
} else {
    Write-Output "0"
}
