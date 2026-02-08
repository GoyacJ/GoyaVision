import { useVueFlow } from '@vue-flow/core'

export function useConnectionValidation() {
  const { onConnect, addEdges } = useVueFlow({ id: 'workflow-editor' })

  onConnect((params) => {
    // Simplified logic: just add the edge
    addEdges([
      {
        ...params,
        id: `e-${params.source}-${params.target}-${Math.random().toString(36).slice(2, 6)}`,
        type: 'conditional',
        data: {
          condition: { type: 'always' }
        }
      }
    ])
  })

  return {}
}
