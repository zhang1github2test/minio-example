package madminclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/madmin-go/v4"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log/slog"
)

func InitAdminClient() *madmin.AdminClient {
	// Use a secure connection.
	ssl := false

	// Initialize minio client object.
	mdmClnt, err := madmin.NewWithOptions("121.43.141.218:9000", &madmin.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: ssl,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return mdmClnt
}

// AddUser 添加用户
func AddUser(client *madmin.AdminClient, username, passwd string) {
	err := client.AddUser(context.Background(), username, passwd)
	if err != nil {
		slog.Error("添加用户失败!", "msg", err)
	}
}

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

// 给用户绑定权限策略
func AttachPolicy(client *madmin.AdminClient, username string, policies []string) {
	client.AttachPolicy(context.Background(), madmin.PolicyAssociationReq{
		User:     username,
		Policies: policies,
	})
}

// DisableUser 禁用用户
func DisableUser(client *madmin.AdminClient, username string) {
	client.SetUserStatus(context.Background(), username, madmin.AccountDisabled)
}

func DeleteUser(client *madmin.AdminClient, username string) {
	err := client.RemoveUser(context.Background(), username)
	if err != nil {
		slog.Error("删除用户失败！", "username", username)
		return
	}
	slog.Info("成功删除用户", "username", username)
}

// GetPolicy 获取对应的权限策略映射
func GetPolicy(client *madmin.AdminClient) {
	policies, err := client.GetPolicyEntities(context.Background(), madmin.PolicyEntitiesQuery{})
	if err != nil {
		slog.Error("获取策略失败", "msg", err)
		return
	}
	marshal, err := json.Marshal(policies)
	if err != nil {
		slog.Error("获取到的policy失败", "msg", err)
	}
	slog.Info("获取到的policy成功", "msg", string(marshal))
}

// ListCannedPolicies 列出所有权限策略
func ListCannedPolicies(client *madmin.AdminClient) {
	policies, err := client.ListCannedPolicies(context.Background())
	if err != nil {
	}
	for key, value := range policies {
		slog.Info("msg", key, value)
	}
}

func AddUserGroup(client *madmin.AdminClient, groupsName string, username ...string) {
	err := client.UpdateGroupMembers(context.Background(), madmin.GroupAddRemove{
		Group:   groupsName,
		Members: username,
	})
	if err != nil {
		slog.Error("添加用户组失败", "msg", err)
	}
}
