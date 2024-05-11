import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const PayInvoice = ({ invoice }) => {
    const [paymentInProgress, setPaymentInProgress] = useState(false);
  
    const payInvoice = async () => {
      try {
        setPaymentInProgress(true);
        const response = await axios.post('http://localhost:8000/pay-invoice', { invoice_id: invoice.id });
      } catch (error) {
        console.error('Error paying invoice:', error);
      } finally {
        setPaymentInProgress(false);
      }
    };
  
    return (
      <button onClick={payInvoice} disabled={paymentInProgress || invoice.status !== 'open'}>
        {paymentInProgress ? 'Processing...' : 'Pay'}
      </button>
    );
  };

export default function UpcomingInvoices() {
    const userId = localStorage.getItem("userId");
    const [invoices, setInvoices] = useState([]);
    const navigate = useNavigate();

    useEffect(() => {
        const getUpcomingInvoices = async () => {
            try {
                const result = await axios.post('http://localhost:8000/upcoming-invoices', {
                    customer_id: userId,
                });
                setInvoices(result?.data?.data);
            } catch (error) {
                console.error('Error fetching upcoming invoices:', error);
            }
        };
        getUpcomingInvoices();
    }, []);

    return (
        <div className="max-w-lg mx-auto py-8">
            <h1 className="text-3xl font-bold mb-8">Upcoming Invoices</h1>

            <div className="grid gap-4">
                {invoices.map((invoice, index) => (
                    <div key={index} className="border border-gray-300 rounded-lg p-4">
                        <h3 className="text-lg font-semibold mb-2"></h3>
                        <a href={`https://dashboard.stripe.com/test/invoices/${invoice.id}`} className="text-blue-500 hover:underline">
                            {invoice.id}
                        </a>

                        <p className="text-gray-600 mb-2">Amount Due: ${invoice.amount_due / 100}</p>

                        <p className="text-gray-600 mb-2">Status: {invoice.status}</p>
                        {invoice.status === 'open' && <PayInvoice invoice={invoice} />}
                    </div>
                ))}
            </div>
        </div>
    );
}
