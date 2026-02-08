import apiClient from './client'

export type AgentSessionStatus = 'running' | 'succeeded' | 'failed' | 'cancelled'

export interface AgentSession {
  id: string
  task_id: string
  status: AgentSessionStatus
  budget?: Record<string, any>
  step_count: number
  started_at: string
  ended_at?: string
  created_at: string
  updated_at: string
}

export interface AgentSessionListResponse {
  items: AgentSession[]
  total: number
}

export interface AgentSessionListQuery {
  task_id?: string
  status?: AgentSessionStatus
  page?: number
  page_size?: number
}

export interface AgentSessionCreateReq {
  task_id: string
  budget?: Record<string, any>
}

export interface AgentSessionStopReq {
  status?: AgentSessionStatus
}

export interface AgentSessionRunReq {
  max_actions?: number
}

export interface AgentSessionEvent {
  id: string
  task_id: string
  session_id?: string
  event_type: string
  source: string
  node_key?: string
  tool_name?: string
  payload?: Record<string, any>
  created_at: string
}

export interface AgentSessionEventListResponse {
  items: AgentSessionEvent[]
  total: number
}

function toListParams(params?: { page?: number; page_size?: number; [key: string]: any }) {
  if (!params) return undefined
  const { page, page_size, ...rest } = params
  const normalized: Record<string, any> = { ...rest }
  if (typeof page_size === 'number' && page_size > 0) {
    normalized.limit = page_size
  }
  if (typeof page === 'number' && page > 0 && typeof page_size === 'number' && page_size > 0) {
    normalized.offset = (page - 1) * page_size
  }
  return normalized
}

export const agentApi = {
  listSessions(params?: AgentSessionListQuery) {
    return apiClient.get<AgentSessionListResponse>('/agent/sessions', { params: toListParams(params) })
  },

  createSession(data: AgentSessionCreateReq) {
    return apiClient.post<AgentSession>('/agent/sessions', data)
  },

  getSession(id: string) {
    return apiClient.get<AgentSession>(`/agent/sessions/${id}`)
  },

  listSessionEvents(id: string, params?: { source?: string; node_key?: string; page?: number; page_size?: number }) {
    return apiClient.get<AgentSessionEventListResponse>(`/agent/sessions/${id}/events`, { params: toListParams(params) })
  },

  stopSession(id: string, status: AgentSessionStatus = 'cancelled') {
    return apiClient.post<AgentSession>(`/agent/sessions/${id}/stop`, { status } as AgentSessionStopReq)
  },

  runSession(id: string, data?: AgentSessionRunReq) {
    return apiClient.post<AgentSession>(`/agent/sessions/${id}/run`, data ?? {})
  },
}
