from __future__ import annotations

import argparse
import asyncio
import random
import statistics
import time
from dataclasses import dataclass

import psutil


@dataclass(frozen=True)
class Source:
    sensor_id: str
    patient_id: str
    metric: str
    unit: str
    base_value: float


SOURCES = [
    Source("hr-001", "p-1001", "heart_rate", "bpm", 72),
    Source("hr-002", "p-1002", "heart_rate", "bpm", 81),
    Source("spo2-001", "p-1001", "spo2", "%", 97),
    Source("spo2-002", "p-1003", "spo2", "%", 95),
    Source("temp-001", "p-1002", "temperature", "C", 36.7),
    Source("temp-002", "p-1003", "temperature", "C", 37.1),
]


def simulate(source: Source) -> dict[str, object]:
    noise = random.gauss(0, 1)
    scale = {"heart_rate": 5, "spo2": 1.2, "temperature": 0.25}[source.metric]
    return {
        "sensor_id": source.sensor_id,
        "patient_id": source.patient_id,
        "metric": source.metric,
        "unit": source.unit,
        "value": round(source.base_value + noise * scale, 2),
        "collected_at": time.time(),
    }


async def collect(duration: float, rate: int) -> dict[str, float]:
    import aiohttp

    process = psutil.Process()
    latencies: list[float] = []
    deadline = time.monotonic() + duration
    async with aiohttp.ClientSession() as session:
        while time.monotonic() < deadline:
            start = time.perf_counter()
            batch = [simulate(random.choice(SOURCES)) for _ in range(rate)]
            async with session.post("http://127.0.0.1:9/blackhole", json=batch) as _:
                pass
            latencies.append(time.perf_counter() - start)
            await asyncio.sleep(1)
    return {
        "records_per_second": rate,
        "avg_batch_latency_ms": statistics.mean(latencies) * 1000 if latencies else 0,
        "rss_mb": process.memory_info().rss / 1024 / 1024,
        "cpu_percent": process.cpu_percent(),
    }


async def dry_collect(duration: float, rate: int) -> dict[str, float]:
    process = psutil.Process()
    process.cpu_percent(interval=None)
    total = 0
    deadline = time.monotonic() + duration
    while time.monotonic() < deadline:
        limit = rate if rate > 0 else 10000
        for i in range(limit):
            simulate(random.choice(SOURCES))
            total += 1
        if rate > 0:
            await asyncio.sleep(1)
        else:
            await asyncio.sleep(0)
    return {
        "records_per_second": total / duration,
        "rss_mb": process.memory_info().rss / 1024 / 1024,
        "cpu_percent": process.cpu_percent(interval=None),
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--duration", type=float, default=30)
    parser.add_argument("--rate", type=int, default=1000)
    args = parser.parse_args()
    print(asyncio.run(dry_collect(args.duration, args.rate)))


if __name__ == "__main__":
    main()
