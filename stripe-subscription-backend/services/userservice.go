package service

import (
	"context"
	"regexp"
	"stripe-subscription/configs"
	"stripe-subscription/configs/crypto"
	"stripe-subscription/models"
	"stripe-subscription/shared/common"
	"stripe-subscription/shared/log"
	"stripe-subscription/shared/message"

	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"golang.org/x/crypto/bcrypt"
)

// call service
func SignUp(req common.UserRegReq) map[string]interface{} {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := configs.NewConnection()

	var IsValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,4}$`).MatchString

	if !IsValidEmail(req.Email) {
		log.GetLog().Info("ERROR : ", "Invalid email format")
		return map[string]interface{}{
			"message": message.EmailInvalid,
			"code":    common.META_FAILED,
		}
	}

	// Check if the user with same email is not there
	var user models.Customer
	conn.GetDB().Where(&models.Customer{Email: req.Email}).First(&user)
	if user.Id != 0 {
		log.GetLog().Info("ERROR : ", "Email in use.")
		return map[string]interface{}{
			"message": message.EmailInUse,
			"code":    common.META_FAILED,
		}
	}

	// Converts the password to hash
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		return map[string]interface{}{
			"message":  message.FailedToHashPassword,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST}
	}

	stripe.Key = configs.StripeSecretKey()

	params := &stripe.CustomerParams{
		Email: stripe.String(req.Email),
		Name:  stripe.String(req.Name),
		Address: &stripe.AddressParams{
			Line1:      stripe.String(req.Address),
			PostalCode: stripe.String(req.PostalCode),
			City:       stripe.String(req.City),
			State:      stripe.String(req.State),
			Country:    stripe.String(req.Country),
		},
	}

	u, err := customer.New(params)
	if err != nil {
		return map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
	}

	// Insert the user details in DB
	user = models.Customer{
		Name:       req.Name,
		Email:      req.Email,
		Address:    req.Address,
		City:       req.City,
		State:      req.State,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		StripeId:   u.ID,
		Password:   string(hashPassword),
	}

	result := conn.GetDB().WithContext(ctx).Create(&user)
	if result.Error != nil {
		log.GetLog().Info("ERROR(Query) : ", result.Error.Error())
		return map[string]interface{}{
			"message":  message.FailedToInsert,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST}
	}

	newId := user.Id
	tokenData := crypto.UserTokenData{
		Id: user.Id,
	}

	token := crypto.GenerateAuthToken(tokenData)
	loginData := common.LoginResponse{
		Id:    newId,
		Name:  user.Name,
		Email: user.Email,
	}

	data := map[string]interface{}{
		"data":     loginData,
		"customer": u,
	}

	data = map[string]interface{}{
		"token": token,
		"data":  data,
	}

	response := common.ResponseSuccessWithToken(message.SignUpSuccess, common.META_SUCCESS, data)

	return response
}

func SignIn(req common.UseLoginReq) map[string]interface{} {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var IsValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,4}$`).MatchString
	if !IsValidEmail(req.Email) {
		log.GetLog().Info("ERROR : ", "Invalid email format.")
		return map[string]interface{}{
			"message": message.EmailInvalid,
			"code":    common.META_FAILED,
		}
	}

	conn := configs.NewConnection()
	var user models.Customer
	conn.GetDB().WithContext(ctx).Where(&models.Customer{Email: req.Email}).First(&user)
	if user.Id == 0 {
		log.GetLog().Info("ERROR : ", message.EmailOrPasswordNotMatched)
		return map[string]interface{}{
			"message": message.EmailOrPasswordNotMatched,
			"code":    common.META_FAILED,
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		return map[string]interface{}{
			"message": message.EmailOrPasswordNotMatched,
			"code":    common.META_FAILED,
		}
	}

	tokenData := crypto.UserTokenData{
		Id: user.Id,
	}

	token := crypto.GenerateAuthToken(tokenData)
	loginData := common.LoginResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	data := map[string]interface{}{
		"token": token,
		"data":  loginData,
	}

	response := common.ResponseSuccessWithToken(message.LoginSuccess, common.META_SUCCESS, data)

	return response
}
