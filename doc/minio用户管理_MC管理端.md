# 🚀MinIO 用户管理与权限控制详解

> 本文面向具备一定 MinIO 运维或开发经验的读者，详细讲解 MinIO 的用户管理机制、访问密钥（Access Key）使用方式，以及如何通过策略控制用户权限。

---

## 📚一、概述

在 MinIO 中，**用户（User）** 是由唯一的访问密钥（Access Key，相当于用户名）和对应的密钥（Secret Key，相当于密码）组成的身份实体。

客户端在访问 MinIO 时，必须同时提供有效的 Access Key 和 Secret Key 来进行身份验证。

MinIO 的权限控制采用基于策略（Policy-Based）的方式：

* 每个用户可以**直接绑定一个或多个权限策略**；
* 用户也可以通过加入**用户组**继承权限策略；
* **默认情况下，MinIO 拒绝一切未被授权的操作或资源访问**，因此必须显式为用户或其所在组绑定权限策略。

## 👑二、MinIO Root 用户

每个 MinIO 实例在启动时都会自动生成一个超级管理员账户（Root User），拥有系统内所有资源的访问权限。

### 设置 Root 用户的方式：

通过以下环境变量设置：

```bash
MINIO_ROOT_USER=<AccessKey>
MINIO_ROOT_PASSWORD=<SecretKey>
```

### 安全建议：

* 使用**长、唯一且随机的字符串**作为 Root 凭证；
* 严格限制 Root 凭证的访问权限，仅授权给信任的管理员；
* **强烈不建议在开发、测试或生产环境中用于普通业务访问**；
* 更推荐为不同客户端创建独立的最小权限用户。

> 若未设置以上环境变量，MinIO 默认会使用 `minioadmin:minioadmin` 作为 Root 凭证，**强烈建议禁止默认凭证的使用。**

---

## 👥三、用户管理实战

准备，配置好mc客户端

```shell
mc alias set myminio http://localhost:9000 minioadmin minioadmin
```

### ✅ 创建用户

使用以下命令在 MinIO 中创建一个新用户：

```bash
# 语法
mc admin user add ALIAS ACCESSKEY SECRETKEY
#实例
mc admin user add myminio zhshl zhshl@pwd
```

* **ALIAS**：MinIO 部署的别名；
* **ACCESSKEY**：访问密钥（用户名）；
* **SECRETKEY**：密码（创建后无法再次获取，请妥善保存）。

⚠️ 建议使用长、唯一、随机的字符串作为 ACCESSKEY 和 SECRETKEY，以符合企业安全或合规要求。

### ✅查看用户及内置策略

新创建的用户没有权限访问桶、需要去绑定对应的权限策略，这里我们展示如何查看权限策略及用户

```shell
# 查看用户
mc admin user list myminio
# 查看策略
mc admin policy list myminio
```

consoleAdmin：管理员相关权限
diagnostics：诊断和监控
readonly：只读权限
readwrite：读写权限
writeonly：只写权限

测试没有绑定权限时，是否能够操作minio集群

```shell
mc alias set new_user_with_out_policy http://localhost:9000 zhshl zhshl@pwd


mc alias set new_user_with_out_policy http://localhost:9000 zhshl zhshl@pwd

root@iZbp1beh21rgw7a61yfkehZ:~/minio-binaries# mc ls new_user_with_out_policy
mc: <ERROR> Unable to list folder. Access Denied.
```

从上面的错误信息看出来，新建的用户无法使用。

### 🔐 绑定权限策略

创建用户后，必须为其绑定权限策略（如内置的 `readwrite` 策略）：

```bash
#语法
mc admin policy attach ALIAS readwrite --user=ACCESSKEY
# 实例
root@iZbp1beh21rgw7a61yfkehZ:~/minio-binaries# ./mc admin policy attach myminio readwrite --user=zhshl
Attached Policies: [readwrite]
To User: zhshl
```

替换 `ACCESSKEY` 为刚刚创建的用户名。 

### 🚫禁用用户

可以使用下面的命令来禁用用户

```shell
#语法
mc admin user disable ALIAS USERNAME
#示例
mc admin user disable myminio zhshl
```

* USERNAME： 要禁用的用户名

* ALIAS:minio集群的别名

### ❌ 删除用户

使用以下命令删除指定用户：

```bash
mc admin user rm ALIAS USERNAME
```

* 替换 `USERNAME` 为要删除的用户 Access Key。

---

## 🧩四、小结

| 功能     | 命令                       | 说明                              |
| ------ | ------------------------ | ------------------------------- |
| 创建用户   | `mc admin user add`      | 创建用户并指定密钥                       |
| 查看用户信息 | `mc admin user info`     | 可查看 Access Key（不可查看 Secret Key） |
| 绑定权限策略 | `mc admin policy attach` | 附加 RBAC 策略                      |
| 禁用用户   | mc admin user disable    | 禁用用户账户                          |
| 删除用户   | `mc admin user rm`       | 移除用户账户                          |

## 📌建议与最佳实践

* **最小权限原则**：为每个客户端分配独立账号，并限制仅能访问其业务所需资源；
* **禁用默认凭证**：部署上线前务必更换默认 `minioadmin`；
* **避免频繁使用 Root**：Root 用户应只用于初始化或紧急运维；
* **结合 IAM 策略精细化控制**：充分利用 MinIO 的策略机制进行资源隔离与操作控制。

📎 参考链接：

* [MinIO 用户管理文档](https://min.io/docs/minio/linux/administration/identity-access-management/minio-user-management.html)
* [MinIO 策略控制文档](https://docs.min.io/community/minio-object-store/administration/identity-access-management/policy-based-access-control.html)
