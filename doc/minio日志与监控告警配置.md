## 一、服务日志管理
### 1.1、服务日志目录配置
我们可以通过环境变量来设置minio服务服务日志存储的目录
```shell
  MINIO_LOG_DIR=/var/log/minio
```
### 1.2、 使用下面的docker-compose文件来创建一个minio实例

```yml
services:
  minio2:
    image: quay.io/minio/minio
    container_name: minio2
    ports:
      - "19000:9000"
      - "19001:9001"
    volumes:
      - ~/minio2/data:/data
      - ~/minio2/log:/var/log/minio
    environment:
      MINIO_LOG_DIR: "/var/log/minio"
    command: server /data --console-address ":9001"
    restart: unless-stopped
```
上面容器中minio服务的日志会存储在`/var/log/minio`目录下,并且还把该日志目录挂载到宿主机~/minio2/log目录下。
### 1.3、验证服务日志目录是否生效
使用docker compose来运行minio实例
```shell
    docker compose up -d
    ```
    查看minio服务的日志是否存储在`/var/log/minio`目录下
    ```shell
      cat ~/minio2/log/minio*.log
      
    {"level":"INFO","time":"2025-08-04T08:39:25.977363653Z","message":"MinIO Object Storage Server"}
      ......
```


## 二、日志配置
### 2.1、审计日志地址配置
   * 使用环境变量配置
```shell
      MINIO_AUDIT_WEBHOOK_ENABLE_first: "on"
      MINIO_AUDIT_WEBHOOK_ENDPOINT_first: "http://8.136.102.104:8095/server-log"
      MINIO_AUDIT_WEBHOOK_AUTH_TOKEN_first: "TOKEN"
```
   * 使用mc命令行工具来配置
```shell
   # 语法
   mc admin config set ALIAS/ audit_webhook:IDENTIFIER  \
   endpoint="https://webhook-1.example.net"     
  # 实例
  mc admin config set minio2/ audit_webhook:first  \
     endpoint="http://172.17.73.139:8095/audit-log" \
     http_timeout="10s"   
```

MINIO_AUDIT_WEBHOOK_ENABLE_{id}: 标识开启审核日志
MINIO_AUDIT_WEBHOOK_ENDPOINT_{id}: 标识审核日志的地址
MINIO_AUDIT_WEBHOOK_AUTH_TOKEN_{id}: 标识审核webhook需要使用的token
可以支持多个webhook，通过id来进行区分。

## 三、事件webhook监听
### 3.1、配置
```shell
    export  MINIO_NOTIFY_WEBHOOK_ENABLE_first: "on"
    export  MINIO_NOTIFY_WEBHOOK_ENDPOINT_first: "http://172.17.73.139:8095/event-log"
    # 重启minio
    mc admin service restart minio2
```
### 3.2、获取事件通知的ARN

```shell
 root@iZbp17vix2j58ya7sc3b9lZ:~/docker-compose-dir# mc admin info --json minio2 | jq  .info.sqsARN
[
  "arn:minio:sqs::first:webhook"
]

```

### 3.3、配置事件通知
```shell
mc event add minio2/mybucket arn:minio:sqs::first:webhook \
  --event put,get,delete
```
[支持的事件类型参考][https://docs.min.io/community/minio-object-store/reference/minio-mc/mc-event-add.html#mc-event-supported-events]

### 四、推送事件信息到mysql中
#### 4.1、配置
```shell
    # 使用环境变量
     MINIO_NOTIFY_MYSQL_ENABLE_first="on"
     MINIO_NOTIFY_MYSQL_DSN_STRING_first="remote_user:RemoteUserPassword@123@tcp(116.62.226.77:3306)/syx"
     MINIO_NOTIFY_MYSQL_TABLE_first="minio-events"
     MINIO_NOTIFY_MYSQL_FORMAT_first="namespace"
     
     
    mc admin config set ALIAS/ notify_mysql:IDENTIFIER \
       dsn_string="<ENDPOINT>" \
       table="<string>" \
       format="<string>" \
       max_open_connections="<string>" \
       queue_dir="<string>" \
       queue_limit="<string>" \
       comment="<string>"
    
    mc admin config set  minio2 notify_mysql:first \
      dsn_string="remote_user:RemoteUserPassword@123@tcp(116.62.226.77:3306)/syx" \
      table="minioevents" \
      format="namespace" 
      # 重启服务
    mc admin service restart minio2
```

#### 4.2、添加事件通知
```shell
root@iZbp17vix2j58ya7sc3b9lZ:~# mc admin info --json minio2 | jq  .info.sqsARN
[
  "arn:minio:sqs::first:webhook",
  "arn:minio:sqs::_:mysql"
]

mc event add minio2/mybucket arn:minio:sqs::_:mysql \
  --event put,get,delete
```

#### 4.3、验证事件通知
```shell
    mc put docker-compose.yml minio2/mybucket
```
然后查看mysql中是否有数据

