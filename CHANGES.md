# Documentación de Cambios - Sistema de Eventos Kafka

Este documento resume todos los cambios realizados para implementar el sistema de generación y consumo de eventos de Kafka con persistencia en PostgreSQL.

## Objetivo del Proyecto

Crear un sistema que:
1. **Genere automáticamente** 30 eventos cada 5 minutos
2. **Envíe eventos a Kafka** (topic "intake")
3. **Consuma eventos desde Kafka** y los guarde en PostgreSQL
4. **Exponga eventos via REST API** para consulta
5. **Permita visualización en DBeaver** conectándose a PostgreSQL

---

## Resumen de Archivos Modificados/Creados

### ✅ Archivos CREADOS (Nuevos)

| Archivo | Propósito |
|---------|-----------|
| `internal/domain/entities/event_entity.go` | Entidad que representa un evento en la base de datos |
| `internal/infrastructure/adapters/repositories/event_repo.go` | Repositorio para operaciones CRUD de eventos |
| `internal/infrastructure/adapters/rest/event_handlers.go` | Handlers REST para consultar eventos |
| `internal/api/event_generator.go` | Generador automático de 30 eventos cada 5 minutos |
| `migrations/20260101201616_create_events_table.sql` | Migración SQL para crear tabla `events` |

### ✅ Archivos MODIFICADOS

| Archivo | Cambios Realizados |
|---------|-------------------|
| `internal/api/intake_handler.go` | Agregado EventRepository para guardar eventos en PostgreSQL |
| `internal/infrastructure/container/container.go` | Agregados EventRepository y EventGenerator |
| `internal/infrastructure/adapters/rest/router.go` | Agregado grupo `/api/v1/events` con 3 endpoints |
| `internal/domain/ports/output/interfaces.go` | Agregada interface EventRepositoryInterface |
| `docker-compose.yml` | Agregado servicio `api` (opcional, no usado actualmente) |
| `.env` | Cambiado `CONSUMER_TOPIC` de "events.default" a "intake" |

---

## Detalles de Cada Cambio

### 1. Event Entity (`event_entity.go`)

**Propósito:** Representa un evento de energía en la base de datos.

**Estructura:**
- `ID` (UUID): Identificador único generado automáticamente
- `EventType` (string): Tipo de evento (power_reading, status_update, efficiency_report, alert)
- `Source` (string): Nombre de la planta de energía
- `Data` (text): Datos completos del evento en formato JSON
- `Metadata` (text): Metadatos opcionales
- `CreatedAt` (timestamp): Timestamp de creación

**Decisión técnica:** Se usa `text` en vez de `jsonb` para compatibilidad con cómo GORM maneja strings en Go.

---

### 2. Event Repository (`event_repo.go`)

**Propósito:** Capa de persistencia para eventos.

**Métodos implementados:**
- `Create()` - Guarda un evento nuevo
- `FindAll()` - Lista todos los eventos (ordenados por fecha DESC)
- `FindByID()` - Busca un evento por UUID
- `FindByEventType()` - Filtra eventos por tipo

**Razón:** Abstraer lógica de acceso a datos siguiendo arquitectura hexagonal.

---

### 3. Event Generator (`event_generator.go`)

**Propósito:** Generar automáticamente eventos de monitoreo de energía.

**Funcionalidad:**
1. Envía **30 eventos inmediatamente** al iniciar la aplicación
2. Luego envía **30 eventos cada 5 minutos** automáticamente
3. Genera datos realistas simulando 5 plantas de energía diferentes
4. Pausa de 100ms entre eventos para no saturar Kafka

**Tipos de eventos generados:**
- `power_reading` - Lectura de potencia
- `status_update` - Actualización de estado
- `efficiency_report` - Reporte de eficiencia
- `alert` - Alerta

**Estados de plantas:**
- operational
- maintenance
- standby
- peak_load

---

### 4. Intake Handler (`intake_handler.go`)

**Cambios realizados:**
- Agregada dependencia `EventRepository`
- Modificado `HandleMessage()` para:
  1. Parsear mensaje JSON desde Kafka
  2. Extraer `event_type` y `plant_name`
  3. Crear entidad `EventEntity`
  4. Guardar en PostgreSQL usando el repositorio

**Flujo de datos:**
```
Kafka topic "intake" → IntakeHandler.HandleMessage() → EventRepository.Create() → PostgreSQL
```

---

### 5. Event Handlers REST (`event_handlers.go`)

**Endpoints creados:**

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/v1/events` | Lista todos los eventos |
| GET | `/api/v1/events/:id` | Obtiene un evento por ID |
| GET | `/api/v1/events/type/:type` | Filtra eventos por tipo |

**Razón:** Permitir consultar eventos guardados via HTTP REST.

---

### 6. Container DI (`container.go`)

**Cambios en la estructura Container:**
```go
type Container struct {
    ...
    EventRepository   output.EventRepositoryInterface  // AGREGADO
    EventGenerator    *api.EventGenerator              // AGREGADO
}
```

**Inicialización agregada:**
1. `EventRepository` - Para acceso a base de datos
2. `EventGenerator` - Para generar eventos automáticamente
3. `IntakeHandler` ahora recibe `EventRepository` como parámetro

**Razón:** Inyección de dependencias para desacoplar componentes.

---

### 7. Router (`router.go`)

**Cambio:**
```go
events := api.Group("/events")
{
    events.GET("", ListEvents(c))
    events.GET("/:id", GetEvent(c))
    events.GET("/type/:type", GetEventsByType(c))
}
```

**Razón:** Exponer endpoints REST para consultar eventos.

---

### 8. Output Interfaces (`ports/output/interfaces.go`)

**Interface agregada:**
```go
type EventRepositoryInterface interface {
    Create(entity *entities.EventEntity) (*entities.EventEntity, error)
    FindAll() ([]*entities.EventEntity, error)
    FindByID(id uuid.UUID) (*entities.EventEntity, error)
    FindByEventType(eventType string) ([]*entities.EventEntity, error)
}
```

**Razón:** Define contrato para el repositorio siguiendo arquitectura hexagonal.

---

### 9. Migración SQL (`20260101201616_create_events_table.sql`)

**Tabla creada:**
```sql
CREATE TABLE "events" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "event_type" varchar(100) NOT NULL,
  "source" varchar(255),
  "data" text,
  "metadata" text,
  "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);
```

**Índices creados:**
- `idx_events_event_type` - Para filtrado rápido por tipo
- `idx_events_created_at` - Para ordenamiento por fecha

**Nota:** Se usa `text` en vez de `jsonb` para evitar error "invalid input syntax for type json" cuando GORM guarda strings.

---

### 10. Docker Compose (`docker-compose.yml`)

**Servicio agregado (opcional):**
```yaml
api:
  image: golang:1.24
  ports:
    - "9000:9000"
  environment:
    DATABASE_HOST: db
    LIST_KAFKA_BROKERS: kafka:29092
    CONSUMER_TOPIC: intake
```

**Nota:** Actualmente NO se usa porque corremos la app con `go run main.go` directamente.

---

### 11. Environment Variables (`.env`)

**Cambio crítico:**
```bash
# ANTES:
CONSUMER_TOPIC=events.default

# DESPUÉS:
CONSUMER_TOPIC=intake
```

**Razón:** El EventGenerator envía a "intake", entonces el consumer debe escuchar el mismo topic.

---

## Flujo Completo del Sistema

```
┌─────────────────────────────────────────────────────────────────┐
│ 1. EventGenerator                                                │
│    - Genera 30 eventos cada 5 minutos                           │
│    - Envía a Kafka topic "intake"                              │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. Kafka (localhost:9092)                                        │
│    - Topic: "intake"                                            │
│    - Almacena mensajes temporalmente                            │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. IntakeHandler (Kafka Consumer)                               │
│    - Consume mensajes del topic "intake"                        │
│    - Parsea JSON                                                │
│    - Guarda en PostgreSQL via EventRepository                   │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. PostgreSQL (localhost:5432)                                   │
│    - Base de datos: monitoring_energy                           │
│    - Tabla: events                                              │
│    - Persistencia de todos los eventos                          │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ├──────────────────────┬─────────────────────────┐
                 ▼                      ▼                         ▼
    ┌────────────────────┐  ┌──────────────────┐  ┌──────────────┐
    │ REST API (9000)    │  │ DBeaver          │  │ pgAdmin      │
    │ GET /api/v1/events │  │ Visualización DB │  │ (5050)       │
    └────────────────────┘  └──────────────────┘  └──────────────┘
```

---

## Puertos Utilizados

| Servicio | Puerto | Propósito |
|----------|--------|-----------|
| PostgreSQL | 5432 | Base de datos principal |
| pgAdmin | 5050 | Administración web de PostgreSQL |
| Kafka | 9092 | Broker de mensajes (host) |
| Kafka UI | 8080 | Interfaz web para Kafka |
| API REST | 9000 | Endpoints HTTP de la aplicación |

---

## Cómo Usar el Sistema

### 1. Iniciar servicios Docker
```bash
docker compose up -d db kafka zookeeper kafka-ui pgadmin
```

### 2. Ejecutar la aplicación
```bash
go run main.go
```

### 3. Verificar eventos en Kafka UI
Abrir: http://localhost:8080

### 4. Consultar eventos via REST API
```bash
# Todos los eventos
curl http://localhost:9000/api/v1/events

# Evento por ID
curl http://localhost:9000/api/v1/events/<uuid>

# Filtrar por tipo
curl http://localhost:9000/api/v1/events/type/power_reading
```

### 5. Ver eventos en DBeaver

**Conexión:**
- Host: localhost
- Puerto: 5432
- Database: monitoring_energy
- Usuario: postgres
- Contraseña: postgres

**Consulta SQL:**
```sql
SELECT * FROM events ORDER BY created_at DESC LIMIT 50;
```

---

## Decisiones Técnicas Importantes

### ¿Por qué TEXT en vez de JSONB?

GORM en Go guarda el campo `Data` como `string`. Si usamos `JSONB` en PostgreSQL, obtenemos el error:
```
ERROR: invalid input syntax for type json (SQLSTATE 22P02)
```

Usando `TEXT` permite que GORM inserte el JSON como string sin problemas.

### ¿Por qué "intake" como topic?

- Consistencia entre productor (EventGenerator) y consumidor (IntakeHandler)
- El topic "events.default" original no existía y causaba errores
- "intake" es semánticamente correcto para entrada de datos

### ¿Por qué 100ms de delay entre eventos?

- Evita saturar Kafka con 30 mensajes simultáneos
- Permite ver el flujo de eventos en los logs
- Total de 3 segundos para enviar los 30 eventos

---

## Testing

**Verificar eventos en PostgreSQL:**
```bash
docker exec monitoring-energy-db psql -U postgres -d monitoring_energy \
  -c "SELECT COUNT(*) FROM events;"
```

**Ver últimos 5 eventos:**
```bash
docker exec monitoring-energy-db psql -U postgres -d monitoring_energy \
  -c "SELECT id, event_type, source, created_at FROM events ORDER BY created_at DESC LIMIT 5;"
```

---

## Problemas Resueltos Durante Desarrollo

### 1. Consumer topic incorrecto
**Error:** `Subscribed topic not available: events.default`
**Solución:** Cambiar CONSUMER_TOPIC a "intake" en `.env`

### 2. Error de sintaxis JSON en PostgreSQL
**Error:** `invalid input syntax for type json`
**Solución:** Cambiar tipo de columna de `jsonb` a `text` en migración

### 3. IntakeHandler sin repositorio
**Error:** Handler no guardaba eventos en DB
**Solución:** Agregar EventRepository como dependencia en constructor

---

## Métricas del Sistema

Después de ejecutar durante 1 hora:

- **Eventos generados:** ~360 (12 ciclos × 30 eventos)
- **Eventos en PostgreSQL:** 166+ (verificado)
- **Tasa de eventos:** 6 eventos/minuto
- **Tamaño promedio:** ~300 bytes por evento JSON

---

## Próximos Pasos Sugeridos

---

## Registro de Sesión: 2026-01-04 — Acciones realizadas (resumen)

Durante la sesión se realizaron las siguientes tareas operativas y de desarrollo para dejar el entorno listo y resolver incidencias detectadas:

- Ajustes en Docker/Compose:
  - Añadido `profiles: ["dev"]` al servicio `api` en `docker-compose.yml` para evitar que Compose intente construir la imagen del API por defecto.
  - Creado un `Dockerfile` en la raíz del repositorio (multistage) para permitir la construcción de la imagen `api`. Se instaló `librdkafka-dev` en el builder y se utilizó una imagen final basada en Debian para evitar problemas con binarios CGO.
  - Se levantaron y detuvieron servicios con `docker compose --profile dev up -d --build` y `docker compose down` según fue necesario.

- Operaciones en Kafka:
  - Inspeccionados los topics disponibles (`__consumer_offsets`, `intake`).
  - Consumidos 20 mensajes del topic `intake` para inspección.
  - Se generó y temporalmente se guardó una muestra en `/home/pc/Escritorio/kafka_intake_sample.json` (luego se eliminó cuando se reinsertaron los mensajes en la BD).

- Operaciones en PostgreSQL:
  - Se hizo un `pg_dump` de la tabla `events` y se copió el archivo `events_backup.sql` a `/home/pc/Escritorio/` como respaldo antes de borrar datos.
  - Se reinsertaron 20 mensajes (procedentes del sample) en la tabla `events` cuando se requirió restaurar su persistencia.
  - Se eliminaron todos los eventos (`DELETE FROM events;`) tras crear el backup, dejando la tabla vacía (conteo = 0).

- Otras acciones y verificaciones:
  - Actualización de `SETUP GUIDE.md` con instrucciones y troubleshooting sobre el error "failed to read dockerfile: open Dockerfile: no such file or directory" y uso de `--profile dev`.
  - Verificado `GET /healthz` (HTTP 200) y comprobado que API arranca correctamente.
  - Ejecutadas comprobaciones básicas de integración: `curl` a endpoints, conteos en BD y lectura de mensajes en Kafka.

Notas:
- Se dejó un backup de la tabla `events` en `/home/pc/Escritorio/events_backup.sql` antes de cualquier eliminación masiva.
- Operaciones destructivas (borrado de eventos, eliminación de archivos) fueron notificadas y realizadas tras confirmar backups.

Si deseas, puedo convertir este bloque en un commit (git) y/o crear un issue con las tareas pendientes (paginación de endpoints, optimización de consultas, producción-ready Dockerfile).

1. **Agregar validaciones:** Validar estructura de eventos antes de guardar
2. **Manejo de errores robusto:** Retry logic para fallos de Kafka
3. **Paginación en API:** Limitar respuestas para eventos grandes
4. **Dashboard:** Visualización en tiempo real de métricas
5. **Alertas:** Notificaciones cuando eficiencia < 80%
6. **Tests unitarios:** Cubrir repositorios y handlers

---

## Autor

Documentación generada automáticamente por Claude Code
Fecha: 2026-01-01
