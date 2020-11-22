drop table if exists spd_address;
drop table if exists spd_image;

CREATE TABLE `spd_address`
(
    `id`             INT UNSIGNED AUTO_INCREMENT COMMENT '自增ID',
    `url`            VARCHAR(511) NOT NULL DEFAULT '' COMMENT 'URL',
    `crawl_status`   TINYINT      NOT NULL DEFAULT 0 COMMENT '爬取状态',
    `crawl_source`   VARCHAR(511) NOT NULL DEFAULT '' COMMENT '爬取来源',
    `crawl_node_num` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '爬取节点号',
    `created_at`     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `url_idx` (url)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='爬虫爬取地址表';


CREATE TABLE `spd_image`
(
    `id`         INT UNSIGNED AUTO_INCREMENT COMMENT '自增ID',
    `url`        VARCHAR(511) NOT NULL DEFAULT '' COMMENT '图片数据来源URL',
    `md5`        VARCHAR(511) NOT NULL DEFAULT '' UNIQUE COMMENT '图片MD5',
    `address_id` INT          NOT NULL DEFAULT 0 COMMENT '来源地址ID',
    `created_at` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='爬虫爬取图片表';

