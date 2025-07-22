package handler

import (
	"exam-system/exam/internal/logic"
	"exam-system/exam/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getExamAttemptsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetExamAttemptsLogic(r.Context(), svcCtx)
		err := l.GetExamAttempts()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
