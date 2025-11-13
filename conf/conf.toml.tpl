# conf.toml 模板

[global]
mode = "debug"                # debug release test
logPath = "log/water.log"
addr = "127.0.0.1:8080        # 127.0.0.1:8080
hostRoot = "/"                # "/" 或 "/abc"
auth = "abc"                  # auth_key

[adminAddress]  # 管理员收款地址: (网络：<地址，图片链接>)
MATIC_MAINNET = [
    { address = "0x0000000000000000000000000000000000000000", url = "https://api.mdavid.cn/gin/src/MATIC.JPG" }
]
ETH_MAINNET = [
    { address = "0x950a4e3beb32d3862272592c8bae79fb5f3475db", url = "https://api.mdavid.cn/gin/src/ETH.JPG" }
]

[mysql]
user = "${DB_USER}"
password = "${DB_PWD}"
host = "gateway01.ap-southeast-1.prod.aws.tidbcloud.com"
port = 4000
dbName = "${DB_NAME}"
maxIdleConns = 10
maxOpenConns = 100
tls = true

[rmq]
user = "${MQ_USER}"
vhost = "/${MQ_VHOST}"
password = "{MQ_PWD}"
host = "armadillo.rmq.cloudamqp.com"
port = 5672
queue = "receipt"
