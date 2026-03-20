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
- [x] Каналы (`make(chan struct{})`, `<-done`, `done <- struct{}{}`)
- [x] Пакеты и импорты (псевдонимы импортов, экспорт через заглавную букву)
- [x] `len([]rune(text))` vs `len(text)` — байты vs символы

### SMPP
- [x] Bind типы: TX, RX, TRX
- [x] Подключение через `go-smpp` библиотеку
- [x] Отправка SMS (`Submit`)
- [x] Длинные сообщения (`SubmitLongMsg`)
- [x] Кодировки: GSM7, Latin1, UCS2 — лимиты и валидация символов
- [x] Поля `ShortMessage`: Register, PriorityFlag, Validity, TON/NPI и др.
- [x] Validity period — относительный (`1h`) и абсолютный (`2006-01-02T15:04:05`)
- [x] DLR (delivery report) — получение через `Handler` в Transceiver
- [x] Формат DLR строки: `id:X sub:001 dlvrd:001 stat:DELIVRD err:000`

### Структура проекта
- [x] Разбивка на пакеты: `smpp/client.go`, `smpp/message.go`, `smpp/encoding.go`, `smpp/dlr.go`
- [x] `main.go` — только точка входа и флаги

---

## In Progress

### DLR парсинг
- [ ] Функция `parseDLR(text string) map[string]string`
- [ ] Вывод распarsенных полей: `id`, `stat`, `err`, `done date`

---

## TODO

### DLR
- [ ] Обработка multipart DLR (несколько частей — несколько DLR)
- [ ] Таймаут ожидания DLR (сейчас ждёт вечно)

### Улучшения
- [ ] Поддержка RX bind типа (сейчас заглушка)
- [ ] TLV fields
- [ ] Multi-destination (`DstList`, `SubmitMulti`)

### Go концепты (впереди)
- [ ] Структуры (`struct`) и методы
- [ ] Тесты (`testing` package)
- [ ] Логирование (`log/slog`)
- [ ] Конфигурационный файл вместо флагов (JSON/YAML)

---

## Библиотеки
- `github.com/fiorix/go-smpp` — SMPP клиент

## Симулятор
- SMPPSim (Java) в папке `SMPPSim-master/`
- Запуск: `cd SMPPSim-master && java -jar target/smppsim.jar conf/logback.xml conf/smppsim.props`
- Порт: `5555`, логин: `pavel`, пароль: `wpsd`
