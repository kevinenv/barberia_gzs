<script setup>
import { ref, onMounted } from 'vue'
import http from '../api/http'

const comisiones = ref([])
const productos = ref([])

const productoId = ref(1)
const cantidad = ref(1)
const metodoPago = ref('efectivo')
const mensaje = ref('')
const error = ref('')

async function cargarCaja() {
  try {
    const [resCom, resProd] = await Promise.all([
      http.get('/caja/comisiones'),
      http.get('/productos')
    ])
    comisiones.value = resCom.data
    productos.value = resProd.data
    if (productos.value.length > 0 && !productoId.value) {
      productoId.value = productos.value[0].id
    }
  } catch (err) {
    console.error('Error al cargar datos de caja:', err)
  }
}

onMounted(cargarCaja)

async function registrarVentaMostrador() {
  mensaje.value = ''
  error.value = ''
  try {
    const { data } = await http.post('/ventas-productos', {
      producto_id: Number(productoId.value),
      cantidad: Number(cantidad.value),
      metodo_pago: metodoPago.value
    })
    mensaje.value = `${data.mensaje}. Total cobrado: Bs. ${data.subtotal_bs} (Stock restante: ${data.stock_restante} uds.)`
    await cargarCaja()
  } catch (err) {
    error.value = err.response?.data?.error || 'Error al registrar venta'
  }
}
</script>

<template>
  <div class="monument-caja">
    <header class="caja-header">
      <span class="caja-category">Administración &amp; Finanzas</span>
      <h1 class="caja-title">CAJA &bull; ARQUEO Y COMISIONES</h1>
      <p class="caja-sub">
        Liquidación de comisiones según el esquema 50/50 por sillón de trabajo y control de salidas de mostrador en tiempo real.
      </p>
    </header>

    <!-- TABLA DE COMISIONES 50/50 -->
    <section class="caja-section">
      <div class="section-top-bar">
        <h2 class="caja-h2">Liquidación de Comisiones por Barbero (50% / 50%)</h2>
        <span class="caja-badge">Reparto Directo</span>
      </div>

      <div class="ledger-table-wrap">
        <table class="monument-table">
          <thead>
            <tr>
              <th>Estación</th>
              <th>Especialista</th>
              <th>Servicios</th>
              <th>Total Generado</th>
              <th>Comisión Barbero (50%)</th>
              <th>Salón GZS (50%)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="b in comisiones" :key="b.barbero_id">
              <td>
                <span class="chair-tag">Sillón #{{ b.silla_numero }}</span>
              </td>
              <td><span class="barber-bold">{{ b.apodo }}</span></td>
              <td>{{ b.total_servicios }} citas</td>
              <td>Bs. {{ b.total_generado.toFixed(2) }}</td>
              <td><span class="highlight-barber">Bs. {{ b.comision_ganada.toFixed(2) }}</span></td>
              <td><span class="highlight-salon">Bs. {{ b.ganancia_barberia.toFixed(2) }}</span></td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- MOSTRADOR Y CONTROL DE INVENTARIO -->
    <section class="caja-section">
      <div class="section-top-bar">
        <h2 class="caja-h2">Ventas de Mostrador y Control de Stock</h2>
        <span class="caja-badge">Salidas Inmediatas</span>
      </div>

      <div v-if="mensaje" class="caja-alert caja-alert--success">{{ mensaje }}</div>
      <div v-if="error" class="caja-alert caja-alert--error">{{ error }}</div>

      <div class="caja-dual-grid">
        <!-- FORMULARIO DE COBRO -->
        <div class="caja-panel-box">
          <h3 class="panel-h3">Registrar Cobro de Producto</h3>
          <p class="panel-p">Seleccione el producto despachado en el mostrador del salón.</p>

          <form @submit.prevent="registrarVentaMostrador" class="caja-form-stack">
            <div class="caja-unit">
              <label class="unit-label" for="producto-sel">Producto</label>
              <select id="producto-sel" v-model="productoId" class="unit-field">
                <option v-for="p in productos" :key="p.id" :value="p.id">
                  {{ p.nombre }} &mdash; Bs. {{ p.precio }} (Stock: {{ p.stock }})
                </option>
              </select>
            </div>

            <div class="caja-unit">
              <label class="unit-label" for="cant-input">Cantidad Despachada</label>
              <input
                id="cant-input"
                v-model="cantidad"
                type="number"
                min="1"
                class="unit-field"
                required
              />
            </div>

            <div class="caja-unit">
              <label class="unit-label">Forma de Pago</label>
              <div class="payment-duo">
                <button
                  type="button"
                  class="pay-choice-btn"
                  :class="{ 'pay-choice-btn--active': metodoPago === 'efectivo' }"
                  @click="metodoPago = 'efectivo'"
                >
                  Efectivo
                </button>
                <button
                  type="button"
                  class="pay-choice-btn"
                  :class="{ 'pay-choice-btn--active': metodoPago === 'qr' }"
                  @click="metodoPago = 'qr'"
                >
                  Transferencia QR
                </button>
              </div>
            </div>

            <button type="submit" class="btn-caja-submit">
              Procesar Venta de Mostrador
            </button>
          </form>
        </div>

        <!-- BALANCE DE STOCK -->
        <div class="caja-panel-box">
          <h3 class="panel-h3">Balance de Inventario en Salón</h3>
          <p class="panel-p">Disponibilidad en góndola y alertas de reabastecimiento.</p>

          <div class="inventory-table-wrap">
            <table class="monument-table">
              <thead>
                <tr>
                  <th>Producto</th>
                  <th>Precio</th>
                  <th>Stock</th>
                  <th>Estado</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="p in productos" :key="p.id">
                  <td><span class="product-name-bold">{{ p.nombre }}</span></td>
                  <td>Bs. {{ p.precio }}</td>
                  <td><span class="stock-qty">{{ p.stock }} uds</span></td>
                  <td>
                    <span
                      class="stock-pill"
                      :class="p.stock <= 5 ? 'stock-pill--low' : 'stock-pill--ok'"
                    >
                      {{ p.stock <= 5 ? 'Reabastecer' : 'Disponible' }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.monument-caja {
  max-width: 1400px;
  margin: 0 auto;
  padding: 3rem 2rem 6rem;
}

.caja-header {
  margin-bottom: 4rem;
}

.caja-category {
  display: block;
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--accent-crimson);
  margin-bottom: 0.5rem;
}

.caja-title {
  font-size: var(--step-head);
  font-weight: 900;
  letter-spacing: -0.05em;
  text-transform: uppercase;
  line-height: 0.95;
  color: var(--text-ink);
  margin-bottom: 1rem;
}

.caja-sub {
  font-size: var(--step-body);
  color: var(--text-muted);
  max-width: 700px;
  line-height: 1.5;
}

.caja-section {
  margin-bottom: 4rem;
}

.section-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 1.5rem;
}

.caja-h2 {
  font-size: var(--step-lead);
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text-ink);
}

.caja-badge {
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--accent-crimson);
}

.ledger-table-wrap,
.inventory-table-wrap {
  background-color: var(--bg-surface);
  padding: 1.5rem;
  overflow-x: auto;
}

.monument-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.monument-table th {
  padding: 1rem 1.25rem;
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-muted);
  background-color: var(--bg-ground);
}

.monument-table td {
  padding: 1.25rem;
  font-size: var(--step-small);
  color: var(--text-ink);
  background-color: var(--bg-surface);
}

.chair-tag {
  display: inline-block;
  padding: 0.35rem 0.65rem;
  background-color: var(--bg-ground);
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  color: var(--accent-crimson);
}

.barber-bold,
.product-name-bold {
  font-weight: 700;
  color: var(--text-ink);
}

.highlight-barber {
  font-weight: 800;
  color: #2ed573;
}

.highlight-salon {
  font-weight: 800;
  color: var(--text-ink);
}

.caja-alert {
  padding: 1.25rem 2rem;
  font-size: var(--step-small);
  font-weight: 700;
  margin-bottom: 2rem;
}

.caja-alert--success {
  background-color: #1b382b;
  color: #a8f5c8;
}

.caja-alert--error {
  background-color: #3d1216;
  color: #ffb3b8;
}

.caja-dual-grid {
  display: grid;
  grid-template-columns: 1fr 1.3fr;
  gap: 2.5rem;
  align-items: start;
}

.caja-panel-box {
  background-color: var(--bg-surface);
  padding: 2.5rem;
}

.panel-h3 {
  font-size: var(--step-lead);
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--text-ink);
  margin-bottom: 0.5rem;
}

.panel-p {
  font-size: var(--step-small);
  color: var(--text-muted);
  margin-bottom: 2rem;
  line-height: 1.45;
}

.caja-form-stack {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.caja-unit {
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

.payment-duo {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
}

.pay-choice-btn {
  padding: 1rem;
  background-color: var(--bg-ground);
  color: var(--text-muted);
  font-size: var(--step-small);
  font-weight: 700;
  text-align: center;
  transition: background-color 140ms ease, color 140ms ease;
}

.pay-choice-btn:hover {
  color: var(--text-ink);
}

.pay-choice-btn--active {
  background-color: var(--bg-surface-elevated);
  color: var(--text-ink);
  font-weight: 800;
}

.btn-caja-submit {
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

.btn-caja-submit:hover {
  opacity: 0.88;
}

.stock-qty {
  font-weight: 700;
  color: var(--text-ink);
}

.stock-pill {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  font-size: var(--step-small);
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stock-pill--ok {
  color: #2ed573;
}

.stock-pill--low {
  color: #ff4757;
}

@media (max-width: 900px) {
  .caja-dual-grid {
    grid-template-columns: 1fr;
  }
}
</style>
