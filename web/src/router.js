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
      props: route => ({
        path: `/${route.params.path}`//add '/' qwq
      })
    },
    {
      path: '/files',
      redirect: '/files/~/',
    },
    {
      path: '/library',
      name: 'library',
      component: () => import('./views/Library.vue'),
    },
    {
      path: '/library/:id',
      name: 'library-reader',
      component: () => import('./views/LibraryReader.vue'),
      props: true
    },
  ],
})

export default router
