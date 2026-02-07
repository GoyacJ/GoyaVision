import apiClient from './client'

export type OperatorCategory = 'analysis' | 'processing' | 'generation' | 'utility'
export type OperatorStatus = 'draft' | 'testing' | 'published' | 'deprecated'
export type OperatorOrigin = 'builtin' | 'custom' | 'marketplace' | 'mcp'
export type OperatorExecMode = 'http' | 'cli' | 'mcp' | 'ai_model'

export interface OperatorVersion {
  id: string
  version: string
  exec_mode: OperatorExecMode
  exec_config?: Record<string, any>
  status: 'draft' | 'active' | 'archived'
}

export interface Operator {
  id: string
  code: string
  name: string
  description?: string
  category: OperatorCategory
  type: string
  origin?: OperatorOrigin
  active_version_id?: string
  exec_mode?: OperatorExecMode
  active_version?: OperatorVersion
  version: string
  endpoint: string
  method: string
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  status: OperatorStatus
  is_builtin: boolean
  tags?: string[]
  visibility?: number
  visible_role_ids?: string[]
  created_at: string
  updated_at: string
}

export interface OperatorListQuery {
  category?: OperatorCategory
  type?: string
  status?: OperatorStatus
  origin?: OperatorOrigin
  exec_mode?: OperatorExecMode
  is_builtin?: boolean
  tags?: string[]
  keyword?: string
  page?: number
  page_size?: number
}

export interface OperatorCreateReq {
  code: string
  name: string
  description?: string
  category: OperatorCategory
  type: string
  origin?: OperatorOrigin
  exec_mode?: OperatorExecMode
  exec_config?: Record<string, any>
  version?: string
  endpoint?: string
  method?: string
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  status?: OperatorStatus
  is_builtin?: boolean
  tags?: string[]
  visibility?: number
  visible_role_ids?: string[]
}

export interface OperatorUpdateReq {
  name?: string
  description?: string
  category?: OperatorCategory
  tags?: string[]
  visibility?: number
  visible_role_ids?: string[]
}

export interface OperatorListResponse {
  items: Operator[]
  total: number
}

export interface OperatorVersionListQuery {
  page?: number
  page_size?: number
}

export interface OperatorVersionListResponse {
  items: OperatorVersion[]
  total: number
}

export interface TestOperatorReq {
  asset_id?: string
  params?: Record<string, any>
}

export interface TestOperatorResponse {
  success: boolean
  message: string
  diagnostics?: Record<string, any>
}

export interface OperatorVersionCreateReq {
  version: string
  exec_mode: OperatorExecMode
  exec_config?: Record<string, any>
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  changelog?: string
  status?: string
}

export interface OperatorVersionActionReq {
  version_id: string
}

export interface ValidateSchemaReq {
  schema: Record<string, any>
}

export interface ValidateConnectionReq {
  upstream_output_spec: Record<string, any>
  downstream_input_schema: Record<string, any>
}

export interface ValidateResultResponse {
  valid: boolean
  message?: string
}

export interface OperatorTemplate {
  id: string
  code: string
  name: string
  description?: string
  category: OperatorCategory
  type: string
  exec_mode: OperatorExecMode
  exec_config?: Record<string, any>
  input_schema?: Record<string, any>
  output_spec?: Record<string, any>
  config?: Record<string, any>
  author?: string
  tags?: string[]
  icon_url?: string
  downloads: number
  created_at: string
  updated_at: string
}

export interface TemplateListQuery {
  category?: OperatorCategory
  type?: string
  exec_mode?: OperatorExecMode
  tags?: string[]
  keyword?: string
  page?: number
  page_size?: number
}

export interface OperatorTemplateListResponse {
  items: OperatorTemplate[]
  total: number
}

export interface InstallTemplateReq {
  template_id: string
  operator_code: string
  operator_name: string
  tags?: string[]
}

export interface OperatorDependency {
  id: string
  operator_id: string
  depends_on_id: string
  min_version?: string
  is_optional: boolean
  created_at: string
}

export interface SetDependenciesReq {
  dependencies: Array<{
    depends_on_id: string
    min_version?: string
    is_optional?: boolean
  }>
}

export interface DependencyCheckResponse {
  satisfied: boolean
  unmet?: string[]
}

export interface MCPServer {
  id: string
  name: string
  description?: string
  status?: string
}

export interface MCPTool {
  name: string
  description?: string
  version?: string
  input_schema?: Record<string, any>
  output_schema?: Record<string, any>
}

export interface MCPInstallReq {
  server_id: string
  tool_name: string
  operator_code: string
  operator_name: string
  category?: OperatorCategory
  type?: string
  timeout_sec?: number
  tags?: string[]
}

export interface SyncMCPTemplatesReq {
  server_id: string
}

export interface SyncMCPTemplatesResponse {
  server_id: string
  total: number
  created: number
  updated: number
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

export const operatorApi = {
  list(params?: OperatorListQuery) {
    return apiClient.get<OperatorListResponse>('/operators', { params: toListParams(params) })
  },

  get(id: string) {
    return apiClient.get<Operator>(`/operators/${id}`)
  },

  create(data: OperatorCreateReq) {
    return apiClient.post<Operator>('/operators', data)
  },

  update(id: string, data: OperatorUpdateReq) {
    return apiClient.put<Operator>(`/operators/${id}`, data)
  },

  delete(id: string) {
    return apiClient.delete(`/operators/${id}`)
  },

  publish(id: string) {
    return apiClient.post<Operator>(`/operators/${id}/publish`)
  },

  deprecate(id: string) {
    return apiClient.post<Operator>(`/operators/${id}/deprecate`)
  },

  test(id: string, data?: TestOperatorReq) {
    return apiClient.post<TestOperatorResponse>(`/operators/${id}/test`, data ?? {})
  },

  listVersions(id: string, params?: OperatorVersionListQuery) {
    return apiClient.get<OperatorVersionListResponse>(`/operators/${id}/versions`, { params: toListParams(params) })
  },

  getVersion(id: string, versionId: string) {
    return apiClient.get<OperatorVersion>(`/operators/${id}/versions/${versionId}`)
  },

  createVersion(id: string, data: OperatorVersionCreateReq) {
    return apiClient.post<OperatorVersion>(`/operators/${id}/versions`, data)
  },

  activateVersion(id: string, data: OperatorVersionActionReq) {
    return apiClient.post<Operator>(`/operators/${id}/versions/activate`, data)
  },

  rollbackVersion(id: string, data: OperatorVersionActionReq) {
    return apiClient.post<Operator>(`/operators/${id}/versions/rollback`, data)
  },

  archiveVersion(id: string, data: OperatorVersionActionReq) {
    return apiClient.post<OperatorVersion>(`/operators/${id}/versions/archive`, data)
  },

  validateSchema(data: ValidateSchemaReq) {
    return apiClient.post<ValidateResultResponse>('/operators/validate-schema', data)
  },

  validateConnection(data: ValidateConnectionReq) {
    return apiClient.post<ValidateResultResponse>('/operators/validate-connection', data)
  },

  listTemplates(params?: TemplateListQuery) {
    return apiClient.get<OperatorTemplateListResponse>('/operators/templates', { params: toListParams(params) })
  },

  getTemplate(templateId: string) {
    return apiClient.get<OperatorTemplate>(`/operators/templates/${templateId}`)
  },

  installTemplate(data: InstallTemplateReq) {
    return apiClient.post<Operator>('/operators/templates/install', data)
  },

  listDependencies(id: string) {
    return apiClient.get<OperatorDependency[]>(`/operators/${id}/dependencies`)
  },

  setDependencies(id: string, data: SetDependenciesReq) {
    return apiClient.put(`/operators/${id}/dependencies`, data)
  },

  checkDependencies(id: string) {
    return apiClient.get<DependencyCheckResponse>(`/operators/${id}/dependencies/check`)
  },

  listMCPServers() {
    return apiClient.get<MCPServer[]>('/operators/mcp/servers')
  },

  listMCPTools(serverId: string) {
    return apiClient.get<MCPTool[]>(`/operators/mcp/servers/${serverId}/tools`)
  },

  previewMCPTool(serverId: string, toolName: string) {
    return apiClient.get<MCPTool>(`/operators/mcp/servers/${serverId}/tools/${encodeURIComponent(toolName)}/preview`)
  },

  installMCPOperator(data: MCPInstallReq) {
    return apiClient.post<Operator>('/operators/mcp/install', data)
  },

  syncMCPTemplates(data: SyncMCPTemplatesReq) {
    return apiClient.post<SyncMCPTemplatesResponse>('/operators/mcp/sync-templates', data)
  },

  listByCategory(category: string, params?: OperatorListQuery) {
    return apiClient.get<Operator[]>(`/operators/category/${category}`, { params: toListParams(params) })
  }
}
