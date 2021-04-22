#!/usr/bin/env pwsh

function Log-Info($message) {
	Write-Output "`e[38;5;61m  ==>`e[0;00m $message"
}
function Log-Critical($message) {
	[Console]::Error.WriteLine("")
	[Console]::Error.WriteLine("  `e[38;5;125m$message`e[0;00m")
	[Console]::Error.WriteLine("")
}
function Get-OS() {
	if ($IsWindows) {
		return "windows"
	} else {
		Log-Critical("Operating system not supported yet")
		Exit 0
	}
	return "unknown"
}
function Get-Architecture() {
	return "$env:PROCESSOR_ARCHITECTURE".ToLower()
}

$OS = Get-OS
$Arch = Get-Architecture

# API endpoint such as "http://localhost:3000"
$API="{{.URL}}"

# package such as "github.com/tj/triage/cmd/triage"
$Pkg="{{.Package}}"

# binary name such as "hello"
$Bin="{{.Binary}}"

# original_version such as "master"
$OriginalVersion="{{.OriginalVersion}}"

# version such as "master"
$Version="{{.Version}}"

$GoBinariesDir = $env:GOBINARIES_DIR
$BinDir = If ($GoBinariesDir) {
  "$GoBinariesDir\bin"
} else {
  "$HOME\.gobinaries\bin"
}

$TempExe = "$env:TEMP\$Bin.exe"
$TargetExe = "$BinDir\$Bin.exe"	

If ($OriginalVersion -ne $Version) {
	Log-Info("Resolved version $OriginalVersion to $Version")
}

# GitHub requires TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

Log-Info("Downloading binary for $OS $arch")
$DownloadUrl = "${API}/binary/${Pkg}?os=${OS}&arch=${Arch}&version=${Version}"
Invoke-WebRequest $DownloadUrl -OutFile $TempExe -UseBasicParsing

If (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

Copy-Item -Path $TempExe -Destination $TargetExe

If (Test-Path $TempExe) {
	Remove-Item $TempExe
}

$User = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable('Path', $User)
If (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
  [Environment]::SetEnvironmentVariable('Path', "$Path;$BinDir", $User)
  $Env:Path += ";$BinDir"
}

Log-Info("Installation complete")