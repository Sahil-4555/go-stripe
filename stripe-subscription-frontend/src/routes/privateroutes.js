import { Routes, Route, BrowserRouter, Navigate } from 'react-router-dom';
import { lazy } from 'react';
import Middleware from './middleware';
import React from 'react';
import AddressForm from '../pages/authentication/addressform';
import Prices from '../pages/prices';
import Subscribe from '../pages/subscribe';
import AccountSubscription from '../pages/account';
import Cancel from '../pages/cancel';
import UpcomingInvoices from '../pages/upcominginvoices';

const Login = lazy(() => import('../pages/authentication/login'));
const Signup = lazy(() => import('../pages/authentication/signup'));


const RouteList = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Middleware />}>
                    <Route path='/prices' exact element={<Prices />} />
                    <Route path='/subscribe' element={<Subscribe />} />
                    <Route path='/account' element={<AccountSubscription />} />
                    <Route path='/cancel' element={<Cancel />} />
                    <Route path='/upcoming-invoices' element={<UpcomingInvoices />} />
                    {/* <Route path="*" element={<Navigate to="/prices" replace />} /> */}
                </Route>
                <Route exact path="/login" element={<Login />} />
                <Route exact path="/signup" element={<Signup />} />
                <Route exact path="/address-form" element={<AddressForm />} />
            </Routes>
        </BrowserRouter>
    );
};

export default RouteList