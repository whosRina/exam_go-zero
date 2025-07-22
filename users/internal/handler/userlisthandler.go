package handler

import (
	"errors"
	"net/http"

	"exam-system/users/internal/logic"
	"exam-system/users/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func userListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取JWT，格式通常为"Bearer <token>"
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httpx.Error(w, errors.New("缺少 Authorization"))
			return
		}

		// 调用逻辑层进行用户创建操作
		l := logic.NewUserListLogic(r.Context(), svcCtx)
		resp, err := l.UserList(tokenString)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON格式的响应给前端
		httpx.OkJson(w, resp)
	}
}
