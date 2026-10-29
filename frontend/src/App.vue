<script setup>
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

function cerrarSesion() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div class="app-container">
    <header class="navbar">
      <div class="brand" @click="router.push('/catalogo')">
        💈 <span>Barbería GZS</span>
      </div>
      <nav class="nav-links">
        <router-link to="/catalogo">Catálogo y Ubicación</router-link>
        <router-link to="/reservas">Reservar Cita</router-link>
        <router-link to="/mantenimiento">Sillas Hidráulicas</router-link>
        <router-link to="/caja">Caja y Comisiones</router-link>
      </nav>
      <div class="user-panel">
        <template v-if="auth.estaLogueado">
          <span class="user-badge">👤 {{ auth.usuario?.nombre }}</span>
          <button class="btn-outline" @click="cerrarSesion">Salir</button>
        </template>
        <router-link v-else to="/login" class="btn-gold-sm">Iniciar Sesión</router-link>
      </div>
    </header>

    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
  background-color: #121316;
  color: #f3f4f6;
  font-family: system-ui, -apple-system, sans-serif;
}
.app-container { max-width: 1050px; margin: 0 auto; padding: 0 1.5rem 3rem; }
.navbar {
  display: flex; justify-content: space-between; align-items: center;
  padding: 1.2rem 0; border-bottom: 2px solid #d4af37; margin-bottom: 2rem; flex-wrap: wrap; gap: 1rem;
}
.brand { font-size: 1.5rem; font-weight: 800; color: #ffffff; cursor: pointer; }
.brand span { color: #d4af37; }
.nav-links { display: flex; gap: 1.2rem; flex-wrap: wrap; }
.nav-links a { color: #d1d5db; text-decoration: none; font-weight: 500; padding: 0.4rem 0.7rem; border-radius: 6px; }
.nav-links a.router-link-active { background: #d4af37; color: #111; font-weight: 700; }
.user-panel { display: flex; align-items: center; gap: 0.8rem; }
.user-badge { font-size: 0.9rem; color: #fbbf24; font-weight: 600; }
.hero-card {
  background: #1c1e24; border: 1px solid #d4af37; border-radius: 10px;
  padding: 1.5rem; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 1rem;
}
.hero-card h2 { margin-top: 0; color: #fbbf24; }
.section-title { margin-top: 2.2rem; border-left: 4px solid #d4af37; padding-left: 0.75rem; color: #fff; }
.grid-cards { display: grid; grid-template-columns: repeat(auto-fill, minmax(230px, 1fr)); gap: 1rem; }
.card { background: #1c1e24; border: 1px solid #2e323b; border-radius: 10px; padding: 1.2rem; }
.card h3 { margin-top: 0; color: #fff; }
.muted { color: #9ca3af; font-size: 0.9rem; }
.price-tag { font-size: 1.15rem; font-weight: 700; color: #fbbf24; margin-top: 0.8rem; }
.price-tag span { font-size: 0.85rem; color: #9ca3af; font-weight: normal; }
.badge-sillon { display: inline-block; background: #d4af37; color: #111; font-weight: 700; font-size: 0.8rem; padding: 0.2rem 0.6rem; border-radius: 4px; margin-bottom: 0.5rem; }
.stock-ok { color: #34d399; font-weight: 700; margin-bottom: 0; }
.stock-low { color: #f87171; font-weight: 700; margin-bottom: 0; }
.btn-gold { background: #d4af37; color: #111; border: none; padding: 0.8rem 1.4rem; border-radius: 8px; font-weight: 700; cursor: pointer; font-size: 1rem; }
.btn-gold:hover { background: #e5c158; }
.btn-gold-sm { background: #d4af37; color: #111; text-decoration: none; padding: 0.5rem 1rem; border-radius: 6px; font-weight: 700; font-size: 0.9rem; }
.btn-outline { background: transparent; color: #f87171; border: 1px solid #f87171; padding: 0.4rem 0.8rem; border-radius: 6px; cursor: pointer; }
.auth-box { max-width: 420px; margin: 2rem auto; }
.form-vertical { display: flex; flex-direction: column; gap: 1rem; margin-top: 1rem; }
.form-vertical label { display: block; margin-bottom: 0.3rem; font-size: 0.9rem; color: #d1d5db; }
.form-vertical input, .form-vertical select, .form-vertical textarea {
  width: 100%; box-sizing: border-box; padding: 0.7rem; border-radius: 6px;
  border: 1px solid #374151; background: #111317; color: #fff;
}
.alert-error { background: rgba(239,68,68,0.15); border: 1px solid #ef4444; color: #fca5a5; padding: 0.75rem; border-radius: 6px; margin-bottom: 1rem; }
.switch-auth { text-align: center; margin-top: 1.2rem; font-size: 0.9rem; }
.switch-auth a { color: #fbbf24; }


.booking-layout { display: grid; grid-template-columns: 1.3fr 1fr; gap: 1.5rem; }
@media (max-width: 768px) { .booking-layout { grid-template-columns: 1fr; } }
.service-selector { display: flex; flex-direction: column; gap: 0.6rem; }
.service-item {
  display: flex; justify-content: space-between; align-items: center;
  padding: 0.8rem 1rem; border: 1px solid #374151; border-radius: 8px; cursor: pointer; background: #14161b;
}
.service-item.selected { border: 2px solid #d4af37; background: rgba(212, 175, 55, 0.12); }
.row-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.summary-line { display: flex; justify-content: space-between; padding: 0.65rem 0; border-bottom: 1px solid #2e323b; }
.highlight-gold { color: #fbbf24; font-size: 1.15rem; }
.qr-placeholder {
  margin: 1.2rem 0; padding: 2rem 1rem; text-align: center;
  border: 2px dashed #d4af37; border-radius: 8px; color: #fbbf24; font-weight: bold; background: #14161b;
}
.alert-success { background: rgba(16,185,129,0.15); border: 1px solid #10b981; color: #6ee7b7; padding: 0.85rem; border-radius: 6px; margin-bottom: 1rem; }
.data-table { width: 100%; border-collapse: collapse; text-align: left; }
.data-table th, .data-table td { padding: 0.85rem; border-bottom: 1px solid #2e323b; }
.data-table th { color: #fbbf24; font-size: 0.9rem; }
.maint-item { padding: 0.85rem 0; border-bottom: 1px solid #2e323b; }



</style>