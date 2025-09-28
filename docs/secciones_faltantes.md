# Análisis Comparativo: Funcionalidades Faltantes en la API de Go

Este documento detalla las funcionalidades que estaban presentes en el proyecto anterior (Deno/Fresh) y que aún no han sido implementadas en la nueva API de Go/Echo.

---

## 1. Módulo de Autenticación y Usuarios

-   **Registro de Usuarios (`/api/register.ts`):**
    -   **Falta:** Un endpoint público para que los nuevos usuarios puedan registrarse por sí mismos.
    -   **Situación Actual:** La API de Go solo permite la creación de usuarios por parte de un administrador.

-   **Cierre de Sesión (`/api/logout.ts`):**
    -   **Falta:** Un endpoint para invalidar el token de sesión del usuario en el servidor.
    -   **Importancia:** Añadir el token a una "lista negra" (blocklist) es una práctica de seguridad robusta para prevenir la reutilización de tokens robados.

-   **Información de Sesión (`/api/session.ts`):**
    -   **Falta:** Un endpoint para que el frontend pueda verificar si un usuario tiene una sesión activa y obtener sus datos.

-   **Eliminación de Usuarios (`/api/admin/users/delete.ts`):**
    -   **Falta:** Un endpoint para que los administradores puedan eliminar usuarios del sistema.

-   **Métricas de Usuario (`/api/users/[id]/metrics.ts`):**
    -   **Falta:** Un endpoint para obtener métricas específicas sobre la actividad o el rendimiento de un usuario.

---

## 2. Módulo de Tareas (Tasks)

-   **Historial de Tareas (`/api/tasks/[id]/history.ts`):**
    -   **Falta:** Funcionalidad para consultar el historial de cambios de una tarea (cambios de estado, asignaciones, etc.). Esto es crucial para la auditoría y el seguimiento.

-   **Seguimiento de Tiempo (`/api/tasks/[id]/time.ts`):**
    -   **Falta:** Capacidad para registrar y consultar el tiempo invertido en una tarea (Time Tracking).

---

## 3. Módulo de Reportes (Totalmente Ausente)

Este módulo completo necesita ser implementado.

-   **Generación de Reportes (`/api/reports/generate.ts`):**
    -   **Falta:** La capacidad de generar reportes personalizados sobre el estado de los proyectos.

-   **Exportación de Reportes (`/api/reports/[id]/export.ts`):**
    -   **Falta:** La funcionalidad para exportar los reportes generados (ej. a CSV o PDF).

-   **Programación de Reportes (`/api/reports/schedule.ts`):**
    -   **Falta:** La capacidad de programar la generación y envío automático de reportes.

---

## 4. Módulo de Conversaciones / Mensajería (Totalmente Ausente)

Este módulo completo necesita ser implementado.

-   **Gestión de Conversaciones (`/api/conversations/index.ts`):**
-   **Envío/Recepción de Mensajes (`/api/conversations/[id]/messages.ts`):**
    -   **Falta:** Todo el sistema de chat o mensajería en tiempo real.

---

## 5. Métricas y Gestión Avanzada

-   **Salud del Proyecto (`/api/projects/[id]/health.ts`):**
    -   **Falta:** Un endpoint que calcule una puntuación o estado general de "salud" del proyecto.

-   **Métricas de Proyecto (`/api/projects/[id]/metrics.ts`):**
    -   **Falta:** Un endpoint para ofrecer métricas más detalladas a nivel de proyecto.

-   **Métricas de Sprint (`/api/sprints/[id]/metrics.ts`):**
    -   **Falta:** Un endpoint para métricas de sprint más completas.

-   **Gestión del Burndown Chart:**
    -   **Falta:** Endpoints para depurar (`burndown-debug.ts`) y recalcular (`recalculate-burndown.ts`) el gráfico de Burndown.
