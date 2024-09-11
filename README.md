<h1 align="center">Auto Receive Crypto Pay</h1>
<h4 align="center">加密货币收款服务，直接的点对点接收，到账后自动执行动作，安全自由，自托管。</h4>


## 简介
通过[流行的节点服务](https://ethereum.org/en/developers/docs/nodes-and-clients/nodes-as-a-service/#popular-node-services)监听地址转账活动，
当 `入账地址ToAddress` 收到加密货币后， `node services` 回调 `https://your_domain/webhook` 接口，
接口先对回调请求落库，然后进行地址匹配，为 `发送地址FromAddress` 生成相关token。

## 使用

```sh
# webhook 回调
curl -X POST http://localhost:8080/webhook \
-H "Content-Type: application/json" \
--data @hook.json
```

```sh
# 根据转账地址查询token
 bash -c 'curl -s "$1" | python -m json.tool' \
 -- "http://localhost:8080/query_token?from_address=0x71660c4005ba85c37ccec55d0c4493e66fe775d3"
```