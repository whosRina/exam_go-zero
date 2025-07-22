package handler

import (
	"errors"
	"exam-system/exam/internal/types"
	"net/http"

	"exam-system/exam/internal/logic"
	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPaperForStudentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPaperForStudentRequest
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
		l := logic.NewGetPaperForStudentLogic(r.Context(), svcCtx)
		resp, err := l.GetPaperForStudent(&req, tokenString)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		// 返回JSON响应
		httpx.OkJson(w, resp)
	}
}
