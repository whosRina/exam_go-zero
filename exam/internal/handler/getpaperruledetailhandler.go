package handler

import (
	"net/http"

	"exam-system/exam/internal/logic"
	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getPaperRuleDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetPaperRuleDetailLogic(r.Context(), svcCtx)
		err := l.GetPaperRuleDetail()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
