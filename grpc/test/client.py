import grpc

import sys
import os

project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
sys.path.append(project_root)

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
   stream_message_enum_pb2,
   stream_message_enum_pb2_grpc,
   request_pb2,
   request_pb2_grpc,
   response_pb2,
   response_pb2_grpc,
   response_enum_pb2,
   response_enum_pb2_grpc,

)



class GRPCClient:
    def __init__(self, host='localhost', port=50051):
        self.channel = grpc.insecure_channel(f'{host}:{port}')

    def test_task_service(self):
        task_stub = task_service_pb2_grpc.TaskServiceStub(self.channel)
        print("Testing TaskService.Subscribe...")

        messages = [
            stream_message_pb2.StreamMessage(
                    enum=stream_message_enum_pb2.StreamMessageEnum.RUN_TASK, 
                    node_key="node1",
                    key="key1",
                    source_id="source1",
                    target_id="target1",
                    data=b"Subscribe 1",
                    error=""
                ),
            stream_message_pb2.StreamMessage(
                    enum=stream_message_enum_pb2.StreamMessageEnum.CANCEL_TASK,  
                    node_key="node2",
                    key="key2",
                    source_id="source2",
                    target_id="target2",
                    data=b"Subscribe 2",
                    error=""
                )
        ]
        try:
            response = task_stub.Subscribe(iter(messages))
            print(f"TaskService.Subscribe response: message={response.message}")
        except grpc.RpcError as e:
            print(f"RPC error: {e.code()} - {e.details()}")


    def test_message_service(self):
        message_stub = message_service_pb2_grpc.MessageServiceStub(self.channel)
        print("Testing MessageService.Contact...")
        
        messages = [
            stream_message_pb2.StreamMessage(data=b"Hello"),
            stream_message_pb2.StreamMessage(data=b"World")
        ]
        
        try:
            responses = message_stub.Contact(iter(messages))
            for response in responses:
                print(f"MessageService response data: {response.data}")
        except grpc.RpcError as e:
            print(f"RPC error: {e.code()} - {e.details()}")

    def test_get_service(self):
        task_stub = task_service_pb2_grpc.TaskServiceStub(self.channel)
        print("Testing TaskService.Subscribe...")
        
        message = request_pb2.Request(
        node="node1",  
        message="Test message",  
        data=b"asdasd"  
    )
        try:
            response = task_stub.Get(message)
            print(f"MessageService response data: {response.enum}")
        except grpc.RpcError as e:
            print(f"RPC error: {e.code()} - {e.details()}")

def run_tests():
   client = GRPCClient()
   print("Running MessageService tests...")
   client.test_message_service()
   
   print("\nRunning Get tests...")
   client.test_get_service()

   print("\nRunning TaskService tests...")
   client.test_task_service()

   
if __name__ == '__main__':
   run_tests()