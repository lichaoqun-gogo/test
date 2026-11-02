-- 给 project 表增加 category 字段，区分开源基础软件和开源工具
-- category: software-开源基础软件, tool-开源工具
ALTER TABLE `project` ADD COLUMN `category` varchar(20) NOT NULL DEFAULT 'software' COMMENT '类别: software-开源基础软件, tool-开源工具' AFTER `country`;

-- 现有数据默认设为 software
UPDATE `project` SET `category` = 'software' WHERE `category` = '' OR `category` IS NULL;
