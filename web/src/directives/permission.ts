import type { Directive, DirectiveBinding } from 'vue'
import { useUserStore } from '../store/user'

export const permissionDirective: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    const { value } = binding

    if (value) {
      const permissions = Array.isArray(value) ? value : [value]
      const hasPermission = userStore.hasAnyPermission(permissions)

      if (!hasPermission) {
        el.parentNode?.removeChild(el)
      }
    }
  }
}

export function checkPermission(permission: string | string[]): boolean {
  const userStore = useUserStore()
  const permissions = Array.isArray(permission) ? permission : [permission]
  return userStore.hasAnyPermission(permissions)
}
