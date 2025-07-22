package handler

import (
	"errors"
	"exam-system/questionBank/internal/types"
	"net/http"

	"exam-system/questionBank/internal/logic"
	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func questionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionListRequest
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
		l := logic.NewQuestionListLogic(r.Context(), svcCtx)
		resp, err := l.QuestionList(&req, tokenString)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON响应
		httpx.OkJson(w, resp)
	}
}
