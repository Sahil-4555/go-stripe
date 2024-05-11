import React, { useState, useEffect, startTransition } from "react";
import { useNavigate, Link } from "react-router-dom";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";



export default function Signup() {
  const navigate = useNavigate();
  const toastOptions = {
    position: "top-right",
    autoClose: 3000,
    pauseOnHover: true,
    draggable: true,
    theme: "dark",
  };
  const [values, setValues] = useState({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  useEffect(() => {
    const token = localStorage.getItem("token_");
    if (token) {
      navigate("/");
    }
  }, [navigate]);

  const handleChange = (event) => {
    setValues({ ...values, [event.target.name]: event.target.value });
  };

  const handleValidation = () => {
    const { name, password, confirmPassword, email } = values;
    if (password !== confirmPassword) {
      toast.error(
        "Password and confirm password should be the same.",
        toastOptions
      );
      return false;
    } else if (name === "") {
      toast.error("Name is required.", toastOptions);
      return false;
    } else if (password.length < 6) {
      toast.error(
        "Password should be equal or greater than 6 characters.",
        toastOptions
      );
      return false;
    } else if (email === "") {
      toast.error("Email is required.", toastOptions);
      return false;
    }
    return true;
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (handleValidation()) {

      navigate('/address-form', {
        state: {
          email: values.email,
          password: values.password,
          name: values.name,
        }
      });
    }
  };

  return (
    <>
      <div className="form-container bg-white shadow-md rounded-md p-8 mx-auto mt-8 max-w-md">
        <form
          action=""
          onSubmit={(event) => handleSubmit(event)}
          className="space-y-6"
        >
          <h2 className="text-3xl font-semibold text-center">Create an Account</h2>
          <input
            type="text"
            placeholder="Name"
            autoComplete="off"
            name="name"
            onChange={(e) => handleChange(e)}
            className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
          />
          <input
            type="email"
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
          <input
            type="password"
            placeholder="Confirm Password"
            name="confirmPassword"
            onChange={(e) => handleChange(e)}
            className="w-full px-4 py-2 border border-gray-300 rounded-md placeholder-gray-500 focus:outline-none focus:border-blue-500"
          />
          <button
            type="submit"
            className="w-full bg-blue-500 text-white px-4 py-2 rounded-full hover:bg-blue-700 transition duration-300"
          >
            Next
          </button>
          <p className="text-center text-gray-700">
            Already have an account?{' '}
            <Link
              to="/login"
              className="text-blue-500 font-semibold hover:underline"
            >
              Login.
            </Link>
          </p>
        </form>
      </div>
      <ToastContainer />
    </>
  );
}