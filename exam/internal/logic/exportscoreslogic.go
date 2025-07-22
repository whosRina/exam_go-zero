package logic

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	_ "errors"
	jwtutil "exam-system/JWT"
	"exam-system/exam/internal/svc"
	"exam-system/exam/internal/types"
	"fmt"
	"html"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ExportScoresLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportScoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportScoresLogic {
	return &ExportScoresLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportScoresLogic) ExportScores(req *types.ExportScoresRequest, tokenString string) (*types.ExportScoresResponse, error) {
	// 解析JWT获取userId和userType
	teacherId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 1 {
		return nil, errors.New("只有教师可以导出成绩")
	}

	// 检查考试是否存在，并验证权限
	exam, err := l.svcCtx.ExamModel.FindOne(l.ctx, req.ExamId)
	if err != nil {
		return nil, errors.New("考试不存在")
	}
	if exam.TeacherId != teacherId {
		return nil, errors.New("你不是该考试的创建者，无权限导出成绩")
	}

	// 查询考试成绩
	attempts, err := l.svcCtx.ExamAttemptModel.FindAttemptByExam(l.ctx, exam.Id)
	if err != nil {
		return nil, errors.New("查询考试记录失败")
	}

	// 4使用bytes.Buffer生成CSV，不创建本地文件
	var buf bytes.Buffer

	// 解决Excel乱码问题，写入UTF-8 BOM
	buf.WriteString("\xEF\xBB\xBF")

	// 创建CSV写入器
	writer := csv.NewWriter(&buf)

	// 写入表头
	headers := []string{"学号", "姓名", "分数"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("写入表头失败: %v", err)
	}

	// 写入数据行
	for _, attempt := range attempts {
		student, err := l.svcCtx.UsersModel.FindOne(l.ctx, attempt.StudentId)
		if err != nil {
			l.Logger.Errorf("查询学生信息失败，studentId=%d: %v", attempt.StudentId, err)
			continue
		}

		// 处理HTML特殊字符，避免CSV解析异常
		username := html.EscapeString(student.Username)
		name := html.EscapeString(student.Name)
		score := fmt.Sprintf("%d", attempt.Score)

		row := []string{username, name, score}
		if err := writer.Write(row); err != nil {
			l.Logger.Errorf("写入CSV失败: %v", err)
			continue
		}
	}
	writer.Flush()

	// 生成安全的文件名（避免非法字符）
	safeFileName := sanitizeFileName(fmt.Sprintf("成绩_%s_%d.csv", exam.Name, time.Now().Unix()))

	// 返回CSV二进制数据
	return &types.ExportScoresResponse{
		Message:    "成绩导出成功",
		FileName:   safeFileName,
		FileStream: buf.Bytes(),
	}, nil
}

// sanitizeFileName清理文件名中的非法字符
func sanitizeFileName(name string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}
	return name
}
