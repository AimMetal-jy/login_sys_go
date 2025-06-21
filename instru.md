## 数据库信息
-- 使用已创建的数据库
USE login_sys;

-- 创建用户表
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID，主键',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名，唯一',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_active BOOLEAN DEFAULT TRUE COMMENT '账户是否激活'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户登录表';

-- 为用户名字段创建索引以提高查询性能
CREATE INDEX idx_username ON users(username);

-- 为激活状态创建索引
CREATE INDEX idx_is_active ON users(is_active);

## 后端要求
现在来设计后端项目，用go语言实现。实现登录和注册功能，以及密码哈希加密。

同时创建了两个Python文件，文件名分别为test_login.py和test_register.py的登录和注册测试用例，

测试用户名：test、测试密码：123456。

## 注意事项
- 一个项目中go语言只能有一个main程序入口。
- 设计项目时仔细检查代码，避免出现错误。