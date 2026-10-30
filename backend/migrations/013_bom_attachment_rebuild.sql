-- 013 统一文件管理：重建 attachment 表（参考项目结构，兼容生命周期治理流程附件复用）+ 新建 bom 表
-- 设计说明：
--   - attachment 统一文件管理：source/filename/md5/store_url，唯一键 (source, filename, store_url) 防重复上传；
--     升级/退出等生命周期流程附件复用同一张表，source 区分业务场景（见 AttachmentSource 枚举）
--   - bom 只存元数据，store_url 引用 attachment.id（文件统一走 attachment）
--   - MinIO 目录规则：{source}/{组织}/{名称}/{版本}/{文件名}（BOM 即 bom/项目组织/项目名称/项目版本/xxx.xml）

DROP TABLE IF EXISTS `bom`;
DROP TABLE IF EXISTS `attachment`;

CREATE TABLE `attachment` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `source` varchar(20) NOT NULL COMMENT '附件来源（AttachmentSource 枚举：BOM_UPLOAD/VULDB_PLAN/VULDB_RESULT/FIRST_VERIFY/FINAL_VERIFY）',
  `filename` varchar(64) NOT NULL COMMENT '文件名',
  `md5` varchar(64) NOT NULL COMMENT '文件md5值',
  `description` mediumtext DEFAULT NULL COMMENT '描述和备注',
  `store_url` varchar(150) NOT NULL COMMENT '对象存储中的相对地址',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_attachment_idx` (`source`,`filename`,`store_url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='附件管理';

CREATE TABLE `bom` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `bom_format` varchar(20) DEFAULT NULL COMMENT 'sbom工具类型',
  `bom_version` int DEFAULT NULL COMMENT 'bom生成版本',
  `spec_version` varchar(20) DEFAULT NULL COMMENT 'sbom工具版本',
  `bom_create` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'bom创建时间',
  `serial_number` varchar(36) DEFAULT NULL COMMENT 'sbom文件uuid',
  `store_url` bigint DEFAULT NULL COMMENT 'sbom文件|附件管理attachmentid',
  `project_id` bigint NOT NULL COMMENT '关联项目',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`),
  KEY `idx_bom_project` (`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='BOM清单';
