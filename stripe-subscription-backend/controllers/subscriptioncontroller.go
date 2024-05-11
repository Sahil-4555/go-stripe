package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"stripe-subscription/configs"
	"stripe-subscription/models"
	"stripe-subscription/shared/common"
	"stripe-subscription/shared/log"
	"stripe-subscription/shared/message"
	"stripe-subscription/validators"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
)

func HandleConfig(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(HandleConfig).")

	stripe.Key = configs.StripeSecretKey()

	params := &stripe.PriceListParams{
		LookupKeys: stripe.StringSlice([]string{"basic_subscription", "premium_subscription", "enterprise_subscription"}),
	}

	prices := make([]*stripe.Price, 0)

	i := price.List(params)

	for i.Next() {
		prices = append(prices, i.Price())
	}

	data := map[string]interface{}{
		"publishableKey": configs.StripePublishableKey(),
		"prices":         prices,
	}

	data = map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyFetchedProductPrices,
		"data":    data,
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyFetchedProductPrices)
}

func HandleCreateSubscription(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(HandleCreateSubscription).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req struct {
		PriceId    string `json:"price_id"`
		CustomerId string `json:"customer_id"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "HandleCreateSubscription"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	conn := configs.NewConnection()
	var user models.Customer
	id, _ := strconv.Atoi(req.CustomerId)
	conn.GetDB().WithContext(ctx).Where(&models.Customer{Id: uint(id)}).First(&user)
	if user.Id == 0 {
		log.GetLog().Error("ERROR : ", message.NoSuchUserExist)
		data := map[string]interface{}{
			"message":  message.NoSuchUserExist,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(user.StripeId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(req.PriceId),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
	}

	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	s, err := sub.New(subscriptionParams)

	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	data := map[string]interface{}{
		"subscriptionId": s.ID,
		"clientSecret":   s.LatestInvoice.PaymentIntent.ClientSecret,
	}

	data = map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyCreatedSubscription,
		"data":    data,
	}
	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyCreatedSubscription)
}

func HandleCancelSubscription(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(HandleCancelSubscription).")

	var req struct {
		SubscriptionID string `json:"subscription_id"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "HandleCancelSubscription"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	s, err := sub.Cancel(req.SubscriptionID, nil)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	data := map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyCancelSubscription,
		"data":    s,
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyCancelSubscription)
}

func HandleInvoicePreview(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(HandleInvoicePreview).")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req struct {
		CustomerId        string `json:"customer_id"`
		SubscriptionId    string `json:"subscription_id"`
		NewPriceLookupKey string `json:"new_price_lookup_key"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "HandleInvoicePreview"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	conn := configs.NewConnection()
	var user models.Customer
	id, _ := strconv.Atoi(req.CustomerId)
	conn.GetDB().WithContext(ctx).Where(&models.Customer{Id: uint(id)}).First(&user)
	if user.Id == 0 {
		log.GetLog().Error("ERROR : ", message.NoSuchUserExist)
		data := map[string]interface{}{
			"message":  message.NoSuchUserExist,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	s, err := sub.Get(req.SubscriptionId, nil)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	params := &stripe.InvoiceParams{
		Customer:     stripe.String(user.StripeId),
		Subscription: stripe.String(req.SubscriptionId),
		SubscriptionItems: []*stripe.SubscriptionItemsParams{{
			ID:    stripe.String(s.Items.Data[0].ID),
			Price: stripe.String(os.Getenv(req.NewPriceLookupKey)),
		}},
	}

	in, err := invoice.GetNext(params)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	data := map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyGetInvoicesForSubscription,
		"data":    in,
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyGetInvoicesForSubscription)
}

func HandleUpdateSubscription(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(HandleUpdateSubscription).")

	var req struct {
		SubscriptionID    string `json:"subscriptionId"`
		NewPriceLookupKey string `json:"newPriceLookupKey"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "HandleUpdateSubscription"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	newPriceID := os.Getenv(strings.ToUpper(req.NewPriceLookupKey))

	s, err := sub.Get(req.SubscriptionID, nil)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{{
			ID:    stripe.String(s.Items.Data[0].ID),
			Price: stripe.String(newPriceID),
		}},
	}

	updatedSubscription, err := sub.Update(req.SubscriptionID, params)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"code":     common.META_FAILED,
			"message":  err.Error(),
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	data := map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyUpdateSubscription,
		"data":    updatedSubscription,
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyUpdateSubscription)
}

func HandleListSubscription(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(HandleListSubscription).")

	stripe.Key = configs.StripeSecretKey()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req struct {
		CustomerId string `json:"customer_id"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "HandleListSubscription"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	conn := configs.NewConnection()
	var user models.Customer
	id, _ := strconv.Atoi(req.CustomerId)
	conn.GetDB().WithContext(ctx).Where(&models.Customer{Id: uint(id)}).First(&user)
	if user.Id == 0 {
		log.GetLog().Info("ERROR : ", message.NoSuchUserExist)
		data := map[string]interface{}{
			"message":  message.NoSuchUserExist,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	params := &stripe.SubscriptionListParams{
		Customer: user.StripeId,
		Status:   "all",
	}

	params.AddExpand("data.default_payment_method")

	i := sub.List(params)

	data := map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyFetchedSubscriptionList,
		"data":    i.SubscriptionList(),
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyFetchedSubscriptionList)
}

func GetUpcommingInvoices(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(GetUpcommingInvoices).")

	stripe.Key = configs.StripeSecretKey()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req struct {
		CustomerId string `json:"customer_id"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "GetUpcommingInvoices"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	conn := configs.NewConnection()
	var user models.Customer
	id, _ := strconv.Atoi(req.CustomerId)
	conn.GetDB().WithContext(ctx).Where(&models.Customer{Id: uint(id)}).First(&user)
	if user.Id == 0 {
		log.GetLog().Error("ERROR : ", message.NoSuchUserExist)
		data := map[string]interface{}{
			"message":  message.NoSuchUserExist,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	params := &stripe.InvoiceListParams{
		Customer: stripe.String(user.StripeId),
	}

	iter := invoice.List(params)

	var data_invoices []interface{}
	for iter.Next() {
		invoice := iter.Invoice()
		data_invoices = append(data_invoices, invoice)
	}

	data := map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyFetchedTheUpcomingInvoices,
		"data":    data_invoices,
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyFetchedTheUpcomingInvoices)
}

// You could able to pay the invoice if any default payment method is set for the customer.
func PayInvoice(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(PayInvoice).")

	stripe.Key = configs.StripeSecretKey()

	var req struct {
		InvoiceId string `json:"invoice_id"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "PayInvoice"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	params := &stripe.InvoicePayParams{}
	inv, err := invoice.Pay(req.InvoiceId, params)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	data := map[string]interface{}{
		"code":    common.META_SUCCESS,
		"message": message.SuccessfullyPaidTheInvoice,
		"data":    inv.Status,
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullyPaidTheInvoice)
}

func SetPaymentDefaultForCustomer(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(SetPaymentDefaultForCustomer).")

	stripe.Key = configs.StripeSecretKey()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var req struct {
		CustomerId      string `json:"customer_id"`
		PaymentMethodId string `json:"payment_method_id"`
	}

	if c.BindJSON(&req) != nil {
		log.GetLog().Error("ERROR : ", message.FailedToReadBody)
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "SetPaymentDefaultForCustomer"); !ok {
		log.GetLog().Error("ERROR : ", resp)
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	conn := configs.NewConnection()
	var user models.Customer
	id, _ := strconv.Atoi(req.CustomerId)
	conn.GetDB().WithContext(ctx).Where(&models.Customer{Id: uint(id)}).First(&user)
	if user.Id == 0 {
		log.GetLog().Info("INFO : ", message.NoSuchUserExist)
		data := map[string]interface{}{
			"message":  message.NoSuchUserExist,
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	cu, err := customer.Get(user.StripeId, nil)
	if err != nil {
		log.GetLog().Info("INFO : ", err.Error())
		data := map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	if cu.InvoiceSettings.DefaultPaymentMethod != nil {
		log.GetLog().Info("INFO : ", message.SuccessfullySetTheDeafultPaymentMethod)
		data := map[string]interface{}{
			"message":  message.SuccessfullySetTheDeafultPaymentMethod,
			"code":     common.META_SUCCESS,
			"res_code": common.STATUS_OK,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	params := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(req.PaymentMethodId),
		},
	}

	cust, err := customer.Update(user.StripeId, params)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		data := map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_BAD_REQUEST,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	data := map[string]interface{}{
		"message": message.SuccessfullySetTheDeafultPaymentMethod,
		"code":    common.META_SUCCESS,
		"data":    cust,
	}

	statusCode := common.GetHTTPStatusCode(data["res_code"])
	common.Respond(c, statusCode, data)
	log.GetLog().Info("INFO : ", message.SuccessfullySetTheDeafultPaymentMethod)
	return
}

func HandleWebhook(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Stripe controller called(HandleWebhook).")
	stripe.Key = configs.StripeSecretKey()

	payload, err := c.GetRawData()
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"message":  fmt.Sprintf("Error reading request body: %v\n", err),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_SERVICE_UNREACHABLE,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	endpointSecret := configs.StripeWebhookKey()
	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		data := map[string]interface{}{
			"message":  fmt.Sprintf("Webhook signature verification failed. %v\n", err),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_UNAUTHORIZED,
		}
		statusCode := common.GetHTTPStatusCode(data["res_code"])
		common.Respond(c, statusCode, data)
		return
	}

	switch event.Type {
	case "customer.subscription.created":
		log.GetLog().Info("INFO : ", fmt.Sprintf("Event Called: %v", event.Type))
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.GetLog().Info("INFO : ", fmt.Sprintf("Subscription created for customer: %s", subscription.Customer.ID))

	case "invoice.payment_succeeded":
		log.GetLog().Info("INFO : ", fmt.Sprintf("Event Called: %v", event.Type))
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.GetLog().Info("INFO : ", fmt.Sprintf("Payment succeeded for invoice: %s", invoice.ID))

	case "payment_intent.succeeded":
		log.GetLog().Info("INFO : ", fmt.Sprintf("Event Called: %v", event.Type))
		var payment stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &payment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.GetLog().Info("INFO : ", fmt.Sprintf("Payment succeeded of amount: %v", payment.Amount))

	case "invoice.payment_failed":
		log.GetLog().Info("INFO : ", fmt.Sprintf("Event Called: %v", event.Type))
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		log.GetLog().Error("ERROR : ", fmt.Sprintf("Payment failed for invoice: %s", invoice.ID))

	default:
		log.GetLog().Info("INFO : ", fmt.Sprintf("Unhandled event type: %s\n", event.Type))
	}

	c.Status(http.StatusOK)
}
