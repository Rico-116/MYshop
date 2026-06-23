package dao

import (
	"MYshop/models"
	"MYshop/util"
	"errors"

	"gorm.io/gorm"
)

func CreateProductCommentTx(tx *gorm.DB, comment *models.ProductComment) error {
	return tx.Table("product_comment").Create(comment).Error
}

func GetProductCommentById(commentId uint) (*models.ProductComment, error) {
	var comment models.ProductComment
	err := util.Db.Table("product_comment").Where("id = ?", commentId).Limit(1).Scan(&comment).Error
	if err != nil {
		return nil, err
	}
	if comment.Id == 0 {
		return nil, nil
	}
	return &comment, nil
}

func HasProductRatingComment(userId uint, orderNo string, productId uint, skuId uint) (bool, error) {
	var count int64
	db := util.Db.Table("product_comment").
		Where("user_id = ? AND order_no = ? AND product_id = ? AND parent_id = 0 AND rating IS NOT NULL", userId, orderNo, productId)
	if skuId > 0 {
		db = db.Where("sku_id = ?", skuId)
	}
	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func CountProductTopComments(productId uint) (int64, error) {
	var total int64
	err := util.Db.Table("product_comment").
		Where("product_id = ? AND parent_id = 0", productId).
		Count(&total).Error
	return total, err
}

func GetProductTopComments(productId uint, page int, pageSize int) ([]models.ProductCommentVO, error) {
	var list []models.ProductCommentVO
	offset := (page - 1) * pageSize
	err := util.Db.Table("product_comment pc").
		Select(`pc.id,
			pc.product_id,
			pc.user_id,
			u.username,
			COALESCE(NULLIF(u.nickname, ''), u.username) AS nickname,
			u.avatar,
			pc.order_no,
			pc.sku_id,
			pc.rating,
			pc.content,
			pc.root_id,
			pc.parent_id,
			pc.reply_to_user_id,
			pc.created_at,
			DATE_FORMAT(pc.created_at, '%Y-%m-%d %H:%i:%s') AS created_time`).
		Joins("LEFT JOIN user u ON pc.user_id = u.id").
		Where("pc.product_id = ? AND pc.parent_id = 0", productId).
		Order("pc.created_at desc").
		Offset(offset).
		Limit(pageSize).
		Scan(&list).Error
	return list, err
}

func GetProductCommentDescendants(productId uint, rootIds []uint) ([]models.ProductCommentVO, error) {
	var list []models.ProductCommentVO
	if len(rootIds) == 0 {
		return list, nil
	}
	err := util.Db.Table("product_comment pc").
		Select(`pc.id,
			pc.product_id,
			pc.user_id,
			u.username,
			COALESCE(NULLIF(u.nickname, ''), u.username) AS nickname,
			u.avatar,
			pc.order_no,
			pc.sku_id,
			pc.rating,
			pc.content,
			pc.root_id,
			pc.parent_id,
			pc.reply_to_user_id,
			COALESCE(NULLIF(ru.nickname, ''), ru.username) AS reply_to_nickname,
			pc.created_at,
			DATE_FORMAT(pc.created_at, '%Y-%m-%d %H:%i:%s') AS created_time`).
		Joins("LEFT JOIN user u ON pc.user_id = u.id").
		Joins("LEFT JOIN user ru ON pc.reply_to_user_id = ru.id").
		Where("pc.product_id = ? AND pc.root_id IN ?", productId, rootIds).
		Order("pc.created_at asc").
		Scan(&list).Error
	return list, err
}

func UpdateProductRatingTx(tx *gorm.DB, productId uint) error {
	result := tx.Exec(`
		UPDATE product
		SET rating = COALESCE((SELECT ROUND(AVG(rating), 1) FROM product_comment WHERE product_id = ? AND parent_id = 0 AND rating IS NOT NULL), 0),
			rating_count = (SELECT COUNT(*) FROM product_comment WHERE product_id = ? AND parent_id = 0 AND rating IS NOT NULL),
			updated_at = NOW()
		WHERE id = ?
	`, productId, productId, productId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("商品不存在")
	}
	return nil
}
