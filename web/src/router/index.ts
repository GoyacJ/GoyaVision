import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const Layout = () => import('../layout/index.vue')
const Login = () => import('../views/login/index.vue')
const StreamList = () => import('../views/stream/index.vue')
const AlgorithmList = () => import('../views/algorithm/index.vue')
const InferenceResultList = () => import('../views/inference/index.vue')
const UserList = () => import('../views/system/user/index.vue')
const RoleList = () => import('../views/system/role/index.vue')
const MenuList = () => import('../views/system/menu/index.vue')

export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登录', hidden: true }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/streams',
    children: [
      {
        path: 'streams',
        name: 'StreamList',
        component: StreamList,
        meta: { title: '视频流管理', icon: 'Monitor' }
      },
      {
        path: 'algorithms',
        name: 'AlgorithmList',
        component: AlgorithmList,
        meta: { title: '算法管理', icon: 'Cpu' }
      },
      {
        path: 'inference-results',
        name: 'InferenceResultList',
        component: InferenceResultList,
        meta: { title: '推理结果', icon: 'DataAnalysis' }
      },
      {
        path: 'system/user',
        name: 'UserList',
        component: UserList,
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: 'system/role',
        name: 'RoleList',
        component: RoleList,
        meta: { title: '角色管理', icon: 'UserFilled' }
      },
      {
        path: 'system/menu',
        name: 'MenuList',
        component: MenuList,
        meta: { title: '菜单管理', icon: 'Menu' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes
})

export default router
