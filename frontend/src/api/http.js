import axios from 'axios'

const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Interceptor profesional: Adjunta el Token JWT automáticamente si el usuario inició sesión
http.interceptors.request.use((config) => {
  const token = localStorage.getItem('gzs_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default http