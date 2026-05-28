from __future__ import annotations

import argparse
import asyncio
import json
from collections import defaultdict, deque
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone

from aiokafka import AIOKafkaConsumer


@dataclass
class SlidingWindow:
    width: timedelta

    def __post_init__(self) -> None:
        self.values: dict[str, deque[tuple[datetime, float]]] = defaultdict(deque)

    def add(self, metric: str, value: float, at: datetime) -> dict[str, float]:
        bucket = self.values[metric]
        bucket.append((at, value))
        cutoff = at - self.width
        while bucket and bucket[0][0] < cutoff:
            bucket.popleft()
        return {metric: sum(v for _, v in bucket) / len(bucket)}


async def consume(brokers: str, topic: str) -> None:
    consumer = AIOKafkaConsumer(
        topic,
        bootstrap_servers=brokers,
        group_id="laba14-python-analyzer",
        auto_offset_reset="latest",
    )
    window = SlidingWindow(timedelta(minutes=5))
    await consumer.start()
    try:
        async for message in consumer:
            aggregate = json.loads(message.value)
            metric = aggregate["metric"]
            avg = float(aggregate["average"])
            at = datetime.now(timezone.utc)
            print(window.add(metric, avg, at))
    finally:
        await consumer.stop()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--brokers", default="localhost:9092")
    parser.add_argument("--topic", default="health-aggregates")
    args = parser.parse_args()
    asyncio.run(consume(args.brokers, args.topic))


if __name__ == "__main__":
    main()
