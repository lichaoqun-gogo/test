package iam

// ResourceType IAM 资源类型（对应 Java 枚举 ResourceType 的 JSON 字符串值）
type ResourceType string

const (
	ResourceTypeData   ResourceType = "DATA"   // 数据
	ResourceTypeMenu   ResourceType = "MENU"   // 菜单
	ResourceTypeAPI    ResourceType = "API"    // API
	// ResourceTypeButton 资源类型「按钮」：源枚举注释被截断（BUT...），名称待与 IAM 确认后补充
	ResourceTypeButton ResourceType = "BUTTON"
)

// Action 权限操作（资源上可执行的动作）
type Action struct {
	Name string `json:"name"` // 操作名称，如 view/edit
	Desc string `json:"desc"` // 操作描述
}

// Condition 资源访问条件（如按组织/部门/角色等约束资源可见范围）
type Condition struct {
	Name     string `json:"name"`     // 条件名称
	Operator string `json:"operator"` // 操作类型，如 in/eq/contains
	Value    string `json:"value"`    // 条件值
}

// Permission IAM 权限接口返回的权限模型：以「资源 + 操作 + 条件」表达细粒度权限。
// 映射到平台权限档案：MENU 类型决定菜单可见性（对应 GetUserMenus），
// API 类型决定接口访问控制，BUTTON/DATA 类型决定按钮与数据范围。
// 适配层需将 IAM Permission 集合翻译为平台权限点编码（permission.Profile.Permissions）。
type Permission struct {
	ID           int64        `json:"id"`           // 资源 ID
	Resource     string       `json:"resource"`     // 资源名称
	ResourceType ResourceType `json:"resourceType"` // 资源类型：DATA/MENU/API/BUTTON
	Visible      bool         `json:"visible"`      // 是否可见
	Actions      []Action     `json:"actions"`      // 操作列表
	Conditions   []Condition  `json:"conditions"`   // 访问条件
}
