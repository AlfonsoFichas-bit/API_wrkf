#!/usr/bin/fish

echo "--- 🚀 INICIANDO PRUEBA COMPLETA DEL WORKFLOW ---"

# --- PASO 1: LOGIN ---
echo "(1/8) Obteniendo token de Admin..."
set ADMIN_TOKEN_JSON (curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"correo": "admin@example.com", "contraseña": "admin123"}')
set ADMIN_TOKEN (echo $ADMIN_TOKEN_JSON | jq -r .token)

if test -z "$ADMIN_TOKEN"
    echo "❌ ERROR: No se pudo obtener el token de Admin. Saliendo." 
    exit 1
end
echo "✅ Token de Admin guardado."

# --- PASO 2: CREAR PROYECTO ---
echo "(2/8) Creando proyecto..."
set PROJECT_JSON (curl -s -X POST http://localhost:8080/api/projects \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"Name": "Proyecto de Prueba Automatizada"}')
set PROJECT_ID (echo $PROJECT_JSON | jq -r .ID)

if test -z "$PROJECT_ID" -o "$PROJECT_ID" = "null"
    echo "❌ ERROR: No se pudo crear el proyecto. Saliendo." 
    exit 1
end
echo "✅ Proyecto creado con ID: $PROJECT_ID"

# --- PASO 3: CREAR USER STORY ---
echo "(3/8) Creando User Story..."
set STORY_JSON (curl -s -X POST http://localhost:8080/api/projects/$PROJECT_ID/userstories \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title": "Historia para Automatización"}')
set STORY_ID (echo $STORY_JSON | jq -r .ID)

if test -z "$STORY_ID" -o "$STORY_ID" = "null"
    echo "❌ ERROR: No se pudo crear la User Story. Saliendo." 
    exit 1
end
echo "✅ User Story creada con ID: $STORY_ID"

# --- PASO 4: CREAR TAREA ---
echo "(4/8) Creando Tarea..."
set TASK_JSON (curl -s -X POST http://localhost:8080/api/userstories/$STORY_ID/tasks \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title": "Tarea para Mover"}')
set TASK_ID (echo $TASK_JSON | jq -r .ID)

if test -z "$TASK_ID" -o "$TASK_ID" = "null"
    echo "❌ ERROR: No se pudo crear la Tarea. Saliendo." 
    exit 1
end
echo "✅ Tarea creada con ID: $TASK_ID"

# --- PASO 5: PROBAR ACTUALIZACIÓN DE ESTADO (ÉXITO) ---
echo "(5/8) Actualizando estado a 'in_progress'..."
set UPDATE_RESPONSE (curl -s -X PUT http://localhost:8080/api/tasks/$TASK_ID/status \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"status": "in_progress"}')
set NEW_STATUS (echo $UPDATE_RESPONSE | jq -r .Status)

if test "$NEW_STATUS" = "in_progress"
    echo "✅ ÉXITO: El estado de la tarea ahora es 'in_progress'."
else
    echo "❌ FALLO: El estado de la tarea no se actualizó correctamente."
    echo "Respuesta recibida:" 
    echo $UPDATE_RESPONSE | jq
    exit 1
end

# --- PASO 6: PROBAR ACTUALIZACIÓN DE ESTADO (FALLO) ---
echo "(6/8) Intentando actualizar a estado inválido 'trabajando'..."
set INVALID_UPDATE_STATUS (curl -s -o /dev/null -w "%{http_code}" -X PUT http://localhost:8080/api/tasks/$TASK_ID/status \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"status": "trabajando"}')

if test "$INVALID_UPDATE_STATUS" = "400"
    echo "✅ ÉXITO: La API rechazó correctamente el estado inválido con un código 400."
else
    echo "❌ FALLO: La API no devolvió un error 400 para un estado inválido. Código recibido: $INVALID_UPDATE_STATUS"
    exit 1
end

# --- PASO 7: PROBAR ASIGNACIÓN ---
echo "(7/8) Creando usuario 'Diana Dev' para asignación..."
set DEV_JSON (curl -s -X POST http://localhost:8080/api/admin/users \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"Nombre": "Diana", "ApellidoPaterno": "Dev", "Correo": "diana.dev@example.com", "Contraseña": "password123"}')
set DEV_ID (echo $DEV_JSON | jq -r .ID)
echo "Añadiendo a Diana al proyecto..."
curl -s -X POST http://localhost:8080/api/admin/projects/$PROJECT_ID/members \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"userId": '$DEV_ID', "role": "team_developer"}' > /dev/null

echo "Asignando tarea a Diana..."
set ASSIGN_RESPONSE (curl -s -X PUT http://localhost:8080/api/tasks/$TASK_ID/assign \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"userId": '$DEV_ID'}')
set ASSIGNED_ID (echo $ASSIGN_RESPONSE | jq -r .AssignedToID)

if test "$ASSIGNED_ID" = "$DEV_ID"
    echo "✅ ÉXITO: La tarea fue asignada correctamente a Diana."
else
    echo "❌ FALLO: La tarea no fue asignada."
    echo "Respuesta recibida:" 
    echo $ASSIGN_RESPONSE | jq
    exit 1
end

# --- PASO 8: LIMPIEZA ---
echo "(8/8) Limpiando recursos (borrando proyecto)..."
set DELETE_STATUS (curl -s -o /dev/null -w "%{http_code}" -X DELETE http://localhost:8080/api/projects/$PROJECT_ID \
    -H "Authorization: Bearer $ADMIN_TOKEN")

if test "$DELETE_STATUS" = "204"
    echo "✅ ÉXITO: El proyecto y todos sus recursos asociados fueron eliminados."
else
    echo "❌ FALLO: No se pudo eliminar el proyecto. Código recibido: $DELETE_STATUS"
    exit 1
end

echo "
--- 🎉 PRUEBA COMPLETA Y EXITOSA --- "
