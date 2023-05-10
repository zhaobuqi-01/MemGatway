package logic

import (
	"encoding/base64"
	"fmt"
	"gateway/globals"
	"gateway/proxy/http_proxy/dto"
	"gateway/proxy/pkg"
	"strings"
	"time"

	"gateway/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// interface
type OAuthLogic interface {
	// Tokens
	Tokens(c *gin.Context, param *dto.TokensInput) (*dto.TokensOutput, error)
}

type oauthLogic struct{}

func NewOAuthLogic() *oauthLogic {
	return &oauthLogic{}
}

// Tokens
func (o *oauthLogic) Tokens(c *gin.Context, param *dto.TokensInput) (*dto.TokensOutput, error) {
	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		return nil, fmt.Errorf("用户名或密码格式错误")
	}
	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("用户名或密码格式错误")
	}
	appList := pkg.Cache.GetAppList()
	for _, appInfo := range appList {
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,
				ExpiresAt: time.Now().Add(globals.JwtExpires * time.Second).Unix(),
			}
			token, err := utils.JwtEncode(claims)
			if err != nil {
				return nil, err
			}
			output := &dto.TokensOutput{
				ExpiresIn:   globals.JwtExpires,
				TokenType:   "Bearer",
				AccessToken: token,
				Scope:       "read_write",
			}
			return output, nil
		}
	}
	return nil, fmt.Errorf("未匹配正确APP信息")
}
