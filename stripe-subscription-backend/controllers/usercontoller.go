package controllers

import (
	"net/http"
	service "stripe-subscription/services"
	"stripe-subscription/shared/common"
	"stripe-subscription/shared/log"
	"stripe-subscription/shared/message"
	"stripe-subscription/validators"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	log.GetLog().Info("INFO : ", "User Controller Called(SignUp).")

	var req common.UserRegReq

	// Decode the request body into struct and failed if any error occur
	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	// struct field validation
	if resp, ok := validators.ValidateStruct(req, "SignUpRequest"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	// call service
	resp := service.SignUp(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Response_SignIn(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "SignUp successfully...")
}

func SignIn(c *gin.Context) {
	log.GetLog().Info("INFO : ", "User Controller Called(SignIn).")

	var req common.UseLoginReq

	// Decode the request body into struct and failed if any error occur
	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	// struct field validation
	if resp, ok := validators.ValidateStruct(req, "SignInRequest"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.SignIn(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Response_SignIn(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "SignIn successfully...")
}
