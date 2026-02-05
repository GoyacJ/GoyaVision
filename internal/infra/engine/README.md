# DAG Workflow Engine

A complete implementation of a Directed Acyclic Graph (DAG) workflow engine for parallel and sequential task execution.

## Features

### 1. Topological Sorting (Kahn's Algorithm)
- Determines correct execution order of nodes based on dependencies
- Ensures all dependencies are executed before dependent nodes
- Validates workflow structure before execution

### 2. Cycle Detection
- Automatically detects cycles in workflow definitions
- Prevents infinite loops by rejecting cyclic workflows
- Provides clear error messages for debugging

### 3. Parallel Execution
- Groups nodes into execution layers
- Nodes within the same layer execute concurrently
- Maximizes throughput for independent operations
- Thread-safe execution with proper synchronization

### 4. Data Flow
- Captures output from each node execution
- Makes upstream node results available to downstream nodes
- Stores results under predictable keys: `{nodeKey}_output`, `{nodeKey}_assets`, etc.
- Supports complex multi-node data processing pipelines

### 5. Node Configuration
- **Retry Logic**: Automatic retry with exponential backoff
- **Timeout Control**: Per-node timeout configuration
- **Parameter Override**: Node config params override task params

### 6. Progress Tracking
- Layer-based progress calculation
- Real-time progress updates in database
- Current node tracking for debugging
- Thread-safe progress retrieval

### 7. Artifact Management
- Automatic artifact creation for each node output
- Supports multiple artifact types: Asset, Result, Timeline, Report
- Node key tagging for traceability
- Separate artifacts per node for granular analysis

## Architecture

```
DAGWorkflowEngine
├── Execute()           - Main execution entry point
├── Cancel()            - Graceful cancellation support
├── GetProgress()       - Real-time progress monitoring
│
├── Algorithm Layer
│   ├── topologicalSort()      - Kahn's algorithm for ordering
│   └── buildExecutionLayers() - Layer construction for parallelism
│
├── Execution Layer
│   ├── executeLayer()  - Parallel layer execution
│   ├── executeNode()   - Single node execution with retry
│   └── prepareNodeInput() - Input preparation with data flow
│
└── Persistence Layer
    ├── saveArtifacts()       - Artifact persistence
    ├── updateTaskProgress()  - Progress updates
    └── updateTaskStatus()    - Task lifecycle management
```

## Usage Example

### Simple Linear Workflow
```
A → B → C
```
Execution: A, then B, then C (sequential)

### Diamond Pattern (Parallel Branches)
```
    A
   / \
  B   C
   \ /
    D
```
Execution Layers:
- Layer 1: A
- Layer 2: B, C (parallel)
- Layer 3: D

### Complex Multi-Branch
```
      start
     /  |  \
    p1  p2  p3
     \  |  /
      join
       |
      end
```
Execution Layers:
- Layer 1: start
- Layer 2: p1, p2, p3 (parallel)
- Layer 3: join
- Layer 4: end

## Node Configuration

```json
{
  "node_key": "process_video",
  "operator_id": "uuid",
  "config": {
    "params": {
      "resolution": "1080p",
      "codec": "h264"
    },
    "retry_count": 3,
    "timeout_seconds": 300
  }
}
```

## Data Flow Example

Given workflow: `extract_frames → detect_objects → generate_report`

**extract_frames output:**
```json
{
  "output_assets": [{"type": "image", "path": "/frames/frame_001.jpg"}],
  "results": [{"type": "frame_info", "data": {"count": 100}}]
}
```

**detect_objects receives:**
```json
{
  "params": {
    "extract_frames_output": {...},
    "extract_frames_assets": [...],
    "extract_frames_results": [...]
  }
}
```

## Error Handling

### Retry with Exponential Backoff
```
Attempt 1: Execute
Attempt 2: Wait 1s, Execute
Attempt 3: Wait 2s, Execute
Attempt 4: Wait 4s, Execute
```

### Timeout Protection
Each node can specify `timeout_seconds` in its config. If execution exceeds timeout, the context is cancelled and the node fails.

### Cancellation
Workflows can be cancelled mid-execution via `Cancel(taskID)`. All running nodes receive cancellation signals through context propagation.

## Performance Characteristics

- **Memory**: O(V + E) where V = nodes, E = edges
- **Topological Sort**: O(V + E) using Kahn's algorithm
- **Layer Building**: O(V + E)
- **Execution**: O(L) where L = number of layers (sequential), nodes within layer execute in parallel

## Thread Safety

- Task execution map protected by RWMutex
- Per-task node results protected by dedicated mutex
- Safe concurrent access to progress and results
- Context cancellation propagates safely

## Testing

Comprehensive test suite covers:
- ✅ Topological sort correctness
- ✅ Cycle detection
- ✅ Parallel execution layer building
- ✅ Diamond pattern execution
- ✅ Wide parallelism (many parallel branches)
- ✅ Retry logic with failure recovery
- ✅ Timeout handling
- ✅ Cancellation propagation
- ✅ Progress tracking
- ✅ Data flow between nodes

Run tests:
```bash
go test ./internal/infra/engine/... -v
```

## Migration from SimpleWorkflowEngine

The DAG engine is a drop-in replacement:

**Before:**
```go
executor := engine.NewHTTPOperatorExecutor()
workflowEngine := engine.NewSimpleWorkflowEngine(repo, executor)
```

**After:**
```go
executor := engine.NewHTTPOperatorExecutor()
workflowEngine := infraengine.NewDAGWorkflowEngine(uow, executor)
```

Key differences:
- Uses UnitOfWork instead of Repository directly (better transaction control)
- Supports parallel execution (SimpleEngine was sequential only)
- Proper cycle detection (SimpleEngine didn't validate DAG)
- Enhanced retry and timeout logic
- Better data flow with upstream node outputs

## Limitations

- No conditional edge execution yet (planned)
- No dynamic branching (all branches defined at workflow creation time)
- No partial workflow restart (must restart entire workflow)

## Future Enhancements

- [ ] Conditional edge execution based on output values
- [ ] Dynamic subworkflow spawning
- [ ] Checkpoint/resume for long-running workflows
- [ ] Node-level parallelism limits
- [ ] Priority-based execution ordering
- [ ] Resource allocation and limits per node
