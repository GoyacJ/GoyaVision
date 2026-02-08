import { useVueFlow, Position } from '@vue-flow/core'
import dagre from 'dagre'

export function useAutoLayout() {
  const { nodes, edges, fitView } = useVueFlow({ id: 'workflow-editor' })

  function layout(direction = 'TB') {
    const dagreGraph = new dagre.graphlib.Graph()
    dagreGraph.setDefaultEdgeLabel(() => ({}))

    const isHorizontal = direction === 'LR'
    dagreGraph.setGraph({ rankdir: direction })

    nodes.value.forEach((node) => {
      dagreGraph.setNode(node.id, { width: 180, height: 80 })
    })

    edges.value.forEach((edge) => {
      dagreGraph.setEdge(edge.source, edge.target)
    })

    dagre.layout(dagreGraph)

    nodes.value = nodes.value.map((node) => {
      const nodeWithPosition = dagreGraph.node(node.id)
      return {
        ...node,
        targetPosition: isHorizontal ? Position.Left : Position.Top,
        sourcePosition: isHorizontal ? Position.Right : Position.Bottom,
        position: {
          x: nodeWithPosition.x - 90,
          y: nodeWithPosition.y - 40,
        },
      }
    })

    setTimeout(() => {
      fitView()
    }, 100)
  }

  return { layout }
}
