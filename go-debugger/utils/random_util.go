package utils

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func GetUUID() string {
	u1, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return u1.String()
}

// 通过时间搓 + 随机数生成的较短的随机code
func GetGenerateUniqueCode() string {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(1000) // 生成一个0到999之间的随机数

	uniqueNumber := strconv.FormatInt(timestamp, 10) + strconv.Itoa(randomNum)
	return uniqueNumber
}

// GetRandomNumber 生成验证码
func GetRandomNumber(number int) string {
	var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 限定随机数生成的范围。这里用了1e9来表示10的9次方
	max := int32(1e9)
	code := rnd.Int31n(max)
	str := fmt.Sprintf("%0*d", number, code) // 为保持结果为9位数，使用%09d保持格式
	return str[len(str)-number:]
}

// GetRandomPassword 生成随机密码
func GetRandomPassword(length int) string {
	baseStr := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	bytes := make([]byte, length, length)
	l := len(baseStr)
	for i := 0; i < length; i++ {
		bytes[i] = baseStr[r.Intn(l)]
	}
	return string(bytes)
}
