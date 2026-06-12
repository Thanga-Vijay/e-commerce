import api from './axios';

const PRODUCT_URL = import.meta.env.VITE_PRODUCT_SERVICE_URL || 'http://localhost:8082';

export const productService = {
  getProducts: async (page = 1, limit = 12, search = '', category = '') => {
    const params = new URLSearchParams({ page, limit });
    if (search) params.append('search', search);
    if (category) params.append('category', category);
    
    const response = await api.get(`${PRODUCT_URL}/api/v1/products?${params}`);
    return response.data;
  },

  getProductById: async (id) => {
    const response = await api.get(`${PRODUCT_URL}/api/v1/products/${id}`);
    return response.data;
  },

  getCategories: async () => {
    const response = await api.get(`${PRODUCT_URL}/api/v1/categories`);
    return response.data;
  },

  createProduct: async (productData) => {
    const response = await api.post(`${PRODUCT_URL}/api/v1/products`, productData);
    return response.data;
  },

  updateProduct: async (id, productData) => {
    const response = await api.put(`${PRODUCT_URL}/api/v1/products/${id}`, productData);
    return response.data;
  },

  deleteProduct: async (id) => {
    const response = await api.delete(`${PRODUCT_URL}/api/v1/products/${id}`);
    return response.data;
  },
};
