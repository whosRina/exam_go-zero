package handler

import (
	"exam-system/exam/internal/logic"
	"exam-system/exam/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func setExamQuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSetExamQuestionsLogic(r.Context(), svcCtx)
		err := l.SetExamQuestions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
