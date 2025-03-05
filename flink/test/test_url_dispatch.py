# flink/test/test_url_dispatch.py
import unittest
from ..src.url_dispatch import run_url_dispatch_job
from ..src import config
from pyflink.datastream import StreamExecutionEnvironment
from ..src.strategies.round_robin_strategy import RoundRobinStrategy
class TestURLDispatch(unittest.TestCase):
    def test_run_url_dispatch_job(self):
        # 创建一个模拟的 StreamExecutionEnvironment
        env = StreamExecutionEnvironment.get_execution_environment()

        # 模拟输入数据（待爬取 URL）
        input_data = ["http://example.com/page1", "http://example.com/page2"]
        input_stream = env.from_collection(input_data)

        # 设置分发策略
        strategy = RoundRobinStrategy(node_count=2)

        # 调用 run_url_dispatch_job 函数
        run_url_dispatch_job(env, config)

        # 执行作业并获取结果
        result = []
        output_stream = env.add_sink(result)
        env.execute("Test URL Dispatch Job")

        # 验证输出结果
        expected_output = ["node-0:http://example.com/page1", "node-1:http://example.com/page2"]
        self.assertEqual(sorted(result), sorted(expected_output))

if __name__ == '__main__':
    unittest.main()