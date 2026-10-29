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
  const [resCom, resProd] = await Promise.all([
    http.get('/caja/comisiones'),
    http.get('/productos')
  ])
  comisiones.value = resCom.data
  productos.value = resProd.data
  if (productos.value.length > 0 && !productoId.value) {
    productoId.value = productos.value[0].id
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
    mensaje.value = `${data.mensaje} — Total cobrado: Bs. ${data.subtotal_bs} (Stock restante: ${data.stock_restante} uds.)`
    await cargarCaja()
  } catch (err) {
    error.value = err.response?.data?.error || 'Error al registrar venta'
  }
}
</script>

<template>
  <div>
    <h2 class="section-title">💰 Liquidación de Comisiones por Barbero (Regla 50% / 50%)</h2>
    <div class="card" style="margin-bottom: 2rem;">
      <table class="data-table">
        <thead>
          <tr>
            <th>Sillón</th>
            <th>Barbero</th>
            <th>Servicios Realizados</th>
            <th>Total Generado</th>
            <th>Comisión Barbero (50%)</th>
            <th>Ganancia Barbería (50%)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in comisiones" :key="b.barbero_id">
            <td><span class="badge-sillon">Sillón #{{ b.silla_numero }}</span></td>
            <td><strong>{{ b.apodo }}</strong></td>
            <td>{{ b.total_servicios }} servicios</td>
            <td>Bs. {{ b.total_generado.toFixed(2) }}</td>
            <td style="color: #34d399; font-weight: bold;">Bs. {{ b.comision_ganada.toFixed(2) }}</td>
            <td style="color: #fbbf24; font-weight: bold;">Bs. {{ b.ganancia_barberia.toFixed(2) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 class="section-title">🧴 Venta Rápida de Productos en Mostrador (Descuento de Stock)</h2>
    <div v-if="mensaje" class="alert-success">{{ mensaje }}</div>
    <div v-if="error" class="alert-error">{{ error }}</div>

    <div class="booking-layout">
      <div class="card">
        <h3>Registrar Salida de Inventario</h3>
        <form @submit.prevent="registrarVentaMostrador" class="form-vertical">
          <div>
            <label>Producto</label>
            <select v-model="productoId">
              <option v-for="p in productos" :key="p.id" :value="p.id">
                {{ p.nombre }} — Bs. {{ p.precio_venta }} (Stock: {{ p.stock_actual }})
              </option>
            </select>
          </div>
          <div class="row-2">
            <div>
              <label>Cantidad</label>
              <input type="number" min="1" v-model="cantidad" required />
            </div>
            <div>
              <label>Método de Pago</label>
              <select v-model="metodoPago">
                <option value="efectivo">💵 Efectivo</option>
                <option value="qr">📱 Pago QR</option>
                <option value="tarjeta">💳 Tarjeta</option>
              </select>
            </div>
          </div>
          <button type="submit" class="btn-gold">🛒 Cobrar y Descontar Stock</button>
        </form>
      </div>

      <div class="card">
        <h3>Estado Actual del Inventario</h3>
        <div v-for="p in productos" :key="p.id" class="summary-line">
          <span>{{ p.nombre }}</span>
          <strong :class="p.stock_actual > p.stock_minimo ? 'stock-ok' : 'stock-low'">
            {{ p.stock_actual }} uds.
          </strong>
        </div>
      </div>
    </div>
  </div>
</template>