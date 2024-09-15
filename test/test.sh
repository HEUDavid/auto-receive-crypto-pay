# webhook 回调
curl -X POST http://localhost:8080/webhook?auth=auth_key \
     -H "Content-Type: application/json" \
     --data @data/Test_ETH_MAINNET.json

# 根据转账地址查询 invoice
bash -c 'curl -s "$1" | python -m json.tool' \
     -- "http://localhost:8080/query_token?from_address=0x71660c4005ba85c37ccec55d0c4493e66fe775d3"

# 查询 invoice 详情
bash -c 'curl -s "$1" | python -m json.tool' \
     -- "http://localhost:8080/token_details?token=942ef637afad1f6ac3860c4dd8a0ff74"
