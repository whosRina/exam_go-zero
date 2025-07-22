package handler

import (
	"net/http"

	"exam-system/users/internal/logic"
	"exam-system/users/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func loginOutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewLoginOutLogic(r.Context(), svcCtx)
		err := l.LoginOut()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
