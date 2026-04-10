# GREEN-API Demo (HTML + Go Backend Proxy)

Тестовое задание: подготовлена HTML-страница для вызова методов **GREEN-API** и отображения результатов в отдельных полях **read-only**, а также реализован небольшой **Go backend (прокси)** для корректной работы из браузера (CORS).

## Что реализовано

На фронтенде доступны кнопки (по ТЗ):

1. `getSettings`
2. `getStateInstance`
3. `sendMessage`
4. `sendFileByUrl`

Параметры подключения:

- `idInstance`
- `ApiTokenInstance`

Ответы методов выводятся справа в отдельных блоках:
- `getSettings (Ответ)`
- `getStateInstance (Ответ)`
- `sendMessage (Ответ)`
- `sendFileByUrl (Ответ)`

## Почему нужен Go backend

Запросы напрямую из браузера к `https://api.green-api.com` блокируются политикой **CORS** (в т.ч. preflight для `POST`), поэтому фронтенд не может вызвать GREEN-API “в лоб”.

Чтобы страница работала корректно, реализован Go-сервер-прокси:
- фронтенд делает запросы на `http(s)://<backend-host>/api/...`
- backend уже отправляет запросы в GREEN-API и возвращает ответ обратно на страницу

## Требования / форматы данных

### chatId (для sendMessage и sendFileByUrl)
GREEN-API требует один из форматов:

- `phone_number@c.us` (для контакта)
- `group_id@g.us` (для группы)
- `lid_id@lid` (для LID)

Пример для отправки “себе” (контакт):
- `79991234567@c.us`

### fileUrl (для sendFileByUrl)
- ссылка должна быть **прямой HTTPS-ссылкой** на файл (который доступен без авторизации)
- backend автоматически извлекает `fileName` из URL пути

## Структура репозитория

В этом проекте используются:
- `index.html` — фронтенд (страница с UI и вызовами)
- `server.go` — backend-прокси на Go

## Локальный запуск

### 1) Запуск backend
```bash
go run server.go
После запуска backend будет слушать:

http://localhost:8080
endpoint’ы:

GET  /api/getSettings
GET  /api/getStateInstance
POST /api/sendMessage
POST /api/sendFileByUrl

```

### 2) Открытие фронтенда
Открывайте в браузере:
````
http://localhost:8080/


Фронтенд раздается backend’ом через статическую отдачу index.html.

Тестирование методов (порядок как в ТЗ)

Ввести idInstance и ApiTokenInstance
Нажать:

getSettings
getStateInstance


Ввести:

chatId
message


Нажать:

sendMessage


Ввести:

fileUrl


Нажать:

sendFileByUrl

````
### Деплой (GitHub Pages + Backend)


GitHub Pages используется для публикации только HTML-страницы (index.html) и получения публичной ссылки (как в ТЗ).
Go backend публикуется отдельно на любом хостинге с HTTPS (Render/Fly/Railway/VPS и т.п.).

После деплоя backend:

в index.html нужно заменить:

BACKEND_BASE на URL вашего backend (например https://your-backend.example.com)