# Registro de Desarrollo del Proyecto

Este documento narra el proceso de desarrollo y las decisiones arquitectónicas clave tomadas durante la creación del proyecto `API-wrkf`.

## 1. Configuración Inicial: De Esquemas TypeScript a Modelos Go

**Objetivo:** Crear una API en Go con una base de datos PostgreSQL, utilizando los esquemas existentes de TypeScript/Drizzle como referencia.

-   **Análisis de Esquemas:** El proceso comenzó analizando los archivos `.ts` proporcionados (`db.ts`, `schema.ts` y los archivos dentro de `schemas/`).
-   **Creación de Modelos GORM:** Cada esquema de TypeScript (por ejemplo, `users`, `projects`, `tasks`) se tradujo manualmente a una estructura Go correspondiente en el directorio `/models`. Se utilizaron etiquetas GORM (`gorm:"..."`) para definir claves primarias, restricciones (no nulo, único) y relaciones (claves foráneas).
-   **Conexión a la Base de Datos:** Se creó un archivo `storage/postgres.go` para manejar la conexión GORM a la base de datos PostgreSQL.
-   **Migración Automática:** Se utilizó la función `db.AutoMigrate()` para crear automáticamente todas las tablas de la base de datos basándose en los modelos de Go. Esto asegura que el esquema de la base de datos esté siempre sincronizado con los modelos de la aplicación.

## 2. Primera Iteración: API Básica con Echo

**Objetivo:** Crear un servidor web funcional que pueda manejar operaciones CRUD básicas.

-   **Elección del Framework:** Se eligió el framework [Echo](https://echo.labstack.com/) por su alto rendimiento y diseño minimalista.
-   **Estructura Inicial:** Se implementó una estructura simple:
    -   `main.go`: Manejaba la conexión a la base de datos, la migración, la inicialización del servidor Echo y las definiciones de rutas.
    -   `handlers/`: Contenía funciones para procesar solicitudes HTTP.
    -   `storage/`: Contenía "Repositorios" responsables de la interacción directa con la base de datos (por ejemplo, `CreateUser`, `GetUserByID`).
-   **Flujo de Datos:** El flujo era simple: `Solicitud` -> `main.go (Router)` -> `Handler` -> `Repository` -> `Base de Datos`.

## 3. Refactorización Arquitectónica: Implementación de una Arquitectura en Capas

**Objetivo:** Mejorar la estructura del proyecto para una mejor escalabilidad, mantenibilidad y capacidad de prueba.

Se identificó que la estructura inicial, aunque funcional, se volvería difícil de manejar a medida que la aplicación creciera. Se llevó a cabo una refactorización importante para establecer una arquitectura limpia y en capas.

-   **Introducción de la Capa de Servicio (`/services`):**
    -   Se creó un nuevo directorio `/services`.
    -   **Propósito:** Contener la lógica de negocio central. Por ejemplo, al crear un usuario, la capa de servicio es donde se agregaría la lógica para enviar un correo electrónico de bienvenida, no en el handler o el repositorio.
    -   **Desacoplamiento:** La capa `Handler` se modificó para llamar a la capa `Service`, que a su vez llama a la capa `Repository`. Esto desacopla la lógica HTTP de la lógica de negocio y la lógica de negocio de la lógica de acceso a datos.

-   **Centralización de Rutas (`/routes`):**
    -   Se creó un nuevo directorio `/routes`.
    -   **Propósito:** Eliminar todas las definiciones de rutas de `main.go`, haciendo que el punto de entrada sea más limpio.
    -   `routes/router.go` ahora contiene una función `SetupRoutes` responsable de definir todos los puntos finales de la API.

-   **Nuevo Flujo de Datos:** `Solicitud` -> `Router` -> `Handler` -> `Service` -> `Repository` -> `Base de Datos`.

## 4. Implementación de Seguridad: Autenticación JWT

**Objetivo:** Proteger la API asegurando que solo los usuarios autenticados puedan acceder a los recursos protegidos.

-   **Hash de Contraseñas:**
    -   Se añadió el paquete `golang.org/x/crypto/bcrypt`.
    -   El `UserService` se modificó para hashear automáticamente las contraseñas de los usuarios con `bcrypt` antes de guardarlas en la base de datos. Esta es una medida de seguridad crítica para evitar almacenar contraseñas en texto plano.

-   **Endpoint de Inicio de Sesión:**
    -   Se creó un endpoint `POST /login`.
    -   El `UserService` recibió un método `Login` que:
        1.  Encuentra un usuario por su correo electrónico.
        2.  Utiliza `bcrypt.CompareHashAndPassword` para verificar de forma segura si la contraseña proporcionada coincide con el hash almacenado.
        3.  Si las credenciales son válidas, genera un JSON Web Token (JWT).

-   **Generación de JWT:**
    -   Se añadió el paquete `github.com/golang-jwt/jwt/v5`.
    -   El JWT generado contiene "claims" como el ID del usuario (`sub`) y un tiempo de expiración (`exp`).
    -   El token se firma con una clave secreta.

-   **Middleware JWT (`/middleware`):**
    -   Se creó un nuevo directorio `/middleware`.
    -   Se construyó un middleware de autenticación para proteger las rutas.
    -   **Lógica:** Para cualquier solicitud entrante a una ruta protegida, el middleware verifica un encabezado `Authorization: Bearer <token>`, valida la firma y la expiración del token, y permite que la solicitud continúe o la bloquea con un error `401 No Autorizado`.

## 5. Pulido Final: Configuración Segura con Variables de Entorno

**Objetivo:** Eliminar secretos codificados del código base, siguiendo las mejores prácticas de la industria para seguridad y flexibilidad.

-   **El Problema:** La clave secreta de JWT estaba inicialmente codificada en los archivos `.go`. Esto es un riesgo de seguridad importante y dificulta la gestión de diferentes claves para desarrollo y producción.

-   **Generación de una Clave Segura:** Se generó una clave criptográficamente segura utilizando la biblioteca `crypto` de Node.js:
    ```sh
    node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
    # Salida: 5a8e02f0b339f9b67f85aa6a5160b7e134e50246e77fc273d78de1af49cfe365
    ```

-   **Configuración Centralizada (`/config`):**
    -   El archivo `config/config.go` se refactorizó para gestionar todas las configuraciones externas.
    -   Se programó para leer el `JWT_SECRET` de una variable de entorno. Si la variable no se encuentra, recurre a un valor predeterminado (la clave segura que generamos).

-   **Inyección de Dependencias:**
    -   El `UserService` y el `JWTAuthMiddleware` se refactorizaron para no tener su propio secreto codificado. En su lugar, el secreto se "inyecta" en ellos desde `main.go` durante la inicialización.

-   **Resultado Final:** La aplicación ahora es altamente segura y configurable. Puede ejecutarse con una configuración predeterminada para facilitar el desarrollo, o con secretos listos para producción pasados a través de variables de entorno, sin requerir cambios en el código.
    ```sh
    # Ejecutar con una clave personalizada y segura en un entorno de producción
    JWT_SECRET="a_different_secret_for_production" go run main.go
    ```
