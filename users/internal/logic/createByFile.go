package logic

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/users/internal/types"
	"exam-system/users/model"
	"github.com/xuri/excelize/v2"
	"strings"
	"time"

	"exam-system/users/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateByFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateByFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateByFileLogic {
	return &CreateByFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateByFileLogic) CreateByFile(req *types.CreateByFileRequest, tokenString string) (*types.CreateByFileResponse, error) {
	// 解析JWT获取userId和userType
	_, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}

	// 只允许管理员创建用户
	if userType != 0 {
		return nil, errors.New("无权限创建用户")
	}

	// 解码base64文件
	fileBytes, err := base64.StdEncoding.DecodeString(req.Base64File)
	if err != nil {
		return nil, errors.New("base64解码失败")
	}

	// 打开Excel文件
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, errors.New("excel文件读取失败")
	}
	defer func(f *excelize.File) {
		_ = f.Close()
	}(f)

	// 读取第一个Sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return nil, errors.New("excel内容不足")
	}

	// 判断用户类型：根据表头第一列内容是"学号"或"工号"
	header := strings.TrimSpace(rows[0][0])
	var inferredUserType int64
	if header == "学号" {
		inferredUserType = 2 // 学生
	} else if header == "工号" {
		inferredUserType = 1 // 教师
	} else {
		return nil, errors.New("无法识别的表头，第一列应为“学号”或“工号”")
	}

	var successCount, failCount int
	var failedItems []types.FailedUserItem

	// 遍历数据行（从第二行开始）
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}

		username := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		defaultPassword := username

		salt := GenerateSalt()
		hashedPasswd := HashPassword(defaultPassword, salt)
		_, err := l.svcCtx.UsersModel.Insert(l.ctx, &model.Users{
			Username:   username,
			Name:       name,
			Passwd:     hashedPasswd,
			Salt:       salt,
			Type:       inferredUserType,
			CreateTime: time.Now(),
			IsDelete:   0,
		})
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				failCount++
				failedItems = append(failedItems, types.FailedUserItem{
					Username: username,
					Name:     name,
					Reason:   "用户名已存在",
				})
			} else {
				failCount++
				failedItems = append(failedItems, types.FailedUserItem{
					Username: username,
					Name:     name,
					Reason:   err.Error(),
				})
			}
			logx.Errorf("创建用户失败:%s，错误:%v", username, err)
			continue
		}
		successCount++
	}

	return &types.CreateByFileResponse{
		Message:      "用户导入完成",
		TotalCount:   len(rows) - 1,
		SuccessCount: successCount,
		FailCount:    failCount,
		FailedUsers:  failedItems,
	}, nil
}
