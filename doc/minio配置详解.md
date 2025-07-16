

### MinIO配置项速查表

本文将MinIO的主要配置项按照功能模块进行分类，并以表格形式清晰地展示了每个配置项的环境变量、`mc`命令、功能描述、默认值和关键注意事项。

> **配置优先级说明**: 当同一配置项同时通过环境变量和 `mc admin config set` 命令设置时，**MinIO会优先采用环境变量的值**。

https://min.io/docs/minio/linux/reference/minio-server/settings.html

#### 核心配置 (Core Settings)

这些配置项控制MinIO服务的核心行为。

| 配置项 (Configuration Item) | 环境变量 (Environment Variable)      | mc admin config 命令 / 顶级键  | 默认值 (Default Value)                                       | 描述 (Description)                                           |
| :-------------------------- | :----------------------------------- | :----------------------------- | :----------------------------------------------------------- | :----------------------------------------------------------- |
| **服务器启动选项**          | `MINIO_OPTS`                         | 不适用                         | 无                                                           | 定义 `minio server` 启动时附加的命令行参数。                 |
| **存储卷**                  | `MINIO_VOLUMES`                      | 不适用                         | 无（启动时必需）                                             | 指定MinIO用作后端的存储目录或驱动器。                        |
| **环境配置文件路径**        | `MINIO_CONFIG_ENV_FILE`              | 不适用                         | 适用于通过 `systemd` 管理 MinIO 服务的情况,一般在/etc/default/minio | 指定MinIO服务器用于加载环境变量的文件的完整路径。            |
| **ILM过期工作线程数**       | `MINIO_ILM_EXPIRY_WORKERS`           | `ilm expiry_workers`           | 最多使用一半的可用CPU核心数                                  | 指定用于处理对象生命周期（ILM）过期规则的工作线程数。        |
| **域名**                    | `MINIO_DOMAIN`                       | `domain`                       | 未设置（使用路径风格请求）                                   | 为部署设置一个FQDN，以启用虚拟主机风格的请求。               |
| **扫描器速度**              | `MINIO_SCANNER_SPEED`                | `scanner speed`                | `default`                                                    | 管理后台扫描器（用于复制、修复等）的速度。                   |
| **启用压缩**                | `MINIO_COMPRESSION_ENABLE`           | `compression enable`           | `off`                                                        | 设置为 `on` 以对新对象启用数据压缩。                         |
| **压缩后加密**              | `MINIO_COMPRESSION_ALLOW_ENCRYPTION` | `compression allow_encryption` | `off`                                                        | 设置为 `on` 以允许在压缩对象后再对其进行加密。（**注意安全风险**） |
| **压缩文件扩展名**          | `MINIO_COMPRESSION_EXTENSIONS`       | `compression extensions`       | `.txt, .log, .csv, .json, .tar, .xml, .bin`                  | 指定要压缩的文件扩展名列表（逗号分隔）。                     |
| **压缩MIME类型**            | `MINIO_COMPRESSION_MIME_TYPES`       | `compression mime_types`       | `text/*, application/json, application/xml, binary/octet-stream` | 指定要压缩的MIME类型列表（逗号分隔）。                       |
| **纠删集驱动器数**          | `MINIO_ERASURE_SET_DRIVE_COUNT`      | `erasure_set drive_count`      | 自动选择                                                     | （**仅限初始化前**）为服务器池中的所有驱动器应用纠删集大小。**请勿随意更改**。 |
| **最大对象版本数**          | `MINIO_API_OBJECT_MAX_VERSIONS`      | `api object_max_versions`      | ~9.2 x 10^18 (Int64最大值)                                   | 定义每个对象允许的最大版本数。建议设置为不超过100。          |

---

#### 根用户访问配置 (Root Access Settings)

这些配置项控制拥有超级管理员权限的根用户的访问。

| 配置项 (Configuration Item) | 环境变量 (Environment Variable) | mc admin config 命令 / 顶级键 | 默认值 (Default Value) | 描述 (Description)                                           |
| :-------------------------- | :------------------------------ | :---------------------------- | :--------------------- | :----------------------------------------------------------- |
| **根用户名**                | `MINIO_ROOT_USER`               | `root_user`                   | `minioadmin`           | 根用户的访问密钥。**切勿在生产中使用默认值**。               |
| **根用户密码**              | `MINIO_ROOT_PASSWORD`           | `root_password`               | `minioadmin`           | 根用户的秘密密钥。**切勿在生产中使用默认值**。               |
| **根用户API访问**           | `MINIO_API_ROOT_ACCESS`         | `api root_access`             | `on`                   | 设置为 `off` 可禁用根用户账户。禁用前务必创建其他管理员账户。 |

---

#### 纠删码与存储类别 (Erasure Code & Storage Class Settings)

这些配置项定义了数据的冗余和可用性策略。

| 配置项 (Configuration Item) | 环境变量 (Environment Variable) | mc admin config 命令 / 顶级键 | 默认值 (Default Value)                          | 描述 (Description)                                           |
| :-------------------------- | :------------------------------ | :---------------------------- | :---------------------------------------------- | :----------------------------------------------------------- |
| **标准存储类**              | `MINIO_STORAGE_CLASS_STANDARD`  | `storage_class standard`      | 依赖于纠删集大小（例如8-16个驱动器默认为 EC:4） | 为`STANDARD`存储类设置奇偶校验级别，格式为 `EC:M`。          |
| **低冗余存储类**            | `MINIO_STORAGE_CLASS_RRS`       | `storage_class rrs`           | 驱动器数>1时为`EC:1`，=1时为`EC:0`              | 为`REDUCED_REDUNDANCY`存储类设置奇偶校验级别。               |
| **奇偶校验保留优化**        | `MINIO_STORAGE_CLASS_OPTIMIZE`  | `storage_class optimize`      | 未设置（默认行为是提升奇偶校验）                | 设置为 `capacity` 可优先保证集群容量，而不是提升对象的奇偶校验级别。 |

---

#### MinIO控制台配置 (MinIO Console Settings)

这些配置项用于管理嵌入式的Web控制台。

| 配置项 (Configuration Item) | 环境变量 (Environment Variable)  | mc admin config 命令 / 顶级键 | 默认值 (Default Value) | 描述 (Description)                                    |
| :-------------------------- | :------------------------------- | :---------------------------- | :--------------------- | :---------------------------------------------------- |
| **启用/禁用控制台**         | `MINIO_BROWSER`                  | `browser`                     | `on` (启用)            | 设置为 `off` 可禁用嵌入式MinIO控制台。                |
| **控制台重定向**            | `MINIO_BROWSER_REDIRECT`         | `browser redirect`            | `true`                 | 是否将来自浏览器的请求自动重定向到控制台。            |
| **控制台重定向URL**         | `MINIO_BROWSER_REDIRECT_URL`     | `browser redirect_url`        | 无（监听所有主机IP）   | 指定控制台监听的FQDN，用于反向代理场景。              |
| **控制台会话时长**          | `MINIO_BROWSER_SESSION_DURATION` | `browser session_duration`    | `12h`                  | 设置控制台登录会话的持续时间。                        |
| **HSTS秒数**                | `MINIO_BROWSER_HSTS_SECONDS`     | `browser hsts_seconds`        | `0` (禁用)             | 启用HSTS（HTTP严格传输安全）并设置 `max-age` 的秒数。 |
| **Prometheus URL**          | `MINIO_PROMETHEUS_URL`           | `prometheus url`              | 无                     | 指定用于抓取MinIO指标的Prometheus服务URL。            |
| **Prometheus Job ID**       | `MINIO_PROMETHEUS_JOB_ID`        | `prometheus job_id`           | `minio-job`            | 指定Prometheus抓取任务的ID。                          |
| **Prometheus认证令牌**      | `MINIO_PROMETHEUS_AUTH_TOKEN`    | `prometheus auth_token`       | 无                     | 指定控制台连接到Prometheus时使用的Basic Auth令牌。    |

---

#### 指标与日志配置 (Metrics & Logging Settings)

这些配置项用于将MinIO的运行状态导出到外部系统。

| 配置项 (Configuration Item) | 环境变量 (Environment Variable)      | mc admin config 命令 / 顶级键  | 默认值 (Default Value) | 描述 (Description)                                           |
| :-------------------------- | :----------------------------------- | :----------------------------- | :--------------------- | :----------------------------------------------------------- |
| **Prometheus认证类型**      | `MINIO_PROMETHEUS_AUTH_TYPE`         | `prometheus auth_type`         | `jwt`                  | `/metrics`端点的认证模式。`jwt`（需要令牌）或 `public`（公开）。 |
| **启用服务器日志Webhook**   | `MINIO_LOGGER_WEBHOOK_ENABLE[_ID]`   | `logger_webhook[_ID] enable`   | `off`                  | 设置为 `on` 以将服务器日志发送到Webhook。                    |
| **服务器日志Webhook端点**   | `MINIO_LOGGER_WEBHOOK_ENDPOINT[_ID]` | `logger_webhook[_ID] endpoint` | 无（若启用则必需）     | Webhook的HTTP(S) URL。                                       |
| **启用审计日志Webhook**     | `MINIO_AUDIT_WEBHOOK_ENABLE[_ID]`    | `audit_webhook[_ID] enable`    | `off`                  | 设置为 `on` 以将审计日志发送到Webhook。                      |
| **审计日志Webhook队列大小** | `MINIO_AUDIT_WEBHOOK_QUEUE_SIZE`     | `audit_webhook queue_size`     | `100000`               | 审计Webhook目标的队列大小。                                  |
| **启用Kafka审计日志**       | `MINIO_AUDIT_KAFKA_ENABLE`           | `audit_kafka enable`           | `off`                  | 设置为 `on` 以将审计日志发送到Kafka。                        |
| **Kafka SASL机制**          | `MINIO_AUDIT_KAFKA_SASL_MECHANISM`   | `audit_kafka sasl_mechanism`   | `plain`                | MinIO用于对Kafka进行身份验证的SASL机制。                     |
| **Kafka队列大小**           | `MINIO_AUDIT_KAFKA_QUEUE_SIZE`       | `audit_kafka queue_size`       | `100000`               | Kafka目标的队列大小。                                        |









