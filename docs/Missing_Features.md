# Análisis de Paridad de Features: Deno (Anterior) vs. Go (Actual) - Análisis Profundo

Este documento describe las funcionalidades presentes en la API original o inferidas a partir de los modelos de datos de la nueva API de Go, que aún no han sido completamente implementadas.

## 1. Módulo de Rúbricas (CRÍTICO Y COMPLEJO)
**Confirmado a través de búsqueda de código.** El proyecto anterior contenía un módulo completo para la creación, gestión y utilización de rúbricas de evaluación. Esta es la funcionalidad más grande y crítica que falta.

**Modelos de Datos a Crear:**
*   `Rubric`: Debe contener campos como `name`, `description`, `status` (`DRAFT`, `ACTIVE`, `ARCHIVED`), `isTemplate` (booleano).
*   `RubricCriterion`: Parte de una `Rubric`. Debe contener `title`, `description`, `maxPoints`, y una lista de `Levels`.
*   `RubricCriterionLevel`: Parte de un `Criterion`. Debe contener `score` y `description`.

**Acciones Requeridas:**
*   Crear los modelos `rubric.go`, `rubric_criterion.go`, etc., en la carpeta `models/`.
*   Implementar un `RubricService` y `RubricHandler`.
*   Añadir rutas para el CRUD completo y funcionalidades adicionales:
    *   `GET /api/rubrics`: Con filtros por `projectId`, `userId`, y `isTemplate`.
    *   `POST /api/rubrics`: Crear una nueva rúbrica.
    *   `GET /api/rubrics/:id`: Obtener detalles de una rúbrica.
    *   `PUT /api/rubrics/:id`: Actualizar una rúbrica.
    *   `DELETE /api/rubrics/:id`: Eliminar una rúbrica.
    *   `POST /api/rubrics/:id/duplicate`: Duplicar una rúbrica.

## 2. Módulo de Evaluaciones (Crítico)
Esta funcionalidad depende directamente del Módulo de Rúbricas. Una evaluación es, en esencia, la aplicación de una rúbrica a un entregable o miembro del equipo.

**Acciones Requeridas:**
*   Asegurar que el modelo `evaluation.go` tenga un campo `RubricID`.
*   Implementar `EvaluationService` y `EvaluationHandler`.
*   Añadir rutas para el CRUD de evaluaciones.

## 3. Módulo de Reportes (Crítico)
El modelo `models/reporting.go` indica la necesidad de generar reportes, pero no hay lógica de negocio ni endpoints.

**Acciones Requeridas:**
*   Implementar `ReportingService` y `ReportingHandler`.
*   Añadir rutas para generar reportes (ej. `GET /api/projects/:id/reports/burndown`).

## 4. Gestión de Roles y Permisos Avanzados (Importante)
El sistema de permisos actual es binario (admin/no-admin). El modelo `roles.go` sugiere la necesidad de un sistema más granular.

**Acciones Requeridas:**
*   Expandir `UserService` para manejar roles.
*   Modificar `POST /api/projects/:id/members` para asignar roles.
*   Implementar middleware de autorización basado en roles.

## 5. Búsqueda Avanzada tipo JQL (Importante)
Funcionalidad de la API original que permitía consultas complejas.

**Acciones Requeridas:**
*   Diseñar e implementar un endpoint de búsqueda (ej. `GET /api/search?q={query}`).

## 6. Gestión Completa de Comentarios de Tareas (Mejora)
Actualmente solo se pueden añadir comentarios.

**Acciones Requeridas:**
*   Añadir endpoints para `GET`, `PUT`, `DELETE` comentarios.

## 7. API para Historial de Tareas (Mejora)
El historial de cambios de las tareas no es accesible.

**Acciones Requeridas:**
*   Añadir endpoint `GET /api/tasks/:id/history`.
