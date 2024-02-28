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

 Date: 28/02/2024 09:30:30
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_zkf_stat_daily_gas
-- ----------------------------
DROP TABLE IF EXISTS `app_zkf_stat_daily_gas`;
CREATE TABLE `app_zkf_stat_daily_gas` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `block_start` bigint(20) DEFAULT NULL COMMENT '启始高度',
  `block_end` bigint(20) DEFAULT NULL COMMENT '截止高度',
  `date_start` datetime DEFAULT NULL COMMENT '开始时间',
  `date_end` datetime DEFAULT NULL COMMENT '结束时间',
  `total_tx_count` int(11) DEFAULT NULL COMMENT '交易总数',
  `total_gas_fee` decimal(30,0) DEFAULT NULL COMMENT '总gas fee',
  `total_gas_price` decimal(30,0) DEFAULT NULL COMMENT '平均每笔交易gas price',
  `calc_status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '状态(1-进行中 2-结束)',
  `updated_at` datetime DEFAULT current_timestamp() COMMENT '更新时间',
  `created_at` datetime DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='zkf统计每日Gas';

-- ----------------------------
-- Table structure for app_zkf_stat_hours_gas
-- ----------------------------
DROP TABLE IF EXISTS `app_zkf_stat_hours_gas`;
CREATE TABLE `app_zkf_stat_hours_gas` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `block_start` bigint(20) DEFAULT NULL COMMENT '启始高度',
  `block_end` bigint(20) DEFAULT NULL COMMENT '截止高度',
  `date_start` datetime DEFAULT NULL COMMENT '开始时间',
  `date_end` datetime DEFAULT NULL COMMENT '结束时间',
  `total_tx_count` int(11) DEFAULT NULL COMMENT '交易总数',
  `total_gas_fee` decimal(30,0) DEFAULT NULL COMMENT '总gas',
  `total_gas_price` decimal(30,0) DEFAULT NULL COMMENT '总的gas价格，用于方便计算平均值',
  `calc_status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '状态(1-进行中 2-结束)',
  `updated_at` datetime DEFAULT current_timestamp() COMMENT '更新时间',
  `created_at` datetime DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='zkf统计每小时Gas';

-- ----------------------------
-- Table structure for app_zkf_stat_weekly_gas
-- ----------------------------
DROP TABLE IF EXISTS `app_zkf_stat_weekly_gas`;
CREATE TABLE `app_zkf_stat_weekly_gas` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '编号',
  `block_start` bigint(20) DEFAULT NULL COMMENT '启始高度',
  `block_end` bigint(20) DEFAULT NULL COMMENT '截止高度',
  `date_start` datetime DEFAULT NULL COMMENT '开始时间',
  `date_end` datetime DEFAULT NULL COMMENT '结束时间',
  `total_tx_count` int(11) DEFAULT NULL COMMENT '交易总数',
  `total_gas_fee` decimal(30,0) DEFAULT NULL COMMENT '总gas',
  `total_gas_price` decimal(30,0) DEFAULT NULL COMMENT '总的gas价格，用于方便计算平均值',
  `calc_status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '状态(1-进行中 2-结束)',
  `updated_at` datetime DEFAULT current_timestamp() COMMENT '更新时间',
  `created_at` datetime DEFAULT current_timestamp() COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='zkf统计每周Gas';

SET FOREIGN_KEY_CHECKS = 1;
