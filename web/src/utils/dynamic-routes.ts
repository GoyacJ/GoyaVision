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
      const component = resolveComponent(menu.component)
      const route: RouteRecordRaw = {
        path: menu.path,
        name: menu.code || menu.name,
        component: component,
        meta: {
          title: menu.name,
          icon: menu.icon,
          permission: menu.permission
        },
        children: children.length > 0 ? children : undefined
      }

      if (!route.component && children.length > 0) {
        route.redirect = children[0].path
      }

      if (!route.component && (!children || children.length === 0)) {
        console.warn(`[Router] Menu "${menu.name}" (${menu.path}) has no component and no children, skipping route registration`)
      }

      return route
    })
}

export function hasRouteComponent(route: RouteRecordRaw): boolean {
  return !!route.component
}