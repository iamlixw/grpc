-- ----------------------------
-- 用户表
-- ----------------------------
CREATE TABLE `arc_user` (
  `id` varchar(50) NOT NULL DEFAULT '' COMMENT '用户ID',
  `username` varchar(50)  NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码',
  `name` varchar(50)  NOT NULL DEFAULT '' COMMENT '姓名',
  `mobile` varchar(15) NOT NULL DEFAULT '0' COMMENT '手机号码',
  `email` varchar(255)  NOT NULL DEFAULT '' COMMENT '邮箱',
  `createTime` int(10) NOT NULL COMMENT '创建时间',
  `updateTime` int(10) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;