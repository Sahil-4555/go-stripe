package crypto

import (
	"stripe-subscription/configs/middleware"
	"stripe-subscription/shared/log"
	"time"
)

type UserTokenData struct {
	Id        uint
	CreatedAt time.Time
}

func (u *UserTokenData) TimeStamp() {
	u.CreatedAt = time.Now()
}

func GenerateAuthToken(tokenData UserTokenData) string {
	tokenData.TimeStamp()
	token, err := middleware.GenerateToken(&tokenData)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
	}
	return token
}
