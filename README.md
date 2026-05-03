# go-vless-client

VLESS-клиент на базе Xray-core, предоставляющий локальный SOCKS5-прокси.

## Быстрый старт

Создайте файл `.env` в корне проекта:

```env
VLESS_URI=vless://YOUR_UUID@YOUR_SERVER_ADDRESS:PORT?encryption=none&security=reality&sni=example.com&fp=chrome&pbk=YOUR_PUBLIC_KEY&sid=YOUR_SHORT_ID&type=tcp
SOCKS_PORT=1080
```

Запустите через Docker Compose:

```shell
docker compose up -d
```

Прокси будет доступен на `127.0.0.1:SOCKS_PORT` (по умолчанию `1080`).

Чтобы остановить:

```shell
docker compose down
```

## Использование

```shell
# Без Docker
go run . -uri "vless://..." -port 1080

# Сборка
make build
```

## Параметры

| Переменная   | Флаг       | Описание                         | По умолчанию |
|-------------|------------|----------------------------------|--------------|
| `VLESS_URI` | `-uri`     | VLESS URI подключения            | —            |
| `SOCKS_PORT`| `-port`    | Порт локального SOCKS5-прокси    | `1080`       |
