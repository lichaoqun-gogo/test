package identity

import (
	"context"

	"osm-mngt/internal/iam"
	"osm-mngt/internal/models"

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

// LocalService 本地实现：从 users 表构造/补齐 account，组织信息暂不提供。
// 对应「没有权限接口时先落 users 表」的临时方案，IAM 接入后整体替换为 OpenAPIService。
type LocalService struct {
	db *gorm.DB
}

func NewLocalService(db *gorm.DB) *LocalService {
	return &LocalService{db: db}
}

func (s *LocalService) GetUserInfo(ctx context.Context, account *iam.Account) (*iam.IAMUser, error) {
	if account == nil {
		account = &iam.Account{}
	}
	// 从 users 表补齐 account 缺失的身份信息（姓名/邮箱/部门）
	if account.AccountCode != "" {
		var user models.User
		if err := s.db.Where("username = ?", account.AccountCode).First(&user).Error; err == nil {
			if account.AccountName == "" {
				account.AccountName = user.Nickname
			}
			if account.Email == "" {
				account.Email = user.Email
			}
			if account.DepName == "" {
				account.DepName = user.Department
			}
			if account.OrgName == "" {
				account.OrgName = user.Team
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
