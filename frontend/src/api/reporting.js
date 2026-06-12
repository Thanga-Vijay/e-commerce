import api from './axios';

const REPORTING_URL = import.meta.env.VITE_REPORTING_SERVICE_URL || 'http://localhost:8089';

export const reportingService = {
  getDashboard: async () => {
    const response = await api.get(`${REPORTING_URL}/api/v1/dashboard`);
    return response.data;
  },

  exportDashboard: async (format = 'csv') => {
    const response = await api.get(`${REPORTING_URL}/api/v1/dashboard/export?format=${format}`, {
      responseType: 'blob',
    });
    return response.data;
  },

  getRevenueReport: async (startDate, endDate, period = 'daily') => {
    const response = await api.get(
      `${REPORTING_URL}/api/v1/reports/revenue?startDate=${startDate}&endDate=${endDate}&period=${period}`
    );
    return response.data;
  },

  exportRevenueReport: async (startDate, endDate, period, format = 'csv') => {
    const response = await api.get(
      `${REPORTING_URL}/api/v1/reports/revenue/export?startDate=${startDate}&endDate=${endDate}&period=${period}&format=${format}`,
      { responseType: 'blob' }
    );
    return response.data;
  },

  getTopProducts: async (limit = 10) => {
    const response = await api.get(`${REPORTING_URL}/api/v1/reports/products?limit=${limit}`);
    return response.data;
  },

  getCustomerReport: async (limit = 10) => {
    const response = await api.get(`${REPORTING_URL}/api/v1/reports/customers?limit=${limit}`);
    return response.data;
  },
};
