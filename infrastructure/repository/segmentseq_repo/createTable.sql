DROP TABLE IF EXISTS `t_segment_sequence`;
CREATE TABLE `t_segment_sequence` (
	`id` BIGINT auto_increment COMMENT '序列号标识',
	`product_id` INT NOT NULL COMMENT '一级业务识别码',
	`business_id` INT NOT NULL COMMENT '二级业务识别码',
	`base_value` BIGINT NOT NULL DEFAULT 1 COMMENT '基值',
	`max_value` BIGINT NOT NULL DEFAULT 9999999 COMMENT '最大值',
	`current_value` BIGINT NOT NULL DEFAULT 1 COMMENT '当前值',
	`step_length` TINYINT NOT NULL DEFAULT 1 COMMENT '步长值',
	`reset_type` TINYINT NOT NULL DEFAULT 4 COMMENT '重置类型标识：1：每日重置 2：每月重置 3：每年重置 4：手动重置',
	`last_reset_time` INT NOT NULL COMMENT '上一次重置时间',
	`product_name` VARCHAR ( 255 ) NOT NULL COMMENT '一级业务名称',
	`business_name` VARCHAR ( 255 ) NOT NULL COMMENT '二级业务名称',
	`remark` VARCHAR ( 255 ) COMMENT '备注',
	`add_time` INT NOT NULL COMMENT '添加时间',
	UNIQUE KEY(`id`),
	INDEX `idx_pid_bid`(`product_id`,`business_id`)
) ENGINE = INNODB DEFAULT CHARSET = utf8 COMMENT '序列号表';
