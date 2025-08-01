# 🚀MinIO 用户管理实战指南（附完整Go语言示例代码）

> 本文以实战方式介绍如何通过 MinIO 官方提供的 Go SDK（`madmin`）实现用户管理与权限操作。涵盖从用户创建、权限绑定到删除、禁用等核心功能，并配有详细代码讲解，适合 MinIO 运维人员和后端开发者参考。

---

## 🧩 背景说明

MinIO 是一款高性能、Kubernetes 原生的对象存储系统。为满足多用户、多角色的权限管理需求，MinIO 提供了强大的 IAM（身份访问管理）机制。

在日常管理中，我们常常需要进行如下操作：

* 添加用户
* 绑定权限策略
* 管理用户状态（禁用/启用）
* 删除用户
* 查看用户-权限策略的映射关系

下面通过 `madmin` 客户端实现相关功能。

---

## 1️⃣ 创建用户

使用 `AddUser` 方法，可以添加一个新用户，需要传入用户名和密码。

```go
// AddUser 添加用户
func AddUser(client *madmin.AdminClient, username, passwd string) {
	err := client.AddUser(context.Background(), username, passwd)
	if err != nil {
		slog.Error("添加用户失败!", "msg", err)
	}
}
```

---

## 2️⃣ 列出所有用户

调用 `ListUsers` 可以列出当前 MinIO 服务中的所有用户信息。

```go
// ListUser 列出所有用户
func ListUser(client *madmin.AdminClient) {
	users, err := client.ListUsers(context.Background())
	if err != nil {
		slog.Error("操作用户列表失败")
		return
	}
	indent, err := json.MarshalIndent(users, "", " ")
	slog.Info("获取到的所有列表", "msg", string(indent))
}
```

---

## 3️⃣ 获取所有已创建的权限策略

使用 `ListCannedPolicies` 获取 MinIO 中所有的内置或自定义权限策略。

```go
// ListCannedPolicies 列出所有权限策略
func ListCannedPolicies(client *madmin.AdminClient) {
	policies, err := client.ListCannedPolicies(context.Background())
	if err != nil {
		slog.Error("获取策略失败", "msg", err)
		return
	}
	for key, value := range policies {
		slog.Info("策略信息", "名称", key, "内容", value)
	}
}
```

---

## 4️⃣ 绑定权限策略给用户

使用 `AttachPolicy` 方法，将一个或多个权限策略绑定到用户。

```go
// AttachPolicy 给用户绑定权限策略
func AttachPolicy(client *madmin.AdminClient, username string, policies []string) {
	client.AttachPolicy(context.Background(), madmin.PolicyAssociationReq{
		User:     username,
		Policies: policies,
	})
}
```

---

## 5️⃣ 禁用用户

通过 `SetUserStatus` 方法设置用户状态为禁用。

```go
// DisableUser 禁用用户
func DisableUser(client *madmin.AdminClient, username string) {
	client.SetUserStatus(context.Background(), username, madmin.AccountDisabled)
}
```

> ✅ 注意：禁用用户后，该用户无法继续访问对象存储服务，权限将失效。

---

## 6️⃣ 查看权限策略与用户/组的映射关系

`GetPolicyEntities` 可以查询当前系统中，哪些用户绑定了哪些权限策略。

```go
// GetPolicy 获取对应的权限策略映射
func GetPolicy(client *madmin.AdminClient) {
	policies, err := client.GetPolicyEntities(context.Background(), madmin.PolicyEntitiesQuery{})
	if err != nil {
		slog.Error("获取策略失败", "msg", err)
		return
	}
	marshal, err := json.Marshal(policies)
	if err != nil {
		slog.Error("序列化失败", "msg", err)
	}
	slog.Info("获取到的policy成功", "msg", string(marshal))
}
```

---

## 7️⃣ 删除用户

彻底删除某个用户账户。

```go
// DeleteUser 删除用户
func DeleteUser(client *madmin.AdminClient, username string) {
	err := client.RemoveUser(context.Background(), username)
	if err != nil {
		slog.Error("删除用户失败！", "username", username)
		return
	}
	slog.Info("成功删除用户", "username", username)
}
```

---

## 🔚 总结

本文基于 MinIO 的 Go 客户端 SDK，实现了常见的用户管理与权限配置功能。可作为运维人员或开发者日常管理 MinIO 的实用工具。

> 📌 建议大家在使用时做好权限策略的统一规划，避免权限过大或过小影响系统安全或使用体验。

---

## 📎 推荐阅读

* [MinIO 官方文档 - 用户管理](https://docs.min.io/community/minio-object-store/administration/identity-access-management/minio-user-management.html)
* [MinIO 权限策略详解](https://docs.min.io/community/minio-object-store/administration/identity-access-management/policy-based-access-control.html)

如果你觉得本文对你有帮助，欢迎点赞 👍 收藏 ⭐ 留言 💬，让更多人看到！
