import { createContext, useContext, useState, useEffect } from 'react';
import { cartService } from '../api/cart';
import { useAuth } from './AuthContext';
import toast from 'react-hot-toast';

const CartContext = createContext(null);

export const CartProvider = ({ children }) => {
  const [cart, setCart] = useState(null);
  const [loading, setLoading] = useState(false);
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    if (isAuthenticated) {
      fetchCart();
    } else {
      setCart(null);
    }
  }, [isAuthenticated]);

  const fetchCart = async () => {
    try {
      setLoading(true);
      const response = await cartService.getCart();
      setCart(response.data);
    } catch (error) {
      console.error('Failed to fetch cart:', error);
    } finally {
      setLoading(false);
    }
  };

  const addItem = async (productId, quantity = 1) => {
    try {
      const response = await cartService.addItem(productId, quantity);
      setCart(response.data);
      toast.success('Item added to cart');
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to add item';
      toast.error(message);
      throw error;
    }
  };

  const updateItem = async (itemId, quantity) => {
    try {
      const response = await cartService.updateItem(itemId, quantity);
      setCart(response.data);
      toast.success('Cart updated');
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to update item';
      toast.error(message);
      throw error;
    }
  };

  const removeItem = async (itemId) => {
    try {
      const response = await cartService.removeItem(itemId);
      setCart(response.data);
      toast.success('Item removed from cart');
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to remove item';
      toast.error(message);
      throw error;
    }
  };

  const clearCart = async () => {
    try {
      await cartService.clearCart();
      setCart(null);
      toast.success('Cart cleared');
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to clear cart';
      toast.error(message);
      throw error;
    }
  };

  const value = {
    cart,
    loading,
    fetchCart,
    addItem,
    updateItem,
    removeItem,
    clearCart,
    itemCount: cart?.items?.length || 0,
    totalAmount: cart?.totalAmount || 0,
  };

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};

export const useCart = () => {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};
