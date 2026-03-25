# SMPP CLI — Learning Plan

## Completed

### Go basics
- [x] Флаги командной строки (`flag` package)
- [x] Указатели и разыменование (`*`, `&`)
- [x] Форматирование строк (`fmt.Sprintf`, `fmt.Printf`, raw strings с бэктиками)
- [x] Логические операторы (`&&`, `||`, `!`)
- [x] `switch/case`
- [x] Интерфейсы — объявление и реализация
- [x] Множественный возврат из функции `(value, error)`
- [x] `defer` — отложенный вызов
- [x] Обработка ошибок (`if err != nil`)
- [x] Горутины (`go func()`)
- [x] Каналы — небуферизованные и буферизованные, `select`, `time.After`
- [x] Пакеты и импорты (псевдонимы импортов, экспорт через заглавную букву)
- [x] `len([]rune(text))` vs `len(text)` — байты vs символы
- [x] Округление вверх при целочисленном делении
- [x] Регулярные выражения (`regexp.MustCompile`, `FindStringSubmatch`)
- [x] Структуры (`struct`) — объявление и возврат из функции

### SMPP
- [x] Bind типы: TX, RX, TRX
- [x] Подключение через `go-smpp` библиотеку
- [x] Отправка SMS (`Submit`)
- [x] Длинные сообщения (`SubmitLongMsg`) + подсчёт частей
- [x] Кодировки: GSM7, Latin1, UCS2 — лимиты и валидация символов
- [x] Поля `ShortMessage`: Register, PriorityFlag, Validity, TON/NPI и др.
- [x] Validity period — относительный (`1h`) и абсолютный (`2006-01-02T15:04:05`)
- [x] DLR — получение через `Handler` в Transceiver
- [x] Парсинг DLR строки — структура `DLR`, функция `ParseDLR`
- [x] Multipart DLR — буферизованный канал + счётчик частей
- [x] Таймаут ожидания DLR — из validity period, дефолт 24h

### Структура проекта
- [x] Разбивка на пакеты: `smpp/client.go`, `smpp/message.go`, `smpp/encoding.go`, `smpp/dlr.go`
- [x] `main.go` — только точка входа и флаги

---

## TODO

### DLR
- [ ] Сопоставление `message_id` из `submit_sm_resp` и из DLR через `map`

### Улучшения
- [ ] Поддержка RX bind типа (сейчас заглушка)
- [ ] TLV fields
- [ ] Multi-destination (`DstList`, `SubmitMulti`)

### Go концепты (впереди)
- [x] Конфигурационный файл вместо флагов (JSON + `os.ReadFile` + `json.Unmarshal`)
- [ ] Тесты (`testing` package) ← **сейчас**
- [ ] Логирование (`log/slog`)

---

## Библиотеки
- `github.com/fiorix/go-smpp` — SMPP клиент

## Симулятор
- SMPPSim (Java) в папке `SMPPSim-master/`
- Запуск: `cd SMPPSim-master && java -jar target/smppsim.jar conf/logback.xml conf/smppsim.props`
- Порт: `5555`, логин: `pavel`, пароль: `wpsd`
