





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





-- 创建用户表
