<script setup>
import { ref, onMounted } from 'vue'
import http from '../api/http'

const solicitudes = ref([])
const marcaModelo = ref('')
const descripcionFalla = ref('')
const mensaje = ref('')
const error = ref('')

async function cargarSolicitudes() {
  const { data } = await http.get('/mantenimiento/mis-solicitudes')
  solicitudes.value = data
}

onMounted(cargarSolicitudes)

async function enviarSolicitud() {
  mensaje.value = ''
  error.value = ''
  try {
    const { data } = await http.post('/mantenimiento', {
      marca_modelo_silla: marcaModelo.value,
      descripcion_falla: descripcionFalla.value
    })
    mensaje.value = data.mensaje
    marcaModelo.value = ''
    descripcionFalla.value = ''
    await cargarSolicitudes()
  } catch (err) {
    error.value = err.response?.data?.error || 'Error al enviar la solicitud'
  }
}
</script>

<template>
  <div>
    <h2 class="section-title">🔧 Servicio Técnico: Mantenimiento de Sillas Hidráulicas</h2>
    <p class="muted">Describe el problema del pistón, bomba hidráulica o tapizado de tu sillón de barbería para recibir cotización y fecha de reparación.</p>

    <div v-if="mensaje" class="alert-success">{{ mensaje }}</div>
    <div v-if="error" class="alert-error">{{ error }}</div>

    <div class="booking-layout">
      <div class="card">
        <h3>Solicitar Diagnóstico y Cotización</h3>
        <form @submit.prevent="enviarSolicitud" class="form-vertical">
          <div>
            <label>Marca / Modelo del Sillón Hidráulico</label>
            <input v-model="marcaModelo" type="text" placeholder="Ej: Sillón Hidráulico Clásico Cromado" required />
          </div>
          <div>
            <label>Descripción detallada de la falla</label>
            <textarea v-model="descripcionFalla" rows="4" placeholder="Ej: El pedal sube pero la silla se baja sola con el peso del cliente, pierde aceite en la base..." required></textarea>
          </div>
          <button type="submit" class="btn-gold">📨 Enviar Solicitud de Mantenimiento</button>
        </form>
      </div>

      <div class="card">
        <h3>Mis Solicitudes Enviadas</h3>
        <p v-if="solicitudes.length === 0" class="muted">No has enviado solicitudes de mantenimiento aún.</p>
        <div v-for="sol in solicitudes" :key="sol.id" class="maint-item">
          <div style="display: flex; justify-content: space-between;">
            <strong>{{ sol.marca_modelo_silla }}</strong>
            <span class="badge-sillon">{{ sol.estado }}</span>
          </div>
          <p class="muted" style="margin: 0.4rem 0;">{{ sol.descripcion_falla }}</p>
          <small>
            Cotización: <strong>{{ sol.cotizacion_estimada ? `Bs. ${sol.cotizacion_estimada}` : 'Pendiente de revisión' }}</strong> |
            Fecha: <strong>{{ sol.fecha_programada || 'Por definir' }}</strong>
          </small>
        </div>
      </div>
    </div>
  </div>
</template>