# Sistema de Notificaciones

Este documento detalla la arquitectura y el funcionamiento del sistema de notificaciones dentro de la API.

## 1. Propósito

El sistema de notificaciones tiene como objetivo mantener a los usuarios informados sobre eventos importantes que ocurren en la plataforma y que son relevantes para ellos, mejorando la colaboración y el seguimiento del trabajo.

## 2. Arquitectura

La funcionalidad se ha implementado siguiendo la arquitectura en capas del proyecto para mantener la separación de responsabilidades y la modularidad.

-   **Modelo (`models/notification.go`):** Define la estructura de datos de una notificación y cómo se mapea a la tabla `notifications` en la base de datos.
-   **Repositorio (`storage/notification_repository.go`):** Contiene la lógica de acceso directo a la base de datos para realizar operaciones CRUD (Crear, Leer, Actualizar) sobre las notificaciones.
-   **Servicio (`services/notification_service.go`):** Orquesta la lógica de negocio. Por ejemplo, encapsula la creación de una notificación y es utilizado por otros servicios para generar notificaciones.
-   **Handler (`handlers/notification_handler.go`):** Expone la funcionalidad a través de la API REST, manejando las solicitudes HTTP, la validación de entradas y las respuestas.
-   **Utilidad JWT (`utils/jwt_utils.go`):** Se creó una función auxiliar para extraer de forma segura el ID del usuario desde el token JWT en el contexto de la solicitud.

## 3. Modelo de Datos (Schema)

El modelo `Notification` tiene la siguiente estructura:

-   `ID`: Identificador único de la notificación.
-   `UserID`: ID del usuario que recibe la notificación. La tabla tiene un índice en esta columna para búsquedas eficientes.
-   `Message`: El contenido textual de la notificación (ej. "Se te ha asignado una nueva tarea").
-   `IsRead`: Un booleano que indica si el usuario ya ha visto la notificación. Por defecto es `false`.
-   `Link`: Una URL relativa que permite al frontend redirigir al usuario al recurso relevante (ej. `/tasks/123`).
-   `CreatedAt`, `UpdatedAt`, `DeletedAt`: Campos estándar de GORM para el seguimiento del ciclo de vida del registro.

## 4. Endpoints de la API

Se han añadido los siguientes endpoints, todos protegidos y que requieren un token JWT válido:

-   **`GET /api/notifications`**
    -   **Descripción:** Obtiene todas las notificaciones para el usuario autenticado, ordenadas por fecha de creación descendente.
    -   **Controlador:** `notificationHandler.GetUserNotifications`

-   **`POST /api/notifications/:id/read`**
    -   **Descripción:** Marca una notificación específica como leída. El usuario solo puede marcar sus propias notificaciones.
    -   **Controlador:** `notificationHandler.MarkAsRead`

-   **`POST /api/notifications/read/all`**
    -   **Descripción:** Marca todas las notificaciones no leídas del usuario como leídas.
    -   **Controlador:** `notificationHandler.MarkAllAsRead`

-   **`POST /api/tasks/:id/comments`**
    -   **Descripción:** Añade un nuevo comentario a una tarea.
    -   **Controlador:** `taskHandler.AddComment`

## 5. Triggers Automáticos de Notificaciones

Actualmente, el sistema genera notificaciones de forma automática en los siguientes casos:

1.  **Asignación de Proyecto:**
    -   **Disparador:** Un administrador añade un usuario a un proyecto (`POST /api/admin/projects/:id/members`).
    -   **Lógica:** La función `AddMemberToProject` en `services/project_service.go` llama al `NotificationService` para crear una notificación para el usuario que ha sido añadido.
    -   **Mensaje:** "Has sido añadido al proyecto '[Nombre del Proyecto]' con el rol de '[Rol]'."

2.  **Asignación de Tarea:**
    -   **Disparador:** Se asigna una tarea a un usuario (`PUT /api/tasks/:taskId/assign`).
    -   **Lógica:** La función `AssignTask` en `services/task_service.go` llama al `NotificationService` para notificar al usuario asignado.
    -   **Mensaje:** "Se te ha asignado la tarea '[Título de la Tarea]'."

3.  **Nuevo Comentario en Tarea:**
    -   **Disparador:** Un usuario publica un comentario en una tarea (`POST /api/tasks/:id/comments`).
    -   **Lógica:** La función `AddCommentToTask` en `services/task_service.go` crea una notificación para el usuario al que está asignada la tarea, siempre y cuando no sea la misma persona que ha escrito el comentario.
    -   **Mensaje:** "Nuevo comentario en la tarea '[Título de la Tarea]'."

## 6. Futuras Mejoras

El sistema de notificaciones puede expandirse para incluir otros eventos, como:
-   Cambios de estado en tareas importantes.
-   Menciones a usuarios en comentarios (ej. `@username`).
-   Recordatorios de fechas de entrega próximas.
