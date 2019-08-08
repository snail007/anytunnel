-- Adminer 4.2.5 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP DATABASE IF EXISTS `anytunnel_system`;
CREATE DATABASE `anytunnel_system` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `anytunnel_system`;

DROP TABLE IF EXISTS `system_privilege`;
CREATE TABLE `system_privilege` (
  `privilege_id` int(10) NOT NULL AUTO_INCREMENT COMMENT '权限id',
  `name` char(30) NOT NULL DEFAULT '' COMMENT '权限名',
  `parent_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上级',
  `type` enum('controller','menu','navigator') DEFAULT 'controller' COMMENT '权限类型：控制器、菜单、导航',
  `controller` char(100) NOT NULL DEFAULT '' COMMENT '控制器',
  `action` char(100) NOT NULL DEFAULT '' COMMENT '动作',
  `icon` char(100) NOT NULL DEFAULT '' COMMENT '图标（用于展示)',
  `is_display` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否显示：0不显示 1显示',
  `sequence` int(10) NOT NULL DEFAULT '0' COMMENT '排序(越小越靠前)',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `target` char(200) NOT NULL DEFAULT '' COMMENT '目标地址',
  PRIMARY KEY (`privilege_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='权限（动作）表';

INSERT INTO `system_privilege` (`privilege_id`, `name`, `parent_id`, `type`, `controller`, `action`, `icon`, `is_display`, `sequence`, `create_time`, `update_time`, `target`) VALUES
(1,	'系统',	0,	'navigator',	'',	'',	'fa fa-windows',	1,	2,	0,	1502181677,	'x'),
(3,	'权限管理',	1,	'menu',	'',	'',	'fa fa-lock',	1,	3,	0,	1502092329,	''),
(4,	'角色管理',	1,	'menu',	'',	'',	'fa fa-users',	1,	2,	0,	1502092279,	''),
(43,	'新增权限',	3,	'controller',	'privilege',	'add',	'glyphicon glyphicon-plus',	1,	1,	0,	1502084989,	''),
(45,	'权限列表',	3,	'controller',	'privilege',	'list',	'glyphicon glyphicon-align-justify',	1,	2,	0,	1502084884,	''),
(53,	'新增角色',	4,	'controller',	'role',	'add',	'glyphicon glyphicon-plus',	1,	1,	0,	1502085060,	''),
(61,	'角色列表',	4,	'controller',	'role',	'list',	'glyphicon glyphicon-align-justify',	1,	1,	1492598986,	1501643127,	''),
(73,	'用户管理',	1,	'menu',	'',	'',	'glyphicon glyphicon-user',	1,	1,	1492672432,	1492672456,	''),
(75,	'新增用户',	73,	'controller',	'user',	'add',	'glyphicon glyphicon-plus',	1,	1,	1492672685,	1502085169,	''),
(81,	'用户列表',	73,	'controller',	'user',	'list',	'glyphicon glyphicon-align-justify',	1,	1,	1492758209,	1502085186,	''),
(115,	'编辑用户',	73,	'controller',	'user',	'edit',	'fa fa-pencil-square',	0,	3,	1493361157,	1502085218,	''),
(138,	'编辑权限',	3,	'controller',	'privilege',	'edit',	'fa fa-pencil-square',	0,	2,	1502084965,	1502085150,	''),
(139,	'编辑角色',	4,	'controller',	'role',	'edit',	'fa fa-pencil-square',	0,	1,	1502085125,	0,	''),
(140,	'禁用用户',	73,	'controller',	'user',	'Forbidden',	'fa fa-times',	0,	4,	1502174927,	0,	''),
(141,	'恢复用户',	73,	'controller',	'user',	'Review',	'fa fa-check',	0,	5,	1502174959,	0,	''),
(142,	'删除角色',	4,	'controller',	'role',	'Delete',	'fa fa-times',	0,	4,	1502175039,	1502175050,	''),
(143,	'角色授权',	4,	'controller',	'role_privilege',	'add',	'fa fa-repeat',	0,	6,	1502175143,	0,	''),
(144,	'删除权限',	3,	'controller',	'privilege',	'delete',	'fa fa-times',	0,	4,	1502175191,	0,	''),
(145,	'我的',	0,	'navigator',	'',	'',	'fa fa-heart',	1,	1,	1502181668,	0,	''),
(146,	'个人中心',	145,	'menu',	'',	'',	'fa fa-leaf',	1,	1,	1502181713,	0,	''),
(147,	'我的资料',	146,	'controller',	'user',	'profile',	'fa fa-certificate',	1,	1,	1502181803,	0,	''),
(148,	'修改密码',	146,	'controller',	'user',	'changepassword',	'fa fa-key',	1,	2,	1502181896,	0,	''),
(149,	'用户',	0,	'navigator',	'',	'',	'fa fa-user',	1,	3,	1503562112,	0,	''),
(150,	'用户管理',	149,	'menu',	'',	'',	'fa fa-users',	1,	1,	1503562205,	0,	''),
(151,	'隧道云',	0,	'navigator',	'',	'',	'fa fa-cloud',	1,	4,	1503562260,	0,	''),
(152,	'角色管理',	149,	'menu',	'',	'',	'fa fa-sitemap',	1,	2,	1503562338,	0,	''),
(153,	'区域管理',	151,	'menu',	'',	'',	'fa fa-th',	1,	5,	1503562479,	0,	''),
(154,	'Server管理',	151,	'menu',	'',	'',	'fa fa-briefcase',	1,	2,	1503562633,	1503658805,	''),
(155,	'Client管理',	151,	'menu',	'',	'',	'fa fa-road',	1,	2,	1503562652,	0,	''),
(156,	'Tunnel管理',	151,	'menu',	'',	'',	'fa fa-university',	1,	3,	1503562699,	0,	''),
(157,	'数据统计',	151,	'menu',	'',	'',	'fa fa-coffee',	1,	4,	1503562839,	0,	''),
(158,	'角色列表',	152,	'controller',	'web/role',	'list',	'fa fa-th-list',	1,	1,	1503566704,	1503639994,	''),
(159,	'添加角色',	152,	'controller',	'web/role',	'add',	'fa fa-thumbs-o-up',	1,	2,	1503640058,	1503640101,	''),
(160,	'用户列表',	150,	'controller',	'web/user',	'list',	'fa fa-th-list',	1,	1,	1503642077,	0,	''),
(162,	'Server列表',	154,	'controller',	'web/server',	'list',	'fa fa-th-large',	1,	1,	1503644690,	0,	''),
(163,	'新增Server',	154,	'controller',	'web/server',	'add',	'fa fa-inbox',	1,	2,	1503644732,	0,	''),
(164,	'Client列表',	155,	'controller',	'web/client',	'list',	'fa fa-align-justify',	1,	1,	1503649600,	0,	''),
(166,	'区域列表',	153,	'controller',	'web/region',	'list',	'fa fa-list-ul',	1,	1,	1503654799,	0,	''),
(167,	'新增区域',	153,	'controller',	'web/region',	'add',	'fa fa-bell-o',	1,	2,	1503654830,	0,	''),
(168,	'Cluster管理',	151,	'menu',	'',	'',	'fa fa-heart',	1,	1,	1503658763,	0,	''),
(169,	'Cluster列表',	168,	'controller',	'web/cluster',	'list',	'fa fa-list-ul',	1,	1,	1503658868,	0,	''),
(170,	'新增Cluster',	168,	'controller',	'web/cluster',	'add',	'fa fa-magnet',	1,	2,	1503658928,	0,	''),
(171,	'Tunnel列表',	156,	'controller',	'web/tunnel',	'list',	'fa fa-random',	1,	1,	1503711836,	0,	''),
(172,	'在线Server',	157,	'controller',	'web/online',	'list?cs=server',	'fa fa-cc-diners-club',	1,	1,	1503712418,	1503739694,	''),
(173,	'在线Client',	157,	'controller',	'web/online',	'list?cs=client',	'fa fa-first-order',	1,	2,	1503712458,	1503739711,	''),
(174,	'用户流量',	157,	'controller',	'web/traffic',	'list',	'fa fa-user-md',	1,	3,	1503738742,	0,	''),
(175,	'修改用户',	150,	'controller',	'/web/user',	'edit',	'fa fa-cog',	0,	1,	1503888185,	0,	''),
(176,	'禁用用户',	150,	'controller',	'web/user',	'forbidden',	'fa fa-power-off',	0,	1,	1503888265,	0,	''),
(177,	'恢复用户',	150,	'controller',	'web/user',	'review',	'fa fa-power-off',	0,	1,	1503888309,	0,	''),
(178,	'修改角色',	152,	'controller',	'web/role',	'edit',	'fa fa-power-off',	0,	1,	1503888423,	0,	''),
(179,	'删除角色',	152,	'controller',	'web/role',	'delete',	'fa fa-power-off',	0,	1,	1503888465,	0,	''),
(180,	'编辑区域',	152,	'controller',	'web/role',	'regions',	'fa fa-power-off',	0,	1,	1503888516,	0,	''),
(181,	'修改Cluster',	168,	'controller',	'web/cluster',	'edit',	'fa fa-power-off',	0,	1,	1503889011,	0,	''),
(182,	'删除Clusetr',	168,	'controller',	'web/cluster',	'delete',	'fa fa-power-off',	0,	1,	1503889039,	1503889056,	''),
(183,	'禁用Cluster',	168,	'controller',	'web/cluster',	'forbidden',	'fa fa-power-off',	0,	1,	1503889104,	0,	''),
(184,	'恢复Cluster',	168,	'controller',	'web/cluster',	'review',	'fa fa-power-off',	0,	1,	1503889153,	0,	''),
(185,	'修改Server',	154,	'controller',	'web/server',	'edit',	'fa fa-power-off',	0,	1,	1503890008,	0,	''),
(186,	'删除Server',	154,	'controller',	'web/server',	'delete',	'fa fa-power-off',	0,	1,	1503890031,	0,	''),
(187,	'重置Token',	154,	'controller',	'web/server',	'reset',	'fa fa-power-off',	0,	1,	1503890081,	0,	'');

DROP TABLE IF EXISTS `system_role`;
CREATE TABLE `system_role` (
  `role_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名称',
  `is_delete` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除，默认0，否；1是',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色表';

INSERT INTO `system_role` (`role_id`, `name`, `is_delete`, `create_time`, `update_time`) VALUES
(3,	'普通用户',	0,	1502181940,	0);

DROP TABLE IF EXISTS `system_role_privilege`;
CREATE TABLE `system_role_privilege` (
  `role_privilege_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `role_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '角色ID',
  `privilege_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '权限ID',
  `is_delete` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除，默认0，否；1是',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`role_privilege_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色权限表';

INSERT INTO `system_role_privilege` (`role_privilege_id`, `role_id`, `privilege_id`, `is_delete`, `create_time`, `update_time`) VALUES
(48,	3,	145,	0,	0,	0),
(49,	3,	146,	0,	0,	0),
(50,	3,	147,	0,	0,	0),
(51,	3,	148,	0,	0,	0);

DROP TABLE IF EXISTS `system_user`;
CREATE TABLE `system_user` (
  `user_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `given_name` varchar(50) NOT NULL DEFAULT '' COMMENT '真实姓名',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `mobile` char(18) NOT NULL DEFAULT '' COMMENT '手机号码',
  `is_forbidden` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否屏蔽，默认0，否；1屏蔽',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

INSERT INTO `system_user` (`user_id`, `username`, `given_name`, `password`, `email`, `mobile`, `is_forbidden`, `create_time`, `update_time`) VALUES
(1,	'root',	'root',	'a906449d5769fa7361d7ecc6aa3f6d28',	'',	'',	0,	1501643981,	1503909597);

DROP TABLE IF EXISTS `system_user_role`;
CREATE TABLE `system_user_role` (
  `user_role_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `role_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '角色ID',
  `is_delete` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除，默认0，否；1是',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`user_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户角色表';


-- 2017-08-28 08:53:00