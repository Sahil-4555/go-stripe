import React, { useState } from 'react';
import { loadStripe } from '@stripe/stripe-js';
import { useLocation, useNavigate } from 'react-router-dom';
import {
    CardElement,
    useStripe,
    useElements,
} from '@stripe/react-stripe-js';
import { ToastContainer, toast } from 'react-toastify';
import axios from 'axios';
import { API_STATUS } from '../helpers/constant';

const Subscribe = () => {
    const userId = localStorage.getItem("userId")
    const [name, setName] = useState("");
    const [messages, setMessages] = useState('');
    const [paymentIntent, setPaymentIntent] = useState();
    const navigate = useNavigate();
    const location = useLocation();
    const subscriptionData = location.state;
    const toastOptions = {
        position: "top-right",
        autoClose: 3000,
        pauseOnHover: true,
        draggable: true,
        theme: "dark",
    };

    const handleChange = (event) => {
        setName(event.target.value);
    };

    const stripe = useStripe();
    const elements = useElements();

    if (!stripe || !elements) {
        return '';
    }

    const handleValidation = () => {
        if (name === "") {
            toast.error("Name is required.", toastOptions);
            return false;
        }
        return true;
    };

    const createPaymentMethod = async () => {
        const { paymentMethod, error } = await stripe.createPaymentMethod({
            type: 'card',
            card: elements.getElement(CardElement),
            billing_details: {
                name: name,
            },
        });
        if (error) {
            toast.error(error.message, toastOptions);
            return;
        }
        console.log('Payment Method created:', paymentMethod);
        handleCardPayment(paymentMethod.id);
    };

    const handleCardPayment = async (paymentMethodId) => {
        if (handleValidation()) {
            let { error, paymentIntent } = await stripe.confirmCardPayment(subscriptionData.clientSecret, {
                payment_method: paymentMethodId,
            });
            if (error) {
                setMessages(error.message);
                return;
            }

            setPaymentIntent(paymentIntent);
            if (paymentIntent && paymentIntent.status === 'succeeded') {

                try {
                    const response = await axios.post('http://localhost:8000/webhook', {
                        eventType: 'invoice.payment_succeeded',
                        data: paymentIntent,
                    });
                    if (response.status === 200) {
                        console.log('Webhook event triggered successfully');
                    }
                } catch (error) {
                    console.error('Error triggering webhook event:', error);
                }

                try {
                    const result = await axios.post('http://localhost:8000/set-payment-default-for-customer ', {
                        payment_method_id: paymentMethodId,
                        customer_id: userId,
                    });
                    if (result?.data?.meta.code === API_STATUS.SUCCESS_CODE) {
                        navigate("/account");
                    }
                } catch (error) {
                    console.log('Error : ', error);
                }
            }

        }
    };

    return (
        <>
            <div className="max-w-md mx-auto">
                <h1 className="text-3xl font-bold mb-4">Subscribe</h1>
                <hr className="mb-4" />

                <form onSubmit={e => e.preventDefault()}>
                    <input
                        type="text"
                        placeholder="Name"
                        autoComplete="off"
                        name="name"
                        value={name}
                        onChange={(e) => handleChange(e)}
                        className="w-full px-4 py-2 mb-4 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
                    />
                    <CardElement className="border border-gray-300 rounded-md p-4 mb-4" />

                    <button className="block w-full bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                        onClick={createPaymentMethod}>
                        Subscribe
                    </button>

                    <div className="text-red-500">{messages}</div>
                </form>
            </div>
            <ToastContainer />
        </>
    );
}

export default Subscribe;
