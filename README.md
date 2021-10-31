# Работа с Access JWT и Refresh токенами (Часть сервиса аутентификации)

## Список используемых технологий:
* [JWT](https://github.com/dgrijalva/jwt-go) - Работа с JWT токенами
* [MongoDB](https://github.com/mongodb/mongo-go-driver) - MongoDB в качестве NoSQL базы данных
* [Docker](https://www.docker.com/) - Docker
* [viper](https://github.com/spf13/viper) - Работа с файлами конфигурации

## Для запуска приложения:

``` bash
make build && make run
```

# Примеры запросов/ответов

[![Run in Postman](https://run.pstmn.io/button.svg)](https://www.postman.com/orange-station-722848/workspace/auth-service/documentation/17406947-2812dfb9-3745-463b-8afc-54bedda4b3dd)

### Получить пару Access, Refresh токенов
``` bash
POST http://localhost:8080/auth/get-tokens?id=617e8dfc47336aa7cdd4e535
```

### Выполнить Refresh операцию на пару Access, Refresh токенов
``` bash
POST http://localhost:8080/auth/refresh
# Body(json)
{
    "id": "617e7fa6c76fb924fd95a879",
    "token": "9fcbdf85f635398084722a271af20942f776ccd4241c77d1c37131348547cf68"
}
```

### Проверить валидность токена 
``` bash
GET http://localhost:8080/api/check-access
```
Для проверки работы запроса нужно выбрать использовать header 
"Authorization: Bearer \<token>"

### Создать пользователя (возвращает guid пользователя)
``` bash
POST http://localhost:8080/auth/create-user
```

## Mongo-express

Чтобы просмотреть содержимое базы данных Mongo: http://localhost:8081/db/auth-db/
