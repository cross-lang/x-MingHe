-- MySQL dump 10.13  Distrib 8.0.33, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: minghe
-- ------------------------------------------------------
-- Server version	8.0.45

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `x_enterprise`
--

DROP TABLE IF EXISTS `x_enterprise`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x_enterprise` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `p_id` int unsigned NOT NULL COMMENT '企业所属园区ID',
  `type` tinyint NOT NULL COMMENT '企业类型（1：企业；2：学院）',
  `full_name` varchar(50) NOT NULL COMMENT '企业全称',
  `short_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '企业简称',
  `description` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '企业简介',
  `icon` varchar(255) NOT NULL COMMENT '企业图标',
  `image` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '企业图片',
  `contact` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '联系方式',
  `website` longtext NOT NULL COMMENT '企业官网',
  `label` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '标签',
  `lease_time` date NOT NULL COMMENT '入驻时间',
  `top` tinyint NOT NULL COMMENT '置顶状态（1：置顶；0：不置顶）',
  `weight` int NOT NULL COMMENT '权重',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '账号状态（1：启用；-1：禁用）',
  `created_by` int unsigned NOT NULL COMMENT '创建人ID',
  `updated_by` int unsigned NOT NULL COMMENT '修改人ID',
  `created_at` timestamp NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NOT NULL COMMENT '修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `p_id` (`p_id`,`type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='企业表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `x_user`
--

DROP TABLE IF EXISTS `x_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(11) NOT NULL COMMENT '姓名',
  `avatar` varchar(255) NOT NULL COMMENT '头像',
  `block_status` tinyint NOT NULL COMMENT '拉黑状态（1：未拉黑；2：已拉黑）',
  `account_status` tinyint NOT NULL COMMENT '账号状态（1：正常；2：禁用）',
  `account` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '账号',
  `phone_number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '手机号',
  `gender` tinyint NOT NULL COMMENT '性别（0：保密；1：男；2：女）',
  `verification_status` tinyint NOT NULL COMMENT '实名认证状态（1：已认证；0：未认证）',
  `graduated_school` varchar(50) NOT NULL COMMENT '毕业院校',
  `profession` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '职业',
  `residence` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '居住地',
  `contact` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '联系方式',
  `work_years` double(10,2) NOT NULL COMMENT '工作年限',
  `introduction` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '个人简介',
  `skills` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '核心技能，多个使用分号分割',
  `deactivate_status` tinyint NOT NULL COMMENT '注销状态（0：未注销；1：注销中；2：已注销）',
  `deactivate_at` timestamp NULL DEFAULT NULL COMMENT '注销时间',
  `created_at` timestamp NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NOT NULL COMMENT '修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `account` (`account`) USING BTREE,
  UNIQUE KEY `phone_number` (`phone_number`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `x_user_enterprise`
--

DROP TABLE IF EXISTS `x_user_enterprise`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x_user_enterprise` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `u_id` int unsigned NOT NULL COMMENT '用户ID',
  `e_id` int unsigned NOT NULL COMMENT '企业ID',
  `p_id` int unsigned NOT NULL COMMENT '园区ID',
  `type` tinyint NOT NULL COMMENT '入职企业类型（1：企业；2：学院）',
  `created_at` timestamp NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NOT NULL COMMENT '修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `u_id` (`u_id`) USING BTREE,
  KEY `e_id` (`e_id`) USING BTREE,
  KEY `p_id` (`p_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户企业关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `x_user_identity_verification`
--

DROP TABLE IF EXISTS `x_user_identity_verification`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `x_user_identity_verification` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `u_id` int unsigned NOT NULL COMMENT '用户ID',
  `channel` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '认证渠道',
  `name` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '姓名',
  `id_card` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '身份证号码',
  `hash` char(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '身份证号码Hash',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_u_id` (`u_id`),
  KEY `idx_hash` (`hash`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户实名认证表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'minghe'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-04-20 15:18:19
