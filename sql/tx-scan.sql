/*
 Navicat Premium Data Transfer

 Source Server         : my-server
 Source Server Type    : MySQL
 Source Server Version : 100703
 Source Host           : 127.0.0.1:3306
 Source Schema         : evm-scan

 Target Server Type    : MySQL
 Target Server Version : 100703
 File Encoding         : 65001

 Date: 28/02/2024 09:33:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_transaction
-- ----------------------------
DROP TABLE IF EXISTS `app_transaction`;
CREATE TABLE `app_transaction` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `block_number` bigint(20) DEFAULT NULL COMMENT '块高度',
  `hash` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'hash',
  `from` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'from',
  `to` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'to',
  `protected` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '2' COMMENT '是否保护(1-是 2-否)',
  `value` decimal(30,0) DEFAULT 0 COMMENT '币数',
  `type` char(1) DEFAULT '0' COMMENT '交易类型',
  `effective_gas_price` decimal(30,0) DEFAULT 0 COMMENT 'gas price',
  `gas_used` bigint(20) DEFAULT 0 COMMENT 'gas使用',
  `nonce` bigint(20) DEFAULT 0 COMMENT 'nonce',
  `transaction_index` int(11) DEFAULT 0 COMMENT '交易索引',
  `receipt_status` int(11) DEFAULT 0 COMMENT '回执状态',
  `created_at` datetime DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_hash` (`hash`) COMMENT '交易哈希'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易记录';

SET FOREIGN_KEY_CHECKS = 1;
