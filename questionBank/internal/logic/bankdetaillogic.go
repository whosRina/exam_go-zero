package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/types"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BankDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBankDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BankDetailLogic {
	return &BankDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BankDetailLogic) BankDetail(req *types.BankDetailRequest, tokenString string) (*types.BankDetailResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 根据用户类型鉴别用户是否有权限查看题库
	if userType != 1 { // 只有教师才有权限查看题库详情
		return nil, errors.New("无权限查看题库")
	}

	// 检查题库是否存在且该教师是否是该题库的创建者
	bank, err := l.svcCtx.QuestionBankModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New("查询题库详情失败")
	}

	// 确保查询到的题库属于该教师
	if bank.CreatedBy != userId {
		return nil, errors.New("无权限查看此题库")
	}

	// 调用数据库模型中的GetBankDetailAndTypeCounts查询函数
	bankDetailResponse, err := l.svcCtx.QuestionBankModel.GetBankDetailAndTypeCounts(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New("获取题库详情和题目类型统计信息失败")
	}

	return bankDetailResponse, nil
}
