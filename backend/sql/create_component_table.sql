-- =====================================================
-- 创建 component 表（自研项目依赖包表）
-- 与 component_snapshot 结构一致，通过 PROJECT_ID 关联 project_inner
-- =====================================================
CREATE TABLE `component` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `current_id` varchar(64) NOT NULL COMMENT '关联source_check表',
  `PARENT_COMPONENT_ID` bigint(20) DEFAULT NULL COMMENT '上层依赖',
  `AUTHOR` varchar(50) DEFAULT NULL COMMENT '作者',
  `PUBLISHER` varchar(100) DEFAULT NULL COMMENT '供应商',
  `_group` varchar(50) DEFAULT NULL COMMENT '组件组织',
  `_name` varchar(100) NOT NULL COMMENT '组件名称',
  `_version` varchar(50) DEFAULT '未知' COMMENT '版本号',
  `DESCRIPTION` varchar(1024) DEFAULT NULL COMMENT '组件描述',
  `CLASSIFIER` varchar(20) DEFAULT NULL COMMENT '组件类型',
  `FILENAME` varchar(255) DEFAULT NULL COMMENT '文件名',
  `EXTENSION` varchar(255) DEFAULT NULL COMMENT '扩展信息',
  `INTERNAL` int(1) DEFAULT NULL COMMENT '是否内部 0-否 1-是',
  `MD5` varchar(32) DEFAULT NULL COMMENT '校验信息',
  `SHA1` varchar(40) DEFAULT NULL COMMENT '校验信息',
  `SHA_256` varchar(64) DEFAULT NULL COMMENT '校验信息',
  `SHA_384` varchar(96) DEFAULT NULL COMMENT '校验信息',
  `SHA_512` varchar(128) DEFAULT NULL COMMENT '校验信息',
  `SHA3_256` varchar(64) DEFAULT NULL COMMENT '校验信息',
  `SHA3_384` varchar(96) DEFAULT NULL COMMENT '校验信息',
  `SHA3_512` varchar(128) DEFAULT NULL COMMENT '校验信息',
  `BLAKE2B_256` varchar(64) DEFAULT NULL COMMENT '校验信息',
  `BLAKE2B_384` varchar(96) DEFAULT NULL COMMENT '校验信息',
  `BLAKE2B_512` varchar(128) DEFAULT NULL COMMENT '校验信息',
  `BLAKE3` varchar(255) DEFAULT NULL COMMENT '校验信息',
  `NOTES` mediumtext DEFAULT NULL COMMENT '详细参考',
  `PURL` varchar(255) DEFAULT NULL COMMENT 'Package-URL',
  `PURLCOORDINATES` varchar(255) DEFAULT NULL COMMENT 'Package-URL Coordinates',
  `CPE` varchar(255) DEFAULT NULL COMMENT 'CPE',
  `SWIDTAGID` varchar(255) DEFAULT NULL COMMENT 'SWID',
  `DIRECT_DEPENDENCIES` mediumtext DEFAULT NULL COMMENT 'SBOM依赖列表',
  `EXTERNAL_REFERENCES` mediumblob DEFAULT NULL COMMENT '扩展说明',
  `COPYRIGHT` varchar(1024) DEFAULT NULL COMMENT '许可证说明',
  `LICENSE` varchar(255) DEFAULT NULL COMMENT '许可证',
  `LICENSE_ID` bigint(20) DEFAULT NULL COMMENT '关联license表',
  `LAST_RISKSCORE` decimal(4,3) DEFAULT '0' COMMENT '最后风险值',
  `PROJECT_ID` bigint(20) DEFAULT NULL COMMENT '关联project_inner表',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(20) DEFAULT NULL COMMENT '创建者',
  `update_time` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(20) DEFAULT NULL COMMENT '更新者',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `idx_project_id` (`PROJECT_ID`),
  KEY `idx_current_id` (`current_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='自研项目依赖包表';

-- =====================================================
-- 插入自研项目依赖包数据
-- =====================================================

-- 订单系统（project_inner id=1,2,3）的依赖包
INSERT INTO `component` (`current_id`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `INTERNAL`, `LICENSE`, `LAST_RISKSCORE`, `PROJECT_ID`, `create_by`) VALUES
('SRC-101', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 0, 'Apache-2.0', 0.120, 1, 'admin'),
('SRC-102', 'org.apache.commons', 'commons-lang3', '3.14.0', 'Apache通用工具', 'LIBRARY', 0, 'Apache-2.0', 0.080, 1, 'admin'),
('SRC-103', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web starter', 'LIBRARY', 0, 'Apache-2.0', 0.090, 2, 'admin'),
('SRC-104', 'org.mybatis', 'mybatis-spring-boot-starter', '3.0.3', 'MyBatis Starter', 'LIBRARY', 0, 'Apache-2.0', 0.070, 3, 'admin');

-- 支付系统（project_inner id=4,5,6）的依赖包
INSERT INTO `component` (`current_id`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `INTERNAL`, `LICENSE`, `LAST_RISKSCORE`, `PROJECT_ID`, `create_by`) VALUES
('SRC-105', 'com.alipay.sdk', 'alipay-sdk-java', '4.38.157.ALL', '支付宝SDK', 'LIBRARY', 0, 'Apache-2.0', 0.150, 4, 'admin'),
('SRC-106', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理库', 'LIBRARY', 0, 'Apache-2.0', 0.110, 5, 'admin'),
('SRC-107', 'org.springframework.cloud', 'spring-cloud-starter-gateway', '4.1.0', 'Spring Cloud Gateway', 'LIBRARY', 0, 'Apache-2.0', 0.130, 6, 'admin');

-- 用户中心（project_inner id=7,8,9）的依赖包
INSERT INTO `component` (`current_id`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `INTERNAL`, `LICENSE`, `LAST_RISKSCORE`, `PROJECT_ID`, `create_by`) VALUES
('SRC-108', 'org.springframework.boot', 'spring-boot-starter-security', '3.2.0', 'Spring Security', 'LIBRARY', 0, 'Apache-2.0', 0.100, 7, 'admin'),
('SRC-109', 'io.jsonwebtoken', 'jjwt', '0.12.3', 'JWT令牌库', 'LIBRARY', 0, 'Apache-2.0', 0.090, 8, 'admin'),
('SRC-110', 'org.springframework.boot', 'spring-boot-starter-data-redis', '3.2.0', 'Spring Data Redis', 'LIBRARY', 0, 'Apache-2.0', 0.080, 9, 'admin');

-- 消息平台（project_inner id=10,11）的依赖包
INSERT INTO `component` (`current_id`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `INTERNAL`, `LICENSE`, `LAST_RISKSCORE`, `PROJECT_ID`, `create_by`) VALUES
('SRC-111', 'org.apache.kafka', 'kafka-clients', '3.6.0', 'Kafka客户端', 'LIBRARY', 0, 'Apache-2.0', 0.120, 10, 'admin'),
('SRC-112', 'com.rabbitmq', 'amqp-client', '5.20.0', 'RabbitMQ客户端', 'LIBRARY', 0, 'Apache-2.0', 0.100, 11, 'admin');

-- 风控系统（project_inner id=12,13）的依赖包
INSERT INTO `component` (`current_id`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `INTERNAL`, `LICENSE`, `LAST_RISKSCORE`, `PROJECT_ID`, `create_by`) VALUES
('SRC-113', 'org.drools', 'drools-core', '8.44.0.Final', 'Drools规则引擎', 'LIBRARY', 0, 'Apache-2.0', 0.140, 12, 'admin'),
('SRC-114', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 0, 'Apache-2.0', 0.120, 13, 'admin');

-- 数据中台（project_inner id=14,15）的依赖包
INSERT INTO `component` (`current_id`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `INTERNAL`, `LICENSE`, `LAST_RISKSCORE`, `PROJECT_ID`, `create_by`) VALUES
('SRC-115', 'org.apache.flink', 'flink-core', '1.18.0', 'Flink核心库', 'LIBRARY', 0, 'Apache-2.0', 0.130, 14, 'admin'),
('SRC-116', 'org.apache.spark', 'spark-core_2.13', '3.5.0', 'Spark核心库', 'LIBRARY', 0, 'Apache-2.0', 0.140, 15, 'admin');
