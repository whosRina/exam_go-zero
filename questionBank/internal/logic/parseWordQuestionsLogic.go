package logic

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	jwtutil "exam-system/JWT"
	"exam-system/questionBank/internal/svc"
	"exam-system/questionBank/internal/types"
	"exam-system/questionBank/model"
	"fmt"
	"github.com/nguyenthenguyen/docx"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
	"strconv"
	"strings"
)

type ParseWordQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParseWordQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseWordQuestionsLogic {
	return &ParseWordQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Option struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

func (l *ParseWordQuestionsLogic) ParseAndCreateQuestions(req *types.ParseWordRequest, tokenString string) (*types.ParseWordResponse, error) {
	// 解析JWT获取userId和userType
	userId, userType, err := jwtutil.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("无效的JWT")
	}
	if userType != 1 {
		return nil, errors.New("无权限操作")
	}

	// 检查题库
	if _, err := l.svcCtx.QuestionBankModel.FindById(l.ctx, int(req.BankId)); err != nil {
		return nil, errors.New("题库不存在")
	}

	// 解码Base64文件内容
	content, err := base64.StdEncoding.DecodeString(req.FileBase64)
	if err != nil {
		return nil, errors.New("base64解码失败")
	}

	// 使用docx库从内存读取文件
	reader := bytes.NewReader(content)
	doc, err := docx.ReadDocxFromMemory(reader, int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("docx解析失败:%v", err)
	}
	defer func(doc *docx.ReplaceDocx) {
		err := doc.Close()
		if err != nil {

		}
	}(doc)

	// 获取Word文档原始内容，并进行清洗
	rawContent := doc.Editable().GetContent()
	cleanedContent := cleanContent(rawContent)
	l.Logger.Infof("清洗后的文本内容：\n%s", cleanedContent)

	// 按换行符拆分成行
	lines := strings.Split(cleanedContent, "\n")
	l.Logger.Infof("共分割出%d行", len(lines))

	// 解析题目
	questions := l.parseQuestions(lines)

	var successCount, failCount int
	var failedItems []types.FailedItem

	for i, q := range questions {
		q.CreatedBy = userId
		q.BankId = req.BankId
		_, err := l.svcCtx.QuestionModel.Insert(l.ctx, q)
		if err != nil {
			failCount++
			// 获取题目内容的前 50 个字符
			contentSnippet := q.Content
			if len(contentSnippet) > 50 {
				contentSnippet = contentSnippet[:50] + "..."
			}
			// 记录失败项
			failedItems = append(failedItems, struct {
				Index   int    `json:"index"`
				Content string `json:"content"`
			}{
				Index:   i + 1, // 题目索引从 1 开始计数
				Content: contentSnippet,
			})
			continue
		}
		successCount++
	}

	return &types.ParseWordResponse{
		Message:      "题目导入完成",
		TotalCount:   len(questions),
		SuccessCount: successCount,
		FailCount:    failCount,
		FailedItems:  failedItems,
	}, nil
}

// cleanContent 清洗文本，确保题干、选项、答案与题型标头正确分行
func cleanContent(content string) string {
	// 去除XML标签
	content = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(content, "")
	// 统一换行符
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	// 替换中文顿号和全角句号为英文句点
	content = strings.ReplaceAll(content, "．", ".")
	content = regexp.MustCompile(`([A-D])、`).ReplaceAllString(content, "$1.")
	// 数字题号中的顿号，如 1、2、3，仅当它们在行首或换行符后时才替换
	content = regexp.MustCompile(`(\d+)、`).ReplaceAllString(content, "$1.")
	// 在数字题号前插入换行符，例如 "1." 前
	content = regexp.MustCompile(`(?m)(\d+\.)`).ReplaceAllString(content, "\n$1")
	// 在中文题型标头前插入换行符，如 "二."、"三."等（允许可能存在空格）
	content = regexp.MustCompile(`(?m)\s*([一二三四五六七八九十]+[.。、])`).ReplaceAllString(content, "\n$1")
	// 在题型关键字前后插入换行符
	content = regexp.MustCompile(`(?i)(单选题|多选题|判断题|简答题)`).ReplaceAllString(content, "\n$1\n")
	// 对“参考答案：”后仅捕获答案部分（A-D），并在其后插入换行符
	content = regexp.MustCompile(`(?i)(参考答案|答案)[：:\s]*([A-D]+)\b`).ReplaceAllString(content, "$1：$2\n")
	// 合并多个换行符为一个
	content = regexp.MustCompile(`\n+`).ReplaceAllString(content, "\n")
	return strings.TrimSpace(content)
}

// parseQuestions 根据行数组解析题目
func (l *ParseWordQuestionsLogic) parseQuestions(lines []string) []*model.Question {
	var questions []*model.Question
	var currentType int64 = 0 // 1:单选,2:多选,3:判断,4:简答
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		// 判断题型标头
		if strings.Contains(line, "单选题") {
			currentType = 1
			i++
			continue
		} else if strings.Contains(line, "多选题") {
			currentType = 2
			i++
			continue
		} else if strings.Contains(line, "判断题") {
			currentType = 3
			i++
			continue
		} else if strings.Contains(line, "简答题") {
			currentType = 4
			i++
			continue
		}
		// 当题型确定且行以数字加点开始，则认为是题目起始
		questionStartPattern := regexp.MustCompile(`^\d+\.\s*(.*)`)
		if currentType > 0 && questionStartPattern.MatchString(line) {
			q, consumed := l.parseQuestion(lines[i:], currentType)
			if q != nil {
				questions = append(questions, q)
			}
			i += consumed
		} else {
			i++
		}
	}
	return questions
}

// parseQuestion 解析单道题目，返回题目对象和消耗的行数
// 若题目、选项、答案都在一行内，则进行特殊拆分
func (l *ParseWordQuestionsLogic) parseQuestion(lines []string, qType int64) (*model.Question, int) {
	if len(lines) == 0 {
		return nil, 0
	}
	// 尝试获取第一行的题目内容（例如 "1. xxx"）
	questionRegex := regexp.MustCompile(`^\d+\.\s*(.*)`)
	firstLine := strings.TrimSpace(lines[0])
	matches := questionRegex.FindStringSubmatch(firstLine)
	if len(matches) < 2 {
		return nil, 1
	}
	remaining := matches[1]
	consumed := 1

	// 如果第一行中包含选项标记"A."，则认为整道题都在一行
	if strings.Contains(remaining, "A.") {
		// 将题干、选项和答案从该行拆分出来
		posA := strings.Index(remaining, "A.")
		questionContent := strings.TrimSpace(remaining[:posA])
		remainder := remaining[posA:]
		// 查找答案标记"参考答案：" 或 "答案："
		ansMarker := ""
		posAns := -1
		if idx := strings.Index(remainder, "参考答案："); idx != -1 {
			posAns = idx
			ansMarker = "参考答案："
		} else if idx := strings.Index(remainder, "答案："); idx != -1 {
			posAns = idx
			ansMarker = "答案："
		}
		var optionPart, answerPart string
		if posAns != -1 {
			optionPart = strings.TrimSpace(remainder[:posAns])
			answerPart = strings.TrimSpace(remainder[posAns+len(ansMarker):])
		} else {
			optionPart = strings.TrimSpace(remainder)
			answerPart = ""
		}
		// 解析选项：用正则匹配每个选项
		optionRegex := regexp.MustCompile(`([A-D])\.\s*([^A-D]+)`)
		optionMatches := optionRegex.FindAllStringSubmatch(optionPart, -1)
		var opts []Option
		for _, m := range optionMatches {
			if len(m) >= 3 {
				opts = append(opts, Option{
					Label: strings.TrimSpace(m[1]),
					Text:  strings.TrimSpace(m[2]),
				})
			}
		}
		// 构造题目对象
		question := &model.Question{
			Type:    qType,
			Content: questionContent,
		}
		if len(opts) > 0 {
			optBytes, _ := json.Marshal(opts)
			question.Options = string(optBytes)
		} else {
			question.Options = "[]"
		}
		// 根据题型处理答案
		switch qType {
		case 1:
			if answerPart != "" {
				question.Answer = strconv.Quote(answerPart) // 直接序列化字符串
				question.Type = 1
			}
		case 2:
			if answerPart != "" {
				// 拆分答案字符串为字符数组（例如"B"或"ABCD"）
				ansArr := strings.Split(answerPart, "")
				ansBytes, _ := json.Marshal(ansArr)
				question.Answer = string(ansBytes)
				question.Type = 2
			}
		case 3:
			boolAns := strings.HasPrefix(answerPart, "A") || strings.Contains(strings.ToLower(answerPart), "正")
			ansBytes, _ := json.Marshal(boolAns)
			question.Answer = string(ansBytes)
			question.Type = 3
		case 4:
			question.Answer = strconv.Quote(answerPart)
			question.Type = 4
		}
		return question, consumed
	}

	// 否则，按逐行解析处理
	var contentBuilder strings.Builder
	var opts []Option
	var answer string
	// 先将第一行内容作为题干
	contentBuilder.WriteString(remaining)

	// 定义正则
	optionRegex := regexp.MustCompile(`^([A-D])\.\s*(.+)$`)
	answerRegex := regexp.MustCompile(`(?i)^(参考答案|答案)[：:\s]*(.+)$`)
	markerRegex := regexp.MustCompile(`[A-D]\.`)

	for i := consumed; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			consumed++
			continue
		}
		// 若遇到下一题起始则结束
		if questionRegex.MatchString(line) {
			break
		}
		// 检查是否含有选项标记但整行不符合选项格式，则拆分处理
		if (qType == 1 || qType == 2) && !optionRegex.MatchString(line) && markerRegex.MatchString(line) {
			loc := markerRegex.FindStringIndex(line)
			if loc != nil && loc[0] > 0 {
				prefix := strings.TrimSpace(line[:loc[0]])
				suffix := strings.TrimSpace(line[loc[0]:])
				if prefix != "" {
					contentBuilder.WriteString(" " + prefix)
				}
				line = suffix
			}
		}
		// 尝试匹配选项
		if (qType == 1 || qType == 2) && optionRegex.MatchString(line) {
			m := optionRegex.FindStringSubmatch(line)
			if len(m) == 3 {
				opts = append(opts, Option{
					Label: m[1],
					Text:  m[2],
				})
			}
			consumed++
			continue
		}
		// 尝试匹配答案行
		if answerRegex.MatchString(line) {
			m := answerRegex.FindStringSubmatch(line)
			if len(m) >= 3 {
				answer = strings.TrimSpace(m[2])
			}
			consumed++
			break
		}
		// 否则，作为题干续行
		contentBuilder.WriteString(" " + line)
		consumed++
	}
	question := &model.Question{
		Type:    qType,
		Content: strings.TrimSpace(contentBuilder.String()),
	}
	if len(opts) > 0 {
		optBytes, _ := json.Marshal(opts)
		question.Options = string(optBytes)
	} else {
		question.Options = "[]"
	}
	switch qType {
	case 1:
		if answer != "" {

			ansBytes, _ := json.Marshal(answer) // 直接序列化字符串，而不是数组
			question.Answer = string(ansBytes)
			question.Type = 1

		}
	case 2:
		if answer != "" {
			// 多选题：将答案字符串拆分为字符数组（例如 "ABCD" -> ["A","B","C","D"]）
			ansArr := strings.Split(answer, "")
			ansBytes, _ := json.Marshal(ansArr)
			question.Answer = string(ansBytes)
			question.Type = 2
		}
	case 3:
		boolAns := strings.HasPrefix(answer, "A") || strings.Contains(strings.ToLower(answer), "正")
		ansBytes, _ := json.Marshal(boolAns)
		question.Answer = string(ansBytes)
		question.Type = 3
	case 4:
		question.Answer = strconv.Quote(answer)
		question.Type = 4
	}
	return question, consumed
}
