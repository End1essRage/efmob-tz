GET http://subs:8080/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&page=2&page_size=10&order_by=start_date&direction=desc

HTTP/1.1 200
[Asserts]
jsonpath("$[0].user_id") == "60601fee-2bf1-4721-ae6f-7636e79a0cba"
