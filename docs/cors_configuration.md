# Configuración de CORS en Echo (Go)

Para permitir que el frontend (ej. corriendo en `localhost:5173`) haga solicitudes a la API Go (ej. en `localhost:8080`), necesitas configurar CORS (Cross-Origin Resource Sharing) en tu servidor Echo.

## Pasos para configurar CORS

### 1. Instalar el paquete CORS

Asegúrate de tener instalado el paquete `github.com/rs/cors` en tu proyecto Go:

```bash
go get github.com/rs/cors
```

### 2. Configurar CORS en tu servidor Echo

En tu archivo `main.go` o donde inicialices el servidor Echo, agrega el middleware de CORS antes de definir las rutas.

Ejemplo de configuración básica:

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/rs/cors"
)

func main() {
    e := echo.New()

    // Configurar CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:5173", // URL del frontend en desarrollo
            "http://localhost:3000", // Si usas otro puerto
            // Agrega más orígenes según necesites (ej. para producción)
        },
        AllowedMethods: []string{
            "GET",
            "POST",
            "PUT",
            "DELETE",
            "OPTIONS",
        },
        AllowedHeaders: []string{
            "Content-Type",
            "Authorization",
            "X-Requested-With",
        },
        AllowCredentials: true, // Si necesitas enviar cookies o headers de auth
    })

    // Aplicar el middleware CORS al servidor Echo
    e.Use(echo.WrapMiddleware(c.Handler))

    // Opcional: Otros middlewares
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Definir rutas aquí
    // e.POST("/login", loginHandler)
    // e.POST("/api/admin/users", createUserHandler)

    e.Logger.Fatal(e.Start(":8080"))
}
```

### 3. Explicación de las opciones

- **AllowedOrigins**: Lista de orígenes permitidos. Agrega la URL de tu frontend (ej. `http://localhost:5173` para Vite).
- **AllowedMethods**: Métodos HTTP permitidos.
- **AllowedHeaders**: Headers que se permiten en las solicitudes.
- **AllowCredentials**: Permite enviar credenciales (importante si usas autenticación con cookies o Authorization headers).

### 4. Para producción

En producción, reemplaza los orígenes con la URL real de tu frontend (ej. `https://tudominio.com`).

### 5. Verificar

- Inicia tu servidor Go.
- Abre el navegador en `http://localhost:5173` y prueba el login/registro.
- El error de CORS debería desaparecer.

Si tienes problemas, verifica que el puerto de la API sea `8080` y el del frontend `5173`, y ajusta según corresponda.