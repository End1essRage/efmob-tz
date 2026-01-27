POST http://subs:8080/subscriptions
Content-Type: application/json

{
  "user_id": 12345
}

HTTP/1.1 400

[Asserts]
jsonpath "$.error" == "invalid request format"
jsonpath "$.code" == "VALIDATION_ERROR"