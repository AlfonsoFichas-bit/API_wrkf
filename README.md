# API-wrkf (API Go con Echo y GORM)

Este proyecto es una API REST robusta y escalable construida con Go, utilizando el framework Echo para manejar las solicitudes HTTP y GORM como ORM para interactuar con una base de datos PostgreSQL. Presenta una arquitectura limpia y en capas, y una autenticación segura mediante JSON Web Tokens (JWT).

## Características

-   **Arquitectura en Capas**: Clara separación de responsabilidades (Handlers, Servicios, Repositorios).
-   **Base de Datos PostgreSQL**: Utiliza GORM para el mapeo objeto-relacional.
-   **Autenticación JWT**: Protege los endpoints utilizando middleware JWT.
-   **Actualizaciones en Tiempo Real**: Utiliza WebSockets para notificar a los clientes de cambios en el tablero Kanban instantáneamente.
-   **Gestión de Configuración**: Configuración centralizada cargada desde variables de entorno.
-   **Hash de Contraseñas**: Almacena de forma segura las contraseñas de los usuarios utilizando `bcrypt`.

## Estructura del Proyecto

```
/
├── config/         # Configuración de la base de datos y la aplicación
├── docs/           # Documentación del proyecto
├── handlers/       # Manejadores de solicitudes HTTP
├── middleware/     # Middleware de autenticación JWT
├── models/         # Modelos de datos GORM
├── routes/         # Definiciones de rutas de la API
├── services/       # Capa de lógica de negocio
├── storage/        # Capa de acceso a la base de datos (repositorios)
├── go.mod
├── go.sum
└── main.go         # Punto de entrada de la aplicación
```

## Primeros Pasos

### Prerrequisitos

-   Go (versión 1.18 o superior)
-   PostgreSQL
-   Node.js (opcional, para generar una clave secreta)

### Instalación y Configuración

1.  **Clonar el repositorio:**
    ```sh
    git clone https://github.com/AlfonsoFichas-bit/API_wrkf.git
    cd API_wrkf
    ```

2.  **Instalar dependencias de Go:**
    ```sh
    go mod tidy
    ```

3.  **Crear la base de datos PostgreSQL:**
    Conéctate a PostgreSQL y ejecuta:
    ```sql
    CREATE DATABASE api_db;
    ```

4.  **Configurar variables de entorno (Opcional pero Recomendado):**
    Este proyecto puede ejecutarse sin ninguna configuración, pero para producción, es altamente recomendable usar variables de entorno.

    Puedes generar una clave secreta segura usando Node.js:
    ```sh
    node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
    ```

### Ejecutar la Aplicación

Puedes iniciar el servidor de dos maneras:

**1. Usando la Configuración Predeterminada:**
Esto utilizará las credenciales de base de datos predeterminadas y el secreto JWT de respaldo.

```sh
go run main.go
```

**2. Usando Variables de Entorno (Recomendado):**
Esta es la forma más segura y flexible de ejecutar la aplicación.

```sh
# Ejemplo para Linux/macOS
export DB_USER="tu_usuario_db"
export DB_PASSWORD="tu_contraseña_db"
export JWT_SECRET="tu_clave_secreta_super_secreta_aqui"

go run main.go
```

El servidor se iniciará en `http://localhost:8080`.

## Endpoints de la API

-   `POST /users`: Crear un nuevo usuario.
-   `POST /login`: Autenticar un usuario y recibir un token JWT.
-   `GET /users/:id`: Obtener detalles del usuario (Ruta Protegida).

### Ejemplo de Uso

1.  **Crear un usuario:**
    ```sh
    curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{
      "Nombre": "Test", "ApellidoPaterno": "User", "Correo": "test@example.com", "Contraseña": "password123"
    }'
    ```

2.  **Iniciar sesión para obtener un token:**
    ```sh
    curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{
      "correo": "test@example.com", "contraseña": "password123"
    }'
    ```

3.  **Acceder a una ruta protegida:**
    Reemplaza `<TU_TOKEN_AQUI>` con el token del paso de inicio de sesión.
    ```sh
    curl http://localhost:8080/users/1 -H "Authorization: Bearer <TU_TOKEN_AQUI>"
    ```
