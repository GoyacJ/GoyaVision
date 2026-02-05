# DAG Workflow Engine - Execution Flow

## High-Level Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    DAGWorkflowEngine.Execute()                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │ Build Execution  │
                    │     Layers       │
                    │  (Kahn's Algo)   │
                    └──────────────────┘
                              │
                              ▼
                   ┌────────────────────┐
                   │  Cycle Detection   │
                   │ (Check sorted len) │
                   └────────────────────┘
                              │
                ┌─────────────┴──────────────┐
                │                            │
                ▼                            ▼
         ┌──────────┐                  ┌──────────┐
         │  Error   │                  │   OK     │
         │  Return  │                  │Continue  │
         └──────────┘                  └──────────┘
                                             │
                                             ▼
                                   ┌──────────────────┐
                                   │ Create Node Map  │
                                   │& Execution State │
                                   └──────────────────┘
                                             │
                                             ▼
                                   ┌──────────────────┐
                                   │  Update Task to  │
                                   │     RUNNING      │
                                   └──────────────────┘
                                             │
                                             ▼
                             ╔═══════════════════════════╗
                             ║  FOR EACH LAYER (i=1..L)  ║
                             ╚═══════════════════════════╝
                                             │
                              ┌──────────────┴───────────────┐
                              │                              │
                              ▼                              ▼
                    ┌──────────────────┐         ┌────────────────────┐
                    │ Check Context    │         │   Execute Layer    │
                    │  Cancellation    │────────▶│   (Parallel)       │
                    └──────────────────┘         └────────────────────┘
                              │                              │
                              │ Cancelled                    │ Success
                              ▼                              ▼
                    ┌──────────────────┐         ┌────────────────────┐
                    │  Update Task to  │         │  Update Progress   │
                    │   CANCELLED      │         │    (i/L * 100)     │
                    └──────────────────┘         └────────────────────┘
                                                              │
                                                              ▼
                                                       ┌──────────────┐
                                                       │   Next Layer │
                                                       └──────────────┘
                                                              │
                              ┌───────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  All Layers Done │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Update Task to  │
                    │     SUCCESS      │
                    └──────────────────┘
```

## Layer Execution Detail

```
┌─────────────────────────────────────────────────────────────────┐
│                    executeLayer(layer, nodeMap)                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Check Layer     │
                    │  Size            │
                    └──────────────────┘
                              │
                ┌─────────────┴──────────────┐
                │                            │
                ▼                            ▼
         ┌──────────┐                ┌──────────────┐
         │ Size = 1 │                │  Size > 1    │
         └──────────┘                └──────────────┘
                │                            │
                ▼                            ▼
    ┌────────────────────┐      ┌─────────────────────────┐
    │  Execute Node      │      │  Create WaitGroup       │
    │  Synchronously     │      │  & Error Channel        │
    └────────────────────┘      └─────────────────────────┘
                                            │
                                            ▼
                              ╔═══════════════════════════╗
                              ║  FOR EACH NODE IN LAYER   ║
                              ╚═══════════════════════════╝
                                            │
                                            ▼
                              ┌─────────────────────────────┐
                              │  Launch Goroutine:          │
                              │  - WaitGroup.Add(1)         │
                              │  - Execute Node             │
                              │  - WaitGroup.Done()         │
                              │  - Send errors to channel   │
                              └─────────────────────────────┘
                                            │
                                            ▼
                              ┌─────────────────────────────┐
                              │  WaitGroup.Wait()           │
                              │  (Wait for all to complete) │
                              └─────────────────────────────┘
                                            │
                                            ▼
                              ┌─────────────────────────────┐
                              │  Close Error Channel        │
                              │  & Check for Errors         │
                              └─────────────────────────────┘
                                            │
                              ┌─────────────┴──────────────┐
                              │                            │
                              ▼                            ▼
                       ┌──────────┐                 ┌──────────┐
                       │  Errors  │                 │   OK     │
                       │  Return  │                 │  Return  │
                       └──────────┘                 └──────────┘
```

## Node Execution Detail

```
┌─────────────────────────────────────────────────────────────────┐
│                    executeNode(node, task, exec)                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Update Current  │
                    │  Node in State   │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Skip if no      │
                    │  Operator ID     │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Get Operator    │
                    │  from Repository │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Prepare Input:  │
                    │  - Task params   │
                    │  - Node config   │
                    │  - Upstream outs │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Apply Timeout   │
                    │  if configured   │
                    └──────────────────┘
                              │
                              ▼
              ╔═══════════════════════════════════╗
              ║  RETRY LOOP (1 + retry_count)    ║
              ╚═══════════════════════════════════╝
                              │
                              ▼
                    ┌──────────────────┐
                    │  Execute         │
                    │  Operator        │
                    └──────────────────┘
                              │
                ┌─────────────┴──────────────┐
                │                            │
                ▼                            ▼
         ┌──────────┐                ┌──────────────┐
         │  Success │                │    Error     │
         └──────────┘                └──────────────┘
                │                            │
                │                            ▼
                │                  ┌──────────────────┐
                │                  │  More Retries?   │
                │                  └──────────────────┘
                │                            │
                │                  ┌─────────┴─────────┐
                │                  │                   │
                │                  ▼                   ▼
                │           ┌──────────┐       ┌──────────┐
                │           │   Wait   │       │  Return  │
                │           │ Backoff  │       │  Error   │
                │           └──────────┘       └──────────┘
                │                  │
                │                  ▼
                │           ┌──────────┐
                │           │  Retry   │
                │           └──────────┘
                │                  │
                └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Store Output    │
                    │  in Exec State   │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Save Artifacts  │
                    │  to Database     │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │     Return       │
                    │     Success      │
                    └──────────────────┘
```

## Topological Sort (Kahn's Algorithm)

```
┌─────────────────────────────────────────────────────────────────┐
│                      topologicalSort(nodes, edges)               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Build Graph:    │
                    │  - Adjacency map │
                    │  - In-degree map │
                    └──────────────────┘
                              │
                              ▼
              ╔═══════════════════════════════════╗
              ║  FOR EACH NODE                    ║
              ║  - graph[node] = []               ║
              ║  - inDegree[node] = 0             ║
              ╚═══════════════════════════════════╝
                              │
                              ▼
              ╔═══════════════════════════════════╗
              ║  FOR EACH EDGE (src → tgt)        ║
              ║  - graph[src].append(tgt)         ║
              ║  - inDegree[tgt]++                ║
              ╚═══════════════════════════════════╝
                              │
                              ▼
                    ┌──────────────────┐
                    │  Find all nodes  │
                    │  with inDegree=0 │
                    │  → queue          │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  sorted = []     │
                    └──────────────────┘
                              │
                              ▼
              ╔═══════════════════════════════════╗
              ║  WHILE queue not empty:           ║
              ║  1. current = queue.pop()         ║
              ║  2. sorted.append(current)        ║
              ║  3. FOR neighbor in graph[curr]:  ║
              ║     - inDegree[neighbor]--        ║
              ║     - if inDegree[neighbor] == 0: ║
              ║       queue.append(neighbor)      ║
              ╚═══════════════════════════════════╝
                              │
                              ▼
                    ┌──────────────────┐
                    │  Check:          │
                    │  len(sorted) ==  │
                    │  len(nodes)?     │
                    └──────────────────┘
                              │
                ┌─────────────┴──────────────┐
                │                            │
                ▼                            ▼
         ┌──────────┐                ┌──────────────┐
         │   NO:    │                │    YES:      │
         │  Cycle!  │                │   Return     │
         │  Error   │                │   sorted     │
         └──────────┘                └──────────────┘
```

## Data Flow Example

```
Workflow: extract_frames → detect_objects → generate_report

┌─────────────────────────────────────────────────────────────────┐
│                         Layer 1: extract_frames                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    Input: { asset_id: uuid }
                              │
                              ▼
                    ┌──────────────────┐
                    │  Execute         │
                    │  Operator        │
                    └──────────────────┘
                              │
                              ▼
                    Output: {
                      output_assets: [
                        {type: "image", path: "/frames/001.jpg"}
                      ],
                      results: [
                        {type: "frame_info", data: {count: 100}}
                      ]
                    }
                              │
                              ▼
                    ┌──────────────────┐
                    │  Store in        │
                    │  exec.nodeResults│
                    │  ["extract_..."] │
                    └──────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Layer 2: detect_objects                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    Input: {
                      asset_id: uuid,
                      extract_frames_output: {...},
                      extract_frames_assets: [...],
                      extract_frames_results: [...]
                    }
                              │
                              ▼
                    ┌──────────────────┐
                    │  Execute         │
                    │  Operator        │
                    └──────────────────┘
                              │
                              ▼
                    Output: {
                      results: [
                        {type: "detection", data: {objects: 5}}
                      ]
                    }
                              │
                              ▼
                    ┌──────────────────┐
                    │  Store in        │
                    │  exec.nodeResults│
                    │  ["detect_..."]  │
                    └──────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Layer 3: generate_report                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    Input: {
                      asset_id: uuid,
                      extract_frames_output: {...},
                      extract_frames_assets: [...],
                      extract_frames_results: [...],
                      detect_objects_output: {...},
                      detect_objects_results: [...]
                    }
                              │
                              ▼
                    ┌──────────────────┐
                    │  Execute         │
                    │  Operator        │
                    └──────────────────┘
                              │
                              ▼
                    Output: {
                      output_assets: [
                        {type: "report", path: "/reports/summary.pdf"}
                      ]
                    }
```

## Parallel Execution Example

```
Workflow Diamond Pattern:

           start
          /     \
      proc1     proc2
          \     /
           join

Layer 1: [start]
┌─────────────────────────────────────────┐
│  Execute start (single node)            │
│  - No parallelism                       │
│  - Updates progress to 33%              │
└─────────────────────────────────────────┘

Layer 2: [proc1, proc2]
┌─────────────────────────────────────────┐
│  Launch goroutine 1:                    │
│  ┌───────────────────────────────────┐  │
│  │  Execute proc1                    │  │
│  │  - Uses start_output              │  │
│  │  - Stores proc1_output            │  │
│  └───────────────────────────────────┘  │
│                                         │
│  Launch goroutine 2:                    │
│  ┌───────────────────────────────────┐  │
│  │  Execute proc2                    │  │
│  │  - Uses start_output              │  │
│  │  - Stores proc2_output            │  │
│  └───────────────────────────────────┘  │
│                                         │
│  WaitGroup.Wait()                       │
│  - Both complete before continuing      │
│  - Updates progress to 66%              │
└─────────────────────────────────────────┘

Layer 3: [join]
┌─────────────────────────────────────────┐
│  Execute join (single node)             │
│  - Uses proc1_output                    │
│  - Uses proc2_output                    │
│  - Updates progress to 100%             │
└─────────────────────────────────────────┘

Speedup: 4 nodes in 3 time units vs 4 time units (25% faster)
```

## Error Handling Flow

```
                    ┌──────────────────┐
                    │  Node Execution  │
                    │      Fails       │
                    └──────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  Retry Logic?    │
                    └──────────────────┘
                              │
                ┌─────────────┴──────────────┐
                │                            │
                ▼                            ▼
         ┌──────────┐                ┌──────────────┐
         │   Yes    │                │      No      │
         └──────────┘                └──────────────┘
                │                            │
                ▼                            ▼
    ┌────────────────────┐          ┌───────────────┐
    │  Exponential       │          │  Return Error │
    │  Backoff & Retry   │          │  to Layer     │
    └────────────────────┘          └───────────────┘
                │                            │
                ▼                            ▼
    ┌────────────────────┐          ┌───────────────┐
    │  Still Failing?    │          │  Layer Fails  │
    └────────────────────┘          └───────────────┘
                │                            │
                ▼                            ▼
    ┌────────────────────┐          ┌───────────────┐
    │  Return Error      │          │  Cancel Other │
    │  with Attempts     │          │  Goroutines   │
    └────────────────────┘          └───────────────┘
                                             │
                                             ▼
                                    ┌───────────────┐
                                    │  Update Task  │
                                    │  to FAILED    │
                                    │  with Error   │
                                    └───────────────┘
                                             │
                                             ▼
                                    ┌───────────────┐
                                    │  Return Error │
                                    │  to Caller    │
                                    └───────────────┘
```

## Concurrency Control

```
┌─────────────────────────────────────────────────────────────────┐
│                         Thread Safety                            │
└─────────────────────────────────────────────────────────────────┘

Engine Level (DAGWorkflowEngine):
┌──────────────────────────────────────┐
│  tasks map[uuid.UUID]*taskExecution  │
│  Protected by: RWMutex               │
│  - Read: GetProgress, check exists   │
│  - Write: Execute start, Cancel      │
└──────────────────────────────────────┘

Task Level (taskExecution):
┌──────────────────────────────────────┐
│  nodeResults map[string]*Output      │
│  progress int                        │
│  currentNode string                  │
│  Protected by: RWMutex               │
│  - Read: prepareNodeInput            │
│  - Write: executeNode (store output) │
└──────────────────────────────────────┘

Context Cancellation:
┌──────────────────────────────────────┐
│  ctx, cancel := context.WithCancel() │
│  - Propagates to all goroutines      │
│  - Checked before each layer         │
│  - Operator executor receives ctx    │
│  - Timeout enforcement per node      │
└──────────────────────────────────────┘
```

## Progress Calculation

```
Total Layers: L = 5
Current Layer: i

Layer 1 complete: progress = (1/5) * 100 = 20%
Layer 2 complete: progress = (2/5) * 100 = 40%
Layer 3 complete: progress = (3/5) * 100 = 60%
Layer 4 complete: progress = (4/5) * 100 = 80%
Layer 5 complete: progress = (5/5) * 100 = 100%

┌─────────────────────────────────────────────────────────────────┐
│  Progress Bar Visualization:                                    │
│                                                                 │
│  [████████████████████                                    ] 40% │
│  Layer 2/5 - Executing: detect_objects                          │
└─────────────────────────────────────────────────────────────────┘
```
