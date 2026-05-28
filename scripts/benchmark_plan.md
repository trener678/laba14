# Performance comparison plan

1. Start the Go collector for 30 seconds at the same source count and interval.
2. Capture throughput from logs and process metrics with `docker stats`.
3. Run `python/collector/async_collector.py --duration 30 --rate 1000`.
4. Store records per second, RSS memory and CPU percent.
5. Build charts in `docs/performance_report.md`.
