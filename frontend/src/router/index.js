import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/home/index.vue'),
    meta: { title: '首页', keepAlive: true }
  },
  {
    path: '/shop-list',
    name: 'ShopList',
    component: () => import('@/views/shop-list/index.vue'),
    meta: { title: '商户列表' }
  },
  {
    path: '/shop-detail/:id',
    name: 'ShopDetail',
    component: () => import('@/views/shop-detail/index.vue'),
    meta: { title: '商户详情' }
  },
  {
    path: '/blog-detail/:id',
    name: 'BlogDetail',
    component: () => import('@/views/blog-detail/index.vue'),
    meta: { title: '笔记详情' }
  },
  {
    path: '/blog-edit',
    name: 'BlogEdit',
    component: () => import('@/views/blog-edit/index.vue'),
    meta: { title: '发布笔记', requireAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/profile/index.vue'),
    meta: { title: '个人主页', requireAuth: true }
  },
  {
    path: '/profile-edit',
    name: 'ProfileEdit',
    component: () => import('@/views/profile-edit/index.vue'),
    meta: { title: '编辑资料', requireAuth: true }
  },
  {
    path: '/user/:id',
    name: 'UserProfile',
    component: () => import('@/views/user-profile/index.vue'),
    meta: { title: '用户主页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', guest: true }
  },
  {
    path: '/login-password',
    name: 'LoginPassword',
    component: () => import('@/views/login/password.vue'),
    meta: { title: '密码登录', guest: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  // 需要登录的页面
  if (to.meta.requireAuth && !userStore.token) {
    next('/login')
    return
  }
  
  // 已登录用户不能访问登录页
  if (to.meta.guest && userStore.token) {
    next('/')
    return
  }
  
  next()
})

export default router
