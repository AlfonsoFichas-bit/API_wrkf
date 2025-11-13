# Análisis de Endpoints Faltantes para la Vista de Detalles de Tarea

Basado en el análisis de la maqueta de la interfaz de usuario `task_details_view_2.png`, se han identificado los siguientes endpoints de API que son necesarios para implementar completamente la funcionalidad mostrada, pero que actualmente no existen en la API.

## 1. Obtener Detalles de una Tarea Específica

*   **Método HTTP:** `GET`
*   **Ruta:** `/api/tasks/:id`
*   **Descripción:** Este endpoint es crucial para la vista de detalles de la tarea. Debe devolver toda la información de una tarea específica, incluyendo su título, descripción, estado, prioridad, fecha de vencimiento y el usuario asignado. Actualmente, no hay una forma directa de obtener los datos de una sola tarea por su ID.

## 2. Obtener el Historial de Actividad de una Tarea

*   **Método HTTP:** `GET`
*   **Ruta:** `/api/tasks/:id/history`
*   **Descripción:** La pestaña "Activity Log" en la interfaz de usuario implica que se debe poder consultar un historial de cambios para una tarea. Este endpoint debería devolver una lista de eventos o cambios que ha sufrido la tarea a lo largo del tiempo, como cambios de estado, actualizaciones en la descripción o cambios de asignación.
