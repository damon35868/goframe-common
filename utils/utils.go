package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/damon35868/goframe-common/types"
	"github.com/damon35868/goframe-common/vars"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func GetTokenUserId(ctx context.Context, jwtKey string) uint {
	tokenString := g.RequestFromCtx(ctx).Request.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	tokenClaims, _ := jwt.ParseWithClaims(tokenString, &types.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if tokenClaims == nil {
		return 0
	}
	if claims, ok := tokenClaims.Claims.(*types.JwtClaims); ok && tokenClaims.Valid {
		return claims.Id
	}

	return 0
}

func GenToken(ctx context.Context, jwtKey string, id uint, registeredClaims jwt.RegisteredClaims) (token string, err error) {
	uc := &types.JwtClaims{
		Id:               id,
		RegisteredClaims: registeredClaims,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	token, err = t.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return
}

func GetClientTokenUserId(ctx context.Context) uint {
	return GetTokenUserId(ctx, g.Cfg().MustGet(ctx, vars.ClientJwtKey).String())
}

func GetAdminTokenUserId(ctx context.Context) uint {
	return GetTokenUserId(ctx, g.Cfg().MustGet(ctx, vars.AdminJwtKey).String())
}

func GenerateOrderNo(prefixs ...string) string {
	prefix := ""
	if len(prefixs) > 0 {
		prefix = prefixs[0]
	}
	// 格式化时间为 YYYYMMDDHHMMSS
	dateStr := time.Now().Format("20060102150405")
	// 生成两个 4 位随机数并拼接
	// rand.Intn(9000) + 1000 确保一定是 4 位数
	return fmt.Sprintf("%s%s%04d%04d", prefix, dateStr, rand.Intn(10000), rand.Intn(10000))
}

func IndexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func GetWeeklyDateTime() (*gtime.Time, *gtime.Time) {
	now := gtime.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.Add(time.Duration(1-weekday) * 24 * time.Hour)
	monday = gtime.NewFromTime(time.Date(monday.Year(), time.Month(monday.Month()), monday.Day(), 0, 0, 0, 0, monday.Location()))

	sunday := now.Add(time.Duration(7-weekday) * 24 * time.Hour)
	sunday = gtime.NewFromTime(time.Date(sunday.Year(), time.Month(sunday.Month()), sunday.Day(), 23, 59, 59, 0, sunday.Location()))

	return monday, sunday
}

func GenRandomCode() string {
	// 设置时间戳种子，避免每次运行生成相同的随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成0-999999之间的随机数，不足6位前面补0
	return fmt.Sprintf("%06v", r.Intn(1000000))
}

func GenerateUserNo(length int) string {
	if length <= 0 {
		length = 6
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var userNo []byte
	// 第一位不能为0
	userNo = append(userNo, byte(r.Intn(9)+1)+'0')
	for i := 1; i < length; i++ {
		userNo = append(userNo, byte(r.Intn(10))+'0')
	}
	return string(userNo)
}
