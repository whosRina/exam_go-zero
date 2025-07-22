package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/types"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type QuestionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuestionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuestionListLogic {
	return &QuestionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuestionListLogic) QuestionList(req *types.QuestionListRequest, tokenString string) (*types.QuestionListResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 根据用户类型鉴别用户是否有权限查看题目列表
	if userType != 1 { // 只有教师才有权限查看题目列表
		return nil, errors.New("无权限查看题目列表")
	}

	// 检查题库是否存在且该教师是否是该题库的创建者
	bank, err := l.svcCtx.QuestionBankModel.FindOne(l.ctx, int64(req.BankId))
	if err != nil {
		return nil, errors.New("查询题库失败")
	}

	// 确保查询到的题库属于该教师
	if bank.CreatedBy != userId {
		return nil, errors.New("无权限查看此题库")
	}

	// 验证分页参数是否合法
	if req.Page <= 0 {
		req.Page = 1 // 默认从第一页开始
	}

	// 限制size的最大值为500
	if req.Size <= 0 {
		req.Size = 10 // 默认每页10条数据
	}
	if req.Size > 500 {
		req.Size = 500 // 如果请求的size大于500，限制为500
	}

	offset := (req.Page - 1) * req.Size // 计算偏移量

	// 调用封装的分页查询方法
	questions, err := l.svcCtx.QuestionModel.FindQuestionsByBankAndType(l.ctx, req.BankId, req.Type, offset, req.Size)
	if err != nil {
		return nil, err
	}

	// 获取该类型题目的总数
	total, err := l.svcCtx.QuestionModel.CountQuestionsByBankAndType(l.ctx, req.BankId, req.Type)
	if err != nil {
		return nil, err
	}

	// 构建返回结果
	return &types.QuestionListResponse{
		Total:     total,     // 查询到的题目总数（通过COUNT查询得到）
		Page:      req.Page,  // 当前页码，从请求中获取
		Size:      req.Size,  // 每页显示的题目数量，从请求中获取
		Questions: questions, // 当前页的题目列表，数据库查询结果
	}, nil
}
