$protoFiles = @(
    "D:/DT/grpc/proto/ApplicationLayer/message_service.proto",
    "D:/DT/grpc/proto/ApplicationLayer/task_service.proto",
    "D:/DT/grpc/proto/DataModels/node.proto",
    "D:/DT/grpc/proto/DataModels/task.proto",
    "D:/DT/grpc/proto/DomainLayer/node_.proto",
    "D:/DT/grpc/proto/DomainLayer/request.proto",
    "D:/DT/grpc/proto/DomainLayer/response.proto",
    "D:/DT/grpc/proto/DomainLayer/response_enum.proto",
    "D:/DT/grpc/proto/DomainLayer/stream_message.proto",
    "D:/DT/grpc/proto/DomainLayer/stream_message_enum.proto"
)

$outputPath = "./proto"


if (-not (Test-Path -Path $outputPath)) {
    New-Item -Path $outputPath -ItemType Directory | Out-Null
}


& protoc -I D:/DT/grpc/proto --go_out=$outputPath --go-grpc_out=$outputPath $protoFiles