import axios from 'axios';

const API_BASE_URL = '/api';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('nabd_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('nabd_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (token) => axios.post('/api/auth/login', { token }),
};

export const containerAPI = {
  getContainers: () => api.get('/containers'),
  getMetrics: () => api.get('/metrics'),
  getMetricsHistory: (id, hours = 24) => api.get(`/metrics/${id}/history?hours=${hours}`),
  getLogs: (container, lines = 100) => api.get(`/logs?container=${container}&lines=${lines}`),
  restartContainer: (name) => api.post(`/containers/${name}/restart`),
};

export const autoHealAPI = {
  getHistory: (limit = 50) => api.get(`/autoheal/history?limit=${limit}`),
  trigger: () => api.post('/autoheal/trigger'),
};

export const alertAPI = {
  getAlerts: () => api.get('/alerts'),
};

export default api;