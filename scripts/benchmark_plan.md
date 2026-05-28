# Performance comparison plan

1. Run `go run ./scripts/go_collect_benchmark.go -duration 10s -rate 0`.
2. Capture Go process CPU seconds with PowerShell and convert them to one-core CPU percent.
3. Run `uv run --with psutil python ./python/collector/async_collector.py --duration 10 --rate 0`.
4. Store records per second, RSS memory and CPU percent.
5. Build charts in `docs/performance_report.md`.
