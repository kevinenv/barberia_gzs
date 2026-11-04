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
  <div class="monument-auth">
    <div class="auth-center-col">
      <header class="auth-header">
        <h1 class="auth-title">
          {{ modoRegistro ? 'NUEVA CUENTA' : 'ACCESO AL SALÓN' }}
        </h1>
        <p class="auth-sub">
          {{ modoRegistro
            ? 'Complete sus datos personales para agendar turnos y revisar su historial en el salón.'
            : 'Ingrese con sus credenciales autorizadas de cliente o personal de Barbería GZS.' }}
        </p>
      </header>

      <div v-if="errorMsg" class="auth-alert">{{ errorMsg }}</div>

      <form @submit.prevent="procesarFormulario" class="auth-form-stack">
        <div v-if="modoRegistro" class="auth-field-unit">
          <label class="field-label">Nombre y Apellido</label>
          <input
            v-model="form.nombre"
            type="text"
            placeholder="Ej: Kevin Nuñez"
            class="field-input"
            required
          />
        </div>

        <div v-if="modoRegistro" class="auth-field-unit">
          <label class="field-label">Teléfono / WhatsApp</label>
          <input
            v-model="form.telefono"
            type="text"
            placeholder="Ej: 70012345"
            class="field-input"
            required
          />
        </div>

        <div class="auth-field-unit">
          <label class="field-label">Correo Electrónico</label>
          <input
            v-model="form.email"
            type="email"
            placeholder="correo@ejemplo.com"
            class="field-input"
            required
          />
        </div>

        <div class="auth-field-unit">
          <label class="field-label">Contraseña de Seguridad</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="Mínimo 6 caracteres"
            class="field-input"
            required
          />
        </div>

        <button
          type="submit"
          class="btn-auth-submit"
          :disabled="cargando"
        >
          {{ cargando ? 'Verificando...' : (modoRegistro ? 'Confirmar Registro' : 'Ingresar al Sistema') }}
        </button>

        <div class="auth-mode-switch">
          <button
            type="button"
            class="switch-text-btn"
            @click="modoRegistro = !modoRegistro; errorMsg = ''"
          >
            {{ modoRegistro
              ? '¿Ya posee una cuenta? Iniciar Sesión'
              : '¿Primera vez en el salón? Crear Cuenta' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.monument-auth {
  max-width: 1400px;
  margin: 0 auto;
  padding: 5rem 2rem 8rem;
  display: flex;
  justify-content: center;
}

.auth-center-col {
  width: 100%;
  max-width: 520px;
  background-color: var(--bg-surface);
  padding: 3.5rem 3rem;
}

.auth-header {
  margin-bottom: 2.5rem;
}

.auth-title {
  font-size: var(--step-head);
  font-weight: 900;
  letter-spacing: -0.05em;
  text-transform: uppercase;
  color: var(--text-ink);
  line-height: 1.05;
  margin-bottom: 0.75rem;
}

.auth-sub {
  font-size: var(--step-small);
  color: var(--text-muted);
  line-height: 1.5;
}

.auth-alert {
  padding: 1rem 1.25rem;
  background-color: #3d1216;
  color: #ffb3b8;
  font-size: var(--step-small);
  font-weight: 700;
  margin-bottom: 2rem;
}

.auth-form-stack {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.auth-field-unit {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field-label {
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-ink);
}

.field-input {
  padding: 1rem 1.25rem;
  background-color: var(--bg-ground);
  color: var(--text-ink);
  border: none;
  font-size: var(--step-body);
  outline: none;
}

.field-input:focus {
  background-color: var(--bg-surface-elevated);
}

.btn-auth-submit {
  width: 100%;
  padding: 1.25rem;
  background-color: var(--accent-crimson);
  color: var(--accent-ink);
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  text-align: center;
  margin-top: 1rem;
  transition: opacity 140ms ease;
}

.btn-auth-submit:hover:not(:disabled) {
  opacity: 0.88;
}

.btn-auth-submit:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.auth-mode-switch {
  text-align: center;
  padding-top: 1.5rem;
}

.switch-text-btn {
  font-size: var(--step-small);
  font-weight: 700;
  color: var(--text-muted);
  transition: color 140ms ease;
}

.switch-text-btn:hover {
  color: var(--text-ink);
}

@media (max-width: 600px) {
  .auth-center-col {
    padding: 2.5rem 1.5rem;
  }
}
</style>
