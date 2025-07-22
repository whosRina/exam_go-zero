package handler

import (
	"errors"
	"exam-system/users/internal/types"
	"net/http"

	"exam-system/users/internal/logic"
	"exam-system/users/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func deleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		// 从请求头中获取JWT，格式通常为"Bearer <token>"
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httpx.Error(w, errors.New("缺少 Authorization"))
			return
		}

		// 创建删除逻辑
		l := logic.NewDeleteLogic(r.Context(), svcCtx)
		resp, err := l.Delete(&req, tokenString) // 传递ID给删除逻辑
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON格式的响应给前端
		httpx.OkJson(w, resp)
	}
}
