package iam

// Menu IAM 权限系统菜单接口返回的菜单模型（树形节点，Children 表达子菜单）。
// 对应 IAM OpenAPI 的菜单查询接口响应，用于将外部权限系统的菜单映射到平台菜单。
// 映射说明：path/icon/parentId/orderNum 与平台 menu 表直接对应；visible 与平台 hidden 取反；
// appCode/appModule 为 IAM 应用维度，平台为单一应用可忽略；
// 平台 menu 表独有的 component（前端组件路径）/group（分组）IAM 菜单不提供，需在适配层本地补齐。
type Menu struct {
	ID        int64   `json:"id"`        // 菜单 ID
	Name      string  `json:"name"`      // 菜单名称
	AppCode   string  `json:"appCode"`   // 所属应用
	AppModule string  `json:"appModule"` // 所属模块
	ParentID  int64   `json:"parentId"`  // 父节点 ID，根节点为 0
	OrderNum  int     `json:"orderNum"`  // 同级排序号
	Path      string  `json:"path"`      // 访问地址
	Icon      string  `json:"icon"`      // 显示图标
	Visible   bool    `json:"visible"`   // 是否可见
	Children  []*Menu `json:"children"`  // 子菜单列表
}
