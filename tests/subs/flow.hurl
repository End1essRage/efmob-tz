# Создаем подписку (новый пользователь)
POST http://subs:8080/subscriptions
Content-Type: application/json

{
  "user_id": "12601fee-2bf1-4721-ae6f-7636e79a0cba",
  "service_name": "Yandex Plus",
  "price": 400,
  "start_date": "07-2025"
}

HTTP/1.1 201
[Asserts]
jsonpath "$.user_id" == "12601fee-2bf1-4721-ae6f-7636e79a0cba"
jsonpath "$.service_name" == "Yandex Plus"
jsonpath "$.price" == 400
[Captures]
sub_id: jsonpath "$.id"


# Обновляем подписку
PATCH http://subs:8080/subscriptions/{{sub_id}}
Content-Type: application/json
{
  "price": 500,
  "start_date": "08-2025"
}

HTTP/1.1 202

# Читаем подписку, чтобы убедиться, что изменения сохранились
GET http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.service_name" == "Yandex Plus"
jsonpath "$.price" == 500

# Удаляем подписку
DELETE http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 204

# Проверяем, что подписка удалена
GET http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 404
[Asserts]
jsonpath "$.error" == "NOT_FOUND"
