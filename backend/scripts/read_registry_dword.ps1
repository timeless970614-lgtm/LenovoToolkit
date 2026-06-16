<# read_registry_dword.ps1
# Reads a DWORD value from a registry path and outputs it as a plain integer.
# Exit code: 0=success, 1=key not found, 2=invalid parameters
#
# Usage:
#   powershell -ExecutionPolicy Bypass -File read_registry_dword.ps1 "<path>" "<name>"
#
# Examples:
#   powershell -ExecutionPolicy Bypass -File read_registry_dword.ps1 "HKLM:\SOFTWARE\..." "EPP"
#   powershell -ExecutionPolicy Bypass -File read_registry_dword.ps1 "HKCU:\SOFTWARE\..." "ValueName"
param(
    [Parameter(Mandatory=$true)][string]$Path,
    [Parameter(Mandatory=$true)][string]$Name
)

$ErrorActionPreference = 'SilentlyContinue'

# Normalize path: convert HKLM:\ to HKLM: etc.
$normalizedPath = $Path -replace '^HKLM:', 'HKLM:' -replace '^HKCU:', 'HKCU:'

try {
    $value = Get-ItemProperty -Path $normalizedPath -Name $Name -ErrorAction SilentlyContinue
    if ($null -ne $value) {
        $val = $value.$Name
        if ($null -ne $val) {
            Write-Output ([int64]$val)
            exit 0
        }
    }
} catch {
    # Fall through to not found
}

# Not found: output 0 and exit 1 so Go side can distinguish from real 0
Write-Output "0"
exit 1
