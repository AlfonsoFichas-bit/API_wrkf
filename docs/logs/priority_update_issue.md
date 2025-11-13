# Problema con la Actualización de Prioridad en User Stories

## Descripción del Problema

Al editar una historia de usuario (HU) desde el frontend, cuando se cambia la prioridad a "high", el backend responde con "medium" en lugar del valor enviado.

## Logs del Navegador

```
Enviando datos actualizados de la user story:
Object { Title: "entrevista", Description: "realizacion de entrevista", AcceptanceCriteria: "manuscrita", Priority: "high", Points: 1 }
Prioridad actual: high
ID de la user story: 1
URL del endpoint: http://localhost:8080/api/userstories/1
Headers:
Object { "Content-Type": "application/json", Authorization: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI3NDM4MTMsImlhdCI6MTc2MjQ4NDYxMywibmFtIjoiQWRtaW4iLCJyb2wiOiJhZG1pbiIsInN1YiI6MX0.fEzEGhNM4V6rj7yKqgUHHgb2gHSw91YGEaL2E-tysww" }
Status de respuesta: 200
Headers de respuesta:
Object { "content-length": "999", "content-type": "application/json" }
Respuesta del backend:
Object { ID: 1, Title: "entrevista", Description: "realizacion de entrevista", AcceptanceCriteria: "manuscrita", Priority: "medium", Status: "backlog", Points: null, ProjectID: 1, Project: {…}, SprintID: 1, … }
Error: Priority not updated correctly. Sent: high Received: medium
```

## Análisis del Código Frontend

El código del frontend en `UserStoriesIsland.tsx` está correcto:

- Envía correctamente `Priority: priority` en el objeto `updatedUserStory`
- Los logs muestran que se envía "high"
- La validación del frontend confirma que se envió "high" pero se recibió "medium"

## Posibles Causas en el Backend (Go/Echo)

1. **Problema de binding JSON**: El manejador `PUT /api/userstories/:storyId` no vincula correctamente el campo `Priority` del JSON entrante.

2. **Lógica de actualización incompleta**: La consulta de actualización no incluye el campo `Priority`.

3. **Valores por defecto**: Hay código que establece `Priority` en "medium" como valor predeterminado.

4. **Middleware de validación**: Algún middleware sobrescribe el valor de `Priority`.

5. **Problemas de base de datos**: Constraints o triggers que cambian el valor.

## Código del Manejador Esperado

```go
func UpdateUserStory(c echo.Context) error {
    storyId := c.Param("storyId")
    
    var updateData struct {
        Title            string `json:"Title,omitempty"`
        Description      string `json:"Description,omitempty"`
        AcceptanceCriteria string `json:"AcceptanceCriteria,omitempty"`
        Priority         string `json:"Priority,omitempty"`
        Points           *int   `json:"Points,omitempty"`
    }
    
    if err := c.Bind(&updateData); err != nil {
        return c.JSON(400, map[string]string{"error": "Solicitud inválida"})
    }
    
    // IMPORTANTE: La lógica de actualización debe incluir TODOS los campos
    // Ejemplo con GORM:
    // if err := db.Model(&userStory).Where("id = ?", storyId).Updates(updateData).Error; err != nil {
    //     return c.JSON(500, map[string]string{"error": "Error al actualizar"})
    // }
    
    return c.JSON(200, userStory)
}
```

## Pasos para Depurar

1. **Revisar el manejador PUT**: Verificar que `c.Bind()` incluya `Priority` y que la actualización lo use.

2. **Logs del backend**: Agregar logs en el backend para ver qué valores se reciben y se actualizan.

3. **Prueba directa con curl**:
   ```bash
   curl -X PUT http://localhost:8080/api/userstories/1 \
     -H "Authorization: Bearer TU_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"Priority": "high"}'
   ```

4. **Verificar la base de datos**: Comprobar si hay triggers o constraints que cambien el valor.

5. **Revisar middlewares**: Buscar código que modifique la prioridad antes de guardar.

## Solución Sugerida

Asegurarse de que el backend:
- Vincule correctamente el JSON con el campo `Priority`
- Incluya `Priority` en la consulta de actualización
- No tenga valores por defecto que sobrescriban el input del usuario

Si el problema persiste, revisar los logs del backend para ver exactamente qué valores se procesan.