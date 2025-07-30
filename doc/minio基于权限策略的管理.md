# MinIO访问管理与策略配置详解

MinIO是一款高性能的分布式对象存储系统，广泛应用于云计算、大数据等领域。在MinIO中，**访问管理**是通过基于策略的访问控制（PBAC）来实现的。本文将详细介绍MinIO的访问管理、内置策略、策略文档结构以及如何根据具体需求定制权限策略。

## 1. 访问管理概述

MinIO使用基于策略的访问控制（PBAC）来定义认证用户可以执行的操作和访问的资源。每个策略描述了一组操作和条件，规定了用户或用户组的权限。MinIO的PBAC兼容AWS IAM策略的语法、结构和行为，因此可以参考AWS的文档来更深入了解相关功能。

MinIO的`mc admin policy`命令用于创建和管理部署中的策略，帮助管理员灵活控制访问权限。

### 1.1 版本更新：Tag-Based Policy Conditions

在版本`RELEASE.2022-10-02T19-29-29Z`中，MinIO新增了基于标签的条件（Tag-Based Conditions）。通过该功能，策略可以限制用户访问具有特定标签的对象。具体实现方式为在策略的`Condition`语句中使用` s3:ExistingObjectTag/<key>`。

## 2. 内置策略

MinIO提供了多个内置策略，管理员可以将其分配给用户或用户组，以简化权限管理。

### 2.1 内置策略列表

#### 2.1.1 `consoleAdmin`

此策略授予用户对MinIO部署中所有S3和管理API操作的完全访问权限。等同于以下操作：

* `s3:*`
* `admin:*`

#### 2.1.2 `readonly`

该策略为用户提供对MinIO部署中任何对象的只读权限，仅支持GET操作，且不要求列出所有对象。等同于以下操作：

* `s3:GetBucketLocation`
* `s3:GetObject`

例如，以下操作在只读策略下支持：

* `mc cp`
* `mc stat`
* `mc head`
* `mc cat`

但**列出权限**被排除在外，因为典型的只读角色并不需要具有完全的发现能力。

#### 2.1.3 `readwrite`

此策略授予用户对MinIO服务器中所有桶和对象的读写权限。等同于 `s3:*`。

#### 2.1.4 `diagnostics`

该策略授予用户执行诊断操作的权限。包括以下操作：

* `admin:ServerTrace`
* `admin:Profiling`
* `admin:ConsoleLog`
* `admin:ServerInfo`
* `admin:TopLocksInfo`
* `admin:OBDInfo`
* `admin:BandwidthMonitor`
* `admin:Prometheus`

#### 2.1.5 `writeonly`

此策略为用户提供对MinIO部署中任意命名空间（桶及对象路径）的写入权限。等同于 `s3:PutObject`。

### 2.2 策略分配与用户操作

通过 `mc admin policy attach` 命令，可以将策略与用户或用户组关联。例如，以下是一个典型的用户和策略分配示例：

| 用户                              | 策略                                        | 操作               |
| ------------------------------- | ----------------------------------------- | ---------------- |
| Operations                      | `readwrite` on finance bucket             |                  |
| <br/>`readonly` on audit bucket | 对`finance`桶进行PUT和GET操作<br/>对audit桶进行GET操作 |                  |
| Auditing                        | `readonly` on audit bucket                | 对`audit`桶进行GET操作 |
| Admin                           | `admin:*`                                 | 所有mc admin命令     |

在该表中，每个用户只能访问明确授权的资源和操作，MinIO默认拒绝访问任何未授权的资源或操作。

### 2.3 访问控制：拒绝优先

MinIO遵循IAM策略评估规则，在同一操作/资源上，**拒绝（Deny）优先于允许（Allow）**。例如，如果用户有一个包含`Allow`规则的策略，并且其所在组有一个包含`Deny`规则的策略，MinIO将仅应用`Deny`规则。

## 3. 策略文档结构

MinIO的策略文档遵循与AWS IAM策略相同的架构。以下是一个示例策略文档：

```json
{
   "Version" : "2012-10-17",
   "Statement" : [
      {
         "Effect" : "Allow",
         "Action" : [ "s3:<ActionName>", ... ],
         "Resource" : "arn:aws:s3:::*",
         "Condition" : { ... }
      },
      {
         "Effect" : "Deny",
         "Action" : [ "s3:<ActionName>", ... ],
         "Resource" : "arn:aws:s3:::*",
         "Condition" : { ... }
      }
   ]
}
```

* **Statement.Action**：指定一个或多个支持的S3 API操作。
* **Statement.Resource**：指定受限的桶或桶前缀，可以使用`*`和`?`通配符。
* **Statement.Condition**：指定一个或多个支持的条件。

## 4. MinIO策略支持的S3操作

MinIO支持一部分IAM S3操作键，包括常见的S3操作，例如：

* `s3:*`：选择所有S3操作。
* `s3:CreateBucket`：控制`CreateBucket`操作。
* `s3:GetObject`：控制`GetObject`操作。
* `s3:PutObject`：控制`PutObject`操作。

### 4.1 高级操作

MinIO还支持更高级的操作，例如**多部分上传**、**版本控制与保留**、**对象加密**、**桶复制**等。每种操作都支持特定的条件键来进一步控制权限。

## 5. 支持的条件键

MinIO的策略文档支持IAM条件语句，允许根据特定条件来控制访问。常见的条件键包括：

* - `aws:Referer`
  
  - `aws:SourceIp`
  
  - `aws:UserAgent`
  
  - `aws:SecureTransport`
  
  - `aws:CurrentTime`
  
  - `aws:EpochTime`
  
  - `aws:PrincipalType`
  
  - `aws:userid`
  
  - `aws:username`
  
  - `x-amz-content-sha256`
  
  - `s3:signatureAge`

```txt
需要注意的是，aws:Referer、aws:SourceIp
和aws:UserAgent等条件可能会被伪造，
因此MinIO建议仅将这些条件键用于访问拒绝的策略中，作为次级安全措施。
```

## 6. 自定义权限策略

* 创建一个策略文件user-home-policy.json
  
  ```json
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": ["s3:ListBucket"],
        "Effect": "Allow",
        "Resource": ["arn:aws:s3:::mybucket"],
        "Condition": {
          "StringLike": {
            "s3:prefix": ["${aws:username}/*"]
          }
        }
      },
      {
        "Action": [
          "s3:GetObject",
          "s3:PutObject"
        ],
        "Effect": "Allow",
        "Resource": ["arn:aws:s3:::mybucket/${aws:username}/*"]
      }
    ]
  }
  ```
  
  策略的含义：
  
  - 允许用户列出 `mybucket` 存储桶下以自己用户名为前缀的对象路径；
  
  - 允许用户读取和上传这些路径下的对象；
  
  - 实现类似「每人一个家目录」的隔离访问控制。

* 创建策略、
  
  ```shell
  mc admin policy create myminio user-home-policy user-home-policy.json
  ```
  
  命令解释：
  
  * admin policy create：为 MinIO 添加一个自定义策略
  
  * user-home-policy.json：存放策略内容的本地文件路径
  
  * myminio：预先配置的 MinIO 服务别名（通过 mc alias set 设置）
  
  * ser-home-policy：策略名称
  
  * user-home-policy.json:策略文件

* 验证自定义策略是否创建成功
  
  ```shell
  root@iZbp1beh21rgw7a61yfkehZ:~# mc admin policy ls myminio
  diagnostics
  readonly
  readwrite
  user-home-policy
  writeonly
  consoleAdmin
  ```
  
  从上面的输出可以看到我们已经自定义好了user-home-policy

* 创建用户
  
  ```shell
  root:~# mc admin user add myminio zhshl strongpasswd
  Added user `zhshl` successfully.
  ```
  
  - `admin user add`：创建新用户
  
  - `zhshl`：新用户的用户名（也是 Access Key）
  
  - `strongpasswd`：密码（Secret Key）

* 绑定策略给用户
  
  ```shell
  mc admin policy attach myminio user-home-policy --user zhshl
  ```
  
  * admin policy attach: 绑定权限
  
  * user-home-policy: 权限名称
  
  * --user zhshl： 指定绑定的用户为zhshl
  
  * 整条命令实现的是给用户zhshl添加user-home-policy权限

* 验证用户权限
  
  ```shell
  # 
  mc alias set zhshl-minio http://127.0.0.1:9000 zhshl strongpasswd
  #
  echo "this is localfile" >> localfile.txt
  
  # 上传一个对象
  mc cp localfile.txt zhshl-minio/mybucket/zhshl/localfile.txt
  
  # 查看桶内的文件
  mc ls zhshl-minio/mybucket/zhshl/
  
  
  # 尝试删除文件
  root@iZbp1beh21rgw7a61yfkehZ:~# mc rm zhshl-minio/mybucket/zhshl/localfile.txt
  mc: <ERROR> Failed to remove `zhshl-minio/mybucket/zhshl/localfile.txt`. Access Denied.
  
  
  ```
  
  

## 7. 总结

MinIO的访问管理为用户提供了灵活而强大的权限控制功能，通过基于策略的访问控制（PBAC），管理员可以精确地控制用户和用户组的访问权限。无论是通过内置策略还是自定义策略，都能帮助企业实现严格的访问管理，保障数据的安全性和合规性。

有关更多详细信息和配置方法，请参考[MinIO官方文档](https://min.io/docs/minio/linux/administration/identity-access-management/policy-based-access-control.html#)。
