### auth service

Auth Service - сервис авторизации. Реализует три метода:

1. CreateUser - регистрация пользователя
2. Login - авторизация пользователя
3. IsAdmin - получение прав пользователя (админ/не админ)

Команда создания админа:

`go run cmd/create-admin/main.go --email="some@mail.com" --password="password" --storage-path="./storage/auth.db"`
