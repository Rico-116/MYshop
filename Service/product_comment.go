package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/util"
	"errors"
	"strings"
)

func CreateProductComment(userId uint, req models.CreateProductCommentRequest) (*models.ProductComment, error) {
	if userId == 0 {
		return nil, errors.New("用户未登录")
	}
	req.Content = strings.TrimSpace(req.Content)
	req.OrderNo = strings.TrimSpace(req.OrderNo)
	if req.ProductId == 0 {
		return nil, errors.New("商品不能为空")
	}
	if req.Content == "" {
		return nil, errors.New("评论内容不能为空")
	}
	if len([]rune(req.Content)) > 500 {
		return nil, errors.New("评论内容不能超过500个字符")
	}

	if req.ParentId > 0 {
		return createProductCommentReply(userId, req)
	}
	return createProductRatingComment(userId, req)
}

func createProductRatingComment(userId uint, req models.CreateProductCommentRequest) (*models.ProductComment, error) {
	if req.OrderNo == "" {
		return nil, errors.New("订单号不能为空")
	}
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("评分只能是1到5分")
	}

	order, err := dao.GetOrderByOrderNoAndUserId(userId, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if order == nil || order.Id == 0 {
		return nil, errors.New("订单不存在")
	}
	if order.Status != models.OrderStatusFinished {
		return nil, errors.New("确认收货后才可以评分")
	}

	product, err := dao.GetProductById(int(req.ProductId))
	if err != nil {
		return nil, err
	}
	if product == nil || product.Id == 0 {
		return nil, errors.New("商品不存在")
	}

	orderItem, err := dao.GetOrderItemByOrderProduct(userId, req.OrderNo, req.ProductId, req.SkuId)
	if err != nil {
		return nil, err
	}
	if orderItem == nil || orderItem.Id == 0 {
		return nil, errors.New("该订单中没有这个商品")
	}
	if req.SkuId == 0 {
		req.SkuId = orderItem.SkuId
	}

	exists, err := dao.HasProductRatingComment(userId, req.OrderNo, req.ProductId, req.SkuId)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("该订单商品已经评价过")
	}

	rating := req.Rating
	comment := &models.ProductComment{
		ProductId: req.ProductId,
		UserId:    userId,
		OrderNo:   req.OrderNo,
		SkuId:     req.SkuId,
		Rating:    &rating,
		Content:   req.Content,
		RootId:    0,
		ParentId:  0,
	}

	tx := util.Db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	if err := dao.CreateProductCommentTx(tx, comment); err != nil {
		return nil, err
	}
	if err := dao.UpdateProductRatingTx(tx, req.ProductId); err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	committed = true

	DeleteProductCache(int(req.ProductId), int(product.CategoryId))
	return comment, nil
}

func createProductCommentReply(userId uint, req models.CreateProductCommentRequest) (*models.ProductComment, error) {
	parent, err := dao.GetProductCommentById(req.ParentId)
	if err != nil {
		return nil, err
	}
	if parent == nil || parent.Id == 0 {
		return nil, errors.New("要回复的评论不存在")
	}
	if parent.ProductId != req.ProductId {
		return nil, errors.New("回复的评论不属于该商品")
	}

	rootId := parent.RootId
	if rootId == 0 {
		rootId = parent.Id
	}
	comment := &models.ProductComment{
		ProductId:     parent.ProductId,
		UserId:        userId,
		Content:       req.Content,
		RootId:        rootId,
		ParentId:      parent.Id,
		ReplyToUserId: parent.UserId,
	}
	if err := util.Db.Table("product_comment").Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func GetProductCommentTree(productId uint, page int, pageSize int) (*models.ProductCommentListResult, error) {
	if productId == 0 {
		return nil, errors.New("商品不能为空")
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	total, err := dao.CountProductTopComments(productId)
	if err != nil {
		return nil, err
	}
	tops, err := dao.GetProductTopComments(productId, page, pageSize)
	if err != nil {
		return nil, err
	}
	rootIds := make([]uint, 0, len(tops))
	for _, item := range tops {
		rootIds = append(rootIds, item.Id)
	}
	descendants, err := dao.GetProductCommentDescendants(productId, rootIds)
	if err != nil {
		return nil, err
	}

	childrenByParent := make(map[uint][]models.ProductCommentVO)
	for _, item := range descendants {
		childrenByParent[item.ParentId] = append(childrenByParent[item.ParentId], item)
	}
	for i := range tops {
		fillCommentChildren(&tops[i], childrenByParent)
	}

	return &models.ProductCommentListResult{
		List:     tops,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func fillCommentChildren(node *models.ProductCommentVO, childrenByParent map[uint][]models.ProductCommentVO) {
	children := childrenByParent[node.Id]
	for i := range children {
		fillCommentChildren(&children[i], childrenByParent)
	}
	node.Children = children
}
