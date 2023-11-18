package nineauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/NineKanokpol/Nine-shop-test/config"
	"github.com/NineKanokpol/Nine-shop-test/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apiKey"
)

type nineAuth struct {
	mapClaims *nineMapClaims
	cfg       config.IJwtConfig
}

type nineMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type INineAuth interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	//* time.Duration มีหน่วยเป็น nano sec ถ้าเอา sec เข้าไปต้อง *10^9
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *nineAuth) SignToken() string {
	//sign token คู่ payload NewWithClaims
	//asimmatic พวก RHA symmatic key key เดียวใช้ทั้ง encryet decryte
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*nineMapClaims, error) {
	//ParsewithClaims มี payload
	//sign token แบบ HMAC
	token, err := jwt.ParseWithClaims(tokenString, &nineMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		//*ตรวจ algrorithum การ sign token
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expried")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	//* แปลง any -> type อื่นใดๆ ต้องทำแบบนี้
	if claims, ok := token.Claims.(*nineMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &nineAuth{
		cfg: cfg,
		mapClaims: &nineMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "nineshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

// factory
func NewNineAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (INineAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	default:
		return nil, fmt.Errorf("unknow token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) INineAuth {
	return &nineAuth{
		cfg: cfg,
		mapClaims: &nineMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "nineshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) INineAuth {
	return &nineAuth{
		cfg: cfg,
		mapClaims: &nineMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "nineshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
