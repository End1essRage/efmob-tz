POST http://subs:8080/subscriptions
Content-Type: application/json

{
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "service_name": "Yandex Plus",
  "price": 400,
  "start_date": "07-2025"
}

HTTP/1.1 201
[Asserts]
jsonpath("$.user_id") == "60601fee-2bf1-4721-ae6f-7636e79a0cba"
jsonpath("$.service_name") == "Yandex Plus"
jsonpath("$.price") == 400

## TODO try read?
