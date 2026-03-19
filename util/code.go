package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateCode() string { //生成验证码
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%0d", rand.Intn(1000000))
}
