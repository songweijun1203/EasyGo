# 1. 自动获取项目根目录
$RootDir = Split-Path -Parent $PSScriptRoot

# 2. 拼接绝对路径
$OutputDir = Join-Path $RootDir "build\bin"
$OutputFile = Join-Path $OutputDir "myapp_debug.exe"
$SourcePath = Join-Path $RootDir "cmd\myapp"

# 3. 确保输出目录存在
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
}

Write-Host "Building: $SourcePath"

# 4. 执行编译
go build -o $OutputFile -gcflags="all=-N -l" $SourcePath

# 5. 检查编译结果
if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Build completed: $OutputFile" -ForegroundColor Green
} else {
    Write-Host "[FAILED] Build failed! Please check errors above." -ForegroundColor Red
}