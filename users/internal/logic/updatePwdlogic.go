package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/users/internal/types"
	"fmt"

	"exam-system/users/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePwdLogic {
	return &UpdatePwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePwdLogic) UpdatePwd(req *types.UpdatePwdRequest, tokenString string) (*types.UpdatePwdResponse, error) {
	// 解析JWT获取userId和userType
	userId, _, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的token")
	}

	// 查询用户信息
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, userId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 校验原密码是否正确
	inputOldHashed := HashPassword(req.OldPassword, user.Salt)
	if inputOldHashed != user.Passwd {
		return nil, errors.New("原密码错误")
	}

	// 生成新密码hash和盐
	newSalt := GenerateSalt()
	newHashed := HashPassword(req.NewPassword, newSalt)

	// 更新用户密码
	user.Passwd = newHashed
	user.Salt = newSalt

	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("密码更新失败:%v", err)
	}

	return &types.UpdatePwdResponse{
		Message: "密码修改成功",
	}, nil
}
