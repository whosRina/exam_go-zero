package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/types"

	"exam-system/questionBank/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type BankListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBankListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BankListLogic {
	return &BankListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BankListLogic) BankList(tokenString string) (*types.QuestionBankListResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 检查是否为管理员或者教师
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限查看题库")
	}

	// 查询所有题库
	banks, err := l.svcCtx.QuestionBankModel.FindAllById(l.ctx, userId)
	if err != nil {
		return nil, errors.New("查询题库列表失败")
	}

	// 构造返回数据
	var bankList []types.QuestionBankInfo
	for _, bank := range banks {
		bankList = append(bankList, types.QuestionBankInfo{
			Id:        int(bank.Id),
			Name:      bank.Name,
			CreatedBy: int(bank.CreatedBy),
			CreatedAt: bank.CreatedAt,
		})
	}

	// 返回响应
	return &types.QuestionBankListResponse{
		BankList: bankList,
	}, nil
}
