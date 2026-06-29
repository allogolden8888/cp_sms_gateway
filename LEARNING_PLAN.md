# SMSC Simulator — Learning Plan

## Цель проекта
Профессиональный симулятор SMSC для тестирования SMS-клиентов (ESME).

**Этот репозиторий:** SMSC Simulator — принимает подключения, симулирует поведение реального шлюза по сценариям.  
**Отдельный репозиторий (позже):** ESME Load Tester — генерирует нагрузку на реальный SMSC, до 20k TPS.

### Ключевые возможности симулятора
- Аккаунты (`system_id`) с привязкой к сценариям
- Конструктор сценариев: правила ответов, DLR, задержки, ошибки, валидация
- Режимы правил: вероятностный, счётчик, случайный порядок
- Полная совместимость с SMPP 3.4
- Web UI для управления аккаунтами и сценариями
- Производительность: держать 20k TPS входящих submit_sm

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
- [x] Регулярные выражения
- [x] Структуры (`struct`)

### SMPP (клиентская сторона, через go-smpp)
- [x] Bind типы: TX, RX, TRX
- [x] Отправка SMS, кодировки: GSM7, Latin1, UCS2
- [x] DLR — получение, парсинг, multipart, таймаут
- [x] Validity period — относительный и абсолютный

### Структура
- [x] Разбивка на пакеты, JSON конфиг, `main.go` как точка входа

---

## Фаза 2: Качество кода ✅ (завершена)

- [x] Тесты (`testing`) — табличные тесты, `t.Run`, хелперы
- [x] Структурированное логирование (`log/slog`)
- [x] Вынести DLR-сопоставление в пакет (`DLRTracker`)
- [ ] Правильный `error` wrapping — `fmt.Errorf("%w")`, `errors.Is`, `errors.As`
- [ ] Graceful shutdown — `os/signal`, `signal.NotifyContext`
- [ ] Бенчмарки (`testing.B`) — `go test -bench`

---

## Фаза 3: Бинарный протокол и базовый сервер ✅ (завершена)

### Go концепты
- [x] `encoding/binary` — чтение и запись бинарных данных вручную
- [x] `bytes.Reader` — курсорное чтение буфера
- [x] `io.ReadFull` — чтение ровно N байт
- [x] Struct embedding — переиспользование полей PDU
- [x] Type switch — `switch v := x.(type)`
- [x] `net.Listen` / `net.Accept` / `net.Conn` — TCP сервер

### PDU парсеры (`smpp/pdu/`)
- [x] Заголовок, C-строки, TLV
- [x] `bind.go` — bind tx/rx/trx + resp + unbind
- [x] `submit.go` — submit_sm + resp
- [x] `deliver.go` — deliver_sm + resp
- [x] `session.go` — enquire_link + resp, generic_nack
- [x] `ParsePDU` — публичная точка входа по command_id

### Сериализация PDU
- [x] `SerializeBindResp`
- [x] `SerializeSubmitSMResp`

### Базовый сервер (`smpp/server.go`)
- [x] Приём TCP подключений, чтение PDU (header → body)
- [x] Ответы на bind, submit_sm

---

## Фаза 4: Полный SMPP сервер

### Протокол
- [ ] EnquireLink / EnquireLinkResp — отвечать автоматически (иначе клиент рвёт соединение)
- [ ] Unbind — корректное завершение сессии
- [ ] Generic NACK — ответ на неизвестный или битый PDU
- [ ] Сериализация deliver_sm — для отправки DLR клиенту
- [ ] `submit_multi` — парсинг и ответ
- [ ] Состояние соединения: bind state machine (unbound → bound_tx/rx/trx)

### DLR генерация
- [ ] После submit_sm — отправить deliver_sm с DLR через `time.AfterFunc`
- [ ] DLR content: SMPP TLV `receipted_message_id`, `message_state`
- [ ] Конфигурируемая задержка (из правила сценария)

### Go концепты
- [ ] `time.AfterFunc` — отложенное выполнение
- [ ] `io.Writer` — запись в соединение
- [ ] `context.Context` — управление временем жизни соединения

---

## Фаза 5: Аккаунты и хранилище

### База данных
- [ ] `database/sql` + SQLite (`modernc.org/sqlite`, pure Go, без CGO)
- [ ] Схема: `accounts`, `scenarios`, `rules`
- [ ] CRUD для всех сущностей

### Модели
- [ ] **Account**: `system_id`, `password`, `scenario_id`
- [ ] **Scenario**: `id`, `name`, `description`
- [ ] **Rule**: `scenario_id`, `pdu_type`, `mode`, `params` (JSON)

### Аутентификация
- [ ] При bind — проверять `system_id` + `password` по БД
- [ ] Загружать сценарий аккаунта для соединения

### Go концепты
- [ ] `database/sql` — prepare, query, scan
- [ ] JSON в БД — хранение `params` как `json.RawMessage`
- [ ] Миграции схемы

---

## Фаза 6: Движок сценариев

### Типы правил
- [ ] **Ответ submit_sm**: статус (OK / ESME_R*), задержка, message_id (нормальный / дублирующийся)
- [ ] **DLR**: статус (DELIVRD / UNDELIV / EXPIRED / ...), задержка, вероятность доставки, контент
- [ ] **Bind**: принять / отклонить с кодом ошибки
- [ ] **Валидация**: разрешённые отправители, разрешённые получатели, проверка кодировки, проверка корректности PDU

### Режимы правил
- [ ] **Вероятностный** — X% запросов получают это поведение
- [ ] **Счётчик** — первые N норм, следующие M получают поведение, потом цикл
- [ ] **Случайный порядок** — из 100 сообщений M получат поведение в случайном порядке

### Движок
- [ ] Загрузка сценария при старте сессии
- [ ] Применение правил к каждому PDU
- [ ] Комбинирование нескольких правил (задержка + ошибка + DLR drop)

---

## Фаза 7: Web UI

### REST API (`net/http`)
- [ ] `GET/POST/PUT/DELETE /api/accounts`
- [ ] `GET/POST/PUT/DELETE /api/scenarios`
- [ ] `GET/POST/PUT/DELETE /api/rules`
- [ ] `GET /api/sessions` — активные соединения

### Фронтенд (HTML + vanilla JS, без фреймворков)
- [ ] Список аккаунтов — создать, привязать сценарий
- [ ] Конструктор сценариев — добавить правило, выбрать тип и режим, задать параметры
- [ ] Активные сессии — кто подключён, какой сценарий активен

### Go концепты
- [ ] `net/http` — роутинг, middleware
- [ ] `encoding/json` — REST responses
- [ ] `//go:embed` — встраивание статики в бинарник
- [ ] Server-Sent Events или WebSocket — live обновление сессий

---

## Фаза 8: Производительность (20k TPS)

### Конкурентность
- [ ] Worker pool — горутины + канал задач для обработки PDU
- [ ] `sync.Mutex` — защита состояния соединения
- [ ] `sync/atomic` — lock-free счётчики TPS и метрик
- [ ] `sync.Pool` — переиспользование буферов, снижение GC pressure
- [ ] `context.Context` — отмена, таймауты

### Профилирование и тесты
- [ ] `go test -bench` — бенчмарки парсеров
- [ ] `go test -race` — детектор гонок данных
- [ ] `pprof` — CPU и memory профилирование под нагрузкой
- [ ] Нагрузочный тест: 20k TPS входящих submit_sm

---

## Фаза 9: Observability

### Метрики симулятора
- [ ] TPS входящих PDU (по типам)
- [ ] Активные сессии, bind/unbind events
- [ ] DLR delivery rate по сценарию
- [ ] Ошибки по типам (ESME_R*)
- [ ] Latency ответов

### Экспорт
- [ ] Prometheus формат (`/metrics` endpoint)
- [ ] `log/slog` — структурированные логи всех PDU событий

---

## Фаза 10: Assertion Engine

### Описание ожиданий
- [ ] Assertion модель в БД: `test_run_id`, `metric`, `operator` (≥, ≤, =, ==0), `expected_value`
- [ ] Встроенные assertion типы:
  - DLR delivery rate (%)
  - Latency p50/p95/p99
  - Error rate по типам (ESME_R*)
  - Уникальность message_id
  - Соответствие кодировки
  - Процент PDU прошедших валидацию

### Оценка результатов
- [ ] После завершения теста — сравнить собранные метрики с assertions
- [ ] Результат: PASS / FAIL + детали по каждому assertion
- [ ] Сохранение в БД: `test_results` (`run_id`, `assertion_id`, `actual_value`, `passed`)

### Отчётность
- [ ] REST endpoint `GET /api/runs/:id/report` — JSON отчёт
- [ ] Web UI: страница результатов теста — таблица assertions, PASS/FAIL, фактические значения
- [ ] Сравнение двух запусков — регрессия

---

## Фаза 11: Полировка и дистрибуция

- [ ] `cobra` — субкоманды: `smsc start`, `account add`, `scenario list`
- [ ] Graceful shutdown — дождаться завершения активных сессий
- [ ] Cross-compilation (`GOOS`, `GOARCH`)
- [ ] GoReleaser — автосборка под Windows/Linux/Mac
- [ ] Docker образ
- [ ] `Makefile`

---

## Следующий репозиторий: ESME Load Tester

Отдельный проект. Подключается к реальному SMSC и генерирует нагрузку по сценариям.
- Connection pool (N параллельных SMPP сессий)
- Rate limiter (`golang.org/x/time/rate`) — точный TPS до 20k
- Сценарии нагрузки: контент, кодировки, multipart, validity period
- Метрики: latency, DLR rate, error rate, p50/p95/p99
- Сравнение запусков

---

## Библиотеки

| Библиотека | Назначение |
|---|---|
| `modernc.org/sqlite` | SQLite (pure Go, без CGO) |
| `github.com/spf13/cobra` | CLI субкоманды |
| `golang.org/x/time/rate` | Rate limiter (для ESME) |
| `go.opentelemetry.io/otel` | OpenTelemetry |

---

## Симулятор для разработки
- SMPPSim (Java) в папке `SMPPSim-master/`
- Запуск: `cd SMPPSim-master && java -jar target/smppsim.jar conf/logback.xml conf/smppsim.props`
- Порт: `5555`, логин: `pavel`, пароль: `wpsd`

---

## Примерный таймлайн

| Фаза | Оценка |
|---|---|
| Фаза 4: Полный SMPP | 3–4 недели |
| Фаза 5: БД и аккаунты | 3–4 недели |
| Фаза 6: Движок сценариев | 1–2 месяца |
| Фаза 7: Web UI | 1–2 месяца |
| Фаза 8: Производительность | 1–2 месяца |
| Фаза 9: Observability | 3–4 недели |
| Фаза 10: Assertion Engine | 3–4 недели |
| Фаза 11: Полировка | 2–3 недели |
| **Итого** | **~9–14 месяцев** |