package handler

import (
	"errors"
	"exam-system/questionBank/internal/logic"
	"exam-system/questionBank/internal/svc"
	"exam-system/questionBank/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func ParseWordQuestionshandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ParseWordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 从请求头中获取JWT
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httpx.Error(w, errors.New("缺少Authorization"))
			return
		}

		// 调用逻辑层创建题库
		l := logic.NewParseWordQuestionsLogic(r.Context(), svcCtx)
		resp, err := l.ParseAndCreateQuestions(&req, tokenString)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON响应
		httpx.OkJson(w, resp)
	}
}
