# Log de Desarrollo: Implementación de RBAC y Endpoint de Miembros

Este documento detalla el proceso de desarrollo para implementar un sistema de Control de Acceso Basado en Roles (RBAC) a nivel de proyecto y la creación de un endpoint para visualizar los miembros de un proyecto.

## 1. Objetivo

El objetivo principal era mejorar la seguridad y funcionalidad de la API mediante la implementación de un sistema de permisos granular que restrinja las acciones basadas en el rol de un usuario (`scrum_master`, `product_owner`, `team_developer`) dentro de un proyecto específico.

## 2. Fase de Análisis

Antes de escribir código, se realizó un análisis exhaustivo del sistema existente:

1.  **Modelos de Roles (`models/roles.go`):** Se confirmó la existencia de dos tipos de roles: `PlatformRole` (global) y `ProjectRole` (específico del proyecto). Esto validó que la estructura de datos fundamental ya estaba presente.

2.  **Autenticación (`middleware/auth_middleware.go`):** Se verificó que el middleware JWT principal ya extraía el rol de la plataforma (`admin`/`user`) del token y lo almacenaba en el contexto de la solicitud (`c.Set("userRole", ...)`).

3.  **Autorización de Administrador (`middleware/admin_auth_middleware.go`):** Se observó cómo este middleware utilizaba el rol de la plataforma para proteger las rutas de administrador, proporcionando un patrón a seguir.

4.  **Capa de Servicio y Repositorio (`services/project_service.go`, `storage/project_repository.go`):** Se descubrió que las funciones necesarias para obtener el rol de un usuario en un proyecto (`GetUserRoleInProject`) ya estaban implementadas, lo que aceleró significativamente el desarrollo.

## 3. Fase de Implementación

### 3.1. Middleware de Autorización por Rol de Proyecto

El núcleo de la implementación fue la creación de un nuevo middleware.

-   **Archivo Creado:** `middleware/project_auth_middleware.go`
-   **Funcionalidad:** Se creó una función `ProjectRoleAuth` que actúa como una fábrica de middleware. Esta función:
    1.  Acepta los roles de proyecto requeridos como parámetros (ej. `models.RoleScrumMaster`).
    2.  Extrae el `userID` del contexto del token JWT.
    3.  Extrae el `projectID` de los parámetros de la URL.
    4.  Llama a `projectService.GetUserRoleInProject` para obtener el rol real del usuario en la base de datos.
    5.  Compara el rol real con los roles requeridos y permite o deniega el acceso (con un error `403 Forbidden`).

### 3.2. Integración y Primer Caso de Uso

Una vez creado el middleware, se integró en la aplicación:

1.  **Modificación del Router (`routes/router.go`):**
    -   Se actualizó la firma de `SetupRoutes` para inyectar el `ProjectService`, haciéndolo disponible para el nuevo middleware.
    -   Se aplicó el middleware a la ruta `POST /api/projects/:id/sprints`, restringiendo la creación de sprints a los roles `RoleScrumMaster` y `RoleProductOwner`.

2.  **Corrección en `main.go`:**
    -   Se actualizó la llamada a `routes.SetupRoutes` en `main.go` para pasar la instancia de `projectService`, solucionando el error de compilación introducido por el cambio en la firma.

### 3.3. Endpoint para Ver Miembros del Proyecto

Para completar la funcionalidad y atender a una solicitud explícita, se creó un endpoint para ver los miembros del equipo.

1.  **Nuevo Manejador (`handlers/project_handler.go`):**
    -   Se añadió la función `GetProjectMembers`.
    -   Esta función utiliza el `projectService` para obtener los detalles del proyecto, que ya incluyen una lista precargada de miembros y sus datos de usuario.
    -   Se definió un modelo de respuesta JSON limpio (`MemberResponse`) para evitar exponer los modelos internos de la base de datos.

2.  **Nueva Ruta (`routes/router.go`):**
    -   Se añadió la ruta `GET /api/projects/:id/members`.
    -   Se protegió esta ruta con el middleware `ProjectRoleAuth`, configurado para permitir el acceso a cualquier usuario que tenga un rol asignado en ese proyecto (`RoleScrumMaster`, `RoleProductOwner`, `RoleTeamDeveloper`).

## 4. Verificación

Todos los cambios se verificaron de forma incremental utilizando `go build ./...` para asegurar que el proyecto compilara correctamente en cada paso, garantizando la estabilidad del código base.
