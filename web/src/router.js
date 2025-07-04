import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('./views/Home.vue'),
    },
    {
      path: '/files/:path(.*)',
      name: 'files',
      component: () => import('./views/Files.vue'),
      props:true
    },
  ],
})

export default router
