import { shouldbeAuthenticated, useAuthentication } from '@/services/authentication'
import type { App } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import AuthenticationView from '../views/AuthenticationView.vue'
import HomeView from '../views/HomeView.vue'

const router = (app: App) => createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      beforeEnter: shouldbeAuthenticated(app)
    },
    {
      path: '/login',
      name: 'login',
      component: AuthenticationView
    },
    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue')
    // }
  ]
})

export default router
