package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListPapersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPapersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPapersLogic {
	return &ListPapersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPapersLogic) ListPapers(tokenString string) (*types.PaperListResponse, error) {
	// 解析 JWT，获取 userId 和 userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的 JWT")
	}

	// 仅教师可查看试卷
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限查看试卷")
	}

	// 查询试卷列表
	papers, err := l.svcCtx.PaperModel.FindAllByUserId(l.ctx, userId)
	if err != nil {
		return nil, errors.New("查询试卷列表失败")
	}

	// 组装返回数据
	var paperList []types.PaperInfo
	for _, paper := range papers {
		paperList = append(paperList, types.PaperInfo{
			Id:         paper.Id,
			Name:       paper.Name,
			TotalScore: int(paper.TotalScore),
			CreatedBy:  paper.CreatedBy,
			CreatedAt:  paper.CreatedAt,
			UpdatedAt:  paper.UpdatedAt,
		})
	}

	// 返回数据
	return &types.PaperListResponse{
		PaperList: paperList,
	}, nil
}
