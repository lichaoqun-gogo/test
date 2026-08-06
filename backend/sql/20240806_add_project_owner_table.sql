-- 版本号: 20240806
-- 变更说明: 新增项目-负责人关联表 (project_owner)
-- 变更日期: 2024-08-06
-- 变更内容:
--   1. 新增 project_owner 表，支持开源基础软件和自研项目的负责人关联
--   2. 支持多负责人，记录分配人和分配时间
--   3. 冗余部门和团队信息，方便查询

CREATE TABLE IF NOT EXISTS `project_owner` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `project_type` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT '项目类型: software-开源基础软件, inner-自研项目',
  `project_id` bigint NOT NULL COMMENT '项目ID',
  `owner_username` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '负责人用户名',
  `dept_name` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '部门',
  `team` varchar(128) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '团队',
  `assigned_by` varchar(64) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '分配人',
  `assigned_at` timestamp NULL DEFAULT NULL COMMENT '分配时间',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_project_type` (`project_type`),
  KEY `idx_project_id` (`project_id`),
  KEY `idx_owner_username` (`owner_username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='项目-负责人关联表';
