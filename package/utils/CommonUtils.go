package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTsecret = []byte(os.Getenv("Secret_key"))

type Claims struct {
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// 生成Token
func GenerateToken(id uint, username, password string, role string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)
	claims := Claims{
		Id:       id,
		UserName: username,
		Password: password,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "to-do-list",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JWTsecret)
	fmt.Println(err)
	return token, err
}

// 校验Token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTsecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// 生成 MD5 哈希值
func GenerateMD5Hash(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ParseID(s string) uint {
	id, _ := strconv.ParseUint(s, 10, 32)
	return uint(id)
}

func Contains(a, b string) bool {
	return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}

func PercentChange(current, previous int) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0
		}
		return 0.0
	}
	return (float64(current-previous) / float64(previous)) * 100
}
