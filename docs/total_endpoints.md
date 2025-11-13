# Documentación Total de Endpoints de la API

Este documento proporciona una lista exhaustiva y detallada de todos los endpoints disponibles en la API de Gestión de Proyectos.

**Autenticación:** Todos los endpoints bajo `/api` requieren un token JWT en la cabecera `Authorization` con el formato `Bearer {token}`.

---

## 1. Autenticación (Público)

### `POST /login`
- **Propósito:** Autenticar a un usuario.
- **Cuerpo (Body):**
  ```json
  {
    "correo": "user@example.com",
    "contraseña": "password123"
  }
  ```

### `POST /create-admin`
- **Propósito:** Crear el usuario administrador inicial (si no existe).
- **Cuerpo (Body):**
  ```json
  {
    "nombre": "Admin",
    "apellidoPaterno": "User",
    "correo": "admin@example.com",
    "contraseña": "password123"
  }
  ```

---

## 2. Sesión (Autenticado)

### `POST /api/logout`
- **Propósito:** Invalida el token del usuario (si la lógica de lista negra está implementada). En este caso, es gestionado por el frontend.

---

## 3. Usuarios

### `GET /api/me`
- **Propósito:** Obtener los detalles del usuario actualmente autenticado.

### `GET /api/users/:id`
- **Propósito:** Obtener los detalles de un usuario específico.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del usuario.

---

## 4. Notificaciones

### `GET /api/notifications`
- **Propósito:** Obtener todas las notificaciones (leídas y no leídas) del usuario autenticado.

### `POST /api/notifications/read/all`
- **Propósito:** Marcar todas las notificaciones del usuario autenticado como leídas.

### `POST /api/notifications/:id/read`
- **Propósito:** Marcar una notificación específica como leída.
- **Parámetros de Ruta:**
    - `:id` (uint): ID de la notificación.

---

## 5. Proyectos

### `POST /api/projects`
- **Propósito:** Crear un nuevo proyecto.
- **Cuerpo (Body):**
  ```json
  {
    "name": "Nuevo Proyecto Alfa",
    "description": "Descripción del proyecto."
  }
  ```

### `GET /api/projects`
- **Propósito:** Obtener una lista de todos los proyectos a los que el usuario tiene acceso.

### `GET /api/projects/:id`
- **Propósito:** Obtener los detalles de un proyecto específico.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

### `PUT /api/projects/:id`
- **Propósito:** Actualizar los detalles de un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.
- **Cuerpo (Body):**
  ```json
  {
    "name": "Nombre del Proyecto Actualizado",
    "description": "Descripción actualizada."
  }
  ```

### `DELETE /api/projects/:id`
- **Propósito:** Eliminar un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

### `GET /api/projects/:id/unassigned-users`
- **Propósito:** Obtener una lista de usuarios que no están asignados al proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

### `GET /api/projects/:id/members`
- **Propósito:** Obtener la lista de miembros de un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

### `GET /api/projects/:id/active-sprint`
- **Propósito:** Obtener el sprint actualmente activo para un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

---

## 6. Reportes

### `GET /api/projects/:id/reports/velocity`
- **Propósito:** Obtener el reporte de velocidad de un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

### `GET /api/sprints/:id/reports/burndown`
- **Propósito:** Obtener los datos del gráfico Burndown para un sprint.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del sprint.

### `GET /api/sprints/:id/reports/commitment`
- **Propósito:** Obtener el reporte de compromiso vs. completado para un sprint.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del sprint.

---

## 7. Rúbricas

### `POST /api/rubrics`
- **Propósito:** Crear una nueva rúbrica.
- **Cuerpo (Body):** (Ejemplo)
  ```json
  {
    "name": "Rúbrica de Calidad de Código",
    "description": "Evalúa la calidad del código de un entregable.",
    "projectId": 1
  }
  ```

### `GET /api/rubrics`
- **Propósito:** Obtener todas las rúbricas.

### `GET /api/rubrics/:id`
- **Propósito:** Obtener una rúbrica específica.
- **Parámetros de Ruta:**
    - `:id` (uint): ID de la rúbrica.

### `PUT /api/rubrics/:id`
- **Propósito:** Actualizar una rúbrica.
- **Parámetros de Ruta:**
    - `:id` (uint): ID de la rúbrica.
- **Cuerpo (Body):** (Campos a actualizar)

### `DELETE /api/rubrics/:id`
- **Propósito:** Eliminar una rúbrica.
- **Parámetros de Ruta:**
    - `:id` (uint): ID de la rúbrica.

### `POST /api/rubrics/:id/duplicate`
- **Propósito:** Duplicar una rúbrica existente.
- **Parámetros de Ruta:**
    - `:id` (uint): ID de la rúbrica a duplicar.

---

## 8. Historias de Usuario (User Stories)

### `POST /api/projects/:id/userstories`
- **Propósito:** Crear una nueva historia de usuario en un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.
- **Cuerpo (Body):**
  ```json
  {
    "title": "Como usuario, quiero poder...",
    "description": "Detalles adicionales."
  }
  ```

### `GET /api/projects/:id/userstories`
- **Propósito:** Obtener todas las historias de usuario de un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

### `GET /api/userstories/:storyId`
- **Propósito:** Obtener una historia de usuario específica.
- **Parámetros de Ruta:**
    - `:storyId` (uint): ID de la historia de usuario.

### `PUT /api/userstories/:storyId`
- **Propósito:** Actualizar una historia de usuario.
- **Parámetros de Ruta:**
    - `:storyId` (uint): ID de la historia de usuario.
- **Cuerpo (Body):** (Campos a actualizar)

### `DELETE /api/userstories/:storyId`
- **Propósito:** Eliminar una historia de usuario.
- **Parámetros de Ruta:**
    - `:storyId` (uint): ID de la historia de usuario.

---

## 9. Tareas (Tasks)

### `POST /api/userstories/:storyId/tasks`
- **Propósito:** Crear una nueva tarea para una historia de usuario.
- **Parámetros de Ruta:**
    - `:storyId` (uint): ID de la historia de usuario.
- **Cuerpo (Body):**
  ```json
  {
    "title": "Implementar el botón de login"
  }
  ```

### `GET /api/userstories/:storyId/tasks`
- **Propósito:** Obtener todas las tareas de una historia de usuario.
- **Parámetros de Ruta:**
    - `:storyId` (uint): ID de la historia de usuario.

### `PUT /api/tasks/:taskId`
- **Propósito:** Actualizar una tarea.
- **Parámetros de Ruta:**
    - `:taskId` (uint): ID de la tarea.
- **Cuerpo (Body):** (Campos a actualizar)

### `DELETE /api/tasks/:taskId`
- **Propósito:** Eliminar una tarea.
- **Parámetros de Ruta:**
    - `:taskId` (uint): ID de la tarea.

### `PUT /api/tasks/:taskId/assign`
- **Propósito:** Asignar una tarea a un usuario.
- **Parámetros de Ruta:**
    - `:taskId` (uint): ID de la tarea.
- **Cuerpo (Body):**
  ```json
  {
    "userId": 5
  }
  ```

### `PUT /api/tasks/:taskId/status`
- **Propósito:** Actualizar el estado de una tarea.
- **Parámetros de Ruta:**
    - `:taskId` (uint): ID de la tarea.
- **Cuerpo (Body):**
  ```json
  {
    "status": "in_progress"
  }
  ```

### `POST /api/tasks/:id/comments`
- **Propósito:** Añadir un comentario a una tarea.
- **Parámetros de Ruta:**
    - `:id` (uint): ID de la tarea.
- **Cuerpo (Body):**
  ```json
  {
    "content": "Este es un comentario de prueba."
  }
  ```

### `GET /api/tasks/:id/comments`
- **Propósito:** Obtener todos los comentarios de una tarea.
- **Parámetros de Ruta:**
    - `:id` (uint): ID de la tarea.

---

## 10. Evaluaciones de Tareas

### `POST /api/tasks/:taskId/evaluations`
- **Propósito:** Crear una nueva evaluación para una tarea (solo rol 'docente').
- **Parámetros de Ruta:**
    - `:taskId` (uint): ID de la tarea.
- **Cuerpo (Body):**
  ```json
  {
    "rubricId": 1,
    "overallFeedback": "¡Buen trabajo!",
    "criterionEvaluations": [
      { "criterionId": 1, "score": 5, "feedback": "Excelente." },
      { "criterionId": 2, "score": 4, "feedback": "Podría mejorar." }
    ]
  }
  ```

### `GET /api/tasks/:taskId/evaluations`
- **Propósito:** Obtener las evaluaciones de una tarea.
- **Parámetros de Ruta:**
    - `:taskId` (uint): ID de la tarea.

---

## 11. Sprints

### `POST /api/projects/:id/sprints`
- **Propósito:** Crear un nuevo sprint en un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.
- **Cuerpo (Body):**
  ```json
  {
    "name": "Sprint 1 - Alfa",
    "startDate": "2025-01-01T00:00:00Z",
    "endDate": "2025-01-15T23:59:59Z"
  }
  ```

### `GET /api/projects/:id/sprints`
- **Propósito:** Obtener todos los sprints de un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.

### `GET /api/sprints/:sprintId`
- **Propósito:** Obtener un sprint específico.
- **Parámetros de Ruta:**
    - `:sprintId` (uint): ID del sprint.

### `PUT /api/sprints/:sprintId`
- **Propósito:** Actualizar un sprint.
- **Parámetros de Ruta:**
    - `:sprintId` (uint): ID del sprint.
- **Cuerpo (Body):** (Campos a actualizar)

### `DELETE /api/sprints/:sprintId`
- **Propósito:** Eliminar un sprint.
- **Parámetros de Ruta:**
    - `:sprintId` (uint): ID del sprint.

### `GET /api/sprints/:sprintId/tasks`
- **Propósito:** Obtener todas las tareas de todas las HU de un sprint.
- **Parámetros de Ruta:**
    - `:sprintId` (uint): ID del sprint.

### `PUT /api/sprints/:sprintId/status`
- **Propósito:** Actualizar el estado de un sprint (ej. 'active', 'completed').
- **Parámetros de Ruta:**
    - `:sprintId` (uint): ID del sprint.
- **Cuerpo (Body):**
  ```json
  {
    "status": "active"
  }
  ```

### `POST /api/sprints/:sprintId/userstories`
- **Propósito:** Asignar una historia de usuario a un sprint.
- **Parámetros de Ruta:**
    - `:sprintId` (uint): ID del sprint.
- **Cuerpo (Body):**
  ```json
  {
    "userStoryId": 10
  }
  ```

---

## 12. Calendario de Eventos

### `POST /api/projects/:id/events`
- **Propósito:** Crear un nuevo evento en un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.
- **Cuerpo (Body):**
  ```json
  {
    "title": "Reunión de Planificación",
    "startDate": "2025-02-01T10:00:00Z",
    "endDate": "2025-02-01T11:00:00Z"
  }
  ```

### `GET /api/projects/:id/events`
- **Propósito:** Obtener los eventos de un proyecto, con filtros opcionales de fecha.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.
- **Parámetros de Consulta (Query):**
    - `?start=YYYY-MM-DD` (Opcional)
    - `?end=YYYY-MM-DD` (Opcional)

### `PUT /api/events/:id`
- **Propósito:** Actualizar un evento.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del evento.
- **Cuerpo (Body):** (Campos a actualizar)

### `DELETE /api/events/:id`
- **Propósito:** Eliminar un evento.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del evento.

---

## 13. Administración (Solo rol 'Admin')

### `GET /api/admin/users`
- **Propósito:** Obtener una lista de todos los usuarios del sistema.

### `POST /api/admin/users`
- **Propósito:** Crear un nuevo usuario.
- **Cuerpo (Body):**
  ```json
  {
    "nombre": "Nuevo",
    "correo": "nuevo.usuario@example.com",
    "contraseña": "password123",
    "role": "user"
  }
  ```

### `POST /api/admin/users/admin`
- **Propósito:** Crear un nuevo usuario con rol de administrador.
- **Cuerpo (Body):** (Igual que para crear un usuario normal)

### `PUT /api/admin/users/:id`
- **Propósito:** Actualizar un usuario.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del usuario.
- **Cuerpo (Body):** (Campos a actualizar)

### `DELETE /api/admin/users/:id`
- **Propósito:** Eliminar un usuario.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del usuario.

### `POST /api/admin/projects/:id/members`
- **Propósito:** Añadir un miembro a un proyecto.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.
- **Cuerpo (Body):**
  ```json
  {
    "userId": 2,
    "role": "team_developer"
  }
  ```

---

## 14. Exportación de Datos

### `GET /api/projects/:id/export`
- **Propósito:** Exportar los datos de un proyecto (historias de usuario y tareas) a un archivo CSV.
- **Parámetros de Ruta:**
    - `:id` (uint): ID del proyecto.
- **Respuesta (200 OK):**
    - El cuerpo de la respuesta es el contenido del archivo CSV.
    - Las cabeceras `Content-Type` y `Content-Disposition` están configuradas para forzar la descarga del archivo en el navegador.
