-- Adminer 4.2.5 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP DATABASE IF EXISTS `anytunnel_base`;
CREATE DATABASE `anytunnel_base` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `anytunnel_base`;

DROP TABLE IF EXISTS `at_area`;
CREATE TABLE `at_area` (
  `area_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) NOT NULL COMMENT '区域名称',
  `cs_type` enum('server','client') NOT NULL COMMENT '类型,server或者client',
  `is_forbidden` tinyint(4) NOT NULL COMMENT '是否禁止,1:禁止0:允许',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `update_time` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`area_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='server和client登录区域黑白名单';


DROP TABLE IF EXISTS `at_client`;
CREATE TABLE `at_client` (
  `client_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'client主键ID',
  `name` varchar(30) NOT NULL COMMENT 'client名称',
  `token` varchar(32) NOT NULL COMMENT 'client的token,最多32个字符',
  `user_id` varchar(32) NOT NULL COMMENT 'client所属用户ID,0代表是系统的client',
  `local_host` varchar(50) NOT NULL COMMENT '需要暴漏的本地网络host，client可以访问的ip域名都可以',
  `local_port` int(11) NOT NULL COMMENT '需要暴漏的本地网络host的端口',
  `create_time` int(11) NOT NULL COMMENT 'client创建时间',
  `update_time` int(11) NOT NULL COMMENT 'client更新时间',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除,1是,0否',
  PRIMARY KEY (`client_id`),
  UNIQUE KEY `token` (`token`),
  UNIQUE KEY `user_id_token` (`user_id`,`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='client表';


DROP TABLE IF EXISTS `at_cluster`;
CREATE TABLE `at_cluster` (
  `cluster_id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `region_id` int(11) NOT NULL COMMENT '区域主键ID',
  `name` char(30) NOT NULL DEFAULT '' COMMENT '名称',
  `ip` char(15) NOT NULL DEFAULT '' COMMENT 'cluster的ip',
  `system_conn_number` int(10) NOT NULL DEFAULT '0' COMMENT '系统连接数',
  `tunnel_conn_number` int(10) NOT NULL DEFAULT '0' COMMENT '隧道连接数',
  `bandwidth` int(10) NOT NULL DEFAULT '0' COMMENT '带宽',
  `create_time` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `is_disable` tinyint(1) NOT NULL COMMENT '是否禁用，0否，1是',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0 否 1 是',
  PRIMARY KEY (`cluster_id`),
  UNIQUE KEY `ip` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='cluster机器表';


DROP TABLE IF EXISTS `at_conn`;
CREATE TABLE `at_conn` (
  `tunnel_id` int(11) NOT NULL COMMENT '用户的隧道主键ID',
  `user_id` int(11) NOT NULL COMMENT '用户主键ID',
  `cluster_id` int(11) NOT NULL COMMENT 'cluster主键ID',
  `server_id` int(11) NOT NULL COMMENT 'server主键ID',
  `count` int(11) NOT NULL COMMENT '用户的隧道的连接数',
  `upload` int(11) NOT NULL COMMENT '上传速度,字节/秒',
  `download` int(11) NOT NULL COMMENT '下载速度,字节/秒',
  `update_time` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`tunnel_id`),
  KEY `count` (`count`),
  KEY `user_id_count` (`user_id`,`count`),
  KEY `server_id` (`server_id`),
  KEY `upload` (`upload`),
  KEY `download` (`download`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='连接数表';


DROP TABLE IF EXISTS `at_ip_list`;
CREATE TABLE `at_ip_list` (
  `ip_list_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `ip` char(15) NOT NULL COMMENT 'IP地址',
  `is_forbidden` tinyint(1) NOT NULL COMMENT '0禁止,1允许',
  `cs_type` enum('server','client') NOT NULL COMMENT 'cs类型,server或client',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `update_time` int(11) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`ip_list_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='cs登录IP黑白名单表';


DROP TABLE IF EXISTS `at_online`;
CREATE TABLE `at_online` (
  `online_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '在线主键ID',
  `user_id` int(11) NOT NULL COMMENT '用户表主键ID',
  `cs_id` int(11) NOT NULL COMMENT 'server或者client的id',
  `cluster_id` int(11) NOT NULL COMMENT 'cluster的主键ID',
  `cs_ip` char(15) NOT NULL COMMENT 'server或者client的外网IP',
  `cs_type` enum('server','client') NOT NULL COMMENT '类型,server或client',
  `create_time` int(11) NOT NULL COMMENT '上线时间',
  PRIMARY KEY (`online_id`),
  UNIQUE KEY `token_type` (`cs_id`,`cs_type`),
  KEY `ip` (`cluster_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='server或client在线表';


DROP TABLE IF EXISTS `at_package`;
CREATE TABLE `at_package` (
  `package_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '流量包主键ID',
  `user_id` int(11) NOT NULL COMMENT '用户主键ID',
  `bytes_total` bigint(20) NOT NULL COMMENT '流量包大小,单位字节',
  `bytes_left` bigint(20) NOT NULL COMMENT '剩余流量,单位字节',
  `start_time` int(11) NOT NULL COMMENT '生效时间',
  `end_time` int(11) NOT NULL COMMENT '过期时间',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `update_time` int(11) NOT NULL COMMENT '更新时间',
  `comment` char(10) NOT NULL COMMENT '流量包来源说明,最多10个字符',
  PRIMARY KEY (`package_id`),
  KEY `user_id_bytes_left` (`user_id`,`bytes_left`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户流量包表';


DROP TABLE IF EXISTS `at_region`;
CREATE TABLE `at_region` (
  `region_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '区域主键ID',
  `name` char(30) NOT NULL COMMENT '区域名称',
  `parent_id` int(11) NOT NULL COMMENT '上级区域的region_id,顶级区域为0',
  `create_time` int(11) NOT NULL COMMENT '区域创建时间',
  `update_time` int(11) NOT NULL COMMENT '区域修改时间',
  `is_delete` tinyint(1) NOT NULL COMMENT '是否删除,1是.0否',
  PRIMARY KEY (`region_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='cluster区域表';

INSERT INTO `at_region` (`region_id`, `name`, `parent_id`, `create_time`, `update_time`, `is_delete`) VALUES
(1,	'香港',	0,	0,	1503657697,	0),
(2,	'香港免费1区',	1,	0,	1503907304,	0),
(3,	'新加坡',	0,	1503655210,	1503657714,	0),
(4,	'新加坡免费1区',	3,	1503656729,	1503907288,	0),
(5,	'大陆',	0,	1503658393,	0,	0),
(7,	'大陆VIP电信',	5,	1503658507,	1503907337,	0),
(8,	'香港VIP1区',	1,	1503674982,	1503907795,	0),
(9,	'大陆VIP联通',	5,	1503907177,	1503907343,	0),
(10,	'大陆VIP多线',	5,	1503907196,	1503907398,	0),
(11,	'大陆VIP移动',	5,	1503907217,	1503907782,	0),
(12,	'大陆免费1区',	5,	1503907253,	1503907771,	0),
(13,	'新加坡VIP1区',	3,	1503907816,	1503907824,	0),
(14,	'香港高速VIP1区',	1,	1503907869,	1503909105,	0),
(15,	'香港高速VIP2区',	1,	1503907873,	0,	0),
(16,	'日本',	0,	1503909138,	0,	0),
(17,	'日本免费1区',	16,	1503909142,	0,	0),
(18,	'日本VIP1区',	16,	1503909151,	0,	0),
(19,	'美国',	0,	1503909156,	0,	0),
(20,	'美国免费1区',	19,	1503909162,	1503909170,	0),
(21,	'美国VIP1区',	19,	1503909195,	0,	0);

DROP TABLE IF EXISTS `at_role`;
CREATE TABLE `at_role` (
  `role_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '角色主键ID',
  `name` char(30) NOT NULL COMMENT '角色名称',
  `bandwidth` bigint(20) NOT NULL COMMENT '带宽限制,单位字节,0不限制',
  `server_area` enum('china','foreign','all') NOT NULL DEFAULT 'china' COMMENT 'server允许登录的区域,china:中国,foreign:国外,all:全部',
  `client_area` enum('china','foreign','all') NOT NULL DEFAULT 'china' COMMENT 'client允许登录的区域,china:中国,foreign:国外,all:全部',
  `tunnel_mode` char(5) NOT NULL COMMENT '允许的隧道模式0:普通1:高级2:特殊',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `update_time` int(11) NOT NULL COMMENT '修改时间',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除,1是,0否',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色表';

INSERT INTO `at_role` (`role_id`, `name`, `bandwidth`, `server_area`, `client_area`, `tunnel_mode`, `create_time`, `update_time`, `is_delete`) VALUES
(1,	'默认组',	100000,	'china',	'china',	'0,2,1',	0,	1504602526,	0),
(2,	'VIP1',	500000,	'china',	'china',	'0,1,2',	0,	1504509792,	0);

DROP TABLE IF EXISTS `at_role_region`;
CREATE TABLE `at_role_region` (
  `role_region_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '角色和区域关系主键ID',
  `role_id` int(11) NOT NULL COMMENT '角色主键ID',
  `region_id` int(11) NOT NULL COMMENT '区域主键ID',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`role_region_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色和区域关系表';

INSERT INTO `at_role_region` (`role_region_id`, `role_id`, `region_id`, `create_time`) VALUES
(35,	1,	7,	1504249580),
(36,	1,	9,	1504249580),
(37,	1,	10,	1504249580),
(38,	1,	11,	1504249580),
(39,	1,	12,	1504249580),
(40,	1,	17,	1504249580),
(41,	1,	18,	1504249580);

DROP TABLE IF EXISTS `at_server`;
CREATE TABLE `at_server` (
  `server_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'server主键ID',
  `name` varchar(30) NOT NULL COMMENT 'server名称',
  `token` varchar(32) NOT NULL COMMENT 'server的token,最多32个字符',
  `user_id` int(11) NOT NULL COMMENT 'server所属用户ID,0代表是系统的server',
  `create_time` int(11) NOT NULL COMMENT 'server创建时间',
  `update_time` int(11) NOT NULL COMMENT 'server更新时间',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除,1是,0否',
  PRIMARY KEY (`server_id`),
  UNIQUE KEY `token` (`token`),
  UNIQUE KEY `user_id_token` (`user_id`,`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='server表';


DROP TABLE IF EXISTS `at_tunnel`;
CREATE TABLE `at_tunnel` (
  `tunnel_id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `mode` tinyint(1) NOT NULL COMMENT '模式，0普通模式，1高级模式，2特殊模式',
  `name` char(255) NOT NULL DEFAULT '' COMMENT '隧道名称',
  `user_id` int(10) NOT NULL DEFAULT '0' COMMENT '隧道所属用户主键ID',
  `cluster_id` int(10) NOT NULL COMMENT 'cluster主键ID',
  `server_id` int(11) NOT NULL DEFAULT '0' COMMENT 'server主键ID',
  `client_id` int(11) NOT NULL COMMENT 'client主键ID',
  `protocol` tinyint(1) NOT NULL COMMENT '协议 1:TCP 2:UDP',
  `server_listen_port` int(10) NOT NULL DEFAULT '0' COMMENT 'server 监听端口',
  `server_listen_ip` char(100) NOT NULL DEFAULT '' COMMENT 'server bind ip',
  `client_local_port` int(10) NOT NULL DEFAULT '0' COMMENT 'client local port',
  `client_local_host` char(100) NOT NULL DEFAULT '' COMMENT 'client localhost',
  `status` tinyint(1) NOT NULL COMMENT '隧道状态,0异常,1:正常',
  `is_open` tinyint(1) NOT NULL COMMENT '是否已经打开,0否,1是,打开过的无论是否异常,必须手动一次关闭',
  `create_time` int(10) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) NOT NULL DEFAULT '0' COMMENT '修改时间',
  `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0 否 1 是',
  PRIMARY KEY (`tunnel_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='隧道表';


DROP TABLE IF EXISTS `at_user`;
CREATE TABLE `at_user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户主键ID',
  `nickname` char(20) NOT NULL COMMENT '昵称',
  `username` char(20) NOT NULL COMMENT '用户登录名',
  `password` char(32) NOT NULL COMMENT '密码',
  `email` char(50) NOT NULL COMMENT 'email地址',
  `is_active` tinyint(1) NOT NULL COMMENT 'email是否激活',
  `is_forbidden` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否禁用用户,1是.0否',
  `forbidden_reason` varchar(50) NOT NULL COMMENT '禁用的原因描述',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `update_time` int(11) NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `username` (`username`),
  KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';


DROP TABLE IF EXISTS `at_user_role`;
CREATE TABLE `at_user_role` (
  `role_user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户和角色关系主键ID',
  `role_id` int(11) NOT NULL COMMENT '角色主键ID',
  `user_id` int(11) NOT NULL COMMENT '用户主键ID',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `update_time` int(11) NOT NULL,
  PRIMARY KEY (`role_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户和角色关系表';


-- 2017-09-08 09:15:00