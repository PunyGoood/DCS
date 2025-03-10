# flink/test/test_url_management.py
import unittest
from ..src.url_management import run_url_management_job
from ..src import config
from pyflink.datastream import StreamExecutionEnvironment


class TestURLManagement(unittest.TestCase):
    def test_run_url_management_job(self):
        # 创建一个模拟的 StreamExecutionEnvironment
        env = StreamExecutionEnvironment.get_execution_environment()

        # 模拟输入数据（种子 URL）
        input_data = ["http://example.com", "http://example.org", "http://example.com"]
        input_stream = env.from_collection(input_data)

        # 调用 run_url_management_job 函数
        run_url_management_job(env, config)

        # 执行作业并获取结果
        result = []
        output_stream = env.add_sink(result)
        env.execute("Test URL Management Job")

        # 验证输出结果
        expected_output = ["http://example.com", "http://example.org"]  # 去重后的 URL 列表
        self.assertEqual(sorted(result), sorted(expected_output))


if __name__ == '__main__':
    unittest.main()