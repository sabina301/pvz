# Пункт выдачи заказов
### Быстрый старт

1) Склонируйте репозиторий
```
https://github.com/sabina301/pvz.git
cd pvz
```
2) Запустите docker-compose
```
docker-compose up --build
```

### Тестирование
```
go test ./...
```

Данный сервис ПВЗ умеет:
1. Совершать регистрацию и вход пользователей:
   ```
   POST http://some_host:some_port/api/v1/register
   POST http://some_host:some_port/api/v1/login
   ```
2. Совершать выдачу токена:
   ```
   POST http://some_host:some_port/api/v1/dummyLogin
   ```
3. Создавать пвз:
   ```
   POST http://some_host:some_port/api/v1/pvz
   ```
5. Добавлять информацию о приемке товаров:
   ```
   POST http://some_host:some_port/api/v1/receptions
   ```
7. Добавлять товары в рамках 1й приёмки:
   ```
   POST http://some_host:some_port/api/v1/product
   ```
9. Удалять товары в открытой приёмке:
    ```
    POST http://some_host:some_port/api/v1/pvz/{pvzId}/delete_last_product
    ```
11. Закрывать приёмку:
    ```
    POST http://some_host:some_port/api/v1/pvz/{pvzId}/close_last_reception
    ```
13. Получать данные о пвз:
    ```
    GET http://some_host:some_port/api/v1/pvz
    ```
### Вопросы
1. Про роли. Непонятно, нужен ли client. Не стала добавлять
