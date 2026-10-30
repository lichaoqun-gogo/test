-- 缺陷采集台账（统一模型，GitHub/Jira 多通道，upsert 导入）
-- defect_id = 软件小写#编号（redis#13401 / kafka#KAFKA-20449），作为 upsert 业务键
CREATE TABLE IF NOT EXISTS `defect` (
  `defect_id`      VARCHAR(100)  NOT NULL COMMENT '唯一缺陷ID：软件小写#编号（redis#13401 / kafka#KAFKA-20449），upsert 业务键',
  `software`       VARCHAR(100)  NOT NULL COMMENT '软件名称（Redis / Kafka）',
  `issue_url`      VARCHAR(500)  NOT NULL COMMENT '缺陷原文链接（GitHub issue / ASF Jira）',
  `title`          VARCHAR(1000) NOT NULL COMMENT '缺陷标题',
  `body`           MEDIUMTEXT    NULL     COMMENT '缺陷正文（GitHub/Jira 原文，含复现步骤等）',
  `labels`         JSON          NULL     COMMENT '原始标签数组（GitHub labels / Jira labels）',
  `reporter`       VARCHAR(200)  NULL     COMMENT '报告人',
  `created_at`     DATETIME      NULL     COMMENT '缺陷创建时间（UTC）',
  `updated_at`     DATETIME      NULL     COMMENT '缺陷更新时间（UTC），增量游标依据',
  `comments_count` INT           NOT NULL DEFAULT 0 COMMENT '评论数',
  `source`         VARCHAR(20)   NOT NULL DEFAULT 'github-issue' COMMENT '来源：github-issue / jira',
  `severity`       VARCHAR(10)   NOT NULL DEFAULT '未知' COMMENT '风险定级：高危/中危/低危/提示/未知',
  `issue_state`    VARCHAR(10)   NOT NULL DEFAULT 'open' COMMENT '缺陷状态：open(未修复) / closed(已修复)',
  `state_reason`   VARCHAR(20)   NULL     COMMENT '关闭原因：completed(已修复) / not_planned(不计划修=隐患) 等',
  `reactions`      JSON          NULL     COMMENT '社区关注度 {total, up(+1), ...}',
  `milestone`      VARCHAR(100)  NULL     COMMENT '官方版本规划/修复版本（Jira 取 fixVersions 首个）',
  `watch_matches`  JSON          NULL     COMMENT '命中观测版本数组（用户传入的 --watch）',
  `summary_zh`     TEXT          NULL     COMMENT '中文概述（LLM 增强，可为空）',
  `repair_advice`  TEXT          NULL     COMMENT '修复建议（LLM 增强，可为空）',
  `imported_at`    DATETIME      NOT NULL COMMENT '本次导入时间',
  PRIMARY KEY (`defect_id`),
  KEY `idx_software_state` (`software`, `issue_state`),
  KEY `idx_severity` (`severity`),
  KEY `idx_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
  COMMENT='缺陷采集台账（统一模型，GitHub/Jira 多通道，upsert 导入）';
