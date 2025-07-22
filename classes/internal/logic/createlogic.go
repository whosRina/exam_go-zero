package logic

import (
	"context"
	"crypto/rand"
	"errors"
	"exam-system/JWT"
	"exam-system/classes/internal/svc"
	"exam-system/classes/internal/types"
	"exam-system/classes/model"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"strings"
)

type CreateClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateClassLogic {
	return &CreateClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateClassLogic) CreateClass(req *types.CreateClassRequest, tokenString string) (*types.CreateClassResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的 JWT")
	}

	// 确保用户是教师才可创建班级
	if userType != 1 { // 1 代表教师
		return nil, errors.New("无权限创建班级")
	}

	// 校验班级名称和是否允许加入
	if req.ClassName == "" {
		return nil, errors.New("创建失败，班级名称或班级加入状态为空")
	}

	// 检查Is_joinable是否为有效值
	if req.IsJoinable != 0 && req.IsJoinable != 1 {
		return nil, errors.New("班级加入状态无效，只能为0或1")
	}

	// 生成唯一邀请码
	classCode, err := l.generateUniqueClassCode()
	if err != nil {
		logx.Error("生成邀请码失败:", err)
		return nil, errors.New("生成邀请码失败")
	}
	// 创建班级
	class := &model.Class{
		Name:       req.ClassName,
		TeacherId:  userId,
		ClassCode:  classCode,
		IsJoinable: int64(req.IsJoinable), // 默认允许加入
	}

	_, err = l.svcCtx.ClassModel.Insert(l.ctx, class)
	if err != nil {
		return nil, errors.New("创建班级失败")
	}

	return &types.CreateClassResponse{
		Message:   "班级创建成功",
		ClassCode: classCode,
	}, nil
}

// Base58 字符集（去掉 0OIl 避免混淆）
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// generateUniqueClassCode 生成 6 位唯一邀请码
func (l *CreateClassLogic) generateUniqueClassCode() (string, error) {
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

// randomBase58String 生成指定长度的Base58随机字符串
func randomBase58String(length int) (string, error) {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(base58Alphabet))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(base58Alphabet[index.Int64()])
	}
	return sb.String(), nil
}
