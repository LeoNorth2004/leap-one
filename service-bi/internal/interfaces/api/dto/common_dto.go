package dto

// UserInfo 用户简要信息（供跨服务使用�?
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
}
