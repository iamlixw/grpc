-- ----------------------------
-- 短信验证码
-- ----------------------------
CREATE TABLE `arc_sms` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `mobile` varchar(15) NOT NULL DEFAULT '0' COMMENT '手机号码',
  `code` varchar(15)  NOT NULL DEFAULT '' COMMENT '验证码',
  `createTime` int(10) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;