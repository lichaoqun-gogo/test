-- =====================================================
-- 开源治理平台 V3 - 数据库初始化脚本
-- 兼容 TiDB
-- =====================================================

-- 1. 开源基础软件表
CREATE TABLE IF NOT EXISTS `project` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ACTIVE` int(1) DEFAULT NULL COMMENT '状态: 0-关闭, 1-正常',
  `AUTHOR` varchar(50) DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) DEFAULT NULL COMMENT '项目组织',
  `_name` varchar(50) NOT NULL COMMENT '项目名称',
  `_version` varchar(50) DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) DEFAULT NULL COMMENT '组件类型',
  `license` varchar(255) DEFAULT NULL COMMENT '许可证',
  `PURL` varchar(255) DEFAULT NULL COMMENT 'Package-URL',
  `CPE` varchar(255) DEFAULT NULL COMMENT 'Common Platform Enumeration',
  `SWIDTAGID` varchar(255) DEFAULT NULL COMMENT 'SWID标识',
  `DIRECT_DEPENDENCIES` mediumtext COMMENT 'SBOM依赖列表',
  `LAST_BOM_IMPORTED` datetime DEFAULT NULL COMMENT '最后导入时间',
  `LAST_BOM_IMPORTED_FORMAT` varchar(50) DEFAULT NULL COMMENT '导入工具格式',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0.000' COMMENT '最后风险评分',
  `UUID` varchar(36) NOT NULL COMMENT 'UUID唯一标识',
  `type` varchar(20) NOT NULL COMMENT '语言类型: java/Python/go等',
  `applicable_level` varchar(50) NOT NULL COMMENT '适用等级 A/B/C/D',
  `maint_dept` varchar(200) NOT NULL COMMENT '维护部门',
  `use_dept` varchar(200) NOT NULL COMMENT '使用部门',
  `country` varchar(50) NOT NULL COMMENT '开源来源国家',
  `category` varchar(20) NOT NULL DEFAULT 'software' COMMENT '类别: software-开源基础软件, tool-开源工具',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `PROJECT_UQ_IDX` (`_group`, `_name`, `_version`),
  UNIQUE KEY `PROJECT_UUID_IDX` (`UUID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='开源基础软件表';

-- 2. 自研项目表
CREATE TABLE IF NOT EXISTS `project_inner` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `pro_code` varchar(60) NOT NULL COMMENT '项目编码',
  `active` tinyint(1) DEFAULT '1' COMMENT '状态: 0-关闭, 1-正常',
  `type` varchar(20) NOT NULL COMMENT '语言类型: java/Python/go等',
  `is_boundary` tinyint(1) DEFAULT '0' COMMENT '是否边界系统: 1-是, 0-否',
  `group_name` varchar(20) DEFAULT NULL COMMENT '业务组',
  `dept_name` varchar(64) DEFAULT NULL COMMENT '部门',
  `module_code` varchar(64) NOT NULL COMMENT '模块编码',
  `branch` varchar(80) NOT NULL COMMENT 'Git分支名',
  `current_ID` varchar(64) NOT NULL COMMENT '本次数据标识ID',
  `current_commit_id` varchar(64) DEFAULT NULL COMMENT '当前Commit ID',
  `production_commit_id` varchar(64) DEFAULT NULL COMMENT '生产Commit ID',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_by` varchar(20) DEFAULT NULL COMMENT '最后修改人',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_pro_code` (`pro_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='自研项目表';

-- 3. 开源依赖包表
CREATE TABLE IF NOT EXISTS `component_snapshot` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `current_id` varchar(64) NOT NULL COMMENT '关联source_check表标识',
  `PARENT_COMPONENT_ID` bigint(20) DEFAULT NULL COMMENT '上层依赖组件ID',
  `AUTHOR` varchar(50) DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) DEFAULT NULL COMMENT '组件组织',
  `_name` varchar(100) NOT NULL COMMENT '组件名称',
  `_version` varchar(50) DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) DEFAULT NULL COMMENT '组件类型',
  `FILENAME` varchar(255) DEFAULT NULL COMMENT '文件名',
  `EXTENSION` varchar(255) DEFAULT NULL COMMENT '扩展信息',
  `INTERNAL` int(1) DEFAULT NULL COMMENT '是否内部: 0-否, 1-是',
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
  `LICENSE_ID` bigint(20) DEFAULT NULL COMMENT '关联license表ID',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0.000' COMMENT '最后风险评分',
  `PROJECT_ID` bigint(20) DEFAULT NULL COMMENT '关联project_inner表ID',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_current_id` (`current_id`),
  KEY `idx_project_id` (`project_id`),
  KEY `idx_name_version` (`_name`, `_version`),
  KEY `idx_license_id` (`license_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='开源依赖包表';

-- 4. 黑名单配置表
CREATE TABLE IF NOT EXISTS `component_core` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `_format` tinyint(2) NOT NULL DEFAULT '1' COMMENT '类型: maven:1、pypi:2、npm:3、yum:4等',
  `type` varchar(50) DEFAULT NULL COMMENT '类型',
  `_name` varchar(50) NOT NULL COMMENT '包名',
  `_group` varchar(50) NOT NULL COMMENT '包分组',
  `version_exp` varchar(200) NOT NULL COMMENT '匹配表达式',
  `applicable_level` varchar(50) NOT NULL COMMENT '适用等级',
  `support_team` varchar(50) NOT NULL COMMENT '技术支持团队',
  `reason` varchar(400) DEFAULT NULL COMMENT '添加理由',
  `advise` varchar(400) DEFAULT NULL COMMENT '建议',
  `_enable` tinyint(1) NOT NULL COMMENT '是否启用, 1: 启用, 0: 禁用',
  `alert_type` tinyint(2) DEFAULT '1' COMMENT '告警类型 1 阻断 2警告',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  `level_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '风险等级生效时间',
  `is_boundary` tinyint(1) DEFAULT '2' COMMENT '2:全部 1:边界系统 0非边界系统',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uniq_black_list_idx` (`_format`,`_name`,`_group`,`applicable_level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=240001 COMMENT='黑名单配置表';

-- 5. 黑名单阻断记录表
CREATE TABLE IF NOT EXISTS `component_black` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `blacklist_id` bigint(20) NOT NULL COMMENT '命中的黑名单规则ID(关联component_core表)',
  `component_name` varchar(100) NOT NULL COMMENT '组件名称',
  `component_version` varchar(50) NOT NULL COMMENT '组件版本',
  `component_group` varchar(50) DEFAULT NULL COMMENT '组件组织',
  `project_id` bigint(20) DEFAULT NULL COMMENT '关联的自研项目ID',
  `project_name` varchar(100) DEFAULT NULL COMMENT '关联的自研项目名称',
  `alert_type` tinyint(2) NOT NULL DEFAULT '1' COMMENT '告警类型: 1-阻断, 2-警告',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '处理状态: pending-未处理, processing-处理中, resolved-处理完成, false_positive-误报, audited-已审核',
  `handler` varchar(50) DEFAULT NULL COMMENT '处理人',
  `handle_plan` varchar(512) DEFAULT NULL COMMENT '处理计划',
  `handle_result` mediumtext COMMENT '处理结果说明',
  `attachment_url` varchar(512) DEFAULT NULL COMMENT '附件URL(MinIO)',
  `auditor` varchar(50) DEFAULT NULL COMMENT '审核人',
  `audit_opinion` varchar(512) DEFAULT NULL COMMENT '审核意见',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `ci_pipeline_id` varchar(100) DEFAULT NULL COMMENT 'CI流水线ID',
  `ci_build_id` varchar(100) DEFAULT NULL COMMENT 'CI构建ID',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_blacklist_id` (`blacklist_id`),
  KEY `idx_status` (`status`),
  KEY `idx_component` (`component_name`, `component_version`),
  KEY `idx_project_id` (`project_id`),
  KEY `idx_ci_pipeline` (`ci_pipeline_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='黑名单阻断记录表';

-- 6. 许可证表
CREATE TABLE IF NOT EXISTS `license` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `_COMMENT` mediumtext DEFAULT NULL COMMENT '详情',
  `IS_DEPRECATED` int(1) DEFAULT '0' COMMENT '是否过期 0-否 1-是',
  `FSFLIBRE` int(1) DEFAULT '0' COMMENT 'Free Software Foundation(FSF) Libre 0-否 1-是',
  `HEADER` mediumtext DEFAULT NULL COMMENT '文件头说明',
  `LICENSE` varchar(255) NOT NULL COMMENT '许可证代号',
  `_name` varchar(255) NOT NULL COMMENT '许可证名称',
  `IS_OSIAPPROVED` int(1) NOT NULL DEFAULT '0' COMMENT '是否OSI审核 0-否 1-是',
  `SEEALSO` mediumblob DEFAULT NULL COMMENT '另请参阅',
  `TEMPLATE` mediumtext DEFAULT NULL COMMENT '许可证模板',
  `_text` mediumtext DEFAULT NULL COMMENT '许可证说明',
  `UUID` varchar(36) NOT NULL COMMENT 'UUID',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `LICENSE_LICENSEID_IDX` (`LICENSE`),
  UNIQUE KEY `LICENSE_UUID_IDX` (`UUID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=60002 COMMENT='许可证表';

-- 7. 许可证分组表
CREATE TABLE IF NOT EXISTS `licensegroup` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `_name` varchar(50) NOT NULL COMMENT '分组名称',
  `RISKWEIGHT` int(3) NOT NULL DEFAULT '0' COMMENT '风险分值',
  `UUID` varchar(36) NOT NULL COMMENT 'UUID',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `LICENSEGROUP_UUID_IDX` (`UUID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=120001 COMMENT='许可证分组';

-- 8. 许可证分组关联表
CREATE TABLE IF NOT EXISTS `licensegroup_license` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `licensegroup_id` bigint(20) NOT NULL COMMENT 'licensegroup表主键',
  `license_id` bigint(20) NOT NULL COMMENT 'license表主键',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uniq_licensegroup_license_idx` (`licensegroup_id`,`license_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin AUTO_INCREMENT=30001 COMMENT='许可证及许可证分组关联表';

-- 9. 缺陷台账表
CREATE TABLE IF NOT EXISTS `defect_ledger` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `defect_no` varchar(50) NOT NULL COMMENT '缺陷编号(如DF-2024-0001)',
  `source` varchar(30) NOT NULL COMMENT '来源: sca-scan-SCA扫描, vuln-intel-漏洞情报, user-report-用户上报',
  `title` varchar(255) NOT NULL COMMENT '缺陷标题',
  `risk_type` varchar(30) NOT NULL COMMENT '风险类型: vulnerability-安全漏洞, malware-恶意代码, license-violation-许可证违规, supply-chain-供应链风险, other-其他',
  `risk_level` varchar(20) NOT NULL COMMENT '风险等级: low, medium, high, critical',
  `affected_component` varchar(100) NOT NULL COMMENT '受影响组件',
  `affected_component_type` varchar(20) NOT NULL COMMENT '组件类型: project-开源基础软件, dependency-依赖包',
  `affected_version` varchar(50) DEFAULT NULL COMMENT '受影响版本',
  `cve_id` varchar(50) DEFAULT NULL COMMENT 'CVE编号',
  `description` mediumtext COMMENT '详细描述',
  `discovered_date` datetime NOT NULL COMMENT '发现日期',
  `status` varchar(20) NOT NULL DEFAULT 'pending_analysis' COMMENT '状态: pending_analysis-待研判, pending_review-待上会, assigned-已分配, ignored-已忽略, closed-已闭环',
  `analyze_result` varchar(500) DEFAULT NULL COMMENT '专员研判结果',
  `review_meeting_no` varchar(50) DEFAULT NULL COMMENT '评审会议编号',
  `review_conclusion` varchar(20) DEFAULT NULL COMMENT '评审结论: no_action-无需处置, need_fix-需要修复',
  `dispatch_type` varchar(20) DEFAULT NULL COMMENT '分发类型: dependency-依赖包缺陷, project-开源基础软件缺陷',
  `dispatch_target` varchar(50) DEFAULT NULL COMMENT '分发目标人员',
  `dispatch_dept` varchar(100) DEFAULT NULL COMMENT '分发目标部门',
  `dispatch_time` datetime DEFAULT NULL COMMENT '分发时间',
  `dispatcher` varchar(20) DEFAULT NULL COMMENT '分发人',
  `related_task_id` bigint(20) DEFAULT NULL COMMENT '关联的处置任务ID',
  `analyze_by` varchar(50) DEFAULT NULL COMMENT '研判人',
  `analyze_time` datetime DEFAULT NULL COMMENT '研判时间',
  `review_by` varchar(50) DEFAULT NULL COMMENT '评审人',
  `review_time` datetime DEFAULT NULL COMMENT '评审时间',
  `close_time` datetime DEFAULT NULL COMMENT '闭环时间',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_defect_no` (`defect_no`),
  KEY `idx_status` (`status`),
  KEY `idx_risk_level` (`risk_level`),
  KEY `idx_source` (`source`),
  KEY `idx_affected_component` (`affected_component`),
  KEY `idx_related_task_id` (`related_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='缺陷台账';

-- 10. 缺陷-任务映射表
CREATE TABLE IF NOT EXISTS `defect_task_mapping` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `defect_id` bigint(20) NOT NULL COMMENT '缺陷ID',
  `defect_no` varchar(50) NOT NULL COMMENT '缺陷编号',
  `task_id` bigint(20) NOT NULL COMMENT '关联任务ID',
  `task_no` varchar(50) NOT NULL COMMENT '任务编号',
  `mapping_type` varchar(20) NOT NULL COMMENT '映射类型: auto-自动, manual-手动',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_defect_id` (`defect_id`),
  KEY `idx_task_id` (`task_id`),
  UNIQUE KEY `uk_defect_task` (`defect_id`, `task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='缺陷-任务映射表';

-- 11. 治理任务主表
CREATE TABLE IF NOT EXISTS `governance_task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_no` varchar(50) NOT NULL COMMENT '任务编号(如TASK-2024-0001)',
  `type` varchar(20) NOT NULL COMMENT '任务类型: introduce-引入, upgrade-升级, withdraw-退出, defect-缺陷处置',
  `action_type` varchar(50) DEFAULT NULL COMMENT '业务类型：introduce(引入), upgrade(升级), withdraw(退出), usage(使用), risk_report(风险上报)',
  `title` varchar(255) NOT NULL COMMENT '任务标题',
  `applicant` varchar(50) NOT NULL COMMENT '申请人',
  `applicant_id` varchar(50) DEFAULT NULL COMMENT '申请人ID',
  `applicant_dept` varchar(100) NOT NULL COMMENT '申请部门',
  `status` varchar(20) NOT NULL DEFAULT 'draft' COMMENT '状态: draft-草稿, approving-审批中, withdrawn-已撤回, approved-已通过, rejected-已驳回',
  `priority` varchar(20) NOT NULL DEFAULT 'normal' COMMENT '优先级: normal-普通, urgent-紧急',
  `current_step` varchar(100) NOT NULL DEFAULT '' COMMENT '当前审批步骤描述',
  `current_node_id` bigint(20) DEFAULT NULL COMMENT '当前审批节点ID',
  `description` mediumtext COMMENT '任务描述',
  `deadline` date DEFAULT NULL COMMENT '期望完成日期',
  `assigned_to` varchar(100) DEFAULT NULL COMMENT '当前处理人',
  `approval_data` mediumtext COMMENT '审批节点数据(JSON)',
  `project_id` bigint(20) DEFAULT NULL COMMENT '关联的项目ID（开源软件）',
  `project_name` varchar(200) DEFAULT NULL COMMENT '关联项目名称',
  `project_dept` varchar(100) DEFAULT NULL COMMENT '项目所属部门',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_task_no` (`task_no`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_applicant` (`applicant`),
  KEY `idx_assigned_to` (`assigned_to`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='治理任务主表';

-- 12. 引入申请扩展表
CREATE TABLE IF NOT EXISTS `task_introduce` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint(20) NOT NULL COMMENT '关联任务主表ID',
  `component_id` bigint(20) DEFAULT NULL COMMENT '组件ID(关联project表)',
  `component_name` varchar(100) NOT NULL COMMENT '组件名称',
  `component_version` varchar(50) NOT NULL COMMENT '组件版本',
  `license_type` varchar(100) DEFAULT NULL COMMENT '许可证类型',
  `introduce_reason` mediumtext COMMENT '引入原因',
  `use_scenario` mediumtext COMMENT '使用场景',
  `project_names` mediumtext COMMENT '关联项目列表(JSON数组)',
  `pre_eval_cache_id` bigint(20) DEFAULT NULL COMMENT '关联预评估缓存ID',
  `pre_eval_conclusion` varchar(20) DEFAULT NULL COMMENT '预评估结论: pass-通过, warning-风险警告, reject-不通过',
  `pre_eval_report` mediumtext COMMENT '预评估报告(JSON)',
  `eval_report_url` varchar(512) DEFAULT NULL COMMENT '评估报告MinIO路径',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_task_id` (`task_id`),
  KEY `idx_component_name` (`component_name`),
  KEY `idx_pre_eval_cache_id` (`pre_eval_cache_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='引入申请扩展表';

-- 13. 升级申请扩展表
CREATE TABLE IF NOT EXISTS `task_upgrade` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint(20) NOT NULL COMMENT '关联任务主表ID',
  `component_id` bigint(20) DEFAULT NULL COMMENT '组件ID(关联project表)',
  `component_name` varchar(100) NOT NULL COMMENT '组件名称',
  `current_version` varchar(50) NOT NULL COMMENT '当前版本',
  `target_version` varchar(50) NOT NULL COMMENT '目标版本',
  `upgrade_reason` mediumtext COMMENT '升级原因',
  `impact_analysis` mediumtext COMMENT '影响分析',
  `test_report_link` varchar(255) DEFAULT NULL COMMENT '测试报告链接',
  `breaking_changes` mediumtext COMMENT '破坏性变更(JSON数组)',
  `security_fixes` mediumtext COMMENT '安全修复列表(JSON数组)',
  `compatibility_risk` varchar(512) DEFAULT NULL COMMENT '兼容性风险',
  `rollback_plan` mediumtext COMMENT '回滚方案',
  `pre_eval_cache_id` bigint(20) DEFAULT NULL COMMENT '关联预评估缓存ID',
  `pre_eval_conclusion` varchar(20) DEFAULT NULL COMMENT '预评估结论',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_task_id` (`task_id`),
  KEY `idx_component_name` (`component_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='升级申请扩展表';

-- 14. 退出申请扩展表
CREATE TABLE IF NOT EXISTS `task_withdraw` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint(20) NOT NULL COMMENT '关联任务主表ID',
  `component_id` bigint(20) DEFAULT NULL COMMENT '组件ID(关联project表)',
  `component_name` varchar(100) NOT NULL COMMENT '组件名称',
  `withdraw_reason` mediumtext COMMENT '退出原因',
  `replacement_component` varchar(255) DEFAULT NULL COMMENT '替代组件',
  `alternative_plan` mediumtext COMMENT '替代方案',
  `code_cleanup_plan` mediumtext COMMENT '代码清理计划',
  `affected_projects` mediumtext COMMENT '受影响项目(JSON数组)',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_task_id` (`task_id`),
  KEY `idx_component_name` (`component_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='退出申请扩展表';

-- 15. 缺陷处置扩展表
CREATE TABLE IF NOT EXISTS `task_defect` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint(20) NOT NULL COMMENT '关联任务主表ID',
  `defect_id` bigint(20) NOT NULL COMMENT '关联缺陷ID',
  `defect_no` varchar(50) NOT NULL COMMENT '缺陷编号',
  `defect_type` varchar(30) NOT NULL COMMENT '缺陷类型',
  `affected_component` varchar(100) NOT NULL COMMENT '受影响组件',
  `fix_plan` mediumtext COMMENT '修复方案',
  `need_fix` tinyint(1) DEFAULT NULL COMMENT '是否需要处置: 1-是, 0-否',
  `fix_reason` varchar(512) DEFAULT NULL COMMENT '不需要处置原因',
  `production_change_no` varchar(50) DEFAULT NULL COMMENT '生产变更号',
  `production_ticket_text` mediumtext COMMENT '工单文本',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_task_id` (`task_id`),
  KEY `idx_defect_id` (`defect_id`),
  KEY `idx_defect_no` (`defect_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='缺陷处置扩展表';

-- 16. 审批节点记录表
CREATE TABLE IF NOT EXISTS `approval_node` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint(20) NOT NULL COMMENT '关联任务ID',
  `node_name` varchar(100) NOT NULL COMMENT '节点名称(如使用部门意见)',
  `node_order` int(11) NOT NULL COMMENT '节点顺序',
  `approver` varchar(50) DEFAULT NULL COMMENT '审批人',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待审批, approved-通过, rejected-驳回, not_required-不需要',
  `opinion` varchar(1024) DEFAULT NULL COMMENT '审批意见',
  `approved_at` datetime DEFAULT NULL COMMENT '审批时间',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_task_id` (`task_id`),
  KEY `idx_approver` (`approver`),
  KEY `idx_task_order` (`task_id`, `node_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批节点记录表';

-- 17. 审批历史表
CREATE TABLE IF NOT EXISTS `approval_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` bigint(20) NOT NULL COMMENT '关联任务ID',
  `task_no` varchar(50) NOT NULL COMMENT '任务编号',
  `node_id` bigint(20) NOT NULL COMMENT '审批节点ID',
  `node_name` varchar(100) NOT NULL COMMENT '节点名称',
  `approver` varchar(50) NOT NULL COMMENT '审批人',
  `action` varchar(20) NOT NULL COMMENT '操作: approve-通过, reject-驳回',
  `opinion` varchar(1024) DEFAULT NULL COMMENT '审批意见',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_task_id` (`task_id`),
  KEY `idx_task_no` (`task_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批历史表';

-- 18. 预评估缓存表
CREATE TABLE IF NOT EXISTS `pre_eval_cache` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `component_id` bigint(20) DEFAULT NULL COMMENT '组件ID',
  `component_name` varchar(100) NOT NULL COMMENT '组件名称',
  `component_version` varchar(50) NOT NULL COMMENT '组件版本',
  `eval_result_json` mediumtext COMMENT '评估结果JSON',
  `eval_time` timestamp NOT NULL COMMENT '评估时间',
  `expire_time` timestamp NOT NULL COMMENT '过期时间(默认7天后)',
  `is_valid` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否有效: 0-无效, 1-有效',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_component` (`component_name`, `component_version`),
  KEY `idx_expire_time` (`expire_time`),
  KEY `idx_is_valid` (`is_valid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='预评估缓存表';

-- 19. 用户风险上报表
CREATE TABLE IF NOT EXISTS `user_risk_report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `report_no` varchar(50) NOT NULL COMMENT '上报编号(如UR-2024-0001)',
  `reporter` varchar(50) NOT NULL COMMENT '上报人',
  `reporter_dept` varchar(100) NOT NULL COMMENT '上报部门',
  `title` varchar(255) NOT NULL COMMENT '风险标题',
  `risk_type` varchar(30) NOT NULL COMMENT '风险类型: vulnerability-安全漏洞, malware-恶意代码, license-violation-许可证违规, supply-chain-供应链风险, other-其他',
  `risk_level` varchar(20) NOT NULL COMMENT '风险等级: low, medium, high, critical',
  `affected_component` varchar(100) NOT NULL COMMENT '受影响组件',
  `affected_component_type` varchar(20) NOT NULL COMMENT '组件类型: project-开源基础软件, dependency-依赖包',
  `affected_version` varchar(50) DEFAULT NULL COMMENT '受影响版本',
  `description` mediumtext COMMENT '详细描述',
  `attachment_url` varchar(512) DEFAULT NULL COMMENT '附件URL(MinIO)',
  `related_defect_id` bigint(20) DEFAULT NULL COMMENT '关联的缺陷ID(系统生成内部缺陷编号后)',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending-待处理, processed-已处理, rejected-已驳回',
  `reject_reason` varchar(512) DEFAULT NULL COMMENT '驳回原因',
  `processed_by` varchar(50) DEFAULT NULL COMMENT '处理人',
  `processed_time` datetime DEFAULT NULL COMMENT '处理时间',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_report_no` (`report_no`),
  KEY `idx_status` (`status`),
  KEY `idx_related_defect_id` (`related_defect_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户风险上报表';

-- 20. 附件表
CREATE TABLE IF NOT EXISTS `attachment` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `file_name` varchar(255) NOT NULL COMMENT '文件名',
  `file_size` bigint(20) DEFAULT NULL COMMENT '文件大小(字节)',
  `file_type` varchar(100) DEFAULT NULL COMMENT '文件类型(MIME)',
  `minio_path` varchar(512) NOT NULL COMMENT 'MinIO存储路径',
  `minio_bucket` varchar(100) DEFAULT 'osm-mngt' COMMENT 'MinIO存储桶',
  `related_type` varchar(50) NOT NULL COMMENT '关联类型: task-任务, defect-缺陷, blacklist-黑名单, report-报告',
  `related_id` bigint(20) DEFAULT NULL COMMENT '关联记录ID',
  `uploader` varchar(50) NOT NULL COMMENT '上传人',
  `description` varchar(512) DEFAULT NULL COMMENT '文件描述',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_related` (`related_type`, `related_id`),
  KEY `idx_uploader` (`uploader`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='附件表';

-- 21. 人员-资源映射表
CREATE TABLE IF NOT EXISTS `specialist_resource_mapping` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `resource_type` varchar(30) NOT NULL COMMENT '资源类型: project-开源基础软件, inner_project-自研项目, component-依赖包',
  `resource_id` bigint(20) NOT NULL COMMENT '资源ID',
  `resource_name` varchar(100) NOT NULL COMMENT '资源名称',
  `owner` varchar(50) NOT NULL COMMENT '负责人',
  `owner_dept` varchar(100) DEFAULT NULL COMMENT '负责人部门',
  `maintainer` varchar(50) DEFAULT NULL COMMENT '维护人',
  `maintainer_dept` varchar(100) DEFAULT NULL COMMENT '维护人部门',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_resource` (`resource_type`, `resource_id`),
  KEY `idx_owner` (`owner`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='人员-资源映射表';
