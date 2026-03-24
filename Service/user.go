package Service

import (
	"MYshop/dao"
	"MYshop/models"
	"MYshop/package/logger"
	"MYshop/util"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		logger.Log.Warn("发送注册验证码失败：邮箱不能为空")
		return errors.New("邮箱不能为空")
	}

	if !isValidEmail(req.Email) {
		logger.Log.Warn("发送注册验证码失败：邮箱格式不正确",
			zap.String("email", req.Email),
		)
		return errors.New("邮箱格式不正确")
	}

	existUser, err := dao.GetByEmail(req.Email)
	if err != nil {
		logger.Log.Error("查询邮箱失败",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return errors.New("查询邮箱失败")
	}

	if existUser != nil {
		logger.Log.Warn("发送注册验证码失败：该邮箱已被注册",
			zap.String("email", req.Email),
		)
		return errors.New("该邮箱已被注册")
	}

	ctx := context.Background()
	codeKey := fmt.Sprintf("register_code:%s", req.Email)
	cooldownKey := fmt.Sprintf("register_cooldown:%s", req.Email)

	cooldownExists, err := util.RDB.Exists(ctx, cooldownKey).Result()
	if err != nil {
		logger.Log.Error("查询发送冷却失败",
			zap.String("email", req.Email),
			zap.String("cooldown_key", cooldownKey),
			zap.Error(err),
		)
		return errors.New("查询失败")
	}

	if cooldownExists > 0 {
		logger.Log.Warn("发送注册验证码失败：操作过于频繁",
			zap.String("email", req.Email),
		)
		return errors.New("操作过于频繁，请稍后再试")
	}

	code := util.GenerateCode()

	err = util.SendRegisterCodeEmail(req.Email, code)
	if err != nil {
		logger.Log.Error("发送注册验证码邮件失败",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return errors.New("邮件发送失败")
	}

	err = util.RDB.Set(ctx, codeKey, code, RegisterCodeExpire).Err()
	if err != nil {
		logger.Log.Error("保存注册验证码失败",
			zap.String("email", req.Email),
			zap.String("code_key", codeKey),
			zap.Error(err),
		)
		return errors.New("验证码保存失败")
	}

	err = util.RDB.Set(ctx, cooldownKey, 1, RegisterCodeCooldown).Err()
	if err != nil {
		logger.Log.Error("保存发送冷却失败",
			zap.String("email", req.Email),
			zap.String("cooldown_key", cooldownKey),
			zap.Error(err),
		)
		return errors.New("发送冷却保存失败")
	}

	logger.Log.Info("注册验证码发送成功",
		zap.String("email", req.Email),
	)

	return nil
}
func SendLoginCode(req models.SendLoginCodeRequest) error {
	if req.Email == "" {
		logger.Log.Warn("发送登录验证码失败：邮箱不能为空")
		return errors.New("邮箱不能为空")
	}

	if !isValidEmail(req.Email) {
		logger.Log.Warn("发送登录验证码失败：邮箱格式不正确",
			zap.String("email", req.Email),
		)
		return errors.New("邮箱格式不正确")
	}

	user, err := dao.GetByEmail(req.Email)
	if err != nil {
		logger.Log.Error("发送登录验证码失败：查询邮箱失败",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return errors.New("查询邮箱失败")
	}

	if user == nil {
		logger.Log.Warn("发送登录验证码失败：该邮箱未注册",
			zap.String("email", req.Email),
		)
		return errors.New("该邮箱未注册")
	}

	ctx := context.Background()
	codeKey := fmt.Sprintf("login_code:%s", req.Email)
	cooldownKey := fmt.Sprintf("login_cooldown:%s", req.Email)

	cooldownExists, err := util.RDB.Exists(ctx, cooldownKey).Result()
	if err != nil {
		logger.Log.Error("发送登录验证码失败：查询发送冷却失败",
			zap.String("email", req.Email),
			zap.String("cooldown_key", cooldownKey),
			zap.Error(err),
		)
		return errors.New("查询失败")
	}

	if cooldownExists > 0 {
		logger.Log.Warn("发送登录验证码失败：操作过于频繁",
			zap.String("email", req.Email),
		)
		return errors.New("操作过于频繁，请稍后再试")
	}

	code := util.GenerateCode()

	err = util.SendLoginCodeEmail(req.Email, code)
	if err != nil {
		logger.Log.Error("发送登录验证码失败：邮件发送失败",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return errors.New("邮件发送失败")
	}

	err = util.RDB.Set(ctx, codeKey, code, 5*time.Minute).Err()
	if err != nil {
		logger.Log.Error("发送登录验证码失败：验证码保存失败",
			zap.String("email", req.Email),
			zap.String("code_key", codeKey),
			zap.Error(err),
		)
		return errors.New("验证码保存失败")
	}

	err = util.RDB.Set(ctx, cooldownKey, 1, 60*time.Second).Err()
	if err != nil {
		logger.Log.Error("发送登录验证码失败：发送冷却保存失败",
			zap.String("email", req.Email),
			zap.String("cooldown_key", cooldownKey),
			zap.Error(err),
		)
		return errors.New("发送冷却保存失败")
	}

	logger.Log.Info("登录验证码发送成功",
		zap.String("email", req.Email),
	)

	return nil
}
func Register(req models.RegisterRequest) error {
	if req.Username == "" {
		logger.Log.Warn("注册失败：用户名不能为空")
		return errors.New("用户名不能为空")
	}
	if req.Password == "" {
		logger.Log.Warn("注册失败：密码不能为空")
		return errors.New("密码不能为空")
	}
	if req.ConfirmPassword == "" {
		logger.Log.Warn("注册失败：确认密码不能为空")
		return errors.New("确认密码不能为空")
	}
	if req.Email == "" {
		logger.Log.Warn("注册失败：邮箱不能为空")
		return errors.New("邮箱不能为空")
	}
	if !isValidEmail(req.Email) {
		logger.Log.Warn("注册失败：邮箱格式不正确",
			zap.String("email", req.Email))
		return errors.New("邮箱格式不正确")
	}
	if req.Password != req.ConfirmPassword {
		logger.Log.Warn("注册失败：两次密码不一样")
		return errors.New("两次密码不一致")
	}
	if len(req.Password) < 6 {
		logger.Log.Warn("注册失败：密码不能小于6位")
		return errors.New("密码不能小于6位")
	}
	ctx := context.Background()
	codeKey := fmt.Sprintf("register_code:%s", req.Email)
	realCode, err := util.RDB.Get(ctx, codeKey).Result()
	if err != nil {
		logger.Log.Error("注册失败：key不存在",
			zap.Error(err))
		return errors.New("验证码已过期或者不存在")
	}
	if realCode != req.Code {
		logger.Log.Warn("注册失败：验证码错误")
		return errors.New("验证码错误")
	}
	existUser, err := dao.GetByUsername(req.Username)
	if err != nil {
		logger.Log.Error("注册失败:查询用户名失败",
			zap.String("username", req.Username))
		return errors.New("查询用户名名失败")
	}
	if existUser != nil {
		logger.Log.Warn("注册失败：用户名已存在",
			zap.String("username", req.Username))
		return errors.New("用户名已存在")
	}
	emailUser, err := dao.GetByEmail(req.Email)
	if err != nil {
		logger.Log.Error("注册失败：邮箱查询失败",
			zap.String("email", req.Email),
			zap.Error(err))
		return errors.New("邮箱查询失败")
	}
	if emailUser != nil {
		logger.Log.Warn("注册失败：该邮箱已被注册",
			zap.String("email", req.Email))
		return errors.New("该邮箱已被注册")
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error("注册失败：密码加密失败",
			zap.Error(err))
		return errors.New("密码加密失败")
	}
	user := &models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Status:    1,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		CreatedAt: time.Now(),
	}
	err = dao.CreateUser(user)
	if err != nil {
		//println(err.Error())
		logger.Log.Warn("注册失败",
			zap.Error(err))
		return errors.New("注册失败")
	}
	logger.Log.Info("用户注册成功",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
		zap.Uint("user_id", user.UserId),
	)

	_ = util.RDB.Del(ctx, codeKey).Err()
	return nil
}

func SendEmailCode(c *gin.Context) error {
	var req models.SendRegisterCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("发送邮箱验证码参数错误",
			zap.Error(err),
			zap.String("ip", c.ClientIP()),
		)
		return errors.New("参数错误")
	}
	err := SendRegisterCode(req)
	if err != nil {
		logger.Log.Warn("发送注册验证码参数错误",
			zap.String("email", req.Email),
			zap.Error(err),
			zap.String("ip", c.ClientIP()))
		return errors.New("发送失败")

	}
	logger.Log.Info("验证码发送成功",
		zap.String("email", req.Email),
		zap.String("ip", c.ClientIP()))
	return nil
}

func EmailLogin(req models.EmailLoginRequest) (string, *models.User, error) {
	if req.Email == "" {
		logger.Log.Warn("邮箱验证码登录失败：邮箱不能为空")
		return "", nil, errors.New("邮箱不能为空")
	}

	if req.Code == "" {
		logger.Log.Warn("邮箱验证码登录失败：验证码不能为空",
			zap.String("email", req.Email),
		)
		return "", nil, errors.New("验证码不能为空")
	}

	if !isValidEmail(req.Email) {
		logger.Log.Warn("邮箱验证码登录失败：邮箱格式不正确",
			zap.String("email", req.Email),
		)
		return "", nil, errors.New("邮箱格式不正确")
	}

	ctx := context.Background()
	codeKey := fmt.Sprintf("login_code:%s", req.Email)

	realCode, err := util.RDB.Get(ctx, codeKey).Result()
	if err != nil {
		logger.Log.Error("邮箱验证码登录失败：查询验证码失败",
			zap.String("email", req.Email),
			zap.String("code_key", codeKey),
			zap.Error(err),
		)
		return "", nil, errors.New("验证码过期或者不存在")
	}

	if realCode != req.Code {
		logger.Log.Warn("邮箱验证码登录失败：验证码错误",
			zap.String("email", req.Email),
		)
		return "", nil, errors.New("验证码错误")
	}

	user, err := dao.GetByEmail(req.Email)
	if err != nil {
		logger.Log.Error("邮箱验证码登录失败：查询用户失败",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return "", nil, errors.New("查询用户失败")
	}

	if user == nil {
		logger.Log.Warn("邮箱验证码登录失败：用户不存在",
			zap.String("email", req.Email),
		)
		return "", nil, errors.New("用户不存在")
	}

	token, err := util.GenerateToken(user.UserId, user.Username)
	if err != nil {
		logger.Log.Error("邮箱验证码登录失败：生成 token 失败",
			zap.String("email", req.Email),
			zap.Uint("user_id", user.UserId),
			zap.String("username", user.Username),
			zap.Error(err),
		)
		return "", nil, errors.New("生成token失败")
	}

	if err = util.RDB.Del(ctx, codeKey).Err(); err != nil {
		logger.Log.Warn("邮箱验证码登录成功，但删除验证码失败",
			zap.String("email", req.Email),
			zap.String("code_key", codeKey),
			zap.Error(err),
		)
	}

	logger.Log.Info("邮箱验证码登录成功",
		zap.String("email", req.Email),
		zap.Uint("user_id", user.UserId),
		zap.String("username", user.Username),
	)

	return token, user, nil
}
