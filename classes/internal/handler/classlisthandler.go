package handler

import (
	"errors"
	"net/http"

	"exam-system/classes/internal/logic"
	"exam-system/classes/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func classListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取JWT，格式通常为 "Bearer <token>"
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httpx.Error(w, errors.New("缺少Authorization"))
			return
		}

		// 调用逻辑层进行用户创建操作
		l := logic.NewClassListLogic(r.Context(), svcCtx)
		resp, err := l.ClassList(tokenString)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON格式的响应给前端
		httpx.OkJson(w, resp)
	}
}
