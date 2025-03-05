# flink/src/strategies/strategy_source.py
from pyflink.datastream import SourceFunction

class StrategySource(SourceFunction):
    def __init__(self, strategy):
        self.strategy = strategy
        self.is_running = True

    def run(self, ctx):
        while self.is_running:
            ctx.collect(self.strategy)  # 将策略推送到广播流中

    def cancel(self):
        self.is_running = False