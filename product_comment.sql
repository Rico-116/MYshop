CREATE TABLE IF NOT EXISTS `product_comment` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '评论用户ID',
  `order_no` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '评分绑定的订单号，回复为空',
  `sku_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评分绑定的SKU，回复为0',
  `rating` TINYINT DEFAULT NULL COMMENT '评分：1-5，回复为空',
  `content` VARCHAR(500) NOT NULL COMMENT '评论内容',
  `root_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '顶层评论ID，顶层评论为0',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID，顶层评论为0',
  `reply_to_user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '被回复用户ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_product_parent` (`product_id`, `parent_id`),
  KEY `idx_product_root` (`product_id`, `root_id`),
  KEY `idx_user_order_product` (`user_id`, `order_no`, `product_id`, `sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品评论与树形回复';
