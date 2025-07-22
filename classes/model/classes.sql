-- 班级表
CREATE TABLE class (
                       id INT NOT NULL AUTO_INCREMENT COMMENT '班级ID（主键）',
                       name VARCHAR(255) NOT NULL COMMENT '班级名称（即课程名称）',
                       teacher_id INT NOT NULL COMMENT '教师ID',
                       class_code VARCHAR(255) NOT NULL UNIQUE COMMENT '班级唯一邀请码',
                       is_joinable TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否允许加入（0-禁止，1-允许）',
                       PRIMARY KEY (id),
                       INDEX idx_teacher_id (teacher_id) -- 代替外键
) COMMENT='班级表';

-- 班级成员表
CREATE TABLE class_member (
                              id INT NOT NULL AUTO_INCREMENT COMMENT '班级成员ID（主键）',
                              class_id INT NOT NULL COMMENT '班级ID',
                              user_id INT NOT NULL COMMENT '用户ID',
                              PRIMARY KEY (id),
                              INDEX idx_class_id (class_id), -- 代替外键
                              INDEX idx_user_id (user_id)    -- 代替外键
) COMMENT='班级成员表';
