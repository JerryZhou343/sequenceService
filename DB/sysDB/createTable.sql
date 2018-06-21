DROP TABLE IF EXISTS `t_sequence_number`;
CREATE TABLE `t_sequence_number` (
	`seq_id` BIGINT auto_increment COMMENT '序列号标识',
	`first_id` INT NOT NULL COMMENT '一级业务识别码',
	`second_id` INT NOT NULL COMMENT '二级业务识别码',
	`base_value` BIGINT NOT NULL DEFAULT 1 COMMENT '基值',
	`max_value` BIGINT NOT NULL DEFAULT 9999999 COMMENT '最大值',
	`current_value` BIGINT NOT NULL DEFAULT 1 COMMENT '当前值',
	`step_length` TINYINT NOT NULL DEFAULT 1 COMMENT '步长值',
	`reset_type` TINYINT NOT NULL DEFAULT 4 COMMENT '重置类型标识：1：每日重置 2：每月重置 3：每年重置 4：手动重置',
	`last_reset_time` INT NOT NULL COMMENT '上一次重置时间',
	`first_name` VARCHAR ( 255 ) NOT NULL COMMENT '一级业务名称',
	`second_name` VARCHAR ( 255 ) NOT NULL COMMENT '二级业务名称',
	`remark` VARCHAR ( 255 ) COMMENT '备注',
	`add_time` INT NOT NULL COMMENT '添加时间',
	UNIQUE KEY(`seq_id`),
	PRIMARY KEY ( `first_id`, `second_id` )
) ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT '序列号表';
