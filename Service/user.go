package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/util"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

const (
	RegisterCodeExpire   = 5 * time.Minute
	RegisterCodeCooldown = 60 * time.Second
)

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	ok, _ := regexp.MatchString(pattern, email)
	return ok
}
func SendRegisterCode(req models.SendRegisterCodeRequest) error {
	if req.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if !isValidEmail(req.Email) {
		return errors.New("邮箱格式不正确")
	}
	existUser, err := dao.GetByEmail(req.Email)
	if err != nil {
		fmt.Printf(err.Error())
		return errors.New("查询邮箱失败")
	}
	if existUser != nil {
		return errors.New("该邮箱已被注册")
	}
	ctx := context.Background() //初始化ctx
	codeKey := fmt.Sprintf("register_code:%s", req.Email)
	cooldownKey := fmt.Sprintf("register_cooldown:%s", req.Email)
	cooldownExists, err := util.RDB.Exists(ctx, cooldownKey).Result()
	if err != nil {
		return errors.New("查询失败")
	}
	if cooldownExists > 0 {
		return errors.New("操作过于频繁，请稍后再试")
	}
	code := util.GenerateCode()
	err = util.SendRegisterCodeEmail(req.Email, code)
	if err != nil {
		return errors.New("邮件发送失败：%v")
	}
	err = util.RDB.Set(ctx, codeKey, code, RegisterCodeExpire).Err()
	if err != nil {
		return errors.New("验证码保存失败")
	}
	err = util.RDB.Set(ctx, cooldownKey, 1, RegisterCodeCooldown).Err()
	if err != nil {
		return errors.New("发送冷却保存失败")
	}
	return nil
}

func Register(req models.RegisterRequest) error {
	if req.Username == "" {
		return errors.New("用户名不能为空")
	}
	if req.Password == "" {
		return errors.New("密码不能为空")
	}
	if req.ConfirmPassword == "" {
		return errors.New("确认密码不能为空")
	}
	if req.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if !isValidEmail(req.Email) {
		return errors.New("邮箱格式不正确")
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("两次密码不一致")
	}
	if len(req.Password) < 6 {
		return errors.New("密码不能小于6位")
	}
	ctx := context.Background()
	codeKey := fmt.Sprintf("register_code:%s", req.Email)
	realCode, err := util.RDB.Get(ctx, codeKey).Result()
	if err != nil {
		return errors.New("验证码已过期或者不存在")
	}
	if realCode != req.Code {
		return errors.New("验证码错误")
	}
	existUser, err := dao.GetByUsername(req.Username)
	if err != nil {
		return errors.New("查询用户名名失败")
	}
	if existUser != nil {
		return errors.New("用户名已存在")
	}
	emailUser, err := dao.GetByEmail(req.Email)
	if err != nil {
		return errors.New("邮箱查询失败")
	}
	if emailUser != nil {
		return errors.New("该邮箱已被注册")
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}
	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Email:    req.Email,
		Status:   1,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		CreateAt: time.Now(),
	}
	err = dao.CreateUser(user)
	if err != nil {
		println(err.Error())
		return errors.New("注册失败")
	}
	_ = util.RDB.Del(ctx, codeKey).Err()
	return nil
}

func SendEmailCode(c *gin.Context) error {
	var req models.SendRegisterCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return errors.New("参数错误")
	}
	err := SendRegisterCode(req)
	if err != nil {
		fmt.Printf(err.Error())

		return errors.New("发送失败")

	}
	return nil
}
