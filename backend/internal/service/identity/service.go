package identity

import (
	"context"

	"osm-mngt/internal/iam"

	"gorm.io/gorm"
)

// Service 用户身份服务：提供登录用户的「全量权限资源包」（IAMUser，account + 组织等权限资源信息）。
// 数据源可切换：LocalService（users 表，IAM 接入前的临时实现）/ OpenAPIService（IAM 开放接口）。
// handler 只依赖本接口、不感知数据源，便于后期剥离本地存储；后续权限（roles/menus/buttons）接入时
// 仅扩展本接口实现，不影响 handler 与 DTO。
type Service interface {
	// GetUserInfo 获取用户全量权限资源包（account + 组织等权限资源信息）
	GetUserInfo(ctx context.Context, account *iam.Account) (*iam.IAMUser, error)
}

// LocalService 本地实现：从静态用户数据补齐 account 身份信息，组织信息暂不提供。
// 对应「没有权限接口时先落本地数据」的临时方案，IAM 接入后整体替换为 OpenAPIService。
type LocalService struct {
	db *gorm.DB
}

func NewLocalService(db *gorm.DB) *LocalService {
	return &LocalService{db: db}
}

// staticLocalUsers 静态用户身份数据（users 表已废弃，替代兜底数据源）。
// 与 mock CAS（users.json）保持一致，保证登录用户（zhangsan 等）能补齐姓名/邮箱/部门。
var staticLocalUsers = []struct {
	Username   string
	Nickname   string
	Email      string
	Department string
	Team       string
}{
	{Username: "zhangsan", Nickname: "张三", Email: "zhangsan@qim.local", Department: "技术部", Team: "前端组"},
	{Username: "lisi", Nickname: "李四", Email: "lisi@qim.local", Department: "技术部", Team: "前端组"},
	{Username: "zhengjie", Nickname: "郑杰", Email: "zhengjie@qim.local", Department: "技术部", Team: "前端组"},
	{Username: "wangwu", Nickname: "王五", Email: "wangwu@qim.local", Department: "平台研发部", Team: "后端组"},
	{Username: "liuyang", Nickname: "刘洋", Email: "liuyang@qim.local", Department: "平台研发部", Team: "后端组"},
	{Username: "zhaoliu", Nickname: "赵六", Email: "zhaoliu@qim.local", Department: "运维部", Team: "后端组"},
	{Username: "chenming", Nickname: "陈明", Email: "chenming@qim.local", Department: "运维部", Team: "后端组"},
	{Username: "sunqi", Nickname: "孙七", Email: "sunqi@qim.local", Department: "安全部", Team: "设计组"},
	{Username: "qianqi", Nickname: "钱七", Email: "qianqi@qim.local", Department: "安全部", Team: "设计组"},
	{Username: "zhouba", Nickname: "周八", Email: "zhouba@qim.local", Department: "开源治理办公室", Team: "运营组"},
	{Username: "wujiu", Nickname: "吴九", Email: "wujiu@qim.local", Department: "开源治理办公室", Team: "运营组"},
}

func (s *LocalService) GetUserInfo(ctx context.Context, account *iam.Account) (*iam.IAMUser, error) {
	if account == nil {
		account = &iam.Account{}
	}
	// 从静态用户数据补齐 account 缺失的身份信息（姓名/邮箱/部门）
	if account.AccountCode != "" {
		for _, u := range staticLocalUsers {
			if u.Username == account.AccountCode {
				if account.AccountName == "" {
					account.AccountName = u.Nickname
				}
				if account.Email == "" {
					account.Email = u.Email
				}
				if account.DepName == "" {
					account.DepName = u.Department
				}
				if account.OrgName == "" {
					account.OrgName = u.Team
				}
				break
			}
		}
	}
	return &iam.IAMUser{Account: account}, nil
}

// OpenAPIService IAM 实现：用户全量权限包 = account（调用方传入，源自 CAS/JWT 身份）+ 组织信息（IAM 接口查询）拼接。
// roles/menus/buttons 待 IAM 权限/菜单接口接入后在此扩展。
type OpenAPIService struct {
	cli *iam.Client
}

func NewOpenAPIService(cli *iam.Client) *OpenAPIService {
	return &OpenAPIService{cli: cli}
}

func (s *OpenAPIService) GetUserInfo(ctx context.Context, account *iam.Account) (*iam.IAMUser, error) {
	if account == nil {
		account = &iam.Account{}
	}
	orgs, err := s.cli.GetOrgs(ctx, "", "")
	if err != nil {
		return nil, err
	}
	return &iam.IAMUser{Account: account, Orgs: orgs}, nil
}
