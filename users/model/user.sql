CREATE TABLE users (
                       id INT NOT NULL AUTO_INCREMENT COMMENT '用户id（主键）',
                       type INT NOT NULL DEFAULT 2 COMMENT '用户类型（0-管理员、1-教师、2-学生）',
                       name VARCHAR(50) NOT NULL DEFAULT '' COMMENT '用户名',
                       username VARCHAR(50) NOT NULL DEFAULT '' COMMENT '账号（工号/学号）',
                       passwd VARCHAR(100) NOT NULL DEFAULT '' COMMENT '用户密码哈希值',
                       salt INT NOT NULL DEFAULT 0 COMMENT '盐',
                       create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       delete_time DATETIME DEFAULT NULL COMMENT '删除时间',
                       is_delete INT NOT NULL DEFAULT 0 COMMENT '是否删除（0-未删除，1-已删除）',
                       PRIMARY KEY (id),
                       CONSTRAINT unique_username UNIQUE (username) -- 唯一约束，使用CONSTRAINT关键字
) COMMENT='用户表';
