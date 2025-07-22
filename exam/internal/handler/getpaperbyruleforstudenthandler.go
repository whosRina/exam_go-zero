package handler

import (
	"net/http"

	"exam-system/exam/internal/logic"
	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getPaperByRuleForStudentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetPaperByRuleForStudentLogic(r.Context(), svcCtx)
		err := l.GetPaperByRuleForStudent()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
