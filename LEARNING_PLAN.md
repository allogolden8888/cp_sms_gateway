# SMPP Testing Tool — Learning Plan

## Цель проекта
Комплексный распределённый инструмент тестирования SMS шлюзов и SMSC.
- Нагрузочное тестирование до 50 000 TPS (горизонтальное масштабирование)
- Автоматическое conformance-тестирование SMPP (PDU, TLV, поведение)
- TCP-прокси монитор для захвата и анализа сессий
- Встроенный SMSC-симулятор
- Обнаружение вендорских отклонений
- Сценарное тестирование
- Master-Worker распределённая архитектура (gRPC)
- Web UI, мультиплатформа

---

## Фаза 1: Фундамент ✅ (завершена)

### Go basics
- [x] Флаги командной строки (`flag` package)
- [x] Указатели и разыменование (`*`, `&`)
- [x] Форматирование строк (`fmt.Sprintf`, `fmt.Printf`, raw strings)
- [x] Логические операторы, `switch/case`
- [x] Интерфейсы — объявление и реализация
- [x] Множественный возврат `(value, error)`
- [x] `defer`
- [x] Обработка ошибок (`if err != nil`)
- [x] Горутины (`go func()`)
- [x] Каналы — небуферизованные, буферизованные, `select`, `time.After`
- [x] Пакеты и импорты
- [x] `len([]rune(text))` vs `len(text)`
- [x] Регулярные выражения
- [x] Структуры (`struct`)

### SMPP
- [x] Bind типы: TX, RX, TRX
- [x] Подключение через `go-smpp`
- [x] Отправка SMS (`Submit`, `SubmitLongMsg`)
- [x] Кодировки: GSM7, Latin1, UCS2
- [x] Поля `ShortMessage`: Register, PriorityFlag, Validity, TON/NPI
- [x] Validity period — относительный и абсолютный
- [x] DLR — получение через `Handler`
- [x] Парсинг DLR строки — структура `DLR`, `ParseDLR`
- [x] Multipart DLR — буферизованный канал + счётчик частей
- [x] Таймаут ожидания DLR

### Структура
- [x] Разбивка на пакеты: `smpp/client.go`, `smpp/message.go`, `smpp/encoding.go`, `smpp/dlr.go`
- [x] JSON конфиг (`encoding/json`, `os.ReadFile`)
- [x] `main.go` — только точка входа

---

## Фаза 2: Качество кода (сейчас)

### Go концепты
- [x] Тесты (`testing`) — табличные тесты, `t.Run`, хелперы
- [ ] Бенчмарки (`testing.B`) — `go test -bench`
- [x] Структурированное логирование (`log/slog`)
- [ ] Правильный `error` wrapping — `fmt.Errorf("%w")`, `errors.Is`, `errors.As`
- [ ] Graceful shutdown — `os/signal`, `signal.NotifyContext`

### Рефакторинг
- [x] Исправить баг: проверка `err` после `LoadConfig` в `main.go`
- [x] Вынести DLR-сопоставление (`respIDs` map) из `main.go` в пакет (`DLRTracker`)

---

## Фаза 3: Протокол (глубокое погружение)

### Binary протокол
- [x] `encoding/binary` — чтение и запись бинарных данных вручную
- [x] PDU заголовок (16 байт) — `ParsePDUHeader`, big-endian, uint32
- [x] C-строки в бинарных протоколах — `readCString`
- [x] Полный PDU парсер: body с полями разных типов (`parseSubmitSMBody`)
- [x] Рефакторинг: субпакет `smpp/pdu/` с группировкой по типам
- [x] TLV теги — парсинг (`parseTLV`, `parseTLVs`)
- [ ] Struct embedding — переиспользование полей PDU
- [ ] Публичный `ParsePDU` — точка входа, объединяет header + body по command_id

### SMPP сервер (режим замены SMSC)
- [ ] `net.Listen` / `net.Accept` — TCP сервер
- [ ] Обработка входящих SMPP-подключений
- [ ] Ответы на `bind`, `submit_sm`, генерация DLR
- [ ] `io.Reader` / `io.Writer` интерфейсы

### TCP-прокси монитор
- [ ] Прозрачный TCP-прокси (`net.Dial` + `net.Listen` + `io.Copy`)
- [ ] Захват raw bytes обеих сторон
- [ ] Декодирование PDU из захваченного трафика
- [ ] Сохранение сессий для анализа

---

## Фаза 4: Движок нагрузки

### Производительность
- [ ] Worker pool паттерн — горутины + канал задач
- [ ] `golang.org/x/time/rate` — token bucket, точный TPS
- [ ] `sync.WaitGroup` — ожидание завершения воркеров
- [ ] `sync.Mutex` — защита общих данных
- [ ] `sync/atomic` — lock-free счётчики метрик
- [ ] `sync.Pool` — переиспользование буферов, снижение GC
- [ ] `context.Context` — отмена, таймауты, propagation

### Метрики
- [ ] Real-time TPS счётчик
- [ ] Latency (min/max/p50/p95/p99)
- [ ] DLR rate, error rate
- [ ] Счётчики по статусам PDU

### Connection pool
- [ ] Пул SMPP-сессий (N параллельных соединений к одному шлюзу)
- [ ] Переподключение при разрыве

---

## Фаза 5: CLI и сценарии

### CLI
- [ ] `cobra` — субкоманды: `run`, `monitor`, `smsc`, `report`, `scenario`
- [ ] `//go:embed` — встраивание web UI и дефолтных сценариев в бинарник

### Сценарный движок
- [ ] YAML-сценарии (`gopkg.in/yaml.v3`)
- [ ] Предустановленные сценарии (burst, sustained, multipart flood)
- [ ] Параметризованные шаблоны сообщений

### Хранение результатов
- [ ] `database/sql` + SQLite (`modernc.org/sqlite`)
- [ ] Сохранение результатов тестов
- [ ] Сравнение запусков, история

---

## Фаза 6: Distributed (Master-Worker)

### gRPC
- [ ] Protocol Buffers — `.proto` файлы, `protoc`
- [ ] gRPC сервер и клиент на Go
- [ ] Unary RPC — раздача задач воркерам
- [ ] Server-side streaming — метрики в реальном времени от воркеров к мастеру
- [ ] Health checks (`grpc/health`)

### Координация
- [ ] Регистрация воркеров у мастера
- [ ] Heartbeat и детектирование падения узла
- [ ] Распределение нагрузки по воркерам
- [ ] Агрегация метрик со всех узлов

### Конфигурация кластера
- [ ] Конфигурация воркер-узлов
- [ ] Динамическое добавление воркеров в рантайме

---

## Фаза 7: Observability и Web UI

### Web UI
- [ ] `net/http` — встроенный HTTP сервер
- [ ] REST API для управления тестами
- [ ] Server-Sent Events или WebSocket для live метрик
- [ ] Базовый фронтенд (HTML + JS, без фреймворков)

### Observability
- [ ] OpenTelemetry — метрики и трейсинг
- [ ] Экспорт в Prometheus формат
- [ ] `pprof` — профилирование CPU и памяти
- [ ] `go test -race` — детектор гонок данных

---

## Фаза 8: Полировка и дистрибуция

### Generics (Go 1.18+)
- [ ] Рефакторинг worker pool в generic
- [ ] Generic ring buffer для метрик

### Дистрибуция
- [ ] Cross-compilation (`GOOS`, `GOARCH`)
- [ ] GoReleaser — автосборка под Windows/Linux/Mac
- [ ] Docker образ для воркер-узлов
- [ ] `Makefile` — стандартные команды сборки

---

## Библиотеки

| Библиотека | Назначение |
|---|---|
| `github.com/fiorix/go-smpp` | SMPP клиент |
| `github.com/spf13/cobra` | CLI субкоманды |
| `google.golang.org/grpc` | gRPC |
| `google.golang.org/protobuf` | Protocol Buffers |
| `golang.org/x/time/rate` | Rate limiter |
| `gopkg.in/yaml.v3` | YAML сценарии |
| `modernc.org/sqlite` | SQLite (pure Go, без CGO) |
| `go.opentelemetry.io/otel` | OpenTelemetry |

## Симулятор для разработки
- SMPPSim (Java) в папке `SMPPSim-master/`
- Запуск: `cd SMPPSim-master && java -jar target/smppsim.jar conf/logback.xml conf/smppsim.props`
- Порт: `5555`, логин: `pavel`, пароль: `wpsd`

---

## Примерный таймлайн

| Фаза | Оценка |
|---|---|
| Фаза 2: Качество | 2–3 недели |
| Фаза 3: Протокол | 1–2 месяца |
| Фаза 4: Движок нагрузки | 1–2 месяца |
| Фаза 5: CLI и сценарии | 3–4 недели |
| Фаза 6: Distributed | 2–3 месяца |
| Фаза 7: UI и Observability | 1–2 месяца |
| Фаза 8: Полировка | 2–4 недели |
| **Итого** | **~9–14 месяцев** |
