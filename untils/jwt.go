package untils

import (
	"errors"
	"github.com/sniperCore/core/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token:")
)

type Jwt struct {
	SecretKey    string
	SecretMethod string
	ExpireTime   int64
	RefreshTime  int64
}

/**
 * init jwt
 */
func NewJwt() (*Jwt, error) {
	var jwtObj *Jwt
	err := config.Conf.Get("auth")
	if err != nil {
		return nil, err
	}

	err = viper.UnmarshalKey("auth", &jwtObj)
	if err != nil {
		return nil, err
	}

	return jwtObj, nil
}

type AuthClaims struct {
	UserId   int64
	Password string
	jwt.StandardClaims
}

/**
 * 根据对应方式生成对应的token
 */
func (j *Jwt) CreateToken(claims AuthClaims) (string, error) {
	//签名生效时间
	claims.NotBefore = int64(time.Now().Unix() - 1000)
	claims.ExpiresAt = int64(time.Now().Unix() + j.ExpireTime*3600)
	claims.Issuer = "Justin Lie."

	var token *jwt.Token
	switch j.SecretMethod {
	case "HMAC256":
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	case "HMAC384":
		token = jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	case "HMAC512":
		token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	case "ECDSA256":
		token = jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	case "ECDSA384":
		token = jwt.NewWithClaims(jwt.SigningMethodES384, claims)
	case "ECDSA512":
		token = jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	case "RSA256":
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	case "RSA384":
		token = jwt.NewWithClaims(jwt.SigningMethodRS384, claims)
	case "RSA512":
		token = jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	default:
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	}

	return token.SignedString([]byte(j.SecretKey))
}

/**
 * 解析token
 */
func (j *Jwt) ParseToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

/**
 * 更新token
 */
func (j *Jwt) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(j.RefreshTime) * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
