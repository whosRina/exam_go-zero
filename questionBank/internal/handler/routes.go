// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package handler

import (
	"net/http"
	"time"

	"exam-system/questionBank/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/bankDetail",
				Handler: bankDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/bankList",
				Handler: bankListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createBank",
				Handler: createBankHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/createQuestion",
				Handler: createQuestionHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/deleteBank",
				Handler: deleteBankHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/deleteQuestion",
				Handler: deleteQuestionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/questionDetail",
				Handler: questionDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/questionList",
				Handler: questionListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/updateBank",
				Handler: updateBankHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/updateQuestion",
				Handler: updateQuestionHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/parseWordQuestions",
				Handler: ParseWordQuestionshandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/questions"),
		rest.WithTimeout(3000*time.Millisecond),
	)
}
