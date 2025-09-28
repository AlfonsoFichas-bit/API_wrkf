# Análisis de Implementaciones Faltantes: Deno/Fresh vs. Go/Echo

Este documento detalla las características y funcionalidades presentes en el proyecto anterior (Deno/Fresh) que aún no se han implementado en la nueva API de Go/Echo.

## 1. Funcionalidades Principales de Scrum

### Gestión del Backlog (Implementado)
- **Descripción:** El proyecto anterior tenía una interfaz de usuario dedicada para el backlog del producto, que incluía filtros, métricas y funcionalidad de arrastrar y soltar para actualizar el estado de las historias de usuario.
- **Componentes Faltantes en la API de Go:**
    - Endpoints de API para ver el backlog.
    - Endpoints para añadir/eliminar historias de usuario del backlog.
    - Endpoints para actualizar el estado de las historias de usuario.

### Planificación de Sprints (Implementado)
- **Descripción:** El proyecto anterior permitía la planificación de sprints, incluida la adición de historias de usuario a un sprint específico.
- **Componentes Faltantes en la API de Go:**
    - Endpoints de API y lógica de servicio para la planificación de sprints, como la asignación de una historia de usuario a un sprint.

### Tablero de Tareas/Kanban (Implementado)
- **Descripción:** El proyecto anterior incluía un tablero de tareas visual, probablemente con funcionalidad de arrastrar y soltar para cambiar el estado de las tareas.
- **Componentes Faltantes en la API de Go:**
    - Endpoints de API para dar soporte a un tablero de estilo Kanban, como la actualización del estado de una tarea.

## 2. Evaluaciones y Métricas

### Evaluaciones (Implementado)
- **Descripción:** El proyecto anterior contaba con un sistema de evaluación completo que permitía crear evaluaciones, ver el historial de evaluaciones, ver estadísticas y gestionar los entregables.
- **Componentes Faltantes en la API de Go:**
    - Handlers y servicios para gestionar todo el ciclo de vida de las evaluaciones, incluyendo su creación, envío y visualización de resultados.

### Métricas e Informes (Implementado)
- **Descripción:** El proyecto anterior tenía un sistema de métricas e informes muy completo. Podía generar gráficos de burndown, gráficos de velocidad del equipo, gráficos de distribución del trabajo y medidores de salud del proyecto. También contaba con un generador de informes.
- **Componentes Faltantes en la API de Go:**
    - Endpoints de API y servicios para calcular y exponer estas métricas.
    - Lógica para generar gráficos de burndown, calcular la velocidad del equipo y crear informes.

## 3. Gestión de Usuarios y Administradores

### Panel de Administración (Implementado)
- **Descripción:** El proyecto anterior tenía un panel de administración para la gestión de usuarios.
- **Componentes Faltantes en la API de Go:**
    - Endpoints de API para funcionalidades específicas de administración, como listar, crear y eliminar usuarios.

### Control de Acceso Basado en Roles (RBAC)
- **Descripción:** El proyecto anterior tenía un sistema de permisos que adaptaba la experiencia del usuario en función de su rol.
- **Componentes Faltantes en la API de Go:**
    - Middleware y lógica de servicio para hacer cumplir el control de acceso basado en roles en los endpoints de la API.

## 4. Otras Funcionalidades

### Chat/Mensajería
- **Descripción:** El proyecto anterior incluía una función de chat en tiempo real.
- **Componentes Faltantes en la API de Go:**
    - Handlers y lógica de WebSocket para implementar el chat en tiempo real.

### Entregables
- **Descripción:** El proyecto anterior contaba con un sistema para la gestión de los entregables del proyecto.
- **Componentes Faltantes en la API de Go:**
    - Modelo, handler y servicio para la gestión de entregables.
