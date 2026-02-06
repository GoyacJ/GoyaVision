import type { MediaAsset } from '@/api/asset'

export interface AssetCardProps {
  /** 资产数据 */
  asset: MediaAsset
  /** 是否可选择 */
  selectable?: boolean
  /** 是否选中 */
  selected?: boolean
  /** 是否允许编辑 */
  canEdit?: boolean
}

export interface AssetCardEmits {
  /** 点击事件 */
  (e: 'click', asset: MediaAsset): void
  /** 详情事件 */
  (e: 'detail', asset: MediaAsset): void
  /** 删除事件 */
  (e: 'delete', asset: MediaAsset): void
  /** 选择状态变化 */
  (e: 'select', asset: MediaAsset, selected: boolean): void
}
