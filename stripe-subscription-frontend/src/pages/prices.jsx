import React, { useState, useEffect, startTransition } from 'react';

import axios from 'axios';
import { API_STATUS } from '../helpers/constant';
import { useNavigate } from 'react-router-dom';

export default function Prices() {
    const userId = localStorage.getItem("userId");
    const [prices, setPrices] = useState([]);
    const [subscriptionData, setSubscriptionData] = useState(null);
    const navigate = useNavigate();
    useEffect(() => {
        const getPriceListing = async () => {
            const result = await axios.get('http://localhost:8000/config');
            setPrices(result?.data?.data?.prices);
        };
        getPriceListing();
    }, [])

    const createSubscription = async (priceId) => {
        try {
            const result = await axios.post('http://localhost:8000/create-subscription', {
                price_id: priceId,
                customer_id: userId,
            });
            if (result?.data?.meta.code === API_STATUS.SUCCESS_CODE) {
                const data = result?.data.data || [];
                setSubscriptionData(data)
                navigate("/subscribe", {
                    state: {
                        subscriptionId: data.subscriptionId,
                        clientSecret: data.clientSecret,
                    }
                });
            }
        } catch (error) {
            console.log('Error : ', error);
        }
    }

   
    return (
        <div>
            <div className="p-4">
                <h1 className="text-2xl font-bold mb-4">Select a plan</h1>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {prices.map((price) => (
                        <div key={price.id} className="bg-white rounded-lg shadow-lg p-6">

                            <h3 className="text-lg font-semibold mb-2">{price.product.name}</h3>

                            <p className="text-gray-600 mb-4">
                                $ {price.unit_amount / 100} / month
                            </p>

                            <button
                                className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                                onClick={() => createSubscription(price.id)}
                            >
                                Select
                            </button>
                        </div>
                    ))}
                </div>
            </div>

        </div>
    );
}
