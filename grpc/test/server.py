import grpc
from concurrent import futures
import logging
import sys
import os
import time
# 将项目根目录添加到 sys.path
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
print("Project root added to sys.path:", project_root)
sys.path.append(project_root)

# from proto.ApplicationLayer import message_service_pb2_grpc, task_service_pb2_grpc, task_service_pb2
# from proto.DomainLayer import stream_message_pb2, request_pb2, response_pb2

from proto.ApplicationLayer import (
   message_service_pb2,
   message_service_pb2_grpc,
   task_service_pb2,
   task_service_pb2_grpc
)
from proto. DataModels import (
   task_pb2,
   task_pb2_grpc,
   node_pb2,
   node_pb2_grpc
)
from proto.DomainLayer import (
   stream_message_pb2,
   stream_message_pb2_grpc,
   request_pb2,
   request_pb2_grpc,
   response_pb2,
   response_pb2_grpc,
   response_enum_pb2,
   response_enum_pb2_grpc,

)


class MessageService(message_service_pb2_grpc.MessageServiceServicer):
   def Contact(self, request_iterator, context):
       for request in request_iterator:
           print(f"Received message: {request.data}")
           response = stream_message_pb2.StreamMessage(
               data=request.data
           )
           yield response

class TaskService(task_service_pb2_grpc.TaskServiceServicer):
    def Get(self, request, context):
        print(f"Received Get request: {request}")
        return response_pb2.Response(
            enum=response_enum_pb2.ResponseEnum.GOOD,
            message="Get response",
            data=b"",
            error=""
        )

    def Subscribe(self, request_iterator, context):
        try:
            # 遍历流式请求
            for request in request_iterator:
                print(f"Received stream request: node_key={request.node_key}, data={request.data}")

            # 处理完所有流式请求后返回响应
            return response_pb2.Response(
                enum=response_enum_pb2.ResponseEnum.GOOD,
                message="All stream requests processed.",
                data=b"qwer",
                error=""
            )
        except Exception as e:
            # 捕获异常并设置 gRPC 错误状态
            print(f"Error processing stream requests: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Internal error: {str(e)}")
            return response_pb2.Response(
                enum=response_enum_pb2.ResponseEnum.BAD,
                message="Error processing stream requests",
                data=b"",
                error=str(e)
            )

    def process_requests(self, requests):
        # 处理请求的逻辑
        # 这里可以根据具体需求实现处理请求的逻辑
        processed_data = b"Processed data"  # 示例处理结果
        return processed_data
    
def serve():
   server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
   message_service_pb2_grpc.add_MessageServiceServicer_to_server(
       MessageService(), server
   )
   task_service_pb2_grpc.add_TaskServiceServicer_to_server(
       TaskService(), server
   )
   server.add_insecure_port('[::]:50051')
   server.start()
   print("Server started on port 50051")
   try:
       while True:
           time.sleep(86400)
   except KeyboardInterrupt:
       server.stop(0)

if __name__ == '__main__':
   serve()

