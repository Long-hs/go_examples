-- 创建数据库
CREATE DATABASE IF NOT EXISTS double_token_example
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_0900_ai_ci;

USE double_token_example;

-- 用户表
CREATE TABLE user (
    id BIGINT AUTO_INCREMENT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password VARCHAR(100) NOT NULL COMMENT '密码',
    salt VARCHAR(32) NOT NULL COMMENT '密码盐值',
    phone VARCHAR(20) NOT NULL COMMENT '手机号',
    email VARCHAR(100) NULL COMMENT '邮箱',
    status TINYINT DEFAULT 1 NULL COMMENT '状态：1-正常，0-禁用',
    last_login_time TIMESTAMP NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(50) NULL COMMENT '最后登录IP',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY idx_email (email),
    UNIQUE KEY idx_phone (phone),
    UNIQUE KEY idx_username (username)
) COMMENT '用户表';

-- 订单表
CREATE TABLE `order` (
    id BIGINT AUTO_INCREMENT COMMENT '支付ID',
    creator_id BIGINT NOT NULL COMMENT '创建人ID',
    goods_id VARCHAR(48) NOT NULL COMMENT '商品ID',
    payment_no VARCHAR(36) NOT NULL COMMENT '支付流水号',
    quantity TINYINT NOT NULL COMMENT '购买数量',
    amount DECIMAL(10, 2) NOT NULL COMMENT '支付金额',
    payment_method TINYINT NOT NULL COMMENT '支付方式：1-支付宝，2-微信',
    status TINYINT NOT NULL COMMENT '支付状态：1-待支付，2-支付成功，3-支付失败',
    pay_time TIMESTAMP NULL COMMENT '支付时间',
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY idx_payment_no (payment_no),
    FOREIGN KEY (creator_id) REFERENCES user (id)
) COMMENT '订单表';

CREATE TABLE `refresh_tokens` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `jti` VARCHAR(36) NOT NULL COMMENT 'Token 唯一标识',
  `expires_at` DATETIME NOT NULL COMMENT '过期时间',
  `created_at` DATETIME NOT NULL COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_jti` (`jti`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='RefreshToken表';