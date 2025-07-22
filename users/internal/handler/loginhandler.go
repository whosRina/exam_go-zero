package handler

import (
	"net/http"

	"exam-system/users/internal/logic"
	"exam-system/users/internal/svc"
	"exam-system/users/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req) // 正确解构返回值
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 返回正确的响应
		httpx.OkJson(w, resp)
	}
}
