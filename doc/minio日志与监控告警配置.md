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
      MINIO_AUDIT_WEBHOOK_ENDPOINT_first: "http://172.17.73.139:8095/server-log"
      MINIO_AUDIT_WEBHOOK_AUTH_TOKEN_first: "TOKEN"

shell
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
    
    mc admin config set minio2/ notify_webhook:IDENTIFIER \
   endpoint="http://172.17.73.139:8095/event-log" 
    
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

### 3.4、验证事件
```shell
     mc put docker-compose.yml minio2/mybucket
```
### 四、推送事件信息到mysql中
#### 4.1、配置
```shell
    # 使用环境变量
     MINIO_NOTIFY_MYSQL_ENABLE_first="on"
     MINIO_NOTIFY_MYSQL_DSN_STRING_first="remote_user:RemoteUserPassword@123@tcp(localhost:3306)/syx"
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
      dsn_string="remote_user:RemoteUserPassword@123@tcp(localhost:3306)/syx" \
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

## 五、监控管理
### 5.1、使用docker安装prometheus
```shell
mkdir -p ~/prometheus/config
cd ~/prometheus/config

# 创建prometheus的配置文件
vi  prometheus.yml
# prometheus.yml
global:
  scrape_interval: 15s  # 每 15 秒抓取一次目标

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
# 创建prometheus容器
docker run -d \
  --name prometheus \
  -p 9090:9090 \
  -v ~/prometheus/config:/etc/prometheus/ \
  prom/prometheus
```

### 5.2、生成prometheus使用的监控时的token
```shell
root@iZbp17vix2j58ya7sc3b9lZ:~/prometheus/config# mc admin prometheus generate minio2 
scrape_configs:
- job_name: minio-job
  bearer_token: eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJwcm9tZXRoZXVzIiwic3ViIjoibWluaW9hZG1pbiIsImV4cCI6NDkwNzkxODMxNX0.AN0QiamO_Np_2eeKPxZWf30ttlak01q--v9KI4mPMiWX3XLt4fkGJ8A8Le2_1Rs0r7QOZ4fA1Hug4_zVtTq7wA
  metrics_path: /minio/v2/metrics/cluster
  scheme: http
  static_configs:
  - targets: ['localhost:19000']
```
这里的localhost需要换成docker0网卡地址


### 5.3、重启prometheus
```shell
  docker restart pormetheus
```

### 5、4 验证数据是否已经采集到

访问对应的地址信息：
http://121.43.30.5:9090/graph?g0.expr=minio_cluster_usage_objects_count&g0.tab=1&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h


### 5.5、添加minio的监控告警规则
#minio-alerting.yml

```yml
groups:
- name: minio-alerts
  rules:
    - alert: NodesOffline
      expr: avg_over_time(minio_cluster_nodes_offline_total{job="minio-job"}[5m]) > 0
      for: 10m
      labels:
      severity: warn
      annotations:
      summary: "Node down in MinIO deployment"
      description: "Node(s) in cluster {{ $labels.instance }} offline for more than 5 minutes"

    - alert: DisksOffline
      expr: avg_over_time(minio_cluster_drive_offline_total{job="minio-job"}[5m]) > 0
      for: 10m
      labels:
      severity: warn
      annotations:
      summary: "Disks down in MinIO deployment"
      description: "Disks(s) in cluster {{ $labels.instance }} offline for more than 5 minutes"
```
在prometheus的配置规则中添加
```yml
    rule_files:
    - minio-alerting.yml
```

### 5.5、配置可视化看板Grafana
* 安装grafana
```shell
  mkdir -p ~/docker-compose-grafana
  cd ~/docker-compose-grafana
 #docker-compose.yml
services:
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    ports:
      - "3000:3000"   # 本地访问 http://localhost:3000
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana

volumes:
  grafana-storage:

docker compose up -d
```
* 打开grafana地址
  http://localhost:3000
* 配置prometheus源
* 导入minio的看板数据
从下面地址下载：https://grafana.com/grafana/dashboards/13502-minio-dashboard/