import React, { useState, useEffect, startTransition } from "react";
import axios from "axios";
import { useNavigate, Link, useLocation } from "react-router-dom";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { API_STATUS } from "../../helpers/constant";

const AddressForm = () => {
    const navigate = useNavigate();
    const [customer, setCustomer] = useState(null);
    const toastOptions = {
        position: "top-right",
        autoClose: 3000,
        pauseOnHover: true,
        draggable: true,
        theme: "dark",
    };
    const location = useLocation();
    const value = location.state;
    const [addressData, setAddressData] = useState({
        address: '',
        city: '',
        state: '',
        postal_code: '',
        country: ''
    });

    const handleChange = (event) => {
        setAddressData({ ...addressData, [event.target.name]: event.target.value });
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        const formData = { ...value, ...addressData };
        try {
            await new Promise((resolve) => {
                startTransition(() => {
                    axios
                        .post("http://localhost:8000/v1/signup", {
                            email: formData.email,
                            name: formData.name,
                            address: formData.address,
                            city: formData.city,
                            state: formData.state,
                            postal_code: formData.postal_code,
                            country: formData.country,
                            password: formData.password,
                        })
                        .then(({ data }) => {
                            if (data.code === API_STATUS.FAILURE_CODE) {
                                toast.error(data.message, toastOptions);
                            } else {
                                localStorage.setItem("token_", data.meta.token);
                                localStorage.setItem("userId", data.data.data.id);
                                localStorage.setItem("name", data.data.data.name);
                                setCustomer(data.data.customer)
                                resolve();
                                startTransition(() => navigate("/prices"));
                            }
                        })
                        .catch((error) => {
                            console.log("Registration Failed: ", error.message);
                            if (error.response && error.response.status === 500) {
                                toast.error(
                                    "Backend is currently unavailable. Please try again later.",
                                    toastOptions
                                );
                            } else {
                                toast.error(
                                    "An error occurred during registration. Please try again.",
                                    toastOptions
                                );
                                navigate("/error");
                            }
                            resolve();
                        });
                });
            });
        } catch (error) {
            console.error("Error during registration:", error);
        }
    };

    return (
        <div className="form-container bg-white shadow-md rounded-md p-8 mx-auto mt-8 max-w-md">
            <form onSubmit={handleSubmit} className="space-y-6">
                <h2 className="text-3xl font-semibold text-center">Enter Address Details</h2>
                <input
                    type="text"
                    placeholder="Address"
                    name="address"
                    onChange={handleChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
                />
                <input
                    type="text"
                    placeholder="City"
                    name="city"
                    onChange={handleChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
                />
                <input
                    type="text"
                    placeholder="State"
                    name="state"
                    onChange={handleChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
                />
                <input
                    type="text"
                    placeholder="Postal Code"
                    name="postal_code"
                    onChange={handleChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
                />
                <input
                    type="text"
                    placeholder="Country"
                    name="country"
                    onChange={handleChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
                />
                <button
                    type="submit"
                    className="w-full bg-blue-500 text-white px-4 py-2 rounded-full hover:bg-blue-700 transition duration-300"
                >
                    Register
                </button>
            </form>
        </div>
    );
};

export default AddressForm;
