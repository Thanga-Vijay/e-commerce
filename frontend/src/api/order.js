import api from './axios';

const ORDER_URL = import.meta.env.VITE_ORDER_SERVICE_URL || 'http://localhost:8085';

export const orderService = {
  createOrder: async (orderData) => {
    const response = await api.post(`${ORDER_URL}/api/v1/orders`, orderData);
    return response.data;
  },

  getOrders: async (page = 1, limit = 10) => {
    const response = await api.get(`${ORDER_URL}/api/v1/orders?page=${page}&limit=${limit}`);
    return response.data;
  },

  getOrderById: async (id) => {
    const response = await api.get(`${ORDER_URL}/api/v1/orders/${id}`);
    return response.data;
  },

  trackOrder: async (id) => {
    const response = await api.get(`${ORDER_URL}/api/v1/orders/${id}/track`);
    return response.data;
  },

  cancelOrder: async (id) => {
    const response = await api.put(`${ORDER_URL}/api/v1/orders/${id}/cancel`);
    return response.data;
  },

  updateOrderStatus: async (id, status) => {
    const response = await api.put(`${ORDER_URL}/api/v1/orders/${id}/status`, { status });
    return response.data;
  },
};
