$scriptPath = $PSCommandPath
$scriptItem = Get-Item $scriptPath

if ($scriptItem.LinkType -eq "SymbolicLink") {
    $linkDirectory = Split-Path -Parent $scriptPath
    $targetRelative = $scriptItem.Target
    $targetPath = Join-Path $linkDirectory $targetRelative
    $resolvedTarget = Resolve-Path $targetPath -ErrorAction SilentlyContinue
    if ($resolvedTarget) {
        $BASE_DIR = Split-Path -Parent $resolvedTarget
    } else {
        $BASE_DIR = Split-Path -Parent $scriptPath
    }
} else {
    $BASE_DIR = Split-Path -Parent $scriptPath
}

$base_path = (Get-Item (Join-Path $BASE_DIR "..")).FullName

Set-Location $base_path

# 更改 `protoc` 命令路径
protoc -I ./proto --go_out=proto --go-grpc_out=proto ./proto/**/*.proto