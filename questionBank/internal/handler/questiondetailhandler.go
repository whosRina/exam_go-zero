package handler

import (
	"net/http"

	"exam-system/questionBank/internal/logic"
	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func questionDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewQuestionDetailLogic(r.Context(), svcCtx)
		err := l.QuestionDetail()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
