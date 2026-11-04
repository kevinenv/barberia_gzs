<script setup>
import { ref, computed, onMounted } from 'vue'
import http from '../api/http'

const fallbackServicios = [
  { id: 1, nombre: 'Corte Clásico de Caballero', descripcion: 'Diagnóstico capilar, corte a tijera o máquina, lavado previo y toalla caliente aromática.', precio: 50, duracion_minutos: 45 },
  { id: 2, nombre: 'Perfilado de Barba & Navaja', descripcion: 'Tratamiento de vaporización, aceites emolientes y afeitado tradicional a filo recto.', precio: 40, duracion_minutos: 35 },
  { id: 3, nombre: 'Servicio Completo Monumental', descripcion: 'Corte de autor, perfilado integral con navaja, toalla caliente y masaje capilar.', precio: 80, duracion_minutos: 75 },
  { id: 4, nombre: 'Afeitado Clásico Completo', descripcion: 'Espuma batida en tazón, triple pasada de navaja y loción balsámica calmante.', precio: 45, duracion_minutos: 40 }
]

const fallbackBarberos = [
  { id: 1, apodo: 'Gonzalo', silla_numero: 1, especialidad: 'Especialista en degradados limpios y tijera clásica' },
  { id: 2, apodo: 'Javier', silla_numero: 2, especialidad: 'Maestro de afeitado con navaja recta y barbas largas' },
  { id: 3, apodo: 'Mateo', silla_numero: 3, especialidad: 'Estilismo contemporáneo y perfilado tradicional' }
]

const servicios = ref(fallbackServicios)
const barberos = ref(fallbackBarberos)
const misCitas = ref([])

const serviciosSeleccionados = ref([])
const barberoId = ref(0)
const fecha = ref(new Date().toISOString().split('T')[0])
const horaInicio = ref('15:00')
const metodoPago = ref('qr')

const mensajeExito = ref('')
const mensajeError = ref('')
const cargando = ref(false)

async function cargarDatos() {
  try {
    const [resServ, resBarb, resCitas] = await Promise.all([
      http.get('/servicios'),
      http.get('/barberos'),
      http.get('/citas/mis-citas')
    ])
    servicios.value = resServ.data
    barberos.value = resBarb.data
    misCitas.value = resCitas.data
  } catch (error) {
    console.error('Error al cargar datos de reservas:', error)
  }
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

const resumenCarrito = computed(() => {
  const elegidos = servicios.value.filter(s => serviciosSeleccionados.value.includes(s.id))
  const totalBs = elegidos.reduce((acc, s) => acc + s.precio, 0)
  const totalMinutos = elegidos.reduce((acc, s) => acc + s.duracion_minutos, 0)
  const garantia30 = Math.round(totalBs * 0.3)
  return { totalBs, totalMinutos, garantia30 }
})

async function confirmarReserva() {
  mensajeExito.value = ''
  mensajeError.value = ''

  if (serviciosSeleccionados.value.length === 0) {
    mensajeError.value = 'Seleccione al menos un servicio para su cita.'
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
    mensajeExito.value = `${data.mensaje}. Horario asignado: ${data.cita.hora_inicio} a ${data.cita.hora_fin}`
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
  <div class="monument-booking">
    <header class="booking-monument-head">
      <div class="head-lead">
        <span class="booking-category">Sistema de Turnos</span>
        <h1 class="booking-headline">RESERVAS &bull; GZS</h1>
      </div>
      <p class="booking-subcopy">
        Seleccione sus servicios, elija su sillón y asegure su turno sin filas de espera. Atención personalizada con toalla caliente y lavado de cortesía.
      </p>
    </header>

    <div v-if="mensajeExito" class="booking-banner banner-success">
      {{ mensajeExito }}
    </div>
    <div v-if="mensajeError" class="booking-banner banner-error">
      {{ mensajeError }}
    </div>

    <div class="booking-layout-grid">
      <!-- FORM STACK -->
      <div class="booking-steps-column">
        <!-- PASO 1 -->
        <section class="booking-card-step">
          <div class="step-head-bar">
            <span class="step-tag">Paso Primero</span>
            <h2 class="step-h2">Selección de Cortes y Servicios</h2>
          </div>
          <p class="step-desc">Haga clic en uno o varios servicios para agregarlos al turno.</p>

          <div class="services-selector">
            <button
              v-for="s in servicios"
              :key="s.id"
              type="button"
              class="service-select-btn"
              :class="{ 'is-selected': serviciosSeleccionados.includes(s.id) }"
              @click="toggleServicio(s.id)"
            >
              <div class="btn-info">
                <span class="btn-name">{{ s.nombre }}</span>
                <span class="btn-duration">{{ s.duracion_minutos }} minutos</span>
              </div>
              <span class="btn-cost">Bs. {{ s.precio }}</span>
            </button>
          </div>
        </section>

        <!-- PASO 2 -->
        <section class="booking-card-step">
          <div class="step-head-bar">
            <span class="step-tag">Paso Segundo</span>
            <h2 class="step-h2">Estación, Barbero y Horario</h2>
          </div>
          <p class="step-desc">Configure el horario exacto de su llegada y el especialista asignado.</p>

          <div class="inputs-stack">
            <div class="input-unit">
              <label class="unit-label">Especialista / Sillón</label>
              <select v-model="barberoId" class="unit-field">
                <option :value="0">Cualquier sillón disponible (Asignación automática)</option>
                <option v-for="b in barberos" :key="b.id" :value="b.id">
                  Sillón #{{ b.silla_numero }} &mdash; {{ b.apodo }} ({{ b.especialidad }})
                </option>
              </select>
            </div>

            <div class="inputs-duo">
              <div class="input-unit">
                <label class="unit-label">Fecha del Turno</label>
                <input type="date" v-model="fecha" class="unit-field" />
              </div>

              <div class="input-unit">
                <label class="unit-label">Hora de Inicio</label>
                <input type="time" v-model="horaInicio" step="1800" class="unit-field" />
              </div>
            </div>

            <div class="input-unit">
              <label class="unit-label">Modalidad de Garantía (30%)</label>
              <div class="payment-tabs">
                <button
                  type="button"
                  class="payment-tab"
                  :class="{ 'payment-tab--active': metodoPago === 'qr' }"
                  @click="metodoPago = 'qr'"
                >
                  Transferencia QR Simple
                </button>
                <button
                  type="button"
                  class="payment-tab"
                  :class="{ 'payment-tab--active': metodoPago === 'efectivo' }"
                  @click="metodoPago = 'efectivo'"
                >
                  Abono en Mostrador
                </button>
              </div>
            </div>
          </div>
        </section>
      </div>

      <!-- SIDE TICKET SUMMARY -->
      <aside class="booking-ticket-sidebar">
        <div class="ticket-monolith">
          <div class="ticket-header">
            <h3 class="ticket-title">Comprobante de Cita</h3>
            <span class="ticket-id">Turno GZS-{{ fecha.replace(/-/g, '') }}</span>
          </div>

          <div class="ticket-body">
            <div class="ticket-line">
              <span class="ticket-label">Servicios elegidos:</span>
              <span class="ticket-val">{{ serviciosSeleccionados.length }} servicio(s)</span>
            </div>

            <div class="ticket-line">
              <span class="ticket-label">Tiempo estimado:</span>
              <span class="ticket-val">{{ resumenCarrito.totalMinutos }} minutos</span>
            </div>

            <div class="ticket-line">
              <span class="ticket-label">Total en Salón:</span>
              <span class="ticket-val-accent">Bs. {{ resumenCarrito.totalBs }}</span>
            </div>

            <div class="ticket-anticipo-box">
              <span class="anticipo-label">Anticipo QR Requerido (30%)</span>
              <span class="anticipo-amount">Bs. {{ resumenCarrito.garantia30 }}</span>
              <p class="anticipo-note">
                El anticipo garantiza la disponibilidad de su sillón y se descuenta del total.
              </p>
            </div>

            <div v-if="metodoPago === 'qr'" class="qr-deposit-box">
              <p class="qr-bank-info">
                Banco Unión &bull; Cta: 100000458921 &bull; Barbería GZS S.R.L.
              </p>
              <div class="qr-visual-stub">
                <span class="qr-mock-code">[ CÓDIGO QR SIMPLE DEPOSITO BS. {{ resumenCarrito.garantia30 }} ]</span>
              </div>
              <a
                href="https://wa.me/59170000000?text=Hola%20Barbería%20GZS,%20adjunto%20comprobante%20de%20anticipo%20QR"
                target="_blank"
                class="btn-whatsapp-proof"
              >
                Enviar Comprobante WhatsApp
              </a>
            </div>

            <button
              class="btn-confirm-booking"
              :disabled="cargando || serviciosSeleccionados.length === 0"
              @click="confirmarReserva"
            >
              <span v-if="cargando">Procesando Turno...</span>
              <span v-else>Confirmar y Asegurar Cita</span>
            </button>
          </div>
        </div>

        <!-- HISTORIAL DE CITAS DEL CLIENTE -->
        <div class="appointments-history-box" v-if="misCitas.length > 0">
          <h4 class="history-title">Mis Turnos Agendados</h4>
          <div class="history-list">
            <div v-for="cita in misCitas" :key="cita.id" class="history-entry">
              <div class="entry-top">
                <span class="entry-date">{{ cita.fecha }} &bull; {{ cita.hora_inicio }}</span>
                <span class="entry-status" :class="'status-' + cita.estado.toLowerCase()">
                  {{ cita.estado }}
                </span>
              </div>
              <span class="entry-barber">Barbero: {{ cita.barbero_apodo || 'Por Asignar' }}</span>
              <span class="entry-total">Bs. {{ cita.total }}</span>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.monument-booking {
  max-width: 1400px;
  margin: 0 auto;
  padding: 3rem 2rem 6rem;
}

.booking-monument-head {
  margin-bottom: 4rem;
}

.booking-category {
  display: block;
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--accent-crimson);
  margin-bottom: 0.5rem;
}

.booking-headline {
  font-size: var(--step-head);
  font-weight: 900;
  letter-spacing: -0.05em;
  text-transform: uppercase;
  line-height: 0.95;
  color: var(--text-ink);
  margin-bottom: 1rem;
}

.booking-subcopy {
  font-size: var(--step-body);
  color: var(--text-muted);
  max-width: 700px;
  line-height: 1.5;
}

.booking-banner {
  padding: 1.25rem 2rem;
  margin-bottom: 2.5rem;
  font-size: var(--step-small);
  font-weight: 700;
}

.banner-success {
  background-color: #1b382b;
  color: #a8f5c8;
}

.banner-error {
  background-color: #3d1216;
  color: #ffb3b8;
}

.booking-layout-grid {
  display: grid;
  grid-template-columns: 1.6fr 1fr;
  gap: 3.5rem;
  align-items: start;
}

.booking-steps-column {
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

.booking-card-step {
  background-color: var(--bg-surface);
  padding: 2.5rem;
}

.step-head-bar {
  margin-bottom: 0.5rem;
}

.step-tag {
  display: block;
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--accent-crimson);
  margin-bottom: 0.25rem;
}

.step-h2 {
  font-size: var(--step-lead);
  font-weight: 800;
  color: var(--text-ink);
  letter-spacing: -0.03em;
}

.step-desc {
  font-size: var(--step-small);
  color: var(--text-muted);
  margin-bottom: 2rem;
}

.services-selector {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.service-select-btn {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  background-color: var(--bg-ground);
  color: var(--text-ink);
  text-align: left;
  transition: background-color 140ms ease, color 140ms ease;
}

.service-select-btn:hover {
  background-color: var(--bg-surface-elevated);
}

.service-select-btn.is-selected {
  background-color: var(--accent-crimson);
  color: var(--accent-ink);
}

.service-select-btn.is-selected .btn-duration {
  color: #ffffff;
}

.btn-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.btn-name {
  font-size: var(--step-body);
  font-weight: 700;
}

.btn-duration {
  font-size: var(--step-small);
  color: var(--text-muted);
}

.btn-cost {
  font-size: var(--step-lead);
  font-weight: 800;
  font-variant-numeric: tabular-nums;
}

.inputs-stack {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.inputs-duo {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.input-unit {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.unit-label {
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-ink);
}

.unit-field {
  padding: 1rem 1.25rem;
  background-color: var(--bg-ground);
  color: var(--text-ink);
  border: none;
  font-size: var(--step-body);
  outline: none;
}

.unit-field:focus {
  background-color: var(--bg-surface-elevated);
}

.payment-tabs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
}

.payment-tab {
  padding: 1rem;
  background-color: var(--bg-ground);
  color: var(--text-muted);
  font-size: var(--step-small);
  font-weight: 700;
  text-align: center;
  transition: background-color 140ms ease, color 140ms ease;
}

.payment-tab:hover {
  color: var(--text-ink);
}

.payment-tab--active {
  background-color: var(--bg-surface-elevated);
  color: var(--text-ink);
  font-weight: 800;
}

/* SIDEBAR TICKET */
.booking-ticket-sidebar {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.ticket-monolith {
  background-color: var(--bg-surface);
  padding: 2.5rem;
}

.ticket-header {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  background-color: var(--bg-surface);
}

.ticket-kicker {
  display: block;
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--accent-crimson);
  margin-bottom: 0.25rem;
}

.ticket-title {
  font-size: var(--step-lead);
  font-weight: 800;
  color: var(--text-ink);
  margin-bottom: 0.25rem;
}

.ticket-id {
  font-size: var(--step-small);
  color: var(--text-muted);
}

.ticket-body {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.ticket-line {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.ticket-label {
  font-size: var(--step-small);
  color: var(--text-muted);
}

.ticket-val {
  font-size: var(--step-small);
  font-weight: 700;
  color: var(--text-ink);
}

.ticket-val-accent {
  font-size: var(--step-lead);
  font-weight: 800;
  color: var(--text-ink);
}

.ticket-anticipo-box {
  background-color: var(--bg-ground);
  padding: 1.5rem;
  margin: 1rem 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.anticipo-label {
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--accent-crimson);
}

.anticipo-amount {
  font-size: var(--step-lead);
  font-weight: 900;
  color: var(--text-ink);
}

.anticipo-note {
  font-size: var(--step-small);
  color: var(--text-muted);
  line-height: 1.4;
}

.qr-deposit-box {
  background-color: var(--bg-ground);
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.qr-bank-info {
  font-size: var(--step-small);
  color: var(--text-muted);
  line-height: 1.4;
}

.qr-visual-stub {
  background-color: var(--bg-surface);
  padding: 1.5rem;
  text-align: center;
}

.qr-mock-code {
  font-size: var(--step-small);
  font-weight: 700;
  letter-spacing: 0.05em;
  color: var(--text-ink);
}

.btn-whatsapp-proof {
  display: block;
  text-align: center;
  padding: 0.75rem 1rem;
  background-color: #123d24;
  color: #a8f5c8;
  font-size: var(--step-small);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  transition: opacity 140ms ease;
}

.btn-whatsapp-proof:hover {
  opacity: 0.9;
}

.btn-confirm-booking {
  width: 100%;
  padding: 1.25rem;
  background-color: var(--accent-crimson);
  color: var(--accent-ink);
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  text-align: center;
  transition: opacity 140ms ease;
}

.btn-confirm-booking:hover:not(:disabled) {
  opacity: 0.88;
}

.btn-confirm-booking:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* HISTORIAL */
.appointments-history-box {
  background-color: var(--bg-surface);
  padding: 2rem;
}

.history-title {
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin-bottom: 1.5rem;
  color: var(--text-ink);
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.history-entry {
  background-color: var(--bg-ground);
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.entry-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.entry-date {
  font-size: var(--step-small);
  font-weight: 700;
  color: var(--text-ink);
}

.entry-status {
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.status-pendiente {
  color: #f7b731;
}

.status-confirmado {
  color: #2ed573;
}

.status-completado {
  color: #a4b0be;
}

.entry-barber {
  font-size: var(--step-small);
  color: var(--text-muted);
}

.entry-total {
  font-size: var(--step-small);
  font-weight: 800;
  color: var(--text-ink);
}

@media (max-width: 900px) {
  .booking-layout-grid {
    grid-template-columns: 1fr;
  }
  .inputs-duo {
    grid-template-columns: 1fr;
  }
}
</style>
