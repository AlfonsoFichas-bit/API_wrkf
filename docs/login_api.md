# Documentación de API: Login de Usuario

Este documento detalla cómo utilizar el endpoint de autenticación para obtener un token de acceso.

---

## Autenticación de Usuario

Permite a un usuario autenticarse en el sistema proporcionando sus credenciales (correo y contraseña) para recibir un token JWT que deberá ser utilizado en las cabeceras de las solicitudes a endpoints protegidos.

- **URL**: `/login`
- **Método**: `POST`
- **Autenticación Requerida**: No

### Cuerpo de la Solicitud (Request Body)

El cuerpo de la solicitud debe ser un objeto JSON con la siguiente estructura:

```json
{
  "correo": "tu_correo@example.com",
  "contraseña": "tu_contraseña"
}
```

**Parámetros:**

| Campo        | Tipo   | Descripción                  | Requerido |
|--------------|--------|------------------------------|-----------|
| `correo`     | string | La dirección de correo del usuario. | Sí        |
| `contraseña` | string | La contraseña del usuario.   | Sí        |

---

### Respuesta Exitosa (Success Response)

Si las credenciales son válidas, el servidor responderá con un código `200 OK` y un cuerpo JSON que contiene el token JWT.

- **Código**: `200 OK`
- **Contenido**:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

---

### Respuestas de Error (Error Responses)

- **Código**: `400 Bad Request`
  - **Causa**: El cuerpo de la solicitud no es un JSON válido o faltan campos requeridos.
  - **Contenido**:
    ```json
    {
      "error": "Invalid input"
    }
    ```

- **Código**: `401 Unauthorized`
  - **Causa**: Las credenciales proporcionadas (correo o contraseña) son incorrectas.
  - **Contenido**:
    ```json
    {
      "error": "invalid credentials"
    }
    ```

---

### Ejemplo de Uso con cURL

```bash
curl -X POST \
  http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{
    "correo": "admin@example.com",
    "contraseña": "admin123"
  }'
```
