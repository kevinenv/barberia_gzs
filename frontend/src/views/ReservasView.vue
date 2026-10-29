<script setup>
import { ref, computed, onMounted } from 'vue'
import http from '../api/http'

const servicios = ref([])
const barberos = ref([])
const misCitas = ref([])

const serviciosSeleccionados = ref([])
const barberoId = ref(0) // 0 = Cualquier barbero disponible (Modelo Híbrido)
const fecha = ref(new Date().toISOString().split('T')[0])
const horaInicio = ref('15:00')
const metodoPago = ref('qr')

const mensajeExito = ref('')
const mensajeError = ref('')
const cargando = ref(false)

async function cargarDatos() {
  const [resServ, resBarb, resCitas] = await Promise.all([
    http.get('/servicios'),
    http.get('/barberos'),
    http.get('/citas/mis-citas')
  ])
  servicios.value = resServ.data
  barberos.value = resBarb.data
  misCitas.value = resCitas.data
}

onMounted(cargarDatos)

function toggleServicio(id) {
  const index = serviciosSeleccionados.value.indexOf(id)
  if (index === -1) {
    serviciosSeleccionados.value.push(id)
  } else {
    serviciosSeleccionados.value.splice(index, 1)
  }
}

// Cálculos automáticos en tiempo real del carrito de servicios
const resumenCarrito = computed(() => {
  const elegidos = servicios.value.filter(s => serviciosSeleccionados.value.includes(s.id))
  const totalBs = elegidos.reduce((acc, s) => acc + s.precio, 0)
  const totalMinutos = elegidos.reduce((acc, s) => acc + s.duracion_minutos, 0)
  const garantia50 = totalBs * 0.5
  return { totalBs, totalMinutos, garantia50 }
})

async function confirmarReserva() {
  mensajeExito.value = ''
  mensajeError.value = ''

  if (serviciosSeleccionados.value.length === 0) {
    mensajeError.value = 'Selecciona al menos un servicio para tu cita.'
    return
  }

  cargando.value = true
  try {
    const { data } = await http.post('/citas', {
      barbero_id: Number(barberoId.value),
      fecha: fecha.value,
      hora_inicio: horaInicio.value,
      servicios_ids: serviciosSeleccionados.value,
      metodo_pago_garantia: metodoPago.value
    })
    mensajeExito.value = `${data.mensaje} (Horario asignado: ${data.cita.hora_inicio} a ${data.cita.hora_fin})`
    serviciosSeleccionados.value = []
    await cargarDatos()
  } catch (err) {
    mensajeError.value = err.response?.data?.error || 'No se pudo registrar la cita'
  } finally {
    cargando.value = false
  }
}
</script>

<template>
  <div>
    <h2 class="section-title">📅 Agendar Cita (Carrito de Servicios)</h2>

    <div v-if="mensajeExito" class="alert-success">{{ mensajeExito }}</div>
    <div v-if="mensajeError" class="alert-error">{{ mensajeError }}</div>

    <div class="booking-layout">
      <!-- Columna Izquierda: Selección de Servicios, Fecha y Barbero -->
      <div class="card">
        <h3>1. Elige uno o varios servicios</h3>
        <div class="service-selector">
          <div
            v-for="s in servicios"
            :key="s.id"
            class="service-item"
            :class="{ selected: serviciosSeleccionados.includes(s.id) }"
            @click="toggleServicio(s.id)"
          >
            <div>
              <strong>{{ s.nombre }}</strong>
              <div class="muted">{{ s.duracion_minutos }} min</div>
            </div>
            <div class="price-tag" style="margin: 0;">Bs. {{ s.precio }}</div>
          </div>
        </div>

        <h3 style="margin-top: 1.5rem;">2. Elige Barbero y Horario</h3>
        <div class="form-vertical">
          <div>
            <label>Especialista / Sillón</label>
            <select v-model="barberoId">
              <option :value="0">⚡ Cualquier barbero disponible (Asignación Automática)</option>
              <option v-for="b in barberos" :key="b.id" :value="b.id">
                Sillón #{{ b.silla_numero }} - {{ b.apodo }} ({{ b.especialidad }})
              </option>
            </select>
          </div>

          <div class="row-2">
            <div>
              <label>Fecha de la Cita</label>
              <input type="date" v-model="fecha" required />
            </div>
            <div>
              <label>Hora de Inicio (09:00 a 20:00)</label>
              <input type="time" v-model="horaInicio" required />
            </div>
          </div>

          <div>
            <label>Método de Pago de Garantía (50%)</label>
            <select v-model="metodoPago">
              <option value="qr">📱 Pago con QR Simple</option>
              <option value="transferencia">🏦 Transferencia Bancaria</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Columna Derecha: Resumen y Garantía QR -->
      <div class="card summary-box">
        <h3>🧾 Resumen y Garantía Anti-Faltas</h3>
        <p class="muted">Para asegurar tu sillón y evitar citas falsas, se abona el 50% de adelanto.</p>

        <div class="summary-line">
          <span>Duración total estimada:</span>
          <strong>{{ resumenCarrito.totalMinutos }} minutos</strong>
        </div>
        <div class="summary-line">
          <span>Costo total de servicios:</span>
          <strong>Bs. {{ resumenCarrito.totalBs.toFixed(2) }}</strong>
        </div>
        <div class="summary-line highlight-gold">
          <span>Adelanto obligatorio (50%):</span>
          <strong>Bs. {{ resumenCarrito.garantia50.toFixed(2) }}</strong>
        </div>

        <div class="qr-placeholder">
          <span>[ QR BARBERÍA GZS - Bs. {{ resumenCarrito.garantia50.toFixed(2) }} ]</span>
        </div>

        <button class="btn-gold" style="width: 100%;" :disabled="cargando" @click="confirmarReserva">
          {{ cargando ? 'Verificando disponibilidad...' : '✅ Confirmar Pago QR y Agendar' }}
        </button>
      </div>
    </div>

    <!-- Historial de Citas del Cliente -->
    <h2 class="section-title">🕒 Mi Historial de Citas</h2>
    <div class="card">
      <p v-if="misCitas.length === 0" class="muted">Aún no tienes citas registradas.</p>
      <table v-else class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Fecha</th>
            <th>Horario</th>
            <th>Barbero / Sillón</th>
            <th>Total</th>
            <th>Garantía Pagada</th>
            <th>Estado</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in misCitas" :key="c.id">
            <td>#{{ c.id }}</td>
            <td>{{ c.fecha }}</td>
            <td>{{ c.hora_inicio }} - {{ c.hora_fin }}</td>
            <td>Sillón #{{ c.silla_numero }} ({{ c.barbero_apodo }})</td>
            <td>Bs. {{ c.monto_total }}</td>
            <td>Bs. {{ c.monto_garantia_pagado }} ({{ c.metodo_pago_garantia.toUpperCase() }})</td>
            <td><span class="badge-sillon">{{ c.estado.toUpperCase() }}</span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>