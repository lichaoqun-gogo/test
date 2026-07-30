-- 清理重复数据
DELETE FROM component WHERE PROJECT_ID <= 9;

-- ORDER-SYS的order-service (id=1)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('a1b2c3d4e5f6', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.google.guava/guava@32.1.3-jre', 0, 0.05, 1),
('a1b2c3d4e5f6', 'org.apache.commons', 'commons-lang3', '3.14.0', 'Apache通用工具', 'LIBRARY', 'Apache-2.0', 'pkg:maven/org.apache.commons/commons-lang3@3.14.0', 0, 0.03, 1),
('a1b2c3d4e5f6', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 1);

-- ORDER-SYS的order-api (id=2)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('a1b2c3d4e5f6', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 2),
('a1b2c3d4e5f6', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.16.0', 0, 0.06, 2);

-- ORDER-SYS的order-dao (id=3)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('a1b2c3d4e5f6', 'org.mybatis.spring.boot', 'mybatis-spring-boot-starter', '3.0.3', 'MyBatis集成', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.mybatis.spring.boot/mybatis-spring-boot-starter@3.0.3', 0, 0.07, 3),
('a1b2c3d4e5f6', 'mysql', 'mysql-connector-java', '8.0.33', 'MySQL驱动', 'LIBRARY', 'GPL-2.0', 'pkg:maven/mysql/mysql-connector-java@8.0.33', 0, 0.25, 3),
('a1b2c3d4e5f6', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.google.guava/guava@32.1.3-jre', 0, 0.05, 3);

-- PAY-SYS的pay-service (id=4)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('b7c8d9e0f1a2', 'com.alipay.sdk', 'alipay-sdk-java', '4.38.90.ALL', '支付宝SDK', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.alipay.sdk/alipay-sdk-java@4.38.90.ALL', 0, 0.04, 4),
('b7c8d9e0f1a2', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 4);

-- PAY-SYS的pay-api (id=5)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('b7c8d9e0f1a2', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.16.0', 0, 0.06, 5),
('b7c8d9e0f1a2', 'org.apache.commons', 'commons-lang3', '3.14.0', 'Apache通用工具', 'LIBRARY', 'Apache-2.0', 'pkg:maven/org.apache.commons/commons-lang3@3.14.0', 0, 0.03, 5);

-- PAY-SYS的pay-gateway (id=6)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('b7c8d9e0f1a2', 'org.springframework.cloud', 'spring-cloud-starter-gateway', '4.1.0', 'Spring Cloud Gateway', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.cloud/spring-cloud-starter-gateway@4.1.0', 0, 0.09, 6),
('b7c8d9e0f1a2', 'io.netty', 'netty-all', '4.1.100.Final', 'Netty网络框架', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/io.netty/netty-all@4.1.100.Final', 0, 0.07, 6);

-- USER-CENTER的user-service (id=7)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('c3d4e5f6a7b8', 'org.springframework.boot', 'spring-boot-starter-security', '3.2.0', 'Spring Security', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-security@3.2.0', 0, 0.06, 7),
('c3d4e5f6a7b8', 'org.springframework.boot', 'spring-boot-starter-web', '3.2.0', 'Spring Boot Web', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-web@3.2.0', 0, 0.08, 7);

-- USER-CENTER的user-api (id=8)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('c3d4e5f6a7b8', 'io.jsonwebtoken', 'jjwt', '0.12.3', 'JWT库', 'LIBRARY', 'Apache-2.0', 'pkg:maven/io.jsonwebtoken/jjwt@0.12.3', 0, 0.04, 8),
('c3d4e5f6a7b8', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 'pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.16.0', 0, 0.06, 8);

-- USER-CENTER的user-auth (id=9)
INSERT INTO component (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, PURL, INTERNAL, LAST_RISKSCORE, PROJECT_ID) VALUES
('c3d4e5f6a7b8', 'org.springframework.boot', 'spring-boot-starter-data-redis', '3.2.0', 'Spring Data Redis', 'FRAMEWORK', 'Apache-2.0', 'pkg:maven/org.springframework.boot/spring-boot-starter-data-redis@3.2.0', 0, 0.05, 9),
('c3d4e5f6a7b8', 'redis.clients', 'jedis', '5.1.0', 'Jedis客户端', 'LIBRARY', 'MIT', 'pkg:maven/redis.clients/jedis@5.1.0', 0, 0.03, 9);

SELECT '数据更新完成' as result;
