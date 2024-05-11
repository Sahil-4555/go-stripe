package routes

import (
	"context"
	"fmt"
	"net/http"
	"stripe-subscription/configs"
	"stripe-subscription/configs/middleware"
	"stripe-subscription/controllers"
	"stripe-subscription/shared/common"
	"stripe-subscription/shared/log"

	"github.com/gin-gonic/gin"
)

var server *http.Server

func Run() {
	port := configs.ServerPort()
	log.GetLog().Info("", "Service listen on "+port)
	router := gin.New()
	server = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	common.SetupProducts()
	SetupRoutes(router)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.GetLog().Fatal("", "listen: %s\n", err)
	}
}

func Close(ctx context.Context) error {
	if server != nil {
		return server.Shutdown(ctx)
	}
	return nil
}

func SetupRoutes(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.GinMiddleware())
	public := r.Group("/v1")
	public.POST("/signup", controllers.SignUp)
	public.POST("/login", controllers.SignIn)

	private := r.Group("/")
	// private.Use(middleware.AuthHandler())
	private.GET("/config", controllers.HandleConfig)
	private.POST("/create-subscription", controllers.HandleCreateSubscription)
	private.POST("/cancel-subscription", controllers.HandleCancelSubscription)
	private.PUT("/update-subscription", controllers.HandleUpdateSubscription)
	private.GET("/invoice-preview", controllers.HandleInvoicePreview)
	private.POST("/subscriptions", controllers.HandleListSubscription)
	private.POST("/upcoming-invoices", controllers.GetUpcommingInvoices)
	private.POST("/pay-invoice", controllers.PayInvoice)
	private.POST("/set-payment-default-for-customer", controllers.SetPaymentDefaultForCustomer)
	private.POST("/webhook", controllers.HandleWebhook)
}
