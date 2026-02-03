import type { MediaAsset } from '@/api/asset'

export interface AssetCardProps {
  /** 资产数据 */
  asset: MediaAsset
  /** 是否可选择 */
  selectable?: boolean
  /** 是否选中 */
  selected?: boolean
}

export interface AssetCardEmits {
  /** 点击事件 */
  (e: 'click', asset: MediaAsset): void
  /** 查看事件 */
  (e: 'view', asset: MediaAsset): void
  /** 编辑事件 */
  (e: 'edit', asset: MediaAsset): void
  /** 删除事件 */
  (e: 'delete', asset: MediaAsset): void
  /** 选择状态变化 */
  (e: 'select', asset: MediaAsset, selected: boolean): void
}
