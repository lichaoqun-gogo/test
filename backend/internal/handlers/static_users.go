package handlers

// ==================== 静态用户数据（替代 users 表兜底） ====================
// users 表已废弃，人员/部门/组织等兜底数据统一由此静态数据提供，
// 与 mock CAS（users.json）保持一致，保证登录用户（zhangsan 等）能在人员列表中找到。
// 对接真实权限系统后，这些接口的数据源切换为权限系统透传（见各 handler 优先分支），
// 本文件保留作为未对接时的开发/演示数据。
//
// 注意：角色（role）不在此定义——登录用户的角色由 CAS 返回的 role 字段决定
// （见 auth.go CasCallback：role -> JWT -> /api/user/info 的 roles），
// 此处的用户仅用于部门/团队/人员列表等参考数据的兜底展示。

// StaticUser 静态用户（对齐 users 表结构）
type StaticUser struct {
	ID         uint   // 自增主键（模拟表自增）
	Username   string // 登录账号
	Nickname   string // 姓名
	Department string // 部门
	Team       string // 团队
	Email      string // 邮箱
}

// staticUsers 静态用户列表（与 mock CAS users.json 一致）
var staticUsers = []StaticUser{
	{ID: 1, Username: "zhangsan", Nickname: "张三", Department: "技术部", Team: "前端组", Email: "zhangsan@qim.local"},
	{ID: 2, Username: "lisi", Nickname: "李四", Department: "技术部", Team: "前端组", Email: "lisi@qim.local"},
	{ID: 3, Username: "zhengjie", Nickname: "郑杰", Department: "技术部", Team: "前端组", Email: "zhengjie@qim.local"},
	{ID: 4, Username: "wangwu", Nickname: "王五", Department: "平台研发部", Team: "后端组", Email: "wangwu@qim.local"},
	{ID: 5, Username: "liuyang", Nickname: "刘洋", Department: "平台研发部", Team: "后端组", Email: "liuyang@qim.local"},
	{ID: 6, Username: "zhaoliu", Nickname: "赵六", Department: "运维部", Team: "后端组", Email: "zhaoliu@qim.local"},
	{ID: 7, Username: "chenming", Nickname: "陈明", Department: "运维部", Team: "后端组", Email: "chenming@qim.local"},
	{ID: 8, Username: "sunqi", Nickname: "孙七", Department: "安全部", Team: "设计组", Email: "sunqi@qim.local"},
	{ID: 9, Username: "qianqi", Nickname: "钱七", Department: "安全部", Team: "设计组", Email: "qianqi@qim.local"},
	{ID: 10, Username: "zhouba", Nickname: "周八", Department: "开源治理办公室", Team: "运营组", Email: "zhouba@qim.local"},
	{ID: 11, Username: "wujiu", Nickname: "吴九", Department: "开源治理办公室", Team: "运营组", Email: "wujiu@qim.local"},
}

// staticUserByUsername 按登录账号查静态用户
func staticUserByUsername(username string) (StaticUser, bool) {
	for _, u := range staticUsers {
		if u.Username == username {
			return u, true
		}
	}
	return StaticUser{}, false
}
