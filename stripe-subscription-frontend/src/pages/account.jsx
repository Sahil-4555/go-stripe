import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const AccountSubscription = ({ subscription, cancelSubscription }) => (
  <section className="border border-gray-300 rounded-lg p-4 mb-4">
    <h4 className="text-lg font-semibold mb-2">
      <a href={`https://dashboard.stripe.com/test/subscriptions/${subscription.ID}`} className="text-blue-500 hover:underline">
        {subscription.ID}
      </a>
    </h4>

    <p>
      Status: {subscription.Status}
    </p>

    <p>
      Card last4: {subscription.DefaultPaymentMethod?.Card?.Last4}
    </p>

    <p>
      Current period end: {(new Date(subscription.CurrentPeriodEnd * 1000).toString())}
    </p>

    <button
      onClick={() => cancelSubscription(subscription.ID)}
      className="text-red-500 hover:underline"
    >
      Cancel
    </button>
  </section>
);

const Account = () => {
  const userId = localStorage.getItem("userId");
  const [subscriptions, setSubscriptions] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const getSubscriptionList = async () => {
      try {
        const result = await axios.post('http://localhost:8000/subscriptions', {
          customer_id: userId,
        });
        setSubscriptions(result?.data?.data?.Data);
      } catch (error) {
        console.error('Error fetching subscriptions:', error);
      }
    };
    getSubscriptionList();
  }, []);

  const cancelSubscription = (subscriptionId) => {
    navigate('/cancel', { state: { subscription: subscriptionId } });
  };

  if (!subscriptions) {
    return null;
  }

  return (
    <div className="max-w-lg mx-auto py-8">
      <h1 className="text-3xl font-bold mb-8">Account</h1>

      <div className="flex justify-between mb-8">
        <a href="/prices" className="text-blue-500 hover:underline">Add a subscription</a>
        <a href="/prices" className="text-blue-500 hover:underline">Restart demo</a>
        <a href="/upcoming-invoices" className="text-blue-500 hover:underline">Upcoming invoices</a>
      </div>

      <h2 className="text-xl font-semibold mb-4">Subscriptions</h2>

      <div id="subscriptions">
        {subscriptions.map(subscription => (
          <AccountSubscription key={subscription.ID} subscription={subscription} cancelSubscription={cancelSubscription} />
        ))}
      </div>
    </div>
  );
}

export default Account;
