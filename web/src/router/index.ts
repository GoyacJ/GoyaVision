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
        path: '/assets',
        name: 'Assets',
        component: () => import('../views/asset/index.vue'),
        meta: { title: '媒体资产库' }
      },
      {
        path: '/operator-marketplace',
        name: 'OperatorMarketplace',
        component: OperatorMarketplace,
        meta: { title: '算子模板市场', hidden: true }
      },
      {
        path: '/workflows',
        name: 'Workflows',
        component: () => import('../views/workflow/index.vue'),
        meta: { title: '工作流管理' }
      },
      {
        path: '/workflows/:id/edit',
        name: 'WorkflowEditor',
        component: () => import('../views/workflow/editor/index.vue'),
        meta: { title: '编辑工作流', hidden: true }
      },
      {
        path: '/tasks',
        name: 'Tasks',
        component: () => import('../views/task/index.vue'),
        meta: { title: '任务管理' }
      },
      {
        path: '/tasks/:id',
        name: 'TaskDetail',
        component: () => import('../views/task/detail.vue'),
        meta: { title: '任务详情', hidden: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes
})

export default router
