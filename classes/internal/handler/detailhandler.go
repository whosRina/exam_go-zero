package handler

import (
	"errors"
	"exam-system/classes/internal/types"
	"net/http"

	"exam-system/classes/internal/logic"
	"exam-system/classes/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func detailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClassDetailRequest
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

		// 创建删除逻辑
		l := logic.NewDetailLogic(r.Context(), svcCtx)
		resp, err := l.Detail(&req, tokenString) // 传递 ID 给删除逻辑
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON格式的响应给前端
		httpx.OkJson(w, resp)
	}
}
