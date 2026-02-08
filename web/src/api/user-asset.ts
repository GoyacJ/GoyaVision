import apiClient from './client'

export interface UserAssetSummary {
  points: number
  balance: number
  subscription: string
  member_level: string
}

export interface TransactionRecord {
  id: string
  type: string
  method: string
  amount: number
  status: string
  created_at: string
}

export interface PointRecord {
  id: number
  type: string
  change: number
  balance: number
  created_at: string
}

export interface UsageStats {
  operator_calls: number
  ai_model_calls: number
  token_usage: number
}

export const userAssetApi = {
  getSummary() {
    return apiClient.get<UserAssetSummary>('/user/assets/summary')
  },
  getTransactions(params?: { limit?: number; offset?: number }) {
    return apiClient.get<TransactionRecord[]>('/user/assets/transactions', { params })
  },
  getPoints(params?: { limit?: number; offset?: number }) {
    return apiClient.get<PointRecord[]>('/user/assets/points', { params })
  },
  getUsage() {
    return apiClient.get<UsageStats>('/user/assets/usage')
  },
  recharge(data: { amount: number; channel: string }) {
    return apiClient.post<{ order_no: string; pay_url?: string; qrcode?: string }>('/user/assets/recharge', data)
  },
  checkIn() {
    return apiClient.post<{ message: string }>('/user/assets/checkin')
  },
  subscribe(planName: string) {
    return apiClient.post<{ message: string }>('/user/assets/subscribe', { plan_name: planName })
  }
}
