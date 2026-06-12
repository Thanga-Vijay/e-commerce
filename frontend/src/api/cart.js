import api from './axios';

const CART_URL = import.meta.env.VITE_CART_SERVICE_URL || 'http://localhost:8083';

export const cartService = {
  getCart: async () => {
    const response = await api.get(`${CART_URL}/api/v1/cart`);
    return response.data;
  },

  addItem: async (productId, quantity) => {
    const response = await api.post(`${CART_URL}/api/v1/cart/items`, {
      productId,
      quantity,
    });
    return response.data;
  },

  updateItem: async (itemId, quantity) => {
    const response = await api.put(`${CART_URL}/api/v1/cart/items/${itemId}`, {
      quantity,
    });
    return response.data;
  },

  removeItem: async (itemId) => {
    const response = await api.delete(`${CART_URL}/api/v1/cart/items/${itemId}`);
    return response.data;
  },

  clearCart: async () => {
    const response = await api.delete(`${CART_URL}/api/v1/cart`);
    return response.data;
  },
};
