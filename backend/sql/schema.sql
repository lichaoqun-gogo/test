-- =====================================================
-- 开源治理平台（OSM）数据库表结构 DDDL
-- 数据库：osm
-- 字符集：utf8mb4
-- 说明：包含全部 38 张业务表的 CREATE TABLE 语句
-- 生成方式：mysqldump --no-data 从运行中的 MySQL 导出
-- 最新导出时间：2026-10-30 16:32:45
-- =====================================================


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
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;
SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ 'f71daf3e-707b-11f1-92f3-be42118d902a:1-8822';
DROP TABLE IF EXISTS `approval_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint NOT NULL COMMENT '关联任务ID',
  `task_no` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '任务编号',
  `node_id` bigint NOT NULL COMMENT '审批节点ID',
  `node_key` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '节点标识',
  `node_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '节点名称',
  `approver` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '审批人',
  `approver_dept` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作人部门',
  `approver_id` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作人ID',
  `action` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '操作: approve-通过, reject-驳回',
  `opinion` varchar(1000) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '审批意见',
  `from_status` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作前状态',
  `to_status` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作后状态',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_task_no` (`task_no`)
) ENGINE=InnoDB AUTO_INCREMENT=104 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批历史表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `approval_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_node` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint NOT NULL COMMENT '关联任务ID',
  `node_key` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '节点标识（来自模板）',
  `node_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '节点名称(如使用部门意见)',
  `node_order` int NOT NULL COMMENT '节点顺序',
  `node_type` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '节点类型：fixed, conditional',
  `approver_id` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '审批人ID',
  `approver` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '审批人',
  `approver_dept` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '审批人部门',
  `status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待审批, approved-通过, rejected-驳回, not_required-不需要',
  `opinion` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '审批意见',
  `approved_at` datetime DEFAULT NULL COMMENT '审批时间',
  `is_skipped` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否被跳过',
  `skip_reason` varchar(200) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '跳过原因',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_approver` (`approver`),
  KEY `idx_task_order` (`task_id`,`node_order`)
) ENGINE=InnoDB AUTO_INCREMENT=147 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批节点记录表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `approval_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_template` (
  `id` bigint NOT NULL,
  `template_key` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '模板唯一标识：introduce, upgrade, withdraw, usage, risk_report',
  `template_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '模板名称，如：引入申请审批流',
  `description` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '模板描述',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用：1-启用，0-禁用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_template_key` (`template_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批流模板表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `approval_template_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_template_node` (
  `id` bigint NOT NULL,
  `template_id` bigint NOT NULL COMMENT '关联的模板ID',
  `node_key` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '节点唯一标识：node_0, node_1, node_2, node_3',
  `node_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '节点名称',
  `node_order` int NOT NULL COMMENT '节点顺序（从0开始）',
  `node_type` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'fixed' COMMENT '节点类型：fixed-固定节点，conditional-条件节点',
  `condition_type` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '条件类型：dept_specialist-本部门专员，usage_dept-使用部门，mgmt_dept-管理部门',
  `is_skipable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否可跳过：1-可跳过，0-不可跳过',
  `skip_rule` varchar(200) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '跳过规则JSON，如：{"type":"same_person","field":"approver_id"}',
  `approver_role` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '审批人角色：initiator-发起人，dept_specialist-部门专员，usage_dept_specialist-使用部门专员，mgmt_dept_specialist-管理部门专员',
  `description` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '节点描述',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批流模板节点定义表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `attachment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `attachment` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `source` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '附件来源（AttachmentSource 枚举：BOM_UPLOAD/VULDB_PLAN/VULDB_RESULT/FIRST_VERIFY/FINAL_VERIFY）',
  `filename` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '文件名',
  `md5` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '文件md5值',
  `description` mediumtext COLLATE utf8mb4_bin COMMENT '描述和备注',
  `store_url` varchar(150) COLLATE utf8mb4_bin NOT NULL COMMENT '对象存储中的相对地址',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_attachment_idx` (`source`,`filename`,`store_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='附件管理';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `bom`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bom` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `bom_format` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'sbom工具类型',
  `bom_version` int DEFAULT NULL COMMENT 'bom生成版本',
  `spec_version` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'sbom工具版本',
  `bom_create` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'bom创建时间',
  `serial_number` varchar(36) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'sbom文件uuid',
  `store_url` bigint DEFAULT NULL COMMENT 'sbom文件|附件管理attachmentid',
  `project_id` bigint NOT NULL COMMENT '关联项目',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  KEY `idx_bom_project` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='BOM清单';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `commit_id_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `commit_id_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `commit_id` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '提交commitId',
  `create_time_history` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '历史创建时间',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '后台归档时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_commit_id` (`commit_id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='历史commit记录';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `component`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `component` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `current_id` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '关联source_check表',
  `PARENT_COMPONENT_ID` bigint DEFAULT NULL COMMENT '上层依赖',
  `AUTHOR` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件组织',
  `_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '组件名称',
  `_version` varchar(50) COLLATE utf8mb4_bin DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件类型',
  `FILENAME` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名',
  `EXTENSION` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '扩展信息',
  `INTERNAL` int DEFAULT NULL COMMENT '是否内部 0-否 1-是',
  `MD5` varchar(32) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `SHA1` varchar(40) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `SHA_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `SHA_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `SHA_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `SHA3_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `SHA3_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `SHA3_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `BLAKE2B_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `BLAKE2B_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `BLAKE2B_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `BLAKE3` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '校验信息',
  `NOTES` mediumtext COLLATE utf8mb4_bin COMMENT '详细参考',
  `PURL` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Package-URL',
  `PURLCOORDINATES` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Package-URL Coordinates',
  `CPE` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'CPE',
  `SWIDTAGID` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SWID',
  `DIRECT_DEPENDENCIES` mediumtext COLLATE utf8mb4_bin COMMENT 'SBOM依赖列表',
  `EXTERNAL_REFERENCES` mediumblob COMMENT '扩展说明',
  `COPYRIGHT` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证说明',
  `LICENSE` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证',
  `LICENSE_ID` bigint DEFAULT NULL COMMENT '关联license表',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0.000' COMMENT '最后风险值',
  `PROJECT_ID` bigint DEFAULT NULL COMMENT '关联project_inner表',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  KEY `idx_project_id` (`PROJECT_ID`),
  KEY `idx_current_id` (`current_id`)
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='自研项目依赖包表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `component_black`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `component_black` (
  `ID` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `PROJECT_ID` bigint NOT NULL COMMENT '关联的项目',
  `COMPONENT_ID` bigint DEFAULT NULL COMMENT '关联的组件',
  `BLACK_ID` bigint NOT NULL COMMENT '关联的黑名单',
  `PROCESS_STATUS` int NOT NULL COMMENT '处理状态 0-未处理|1-处理中|2-处理完成|3-误报',
  `reason` varchar(600) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '阻断原因',
  `jira_no` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'jira工单号',
  `pro_time` datetime DEFAULT NULL COMMENT '上线修改时间',
  `test_time` datetime DEFAULT NULL COMMENT '测试修改时间',
  `DESCRIPTION` mediumtext COLLATE utf8mb4_bin COMMENT '处理描述和备注',
  `PROCESS_PLAN_URL` bigint DEFAULT NULL COMMENT '处理计划附件attachmentId',
  `PROCESS_RESULT_URL` bigint DEFAULT NULL COMMENT '处理结果附件attachmentId',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `CURRENT_ID` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '本次数据标识ID',
  `internal` int NOT NULL DEFAULT '1' COMMENT '是否内部项目 0-否|1-是',
  `alert_type` tinyint DEFAULT '1' COMMENT '告警类型 1阻断 2警告',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uniq_project_balck_commit_idx` (`PROJECT_ID`,`COMPONENT_ID`,`BLACK_ID`,`CURRENT_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=210027 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='组件黑名单处理列表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `component_core`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `component_core` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `_format` tinyint NOT NULL DEFAULT '1' COMMENT '类型: maven:1、pypi:2、npm:3、yum:4等',
  `type` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '类型',
  `_name` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '包名',
  `_group` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '包分组',
  `version_exp` varchar(200) COLLATE utf8mb4_bin NOT NULL COMMENT '匹配表达式',
  `applicable_level` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '适用等级',
  `support_team` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '技术支持团队',
  `reason` varchar(400) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '添加理由',
  `advise` varchar(400) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '建议',
  `_enable` tinyint(1) NOT NULL COMMENT '是否启用, 1: 启用, 0: 禁用',
  `alert_type` tinyint DEFAULT '1' COMMENT '告警类型 1 阻断 2警告',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  `level_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '风险等级生效时间',
  `is_boundary` tinyint(1) DEFAULT '2' COMMENT '2:全部 1:边界系统 0非边界系统',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_black_list_idx` (`_format`,`_name`,`_group`,`applicable_level`)
) ENGINE=InnoDB AUTO_INCREMENT=240006 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='黑名单配置表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `component_snapshot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `component_snapshot` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `current_id` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '关联source_check表标识',
  `PARENT_COMPONENT_ID` bigint DEFAULT NULL COMMENT '上层依赖组件ID',
  `AUTHOR` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件组织',
  `_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '组件名称',
  `_version` varchar(50) COLLATE utf8mb4_bin DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件类型',
  `FILENAME` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名',
  `EXTENSION` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '扩展信息',
  `INTERNAL` int DEFAULT NULL COMMENT '是否内部: 0-否, 1-是',
  `MD5` varchar(32) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'MD5校验值',
  `SHA1` varchar(40) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA1校验值',
  `SHA_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA256校验值',
  `SHA_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA384校验值',
  `SHA_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA512校验值',
  `SHA3_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA3-256校验值',
  `SHA3_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA3-384校验值',
  `SHA3_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA3-512校验值',
  `BLAKE2B_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE2B-256校验值',
  `BLAKE2B_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE2B-384校验值',
  `BLAKE2B_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE2B-512校验值',
  `BLAKE3` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE3校验值',
  `NOTES` mediumtext COLLATE utf8mb4_bin COMMENT '详细参考',
  `PURL` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Package-URL',
  `PURLCOORDINATES` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Package-URL Coordinates',
  `CPE` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'CPE标识',
  `SWIDTAGID` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SWID标识',
  `DIRECT_DEPENDENCIES` mediumtext COLLATE utf8mb4_bin COMMENT 'SBOM依赖列表',
  `EXTERNAL_REFERENCES` mediumblob COMMENT '扩展说明',
  `COPYRIGHT` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证说明',
  `LICENSE` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证',
  `LICENSE_ID` bigint DEFAULT NULL COMMENT '关联license表ID',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0.000' COMMENT '最后风险评分',
  `PROJECT_ID` bigint DEFAULT NULL COMMENT '关联project_inner表ID',
  `is_delete` tinyint NOT NULL DEFAULT '1' COMMENT '逻辑删除标识: 0-已删除, 1-未删除',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_current_id` (`current_id`),
  KEY `idx_project_id` (`PROJECT_ID`),
  KEY `idx_name_version` (`_name`,`_version`),
  KEY `idx_license_id` (`LICENSE_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=230 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='开源依赖包表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `component_snapshot_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `component_snapshot_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `current_id` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '关联source_check表标识',
  `PARENT_COMPONENT_ID` bigint DEFAULT NULL COMMENT '上层依赖组件ID',
  `AUTHOR` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件组织',
  `_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '组件名称',
  `_version` varchar(50) COLLATE utf8mb4_bin DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件类型',
  `FILENAME` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '文件名',
  `EXTENSION` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '扩展信息',
  `INTERNAL` int DEFAULT NULL COMMENT '是否内部: 0-否, 1-是',
  `MD5` varchar(32) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'MD5校验值',
  `SHA1` varchar(40) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA1校验值',
  `SHA_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA256校验值',
  `SHA_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA384校验值',
  `SHA_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA512校验值',
  `SHA3_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA3-256校验值',
  `SHA3_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA3-384校验值',
  `SHA3_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SHA3-512校验值',
  `BLAKE2B_256` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE2B-256校验值',
  `BLAKE2B_384` varchar(96) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE2B-384校验值',
  `BLAKE2B_512` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE2B-512校验值',
  `BLAKE3` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'BLAKE3校验值',
  `NOTES` mediumtext COLLATE utf8mb4_bin COMMENT '详细参考',
  `PURL` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Package-URL',
  `PURLCOORDINATES` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Package-URL Coordinates',
  `CPE` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'CPE标识',
  `SWIDTAGID` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SWID标识',
  `DIRECT_DEPENDENCIES` mediumtext COLLATE utf8mb4_bin COMMENT 'SBOM依赖列表',
  `EXTERNAL_REFERENCES` mediumblob COMMENT '扩展说明',
  `COPYRIGHT` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证说明',
  `LICENSE` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证',
  `LICENSE_ID` bigint DEFAULT NULL COMMENT '关联license表ID',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0.000' COMMENT '最后风险评分',
  `PROJECT_ID` bigint DEFAULT NULL COMMENT '关联project_inner表ID',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_current_id` (`current_id`),
  KEY `idx_project_id` (`PROJECT_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=183 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='开源依赖包历史归档表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `defect`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `defect` (
  `defect_id` varchar(100) COLLATE utf8mb4_general_ci NOT NULL COMMENT '唯一缺陷ID：软件小写#编号（redis#13401 / kafka#KAFKA-20449），upsert 业务键',
  `software` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
  `issue_url` varchar(500) COLLATE utf8mb4_general_ci NOT NULL,
  `title` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL,
  `body` mediumtext COLLATE utf8mb4_general_ci,
  `labels` json DEFAULT NULL,
  `reporter` varchar(200) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `comments_count` bigint NOT NULL DEFAULT '0',
  `source` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'github-issue',
  `severity` varchar(10) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '未知',
  `issue_state` varchar(10) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'open',
  `internal_status` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'pending_analysis',
  `state_reason` varchar(20) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `reactions` json DEFAULT NULL,
  `milestone` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `watch_matches` json DEFAULT NULL,
  `summary_zh` text COLLATE utf8mb4_general_ci,
  `repair_advice` text COLLATE utf8mb4_general_ci,
  `org` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `maint_dept` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `imported_at` datetime(3) NOT NULL,
  `title_zh` varchar(1000) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`defect_id`),
  KEY `idx_software_state` (`software`,`issue_state`),
  KEY `idx_severity` (`severity`),
  KEY `idx_updated_at` (`updated_at`),
  KEY `idx_internal_status` (`internal_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='缺陷采集台账（统一模型，GitHub/Jira 多通道，upsert 导入）';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `defect_disposal`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `defect_disposal` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `disposal_no` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `dept_name` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `disposal_type` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `assignee_username` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `assignee_name` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `ticket_no` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `disposal_note` text COLLATE utf8mb4_bin,
  `deadline` datetime(3) DEFAULT NULL,
  `status` varchar(20) COLLATE utf8mb4_bin DEFAULT 'pending_assign',
  `assigned_by` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `assigned_at` datetime(3) DEFAULT NULL,
  `completed_at` datetime(3) DEFAULT NULL,
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  `project_id` bigint DEFAULT '0',
  `project_name` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `project_type` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `defect_no` varchar(100) COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_assignee` (`assignee_username`),
  KEY `idx_defect_no` (`defect_no`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `defect_flow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `defect_flow` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `defect_no` varchar(100) COLLATE utf8mb4_bin NOT NULL,
  `flow_type` varchar(20) COLLATE utf8mb4_bin DEFAULT 'defect',
  `maint_dept` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `deadline` datetime(3) DEFAULT NULL,
  `analyze_result` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL,
  `review_meeting_no` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `review_conclusion` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `dispatch_type` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `dispatch_target` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `dispatch_dept` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL,
  `dispatch_time` datetime(3) DEFAULT NULL,
  `dispatcher` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `analyze_by` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `analyze_time` datetime(3) DEFAULT NULL,
  `review_by` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
  `review_time` datetime(3) DEFAULT NULL,
  `close_time` datetime(3) DEFAULT NULL,
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_defect_no` (`defect_no`),
  UNIQUE KEY `defect_no` (`defect_no`)
) ENGINE=InnoDB AUTO_INCREMENT=550 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='缺陷处置流程表（缺陷级状态机，按缺陷编号关联实体）';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `defect_ledger`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `defect_ledger` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `defect_no` varchar(100) COLLATE utf8mb4_bin NOT NULL,
  `source` varchar(30) COLLATE utf8mb4_bin NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `risk_level` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '未知',
  `description` mediumtext COLLATE utf8mb4_bin,
  `affected_component` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `discovered_date` datetime(3) NOT NULL,
  `internal_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'pending_analysis',
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_defect_no` (`defect_no`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_source` (`source`),
  KEY `idx_internal_status` (`internal_status`)
) ENGINE=InnoDB AUTO_INCREMENT=279 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='缺陷台账';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `defect_status_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `defect_status_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `disposal_id` bigint DEFAULT '0',
  `from_status` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `to_status` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `action` varchar(30) COLLATE utf8mb4_bin DEFAULT NULL,
  `operator` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `operator_name` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `remark` varchar(512) COLLATE utf8mb4_bin DEFAULT NULL,
  `create_time` datetime(3) DEFAULT NULL,
  `defect_no` varchar(100) COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_defect_no` (`defect_no`)
) ENGINE=InnoDB AUTO_INCREMENT=80 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `dept_specialist_mapping`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `dept_specialist_mapping` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `dept` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '部门名称',
  `specialist_id` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '专员ID',
  `specialist_name` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '专员姓名',
  `specialist_role` varchar(30) COLLATE utf8mb4_bin NOT NULL COMMENT '专员角色：dept_specialist-部门开源治理专员，usage_dept_specialist-使用部门专员，mgmt_dept_specialist-管理部门专员',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否有效：1-有效，0-无效',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dept_role_specialist` (`dept`,`specialist_role`,`specialist_id`),
  KEY `idx_specialist_id` (`specialist_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1676 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='部门开源治理专员映射表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `governance_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `governance_task` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_no` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '任务编号(如TASK-2024-0001)',
  `type` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '任务类型: introduce-引入, upgrade-升级, withdraw-退出, defect-缺陷处置',
  `action_type` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '业务类型：introduce(引入), upgrade(升级), withdraw(退出), usage(使用), risk_report(风险上报)',
  `title` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '任务标题',
  `applicant` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '申请人',
  `applicant_id` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '申请人ID',
  `applicant_dept` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '申请部门',
  `status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'draft' COMMENT '状态: draft-草稿, approving-审批中, withdrawn-已撤回, approved-已通过, rejected-已驳回',
  `priority` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'normal' COMMENT '优先级: normal-普通, urgent-紧急',
  `current_step` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '当前审批步骤描述',
  `current_node_id` bigint DEFAULT NULL COMMENT '当前审批节点ID',
  `description` mediumtext COLLATE utf8mb4_bin COMMENT '任务描述',
  `deadline` date DEFAULT NULL COMMENT '期望完成日期',
  `assigned_to` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '当前处理人',
  `approval_data` mediumtext COLLATE utf8mb4_bin COMMENT '审批节点数据(JSON)',
  `project_id` bigint DEFAULT NULL COMMENT '关联的项目ID（开源软件）',
  `project_name` varchar(200) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '关联项目名称',
  `project_dept` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '项目所属部门',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_no` (`task_no`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_applicant` (`applicant`),
  KEY `idx_assigned_to` (`assigned_to`)
) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='治理任务主表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `license`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `license` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `_COMMENT` mediumtext COLLATE utf8mb4_bin COMMENT '详情',
  `IS_DEPRECATED` int DEFAULT '0' COMMENT '是否过期 0-否 1-是',
  `FSFLIBRE` int DEFAULT '0' COMMENT 'Free Software Foundation(FSF) Libre 0-否 1-是',
  `HEADER` mediumtext COLLATE utf8mb4_bin COMMENT '文件头说明',
  `LICENSE` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '许可证代号',
  `_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '许可证名称',
  `IS_OSIAPPROVED` int NOT NULL DEFAULT '0' COMMENT '是否OSI审核 0-否 1-是',
  `SEEALSO` mediumblob COMMENT '另请参阅',
  `TEMPLATE` mediumtext COLLATE utf8mb4_bin COMMENT '许可证模板',
  `_text` mediumtext COLLATE utf8mb4_bin COMMENT '许可证说明',
  `UUID` varchar(36) COLLATE utf8mb4_bin NOT NULL COMMENT 'UUID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `LICENSE_LICENSEID_IDX` (`LICENSE`),
  UNIQUE KEY `LICENSE_UUID_IDX` (`UUID`)
) ENGINE=InnoDB AUTO_INCREMENT=60009 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='许可证表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `licensegroup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `licensegroup` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `_name` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '分组名称',
  `RISKWEIGHT` int NOT NULL DEFAULT '0' COMMENT '风险分值',
  `UUID` varchar(36) COLLATE utf8mb4_bin NOT NULL COMMENT 'UUID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `LICENSEGROUP_UUID_IDX` (`UUID`)
) ENGINE=InnoDB AUTO_INCREMENT=120005 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='许可证分组';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `licensegroup_license`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `licensegroup_license` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `licensegroup_id` bigint NOT NULL COMMENT 'licensegroup表主键',
  `license_id` bigint NOT NULL COMMENT 'license表主键',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb3_bin DEFAULT NULL COMMENT '创建者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb3_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_licensegroup_license_idx` (`licensegroup_id`,`license_id`)
) ENGINE=InnoDB AUTO_INCREMENT=30008 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin COMMENT='许可证及许可证分组关联表';
/*!40101 SET character_set_client = @saved_cs_client */;
-- menus 表已废弃：菜单统一由「统一权限系统」按角色下发，不再存表（见 PLATFORM_DESIGN.md 第 11 章）
DROP TABLE IF EXISTS `pre_eval_cache`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pre_eval_cache` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `component_id` bigint DEFAULT NULL COMMENT '组件ID',
  `component_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '组件名称',
  `component_version` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '组件版本',
  `eval_result_json` mediumtext COLLATE utf8mb4_bin COMMENT '评估结果JSON',
  `eval_time` timestamp NOT NULL COMMENT '评估时间',
  `expire_time` timestamp NOT NULL COMMENT '过期时间(默认7天后)',
  `is_valid` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否有效: 0-无效, 1-有效',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_component` (`component_name`,`component_version`),
  KEY `idx_expire_time` (`expire_time`),
  KEY `idx_is_valid` (`is_valid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='预评估缓存表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `project` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ACTIVE` int DEFAULT NULL COMMENT '状态: 0-关闭, 1-正常',
  `AUTHOR` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '项目组织',
  `_name` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '项目名称',
  `_version` varchar(50) COLLATE utf8mb4_bin DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '组件类型',
  `license` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证',
  `PURL` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Package-URL',
  `CPE` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Common Platform Enumeration',
  `SWIDTAGID` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'SWID标识',
  `DIRECT_DEPENDENCIES` mediumtext COLLATE utf8mb4_bin COMMENT 'SBOM依赖列表',
  `LAST_BOM_IMPORTED` datetime DEFAULT NULL COMMENT '最后导入时间',
  `LAST_BOM_IMPORTED_FORMAT` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '导入工具格式',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0.000' COMMENT '最后风险评分',
  `UUID` varchar(36) COLLATE utf8mb4_bin NOT NULL COMMENT 'UUID唯一标识',
  `type` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '语言类型: java/Python/go等',
  `applicable_level` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '适用等级 A/B/C/D',
  `maint_dept` varchar(200) COLLATE utf8mb4_bin NOT NULL COMMENT '维护部门',
  `use_dept` varchar(200) COLLATE utf8mb4_bin NOT NULL COMMENT '使用部门',
  `country` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '开源来源国家',
  `category` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'software' COMMENT '类别: software-开源基础软件, tool-开源工具',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `PROJECT_UUID_IDX` (`UUID`),
  UNIQUE KEY `PROJECT_UQ_IDX` (`_group`,`_name`,`_version`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='开源基础软件表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `project_inner`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `project_inner` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `pro_code` varchar(60) COLLATE utf8mb4_bin NOT NULL COMMENT '项目编码',
  `active` tinyint(1) DEFAULT '1' COMMENT '状态: 0-关闭, 1-正常',
  `type` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '语言类型: java/Python/go等',
  `is_boundary` tinyint(1) DEFAULT '0' COMMENT '是否边界系统: 1-是, 0-否',
  `group_name` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '业务组',
  `dept_name` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '部门',
  `module_code` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '模块编码',
  `branch` varchar(80) COLLATE utf8mb4_bin NOT NULL COMMENT 'Git分支名',
  `current_ID` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '本次数据标识ID',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '最后修改人',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='自研项目表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `project_inner_pro_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `project_inner_pro_info` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `daily_period` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '时期如:2020-01-10 版本，信息同步日期',
  `pro_code` varchar(60) COLLATE utf8mb4_bin NOT NULL COMMENT '项目名称',
  `commit_id` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '线上版本号',
  `is_deleted` int NOT NULL DEFAULT '1' COMMENT '项目是否删除 0-在用|1-已删除',
  `last_release_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最新发布时间',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_period_index` (`pro_code`),
  KEY `commit_id_idx` (`commit_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='自研项目线上信息';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `project_owner`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `project_owner` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_type` varchar(32) COLLATE utf8mb4_bin NOT NULL,
  `project_id` bigint NOT NULL,
  `owner_username` varchar(64) COLLATE utf8mb4_bin NOT NULL,
  `dept_name` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `team` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `assigned_by` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `assigned_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_type_project` (`project_type`,`project_id`),
  KEY `idx_owner` (`owner_username`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
-- role_menus 表已废弃：角色-菜单映射由「统一权限系统」下发或前端本地映射（见 PLATFORM_DESIGN.md 第 11 章）
DROP TABLE IF EXISTS `specialist_resource_mapping`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `specialist_resource_mapping` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `resource_type` varchar(30) COLLATE utf8mb4_bin NOT NULL COMMENT '资源类型: project-开源基础软件, inner_project-自研项目, component-依赖包',
  `resource_id` bigint NOT NULL COMMENT '资源ID',
  `resource_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '资源名称',
  `owner` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '负责人',
  `owner_dept` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '负责人部门',
  `maintainer` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '维护人',
  `maintainer_dept` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '维护人部门',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_owner` (`owner`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='人员-资源映射表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `task_applicability_assessment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_applicability_assessment` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint NOT NULL COMMENT '关联任务ID',
  `functionality_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '功能特性: fully_meet/deviation',
  `functionality_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '功能特性偏差说明',
  `security_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '安全性',
  `security_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '安全性偏差说明',
  `reliability_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '可靠性',
  `reliability_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '可靠性偏差说明',
  `performance_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '性能效率',
  `performance_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '性能效率偏差说明',
  `compatibility_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '兼容性',
  `compatibility_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '兼容性偏差说明',
  `maintainability_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '可维护性',
  `maintainability_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '可维护性偏差说明',
  `extensibility_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '可扩展性',
  `extensibility_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '可扩展性偏差说明',
  `usability_status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '易用性',
  `usability_deviation` mediumtext COLLATE utf8mb4_bin COMMENT '易用性偏差说明',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_id` (`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='适用性评估表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `task_defect`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_defect` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint NOT NULL COMMENT '关联任务主表ID',
  `risk_description` mediumtext COLLATE utf8mb4_bin,
  `internal_defect_no` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '内部缺陷编号（系统自动生成 DF-xxx）',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_id` (`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='缺陷处置扩展表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `task_introduce`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_introduce` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint NOT NULL COMMENT '关联任务主表ID',
  `component_id` bigint DEFAULT NULL COMMENT '组件ID(关联project表)',
  `component_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '组件名称',
  `component_version` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '组件版本',
  `license_type` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '许可证类型',
  `introduce_reason` mediumtext COLLATE utf8mb4_bin COMMENT '引入原因',
  `use_scenario` mediumtext COLLATE utf8mb4_bin COMMENT '使用场景',
  `project_names` mediumtext COLLATE utf8mb4_bin COMMENT '关联项目列表(JSON数组)',
  `pre_eval_cache_id` bigint DEFAULT NULL COMMENT '关联预评估缓存ID',
  `pre_eval_conclusion` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '预评估结论: pass-通过, warning-风险警告, reject-不通过',
  `pre_eval_report` mediumtext COLLATE utf8mb4_bin COMMENT '预评估报告(JSON)',
  `remark` varchar(512) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '备注说明',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_id` (`task_id`),
  KEY `idx_component_name` (`component_name`),
  KEY `idx_pre_eval_cache_id` (`pre_eval_cache_id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='引入申请扩展表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `task_upgrade`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_upgrade` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint NOT NULL COMMENT '关联任务主表ID',
  `component_id` bigint DEFAULT NULL COMMENT '组件ID(关联project表)',
  `component_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '组件名称',
  `current_version` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '当前版本',
  `target_version` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '目标版本',
  `upgrade_reason` mediumtext COLLATE utf8mb4_bin COMMENT '升级原因',
  `impact_analysis` mediumtext COLLATE utf8mb4_bin COMMENT '影响分析',
  `test_report_link` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '测试报告链接',
  `breaking_changes` mediumtext COLLATE utf8mb4_bin COMMENT '破坏性变更(JSON数组)',
  `security_fixes` mediumtext COLLATE utf8mb4_bin COMMENT '安全修复列表(JSON数组)',
  `compatibility_risk` varchar(512) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '兼容性风险',
  `rollback_plan` mediumtext COLLATE utf8mb4_bin COMMENT '回滚方案',
  `pre_eval_cache_id` bigint DEFAULT NULL COMMENT '关联预评估缓存ID',
  `pre_eval_conclusion` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '预评估结论',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_id` (`task_id`),
  KEY `idx_component_name` (`component_name`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='升级申请扩展表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `task_usage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_usage` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL,
  `component_name` varchar(255) DEFAULT '',
  `version` varchar(100) DEFAULT '',
  `is_production` int DEFAULT '0',
  `purpose` mediumtext,
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `task_withdraw`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_withdraw` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint NOT NULL COMMENT '关联任务主表ID',
  `component_id` bigint DEFAULT NULL COMMENT '组件ID(关联project表)',
  `component_name` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '组件名称',
  `component_version` varchar(100) COLLATE utf8mb4_bin DEFAULT '' COMMENT '当前版本',
  `withdraw_reason` mediumtext COLLATE utf8mb4_bin COMMENT '退出原因',
  `replacement_component` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '替代组件',
  `alternative_plan` mediumtext COLLATE utf8mb4_bin COMMENT '替代方案',
  `code_cleanup_plan` mediumtext COLLATE utf8mb4_bin COMMENT '代码清理计划',
  `affected_projects` mediumtext COLLATE utf8mb4_bin COMMENT '受影响项目(JSON数组)',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_id` (`task_id`),
  KEY `idx_component_name` (`component_name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='退出申请扩展表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_risk_report`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_risk_report` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `report_no` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '上报编号(如UR-2024-0001)',
  `reporter` varchar(50) COLLATE utf8mb4_bin NOT NULL COMMENT '上报人',
  `reporter_dept` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '上报部门',
  `title` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '风险标题',
  `risk_type` varchar(30) COLLATE utf8mb4_bin NOT NULL COMMENT '风险类型: vulnerability-安全漏洞, malware-恶意代码, license-violation-许可证违规, supply-chain-供应链风险, other-其他',
  `risk_level` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '风险等级: low, medium, high, critical',
  `affected_component` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '受影响组件',
  `affected_component_type` varchar(20) COLLATE utf8mb4_bin NOT NULL COMMENT '组件类型: project-开源基础软件, dependency-依赖包',
  `affected_version` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '受影响版本',
  `description` mediumtext COLLATE utf8mb4_bin COMMENT '详细描述',
  `attachment_url` varchar(512) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '附件URL(MinIO)',
  `related_defect_id` bigint DEFAULT NULL COMMENT '关联的缺陷ID(系统生成内部缺陷编号后)',
  `status` varchar(20) COLLATE utf8mb4_bin NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待处理, processed-已处理, rejected-已驳回',
  `reject_reason` varchar(512) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '驳回原因',
  `processed_by` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '处理人',
  `processed_time` datetime DEFAULT NULL COMMENT '处理时间',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_report_no` (`report_no`),
  KEY `idx_status` (`status`),
  KEY `idx_related_defect_id` (`related_defect_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户风险上报表';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) COLLATE utf8mb4_bin NOT NULL,
  `nickname` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL,
  `email` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `phone` varchar(32) COLLATE utf8mb4_bin DEFAULT NULL,
  `department` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `team` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL,
  `role` varchar(32) COLLATE utf8mb4_bin DEFAULT 'operator',
  `can_toggle_admin_view` tinyint(1) DEFAULT '0',
  `status` bigint DEFAULT '1',
  `create_at` datetime(3) DEFAULT NULL,
  `update_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

