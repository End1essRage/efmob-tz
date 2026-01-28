# Создаем подписку
POST http://subs:8080/subscriptions
Content-Type: application/json

{
  "user_id": "32321fee-2bf1-4721-ae6f-7636e79a0cba",
  "service_name": "Yandex Plus",
  "price": 400,
  "start_date": "07-2025",
  "end_date": "09-2025"
}

HTTP/1.1 201
[Captures]
sub_id: jsonpath "$.id"

# проверяем update без end_date
PATCH http://subs:8080/subscriptions/{{sub_id}}
Content-Type: application/json
{
  "price": 500,
  "start_date": "08-2025"
}

HTTP/1.1 202
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == "09-2025"

GET http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == "09-2025"

# проверяем зануление end_date 1
PATCH http://subs:8080/subscriptions/{{sub_id}}
Content-Type: application/json
{
  "end_date": ""
}

HTTP/1.1 202
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == ""

GET http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == ""

# проверяем зануление end_date 2
PATCH http://subs:8080/subscriptions/{{sub_id}}
Content-Type: application/json
{
  "end_date": null
}

HTTP/1.1 202
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == ""

GET http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == ""

# проверяем зануление end_date 3
PATCH http://subs:8080/subscriptions/{{sub_id}}
Content-Type: application/json
{
  "end_date": "null"
}

HTTP/1.1 202
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == ""

GET http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == ""

# проверяем обновление end_date
PATCH http://subs:8080/subscriptions/{{sub_id}}
Content-Type: application/json
{
  "end_date": "10-2025"
}

HTTP/1.1 202
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == "10-2025"

GET http://subs:8080/subscriptions/{{sub_id}}

HTTP/1.1 200
[Asserts]
jsonpath "$.price" == 500
jsonpath "$.start_date" == "08-2025"
jsonpath "$.end_date" == "10-2025"