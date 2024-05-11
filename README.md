# Fixed Subscription Demo

This demo project showcases a fixed subscription model implemented with a frontend in React, a backend in Golang, and integrated with Stripe for payment processing. It also includes integration of Stripe webhooks for handling events.

## Features

- Frontend built with React for user interface.
- Backend built with Golang to handle API requests and interact with Stripe.
- Fixed subscription model implemented using Stripe Billing.
- Integration of Stripe Elements for secure payment processing.
- Implementation of Stripe webhooks to handle events such as payment success and subscription status changes.

## Prerequisites

- npm for Frontend
- Golang for Backend
- A Stripe account. You can sign up for free [here](https://dashboard.stripe.com/register)

## Stripe CLI Integration

1. Install the Stripe CLI globally via apt

```
sudo apt install stripe
```
you can download Stripe CLI from here: [Download Stripe CLI](https://docs.stripe.com/stripe-cli)

2. Authenticate the Stripe CLI with your Stripe account:

```
stripe login
stripe fixtures seed.json
```

3. Once installed and authenticated, you can use the Stripe CLI to interact with your Stripe account, test webhooks and more. Here are some common commands:

- Start the Stripe webhook forwarding:
 
 ```
 stripe listen --forward-to http://localhost:8000/webhook
 ```
 Replace http://localhost:3000/webhook with the URL of your webhook endpoint.

- Test a webhook locally:

Use the stripe trigger command to trigger a specific webhook event for testing purposes.
```
stripe trigger payment_intent.succeeded
```

- Inspect Stripe events:

Use the stripe events command to list recent Stripe events.
```
stripe events list
```

- Create a `.env` file in the backend directory based on the provided `.env.example` file.

## Configure Stripe Webhook
Ngrok allows you to expose your local server to the internet securely. Follow these steps to install Ngrok:

- Download Ngrok from the official website: [Download Ngrok](https://ngrok.com/download)


### Start Ngrok Tunnel

```
ngrok http 5000
```
Ngrok will generate a public URL (e.g., https://randomstring.ngrok.io) that forwards HTTP traffic to your local server.

### Setup Webhook in Stripe

- Log in to your Stripe Dashboard.
- Navigate to Developers > Webhooks.
- Click on the Add endpoint button.
- Enter the Ngrok URL generated in step 2 as the endpoint URL (e.g., https://randomstring.ngrok.io/webhook).
- Select the events you want to listen to (e.g., payment success, subscription status changes).
- Click Add endpoint to save the webhook configuration.