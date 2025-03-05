# round_robin_strategy.py

from .base_strategy import BaseStrategy


class RoundRobinStrategy(BaseStrategy):
    """
    轮询策略：按顺序循环分配任务到工作节点。
    """

    def distribute(self, tasks, workers):
        """
        实现轮询分发逻辑。

        :param tasks: 待分发的任务列表
        :param workers: 可用的工作节点列表
        :return: 任务与工作节点的映射（字典）
        """
        distribution = {}
        num_workers = len(workers)

        for i, task in enumerate(tasks):
            worker = workers[i % num_workers]  # 轮询选择工作节点
            if worker not in distribution:
                distribution[worker] = []
            distribution[worker].append(task)

        return distribution