# btcwallet_create.ps1

# 로그 파일 경로 설정 (디버깅 용도)
$logPath = "C:\dev\workspace\CSE-416\btcwallet\btcwallet_create.log"

# 로그 기록 함수
function Log-Message {
    param([string]$message)
    Add-Content -Path $logPath -Value "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') - $message"
}

Log-Message "Starting btcwallet creation process."

# 현재 실행 정책 저장
$currentExecutionPolicy = Get-ExecutionPolicy -Scope CurrentUser
Log-Message "Current Execution Policy: $currentExecutionPolicy"

# 실행 정책 변경 (자동으로 'Yes' 응답)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
Log-Message "Execution Policy set to RemoteSigned."

# 필요한 어셈블리 로드
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

# btcwallet 실행 경로
$btcwalletPath = "C:\dev\workspace\CSE-416\btcwallet\btcwallet.exe"

# btcwallet 프로세스 시작
$process = Start-Process -FilePath $btcwalletPath -ArgumentList "--create" -WindowStyle Normal -PassThru

Log-Message "btcwallet process started."

# btcwallet 창이 열릴 시간을 기다림
Start-Sleep -Seconds 3

# btcwallet 창 핸들 찾기 (창 제목을 실제로 나타나는 제목으로 변경하세요.)
$hwnd = [Win32]::FindWindow($null, "Btcwallet") 

if ($hwnd -ne 0) {
    # 창을 포그라운드로 설정
    [Win32]::SetForegroundWindow($hwnd)
    Log-Message "btcwallet window found and activated."
} else {
    Log-Message "btcwallet window not found."
}

# 환경 변수에서 비밀번호 가져오기 (보안 강화)
$passphrase = $env:BTCWALLET_PASSPHRASE
$confirmPassphrase = $env:BTCWALLET_PASSPHRASE
$addEncryption = "no"
$existingSeed = "no"
$confirmOK = "OK"

# 입력 시뮬레이션
$inputs = @(
    "$passphrase{ENTER}"        # 비밀번호 입력 및 엔터
    "$confirmPassphrase{ENTER}" # 비밀번호 확인 및 엔터
    "$addEncryption{ENTER}"     # 추가 암호화 없음
    "$existingSeed{ENTER}"      # 기존 지갑 시드 없음
    "$confirmOK{ENTER}"         # 시드 저장 확인
)

foreach ($input in $inputs) {
    [System.Windows.Forms.SendKeys]::SendWait($input)
    Log-Message "Sent input: $input"
    Start-Sleep -Seconds 2 # 각 입력 사이에 더 긴 지연을 줌
}

# 프로세스가 종료될 때까지 대기
$process.WaitForExit()

Log-Message "btcwallet creation process completed."

# 실행 정책 복원
Set-ExecutionPolicy -ExecutionPolicy $currentExecutionPolicy -Scope CurrentUser -Force
Log-Message "Execution Policy restored to $currentExecutionPolicy."
