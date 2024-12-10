# btcwallet_create.ps1

# Set the path for the log file (for debugging purposes)
$logPath = "C:\dev\workspace\CSE-416\btcwallet\btcwallet_create.log"

# Function to log messages
function Log-Message {
    param([string]$message)
    Add-Content -Path $logPath -Value "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') - $message"
}

Log-Message "Starting btcwallet creation process."

# Save the current execution policy
$currentExecutionPolicy = Get-ExecutionPolicy -Scope CurrentUser
Log-Message "Current Execution Policy: $currentExecutionPolicy"

# Change the execution policy (auto-confirm with 'Yes')
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
Log-Message "Execution Policy set to RemoteSigned."

# Load required assemblies
Add-Type @"
using System;
using System.Runtime.InteropServices;
public class Win32 {
    [DllImport("user32.dll", SetLastError=true)]
    public static extern IntPtr FindWindow(string lpClassName, string lpWindowName);
    
    [DllImport("user32.dll")]
    [return: MarshalAs(UnmanagedType.Bool)]
    public static extern bool SetForegroundWindow(IntPtr hWnd);
}
"@

Add-Type -AssemblyName System.Windows.Forms

# Path to the btcwallet executable
$btcwalletPath = "../btcwallet/btcwallet"

# Start the btcwallet process
$process = Start-Process -FilePath $btcwalletPath -ArgumentList "--create" -WindowStyle Normal -PassThru

Log-Message "btcwallet process started."

# Wait for the btcwallet window to open
Start-Sleep -Seconds 3

# Find the btcwallet window handle (replace with the actual window title)
$hwnd = [Win32]::FindWindow($null, "Btcwallet") 

if ($hwnd -ne 0) {
    # Set the window to the foreground
    [Win32]::SetForegroundWindow($hwnd)
    Log-Message "btcwallet window found and activated."
} else {
    Log-Message "btcwallet window not found."
}

# Retrieve the passphrase from environment variables (for security)
$passphrase = $env:BTCWALLET_PASSPHRASE
$confirmPassphrase = $env:BTCWALLET_PASSPHRASE
$addEncryption = "no"
$existingSeed = "no"
$confirmOK = "OK"

# Simulate inputs
$inputs = @(
    "$passphrase{ENTER}"        # Enter the passphrase and press Enter
    "$confirmPassphrase{ENTER}" # Confirm the passphrase and press Enter
    "$addEncryption{ENTER}"     # Select no additional encryption
    "$existingSeed{ENTER}"      # Select no existing wallet seed
    "$confirmOK{ENTER}"         # Confirm seed saving
)

foreach ($input in $inputs) {
    [System.Windows.Forms.SendKeys]::SendWait($input)
    Log-Message "Sent input: $input"
    Start-Sleep -Seconds 2 # Add a longer delay between inputs
}

# Wait for the process to exit
$process.WaitForExit()

Log-Message "btcwallet creation process completed."

# Restore the original execution policy
Set-ExecutionPolicy -ExecutionPolicy $currentExecutionPolicy -Scope CurrentUser -Force
Log-Message "Execution Policy restored to $currentExecutionPolicy."
