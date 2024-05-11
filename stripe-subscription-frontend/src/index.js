import React from 'react';
import { createRoot } from 'react-dom';
import App from './App';
import './index.css';
import { Elements } from '@stripe/react-stripe-js';
import { loadStripe } from '@stripe/stripe-js';

fetch('http://localhost:8000/config')
  .then((response) => response.json())
  .then((data) => {
    const stripePromise = loadStripe(data.data.publishableKey);

    const root = createRoot(document.getElementById('root'));
    root.render(
      <React.StrictMode>
        <Elements stripe={stripePromise}>
          <App />
        </Elements>
      </React.StrictMode>
    );
  })
  .catch((error) => {
    console.error('Error:', error);
  });
