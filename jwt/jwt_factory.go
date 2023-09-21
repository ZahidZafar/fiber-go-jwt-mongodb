package jwt

import (
	"greens-basket/data"
	"greens-basket/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTFactory struct {
	Config *utils.Config
}

func (jf *JWTFactory) CreateJWTTempToken(user *data.UserPrincipal) (string, error) {
	exp := time.Now().Add(jf.Config.TempTokenDuration).Unix()
	return jf.buildToken(user.Subject, exp, utils.TempToken, user.Roles)
}

func (jf *JWTFactory) CreateJWTAccessToken(user *data.UserPrincipal) (string, error) {
	exp := time.Now().Add(time.Minute * jf.Config.AccessTokenDuration).Unix()
	return jf.buildToken(user.Subject, exp, utils.AccessToken, user.Roles)
}

func (jf *JWTFactory) CreateJWTRefreshToken(user *data.UserPrincipal) (string, error) {
	exp := time.Now().Add(time.Minute * jf.Config.RefreshTokenDuration).Unix()
	return jf.buildToken(user.Subject, exp, utils.RefreshToken, user.Roles)
}

func (jf *JWTFactory) buildToken(sub string, exp int64, scope string, roles string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims.Valid()

	claims := token.Claims.(jwt.MapClaims)

	claims[utils.Subject] = sub
	claims[utils.Scope] = scope
	claims[utils.Roles] = roles
	claims[utils.Expiration] = exp

	t, err := token.SignedString([]byte(jf.Config.TokenSecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}
