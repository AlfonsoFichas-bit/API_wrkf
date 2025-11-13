# WebSocket Implementation for Real-Time Kanban Updates

---

**Status: Implemented**

*This document originally served as the design plan for the WebSocket feature. The core backend functionality has now been implemented as of November 2025. The backend successfully broadcasts events for task creation, status updates, assignments, and deletions. This document can now be used as a reference for the existing implementation and for future extensions.*

---

## Overview
This document outlines the WebSocket implementation requirements for real-time Kanban board updates, allowing multiple users to see changes instantly when tasks are moved, updated, or modified.

## 🔄 WebSocket Events

### Connection Events

#### 1. Client Connection
```javascript
// Client connects with JWT token
const ws = new WebSocket('ws://localhost:8080/ws?token=<jwt_token>')
```

#### 2. Connection Established
```json
{
  "type": "connection_established",
  "payload": {
    "userId": 123,
    "message": "Connected to Kanban real-time updates"
  }
}
```

### Task Events

#### 1. Task Status Updated
```json
{
  "type": "task_status_updated",
  "payload": {
    "taskId": 456,
    "oldStatus": "todo",
    "newStatus": "in_progress",
    "updatedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:30:00Z"
  }
}
```

#### 2. Task Assigned
```json
{
  "type": "task_assigned",
  "payload": {
    "taskId": 456,
    "assignedTo": {
      "id": 789,
      "name": "Jane Smith",
      "email": "jane@example.com"
    },
    "assignedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:31:00Z"
  }
}
```

#### 3. Task Updated
```json
{
  "type": "task_updated",
  "payload": {
    "taskId": 456,
    "changes": {
      "title": "Updated Task Title",
      "description": "Updated description"
    },
    "updatedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:32:00Z"
  }
}
```

#### 4. Task Created
```json
{
  "type": "task_created",
  "payload": {
    "task": {
      "id": 789,
      "title": "New Task",
      "description": "Task description",
      "status": "todo",
      "userStoryId": 123,
      "sprintId": 456,
      "createdBy": {
        "id": 123,
        "name": "John Doe"
      }
    },
    "timestamp": "2023-11-05T10:33:00Z"
  }
}
```

#### 5. Task Deleted
```json
{
  "type": "task_deleted",
  "payload": {
    "taskId": 789,
    "deletedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:34:00Z"
  }
}
```

### Sprint Events

#### 1. Sprint Status Updated
```json
{
  "type": "sprint_status_updated",
  "payload": {
    "sprintId": 456,
    "oldStatus": "active",
    "newStatus": "completed",
    "updatedBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:35:00Z"
  }
}
```

#### 2. Sprint Started/Completed
```json
{
  "type": "sprint_lifecycle_changed",
  "payload": {
    "sprintId": 456,
    "action": "started", // or "completed"
    "sprint": {
      "id": 456,
      "name": "Sprint 1",
      "status": "active"
    },
    "triggeredBy": {
      "id": 123,
      "name": "John Doe"
    },
    "timestamp": "2023-11-05T10:36:00Z"
  }
}
```

### User Presence Events

#### 1. User Joined Project
```json
{
  "type": "user_joined_project",
  "payload": {
    "projectId": 123,
    "user": {
      "id": 789,
      "name": "Jane Smith"
    },
    "timestamp": "2023-11-05T10:37:00Z"
  }
}
```

#### 2. User Left Project
```json
{
  "type": "user_left_project",
  "payload": {
    "projectId": 123,
    "user": {
      "id": 789,
      "name": "Jane Smith"
    },
    "timestamp": "2023-11-05T10:38:00Z"
  }
}
```

## 🏗️ Backend WebSocket Implementation

### WebSocket Manager Structure
```go
// websocket/manager.go
type WebSocketManager struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.RWMutex
}

type Client struct {
    hub      *WebSocketManager
    conn     *websocket.Conn
    send     chan []byte
    userID   uint
    projects map[uint]bool // Projects this user is subscribed to
}

type Message struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}
```

### Event Broadcasting Functions
```go
// websocket/events.go
func (m *WebSocketManager) BroadcastTaskStatusUpdated(taskID uint, oldStatus, newStatus string, updatedBy User) {
    message := Message{
        Type: "task_status_updated",
        Payload: map[string]interface{}{
            "taskId":    taskID,
            "oldStatus": oldStatus,
            "newStatus": newStatus,
            "updatedBy": updatedBy,
            "timestamp": time.Now().UTC().Format(time.RFC3339),
        },
    }
    m.broadcastMessage(message)
}

func (m *WebSocketManager) BroadcastTaskAssigned(taskID uint, assignedTo, assignedBy User) {
    message := Message{
        Type: "task_assigned",
        Payload: map[string]interface{}{
            "taskId":     taskID,
            "assignedTo": assignedTo,
            "assignedBy": assignedBy,
            "timestamp":  time.Now().UTC().Format(time.RFC3339),
        },
    }
    m.broadcastMessage(message)
}

func (m *WebSocketManager) BroadcastToProject(projectID uint, message Message) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    for client := range m.clients {
        if client.projects[projectID] {
            select {
            case client.send <- messageToBytes(message):
            default:
                close(client.send)
                delete(m.clients, client)
            }
        }
    }
}
```

### Integration with Existing Handlers
```go
// handlers/task_handler.go (updated)
func (h *TaskHandler) UpdateTaskStatus(c echo.Context) error {
    // ... existing logic ...
    
    // After successful update, broadcast WebSocket event
    if err := h.taskService.UpdateTaskStatus(taskID, newStatus); err == nil {
        task, _ := h.taskService.GetTask(taskID)
        updatedBy := getUserFromContext(c)
        
        h.websocketManager.BroadcastTaskStatusUpdated(
            taskID, 
            oldStatus, 
            newStatus, 
            updatedBy,
        )
    }
    
    return c.JSON(200, updatedTask)
}
```

## 🌐 Frontend WebSocket Integration

### WebSocket Service
```javascript
// services/websocketService.js
class KanbanWebSocket {
  constructor(token) {
    this.token = token
    this.ws = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectInterval = 5000
    this.listeners = new Map()
    this.projectSubscriptions = new Set()
  }

  connect() {
    try {
      this.ws = new WebSocket(`ws://localhost:8080/ws?token=${this.token}`)
      
      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0
        this.emit('connection_established', { message: 'Connected' })
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          this.handleMessage(data)
        } catch (error) {
          console.error('Error parsing WebSocket message:', error)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.handleReconnect()
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

    } catch (error) {
      console.error('Failed to connect WebSocket:', error)
      this.handleReconnect()
    }
  }

  handleMessage(data) {
    const { type, payload } = data
    
    // Emit to specific listeners
    if (this.listeners.has(type)) {
      this.listeners.get(type).forEach(callback => {
        try {
          callback(payload)
        } catch (error) {
          console.error(`Error in ${type} listener:`, error)
        }
      })
    }

    // Emit to general message listeners
    this.emit('message', data)
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, [])
    }
    this.listeners.get(event).push(callback)
  }

  off(event, callback) {
    if (this.listeners.has(event)) {
      const callbacks = this.listeners.get(event)
      const index = callbacks.indexOf(callback)
      if (index > -1) {
        callbacks.splice(index, 1)
      }
    }
  }

  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => callback(data))
    }
  }

  subscribeToProject(projectId) {
    this.projectSubscriptions.add(projectId)
    // Send subscription message to server if needed
    this.send({
      type: 'subscribe_project',
      payload: { projectId }
    })
  }

  unsubscribeFromProject(projectId) {
    this.projectSubscriptions.delete(projectId)
    this.send({
      type: 'unsubscribe_project',
      payload: { projectId }
    })
  }

  send(data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }

  handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)
      
      setTimeout(() => {
        this.connect()
      }, this.reconnectInterval)
    } else {
      console.error('Max reconnection attempts reached')
      this.emit('connection_failed', { message: 'Unable to reconnect' })
    }
  }

  disconnect() {
    this.reconnectAttempts = this.maxReconnectAttempts // Prevent reconnection
    if (this.ws) {
      this.ws.close()
    }
  }
}
```

### Integration with React Context
```javascript
// contexts/KanbanContext.js (updated)
export const KanbanProvider = ({ children }) => {
  const [websocket, setWebsocket] = useState(null)
  const [connectionStatus, setConnectionStatus] = useState('disconnected')

  useEffect(() => {
    const token = localStorage.getItem('jwt_token')
    if (token) {
      const ws = new KanbanWebSocket(token)
      
      // Set up event listeners
      ws.on('connection_established', () => {
        setConnectionStatus('connected')
      })

      ws.on('connection_failed', () => {
        setConnectionStatus('failed')
      })

      ws.on('task_status_updated', (payload) => {
        // Update local state when task status changes
        setTasks(prev => prev.map(task => 
          task.ID === payload.taskId 
            ? { ...task, Status: payload.newStatus }
            : task
        ))
        
        // Show notification
        showNotification(`Task moved to ${payload.newStatus}`, 'info')
      })

      ws.on('task_assigned', (payload) => {
        setTasks(prev => prev.map(task => 
          task.ID === payload.taskId 
            ? { ...task, AssignedTo: payload.assignedTo }
            : task
        ))
        
        showNotification(`Task assigned to ${payload.assignedTo.name}`, 'info')
      })

      ws.on('task_created', (payload) => {
        setTasks(prev => [...prev, payload.task])
        showNotification('New task created', 'success')
      })

      ws.on('task_deleted', (payload) => {
        setTasks(prev => prev.filter(task => task.ID !== payload.taskId))
        showNotification('Task deleted', 'warning')
      })

      ws.connect()
      setWebsocket(ws)

      return () => {
        ws.disconnect()
      }
    }
  }, [])

  const subscribeToProject = (projectId) => {
    if (websocket) {
      websocket.subscribeToProject(projectId)
    }
  }

  return (
    <KanbanContext.Provider value={{
      // ... existing values ...
      websocket,
      connectionStatus,
      subscribeToProject
    }}>
      {children}
    </KanbanContext.Provider>
  )
}
```

## 🔧 Backend Implementation Steps

### 1. Add WebSocket Dependencies
```bash
go get github.com/gorilla/websocket
```

### 2. Create WebSocket Package Structure
```
websocket/
├── manager.go      # WebSocket connection manager
├── events.go       # Event broadcasting functions
├── client.go       # Client connection handling
└── handlers.go     # WebSocket HTTP handlers
```

### 3. Update Main Server
```go
// main.go (add WebSocket upgrade endpoint)
e.GET("/ws", func(c echo.Context) error {
    return websocketHandler.HandleConnection(c)
})
```

### 4. Integrate with Existing Services
- Update task handlers to broadcast events
- Update sprint handlers to broadcast events
- Add WebSocket manager to service constructors

## 🧪 Testing WebSocket Implementation

### Backend Tests
```go
// tests/websocket_test.go
func TestWebSocketConnection(t *testing.T) {
    // Test WebSocket connection establishment
    // Test message broadcasting
    // Test client subscription management
}
```

### Frontend Tests
```javascript
// tests/websocket.test.js
describe('KanbanWebSocket', () => {
  test('should connect successfully', () => {
    // Mock WebSocket and test connection
  })

  test('should handle task status updates', () => {
    // Test event handling
  })

  test('should reconnect on disconnection', () => {
    // Test reconnection logic
  })
})
```

## 🚨 Error Handling & Edge Cases

### Connection Issues
- Network interruptions
- Server restarts
- Invalid/expired JWT tokens
- Rate limiting

### Event Handling
- Malformed messages
- Missing event listeners
- Concurrent updates
- Permission validation

### Performance Considerations
- Connection pooling
- Message batching
- Memory management
- Scalability

## 📊 Monitoring & Analytics

### Metrics to Track
- Active WebSocket connections
- Message throughput
- Connection duration
- Error rates
- Reconnection attempts

### Logging
- Connection events
- Message broadcasts
- Error conditions
- Performance metrics

## 🔒 Security Considerations

### Authentication
- JWT token validation
- Token expiration handling
- User permission checks

### Authorization
- Project access validation
- Event filtering by user permissions
- Preventing unauthorized event broadcasting

### Rate Limiting
- Connection rate limits
- Message frequency limits
- DDoS protection