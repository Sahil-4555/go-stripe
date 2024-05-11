package configs

import (
	"os"
	"stripe-subscription/shared/log"

	"github.com/joho/godotenv"
)

func Username() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("DB_USER")
}

func Password() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv("DB_PASSWORD")
}

func Host() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("DB_HOST"))
}

func DBPort() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("DB_PORT"))
}

func ServerPort() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("SERVER_PORT"))
}

func DbName() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("DBNAME"))
}

func JwtApiAuthKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("JWT_API_AUTH_KEY"))
}

func StripePublishableKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("STRIPE_PUBLISHABLE_KEY"))
}

func StripeWebhookKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("STRIPE_WEBHOOK_SECRET"))
}

func StripeSecretKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("STRIPE_SECRET_KEY"))
}

func ReactStripePublishableKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error loading .env file.")
	}

	return os.Getenv(("REACT_APP_STRIPE_PUBLISHABLE_KEY"))
}
