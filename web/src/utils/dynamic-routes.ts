import type { RouteRecordRaw } from 'vue-router'
import type { MenuInfo } from '../api/auth'

const viewModules = import.meta.glob('../views/**/*.vue')

function resolveComponent(component?: string) {
  if (!component) {
    return undefined
  }
  const viewPath = `../views/${component}.vue`
  return viewModules[viewPath]
}

export function buildRoutesFromMenus(menus: MenuInfo[]): RouteRecordRaw[] {
  return menus
    .filter(menu => menu.visible)
    .map((menu) => {
      const children = menu.children ? buildRoutesFromMenus(menu.children) : []
      const route: RouteRecordRaw = {
        path: menu.path,
        name: menu.code || menu.name,
        component: resolveComponent(menu.component),
        meta: {
          title: menu.name,
          icon: menu.icon,
          permission: menu.permission
        },
        children
      }

      if (!route.component && children.length > 0) {
        route.redirect = children[0].path
      }

      return route
    })
}

export function hasRouteComponent(route: RouteRecordRaw): boolean {
  return !!route.component
}