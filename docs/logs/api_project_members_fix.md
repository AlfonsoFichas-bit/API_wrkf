# Corrección API - Asignación de Tareas

## Problema Identificado
Los usuarios no pueden asignar tareas debido al error: **"assignment failed: user is not a member of this project"**

## Causa Raíz
Cuando se crea un proyecto, el **creador del proyecto no se agrega automáticamente como miembro** en la tabla `project_members`. La API requiere que los usuarios sean miembros del proyecto para poder asignar tareas.

### Consulta que Falla
```sql
SELECT * FROM "project_members"
WHERE user_id = 1 AND project_id = 1
ORDER BY "project_members"."id" LIMIT 1
```
**Resultado:** 0 filas encontradas

## Solución Recomendada
**Agregar automáticamente al creador del proyecto como miembro con rol `product_owner`** cuando se crea un nuevo proyecto.

## Implementación

### 1. Modificar el Handler de Creación de Proyectos
**Archivo:** `handlers/project_handler.go`

**Código actual (aproximado):**
```go
func CreateProject(c echo.Context) error {
    // ... código existente para crear proyecto ...

    project := models.Project{
        Name: req.Name,
        Description: req.Description,
        CreatedByID: userID, // ID del usuario que crea el proyecto
        // ... otros campos ...
    }

    if err := config.DB.Create(&project).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Could not create project",
        })
    }

    return c.JSON(http.StatusCreated, project)
}
```

**Código corregido:**
```go
func CreateProject(c echo.Context) error {
    userID, err := utils.GetUserIDFromContext(c)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{
            "error": err.Error(),
        })
    }

    // ... parsing del request ...

    project := models.Project{
        Name: req.Name,
        Description: req.Description,
        CreatedByID: userID,
        // ... otros campos ...
    }

    // Crear el proyecto
    if err := config.DB.Create(&project).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Could not create project",
        })
    }

    // 🔥 AGREGAR AL CREADOR COMO MIEMBRO DEL PROYECTO
    projectMember := models.ProjectMember{
        UserID:    userID,
        ProjectID: project.ID,
        Role:      "product_owner", // Rol por defecto para el creador
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := config.DB.Create(&projectMember).Error; err != nil {
        // Log del error pero no fallar la creación del proyecto
        log.Printf("Error adding creator as project member: %v", err)
    }

    return c.JSON(http.StatusCreated, project)
}
```

### 2. Verificar Modelo ProjectMember
**Archivo:** `models/project_member.go`

Asegurarse de que el modelo tenga los campos correctos:
```go
type ProjectMember struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id" gorm:"not null"`
    ProjectID uint      `json:"project_id" gorm:"not null"`
    Role      string    `json:"role" gorm:"not null"` // "product_owner", "scrum_master", "team_developer"
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    // Relaciones
    User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}
```

### 3. Verificar Migración de Base de Datos
**Archivo:** `migrations/xxxxx_create_project_members.go` o similar

Asegurarse de que la tabla `project_members` tenga las columnas correctas:
```go
func (m *CreateProjectMembersTable) Up() {
    m.CreateTable("project_members", func(t *migration.Table) {
        t.Column("id", "serial", map[string]string{"primary": "true"})
        t.Column("user_id", "integer", map[string]string{"null": "false"})
        t.Column("project_id", "integer", map[string]string{"null": "false"})
        t.Column("role", "varchar(50)", map[string]string{"null": "false"})
        t.Column("created_at", "timestamp", map[string]string{"null": "false"})
        t.Column("updated_at", "timestamp", map[string]string{"null": "false"})

        // Índices y constraints
        t.Index("idx_project_members_user_project", "user_id", "project_id")
        t.ForeignKey("user_id", "users", "id", "CASCADE", "CASCADE")
        t.ForeignKey("project_id", "projects", "id", "CASCADE", "CASCADE")
    })
}
```

## Verificación

### 1. Crear un Nuevo Proyecto
```bash
POST /api/projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Test Project",
  "description": "Testing automatic member addition"
}
```

### 2. Verificar que el Creador sea Miembro
```sql
SELECT * FROM project_members WHERE project_id = <new_project_id>;
-- Debería mostrar al creador como miembro con rol 'product_owner'
```

### 3. Probar Asignación de Tarea
Ahora debería funcionar asignar tareas en el proyecto recién creado.

## Archivos a Modificar
- `handlers/project_handler.go` - Agregar lógica para añadir creador como miembro
- `models/project_member.go` - Verificar modelo (si es necesario)
- `migrations/` - Verificar migración (si es necesario)

## Fecha de Implementación
Pendiente - 06 de noviembre de 2025

## Notas Adicionales
- Esta solución asegura que todos los creadores de proyectos tengan automáticamente permisos para gestionar sus proyectos
- El rol `product_owner` es apropiado para el creador del proyecto
- Si ya existen proyectos sin miembros, se necesitará un script de migración para añadir a los creadores como miembros