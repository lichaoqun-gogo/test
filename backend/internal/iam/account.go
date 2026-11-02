package iam

// Account IAM 帐户接口返回的用户模型（登录用户 profile，即 CAS 登录后从 IAM 获取的用户信息）。
// 对应 IAM OpenAPI 的用户信息查询接口响应，用于将 IAM 用户映射到平台用户身份。
// 与平台 users 表的映射（由适配层完成）：AccountCode→username、AccountName→nickname、
// Email→email、DepName→department。
type Account struct {
	AccountCode    string   `json:"accountCode"`    // 帐户（登录账号）
	AccountName    string   `json:"accountName"`    // 帐户名称（姓名）
	Email          string   `json:"email"`          // 邮箱
	DepCode        string   `json:"depCode"`        // 部门编码
	DepName        string   `json:"depName"`        // 部门名称
	OrgCode        string   `json:"orgCode"`        // 组织编码
	CustomOrgCode  string   `json:"customOrgCode"`  // 自定义组织编码
	OrgName        string   `json:"orgName"`        // 组织名称
	OrgDisplayName string   `json:"orgDisplayName"` // 组织机构展示名称
	PartOrgCodes   []string `json:"partOrgCodes"`   // 所属组织编码列表
}
