





-- 期刊表
CREATE TABLE `periodicals` (
                               `periodical_id` int NOT NULL AUTO_INCREMENT COMMENT '期刊ID',
                               `name` varchar(255) NOT NULL COMMENT '期刊名字',
                               `publication_time` varchar(255) NOT NULL COMMENT '出刊时间',
                               `updated_periodical_at` int NOT NULL COMMENT '更新时间',
                               `competent_unit` varchar(255) DEFAULT NULL COMMENT '主管单位',
                               `organization` varchar(255) DEFAULT NULL COMMENT '主办单位',
                               `domestic_number` varchar(50) DEFAULT NULL COMMENT '国内刊号',
                               `international_number` varchar(50) DEFAULT NULL COMMENT '国际刊号',
                               `founding_time` varchar(50) DEFAULT NULL COMMENT '创刊时间',
                               `citation_rate` varchar(50) DEFAULT NULL COMMENT '引用率',
                               `periodical_description` text COMMENT '期刊描述',
                               `column_setting` text COMMENT '栏目设置',
                               `composite_influence_factor` float DEFAULT NULL COMMENT '复合影响因子',
                               `synthetic_influence_factor` float DEFAULT NULL COMMENT '综合影响因子',
                               `audit_time` varchar(50) DEFAULT NULL COMMENT '审核周期',
                               `invoice_receipt` tinyint(1) DEFAULT NULL COMMENT '是否可开杂志社发票',
                               `attention_matter` text COMMENT '文章注意事项',
                               `article_naming` text COMMENT '文章命名',
                               `internal_processes` text COMMENT '内部流程',
                               `price` float DEFAULT NULL COMMENT '期刊发表服务费',
                               `is_warp` tinyint(1) DEFAULT NULL COMMENT '是否全包',
                               `periodical_page` varchar(24) DEFAULT NULL COMMENT '期刊页码',
                               `periodical_batch` int DEFAULT NULL COMMENT '期刊批次',
                               `works` tinyint(1) DEFAULT NULL COMMENT '可发作品',
                               `created_at` int DEFAULT NULL COMMENT '创建时间',
                               `updated_at` int DEFAULT NULL COMMENT '更新时间',
                               PRIMARY KEY (`periodical_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci



-- 通用分类表
CREATE TABLE `categorizes` (
                               `categorize_id` int NOT NULL AUTO_INCREMENT COMMENT '主键，唯一标识方向',
                               `parent_id` bigint DEFAULT NULL COMMENT '父级方向 ID，支持树状结构',
                               `name` varchar(255) NOT NULL COMMENT '名称，例如方向名称、级别名称等',
                               `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态，0-禁用，1-启用',
                               `app_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '应用状态，备用字段',
                               `type` tinyint(1) NOT NULL COMMENT '类型：1-收稿方向，2-期刊级别，3-网站收录，4-出版周期',
                               `sort` int NOT NULL DEFAULT '0' COMMENT '排序值，用于前端排序',
                               `link` varchar(255) DEFAULT NULL COMMENT '相关链接，备用字段',
                               `extra_info` varchar(255) DEFAULT NULL COMMENT '额外信息，备用字段',
                               `created_at` int DEFAULT NULL COMMENT '创建时间',
                               `updated_at` int DEFAULT NULL COMMENT '更新时间',
                               PRIMARY KEY (`categorize_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='收稿方向表';


-- 期刊表和通用分类表 关联
CREATE TABLE `periodical_categorize` (
                                         `categorize_id` int NOT NULL,
                                         `periodical_id` int NOT NULL,
                                         PRIMARY KEY (`categorize_id`,`periodical_id`),
                                         KEY `fk_periodical_categorize` (`periodical_id`),
                                         CONSTRAINT `fk_periodical_categorize_categorize` FOREIGN KEY (`categorize_id`) REFERENCES `categorizes` (`categorize_id`),
                                         CONSTRAINT `fk_periodical_categorize_periodical` FOREIGN KEY (`periodical_id`) REFERENCES `periodicals` (`periodical_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- 创建用户表
CREATE TABLE `users` (
     `id` int unsigned NOT NULL AUTO_INCREMENT,
     `created_at` int NOT NULL COMMENT '创建时间',
     `updated_at` int NOT NULL COMMENT '更新时间',
     `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名, 用户名不允许重复的',
     `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '不能保持用户的明文密码',
     `label` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户标签',
     `role` tinyint NOT NULL COMMENT '用户的角色',
     PRIMARY KEY (`id`) USING BTREE,
     UNIQUE KEY `idx_user` (`username`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


CREATE TABLE `tokens` (
                          `created_at` int NOT NULL COMMENT '创建时间',
                          `updated_at` int NOT NULL COMMENT '更新时间',
                          `user_id` int NOT NULL COMMENT '用户的Id',
                          `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名, 用户名不允许重复的',
                          `access_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户的访问令牌',
                          `access_token_expired_at` int NOT NULL COMMENT '令牌过期时间',
                          `refresh_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '刷新令牌',
                          `refresh_token_expired_at` int NOT NULL COMMENT '刷新令牌过期时间',
                          `is_enable` bigint NOT NULL COMMENT 'Token是否有效;0无效,1有效',
                          PRIMARY KEY (`access_token`) USING BTREE,
                          UNIQUE KEY `idx_token` (`access_token`) USING BTREE,
                          INDEX  `idx_userid_username` (`user_id`,`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 收稿方向 type 1
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('艺术类', 1, 1, 0, 2);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('经济类', 1, 1, 0, 3);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('文学类', 1, 1, 0, 1);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('综合类', 1, 1, 0, 4);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('社会学/政治/党建', 1, 1, 0, 5);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('旅游类', 1, 1, 0, 6);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('基础教育类', 1, 1, 0, 7);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('职教高教类', 1, 1, 0, 8);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('思想教育', 1, 1, 0, 9);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('新闻/档案/文化', 1, 1, 0, 10);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('法律类', 1, 1, 0, 11);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('体育类', 1, 1, 0, 12);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('农业/畜禽/园艺/食品', 1, 1, 0, 13);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('医学类', 1, 1, 0, 14);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('计算机与信息通信', 1, 1, 0, 15);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('水利水电/电力/地质/自然科学', 1, 1, 0, 16);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('建筑建材类/机械/工程技术/交通', 1, 1, 0, 17);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('工业化工/科技类', 1, 1, 0, 18);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('低价刊物', 1, 1, 0, 19);




-- 期刊级别 type 2
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('省级', 2, 1, 0, 101);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('国家级', 2, 1, 0, 102);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('学报', 2, 1, 0, 103);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('国际期刊', 2, 1, 0, 104);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('会议', 2, 1, 0, 105);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('报纸', 2, 1, 0, 106);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('专利', 2, 1, 0, 107);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('软著', 2, 1, 0, 108);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('图书出版', 2, 1, 0, 109);




-- 网站收录 type 3
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('知网', 3, 1, 0, 201);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('维普', 3, 1, 0, 202);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('万方', 3, 1, 0, 203);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('AMI', 3, 1, 0, 204);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('SCD', 3, 1, 0, 205);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('RCCSE', 3, 1, 0, 206);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('中刊网/期刊网', 3, 1, 0, 207);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('龙源新媒体', 3, 1, 0, 208);



-- 出版周期 type 4
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('半月刊及以下', 4, 1, 0, 301);
INSERT INTO `categorizes` (`name`, `type`, `status`, `app_status`, `sort`) VALUES ('月刊及以上', 4, 1, 0, 302);


-- ==========================================
-- 动态角色与权限管理 (RBAC) 表结构及初始数据
-- ==========================================

-- 1. 角色表
CREATE TABLE IF NOT EXISTS `roles` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(64) NOT NULL COMMENT '角色名称',
  `code` varchar(64) NOT NULL COMMENT '角色唯一编码',
  `description` varchar(255) DEFAULT '' COMMENT '角色描述',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否内置角色：1-是（不可删除），0-否',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` int NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_at` int NOT NULL DEFAULT 0 COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色表';

-- 2. 权限资源表（菜单 + 按钮/接口操作）
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `parent_id` int unsigned NOT NULL DEFAULT 0 COMMENT '父级权限ID，0为顶级',
  `name` varchar(64) NOT NULL COMMENT '权限名称',
  `code` varchar(100) NOT NULL COMMENT '权限唯一标识码，如 periodical:create',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '类型：1-菜单/目录，2-按钮/接口操作',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `created_at` int NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_at` int NOT NULL DEFAULT 0 COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_perm_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='权限资源表';

-- 3. 角色-权限关联表（多对多）
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `role_id` int unsigned NOT NULL COMMENT '角色ID',
  `permission_id` int unsigned NOT NULL COMMENT '权限ID',
  PRIMARY KEY (`role_id`, `permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='角色权限关联表';

-- 4. 用户-角色关联表（多对多，支持单用户绑定多个角色）
CREATE TABLE IF NOT EXISTS `user_roles` (
  `user_id` int unsigned NOT NULL COMMENT '用户ID',
  `role_id` int unsigned NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`, `role_id`),
  KEY `idx_user_roles_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户角色关联表';

-- 初始化默认角色（兼容原枚举：1-管理员，2-普通用户，3-超级管理员）
INSERT IGNORE INTO `roles` (`id`, `name`, `code`, `description`, `status`, `is_system`, `sort`, `created_at`, `updated_at`) VALUES
(1, '管理员', 'admin', '业务管理员，负责期刊管理及用户日常维护', 1, 0, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, '普通用户', 'author', '普通用户，可查询、浏览期刊详情及复制期刊描述', 1, 0, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, '超级管理员', 'super_admin', '系统最高权限管理员，拥有系统全部权限', 1, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 初始化预设权限树
INSERT IGNORE INTO `permissions` (`id`, `parent_id`, `name`, `code`, `type`, `sort`, `created_at`, `updated_at`) VALUES
(1, 0, '首页看板', 'dashboard:view', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 0, '期刊管理', 'periodical:view', 1, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 2, '查询与查看详情', 'periodical:query', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 2, '新增期刊', 'periodical:create', 2, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(5, 2, '编辑期刊', 'periodical:update', 2, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(6, 2, '删除期刊', 'periodical:delete', 2, 4, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(7, 2, '复制期刊描述', 'periodical:copy', 2, 5, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(8, 0, '系统管理', 'system:view', 1, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(9, 8, '用户管理', 'user:view', 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(10, 9, '查询用户列表', 'user:query', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(11, 9, '新增用户', 'user:create', 2, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(12, 9, '编辑用户/重置密码', 'user:update', 2, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(13, 9, '删除用户', 'user:delete', 2, 4, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(14, 8, '角色与权限管理', 'role:view', 1, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(15, 14, '查询角色与权限树', 'role:query', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(16, 14, '新增角色', 'role:create', 2, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(17, 14, '编辑角色与分配权限', 'role:update', 2, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(18, 14, '删除角色', 'role:delete', 2, 4, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 为普通用户 (role_id=2) 分配默认权限：首页、期刊查看、查询详情、复制描述
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(2, 1), (2, 2), (2, 3), (2, 7);

-- 为管理员 (role_id=1) 分配默认权限：首页、期刊全权限、用户管理全权限、角色查询权限
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7),
(1, 8), (1, 9), (1, 10), (1, 11), (1, 12), (1, 13),
(1, 14), (1, 15);

-- 为超级管理员 (role_id=3) 分配所有权限
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(3, 1), (3, 2), (3, 3), (3, 4), (3, 5), (3, 6), (3, 7),
(3, 8), (3, 9), (3, 10), (3, 11), (3, 12), (3, 13),
(3, 14), (3, 15), (3, 16), (3, 17), (3, 18);

-- 迁移存量用户的单角色字段到 user_roles 多对多关联表
INSERT IGNORE INTO `user_roles` (`user_id`, `role_id`)
SELECT `id`, CASE WHEN `role` IN (1, 2, 3) THEN `role` ELSE 2 END FROM `users`;

-- ==========================================
-- AI Agent 智能助手与多模型管理表结构及权限
-- ==========================================

INSERT IGNORE INTO `permissions` (`id`, `parent_id`, `name`, `code`, `type`, `sort`, `created_at`, `updated_at`) VALUES
(19, 0, 'AI 智能工作台', 'ai:workspace', 1, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(20, 19, 'AI 对话与智能助手', 'ai:chat', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(21, 8, 'AI 模型与审计管理', 'ai:config', 1, 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(22, 21, '全员对话审计查看', 'ai:audit', 2, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(1, 19), (1, 20),
(2, 19), (2, 20),
(3, 19), (3, 20), (3, 21), (3, 22);

CREATE TABLE IF NOT EXISTS `ai_model_configs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '渠道显示名称',
  `provider` varchar(50) NOT NULL COMMENT '供应商标识：qwen/deepseek/openai/custom',
  `base_url` varchar(255) NOT NULL COMMENT 'API Base URL',
  `api_key` varchar(255) DEFAULT '' COMMENT 'API密钥',
  `model_name` varchar(100) NOT NULL COMMENT '模型标识',
  `system_prompt` text COMMENT '自定义系统提示词',
  `temperature` decimal(3,2) NOT NULL DEFAULT 0.20 COMMENT '采样温度',
  `max_tokens` int NOT NULL DEFAULT 4096 COMMENT '最大Token数',
  `is_default` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否全局默认：1-是，0-否',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  `remark` varchar(255) DEFAULT '' COMMENT '备注',
  `created_at` int NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_at` int NOT NULL DEFAULT 0 COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI大模型渠道配置表';

CREATE TABLE IF NOT EXISTS `ai_chat_sessions` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `session_id` varchar(64) NOT NULL COMMENT '会话唯一UUID',
  `user_id` int unsigned NOT NULL COMMENT '所属用户ID',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '所属用户名',
  `title` varchar(200) NOT NULL DEFAULT '新对话' COMMENT '会话标题',
  `scene` varchar(32) NOT NULL DEFAULT 'workspace' COMMENT '来源场景：workspace/copilot',
  `model_name` varchar(100) DEFAULT '' COMMENT '使用模型名',
  `message_count` int NOT NULL DEFAULT 0 COMMENT '消息总数',
  `is_pinned` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否置顶',
  `created_at` int NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_at` int NOT NULL DEFAULT 0 COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_session_id` (`session_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI对话会话表';

CREATE TABLE IF NOT EXISTS `ai_chat_messages` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `session_id` varchar(64) NOT NULL COMMENT '会话UUID',
  `user_id` int unsigned NOT NULL COMMENT '所属用户ID',
  `username` varchar(64) NOT NULL DEFAULT '' COMMENT '所属用户名',
  `role` varchar(32) NOT NULL COMMENT '角色：user/assistant',
  `content` longtext COMMENT '消息文本内容',
  `model_name` varchar(100) DEFAULT '' COMMENT '响应模型名称',
  `tool_steps` json DEFAULT NULL COMMENT '工具调用步骤记录',
  `action_card` json DEFAULT NULL COMMENT '结构化交互卡片与确认状态',
  `created_at` int NOT NULL DEFAULT 0 COMMENT '创建时间',
  `updated_at` int NOT NULL DEFAULT 0 COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_msg_session_id` (`session_id`),
  KEY `idx_msg_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI对话消息记录表';


