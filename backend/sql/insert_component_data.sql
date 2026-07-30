-- 为Redis(project_id=1)插入依赖数据
INSERT INTO component_snapshot (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-001', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 'Apache-2.0', 0.120, 1),
('DEP-002', 'org.apache.commons', 'commons-lang3', '3.14.0', 'Apache通用工具', 'LIBRARY', 'Apache-2.0', 0.080, 1),
('DEP-003', 'io.netty', 'netty-all', '4.1.100.Final', '网络通信框架', 'LIBRARY', 'Apache-2.0', 0.110, 1),
('DEP-004', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理库', 'LIBRARY', 'Apache-2.0', 0.150, 1),
('DEP-005', 'org.slf4j', 'slf4j-api', '2.0.9', '日志门面', 'LIBRARY', 'MIT', 0.050, 1),
('DEP-006', 'ch.qos.logback', 'logback-classic', '1.4.14', '日志实现', 'LIBRARY', 'EPL-1.0', 0.060, 1),
('DEP-007', 'org.apache.commons', 'commons-pool2', '2.12.0', '对象池', 'LIBRARY', 'Apache-2.0', 0.070, 1),
('DEP-008', 'redis.clients', 'jedis', '5.1.0', 'Redis客户端', 'LIBRARY', 'MIT', 0.090, 1);

-- 为Spring Boot(project_id=2)插入依赖数据
INSERT INTO component_snapshot (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-009', 'org.springframework', 'spring-core', '6.1.1', 'Spring核心', 'LIBRARY', 'Apache-2.0', 0.090, 2),
('DEP-010', 'org.springframework.boot', 'spring-boot-starter', '3.2.0', 'Spring Boot Starter', 'LIBRARY', 'Apache-2.0', 0.080, 2),
('DEP-011', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 0.150, 2),
('DEP-012', 'org.apache.tomcat.embed', 'tomcat-embed-core', '10.1.16', '内嵌Tomcat', 'LIBRARY', 'Apache-2.0', 0.120, 2),
('DEP-013', 'org.hibernate.validator', 'hibernate-validator', '8.0.1.Final', '参数校验', 'LIBRARY', 'Apache-2.0', 0.100, 2),
('DEP-014', 'com.zaxxer', 'HikariCP', '5.1.0', '连接池', 'LIBRARY', 'Apache-2.0', 0.070, 2),
('DEP-015', 'org.yaml', 'snakeyaml', '2.2', 'YAML解析', 'LIBRARY', 'Apache-2.0', 0.060, 2);

-- 为MySQL(project_id=3)插入依赖数据
INSERT INTO component_snapshot (current_id, _group, _name, _version, DESCRIPTION, CLASSIFIER, LICENSE, LAST_RISKSCORE, PROJECT_ID) VALUES
('DEP-016', 'com.google.protobuf', 'protobuf-java', '3.25.1', '序列化', 'LIBRARY', 'BSD-3-Clause', 0.130, 3),
('DEP-017', 'io.grpc', 'grpc-netty', '1.60.0', 'RPC框架', 'LIBRARY', 'Apache-2.0', 0.140, 3),
('DEP-018', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 'Apache-2.0', 0.120, 3),
('DEP-019', 'org.slf4j', 'slf4j-api', '2.0.9', '日志门面', 'LIBRARY', 'MIT', 0.050, 3),
('DEP-020', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理', 'LIBRARY', 'Apache-2.0', 0.150, 3);
