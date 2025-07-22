-- 题库表
CREATE TABLE question_bank (
                               id INT AUTO_INCREMENT  COMMENT '题库ID',
                               name VARCHAR(255) NOT NULL COMMENT '题库名称',
                               created_by INT NOT NULL COMMENT '创建者ID（关联user表）',
                               created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               PRIMARY KEY (id)  -- 这里显式指定主键
) COMMENT='题库表';

-- 题目表
CREATE TABLE question (
                          id INT AUTO_INCREMENT COMMENT '题目ID',
                          bank_id INT NOT NULL COMMENT '题库ID（关联question_bank表）',
                          content TEXT NOT NULL COMMENT '题目内容',
                          type TINYINT(1) NOT NULL COMMENT '题型 (1=单选，2=多选，3=判断，4=简答)',
                          options JSON NOT NULL COMMENT '可选项（存储JSON格式）',
                          answer JSON NOT NULL COMMENT '正确答案（存储JSON格式）',
                          created_by INT NOT NULL COMMENT '创建者ID（关联user表）',
                          created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          PRIMARY KEY (id)  -- 这里显式指定主键
) COMMENT='题目表';

-- 添加索引，提高查询性能
CREATE INDEX idx_question_bank_id ON question(bank_id);
CREATE INDEX idx_question_created_by ON question(created_by);
