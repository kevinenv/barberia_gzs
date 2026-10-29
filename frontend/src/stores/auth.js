import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import http from '../api/http'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('gzs_token') || null)
  const usuario = ref(JSON.parse(localStorage.getItem('gzs_usuario') || 'null'))

  const estaLogueado = computed(() => !!token.value)

  function guardarSesion(nuevoToken, datosUsuario) {
    token.value = nuevoToken
    usuario.value = datosUsuario
    localStorage.setItem('gzs_token', nuevoToken)
    localStorage.setItem('gzs_usuario', JSON.stringify(datosUsuario))
  }

  async function login(email, password) {
    const { data } = await http.post('/auth/login', { email, password })
    guardarSesion(data.token, data.usuario)
  }

  async function registro(payload) {
    const { data } = await http.post('/auth/registro', payload)
    guardarSesion(data.token, data.usuario)
  }

  function logout() {
    token.value = null
    usuario.value = null
    localStorage.removeItem('gzs_token')
    localStorage.removeItem('gzs_usuario')
  }

  return { token, usuario, estaLogueado, login, registro, logout }
})