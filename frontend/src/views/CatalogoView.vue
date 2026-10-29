<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import http from '../api/http'

const router = useRouter()
const config = ref(null)
const servicios = ref([])
const productos = ref([])
const barberos = ref([])

onMounted(async () => {
  const [resConfig, resServ, resProd, resBarb] = await Promise.all([
    http.get('/configuracion'),
    http.get('/servicios'),
    http.get('/productos'),
    http.get('/barberos')
  ])
  config.value = resConfig.data
  servicios.value = resServ.data
  productos.value = resProd.data
  barberos.value = resBarb.data
})
</script>

<template>
  <div>
    <!-- Banner de información esencial y ubicación -->
    <div v-if="config" class="hero-card">
      <div>
        <h2>Bienvenido a {{ config.nombre_negocio }}</h2>
        <p>📍 <strong>Ubicación:</strong> {{ config.direccion || 'Consultar en recepción' }}</p>
        <p>🕒 <strong>Horario de Atención:</strong> {{ config.hora_apertura }} a {{ config.hora_cierre }} (Descanso: {{ config.hora_inicio_descanso }} - {{ config.hora_fin_descanso }})</p>
        <p>🔒 <strong>Garantía de Reserva:</strong> Adelanto del {{ config.porcentaje_garantia_reserva }}% vía QR</p>
      </div>
      <button class="btn-gold" @click="router.push('/reservas')">
        📅 Agendar Cita Ahora
      </button>
    </div>

    <!-- Servicios de Barbería -->
    <h2 class="section-title">✂️ Servicios de Corte, Barba y Cejas</h2>
    <div class="grid-cards">
      <div v-for="s in servicios" :key="s.id" class="card">
        <h3>{{ s.nombre }}</h3>
        <p class="muted">{{ s.descripcion }}</p>
        <div class="price-tag">Bs. {{ s.precio }} <span>({{ s.duracion_minutos }} min)</span></div>
      </div>
    </div>

    <!-- Sillones y Barberos -->
    <h2 class="section-title">💈 Sillones y Especialistas</h2>
    <div class="grid-cards">
      <div v-for="b in barberos" :key="b.id" class="card barber-card">
        <div class="badge-sillon">Sillón #{{ b.silla_numero }}</div>
        <h3>{{ b.apodo }}</h3>
        <p class="muted">Especialidad: {{ b.especialidad }}</p>
      </div>
    </div>

    <!-- Productos en Mostrador -->
    <h2 class="section-title">🧴 Productos en Mostrador (Stock en Tiempo Real)</h2>
    <div class="grid-cards">
      <div v-for="p in productos" :key="p.id" class="card">
        <h3>{{ p.nombre }}</h3>
        <p class="muted">{{ p.descripcion }}</p>
        <p class="price-tag">Bs. {{ p.precio_venta }}</p>
        <p :class="p.stock_actual > p.stock_minimo ? 'stock-ok' : 'stock-low'">
          Disponibles: {{ p.stock_actual }} uds.
        </p>
      </div>
    </div>
  </div>
</template>