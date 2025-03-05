 # flink/src/url_dispatch.py

from pyflink.datastream.connectors.kafka import KafkaSource, KafkaSink
from pyflink.common.serialization import SimpleStringSchema
from pyflink.datastream.connectors import DeliveryGuarantee
from pyflink.datastream.state import MapStateDescriptor
from pyflink.common.typeinfo import Types
import config
from strategies.round_robin_strategy import RoundRobinStrategy
from strategies.strategy_source import StrategySource
from pyflink.datastream.functions import BroadcastProcessFunction

class URLDispatchFunction(BroadcastProcessFunction):
    def __init__(self, strategy_descriptor):
        self.strategy_descriptor = strategy_descriptor

    def process_broadcast_element(self, strategy, ctx, out):
        broadcast_state = ctx.get_broadcast_state(self.strategy_descriptor)
        broadcast_state.put("current_strategy", strategy)

    def process_element(self, url, ctx, out):
        broadcast_state = ctx.get_broadcast_state(self.strategy_descriptor)
        current_strategy = broadcast_state.get("current_strategy")
        if current_strategy:
            dispatched_url = current_strategy.dispatch(url)
            out.collect(dispatched_url)

def run_url_dispatch_job(env, config):
    # 1. 配置 KafkaSource，读取待爬取的 URL
    kafka_source = KafkaSource.builder() \
        .set_bootstrap_servers(config.KAFKA_BOOTSTRAP_SERVERS) \
        .set_topics(config.KAFKA_PENDING_TOPIC) \
        .set_start_from_earliest() \
        .set_value_only_deserializer(SimpleStringSchema()) \
        .build()

    # 2. 从 Kafka 读取 URL 并预处理
    pending_urls = env.from_source(
        kafka_source,
        watermark_strategy=None,
        source_name="pending-urls-source"
    ).map(lambda x: x.strip())

    # 3. 定义广播状态描述符（键为字符串，值为序列化的 Python 对象）
    strategy_descriptor = MapStateDescriptor(
        "strategy",
        Types.STRING(),
        Types.PICKLED_BYTE_ARRAY()
    )

    # 4. 创建并广播分发策略（例如 RoundRobinStrategy）
    round_robin_strategy = RoundRobinStrategy(node_count=config.NUM_WORKERS)
    strategy_source = StrategySource(round_robin_strategy)
    broadcast_strategy = env.add_source(strategy_source).broadcast(strategy_descriptor)

    # 5. 分发逻辑：连接主流和广播流
    dispatched_urls = pending_urls \
        .connect(broadcast_strategy) \
        .process(URLDispatchFunction(strategy_descriptor))

    # 6. 配置 KafkaSink，将分发后的 URL 写入 Kafka
    kafka_sink = KafkaSink.builder() \
        .set_bootstrap_servers(config.KAFKA_BOOTSTRAP_SERVERS) \
        .set_topic(config.KAFKA_DISPATCHED_TOPIC) \
        .set_value_serializer(SimpleStringSchema()) \
        .set_deliver_guarantee(DeliveryGuarantee.AT_LEAST_ONCE) \
        .build()

    # 7. 将结果写入 Kafka
    dispatched_urls.sink_to(kafka_sink)

    # 8. 执行 Flink 作业
    env.execute("URL Dispatch Job")