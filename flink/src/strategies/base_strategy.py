# base_strategy.py

from abc import ABC, abstractmethod


class BaseStrategy(ABC):
    """
    基础策略类，所有分发策略都应继承此类。
    """

    @abstractmethod
    def distribute(self, tasks, workers):
        """
        分发任务的核心方法。

        :param tasks: 待分发的任务列表
        :param workers: 可用的工作节点列表
        :return: 分发结果，通常是任务与工作节点的映射
        """
        pass