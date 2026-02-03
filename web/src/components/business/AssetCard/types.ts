/**
 * AssetCard 业务组件类型定义
 * 资产卡片组件
 */

/**
 * 资产数据
 */
export interface Asset {
  id: string | number
  name: string
  type: string
  status: string
  thumbnail?: string
  description?: string
  size?: string
  duration?: string
  createdAt?: string
  updatedAt?: string
  [key: string]: any
}

/**
 * AssetCard Props
 */
export interface AssetCardProps {
  /**
   * 资产数据
   */
  asset: Asset
  
  /**
   * 是否可选择
   * @default false
   */
  selectable?: boolean
  
  /**
   * 是否选中
   * @default false
   */
  selected?: boolean
  
  /**
   * 是否显示操作按钮
   * @default true
   */
  showActions?: boolean
  
  /**
   * 是否显示详情按钮
   * @default true
   */
  showDetail?: boolean
  
  /**
   * 是否显示编辑按钮
   * @default true
   */
  showEdit?: boolean
  
  /**
   * 是否显示删除按钮
   * @default true
   */
  showDelete?: boolean
}

/**
 * AssetCard Emits
 */
export interface AssetCardEmits {
  /**
   * 点击卡片
   */
  (e: 'click', asset: Asset): void
  
  /**
   * 选择状态变化
   */
  (e: 'select', selected: boolean, asset: Asset): void
  
  /**
   * 点击详情
   */
  (e: 'detail', asset: Asset): void
  
  /**
   * 点击编辑
   */
  (e: 'edit', asset: Asset): void
  
  /**
   * 点击删除
   */
  (e: 'delete', asset: Asset): void
}
