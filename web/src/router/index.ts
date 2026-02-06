import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const Layout = () => import('../layout/index.vue')
const Login = () => import('../views/login/index.vue')
const ComponentDemo = () => import('../views/ComponentDemo.vue')
const StateDemo = () => import('../views/StateDemo.vue')
const OperatorMarketplace = () => import('../views/operator-marketplace/index.vue')

export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登录', hidden: true }
  },
  {
    path: '/component-demo',
    name: 'ComponentDemo',
    component: ComponentDemo,
    meta: { title: '组件展示', hidden: true }
  },
  {
    path: '/state-demo',
    name: 'StateDemo',
    component: StateDemo,
    meta: { title: '状态组件演示', hidden: true }
  },
  {
    path: '/',
    name: 'Root',
    component: Layout,
    children: [
      {
        path: '/operator-marketplace',
        name: 'OperatorMarketplace',
        component: OperatorMarketplace,
        meta: { title: '算子模板市场', hidden: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes
})

export default router
