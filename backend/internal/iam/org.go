package iam

// Org IAM 组织接口返回的组织模型（树形节点，Children 表达下级组织）。
// 对应 IAM OpenAPI 的组织查询接口响应，用于将外部组织体系映射到平台「部门」等组织维度。
// 该 DTO 仅描述 IAM 侧数据结构，与平台本地存储解耦；后续切换数据源时仅替换 service 实现。
type Org struct {
	AppCode       string `json:"appCode"`       // 应用编号
	AppOrg        string `json:"appOrg"`        // 应用组织
	OrgCode       string `json:"orgCode"`       // 组织编码（IAM 侧稳定标识）
	OrgName       string `json:"orgName"`       // 组织名称
	OrgPrincipal  string `json:"orgPrincipal"`  // 组织负责人
	ParentOrgCode string `json:"parentOrgCode"` // 上级组织编码
	ParentOrgName string `json:"parentOrgName"` // 上级组织名称
	OrgType       string `json:"orgType"`       // 组织类型
	OrgLevel      string `json:"orgLevel"`      // 组织层级
	DepCode       string `json:"depCode"`       // 部门编码
	DepName       string `json:"depName"`       // 部门名称
	Children      []*Org `json:"children"`      // 下级组织
}
