-- 014 组件快照定时清理（每小时）：归档 component_snapshot -> component_snapshot_history
-- 设计说明：
--   - component_snapshot_history 与 component_snapshot 结构一致（无 is_delete 逻辑删除列，含 create_time）
--   - commit_id_history 记录每次归档的 commit_id（唯一键，防重复归档）
--   - project_inner_pro_info 线上依赖信息（is_deleted=0 在用，其 commit_id 不清理）
--   - is_delete / create_time 字段的添加在 database.go 的 Go 迁移函数中做（幂等守卫），本文件只建表

-- 组件快照历史表（= component_snapshot 去掉 is_delete，另含 create_time）
CREATE TABLE IF NOT EXISTS `component_snapshot_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `current_id` varchar(64) NOT NULL COMMENT '关联source_check表标识',
  `PARENT_COMPONENT_ID` bigint DEFAULT NULL COMMENT '上层依赖组件ID',
  `AUTHOR` varchar(50) DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) DEFAULT NULL COMMENT '组件组织',
  `_name` varchar(100) NOT NULL COMMENT '组件名称',
  `_version` varchar(50) DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) DEFAULT NULL COMMENT '组件类型',
  `FILENAME` varchar(255) DEFAULT NULL COMMENT '文件名',
  `EXTENSION` varchar(255) DEFAULT NULL COMMENT '扩展信息',
  `INTERNAL` int DEFAULT NULL COMMENT '是否内部: 0-否, 1-是',
  `MD5` varchar(32) DEFAULT NULL COMMENT 'MD5校验值',
  `SHA1` varchar(40) DEFAULT NULL COMMENT 'SHA1校验值',
  `SHA_256` varchar(64) DEFAULT NULL COMMENT 'SHA256校验值',
  `SHA_384` varchar(96) DEFAULT NULL COMMENT 'SHA384校验值',
  `SHA_512` varchar(128) DEFAULT NULL COMMENT 'SHA512校验值',
  `SHA3_256` varchar(64) DEFAULT NULL COMMENT 'SHA3-256校验值',
  `SHA3_384` varchar(96) DEFAULT NULL COMMENT 'SHA3-384校验值',
  `SHA3_512` varchar(128) DEFAULT NULL COMMENT 'SHA3-512校验值',
  `BLAKE2B_256` varchar(64) DEFAULT NULL COMMENT 'BLAKE2B-256校验值',
  `BLAKE2B_384` varchar(96) DEFAULT NULL COMMENT 'BLAKE2B-384校验值',
  `BLAKE2B_512` varchar(128) DEFAULT NULL COMMENT 'BLAKE2B-512校验值',
  `BLAKE3` varchar(64) DEFAULT NULL COMMENT 'BLAKE3校验值',
  `NOTES` mediumtext COMMENT '详细参考',
  `PURL` varchar(255) DEFAULT NULL COMMENT 'Package-URL',
  `PURLCOORDINATES` varchar(255) DEFAULT NULL COMMENT 'Package-URL Coordinates',
  `CPE` varchar(255) DEFAULT NULL COMMENT 'CPE标识',
  `SWIDTAGID` varchar(255) DEFAULT NULL COMMENT 'SWID标识',
  `DIRECT_DEPENDENCIES` mediumtext COMMENT 'SBOM依赖列表',
  `EXTERNAL_REFERENCES` mediumblob COMMENT '扩展说明',
  `COPYRIGHT` varchar(1024) DEFAULT NULL COMMENT '许可证说明',
  `LICENSE` varchar(255) DEFAULT NULL COMMENT '许可证',
  `LICENSE_ID` bigint DEFAULT NULL COMMENT '关联license表ID',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0.000' COMMENT '最后风险评分',
  `PROJECT_ID` bigint DEFAULT NULL COMMENT '关联project_inner表ID',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_current_id` (`current_id`),
  KEY `idx_project_id` (`PROJECT_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='开源依赖包历史归档表';

-- 历史 commit 记录（记录每次归档的 commit_id，唯一键防重复）
CREATE TABLE IF NOT EXISTS `commit_id_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `commit_id` varchar(64) NOT NULL COMMENT '提交commitId',
  `create_time_history` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '历史创建时间',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '后台归档时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_commit_id` (`commit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='历史commit记录';

-- 自研项目线上信息（is_deleted=0 在用：其 commit_id 为线上版本，不清理）
CREATE TABLE IF NOT EXISTS `project_inner_pro_info` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `daily_period` varchar(20) NOT NULL COMMENT '时期如:2020-01-10 版本，信息同步日期',
  `pro_code` varchar(60) NOT NULL COMMENT '项目名称',
  `commit_id` varchar(64) NOT NULL COMMENT '线上版本号',
  `is_deleted` int NOT NULL DEFAULT '1' COMMENT '项目是否删除 0-在用|1-已删除',
  `last_release_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '最新发布时间',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_period_index` (`pro_code`),
  KEY `commit_id_idx` (`commit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='自研项目线上信息';
