package logic

import (
	"context"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/classes/internal/svc"
	"exam-system/classes/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateClassLogic {
	return &UpdateClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateClassLogic) UpdateClass(req *types.UpdateClassRequest, tokenString string) (*types.UpdateClassResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 确保用户是教师才可更新班级
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限更新班级")
	}

	// 校验班级ID和班级名称
	if req.ClassID == 0 || req.ClassName == "" {
		return nil, errors.New("更新失败，班级ID或班级名称为空")
	}

	// 检查Is_joinable是否为有效值
	if req.IsJoinable != 0 && req.IsJoinable != 1 {
		return nil, errors.New("班级加入状态无效，只能为0或1")
	}

	// 获取班级信息，检查班级是否存在
	class, err := l.svcCtx.ClassModel.FindOne(l.ctx, req.ClassID)
	if err != nil {
		logx.Error("查询班级信息失败:", err)
		return nil, errors.New("班级不存在")
	}

	// 确保更新的是当前用户创建的班级
	if class.TeacherId != userId {
		return nil, errors.New("无权限更新该班级")
	}

	// 更新班级名称和是否允许加入
	class.Name = req.ClassName
	class.IsJoinable = int64(req.IsJoinable)

	// 如果需要刷新邀请码
	if req.RefreshCode {
		// 生成新的唯一邀请码
		newClassCode, err := l.generateUniqueClassCode()
		if err != nil {
			logx.Error("生成邀请码失败:",err)
			return nil, errors.New("生成邀请码失败")
		}

		// 更新班级邀请码
		class.ClassCode = newClassCode
	}

	// 更新班级信息到数据库
	err = l.svcCtx.ClassModel.Update(l.ctx, class)
	if err != nil {
		logx.Error("更新班级信息失败:",err)
		return nil, errors.New("更新班级失败")
	}

	return &types.UpdateClassResponse{
		Message:   "班级更新成功",
		ClassCode: class.ClassCode,
	}, nil
}

// generateUniqueClassCode 生成6位唯一邀请码
func (l *UpdateClassLogic) generateUniqueClassCode() (string, error) {
	for i := 0; i < 5; i++ { // 最多尝试 5 次避免重复
		code, err := randomBase58String(8)
		if err != nil {
			return "", err
		}

		// 检查数据库中是否已存在
		exist, _ := l.svcCtx.ClassModel.FindByClassCode(l.ctx, code)
		if exist == nil { // 不存在，返回该邀请码
			return code, nil
		}
	}

	return "", errors.New("生成唯一邀请码失败")
}
