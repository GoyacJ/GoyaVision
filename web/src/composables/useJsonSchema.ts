import { ref } from 'vue'
import { operatorApi } from '@/api/operator'

export function useJsonSchema() {
  const validating = ref(false)

  function parseJsonObject(text: string): { valid: boolean; data?: Record<string, any>; message?: string } {
    try {
      const parsed = JSON.parse(text || '{}')
      if (parsed === null || Array.isArray(parsed) || typeof parsed !== 'object') {
        return { valid: false, message: 'JSON 必须是对象类型' }
      }
      return { valid: true, data: parsed }
    } catch {
      return { valid: false, message: 'JSON 格式不合法，请检查后重试' }
    }
  }

  async function validateSchema(schema: Record<string, any>): Promise<{ valid: boolean; message?: string }> {
    validating.value = true
    try {
      const res = await operatorApi.validateSchema({ schema })
      return { valid: !!res.data?.valid, message: res.data?.message }
    } catch (error: any) {
      return {
        valid: false,
        message: error?.response?.data?.message || 'Schema 校验失败'
      }
    } finally {
      validating.value = false
    }
  }

  return {
    validating,
    parseJsonObject,
    validateSchema
  }
}
