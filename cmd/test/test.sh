# webhook 回调
curl -X POST http://localhost:8080/webhook \
-H "Content-Type: application/json" \
--data @hook.json

# 根据转账地址查询token
 bash -c 'curl -s "$1" | python -m json.tool' \
 -- "http://localhost:8080/query_token?from_address=0x71660c4005ba85c37ccec55d0c4493e66fe775d3"
