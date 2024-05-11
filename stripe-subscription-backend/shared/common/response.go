package common

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	Id    uint   `json:"id" structs:"id"`
	Email string `json:"email" structs:"email"`
	Name  string `json:"name" structs:"name"`
}

type BaseSuccessResponse struct {
	Data interface{} `json:"data" bson:"data"  structs:"data"`
	Meta Meta        `json:"meta" bson:"meta"  structs:"meta"`
}

type Meta struct {
	Message string `json:"message" structs:"message"`
	Code    int    `json:"code" structs:"code"`
	Token   string `json:"token,omitempty" structs:"token,omitempty" bson:"token,omitempty"`
}

func MessageWithCode(status int, message string) map[string]interface{} {
	return map[string]interface{}{"res_code": status, "message": message}
}

func ResponseErrorWithCode(status int, message string) map[string]interface{} {
	return MessageWithCode(status, message)
}

func GetHTTPStatusCode(resCode interface{}) int {
	if resCode != nil {
		return resCode.(int)
	}
	return http.StatusOK
}

func ResponseSuccessWithToken(message string, code int, resData map[string]interface{}) map[string]interface{} {
	response := BaseSuccessResponse{
		Meta: Meta{
			Message: message,
			Code:    code,
		},
	}
	if token, ok := resData["token"]; ok {
		response.Meta.Token = token.(string)
	}
	if rData, ok := resData["data"]; ok {
		response.Data = rData

	} else {
		response.Data = nil
	}

	m := structs.Map(response)

	return m
}

func ConvertToInterface(message string, code int, data interface{}) map[string]interface{} {
	d := map[string]interface{}{
		"message": message,
		"code":    code,
		"data":    data,
	}
	d = FinalResponse(d)
	return d
}

func Response(resData interface{}) map[string]interface{} {
	data := resData.(map[string]interface{})
	response := BaseSuccessResponse{
		Meta: Meta{
			Message: data["message"].(string),
			Code:    data["code"].(int),
		},
	}
	if resData != nil {
		if rData, ok := data["data"]; ok {
			response.Data = rData

		} else {
			response.Data = nil
		}
	} else {
		response.Data = nil
	}

	m := structs.Map(response)
	return m
}

func ResponseSuccessWithCode(message string, data ...interface{}) map[string]interface{} {
	return ConvertToInterface(message, META_SUCCESS, data)
}

func FinalResponse(data map[string]interface{}) map[string]interface{} {
	response := BaseSuccessResponse{
		Meta: Meta{
			Message: data["message"].(string),
			Code:    data["code"].(int),
		},
	}

	if rData, ok := data["data"]; ok {
		response.Data = rData

	} else {
		response.Data = nil
	}

	m := structs.Map(response)
	return m
}

func Response_SignIn(c *gin.Context, status int, data map[string]interface{}) {
	if status != 200 {
		data = FinalResponse(data)
	}
	c.JSON(status, data)
}

func Respond(c *gin.Context, status int, data map[string]interface{}) {
	d := FinalResponse(data)
	c.JSON(status, d)
}
