CREATE TABLE generated_paper (
                                 id INT AUTO_INCREMENT COMMENT '记录ID',
                                 exam_id INT NOT NULL COMMENT '关联exam表',
                                 student_id INT NOT NULL COMMENT '关联学生（user表）',
                                 questions JSON NOT NULL COMMENT '学生随机生成的试卷题目列表（题目ID和分值的JSON）',
                                 total_score INT NOT NULL COMMENT '试卷总分',
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
                                 PRIMARY KEY (id)
) COMMENT='学生随机试卷记录表';

CREATE TABLE exam (
                      id INT AUTO_INCREMENT COMMENT '考试ID',
                      name VARCHAR(255) NOT NULL COMMENT '考试名称',
                      teacher_id INT NOT NULL COMMENT '发布考试的教师ID',
                      class_id INT NOT NULL COMMENT '关联班级ID',
                      exam_type ENUM('fixed', 'random') NOT NULL COMMENT '试题生成方式（固定or随机）',
                      total_score INT NOT NULL COMMENT '总分',
                      requires_manual_grading BOOLEAN DEFAULT FALSE COMMENT '是否包含需要人工阅卷的题目',
                      start_time DATETIME NOT NULL COMMENT '考试开始时间',
                      end_time DATETIME NOT NULL COMMENT '考试结束时间',
                      paper_id INT DEFAULT NULL COMMENT '固定试卷ID，exam_type=fixed 时使用',
                      paper_rule_id INT DEFAULT NULL COMMENT '随机试卷生成规则ID，exam_type=random时使用',
                      can_view_results BOOLEAN DEFAULT FALSE COMMENT '考试结束后是否公开试题',
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                      PRIMARY KEY (id)
) COMMENT='考试表';

CREATE TABLE exam_attempt (
                              id INT AUTO_INCREMENT COMMENT '记录ID',
                              exam_id INT NOT NULL COMMENT '关联exam表',
                              student_id INT NOT NULL COMMENT '关联学生（user表）',
                              paper_id INT DEFAULT NULL COMMENT '固定试卷ID（exam_type=fixed时使用）',
                              generated_paper_id INT DEFAULT NULL COMMENT '随机试卷ID（exam_type=random时使用）',
                              score INT DEFAULT NULL COMMENT '学生成绩（未评定前为 NULL）',
                              status  ENUM('not_started', 'ongoing', 'submitted', 'graded') DEFAULT 'not_started' COMMENT '考试状态',
                              start_time DATETIME DEFAULT NULL COMMENT '学生开始考试时间',
                              submit_time DATETIME DEFAULT NULL COMMENT '提交时间',
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
                              PRIMARY KEY (id)
) COMMENT='学生考试记录表';



CREATE TABLE exam_answer (
                             id INT AUTO_INCREMENT COMMENT '记录ID',
                             attempt_id INT NOT NULL COMMENT '关联exam_attempt表',
                             answer JSON NOT NULL COMMENT '学生答案（JSON格式）',
                             score_details JSON DEFAULT NULL COMMENT '每题得分详情',
                             grading_status ENUM('pending', 'auto_scored', 'manual_scored') DEFAULT 'pending' COMMENT '评分状态',
                             submit_time DATETIME DEFAULT NULL COMMENT '提交时间',
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
                             PRIMARY KEY (id)
) COMMENT='学生答案表';