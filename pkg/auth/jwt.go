package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fuloge/basework/api"
	cfg "github.com/fuloge/basework/configs"
	"github.com/fuloge/basework/pkg/log"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type JWT struct {
	signingKey string
	subject    string //主题
}

var authLog = log.New()

type Token struct {
	Token string `json:"token"`
}

func New() (jwt *JWT) {
	jwt = &JWT{
		signingKey: cfg.EnvConfig.Authkey.Key,
		subject:    cfg.EnvConfig.Authkey.Subject,
	}

	return
}

func (j *JWT) CreateToken(userid int64, exptime int64) (res Token, errno *api.Errno) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	//claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix() //过期时间
	claims["exp"] = exptime           //过期时间
	claims["iat"] = time.Now().Unix() //签发时间
	claims["sub"] = j.subject         //主题
	claims["uid"] = strconv.FormatInt(userid, 10)
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(j.signingKey))
	if err != nil {
		fmt.Print("Error while signing the token")
		authLog.Error("CreateToken", zap.String("Error while signing the token", err.Error()))
		errno = api.AuthErr
		return
	}

	res = Token{tokenString}
	return
}

func (j *JWT) ParseToken(tokenString string) (jwt.MapClaims, *api.Errno) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse token err %v", token.Header["alg"])
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		authLog.Error("ParseToken", zap.String("parse token error", err.Error()))
		return nil, api.AuthParseErr
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { // 校验token
		return claims, nil
	}

	return nil, api.AuthParseErr
}

//if token is invalid, method will return true
func (j *JWT) TokenIsInvalid(tokenString string) bool {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		authLog.Error("TokenIsInvalid", zap.String("valid token error", api.AuthParseErr.Message))
	} else {
		//校验下token是否过期
		if res := claims.VerifyExpiresAt(time.Now().Unix(), true); res == false {
			authLog.Error("TokenIsInvalid", zap.String("token is expired", api.AuthExp.Message))
			return true
		}

		if res := claims.VerifyIssuedAt(time.Now().Unix(), true); res == false {
			authLog.Error("TokenIsInvalid", zap.String("token is expired", api.AuthExp.Message))
			return true
		}

		if res := claims["sub"].(string); res != j.subject {
			return true
		}

		uid := claims["uid"].(string)
		if res, _ := strconv.ParseInt(uid, 10, 64); res == 0 {
			return true
		}
	}

	return true
}
