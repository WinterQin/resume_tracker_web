import { createRouter, createWebHistory } from 'vue-router'

// 定义需要登录才能访问的路由
const authRoutes = ['Dashboard', 'Profile', 'Applications'] // 需要登录才能访问的路由名称

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/applications',
    name: 'Applications',
    component: () => import('../views/Applications.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 定义路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const isAuthRequired = authRoutes.includes(to.name)

  if (isAuthRequired && !token) {
    // 需要认证但未登录，重定向到登录页
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (token && (to.name === 'Login' || to.name === 'Register')) {
    // 已登录用户不能访问登录和注册页
    next({ name: 'Dashboard' })
  } else {
    // 其他情况正常放行
    next()
  }
})

export default router 