/*
 Navicat Premium Data Transfer

 Source Server         : my-server
 Source Server Type    : MySQL
 Source Server Version : 100703 (10.7.3-MariaDB-log)
 Source Host           : 127.0.0.1:3306
 Source Schema         : app

 Target Server Type    : MySQL
 Target Server Version : 100703 (10.7.3-MariaDB-log)
 File Encoding         : 65001

 Date: 01/02/2024 23:34:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_transaction
-- ----------------------------
DROP TABLE IF EXISTS `app_transaction`;
CREATE TABLE `app_transaction` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `block_number` int(11) DEFAULT NULL COMMENT '块高度',
  `hash` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'hash',
  `from` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'from',
  `to` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'to',
  `protected` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '2' COMMENT '是否保护(1-是 2-否)',
  `value` decimal(30,0) DEFAULT 0 COMMENT '币数',
  `type` char(1) DEFAULT '0' COMMENT '交易类型',
  `gas_price` decimal(30,0) DEFAULT 0 COMMENT 'gas price',
  `cost` decimal(30,0) DEFAULT 0 COMMENT '手续费',
  `created_at` datetime DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_hash` (`hash`) COMMENT '交易哈希'
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='交易记录';

SET FOREIGN_KEY_CHECKS = 1;
