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

## Прокси для исходящих подключений приложения

Если приложению (app) для подключения к VLESS-серверу требуется SOCKS5-прокси (например, через цепочку прокси), задайте переменные окружения `https_proxy` и `http_proxy`:

```shell
https_proxy=socks5://127.0.0.1:1080 http_proxy=socks5://127.0.0.1:1080 app
```

## Параметры

| Переменная   | Флаг       | Описание                         | По умолчанию |
|-------------|------------|----------------------------------|--------------|
| `VLESS_URI` | `-uri`     | VLESS URI подключения            | —            |
| `SOCKS_PORT`| `-port`    | Порт локального SOCKS5-прокси    | `1080`       |

`.env.example`- пример конфигурационного файла, переименуйте в .env и добавьте свои параметры (uri и port).
