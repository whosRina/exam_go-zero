package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/types"

	"exam-system/exam/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePaperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePaperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePaperLogic {
	return &DeletePaperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePaperLogic) DeletePaper(req *types.DeletePaperRequest, tokenString string) (*types.DeletePaperResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只有教师才能删除试卷
	if userType != 1 {
		return nil, errors.New("无权限删除试卷")
	}

	// 查询试卷是否存在
	paper, err := l.svcCtx.PaperModel.FindOne(l.ctx, int64(req.PaperID))
	if err != nil {
		return nil, errors.New("试卷不存在")
	}

	// 确保只有试卷创建者可以删除
	if paper.CreatedBy != userId {
		return nil, errors.New("无权限删除他人创建的试卷")
	}

	// 执行删除操作
	err = l.svcCtx.PaperModel.Delete(l.ctx, int64(req.PaperID))
	if err != nil {
		return nil, errors.New("删除试卷失败")
	}

	return &types.DeletePaperResponse{
		Message: "试卷删除成功",
	}, nil
}
