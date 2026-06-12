import axios from 'axios';

const AUTH_URL = import.meta.env.VITE_AUTH_SERVICE_URL || 'http://localhost:8081';

export const authService = {
  register: async (userData) => {
    const response = await axios.post(`${AUTH_URL}/api/v1/auth/register`, userData);
    return response.data;
  },

  login: async (credentials) => {
    const response = await axios.post(`${AUTH_URL}/api/v1/auth/login`, credentials);
    return response.data;
  },

  logout: async (refreshToken) => {
    const response = await axios.post(`${AUTH_URL}/api/v1/auth/logout`, { refreshToken });
    return response.data;
  },

  refreshToken: async (refreshToken) => {
    const response = await axios.post(`${AUTH_URL}/api/v1/auth/refresh`, { refreshToken });
    return response.data;
  },

  getProfile: async (token) => {
    const response = await axios.get(`${AUTH_URL}/api/v1/auth/profile`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    return response.data;
  },

  updateProfile: async (token, userData) => {
    const response = await axios.put(`${AUTH_URL}/api/v1/auth/profile`, userData, {
      headers: { Authorization: `Bearer ${token}` }
    });
    return response.data;
  },
};
