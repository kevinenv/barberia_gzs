<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import http from '../api/http'

import heroImg from '../assets/hero.png'
import razorCraftImg from '../assets/a-close-shave-in-antigua-guatemala.jpg'
import vintageRazorImg from '../assets/vintage-gem-micromatic-safety-razor-american-safety.jpg'
import historyEtchingImg from '../assets/razor-s-levee-or-ye-heads-of.jpg'
import barberpoleVideo from '../assets/barberpole-japan.webm'
import barberpolePoster from '../assets/barberpole-japan-poster.jpg'

const router = useRouter()

const fallbackConfig = {
  hora_apertura: '09:00',
  hora_cierre: '20:00',
  porcentaje_garantia_reserva: 30,
  telefono_contacto: '+591 70000000'
}

const fallbackServicios = [
  { id: 1, nombre: 'Corte Clásico de Caballero', descripcion: 'Diagnóstico capilar, corte a tijera o máquina, lavado previo y toalla caliente aromática.', precio: 50, duracion_minutos: 45 },
  { id: 2, nombre: 'Perfilado de Barba & Navaja', descripcion: 'Tratamiento de vaporización, aceites emolientes y afeitado tradicional a filo recto.', precio: 40, duracion_minutos: 35 },
  { id: 3, nombre: 'Servicio Completo Monumental', descripcion: 'Corte de autor, perfilado integral con navaja, toalla caliente y masaje capilar.', precio: 80, duracion_minutos: 75 },
  { id: 4, nombre: 'Afeitado Clásico Completo', descripcion: 'Espuma batida en tazón, triple pasada de navaja y loción balsámica calmante.', precio: 45, duracion_minutos: 40 }
]

const fallbackProductos = [
  { id: 1, nombre: 'Pomada Arcilla Acabado Mate', categoria: 'Fijación', descripcion: 'Acabado natural sin brillo y fijación flexible para todo el día.', precio: 85, stock: 12 },
  { id: 2, nombre: 'Aceite Esencial para Barba', categoria: 'Cuidado Facial', descripcion: 'Fórmula hidratante con jojoba y cedro para suavizar el vello facial.', precio: 70, stock: 8 },
  { id: 3, nombre: 'Tónico Refrescante Capilar', categoria: 'Tratamiento', descripcion: 'Estimulación del cuero cabelludo con extracto de romero y mentol.', precio: 65, stock: 15 }
]

const fallbackBarberos = [
  { id: 1, apodo: 'Gonzalo', silla_numero: 1, especialidad: 'Especialista en degradados limpios y tijera clásica' },
  { id: 2, apodo: 'Javier', silla_numero: 2, especialidad: 'Maestro de afeitado con navaja recta y barbas largas' },
  { id: 3, apodo: 'Mateo', silla_numero: 3, especialidad: 'Estilismo contemporáneo y perfilado tradicional' }
]

const config = ref(fallbackConfig)
const servicios = ref(fallbackServicios)
const productos = ref(fallbackProductos)
const barberos = ref(fallbackBarberos)

function updateScrollMotion() {
  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    document.querySelectorAll('.scroll-target').forEach((el) => {
      el.style.transform = 'none'
      el.style.opacity = '1'
    })
    return
  }

  const targets = document.querySelectorAll('.scroll-target')
  targets.forEach((el) => {
    const rect = el.getBoundingClientRect()
    const inView = rect.top < window.innerHeight && rect.bottom > 0
    if (inView) {
      const factor = Math.min(1, Math.max(0, (window.innerHeight - rect.top) / (window.innerHeight + rect.height)))
      const translateY = Math.round((1 - factor) * 12)
      const opacity = Math.min(1, Math.max(0.5, factor * 1.4))
      el.style.transform = `translateY(${translateY}px)`
      el.style.opacity = `${opacity}`
    }
  })
}

onMounted(async () => {
  try {
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
  } catch (error) {
    console.error('Error al cargar datos del catálogo:', error)
  }

  window.addEventListener('scroll', updateScrollMotion, { passive: true })
  window.addEventListener('resize', updateScrollMotion, { passive: true })
  // Initial motion trigger
  setTimeout(updateScrollMotion, 100)
})

onUnmounted(() => {
  window.removeEventListener('scroll', updateScrollMotion)
  window.removeEventListener('resize', updateScrollMotion)
})
</script>

<template>
  <div class="monument-catalog">
    <!-- HERO: NAME ACROSS THE MIDDLE (direction.json) -->
    <section class="monument-hero-stage">
      <div class="hero-top-flank scroll-target">
        <div class="flank-left scroll-target">
          <p class="hero-descriptor">
            Salón de corte clásico, navaja de afeitar y toalla caliente en La Paz. Citas programadas en sillón individual sin esperas.
          </p>
          <button class="btn-hero-action scroll-target" @click="router.push('/reservas')">
            Agendar turno en línea
          </button>
        </div>

        <div class="flank-right">
          <figure class="hero-video-box scroll-target">
            <video
              autoplay
              muted
              loop
              playsinline
              preload="metadata"
              :poster="barberpolePoster"
              class="barberpole-loop-video"
              width="240"
              height="135"
              aria-hidden="true"
            >
              <source :src="barberpoleVideo" type="video/webm" />
            </video>
          </figure>
        </div>
      </div>

      <!-- MASSIVE ARCHITECTURAL WORDMARK WITH FULL EXPLANATION LINE -->
      <div class="monument-wordmark-container scroll-target">
        <h1 class="monument-wordmark">BARBERÍA GZS</h1>
        <p class="monument-subheading">
          Salón de corte clásico, afeitado a navaja y cuidado capilar en La Paz
        </p>
      </div>

      <div class="hero-bottom-flank scroll-target">
        <div class="baseline-left scroll-target">
          <span class="baseline-strong">La Paz, Bolivia</span>
          <span class="baseline-sub">Avenida 20 de Octubre y Rosendo Gutiérrez</span>
        </div>

        <div class="baseline-right scroll-target">
          <span class="baseline-strong">Lunes a sábado de 09:00 a 20:00</span>
          <span class="baseline-sub">Garantía de turno con anticipo QR 30%</span>
        </div>
      </div>
    </section>

    <!-- SECTION 1: PANORAMA DEL SALÓN Y CIFRAS BASELINE -->
    <section class="monument-section scroll-target">
      <div class="section-headline-group scroll-target">
        <h2 class="section-display-title">El salón</h2>
        <p class="section-lead-copy">
          Un espacio diseñado para la concentración del barbero y la comodidad del cliente. Cada estación cuenta con sillón hidráulico de ajuste milimétrico, iluminación calibrada e higiene rigurosa entre servicios.
        </p>
      </div>

      <figure class="panorama-frame scroll-target">
        <img :src="heroImg" alt="Interior del salón de Barbería GZS en La Paz" class="panorama-img" />
      </figure>

      <div class="facts-strip scroll-target" v-if="config">
        <div class="fact-node scroll-target">
          <span class="fact-num">{{ config.hora_apertura }}</span>
          <span class="fact-text">Apertura del primer turno del día</span>
        </div>
        <div class="fact-node scroll-target">
          <span class="fact-num">{{ config.hora_cierre }}</span>
          <span class="fact-text">Cierre de estaciones y limpieza</span>
        </div>
        <div class="fact-node scroll-target">
          <span class="fact-num">{{ config.porcentaje_garantia_reserva }}%</span>
          <span class="fact-text">Garantía de confirmación vía QR</span>
        </div>
      </div>
    </section>

    <!-- SECTION 2: ESTACIONES Y BARBEROS (SPLIT-STAGE COMPOSITION) -->
    <section class="monument-section scroll-target">
      <div class="section-headline-group scroll-target">
        <h2 class="section-display-title">Las estaciones de trabajo</h2>
        <p class="section-lead-copy">
          Tres sillones hidráulicos asignados a maestros barberos. Elija a su especialista de confianza al momento de solicitar su turno.
        </p>
      </div>

      <div class="split-stations-stage scroll-target">
        <div class="stage-visual">
          <figure class="craft-figure scroll-target">
            <img
              :src="razorCraftImg"
              alt="Afeitado clásico a navaja de filo recto"
              class="craft-photo"
              loading="lazy"
            />
            <figcaption class="craft-caption">
              Afeitado tradicional con navaja de filo recto y toalla caliente.
            </figcaption>
          </figure>
        </div>

        <div class="stage-stations-list">
          <article v-for="b in barberos" :key="b.id" class="station-row scroll-target">
            <div class="station-data">
              <h3 class="station-title">{{ b.apodo }} &mdash; Sillón #{{ b.silla_numero }}</h3>
              <p class="station-specialty">{{ b.especialidad }}</p>
            </div>
            <button class="btn-station-reserve" @click="router.push('/reservas')">
              Reservar con {{ b.apodo }}
            </button>
          </article>
        </div>
      </div>
    </section>

    <!-- SECTION 3: CARTA DE SERVICIOS (TWO-ENDED ROWS) -->
    <section id="carta-servicios" class="monument-section scroll-target">
      <div class="section-headline-group scroll-target">
        <h2 class="section-display-title">Carta de cortes y servicios</h2>
        <p class="section-lead-copy">
          Técnicas de corte contemporáneo y afeitado tradicional a navaja. Todos los turnos incluyen toalla caliente, lavado de cortesía y productos de terminación.
        </p>
      </div>

      <div class="service-ledger scroll-target">
        <article
          v-for="s in servicios"
          :key="s.id"
          class="service-line scroll-target"
          @click="router.push('/reservas')"
        >
          <div class="service-left">
            <h3 class="service-title">{{ s.nombre }}</h3>
            <p class="service-meta">{{ s.descripcion }}</p>
          </div>

          <div class="service-right">
            <span class="service-duration">{{ s.duracion_minutos }} minutos</span>
            <span class="service-price">Bs. {{ s.precio }}</span>
          </div>
        </article>
      </div>
    </section>

    <!-- SECTION 4: PRODUCTOS E INSTRUMENTAL (EDITORIAL FEATURE) -->
    <section class="monument-section scroll-target">
      <div class="section-headline-group scroll-target">
        <h2 class="section-display-title">Góndola de productos y cuidado</h2>
        <p class="section-lead-copy">
          Fórmulas profesionales seleccionadas para el mantenimiento de barba, cabello y cuero cabelludo en el hogar.
        </p>
      </div>

      <div class="editorial-products-stage scroll-target">
        <figure class="tooling-figure scroll-target">
          <img
            :src="vintageRazorImg"
            alt="Maquinaria de navaja clásica para cuidado de precisión"
            class="tooling-photo"
            loading="lazy"
          />
          <figcaption class="craft-caption">
            Instrumental clásico y hojas de precisión calibradas para cada servicio.
          </figcaption>
        </figure>

        <div class="products-stack-list">
          <article v-for="p in productos" :key="p.id" class="product-item-row scroll-target">
            <div class="product-desc-box">
              <span class="product-category-text">{{ p.categoria || 'Cuidado' }}</span>
              <h3 class="product-name">{{ p.nombre }}</h3>
              <p class="product-p">{{ p.descripcion }}</p>
            </div>
            <div class="product-commercial-box">
              <span class="product-price">Bs. {{ p.precio }}</span>
              <span class="product-inventory">{{ p.stock }} unidades disponibles</span>
            </div>
          </article>
        </div>
      </div>
    </section>

    <!-- SECTION 5: HISTORIA DEL OFICIO Y PROTOCOLO -->
    <section class="monument-section scroll-target">
      <div class="section-headline-group scroll-target">
        <h2 class="section-display-title">Oficio y protocolo de salón</h2>
        <p class="section-lead-copy">
          El respeto por el tiempo del cliente y la herencia de la navaja artesanal forman los cimientos de nuestra atención diaria.
        </p>
      </div>

      <div class="protocol-columns-grid scroll-target">
        <div class="protocol-column scroll-target">
          <h3 class="column-title">Puntualidad de agenda</h3>
          <p class="column-body">
            Trabajamos con bloques de tiempo reservados. Si llega 10 minutos antes, podrá relajarse en nuestra sala de espera con café de cortesía.
          </p>
        </div>

        <div class="protocol-column scroll-target">
          <h3 class="column-title">Anticipo QR previo</h3>
          <p class="column-body">
            Para asegurar su sillón se solicita el 30% de anticipo. El saldo restante de su servicio se cancela en caja al finalizar.
          </p>
        </div>

        <div class="protocol-column scroll-target">
          <h3 class="column-title">Asepsia de navajas</h3>
          <p class="column-body">
            Utilizamos navajas con hojas descartables individuales abiertas frente al cliente y esterilización ultravioleta de instrumental.
          </p>
        </div>
      </div>

      <!-- HISTORIC ARCHIVAL ETCHING -->
      <figure class="etching-banner scroll-target">
        <img
          :src="historyEtchingImg"
          alt="Grabado histórico del oficio de barbero del siglo XVIII, The Metropolitan Museum of Art"
          class="etching-img"
          loading="lazy"
        />
        <figcaption class="craft-caption">
          Grabado histórico del arte del afeitado y peluquería (Colección The Metropolitan Museum of Art, Nueva York, Dominio Público CC0).
        </figcaption>
      </figure>
    </section>
  </div>
</template>

<style scoped>
.monument-catalog {
  width: 100%;
  max-width: 100%;
  overflow-x: clip;
}

/* HERO: NAME ACROSS THE MIDDLE */
.monument-hero-stage {
  min-height: calc(100svh - 80px);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 2.5rem clamp(1.25rem, 4vw, 3rem);
  background-color: var(--bg-ground);
  width: 100%;
  box-sizing: border-box;
}

.hero-top-flank {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 2rem;
  width: 100%;
}

.flank-left {
  max-width: 520px;
}

.hero-descriptor {
  font-size: var(--step-body);
  font-weight: 500;
  line-height: 1.5;
  color: var(--text-ink);
  margin-bottom: 1.5rem;
}

.btn-hero-action {
  display: inline-block;
  padding: 1rem 2rem;
  background-color: var(--accent-crimson);
  color: var(--accent-ink);
  font-size: var(--step-body);
  font-weight: 700;
  transition: transform 140ms ease, opacity 140ms ease;
}

.btn-hero-action:hover {
  opacity: 0.9;
  transform: translateY(-2px);
}

.flank-right {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
}

.hero-video-box {
  margin: 0;
  overflow: hidden;
  background-color: var(--bg-surface);
}

.barberpole-loop-video {
  width: 240px;
  height: 135px;
  object-fit: cover;
  display: block;
}

.monument-wordmark-container {
  padding: 2.5rem 0;
  overflow: hidden;
  width: 100%;
}

.monument-wordmark {
  font-size: clamp(2.8rem, 12.6vw, 14.2rem);
  font-weight: 900;
  letter-spacing: -0.06em;
  line-height: 0.84;
  text-transform: uppercase;
  color: var(--text-ink);
  margin: 0;
  white-space: nowrap;
}

.monument-subheading {
  font-size: var(--step-lead);
  font-weight: 500;
  color: var(--text-muted);
  margin-top: 1.25rem;
  line-height: 1.4;
}

.hero-bottom-flank {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 2rem;
  padding-top: 2rem;
  width: 100%;
}

.baseline-left,
.baseline-right {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.baseline-right {
  text-align: right;
}

.baseline-strong {
  font-size: var(--step-body);
  font-weight: 700;
  color: var(--text-ink);
}

.baseline-sub {
  font-size: var(--step-small);
  color: var(--text-muted);
}

/* SECTIONS */
.monument-section {
  max-width: 1400px;
  margin: 0 auto;
  padding: 6rem clamp(1.25rem, 4vw, 3rem) 0;
  width: 100%;
  box-sizing: border-box;
}

/* MOTION LAYER */
.scroll-target {
  transition: transform 300ms cubic-bezier(0.16, 1, 0.3, 1), opacity 300ms cubic-bezier(0.16, 1, 0.3, 1);
  will-change: transform, opacity;
}

.section-headline-group {
  margin-bottom: 3.5rem;
  max-width: 720px;
}

.section-display-title {
  font-size: var(--step-head);
  font-weight: 800;
  letter-spacing: -0.04em;
  line-height: 1.05;
  margin-bottom: 1.25rem;
  color: var(--text-ink);
}

.section-lead-copy {
  font-size: var(--step-body);
  color: var(--text-muted);
  line-height: 1.55;
}

/* PANORAMA */
.panorama-frame {
  width: 100%;
  max-width: 100%;
  max-height: 600px;
  overflow: hidden;
  background-color: var(--bg-surface);
  margin-bottom: 3rem;
  margin-inline: 0;
}

.panorama-img {
  width: 100%;
  max-width: 100%;
  height: auto;
  object-fit: cover;
  display: block;
}

/* FACTS STRIP (VARY COMPOSITION) */
.facts-strip {
  display: flex;
  justify-content: space-between;
  gap: 2rem;
  padding: 2.5rem 0;
  width: 100%;
}

.fact-node {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.fact-num {
  font-size: clamp(2.2rem, 3.8vw, 3.8rem);
  font-weight: 800;
  letter-spacing: -0.04em;
  color: var(--text-ink);
}

.fact-text {
  font-size: var(--step-body);
  color: var(--text-muted);
}

/* SPLIT-STAGE STATIONS */
.split-stations-stage {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 3.5rem;
  align-items: start;
  width: 100%;
}

.craft-figure,
.tooling-figure {
  margin: 0;
  width: 100%;
}

.craft-photo,
.tooling-photo {
  width: 100%;
  max-width: 100%;
  height: 480px;
  object-fit: cover;
  display: block;
  background-color: var(--bg-surface);
}

.craft-caption {
  font-size: var(--step-small);
  color: var(--text-muted);
  margin-top: 0.75rem;
  line-height: 1.4;
}

.stage-stations-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
}

.station-row {
  background-color: var(--bg-surface);
  padding: 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
  transition: transform 180ms ease, background-color 180ms ease;
}

.station-row:hover {
  transform: translateX(6px);
  background-color: var(--bg-surface-elevated);
}

.station-title {
  font-size: var(--step-lead);
  font-weight: 800;
  color: var(--text-ink);
  margin-bottom: 0.35rem;
}

.station-specialty {
  font-size: var(--step-body);
  color: var(--text-muted);
}

.btn-station-reserve {
  padding: 0.85rem 1.5rem;
  background-color: var(--bg-surface-elevated);
  color: var(--text-ink);
  font-size: var(--step-body);
  font-weight: 700;
  transition: background-color 140ms ease, color 140ms ease;
}

.station-row:hover .btn-station-reserve {
  background-color: var(--accent-crimson);
  color: var(--accent-ink);
}

/* SERVICE LEDGER (TWO-ENDED ROWS) */
.service-ledger {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  width: 100%;
}

.service-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.75rem 2rem;
  background-color: var(--bg-surface);
  cursor: pointer;
  width: 100%;
  box-sizing: border-box;
  transition: transform 180ms ease, background-color 180ms ease;
}

.service-line:hover {
  transform: translateX(6px);
  background-color: var(--bg-surface-elevated);
}

.service-left {
  max-width: 60%;
}

.service-title {
  font-size: var(--step-lead);
  font-weight: 700;
  margin-bottom: 0.35rem;
  color: var(--text-ink);
}

.service-meta {
  font-size: var(--step-body);
  color: var(--text-muted);
  line-height: 1.45;
}

.service-right {
  display: flex;
  align-items: baseline;
  gap: 2rem;
}

.service-duration {
  font-size: var(--step-body);
  font-weight: 600;
  color: var(--text-muted);
}

.service-price {
  font-size: var(--step-lead);
  font-weight: 800;
  color: var(--text-ink);
  font-variant-numeric: tabular-nums;
}

/* EDITORIAL PRODUCTS STAGE */
.editorial-products-stage {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 3.5rem;
  align-items: start;
  width: 100%;
}

.products-stack-list {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  width: 100%;
}

.product-item-row {
  background-color: var(--bg-surface);
  padding: 1.75rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  width: 100%;
  box-sizing: border-box;
}

.product-category-text {
  font-size: var(--step-small);
  font-weight: 700;
  color: var(--accent-bright);
  display: block;
  margin-bottom: 0.25rem;
}

.product-name {
  font-size: var(--step-lead);
  font-weight: 700;
  color: var(--text-ink);
  margin-bottom: 0.35rem;
}

.product-p {
  font-size: var(--step-small);
  color: var(--text-muted);
  line-height: 1.45;
}

.product-commercial-box {
  text-align: right;
  min-width: 140px;
}

.product-price {
  font-size: var(--step-lead);
  font-weight: 800;
  color: var(--text-ink);
  display: block;
  margin-bottom: 0.25rem;
}

.product-inventory {
  font-size: var(--step-small);
  color: var(--text-muted);
}

/* PROTOCOL */
.protocol-columns-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2.5rem;
  margin-bottom: 4rem;
  width: 100%;
}

.protocol-column {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.column-title {
  font-size: var(--step-lead);
  font-weight: 700;
  color: var(--text-ink);
}

.column-body {
  font-size: var(--step-body);
  color: var(--text-muted);
  line-height: 1.5;
}

/* ETCHING */
.etching-banner {
  margin: 0;
  background-color: var(--bg-surface);
  padding: 2.5rem;
  width: 100%;
  box-sizing: border-box;
}

.etching-img {
  width: 100%;
  max-width: 100%;
  max-height: 520px;
  object-fit: contain;
  display: block;
}

@media (max-width: 900px) {
  .monument-hero-stage {
    padding: 2rem clamp(1.25rem, 4vw, 2.5rem);
  }
  .hero-top-flank {
    flex-direction: column;
    align-items: stretch;
  }
  .hero-bottom-flank {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
  }
  .facts-strip {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
  .split-stations-stage,
  .editorial-products-stage {
    grid-template-columns: 1fr;
  }
  .protocol-columns-grid {
    grid-template-columns: 1fr;
  }
  .service-line {
    flex-direction: row;
    justify-content: space-between;
  }
  .product-item-row {
    flex-direction: row;
    justify-content: space-between;
  }
}

@media (max-width: 500px) {
  .hero-bottom-flank {
    flex-direction: column;
    align-items: stretch;
  }
  .baseline-right {
    text-align: left;
  }
  .service-line {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  .service-left {
    max-width: 100%;
  }
  .product-item-row {
    flex-direction: column;
    gap: 1rem;
  }
  .product-commercial-box {
    text-align: left;
  }
}
</style>
