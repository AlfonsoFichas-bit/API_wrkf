# Nuevos Endpoints para Kanban de Sprint

## Endpoints Implementados

### 1. GET /api/sprints/:sprintId/tasks
**Descripción**: Obtiene todas las tareas de un sprint específico con sus relaciones completas.

**Parámetros**:
- `sprintId` (path): ID del sprint

**Respuesta**:
```json
[
  {
    "id": 1,
    "title": "Task 1",
    "description": "Description of task 1",
    "status": "todo",
    "userStoryId": 5,
    "assignedToId": 3,
    "assignedTo": {
      "id": 3,
      "nombre": "Juan",
      "apellidoPaterno": "Pérez"
    },
    "createdBy": {
      "id": 2,
      "nombre": "María",
      "apellidoPaterno": "García"
    },
    "userStory": {
      "id": 5,
      "title": "User Story 1",
      "project": {
        "id": 1,
        "name": "Project Alpha"
      }
    }
  }
]
```

**Uso para Kanban**: Este endpoint proporciona todas las tareas necesarias para el tablero Kanban de un sprint.

---

### 2. GET /api/projects/:id/active-sprint
**Descripción**: Obtiene el sprint actualmente activo de un proyecto.

**Parámetros**:
- `id` (path): ID del proyecto

**Respuesta**:
```json
{
  "id": 2,
  "name": "Sprint 2",
  "goal": "Complete user authentication",
  "status": "active",
  "startDate": "2025-01-15T00:00:00Z",
  "endDate": "2025-01-29T00:00:00Z",
  "projectId": 1,
  "project": {
    "id": 1,
    "name": "Project Alpha"
  },
  "createdBy": {
    "id": 2,
    "nombre": "María",
    "apellidoPaterno": "García"
  }
}
```

**Uso para Kanban**: Permite identificar automáticamente qué sprint mostrar en el tablero Kanban.

---

### 3. PUT /api/sprints/:sprintId/status
**Descripción**: Actualiza el estado de un sprint.

**Parámetros**:
- `sprintId` (path): ID del sprint

**Body**:
```json
{
  "status": "active"
}
```

**Estados Válidos**:
- `planned`: Sprint planificado
- `active`: Sprint activo actualmente
- `completed`: Sprint finalizado
- `cancelled`: Sprint cancelado

**Respuesta**:
```json
{
  "id": 2,
  "name": "Sprint 2",
  "status": "active",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

**Uso para Kanban**: Permite cambiar el estado del sprint (ej: de planned a active).

---

## Flujo de Uso para Kanban

### 1. Obtener Sprint Activo
```bash
GET /api/projects/1/active-sprint
```

### 2. Cargar Tareas del Sprint
```bash
GET /api/sprints/2/tasks
```

### 3. Actualizar Estados de Tareas (usando endpoints existentes)
```bash
PUT /api/tasks/123/status
{
  "status": "in_progress"
}
```

### 4. Actualizar Estado del Sprint (opcional)
```bash
PUT /api/sprints/2/status
{
  "status": "completed"
}
```

---

## Optimizaciones Implementadas

### Query Optimizada para Sprint Tasks
- **JOIN directo** entre tasks y user_stories
- **Preloading** de relaciones necesarias (AssignedTo, CreatedBy, UserStory, Project)
- **Ordenamiento** por fecha de creación

### Estado del Sprint
- **Búsqueda eficiente** por project_id y status
- **Preloading** de relaciones (CreatedBy, Project)

---

## Pruebas Implementadas

### Tests Unitarios
- `TestSprintService_GetSprintTasks_Structure`: Verifica estructura del método
- `TestSprintService_UpdateSprintStatus_Structure`: Verifica firma del método
- `TestUpdateSprintStatusRequest_Validation`: Validación de estados

### Calidad de Código
- **golangci-lint**: 0 issues
- **Compilación**: Exitosa
- **Tests**: Pasando correctamente

---

## Próximos Pasos para Kanban

Con estos endpoints implementados, la base para el Kanban está lista:

1. ✅ **Obtener sprint activo**
2. ✅ **Cargar tareas del sprint**
3. ✅ **Actualizar estados de tareas**
4. ✅ **Actualizar estado del sprint**

### Actualizaciones en Tiempo Real (WebSockets)

- ✅ **WebSocket integration para actualizaciones en tiempo real**: ¡Implementado! Las acciones como crear, actualizar estado, asignar o eliminar una tarea ahora emiten eventos a través de WebSockets para actualizar el tablero en todos los clientes conectados. Para más detalles sobre los eventos, ver `api_documentation.md`.

**Lo que falta para Kanban completo**:
- Frontend del tablero Kanban
- Drag & drop functionality
- Métricas del sprint

---

## Ejemplo de Flujo Completo

```javascript
// 1. Obtener sprint activo
const activeSprint = await fetch('/api/projects/1/active-sprint');

// 2. Cargar tareas del sprint
const tasks = await fetch(`/api/sprints/${activeSprint.id}/tasks`);

// 3. Organizar tareas por estado para Kanban
const kanbanColumns = {
  todo: tasks.filter(task => task.status === 'todo'),
  in_progress: tasks.filter(task => task.status === 'in_progress'),
  in_review: tasks.filter(task => task.status === 'in_review'),
  done: tasks.filter(task => task.status === 'done')
};

// 4. Actualizar estado cuando se mueve una tarea
await fetch(`/api/tasks/${taskId}/status`, {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ status: 'in_progress' })
});
```