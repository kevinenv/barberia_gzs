<script setup>
import { ref, onMounted } from 'vue'
import http from '../api/http'

const solicitudes = ref([])
const marcaModelo = ref('')
const descripcionFalla = ref('')
const mensaje = ref('')
const error = ref('')

async function cargarSolicitudes() {
  try {
    const { data } = await http.get('/mantenimiento/mis-solicitudes')
    solicitudes.value = data
  } catch (err) {
    console.error('Error al cargar solicitudes:', err)
  }
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
  <div class="monument-maint">
    <header class="maint-header">
      <span class="maint-category">Ingeniería del Salón</span>
      <h1 class="maint-title">MANTENIMIENTO &bull; SILLONES</h1>
      <p class="maint-sub">
        Servicio técnico y calibración de sillones hidráulicos de barbería. Diagnóstico de pistones, bombas de aceite y ajuste mecánico para barberías asociadas.
      </p>
    </header>

    <div v-if="mensaje" class="maint-alert maint-alert--success">{{ mensaje }}</div>
    <div v-if="error" class="maint-alert maint-alert--error">{{ error }}</div>

    <div class="maint-layout-grid">
      <!-- FORMULARIO -->
      <section class="maint-form-card">
        <h2 class="card-title">Solicitud de Diagnóstico Técnico</h2>
        <p class="card-desc">Detalle la anomalía detectada en el mecanismo para coordinar la inspección técnica en taller.</p>

        <form @submit.prevent="enviarSolicitud" class="maint-inputs-stack">
          <div class="maint-unit">
            <label class="unit-label" for="marca-silla">Marca o Modelo del Sillón</label>
            <input
              id="marca-silla"
              v-model="marcaModelo"
              type="text"
              placeholder="Ej: Belmont Clásico Hidráulico Cromado"
              class="unit-field"
              required
            />
          </div>

          <div class="maint-unit">
            <label class="unit-label" for="falla-desc">Descripción de la Falla Técnica</label>
            <textarea
              id="falla-desc"
              v-model="descripcionFalla"
              rows="5"
              placeholder="Ej: Pérdida gradual de presión al sentar al cliente, fuga de fluido hidráulico en la base giratoria..."
              class="unit-field unit-textarea"
              required
            ></textarea>
          </div>

          <button type="submit" class="btn-maint-submit">
            Enviar Solicitud a Taller
          </button>
        </form>
      </section>

      <!-- HISTORIAL / ORDENES -->
      <section class="maint-history-card">
        <h2 class="card-title">Órdenes de Reparación</h2>
        <p class="card-desc">Estado de diagnósticos e informes técnicos emitidos por los mecánicos.</p>

        <p v-if="solicitudes.length === 0" class="blank-orders-notice">
          No existen órdenes de mantenimiento activas en el registro técnico.
        </p>

        <div class="orders-stack">
          <article
            v-for="sol in solicitudes"
            :key="sol.id"
            class="order-item"
          >
            <div class="order-top">
              <span class="order-equipment">{{ sol.marca_modelo_silla }}</span>
              <span class="order-badge" :class="'badge-' + sol.estado.toLowerCase()">
                {{ sol.estado }}
              </span>
            </div>

            <p class="order-defect">{{ sol.descripcion_falla }}</p>

            <div class="order-data-grid">
              <div class="order-data-unit">
                <span class="data-tag">Presupuesto</span>
                <span class="data-value">
                  {{ sol.costo_reparacion ? `Bs. ${sol.costo_reparacion}` : 'Pendiente' }}
                </span>
              </div>

              <div class="order-data-unit">
                <span class="data-tag">Fecha Registro</span>
                <span class="data-value">
                  {{ new Date(sol.created_at).toLocaleDateString() }}
                </span>
              </div>
            </div>

            <div v-if="sol.informe_tecnico" class="order-technical-report">
              <span class="report-tag">Informe Técnico de Taller:</span>
              <p class="report-text">{{ sol.informe_tecnico }}</p>
            </div>
          </article>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.monument-maint {
  max-width: 1400px;
  margin: 0 auto;
  padding: 3rem 2rem 6rem;
}

.maint-header {
  margin-bottom: 4rem;
}

.maint-category {
  display: block;
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--accent-crimson);
  margin-bottom: 0.5rem;
}

.maint-title {
  font-size: var(--step-head);
  font-weight: 900;
  letter-spacing: -0.05em;
  text-transform: uppercase;
  line-height: 0.95;
  color: var(--text-ink);
  margin-bottom: 1rem;
}

.maint-sub {
  font-size: var(--step-body);
  color: var(--text-muted);
  max-width: 700px;
  line-height: 1.5;
}

.maint-alert {
  padding: 1.25rem 2rem;
  font-size: var(--step-small);
  font-weight: 700;
  margin-bottom: 2.5rem;
}

.maint-alert--success {
  background-color: #1b382b;
  color: #a8f5c8;
}

.maint-alert--error {
  background-color: #3d1216;
  color: #ffb3b8;
}

.maint-layout-grid {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 3.5rem;
  align-items: start;
}

.maint-form-card,
.maint-history-card {
  background-color: var(--bg-surface);
  padding: 2.5rem;
}

.card-title {
  font-size: var(--step-lead);
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text-ink);
  margin-bottom: 0.5rem;
}

.card-desc {
  font-size: var(--step-small);
  color: var(--text-muted);
  margin-bottom: 2rem;
  line-height: 1.45;
}

.maint-inputs-stack {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.maint-unit {
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

.unit-textarea {
  resize: vertical;
  line-height: 1.45;
}

.btn-maint-submit {
  padding: 1.25rem;
  background-color: var(--accent-crimson);
  color: var(--accent-ink);
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  text-align: center;
  transition: opacity 140ms ease;
  margin-top: 0.5rem;
}

.btn-maint-submit:hover {
  opacity: 0.88;
}

.blank-orders-notice {
  font-size: var(--step-body);
  color: var(--text-muted);
  padding: 2rem 0;
}

.orders-stack {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.order-item {
  background-color: var(--bg-ground);
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.order-top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.order-equipment {
  font-size: var(--step-body);
  font-weight: 700;
  color: var(--text-ink);
}

.order-badge {
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.badge-pendiente {
  color: #f7b731;
}

.badge-en_revision {
  color: #3867d6;
}

.badge-reparado {
  color: #2ed573;
}

.order-defect {
  font-size: var(--step-small);
  color: var(--text-muted);
  line-height: 1.45;
}

.order-data-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  background-color: var(--bg-surface);
  padding: 1rem 1.25rem;
}

.order-data-unit {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.data-tag {
  font-size: var(--step-small);
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.data-value {
  font-size: var(--step-small);
  font-weight: 700;
  color: var(--text-ink);
}

.order-technical-report {
  background-color: var(--bg-surface);
  padding: 1rem 1.25rem;
}

.report-tag {
  display: block;
  font-size: var(--step-small);
  font-weight: 700;
  color: var(--accent-crimson);
  margin-bottom: 0.25rem;
}

.report-text {
  font-size: var(--step-small);
  color: var(--text-ink);
  line-height: 1.45;
}

@media (max-width: 900px) {
  .maint-layout-grid {
    grid-template-columns: 1fr;
  }
}
</style>
