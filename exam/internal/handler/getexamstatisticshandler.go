package handler

import (
	"exam-system/exam/internal/logic"
	"exam-system/exam/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getExamStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetExamStatisticsLogic(r.Context(), svcCtx)
		err := l.GetExamStatistics()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
