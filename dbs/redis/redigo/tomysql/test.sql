CREATE TABLE `hashset` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增序号',
	`rkey` varchar(100) NOT NULL COMMENT '存储redis中hashset key',
	`rval` MediumText NOT NULL COMMENT '存储redis中hashset value',
	`create_time` timestamp NOT NULL DEFAULT now() COMMENT '创建时间',
  	`update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY (`rkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'hashset存储';

CREATE TABLE `list` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增序号',
	`rkey` varchar(100) NOT NULL COMMENT '存储redis中list key',
	`rval` MediumText NOT NULL COMMENT '存储redis中list value',
	`create_time` timestamp NOT NULL DEFAULT now() COMMENT '创建时间',
  	`update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY (`rkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'list存储';
