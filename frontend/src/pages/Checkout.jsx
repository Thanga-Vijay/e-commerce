import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { orderService } from '../api/order';
import toast from 'react-hot-toast';

const Checkout = () => {
  const { cart, clearCart } = useCart();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    shippingAddress: {
      street: '',
      city: '',
      state: '',
      zipCode: '',
      country: '',
    },
    billingAddress: {
      street: '',
      city: '',
      state: '',
      zipCode: '',
      country: '',
    },
    sameAsShipping: true,
  });

  const handleChange = (addressType, field, value) => {
    setFormData({
      ...formData,
      [addressType]: {
        ...formData[addressType],
        [field]: value,
      },
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const orderData = {
        shippingAddress: formData.shippingAddress,
        billingAddress: formData.sameAsShipping
          ? formData.shippingAddress
          : formData.billingAddress,
      };

      const response = await orderService.createOrder(orderData);
      await clearCart();
      toast.success('Order placed successfully!');
      navigate(`/orders/${response.data.id}`);
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to place order';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  if (!cart || cart.items.length === 0) {
    navigate('/cart');
    return null;
  }

  const tax = cart.totalAmount * 0.08; // 8% tax
  const shipping = 10.0;
  const total = cart.totalAmount + tax + shipping;

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-4xl font-bold mb-8">Checkout</h1>

      <form onSubmit={handleSubmit} className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Checkout Form */}
        <div className="lg:col-span-2 space-y-6">
          {/* Shipping Address */}
          <div className="card">
            <h2 className="text-2xl font-bold mb-6">Shipping Address</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="md:col-span-2">
                <label className="block text-sm font-medium mb-2">Street Address</label>
                <input
                  type="text"
                  required
                  className="input"
                  value={formData.shippingAddress.street}
                  onChange={(e) => handleChange('shippingAddress', 'street', e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">City</label>
                <input
                  type="text"
                  required
                  className="input"
                  value={formData.shippingAddress.city}
                  onChange={(e) => handleChange('shippingAddress', 'city', e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">State</label>
                <input
                  type="text"
                  required
                  className="input"
                  value={formData.shippingAddress.state}
                  onChange={(e) => handleChange('shippingAddress', 'state', e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">ZIP Code</label>
                <input
                  type="text"
                  required
                  className="input"
                  value={formData.shippingAddress.zipCode}
                  onChange={(e) => handleChange('shippingAddress', 'zipCode', e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">Country</label>
                <input
                  type="text"
                  required
                  className="input"
                  value={formData.shippingAddress.country}
                  onChange={(e) => handleChange('shippingAddress', 'country', e.target.value)}
                />
              </div>
            </div>
          </div>

          {/* Billing Address */}
          <div className="card">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold">Billing Address</h2>
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  checked={formData.sameAsShipping}
                  onChange={(e) =>
                    setFormData({ ...formData, sameAsShipping: e.target.checked })
                  }
                  className="rounded"
                />
                <span className="text-sm">Same as shipping</span>
              </label>
            </div>

            {!formData.sameAsShipping && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium mb-2">Street Address</label>
                  <input
                    type="text"
                    required
                    className="input"
                    value={formData.billingAddress.street}
                    onChange={(e) => handleChange('billingAddress', 'street', e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">City</label>
                  <input
                    type="text"
                    required
                    className="input"
                    value={formData.billingAddress.city}
                    onChange={(e) => handleChange('billingAddress', 'city', e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">State</label>
                  <input
                    type="text"
                    required
                    className="input"
                    value={formData.billingAddress.state}
                    onChange={(e) => handleChange('billingAddress', 'state', e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">ZIP Code</label>
                  <input
                    type="text"
                    required
                    className="input"
                    value={formData.billingAddress.zipCode}
                    onChange={(e) => handleChange('billingAddress', 'zipCode', e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Country</label>
                  <input
                    type="text"
                    required
                    className="input"
                    value={formData.billingAddress.country}
                    onChange={(e) => handleChange('billingAddress', 'country', e.target.value)}
                  />
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Order Summary */}
        <div className="lg:col-span-1">
          <div className="card sticky top-20">
            <h2 className="text-2xl font-bold mb-6">Order Summary</h2>

            <div className="space-y-3 mb-6">
              <div className="flex justify-between">
                <span className="text-gray-600">Subtotal</span>
                <span className="font-semibold">${cart.totalAmount.toFixed(2)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Tax (8%)</span>
                <span className="font-semibold">${tax.toFixed(2)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Shipping</span>
                <span className="font-semibold">${shipping.toFixed(2)}</span>
              </div>
              <div className="border-t pt-3 flex justify-between text-lg">
                <span className="font-bold">Total</span>
                <span className="font-bold text-primary-600">${total.toFixed(2)}</span>
              </div>
            </div>

            <button type="submit" disabled={loading} className="btn btn-primary w-full">
              {loading ? 'Placing Order...' : 'Place Order'}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
};

export default Checkout;
