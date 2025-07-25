// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package types

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Token string `json:"token"`
}

// CreateRequest 结构体（POST 传入）
type CreateRequest struct {
	Users []struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
		Type     int    `json:"type"`
	} `json:"users"`
}

// CreateResponse 定义创建用户后的响应结构体
type CreateResponse struct {
	Message string `json:"message"`
}

type UserListResponse struct {
	UserList []UserInfo `json:"userList"`
}

type UserInfo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Type     int    `json:"type"`
}
type UpdateUserRequest struct {
	Id       int64  `json:"id"`                 // 用户ID
	Name     string `json:"name"`               // 姓名
	Username string `json:"username"`           // 用户名
	Password string `json:"password,omitempty"` // 新密码（可选）
	Type     int    `json:"type"`               // 用户类型
}
type UpdateUserResponse struct {
	Message string `json:"message"`
}

type DeleteRequest struct {
	Id int `json:"id"`
}
type DeleteResponse struct {
	Message string `json:"message"`
}

type CreateByFileRequest struct {
	Base64File string `json:"base64File"` // 前端传的base64编码Excel文件
}
type FailedUserItem struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Reason   string `json:"reason"`
}

type CreateByFileResponse struct {
	Message      string           `json:"message"`
	TotalCount   int              `json:"totalCount"`
	SuccessCount int              `json:"successCount"`
	FailCount    int              `json:"failCount"`
	FailedUsers  []FailedUserItem `json:"failedUsers"`
}

type UpdatePwdRequest struct {
	OldPassword string `json:"oldPassword"` // 原密码
	NewPassword string `json:"newPassword"` // 新密码
}
type UpdatePwdResponse struct {
	Message string `json:"message"`
}

type UserDetailResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	UserType int64  `json:"userType"`
}
