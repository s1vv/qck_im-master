import axios from 'axios';

// Создаём экземпляр Axios
const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' }
});

// Функция для обновления токена
async function refreshToken() {
  try {
    const refresh = localStorage.getItem('refresh_token');
    if (!refresh) throw new Error('Refresh token отсутствует');

    const response = await axios.post('/api/users/refresh-token', { "refresh_token": refresh });

    const newAccessToken = response.data.access_token;
    localStorage.setItem('token', newAccessToken);
    api.defaults.headers['Authorization'] = `Bearer ${newAccessToken}`;

    return newAccessToken;
  } catch (error) {
    console.error('Ошибка обновления токена:', error);
    localStorage.removeItem('token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/login'; // Перенаправляем на логин
    throw error;
  }
}

// Добавляем интерцептор запроса
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
}, error => Promise.reject(error));

// Добавляем интерцептор ответа
api.interceptors.response.use(response => response, async error => {
  if (error.response?.status === 401) {
    try {
      const newAccessToken = await refreshToken();
      error.config.headers['Authorization'] = `Bearer ${newAccessToken}`;
      return api(error.config); // Повторяем запрос
    } catch (refreshError) {
      return Promise.reject(refreshError);
    }
  }
  return Promise.reject(error);
});

export default api;
