package util

import (
	"MYshop/config"
	"fmt"
	"gopkg.in/gomail.v2"
)

func SendEmail(toEmail, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.AppConfig.Email.From)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(
		config.AppConfig.Email.Host,
		config.AppConfig.Email.Port,
		config.AppConfig.Email.Username,
		config.AppConfig.Email.Password,
	)
	return d.DialAndSend(m)
}
func SendRegisterCodeEmail(toEmail, code string) error {
	return SendEmail(
		toEmail,
		"商城注册验证码",
		fmt.Sprintf(`
		<div>
			<h2>商城注册验证码</h2>
			<p>您的验证码是：<b style="color: blue; font-size: 24px;">%s</b></p>
			<p>验证码 5 分钟内有效，请勿泄露给他人。</p>
		</div>
	`, code),
	)
}

func SendLoginCodeEmail(toEmail, code string) error {
	return SendEmail(
		toEmail,
		"商城登录验证码",
		fmt.Sprintf(`
		<div>
			<h2>商城登录验证码</h2>
			<p>您的验证码是：<b style="color: green; font-size: 24px;">%s</b></p>
			<p>验证码 5 分钟内有效，请勿泄露给他人。</p>
		</div>
	`, code),
	)
}
func SendResetPasswordEmail(toEmail, code string) error {
	return SendEmail(
		toEmail,
		"商城找回密码验证码",
		fmt.Sprintf(`
		<div>
			<h2>商城找回密码验证码</h2>
			<p>您的验证码是：<b style="color: green; font-size: 24px;">%s</b></p>
			<p>验证码 5 分钟内有效，请勿泄露给他人。</p>
		</div>
	`, code),
	)
}
