import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';

const Middleware = () => {
  const auth_token = localStorage.getItem('token_');
  return auth_token ? <Outlet /> : <Navigate to='/login' />;
};

export default Middleware;
