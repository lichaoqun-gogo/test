-- 更新project_inner的current_commit_id和branch
UPDATE project_inner SET current_commit_id = 'a1b2c3d4e5f6', branch = 'main' WHERE id = 1;
UPDATE project_inner SET current_commit_id = 'a1b2c3d4e5f6', branch = 'main' WHERE id = 2;
UPDATE project_inner SET current_commit_id = 'a1b2c3d4e5f6', branch = 'main' WHERE id = 3;
UPDATE project_inner SET current_commit_id = 'b7c8d9e0f1a2', branch = 'main' WHERE id = 4;
UPDATE project_inner SET current_commit_id = 'b7c8d9e0f1a2', branch = 'main' WHERE id = 5;
UPDATE project_inner SET current_commit_id = 'b7c8d9e0f1a2', branch = 'main' WHERE id = 6;
UPDATE project_inner SET current_commit_id = 'c3d4e5f6a7b8', branch = 'develop' WHERE id = 7;
UPDATE project_inner SET current_commit_id = 'c3d4e5f6a7b8', branch = 'develop' WHERE id = 8;
UPDATE project_inner SET current_commit_id = 'c3d4e5f6a7b8', branch = 'develop' WHERE id = 9;

-- 为ORDER-SYS的order-service (id=1) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-101', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.google.guava/guava@32.1.3-jre', 0, 0.05, 1),
('DEP-102', 'org.apache.commons', 'commons-lang3', '3.14.0', 'Apache通用工具', 'LIBRARY', 'Apache-2.0', 'pkg:maven/org.apache.commons/commons-lang3@3.14.0', 0, 0.03, 1),
('DEP-103', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 1);

-- 为ORDER-SYS的order-api (id=2) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-201', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 2),
('DEP-202', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.16.0', 0, 0.06, 2);

-- 为ORDER-SYS的order-dao (id=3) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-301', 'org.mybatis.spring.boot', 'mybatis-spring-boot-starter', '3.0.3', 'MyBatis集成', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.mybatis.spring.boot/mybatis-spring-boot-starter@3.0.3', 0, 0.07, 3),
('DEP-302', 'mysql', 'mysql-connector-java', '8.0.33', 'MySQL驱动', 'LIBRARY', 'GPL-2.0', 'pkg:maven/mysql/mysql-connector-java@8.0.33', 0, 0.25, 3),
('DEP-303', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.google.guava/guava@32.1.3-jre', 0, 0.05, 3);

-- 为PAY-SYS的pay-service (id=4) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-401', 'com.alipay.sdk', 'alipay-sdk-java', '4.38.90.ALL', '支付宝SDK', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.alipay.sdk/alipay-sdk-java@4.38.90.ALL', 0, 0.04, 4),
('DEP-402', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 4);

-- 为PAY-SYS的pay-api (id=5) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-501', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.16.0', 0, 0.06, 5),
('DEP-502', 'org.apache.commons', 'commons-lang3', '3.14.0', 'Apache通用工具', 'LIBRARY', 'Apache-2.0', 'pkg:maven/org.apache.commons/commons-lang3@3.14.0', 0, 0.03, 5);

-- 为PAY-SYS的pay-gateway (id=6) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-601', 'org.springframework.cloud', 'spring-cloud-starter-gateway', '4.1.0', 'Spring Cloud Gateway', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.cloud/spring-cloud-starter-gateway@4.1.0', 0, 0.09, 6),
('DEP-602', 'io.netty', 'netty-all', '4.1.100.Final', 'Netty网络框架', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/io.netty/netty-all@4.1.100.Final', 0, 0.07, 6);

-- 为USER-CENTER的user-service (id=7) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-701', 'org.springframework.boot', 'spring-boot-starter-security', '3.2.0', 'Spring Security', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-security@3.2.0', 0, 0.06, 7),
('DEP-702', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 7);

-- 为USER-CENTER的user-api (id=8) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-801', 'io.jsonwebtoken', 'jjwt', '0.12.3', 'JWT库', 'LIBRARY', 'Apache-2.0', 'pkg:maven/io.jsonwebtoken/jjwt@0.12.3', 0, 0.04, 8),
('DEP-802', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.16.0', 0, 0.06, 8);

-- 为USER-CENTER的user-auth (id=9) 插入依赖
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-901', 'org.springframework.boot', 'spring-boot-starter-data-redis', '3.2.0', 'Spring Data Redis', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-data-redis@3.2.0', 0, 0.05, 9),
('DEP-902', 'redis.clients', 'jedis', '5.1.0', 'Jedis客户端', 'LIBRARY', 'MIT', 'pkg:maven/redis.clients/jedis@5.1.0', 0, 0.03, 9);

SELECT '数据插入完成' as result;
