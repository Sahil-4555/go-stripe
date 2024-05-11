import React, { useState, useEffect, startTransition } from "react";
import axios from "axios";
import { useNavigate, Link } from "react-router-dom";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { API_STATUS } from "../../helpers/constant";


export default function Login() {
  const navigate = useNavigate();
  const toastOptions = {
    position: "top-right",
    autoClose: 3000,
    pauseOnHover: true,
    draggable: true,
    theme: "dark",
  };
  const [values, setValues] = useState({ email: "", password: "" });

  useEffect(() => {
    const token = localStorage.getItem("token_");
    if (token) {
      navigate("/");
    }
  }, [navigate]);

  const handleChange = (event) => {
    setValues({ ...values, [event.target.name]: event.target.value });
  };

  const validateForm = () => {
    const { email, password } = values;
    if (email === "" || password === "") {
      toast.error("Email and Password are required.", toastOptions);
      return false;
    }
    return true;
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (validateForm()) {
      const { email, password } = values;

      try {
        await new Promise((resolve) => {
          startTransition(() => {
                axios
                .post("http://localhost:8000/v1/login", { email, password })
                .then(({ data }) => {
                    if (data.code === API_STATUS.FAILURE_CODE) {
                    toast.error(data.message, toastOptions);
                    } else {
                    localStorage.setItem("token_", data.meta.token);
                    localStorage.setItem("userId", data.data._id);
                    localStorage.setItem("name", data.data.name)
                    resolve();
                    startTransition(() => navigate("/prices"));
                    }
                })
              .catch((error) => {   
                  console.log("Login Failed: ", error.message);
                if (error.response && error.response.status === 500) {
                  toast.error(
                    "Backend is currently unavailable. Please try again later.",
                    toastOptions
                  );
                } else {
                  toast.error(
                    "An error occurred during login. Please try again.",
                    toastOptions
                  );
                  navigate("/error");
                }
                resolve();
              });
          });
        });
      } catch (error) {
        console.error("Error during login:", error);
      }
    }
  };

  return (
    <>
      <div className="form-container bg-white shadow-md rounded-md p-8 mx-auto mt-8 max-w-md">
        <form
          onSubmit={(event) => handleSubmit(event)}
          className="space-y-6"
        >
          <h2 className="text-3xl font-semibold text-center">Log In</h2>
          <input
            type="text"
            placeholder="Email"
            autoComplete="off"
            name="email"
            onChange={(e) => handleChange(e)}
            className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
          />
          <input
            type="password"
            placeholder="Password"
            name="password"
            onChange={(e) => handleChange(e)}
            className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
          />
          <button
            type="submit"
            className="w-full bg-blue-500 text-white px-4 py-2 rounded-full hover:bg-blue-700 transition duration-300"
          >
            Log In
          </button>
          <p className="text-center text-gray-700">
            Don't have an account?{' '}
            <Link
              to="/signup"
              className="text-blue-500 font-semibold hover:underline"
            >
              Create One.
            </Link>
          </p>
        </form>
      </div>
      <ToastContainer />
    </>
  );
}