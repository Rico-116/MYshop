package controller

import (
	"MYshop/Service"
	"MYshop/models"
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AdminGetCategoryList(c *gin.Context) {
	var status *uint
	if statusStr := c.Query("status"); statusStr != "" {
		value, err := strconv.Atoi(statusStr)
		if err != nil || value < 0 {
			util.Fail(c, 400, "分类状态参数错误")
			return
		}
		statusValue := uint(value)
		status = &statusValue
	}
	level := 0
	if levelStr := c.Query("level"); levelStr != "" {
		value, err := strconv.Atoi(levelStr)
		if err != nil {
			util.Fail(c, 400, "分类层级参数错误")
			return
		}
		level = value
	}
	list, err := Service.GetAdminCategoryList(status, level)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "查询分类列表成功", gin.H{"list": list})
}

func AdminCreateCategory(c *gin.Context) {
	var req models.AdminCreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	category, err := Service.CreateAdminCategory(req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "新增分类成功", gin.H{"category": category})
}

func AdminUpdateCategory(c *gin.Context) {
	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil || categoryId <= 0 {
		util.Fail(c, 400, "分类id错误")
		return
	}
	var req models.AdminUpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, 400, "参数错误")
		return
	}
	category, err := Service.UpdateAdminCategory(uint(categoryId), req)
	if err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "修改分类成功", gin.H{"category": category})
}

func AdminDeleteCategory(c *gin.Context) {
	categoryId, err := strconv.Atoi(c.Param("id"))
	if err != nil || categoryId <= 0 {
		util.Fail(c, 400, "分类id错误")
		return
	}
	if err := Service.DeleteAdminCategory(uint(categoryId)); err != nil {
		util.Fail(c, 400, err.Error())
		return
	}
	util.Success(c, "删除分类成功", nil)
}
