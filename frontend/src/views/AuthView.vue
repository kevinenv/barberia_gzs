<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const modoRegistro = ref(false)
const errorMsg = ref('')
const cargando = ref(false)

const form = ref({
  nombre: '',
  telefono: '',
  email: '',
  password: ''
})

async function procesarFormulario() {
  errorMsg.value = ''
  cargando.value = true
  try {
    if (modoRegistro.value) {
      await auth.registro(form.value)
    } else {
      await auth.login(form.value.email, form.value.password)
    }
    router.push('/catalogo')
  } catch (err) {
    errorMsg.value = err.response?.data?.error || 'Error al procesar la solicitud'
  } finally {
    cargando.value = false
  }
}
</script>

<template>
  <div class="auth-box card">
    <h2>{{ modoRegistro ? '📝 Crear Cuenta en Barbería GZS' : '🔑 Iniciar Sesión' }}</h2>
    <p class="muted">
      {{ modoRegistro ? 'Regístrate para llevar tu historial de cortes y reservar con seguridad.' : 'Ingresa con tu correo y contraseña.' }}
    </p>

    <div v-if="errorMsg" class="alert-error">{{ errorMsg }}</div>

    <form @submit.prevent="procesarFormulario" class="form-vertical">
      <div v-if="modoRegistro">
        <label>Nombre Completo</label>
        <input v-model="form.nombre" type="text" placeholder="Ej: Kevin Nuñez" required />
      </div>

      <div v-if="modoRegistro">
        <label>Teléfono / WhatsApp</label>
        <input v-model="form.telefono" type="text" placeholder="Ej: 70012345" required />
      </div>

      <div>
        <label>Correo Electrónico</label>
        <input v-model="form.email" type="email" placeholder="correo@ejemplo.com" required />
      </div>

      <div>
        <label>Contraseña</label>
        <input v-model="form.password" type="password" placeholder="Mínimo 6 caracteres" required />
      </div>

      <button type="submit" class="btn-gold" :disabled="cargando">
        {{ cargando ? 'Procesando...' : (modoRegistro ? 'Registrarme e Ingresar' : 'Entrar') }}
      </button>
    </form>

    <p class="switch-auth">
      {{ modoRegistro ? '¿Ya tienes una cuenta?' : '¿Eres cliente nuevo?' }}
      <a href="#" @click.prevent="modoRegistro = !modoRegistro">
        {{ modoRegistro ? 'Inicia sesión aquí' : 'Créate una cuenta gratis' }}
      </a>
    </p>
  </div>
</template>