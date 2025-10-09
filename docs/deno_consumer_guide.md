# Guía para Consumir la API de Go con Deno

Este documento proporciona una guía práctica y ejemplos de código para interactuar con la API de gestión de proyectos desde un cliente escrito en Deno/TypeScript.

## 1. Prerrequisitos

-   Tener [Deno](https://deno.land/) instalado.
-   El servidor de la API de Go debe estar en ejecución en `http://localhost:8080`.

## 2. Flujo de Trabajo Básico

El flujo para interactuar con la API es el siguiente:

1.  **Autenticarse:** Enviar una petición `POST` al endpoint `/login` con las credenciales de usuario para obtener un token JWT.
2.  **Almacenar el Token:** Guardar el token JWT recibido en una variable.
3.  **Realizar Peticiones Autenticadas:** Para acceder a los endpoints protegidos, incluir el token JWT en el encabezado `Authorization` de cada petición, con el formato `Bearer <token>`.

---

## 3. Implementación en Deno

### Paso 1: Autenticación y Obtención del Token

Primero, creamos una función para obtener el token. Usaremos las credenciales del administrador por defecto.

```typescript
// api_client.ts

const API_BASE_URL = "http://localhost:8080";

interface LoginResponse {
  token: string;
}

/**
 * Authenticates with the API and returns a JWT token.
 */
async function login(): Promise<string | null> {
  const loginUrl = `${API_BASE_URL}/login`;
  const credentials = {
    correo: "admin@example.com",
    contraseña: "admin123",
  };

  try {
    const response = await fetch(loginUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      console.error(`Error en el login: ${response.status} ${response.statusText}`);
      const errorBody = await response.json();
      console.error("Detalles:", errorBody);
      return null;
    }

    const data: LoginResponse = await response.json();
    console.log("¡Login exitoso! Token recibido.");
    return data.token;
  } catch (error) {
    console.error("Error de red o al realizar la petición de login:", error);
    return null;
  }
}
```

### Paso 2: Realizar una Petición Autenticada (GET)

Una vez que tenemos el token, podemos usarlo para hacer peticiones a endpoints protegidos, como el que obtiene la lista de proyectos.

```typescript
// api_client.ts (continuación)

/**
 * Fetches all projects using the provided JWT token.
 * @param token The JWT token for authentication.
 */
async function getAllProjects(token: string): Promise<void> {
  const projectsUrl = `${API_BASE_URL}/api/projects`;

  try {
    const response = await fetch(projectsUrl, {
      method: "GET",
      headers: {
        "Authorization": `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      console.error(`Error al obtener proyectos: ${response.status} ${response.statusText}`);
      return;
    }

    const projects = await response.json();
    console.log("\n--- Proyectos Recuperados ---");
    console.log(projects);
  } catch (error) {
    console.error("Error de red al obtener proyectos:", error);
  }
}
```

### Paso 3: Crear un Recurso (POST)

Ahora, un ejemplo de cómo crear un nuevo proyecto enviando datos en el cuerpo de la petición.

```typescript
// api_client.ts (continuación)

/**
 * Creates a new project.
 * @param token The JWT token for authentication.
 */
async function createProject(token: string): Promise<void> {
  const projectsUrl = `${API_BASE_URL}/api/projects`;
  const newProject = {
    Name: "Proyecto desde Deno",
    Description: "Este es un proyecto de prueba creado con un cliente Deno.",
  };

  try {
    const response = await fetch(projectsUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
      },
      body: JSON.stringify(newProject),
    });

    if (response.status !== 201) {
      console.error(`Error al crear el proyecto: ${response.status} ${response.statusText}`);
      const errorBody = await response.json();
      console.error("Detalles:", errorBody);
      return;
    }

    const createdProject = await response.json();
    console.log("\n--- Nuevo Proyecto Creado ---");
    console.log(createdProject);
  } catch (error) {
    console.error("Error de red al crear el proyecto:", error);
  }
}
```

### Paso 4: Script Completo

Aquí está el script completo que puedes guardar como `api_client.ts` y ejecutar con Deno.

```typescript
// api_client.ts

const API_BASE_URL = "http://localhost:8080";

// --- Tipos ---
interface LoginResponse {
  token: string;
}

// --- Funciones de API ---

async function login(): Promise<string | null> {
  // ... (código de la función login de arriba)
  const loginUrl = `${API_BASE_URL}/login`;
  const credentials = {
    correo: "admin@example.com",
    contraseña: "admin123",
  };

  try {
    const response = await fetch(loginUrl, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      console.error(`Error en el login: ${response.status} ${response.statusText}`);
      return null;
    }
    const data: LoginResponse = await response.json();
    console.log("¡Login exitoso! Token recibido.");
    return data.token;
  } catch (error) {
    console.error("Error de red o al realizar la petición de login:", error);
    return null;
  }
}

async function getAllProjects(token: string): Promise<void> {
  // ... (código de la función getAllProjects de arriba)
  const projectsUrl = `${API_BASE_URL}/api/projects`;
  try {
    const response = await fetch(projectsUrl, {
      method: "GET",
      headers: { "Authorization": `Bearer ${token}` },
    });
    if (!response.ok) {
      console.error(`Error al obtener proyectos: ${response.status} ${response.statusText}`);
      return;
    }
    const projects = await response.json();
    console.log("\n--- Proyectos Recuperados ---");
    console.log(projects);
  } catch (error) {
    console.error("Error de red al obtener proyectos:", error);
  }
}

async function createProject(token: string): Promise<void> {
  // ... (código de la función createProject de arriba)
  const projectsUrl = `${API_BASE_URL}/api/projects`;
  const newProject = {
    Name: "Proyecto desde Deno",
    Description: "Este es un proyecto de prueba creado con un cliente Deno.",
  };
  try {
    const response = await fetch(projectsUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
      },
      body: JSON.stringify(newProject),
    });
    if (response.status !== 201) {
      console.error(`Error al crear el proyecto: ${response.status} ${response.statusText}`);
      return;
    }
    const createdProject = await response.json();
    console.log("\n--- Nuevo Proyecto Creado ---");
    console.log(createdProject);
  } catch (error) {
    console.error("Error de red al crear el proyecto:", error);
  }
}


// --- Ejecución Principal ---
async function main() {
  console.log("Iniciando cliente de API...");

  const token = await login();

  if (token) {
    await getAllProjects(token);
    await createProject(token);
    await getAllProjects(token); // Volver a llamar para ver el nuevo proyecto en la lista
  } else {
    console.log("No se pudo obtener el token. Finalizando el script.");
  }
}

// Ejecutar el script principal
main();

```

### Cómo Ejecutar el Script

1.  Guarda el código anterior en un archivo llamado `api_client.ts`.
2.  Abre tu terminal y navega hasta el directorio donde guardaste el archivo.
3.  Ejecuta el script con Deno, otorgando los permisos de red necesarios:

    ```sh
    deno run --allow-net=localhost:8080 api_client.ts
    ```

Verás en la consola el resultado de cada paso: el mensaje de login, la lista de proyectos, el nuevo proyecto creado y finalmente la lista actualizada de proyectos.
