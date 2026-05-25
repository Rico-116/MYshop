package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

func ValidateProductCategory(categoryId int) error {
	category, err := dao.GetCategoryById(categoryId)
	if err != nil {
		return err
	}
	if category == nil || category.Id == 0 {
		return errors.New("分类不存在")
	}
	if category.ParentId == 0 {
		return errors.New("新增商品只能选择二级分类")
	}
	return nil
}

func CreateAdminProduct(req models.AdminCreateProductRequest) (*models.Product, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Subtitle = strings.TrimSpace(req.Subtitle)
	req.MainImage = strings.TrimSpace(req.MainImage)
	req.Description = strings.TrimSpace(req.Description)
	if req.CategoryId == 0 {
		return nil, errors.New("请选择商品分类")
	}
	if req.Name == "" {
		return nil, errors.New("商品名称不能为空")
	}
	if err := ValidateProductCategory(int(req.CategoryId)); err != nil {
		return nil, err
	}
	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	if status != 0 && status != 1 {
		return nil, errors.New("商品状态只能是0或者1")
	}
	now := time.Now()
	product := &models.Product{
		CategoryId:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
		Subtitle:    req.Subtitle,
		MainImage:   req.MainImage,
		Status:      status,
		Price:       0,
		Rating:      0,
		RatingCount: 0,
		ClickCount:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := dao.AdminCreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

func GetAdminProductList(keyword string, categoryId *uint, status *int, page, pageSize int) ([]models.Product, int64, error) {
	keyword = strings.TrimSpace(keyword)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	if status != nil && *status != 0 && *status != 1 {
		return nil, 0, errors.New("商品状态只能是0或者1")
	}
	return dao.AdminGetProductList(keyword, categoryId, status, page, pageSize)
}

func GetAdminProductDetail(productId uint) (*models.Product, []models.ProductSku, error) {
	if productId == 0 {
		return nil, nil, errors.New("商品id不能为空")
	}
	product, err := dao.AdminGetProductById(int(productId))
	if err != nil {
		return nil, nil, err
	}
	if product == nil || product.Id == 0 {
		return nil, nil, errors.New("商品不存在")
	}
	skuList, err := dao.AdminGetSkuListByProductId(product.Id)
	if err != nil {
		return nil, nil, err
	}
	return product, skuList, nil
}

func UpdateAdminProduct(productId uint, req models.AdminUpdateProductRequest) (*models.Product, error) {
	if productId == 0 {
		return nil, errors.New("商品id不能为空")
	}
	oldProduct, err := dao.AdminGetProductById(int(productId))
	if err != nil {
		return nil, err
	}
	if oldProduct == nil || oldProduct.Id == 0 {
		return nil, errors.New("商品不存在")
	}
	updates := make(map[string]interface{})
	if req.CategoryId != nil {
		if *req.CategoryId == 0 {
			return nil, errors.New("请选择商品分类")
		}
		if err := ValidateProductCategory(int(*req.CategoryId)); err != nil {
			return nil, err
		}
		updates["category_id"] = *req.CategoryId
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("商品名称不能为空")
		}
		updates["name"] = name
	}
	if req.Subtitle != nil {
		updates["subtitle"] = strings.TrimSpace(*req.Subtitle)
	}
	if req.MainImage != nil {
		updates["main_image"] = strings.TrimSpace(*req.MainImage)
	}
	if req.Description != nil {
		updates["description"] = strings.TrimSpace(*req.Description)
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return nil, errors.New("商品状态只能是0或者1")
		}
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return nil, errors.New("没有需要修改的商品信息")
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	if err := dao.AdminUpdateProduct(productId, updates); err != nil {
		return nil, err
	}
	product, err := dao.AdminGetProductById(int(productId))
	if err != nil {
		return nil, err
	}
	categoryId := oldProduct.CategoryId
	if product != nil {
		categoryId = product.CategoryId
	}
	DeleteProductCache(int(productId), int(categoryId))
	return product, nil
}

func DeleteAdminProduct(productId uint) error {
	if productId == 0 {
		return errors.New("商品id不能为空")
	}
	product, err := dao.AdminGetProductById(int(productId))
	if err != nil {
		return err
	}
	if product == nil || product.Id == 0 {
		return errors.New("商品不存在")
	}
	if err := dao.AdminDeleteProduct(productId); err != nil {
		return err
	}
	DeleteProductCache(int(productId), int(product.CategoryId))
	return nil
}

func CreateAdminProductSku(req models.AdminCreateSkuRequest) (*models.ProductSku, error) {
	req.SkuCode = strings.TrimSpace(req.SkuCode)
	req.SkuName = strings.TrimSpace(req.SkuName)
	req.Image = strings.TrimSpace(req.Image)
	if req.ProductId == 0 {
		return nil, errors.New("商品id不能为空")
	}
	if req.SkuCode == "" {
		return nil, errors.New("SKU编码不能为空")
	}
	if req.SkuName == "" {
		return nil, errors.New("SKU名称不能为空")
	}
	if req.Price <= 0 {
		return nil, errors.New("SKU价格必须大于0")
	}
	if req.Stock < 0 {
		return nil, errors.New("SKU库存不能小于0")
	}
	product, err := dao.AdminGetProductById(int(req.ProductId))
	if err != nil {
		return nil, err
	}
	if product == nil || product.Id == 0 {
		return nil, errors.New("商品不存在")
	}
	count, err := dao.AdminCountSkuCode(req.ProductId, req.SkuCode)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该商品下SKU编码已存在")
	}
	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	if status != 0 && status != 1 {
		return nil, errors.New("SKU状态只能是0或1")
	}
	now := time.Now()
	sku := &models.ProductSku{
		ProductId: req.ProductId,
		SkuCode:   req.SkuCode,
		SkuName:   req.SkuName,
		Price:     req.Price,
		Stock:     req.Stock,
		Image:     req.Image,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := dao.AdminCreateProductSku(sku); err != nil {
		return nil, err
	}
	if err := dao.SyncProductMinPrice(int(req.ProductId)); err != nil {
		return nil, err
	}
	DeleteProductCache(int(req.ProductId), int(product.CategoryId))
	return sku, nil
}

func UpdateAdminProductSku(skuId uint, req models.AdminUpdateSkuRequest) (*models.ProductSku, error) {
	if skuId == 0 {
		return nil, errors.New("SKU id不能为空")
	}
	oldSku, err := dao.AdminGetSkuById(skuId)
	if err != nil {
		return nil, err
	}
	if oldSku == nil || oldSku.Id == 0 {
		return nil, errors.New("SKU不存在")
	}
	updates := make(map[string]interface{})
	if req.SkuCode != nil {
		skuCode := strings.TrimSpace(*req.SkuCode)
		if skuCode == "" {
			return nil, errors.New("SKU编码不能为空")
		}
		count, err := dao.AdminCountSkuCodeExcludeId(oldSku.ProductId, skuCode, skuId)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("该商品下SKU编码已存在")
		}
		updates["sku_code"] = skuCode
	}
	if req.SkuName != nil {
		skuName := strings.TrimSpace(*req.SkuName)
		if skuName == "" {
			return nil, errors.New("SKU名称不能为空")
		}
		updates["sku_name"] = skuName
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, errors.New("SKU价格必须大于0")
		}
		updates["price"] = *req.Price
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("SKU库存不能小于0")
		}
		updates["stock"] = *req.Stock
	}
	if req.Image != nil {
		updates["image"] = strings.TrimSpace(*req.Image)
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return nil, errors.New("SKU状态只能是0或1")
		}
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return nil, errors.New("没有需要修改的SKU信息")
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	if err := dao.AdminUpdateSku(skuId, updates); err != nil {
		return nil, err
	}
	if err := dao.SyncProductMinPrice(int(oldSku.ProductId)); err != nil {
		return nil, err
	}
	product, _ := dao.AdminGetProductById(int(oldSku.ProductId))
	if product != nil {
		DeleteProductCache(int(oldSku.ProductId), int(product.CategoryId))
	}
	return dao.AdminGetSkuById(skuId)
}

func DeleteAdminProductSku(skuId uint) error {
	if skuId == 0 {
		return errors.New("SKU id不能为空")
	}
	sku, err := dao.AdminGetSkuById(skuId)
	if err != nil {
		return err
	}
	if sku == nil || sku.Id == 0 {
		return errors.New("SKU不存在")
	}
	if err := dao.AdminDeleteSku(skuId); err != nil {
		return err
	}
	if err := dao.SyncProductMinPrice(int(sku.ProductId)); err != nil {
		return err
	}
	product, _ := dao.AdminGetProductById(int(sku.ProductId))
	if product != nil {
		DeleteProductCache(int(sku.ProductId), int(product.CategoryId))
	}
	return nil
}
