# flink/src/url_management.py
from pyflink.datastream import StreamExecutionEnvironment
from pyflink.datastream.functions import KeyedProcessFunction, RuntimeContext
from pyflink.common.typeinfo import Types
from pyflink.datastream.state import ValueStateDescriptor
import config
from utils.url_validator import is_valid_url

from pyflink.datastream import StreamExecutionEnvironment
from pyflink.datastream.connectors.kafka import KafkaSource, KafkaSink
from pyflink.common.serialization import SimpleStringSchema
from pyflink.datastream.connectors import DeliveryGuarantee
import config

class URLDeduplicationFunction(KeyedProcessFunction):
    def open(self, runtime_context: RuntimeContext):
        self.seen_urls = runtime_context.get_state(ValueStateDescriptor("seen_urls", Types.STRING()))

    def process_element(self, url, ctx, out):
        if not self.seen_urls.value():
            self.seen_urls.update(url)
            out.collect(url)

def run_url_management_job(env, config):
    # 配置 KafkaSource，从 Kafka 读取种子 URL
    kafka_source = KafkaSource.builder() \
        .set_bootstrap_servers(config.KAFKA_BOOTSTRAP_SERVERS) \
        .set_topics(config.KAFKA_SEED_TOPIC) \
        .set_start_from_earliest() \
        .set_value_only_deserializer(SimpleStringSchema()) \
        .build()

    # 从 Kafka 获取种子 URL 并进行处理
    seed_urls = env.from_source(
        kafka_source,
        watermark_strategy=None,
        source_name="seed-url-source"
    )

    # 示例：简单地将 URL 进行去重和过滤（这里只是一个示例）
    valid_urls = seed_urls.map(lambda x: x.strip()).filter(is_valid_url)

    # 配置 KafkaSink，将处理后的 URL 写入 Kafka
    kafka_sink = KafkaSink.builder() \
        .set_bootstrap_servers(config.KAFKA_BOOTSTRAP_SERVERS) \
        .set_topic(config.KAFKA_PENDING_TOPIC) \
        .set_value_serializer(SimpleStringSchema()) \
        .set_deliver_guarantee(DeliveryGuarantee.AT_LEAST_ONCE) \
        .build()

    # 将处理后的 URL 写入 Kafka
    valid_urls.sink_to(kafka_sink)