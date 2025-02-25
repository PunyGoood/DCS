# Define the root directory and output directory
$protoRoot = "D:\DT\grpc\proto"
$outputDir = "$protoRoot\output"

# Ensure the output directory exists
if (-Not (Test-Path -Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# Find all .proto files in the root directory and its subdirectories
$protoFiles = Get-ChildItem -Path $protoRoot -Filter *.proto -Recurse

if ($protoFiles -eq $null) {
    Write-Host "No .proto files found in the specified directory."
} else {
    foreach ($protoFile in $protoFiles) {
        $protoFilePath = Resolve-Path $protoFile.FullName
        Write-Host "Compiling $protoFilePath..."
        python -m grpc_tools.protoc --proto_path="$protoRoot" --python_out="$outputDir" --grpc_python_out="$outputDir" "$protoFilePath"
    }
}

Write-Host "All .proto files have been compiled to Python."