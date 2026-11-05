-- 016 外部缺陷中文内容字段补充
-- 背景：缺陷采集台账 defect 表新增中文标题 title_zh（此前由 GORM AutoMigrate 运行时建列，此处补独立增量迁移便于新环境按增量序列执行）
-- 1) defect 表补 title_zh：外部缺陷 kafka-all.json 对齐，存中文标题，用于详情双 tab（基本信息 + 内容详情）
ALTER TABLE `defect` ADD COLUMN `title_zh` varchar(1000) DEFAULT NULL COMMENT '中文标题（导入 kafka-all.json 的 title_zh 字段）' AFTER `imported_at`;
-- 2) defect_ledger 表补 affected_component：漏洞裂变按受影响组件拆分处置单（组件如 fastjson / 基础软件如 redis），该列此前被误删，已由 Go 侧幂等迁移恢复，此处补独立 SQL 记录
ALTER TABLE `defect_ledger` ADD COLUMN `affected_component` varchar(128) DEFAULT NULL COMMENT '受影响组件（漏洞裂变拆分依据）' AFTER `description`;
