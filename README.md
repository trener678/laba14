# Лабораторная работа №14

Лабораторная работа: 14  
ФИО: Ваняев Дмитрий Александрович  
Группа: 221331  
Вариант: 22  
Сложность: повышенная

## Вариант

Предметная область: сбор и анализ медицинских данных.  
Источник данных: датчики здоровья, эмуляция.  
Тип заданий: повышенный.

## Что реализовано

- Распределенный Go-сборщик: несколько экземпляров регистрируются в etcd и детерминированно делят источники данных.
- Оконная агрегация tumbling window в Go: сборщик отправляет не сырые записи, а агрегаты `count/min/max/sum/average`.
- Передача через Apache Arrow: HTTP endpoint `/arrow/aggregates` отдает Arrow IPC stream с RecordBatch.
- Rust-библиотека валидации: `rust/health_validator` проверяет поля и диапазоны, есть C ABI для cgo-сборки.
- Kubernetes-развертывание: Dockerfile, Deployment, Service, etcd и HPA по CPU.
- Потоковая обработка: Go-сборщик публикует агрегаты в Kafka через `segmentio/kafka-go`, Python-анализатор читает тот же топик и поддерживает скользящее окно 5 минут.
- Python-сборщик для сравнения производительности: `python/collector/async_collector.py`.
- Веб-дашборд Streamlit: `python/dashboard/app.py`, обновляет графики по Arrow endpoint.

## Структура

```text
cmd/collector              Go collector entrypoint
internal/health            доменная модель и эмулятор датчиков
internal/sharding          распределение источников и etcd-координация
internal/window            tumbling window aggregation
internal/arrowipc          Apache Arrow IPC encoder and HTTP server
internal/kafkastream       Kafka producer
internal/validator         Go fallback + cgo hook for Rust validator
rust/health_validator      Rust validation library
python/analyzer            Arrow and Kafka clients
python/collector           asyncio collector benchmark
python/dashboard           Streamlit dashboard
k8s                        Kubernetes manifests
docs/performance_report.md отчет по сравнению Go/Python
PROMPT_LOG.md              журнал запросов к ИИ
```

## Быстрый запуск

```powershell
go mod tidy
go test ./...
go run ./cmd/collector -collector-id collector-a -http :8080 -window 10s
```

После запуска Arrow stream доступен по адресу:

```text
http://localhost:8080/arrow/aggregates
```

Python-клиент:

```powershell
cd python
python -m venv .venv
.\.venv\Scripts\activate
pip install -r requirements.txt
python -m analyzer.arrow_client --url http://localhost:8080/arrow/aggregates
streamlit run dashboard/app.py
```

## Распределенный режим

```powershell
docker compose up --build
```

В compose поднимаются etcd, Kafka и два экземпляра Go-сборщика. Каждый collector регистрируется в etcd, получает свою часть источников и публикует оконные агрегаты в Kafka topic `health-aggregates`.

Проверка сквозной цепочки Go -> Kafka -> Python:

```powershell
docker compose up --build
cd python
python -m analyzer.kafka_stream --brokers localhost:29092 --topic health-aggregates
```

## Kubernetes

```powershell
minikube start
minikube image build -t laba14-health-collector:latest .
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/etcd.yaml
kubectl apply -f k8s/collector-deployment.yaml
kubectl apply -f k8s/hpa.yaml
kubectl -n laba14 get pods,hpa
```

Kafka для Kubernetes рекомендуется ставить Helm-чартом Bitnami, подсказка оставлена в `k8s/kafka-placeholder.yaml`.

## Rust-валидатор

```powershell
cd rust/health_validator
cargo test
cargo build --release
```

Для cgo-варианта Go-сборки после сборки Rust-библиотеки:

```powershell
go test -tags rustffi ./internal/validator
go build -tags rustffi ./cmd/collector
```

## Тесты

```powershell
go test ./...
cargo test --manifest-path rust/health_validator/Cargo.toml
```

Python-тестовый запуск требует установленный Python 3.12+ и зависимости из `python/requirements.txt`.
