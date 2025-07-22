package handler

import (
	"errors"
	"exam-system/questionBank/internal/types"
	"net/http"

	"exam-system/questionBank/internal/logic"
	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func updateQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateQuestionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 从请求头中获取 JWT
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httpx.Error(w, errors.New("缺少 Authorization"))
			return
		}

		// 调用逻辑层创建题库
		l := logic.NewUpdateQuestionLogic(r.Context(), svcCtx)
		resp, err := l.UpdateQuestion(&req, tokenString)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回 JSON 响应
		httpx.OkJson(w, resp)
	}
}
