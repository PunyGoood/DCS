# priority_strategy.py

from .base_strategy import BaseStrategy


class PriorityStrategy(BaseStrategy):
    """
    优先级策略：根据工作节点的优先级分配任务。
    """

    def distribute(self, tasks, workers):
        """
        实现优先级分发逻辑。

        :param tasks: 待分发的任务列表
        :param workers: 可用的工作节点列表（假设每个节点有一个优先级属性）
        :return: 任务与工作节点的映射（字典）
        """
        # 假设每个工作节点有一个 'priority' 属性，优先级越高，分配的任务越多
        workers.sort(key=lambda worker: worker.get('priority', 0), reverse=True)

        distribution = {worker['id']: [] for worker in workers}
        total_priority = sum(worker.get('priority', 1) for worker in workers) or 1

        for task in tasks:
            # 根据优先级比例分配任务
            for worker in workers:
                priority_ratio = worker.get('priority', 1) / total_priority
                if priority_ratio > 0:  # 简单模拟分配逻辑
                    distribution[worker['id']].append(task)
                    break

        return distribution