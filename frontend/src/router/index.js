import { createRouter, createWebHistory } from 'vue-router'
import CatalogoView from '../views/CatalogoView.vue'
import AuthView from '../views/AuthView.vue'
import ReservasView from '../views/ReservasView.vue'
import MantenimientoView from '../views/MantenimientoView.vue'
import CajaView from '../views/CajaView.vue'

const routes = [
  { path: '/', redirect: '/catalogo' },
  { path: '/catalogo', component: CatalogoView },
  { path: '/login', component: AuthView },
  { path: '/reservas', component: ReservasView, meta: { requiereAuth: true } },
  { path: '/mantenimiento', component: MantenimientoView, meta: { requiereAuth: true } },
  { path: '/caja', component: CajaView, meta: { requiereAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Guardián de navegación profesional
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('gzs_token')
  if (to.meta.requiereAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router