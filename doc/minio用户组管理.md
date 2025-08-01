以下是一篇适合技术文档风格的内容，基于你提供的 MinIO `mc admin group` 指令文档内容整理而成：

---

# MinIO 用户组管理命令使用指南（`mc admin group`）

MinIO 提供了基于策略的访问控制（Policy-Based Access Control, PBAC），以实现对用户权限的细粒度管理。`mc admin group` 命令用于在 MinIO 部署中管理用户组（Group）。用户组可集中分配访问策略，使权限管理更加简洁高效。

> ⚠️ 仅适用于 MinIO 部署，不适用于其他 S3 兼容服务。

---

## 一、用户组与策略控制

用户组是用户的集合，每个组可以分配一个或多个访问策略（Policy）。用户一旦加入某个组，将继承该组所绑定的策略权限。MinIO 的权限评估逻辑遵循如下规则：

* 用户的最终权限 = 显式绑定的用户策略 + 所属组继承的策略
* 当某个资源或操作上存在同时“允许”和“拒绝”的规则时，**拒绝（Deny）优先生效**

详细规则参考 IAM 文档：[Determining Whether a Request is Allowed or Denied Within an Account](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_evaluation-logic.html)

---

## 二、常用命令

### 1. 创建用户组

```bash
mc admin group add TARGET GROUPNAME USER1 [USER2 USER3...]

root@iZbp17vix2j58ya7sc3b9lZ:~# mc admin group add src-minio default-group foobar foobar2
Added members `foobar,foobar2` to group `default-group` successfully.

```

* `TARGET`：MinIO 部署的别名
* `GROUPNAME`：要创建的组名
* `USER1...`：一个或多个已有的用户名

> 若组不存在，将自动创建。

---

### 2. 查看所有用户组

```bash
mc admin group ls TARGET
root@iZbp17vix2j58ya7sc3b9lZ:~# mc admin group ls src-minio
default-group
```

列出指定 MinIO 部署下的所有组。

---
### 3.给组添加策略
```shell
    mc admin policy attach TARGET POLICYNAME --group GROUPNAME --policy 
    root@iZbp17vix2j58ya7sc3b9lZ:~# mc admin policy attach src-minio readonly --group default-group
Attached Policies: [readonly]
To Group: default-group

```

### 4. 查看组详情

```bash
mc admin group info TARGET GROUPNAME
root@iZbp17vix2j58ya7sc3b9lZ:~# mc admin group info src-minio default-group
Group: default-group
Status: enabled
Policy: readonly
Members: foobar,foobar2
```

显示指定组的详细信息，包括组成员及已绑定的策略。

---
### 5. 禁用用户组

```bash
mc admin group disable TARGET GROUPNAME
root@iZbp17vix2j58ya7sc3b9lZ:~# mc admin group disable src-minio default-group
Disabled group `default-group` successfully.

```

禁用组，组内成员将**无法继承该组策略**。

---

### 6. 启用用户组

```bash
  mc admin group enable TARGET GROUPNAME
  root@iZbp17vix2j58ya7sc3b9lZ:~# mc admin group enable src-minio default-group
Enabled group `default-group` successfully.

```

启用组，组内成员可正常继承绑定的策略。新建的组默认即为启用状态。

---
### 7. 删除用户组

```bash
  mc admin group rm TARGET GROUPNAME
```

删除指定组，不会删除组内的用户。如果需要移除用户，请使用 `mc admin user` 命令。

---

## 三、命令速查表

| 操作       | 命令                       | 说明              |
| -------- | ------------------------ | --------------- |
| 创建组并添加成员 | `mc admin group add`     | 创建不存在的组，并将用户加入组 |
| 查看组详情    | `mc admin group info`    | 显示组成员和策略        |
| 列出所有组    | `mc admin group ls`      | 列出 MinIO 中所有用户组 |
| 删除组      | `mc admin group rm`      | 移除组，不影响组内用户     |
| 启用组      | `mc admin group enable`  | 启用组，使组策略生效      |
| 禁用组      | `mc admin group disable` | 禁用组，暂停组策略继承     |

---

## 四、注意事项

* 用户组名称不能包含 `=` 或 `,` 字符；
* 用户必须在目标部署中存在，才能加入组；
* 创建组后默认不绑定任何策略，需使用以下命令手动绑定：

```bash
mc admin policy attach TARGET --group GROUPNAME --policy POLICYNAME
```

> 详细策略管理请参考官方文档：[MinIO Policy Based Access Control](https://docs.min.io/community/minio-object-store/administration/identity-access-management/policy-based-access-control.html)

---

## 五、相关文档

* [MinIO 用户管理](https://docs.min.io/community/minio-object-store/administration/identity-access-management/minio-user-management.html)
* [MinIO 组管理](https://docs.min.io/community/minio-object-store/administration/identity-access-management/minio-group-management.html)
* [MinIO 策略管理](https://docs.min.io/community/minio-object-store/administration/identity-access-management/policy-based-access-control.html)