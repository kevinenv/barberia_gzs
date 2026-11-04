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
  <div class="monument-shell">
    <!-- TOP STRIP -->
    <header class="monument-nav-strip">
      <div class="nav-left">
        <span class="nav-brand" @click="router.push('/catalogo')">GZS</span>
        <span class="nav-tagline">Salón masculino en La Paz</span>
      </div>

      <nav class="nav-links">
        <router-link to="/catalogo" class="nav-link">El salón</router-link>
        <router-link to="/reservas" class="nav-link">Turnos</router-link>
        <router-link to="/mantenimiento" class="nav-link">Sillones</router-link>
        <router-link to="/caja" class="nav-link">Caja</router-link>
      </nav>

      <div class="nav-auth">
        <template v-if="auth.estaLogueado">
          <span class="auth-name">{{ auth.usuario?.nombre }}</span>
          <button class="btn-text-action" @click="cerrarSesion">Salir</button>
        </template>
        <router-link v-else to="/login" class="btn-text-action">Acceso</router-link>
      </div>
    </header>

    <!-- MAIN VIEW -->
    <main class="monument-main">
      <router-view />
    </main>

    <!-- INVERTED MONUMENTAL FOOTER (direction.json: ground and ink swapped) -->
    <footer class="monument-footer-field">
      <div class="footer-wrap">
        <div class="footer-top-grid">
          <div class="footer-col">
            <span class="footer-col-title">Dirección y salón</span>
            <p class="footer-p">Avenida 20 de Octubre y Rosendo Gutiérrez, La Paz, Bolivia.</p>
            <p class="footer-sub">Tres estaciones de sillón hidráulico individual con servicio personalizado.</p>
          </div>

          <div class="footer-col">
            <span class="footer-col-title">Horario de atención</span>
            <p class="footer-p">Lunes a sábado de 09:00 a 20:00.</p>
            <p class="footer-sub">Atención puntual y programada sin filas de espera.</p>
          </div>

          <div class="footer-col">
            <span class="footer-col-title">Garantía de turno</span>
            <p class="footer-p">30% anticipo mediante QR para confirmar el sillón.</p>
            <p class="footer-sub">Comprobante validado en línea vía WhatsApp.</p>
          </div>

          <div class="footer-col">
            <span class="footer-col-title">Accesos directos</span>
            <div class="footer-nav-stack">
              <router-link to="/reservas" class="footer-nav-link">Agendar turno en línea</router-link>
              <router-link to="/catalogo" class="footer-nav-link">Carta de cortes</router-link>
              <router-link to="/mantenimiento" class="footer-nav-link">Taller de sillones</router-link>
            </div>
          </div>
        </div>

        <!-- GIANT CLOSING WORDMARK BRACKET -->
        <div class="footer-bracket">
          <h2 class="footer-giant-title">BARBERÍA GZS</h2>
        </div>

        <div class="footer-baseline">
          <span>&copy; {{ new Date().getFullYear() }} Barbería GZS. Salón masculino y navaja artesanal.</span>
          <span class="footer-credits-note">Archivos históricos: The Metropolitan Museum of Art (CC0) y Wikimedia Commons (CC BY-SA).</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<style>
:root {
  --bg-ground: #0a0a0a;
  --bg-surface: #151515;
  --bg-surface-elevated: #222222;
  --text-ink: #efefef;
  --text-muted: #a0a0a0;
  --accent-crimson: #b11226;
  --accent-ink: #ffffff;
  --accent-bright: #ff6b78;
  --field-bg: #efefef;
  --field-ink: #0a0a0a;

  --font-family: 'Schibsted Grotesk', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  
  --step-small: clamp(0.9375rem, 0.9rem + 0.1vw, 1rem);
  --step-body: clamp(1.0625rem, 1rem + 0.25vw, 1.1875rem);
  --step-lead: clamp(1.25rem, 1.1rem + 0.7vw, 1.625rem);
  --step-head: clamp(2rem, 1.3rem + 3.2vw, 4rem);
  --step-monument: clamp(3.2rem, 15.5vw, 16rem);
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  background-color: var(--bg-ground);
  color: var(--text-ink);
  font-family: var(--font-family);
  font-size: var(--step-body);
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
}

.monument-shell {
  min-height: 100svh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-ground);
  color: var(--text-ink);
}

/* TOP STRIP */
.monument-nav-strip {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem clamp(1.25rem, 4vw, 3rem);
  background-color: var(--bg-ground);
}

.nav-left {
  display: flex;
  align-items: baseline;
  gap: 1.25rem;
}

.nav-brand {
  font-size: 1.5rem;
  font-weight: 800;
  letter-spacing: -0.05em;
  color: var(--text-ink);
  cursor: pointer;
}

.nav-tagline {
  font-size: var(--step-body);
  color: var(--text-muted);
}

.nav-links {
  display: flex;
  gap: 2.25rem;
}

.nav-link {
  font-size: var(--step-body);
  font-weight: 600;
  color: var(--text-muted);
  transition: color 140ms ease;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--text-ink);
}

.nav-auth {
  display: flex;
  align-items: center;
  gap: 1.25rem;
}

.auth-name {
  font-size: var(--step-body);
  font-weight: 600;
  color: var(--text-ink);
}

.btn-text-action {
  font-size: var(--step-body);
  font-weight: 700;
  color: var(--accent-bright);
  transition: opacity 140ms ease;
}

.btn-text-action:hover {
  opacity: 0.8;
}

.monument-main {
  flex: 1;
}

/* INVERTED FOOTER */
.monument-footer-field {
  background-color: var(--field-bg);
  color: var(--field-ink);
  padding: 6rem clamp(1.25rem, 4vw, 3rem) 2.5rem;
  margin-top: 6rem;
}

.footer-wrap {
  max-width: 1400px;
  margin: 0 auto;
}

.footer-top-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 3rem;
  margin-bottom: 5rem;
}

.footer-col-title {
  display: block;
  font-size: var(--step-body);
  font-weight: 800;
  margin-bottom: 1rem;
  color: var(--field-ink);
}

.footer-p {
  font-size: var(--step-body);
  font-weight: 600;
  line-height: 1.4;
  margin-bottom: 0.5rem;
  color: var(--field-ink);
}

.footer-sub {
  font-size: var(--step-small);
  color: #444444;
  line-height: 1.45;
}

.footer-nav-stack {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.footer-nav-link {
  font-size: var(--step-body);
  font-weight: 700;
  color: var(--field-ink);
  transition: color 140ms ease;
}

.footer-nav-link:hover {
  color: var(--accent-crimson);
}

.footer-bracket {
  padding: 2rem 0;
  overflow: hidden;
}

.footer-giant-title {
  font-size: clamp(3rem, 13vw, 14rem);
  font-weight: 900;
  letter-spacing: -0.06em;
  line-height: 0.85;
  text-transform: uppercase;
  color: var(--field-ink);
  margin: 0;
  word-break: break-all;
}

.footer-baseline {
  display: flex;
  justify-content: space-between;
  padding-top: 2rem;
  font-size: var(--step-small);
  font-weight: 600;
  color: #555555;
}

.footer-credits-note {
  color: #666666;
}

@media (max-width: 768px) {
  .monument-nav-strip {
    flex-direction: column;
    gap: 1.25rem;
    align-items: flex-start;
  }
  .nav-links {
    gap: 1.5rem;
    flex-wrap: wrap;
  }
  .footer-baseline {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style>
