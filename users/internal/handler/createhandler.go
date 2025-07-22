package handler

import (
	"errors"
	"net/http"

	"exam-system/users/internal/logic"
	"exam-system/users/internal/svc"
	"exam-system/users/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func createHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 从请求头中获取JWT，格式通常为 "Bearer <token>"
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httpx.Error(w, errors.New("缺少Authorization"))
			return
		}

		// 调用逻辑层进行用户创建操作
		l := logic.NewCreateLogic(r.Context(), svcCtx)
		resp, err := l.Create(&req, tokenString)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON格式的响应给前端
		httpx.OkJson(w, resp)
	}
}
