import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { API_STATUS } from '../helpers/constant';
import axios from 'axios';

const Cancel = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const { subscription } = location.state;

    const handleClick = async (e) => {
        const result = await axios.post('http://localhost:8000/cancel-subscription', {
            subscription_id: location.state.subscription,
        });
        if (result?.data?.meta.code === API_STATUS.SUCCESS_CODE) {
            navigate("/account");
        }
    };

    return (
        <div className="max-w-md mx-auto py-8">
            <h1 className="text-3xl font-bold mb-8">Cancel</h1>
            <button
                onClick={handleClick}
                className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
            >
                Cancel
            </button>
        </div>
    );
}

export default Cancel;
