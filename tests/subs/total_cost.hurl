GET http://subs:8080/subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&from=07-2025&to=07-2026

HTTP/1.1 200
[Asserts]
jsonpath("$.total") >= 0
