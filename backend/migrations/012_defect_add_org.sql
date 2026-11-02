-- 012 外部缺陷表增加组织列（漏洞情报等使用，如 com.alibaba）
ALTER TABLE `defect` ADD COLUMN `org` varchar(128) NULL COMMENT '组织（开源软件/组件组织命名空间，如 com.alibaba）' AFTER `repair_advice`;
