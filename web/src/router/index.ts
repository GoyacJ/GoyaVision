import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const Layout = () => import('../layout/index.vue')
const Login = () => import('../views/login/index.vue')
const AssetList = () => import('../views/asset/index.vue')
const OperatorList = () => import('../views/operator/index.vue')
const WorkflowList = () => import('../views/workflow/index.vue')
const TaskList = () => import('../views/task/index.vue')
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
    redirect: '/assets',
    children: [
      {
        path: 'assets',
        name: 'AssetList',
        component: AssetList,
        meta: { title: '媒体资产', icon: 'Film' }
      },
      {
        path: 'operators',
        name: 'OperatorList',
        component: OperatorList,
        meta: { title: '算子中心', icon: 'Grid' }
      },
      {
        path: 'workflows',
        name: 'WorkflowList',
        component: WorkflowList,
        meta: { title: '工作流', icon: 'Connection' }
      },
      {
        path: 'tasks',
        name: 'TaskList',
        component: TaskList,
        meta: { title: '任务中心', icon: 'List' }
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
