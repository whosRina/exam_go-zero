-- 试卷表
CREATE TABLE paper (
                       id INT AUTO_INCREMENT COMMENT '试卷ID',
                       name VARCHAR(255) NOT NULL COMMENT '试卷名称',
                       total_score INT NOT NULL COMMENT '试卷总分',
                       created_by INT NOT NULL COMMENT '创建者ID（关联user表）',
                       questions JSON NOT NULL COMMENT '试卷题目列表（存储题目ID和分值的JSON）',
                       created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                       PRIMARY KEY (id)
) COMMENT='试卷表';

-- 随机组卷规则表
CREATE TABLE paper_rule (
                      id INT AUTO_INCREMENT COMMENT '规则ID',
                      name VARCHAR(255) NOT NULL COMMENT '规则名称',
                      total_score INT NOT NULL COMMENT '随机生成的试卷总分',
                      created_by INT NOT NULL COMMENT '创建者ID（关联user表）',
                      bank_id INT NOT NULL COMMENT '关联题库ID',
                      num_questions JSON NOT NULL COMMENT '各题型随机抽取数量（{"1": 5, "2": 3}）',
                      score_config JSON NOT NULL COMMENT '各题型分值（{"1": 2, "2": 3}）',
                      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                      updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                      PRIMARY KEY (id)

) COMMENT='随机组卷规则表';