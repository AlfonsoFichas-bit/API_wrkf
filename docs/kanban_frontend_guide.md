# Frontend Kanban Implementation Guide

## Overview
This guide provides complete instructions for implementing a Kanban board frontend that integrates with the existing backend API.

## Prerequisites
- Backend API running on `http://localhost:8080`
- Valid JWT token for authentication
- Modern frontend framework (React, Vue, Angular recommended)

## 🏗️ Architecture

### 1. Core Components Needed

#### KanbanBoard Component
```javascript
// Main container for the Kanban board
const KanbanBoard = ({ projectId }) => {
  const [sprint, setSprint] = useState(null)
  const [tasks, setTasks] = useState([])
  const [columns, setColumns] = useState([
    { id: 'todo', title: 'To Do', status: 'todo' },
    { id: 'in_progress', title: 'In Progress', status: 'in_progress' },
    { id: 'done', title: 'Done', status: 'done' }
  ])
  
  // Load active sprint and tasks
  useEffect(() => {
    loadActiveSprint()
    loadSprintTasks()
  }, [projectId])
}
```

#### KanbanColumn Component
```javascript
// Individual column for each status
const KanbanColumn = ({ column, tasks, onTaskMove }) => {
  const columnTasks = tasks.filter(task => task.Status === column.status)
  
  return (
    <div className="kanban-column">
      <h3>{column.title}</h3>
      {columnTasks.map(task => (
        <KanbanCard key={task.ID} task={task} onMove={onTaskMove} />
      ))}
    </div>
  )
}
```

#### KanbanCard Component
```javascript
// Individual task card
const KanbanCard = ({ task, onMove }) => {
  return (
    <div 
      className="kanban-card"
      draggable
      onDragStart={(e) => handleDragStart(e, task)}
      onDragOver={handleDragOver}
      onDrop={(e) => handleDrop(e, task)}
    >
      <h4>{task.Title}</h4>
      <p>{task.Description}</p>
      <div className="task-meta">
        <span className="priority">{task.UserStory?.Priority}</span>
        <span className="points">{task.UserStory?.Points} pts</span>
        {task.AssignedTo && (
          <span className="assignee">{task.AssignedTo.Nombre}</span>
        )}
      </div>
    </div>
  )
}
```

### 2. API Integration

#### API Service Layer
```javascript
// api/kanbanService.js
const API_BASE = 'http://localhost:8080'

class KanbanService {
  constructor(token) {
    this.token = token
  }

  async getActiveSprint(projectId) {
    const response = await fetch(`${API_BASE}/api/projects/${projectId}/active-sprint`, {
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!response.ok) {
      if (response.status === 404) return null
      throw new Error('Failed to fetch active sprint')
    }
    
    return response.json()
  }

  async getSprintTasks(sprintId) {
    const response = await fetch(`${API_BASE}/api/sprints/${sprintId}/tasks`, {
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch sprint tasks')
    }
    
    return response.json()
  }

  async updateTaskStatus(taskId, newStatus) {
    const response = await fetch(`${API_BASE}/api/tasks/${taskId}/status`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ status: newStatus })
    })
    
    if (!response.ok) {
      throw new Error('Failed to update task status')
    }
    
    return response.json()
  }

  async updateSprintStatus(sprintId, newStatus) {
    const response = await fetch(`${API_BASE}/api/sprints/${sprintId}/status`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ status: newStatus })
    })
    
    if (!response.ok) {
      throw new Error('Failed to update sprint status')
    }
    
    return response.json()
  }
}
```

### 3. State Management

#### React Context Example
```javascript
// contexts/KanbanContext.js
const KanbanContext = createContext()

export const KanbanProvider = ({ children }) => {
  const [sprint, setSprint] = useState(null)
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const kanbanService = new KanbanService(localStorage.getItem('jwt_token'))

  const loadActiveSprint = async (projectId) => {
    setLoading(true)
    try {
      const activeSprint = await kanbanService.getActiveSprint(projectId)
      setSprint(activeSprint)
      if (activeSprint) {
        const sprintTasks = await kanbanService.getSprintTasks(activeSprint.ID)
        setTasks(sprintTasks)
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const moveTask = async (taskId, newStatus) => {
    try {
      const updatedTask = await kanbanService.updateTaskStatus(taskId, newStatus)
      setTasks(prev => prev.map(task => 
        task.ID === taskId ? updatedTask : task
      ))
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <KanbanContext.Provider value={{
      sprint,
      tasks,
      loading,
      error,
      loadActiveSprint,
      moveTask
    }}>
      {children}
    </KanbanContext.Provider>
  )
}
```

### 4. Drag & Drop Implementation

#### HTML5 Drag & Drop
```javascript
// hooks/useDragAndDrop.js
export const useDragAndDrop = (onTaskMove) => {
  const [draggedTask, setDraggedTask] = useState(null)

  const handleDragStart = (e, task) => {
    setDraggedTask(task)
    e.dataTransfer.effectAllowed = 'move'
  }

  const handleDragOver = (e) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
  }

  const handleDrop = (e, targetStatus) => {
    e.preventDefault()
    if (draggedTask && draggedTask.Status !== targetStatus) {
      onTaskMove(draggedTask.ID, targetStatus)
    }
    setDraggedTask(null)
  }

  return {
    handleDragStart,
    handleDragOver,
    handleDrop,
    draggedTask
  }
}
```

### 5. Styling

#### CSS for Kanban Board
```css
/* KanbanBoard.css */
.kanban-board {
  display: flex;
  gap: 20px;
  padding: 20px;
  min-height: 600px;
  overflow-x: auto;
}

.kanban-column {
  flex: 1;
  min-width: 300px;
  background: #f5f5f5;
  border-radius: 8px;
  padding: 16px;
}

.kanban-column h3 {
  margin: 0 0 16px 0;
  color: #333;
  font-weight: 600;
}

.kanban-card {
  background: white;
  border: 1px solid #ddd;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 8px;
  cursor: move;
  transition: all 0.2s ease;
}

.kanban-card:hover {
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.kanban-card.dragging {
  opacity: 0.5;
  transform: rotate(5deg);
}

.task-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  font-size: 12px;
}

.priority {
  padding: 2px 6px;
  border-radius: 12px;
  font-weight: 600;
}

.priority.high { background: #ffebee; color: #c62828; }
.priority.medium { background: #fff3e0; color: #ef6c00; }
.priority.low { background: #e8f5e8; color: #2e7d32; }

.points {
  background: #e3f2fd;
  color: #1565c0;
  padding: 2px 6px;
  border-radius: 12px;
  font-weight: 600;
}

.assignee {
  background: #f3e5f5;
  color: #7b1fa2;
  padding: 2px 6px;
  border-radius: 12px;
  font-size: 11px;
}
```

## 🚀 Implementation Steps

### Step 1: Setup Project Structure
```
frontend/
├── src/
│   ├── components/
│   │   ├── KanbanBoard.js
│   │   ├── KanbanColumn.js
│   │   └── KanbanCard.js
│   ├── services/
│   │   └── kanbanService.js
│   ├── contexts/
│   │   └── KanbanContext.js
│   ├── hooks/
│   │   └── useDragAndDrop.js
│   └── styles/
│       └── KanbanBoard.css
```

### Step 2: Install Dependencies
```bash
# For React
npm install axios

# For Vue
npm install axios

# For Angular
npm install axios
```

### Step 3: Implement Core Components
1. Create `KanbanService` for API calls
2. Build `KanbanCard` component with drag & drop
3. Create `KanbanColumn` component
4. Implement main `KanbanBoard` component

### Step 4: Add State Management
1. Create context/store for Kanban state
2. Implement loading and error states
3. Add optimistic updates for better UX

### Step 5: Style the Components
1. Add responsive design
2. Implement hover states and transitions
3. Add loading indicators
4. Style error states

### Step 6: Add Advanced Features
1. Task filtering and search
2. Sprint metrics and burndown chart
3. Real-time updates (WebSocket)
4. Task detail modal/panel

## 🔄 WebSocket Integration (Optional but Recommended)

### WebSocket Service
```javascript
// services/websocketService.js
class WebSocketService {
  constructor(token) {
    this.token = token
    this.ws = null
    this.listeners = {}
  }

  connect() {
    this.ws = new WebSocket(`ws://localhost:8080/ws?token=${this.token}`)
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.handleMessage(data)
    }
  }

  on(event, callback) {
    if (!this.listeners[event]) {
      this.listeners[event] = []
    }
    this.listeners[event].push(callback)
  }

  handleMessage(data) {
    const { type, payload } = data
    if (this.listeners[type]) {
      this.listeners[type].forEach(callback => callback(payload))
    }
  }
}
```

## 📱 Mobile Considerations

### Responsive Design
- Use CSS Grid/Flexbox for responsive layout
- Touch-friendly drag & drop for mobile
- Horizontal scroll for columns on small screens
- Collapsible columns on mobile

### Touch Events
```javascript
// Add touch support for mobile
const handleTouchStart = (e, task) => {
  // Implement touch-based drag
}

const handleTouchMove = (e) => {
  // Handle touch movement
}

const handleTouchEnd = (e) => {
  // Handle touch end
}
```

## 🧪 Testing

### Unit Tests
```javascript
// Example test for KanbanService
describe('KanbanService', () => {
  test('should fetch active sprint', async () => {
    const service = new KanbanService('fake-token')
    const sprint = await service.getActiveSprint(1)
    expect(sprint).toBeDefined()
  })
})
```

### Integration Tests
- Test drag & drop functionality
- Test API integration
- Test error handling
- Test loading states

## 🎯 Best Practices

1. **Error Handling**: Always handle API errors gracefully
2. **Loading States**: Show loading indicators during API calls
3. **Optimistic Updates**: Update UI immediately, rollback on error
4. **Accessibility**: Add ARIA labels and keyboard navigation
5. **Performance**: Implement virtual scrolling for large task lists
6. **Caching**: Cache API responses to reduce network calls

## 🔧 Troubleshooting

### Common Issues
1. **CORS Errors**: Ensure backend allows frontend origin
2. **Authentication**: Check JWT token is properly stored and sent
3. **Drag & Drop**: Test on different browsers and devices
4. **Performance**: Monitor bundle size and implement code splitting

### Debug Tips
- Use browser dev tools to inspect API calls
- Add console logging for drag & drop events
- Test with different screen sizes
- Monitor network requests and responses

## 📚 Additional Resources

- [React DnD](https://react-dnd.github.io/react-dnd/) - Advanced drag & drop for React
- [Vue.Draggable](https://github.com/SortableJS/Vue.Draggable) - Drag & drop for Vue
- [ng2-dragula](https://github.com/valor-software/ng2-dragula) - Drag & drop for Angular
- [HTML5 Drag & Drop API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API) - Native browser API