package logic

import (
	"context"
	"database/sql"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/users/internal/types"
	"time"

	"exam-system/users/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteRequest, tokenString string) (*types.DeleteResponse, error) {
	// 解析JWT获取userId和userType
	_, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 检查是否为管理员
	if userType != 0 { // 0 代表管理员
		return nil, errors.New("无权限删除用户")
	}

	// 查询用户是否存在
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, int64(req.Id))
	if err != nil {
		return nil, errors.New("删除用户失败，用户不存在或已删除")
	}

	// 更新用户的删除标志和删除时间
	user.IsDelete = 1
	user.DeleteTime = sql.NullTime{
		Time:  time.Now(), // 设置当前时间
		Valid: true,       // 标记这个时间是有效的
	}

	// 更新数据库（逻辑删除）
	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, errors.New("删除用户失败")

	}
	// 成功删除用户
	return &types.DeleteResponse{
		Message: "用户删除成功",
	}, nil
}
