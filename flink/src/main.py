# flink/src/main.py
from pyflink.datastream import StreamExecutionEnvironment


from pyflink.datastream import StreamExecutionEnvironment
from pyflink.datastream.checkpointing_mode import CheckpointingMode
import config

def main():
    # 初始化 Flink 环境
    env = StreamExecutionEnvironment.get_execution_environment()

    # 配置检查点
    env.get_checkpoint_config().set_checkpointing_mode(CheckpointingMode.EXACTLY_ONCE)
    env.get_checkpoint_config().set_checkpoint_interval(10000)

    # 运行 URL 管理作业
    from url_management import run_url_management_job
    run_url_management_job(env, config)

    # 运行 URL 分发作业
    from url_dispatch import run_url_dispatch_job
    run_url_dispatch_job(env, config)

    # 执行作业
    env.execute("Distributed Crawler Flink Job")

if __name__ == "__main__":
    main()