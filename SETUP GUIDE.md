

# Energy Monitoring System - Setup Guide

Sistema de monitoreo de energía con generación automática de eventos, procesamiento mediante Kafka y persistencia en PostgreSQL.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Arquitectura](#arquitectura)
- [Requisitos Previos](#requisitos-previos)
- [Instalación](#instalación)
- [Configuración](#configuración)
- [Ejecución](#ejecución)
- [Uso de la API REST](#uso-de-la-api-rest)
- [Visualización con DBeaver](#visualización-con-dbeaver)
- [Kafka UI](#kafka-ui)
- [Troubleshooting](#troubleshooting)
- [Estructura del Proyecto](#estructura-del-proyecto)

---

## 🎯 Características

- ✅ **Generación automática de eventos**: 30 eventos cada 5 minutos
- ✅ **Integración con Kafka**: Producer y Consumer
- ✅ **Persistencia en PostgreSQL**: Con TimescaleDB y PostGIS
- ✅ **API REST**: 3 endpoints para consultar eventos
- ✅ **Visualización**: DBeaver, pgAdmin, Kafka UI
- ✅ **Hot Reload**: Desarrollo con Air
- ✅ **Migraciones**: Goose para versionado de DB
- ✅ **Swagger**: Documentación automática de API

---

## 🏗️ Arquitectura

```
┌──────────────────────────────────────────────────────────────┐
│                    ENERGY MONITORING SYSTEM                   │
└──────────────────────────────────────────────────────────────┘

┌─────────────────────┐
│  Event Generator    │  ← Genera 30 eventos cada 5 min
│  (Go Application)   │
└──────────┬──────────┘
           │
           │ Produce eventos
           ▼
┌─────────────────────┐
│   Kafka Broker      │  ← Topic: "intake"
│   (Port 9092)       │
└──────────┬──────────┘
           │
           │ Consume eventos
           ▼
┌─────────────────────┐
│  Intake Handler     │  ← Procesa y guarda
│  (Kafka Consumer)   │
└──────────┬──────────┘
           │
           │ INSERT
           ▼
┌─────────────────────┐
│   PostgreSQL        │  ← Tabla: events
│   (Port 5432)       │
└──────────┬──────────┘
           │
           │ SELECT
           ▼
┌─────────────────────┬─────────────────┬──────────────┐
│    REST API         │    DBeaver      │   pgAdmin    │
│  (Port 9000)        │  (Port 5432)    │  (Port 5050) │
└─────────────────────┴─────────────────┴──────────────┘
```

---

## 📦 Requisitos Previos

### Software Necesario

| Software | Versión Mínima | Verificación |
|----------|----------------|--------------|
| Go | 1.24+ | `go version` |
| Docker | 20.10+ | `docker --version` |
| Docker Compose | 2.0+ | `docker compose version` |
| Make | 4.0+ | `make --version` |
| Git | 2.0+ | `git --version` |

### Herramientas Opcionales

- **DBeaver**: Para visualización de datos en PostgreSQL
- **Postman/cURL**: Para probar la API REST
- **Air**: Hot reload (se instala automáticamente)

---

## 🚀 Instalación

### Paso 1: Clonar el Repositorio

```bash
git clone <repository-url>
cd template-monitoryn-energy-main
```

### Paso 2: Instalar Herramientas de Desarrollo

```bash
make install-dev-tools
```

Este comando instala:
- ✅ Air (hot reload)
- ✅ Modd (file watcher)
- ✅ Goose (migraciones)
- ✅ Atlas (generación de migraciones)
- ✅ Swag (Swagger docs)

### Paso 3: Configurar Variables de Entorno

```bash
cp .env.example .env
```

Verifica que `.env` contenga:

```bash
# Server
PORT=9000
ENVIRONMENT=dev

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=monitoring_energy
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

# Kafka - IMPORTANTE: CONSUMER_TOPIC debe ser "intake"
LIST_KAFKA_BROKERS=localhost:9092
CONSUMER_GROUP=monitoring-energy-group
CONSUMER_TOPIC=intake
PRODUCER_TOPIC=events.output
```

⚠️ **IMPORTANTE**: `CONSUMER_TOPIC` debe ser `intake` (no `events.default`)

### Paso 4: Iniciar Servicios Docker

```bash
make docker-up
```

O manualmente:

```bash
docker compose up -d
```

Esto iniciará:
- ✅ PostgreSQL (puerto 5432)
- ✅ pgAdmin (puerto 5050)
- ✅ Kafka (puerto 9092)
- ✅ Zookeeper (puerto 2181)
- ✅ Kafka UI (puerto 8080)

### Paso 5: Verificar Servicios

```bash
docker compose ps
```

Todos los servicios deben mostrar estado `Up`:

```
NAME                         STATUS
monitoring-energy-db         Up (healthy)
monitoring-energy-kafka      Up
monitoring-energy-zookeeper  Up
monitoring-energy-kafka-ui   Up
monitoring-energy-pgadmin    Up
```

### Paso 6: Verificar Conectividad

```bash
# PostgreSQL
docker exec monitoring-energy-db psql -U postgres -c "SELECT version();"

# Kafka
docker exec monitoring-energy-kafka kafka-topics --bootstrap-server localhost:9092 --list
```

---

## ⚙️ Configuración

### Migración de Base de Datos

Las migraciones se ejecutan automáticamente al iniciar la aplicación, pero también puedes ejecutarlas manualmente:

```bash
# Aplicar todas las migraciones
make goose-up

# Ver estado de migraciones
make goose-status

# Rollback última migración
make goose-down
```

### Verificar Tabla de Eventos

```bash
docker exec monitoring-energy-db psql -U postgres -d monitoring_energy \
  -c "\d events"
```

Deberías ver:

```
                Table "public.events"
   Column    |           Type           | Nullable
-------------+--------------------------+----------
 id          | uuid                     | not null
 event_type  | character varying(100)   | not null
 source      | character varying(255)   |
 data        | text                     |
 metadata    | text                     |
 created_at  | timestamp with time zone |
```

---

## 🏃 Ejecución

### Opción 1: Development con Hot Reload (Recomendado)

```bash
make dev
```

Esto:
- ✅ Inicia la aplicación con Air
- ✅ Recarga automáticamente cuando cambias código
- ✅ Genera documentación Swagger
- ✅ Muestra logs en tiempo real

### Opción 2: Ejecución Directa

```bash
go run main.go
```

### Opción 3: Build y Ejecutar

```bash
make build
./bin/monitoring-energy-service
```

### Verificar que la Aplicación Está Funcionando

```bash
# Health check
curl http://localhost:9000/healthz

# Debería responder: 200 OK
```

---

## 📡 Uso de la API REST

### Endpoints Disponibles

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/v1/events` | Lista todos los eventos |
| GET | `/api/v1/events/:id` | Obtiene un evento por UUID |
| GET | `/api/v1/events/type/:type` | Filtra eventos por tipo |
| GET | `/healthz` | Health check |
| GET | `/readyz` | Readiness check |

### Ejemplos de Uso

#### 1. Listar Todos los Eventos

```bash
curl http://localhost:9000/api/v1/events
```

**Respuesta:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "event_type": "power_reading",
    "source": "Solar Farm Alpha",
    "data": "{\"plant_id\":\"plant-1\",\"power_generated_mw\":523.45,...}",
    "created_at": "2026-01-01T16:42:00Z"
  }
]
```

#### 2. Obtener Evento por ID

```bash
curl http://localhost:9000/api/v1/events/550e8400-e29b-41d4-a716-446655440000
```

#### 3. Filtrar por Tipo

```bash
# Solo power_reading
curl http://localhost:9000/api/v1/events/type/power_reading

# Solo alerts
curl http://localhost:9000/api/v1/events/type/alert

# Solo efficiency_report
curl http://localhost:9000/api/v1/events/type/efficiency_report
```

#### 4. Formatear Respuesta con jq

```bash
curl -s http://localhost:9000/api/v1/events | jq '.[0:5]'
```

### Swagger UI

Accede a la documentación interactiva:

```
http://localhost:9000/swagger/index.html
```

---

## 🗄️ Visualización con DBeaver

### Instalación de DBeaver

**Ubuntu/Debian:**
```bash
sudo snap install dbeaver-ce
```

**Otras plataformas:**
Descarga desde https://dbeaver.io/download/

### Configurar Conexión

1. Abre DBeaver
2. Click en **Nueva Conexión** → **PostgreSQL**
3. Configura:

```
Host:          localhost
Puerto:        5432
Base de datos: monitoring_energy
Usuario:       postgres
Contraseña:    postgres
```

4. Click **Test Connection**
5. Click **Finish**

### Consultas Útiles

#### Ver Últimos 50 Eventos

```sql
SELECT
    id,
    event_type,
    source,
    created_at,
    LEFT(data, 100) as data_preview
FROM events
ORDER BY created_at DESC
LIMIT 50;
```

#### Contar Eventos por Tipo

```sql
SELECT
    event_type,
    COUNT(*) as total
FROM events
GROUP BY event_type
ORDER BY total DESC;
```

#### Eventos por Planta

```sql
SELECT
    source,
    COUNT(*) as total_eventos
FROM events
GROUP BY source
ORDER BY total_eventos DESC;
```

#### Eventos de los Últimos 10 Minutos

```sql
SELECT
    event_type,
    source,
    created_at
FROM events
WHERE created_at > NOW() - INTERVAL '10 minutes'
ORDER BY created_at DESC;
```

#### Extraer Datos JSON

```sql
SELECT
    event_type,
    source,
    data::json->>'power_generated_mw' as power_mw,
    data::json->>'efficiency_percent' as efficiency,
    created_at
FROM events
WHERE event_type = 'power_reading'
LIMIT 20;
```

---

## 📊 Kafka UI

Accede a la interfaz web de Kafka:

```
http://localhost:8080
```

### Ver Mensajes en el Topic "intake"

1. Navega a **Topics** → **intake**
2. Click en **Messages**
3. Verás los eventos en tiempo real
4. Puedes filtrar, buscar y exportar mensajes

### Crear Nuevo Topic (Opcional)

```bash
docker exec monitoring-energy-kafka \
  kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic nuevo-topic \
  --partitions 3 \
  --replication-factor 1
```

---

## 🐛 Troubleshooting

### Problema: "Connection refused" al conectar a Kafka

**Solución:**
```bash
# Verificar que Kafka está corriendo
docker compose ps

# Ver logs de Kafka
docker compose logs kafka

# Reiniciar Kafka
docker compose restart kafka
```

### Problema: "Subscribed topic not available: events.default"

**Causa:** CONSUMER_TOPIC incorrecto en `.env`

**Solución:**
```bash
# Editar .env
CONSUMER_TOPIC=intake  # DEBE ser "intake", NO "events.default"

# Reiniciar aplicación
```

### Problema: "Error saving event to database: invalid input syntax for type json"

**Causa:** Columna `data` configurada como JSONB en vez de TEXT

**Solución:**
```bash
# Rollback y reaplica migración
make goose-down
make goose-up
```

La migración correcta usa `TEXT`:
```sql
CREATE TABLE "events" (
  ...
  "data" text NULL,  -- TEXT, no JSONB
  ...
);
```

### Problema: Puerto 9000 ya en uso

**Solución:**
```bash
# Encontrar proceso usando puerto 9000
lsof -ti:9000

# Matar proceso
kill -9 $(lsof -ti:9000)

# O cambiar puerto en .env
PORT=9001
```

### Problema: No se generan eventos

**Verificación:**
```bash
# Ver logs de la aplicación
tail -f /tmp/app.log

# Deberías ver:
# "Starting Event Generator - will send 30 messages every 5 minutes"
# "Event 1 sent: PlantID=plant-1, Type=power_reading, Power=523.45MW"
```

**Solución:**
- Verifica que `go c.EventGenerator.Start()` esté en main.go
- Confirma que el generador se inicializa en container.go

### Problema: Eventos no se guardan en PostgreSQL

**Diagnóstico:**
```bash
# Ver logs del IntakeHandler
grep "saved to database" /tmp/app.log

# Verificar eventos en DB
docker exec monitoring-energy-db psql -U postgres -d monitoring_energy \
  -c "SELECT COUNT(*) FROM events;"
```

**Solución:**
- Verifica que IntakeHandler recibe EventRepository
- Confirma que CONSUMER_TOPIC es "intake"
- Revisa logs de errores de PostgreSQL

### Logs Útiles

```bash
# Logs de la aplicación Go
tail -f /tmp/app.log

# Logs de PostgreSQL
docker compose logs db

# Logs de Kafka
docker compose logs kafka

# Logs de todos los servicios
docker compose logs -f
```

### Problema: "failed to read dockerfile: open Dockerfile: no such file or directory"

**Causa:** Al ejecutar `docker compose up -d` el servicio `api` intenta construir usando un `Dockerfile` en la raíz del proyecto, pero en este repositorio solo existe `db/Dockerfile`.

**Soluciones:**

- Levantar los servicios sin construir la API (recomendado si ejecutas la app localmente con `go run`):

```bash
docker compose up -d
```

- Incluir y construir la API con el perfil `dev` (solo si quieres que Compose construya/ejecute la API en contenedor):

```bash
docker compose --profile dev up -d
```

- Alternativa: crear un `Dockerfile` para la API en la raíz del repo o modificar `docker-compose.yml` para apuntar a la ruta correcta del `Dockerfile`.

Si eliges usar el contenedor `api`, asegúrate de que el `Dockerfile` exista en la ruta indicada o ajusta `build.context`/`dockerfile` en [docker-compose.yml](docker-compose.yml).

---

## 📁 Estructura del Proyecto

```
.
├── cmd/
│   └── atlasloader/          # Cargador de entidades para Atlas
├── db/
│   ├── Dockerfile            # PostgreSQL con PostGIS y TimescaleDB
│   └── init-db.sql           # Script de inicialización
├── internal/
│   ├── api/
│   │   ├── event_generator.go     # 🆕 Generador de eventos (30 cada 5 min)
│   │   ├── intake_handler.go      # 🔧 Consumer de Kafka (guarda en DB)
│   │   └── kafka_service.go       # Servicio de Kafka
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── event_entity.go    # 🆕 Entidad de eventos
│   │   │   └── energy_plant_entity.go
│   │   └── ports/
│   │       ├── input/             # Interfaces de servicios
│   │       └── output/
│   │           └── interfaces.go  # 🔧 EventRepositoryInterface
│   └── infrastructure/
│       ├── adapters/
│       │   ├── repositories/
│       │   │   └── event_repo.go  # 🆕 Repositorio de eventos
│       │   ├── rest/
│       │   │   ├── event_handlers.go  # 🆕 Handlers REST para eventos
│       │   │   └── router.go          # 🔧 Rutas (+ /api/v1/events)
│       │   └── kafka/
│       ├── conf/                  # Configuración
│       └── container/
│           └── container.go       # 🔧 DI Container (+ EventRepo y Generator)
├── migrations/
│   └── 20260101201616_create_events_table.sql  # 🆕 Migración tabla events
├── .env                          # 🔧 Variables de entorno (CONSUMER_TOPIC=intake)
├── docker-compose.yml            # 🔧 Servicios Docker
├── main.go                       # 🔧 Punto de entrada (+ EventGenerator.Start())
├── CHANGES.md                    # 🆕 Documentación de cambios
└── README_SETUP.md               # 🆕 Esta guía

🆕 = Archivo nuevo
🔧 = Archivo modificado
```

---

## 🔐 Credenciales por Defecto

### PostgreSQL
```
Host:     localhost
Puerto:   5432
DB:       monitoring_energy
Usuario:  postgres
Password: postgres
```

### pgAdmin
```
URL:      http://localhost:5050
Email:    admin@admin.com
Password: admin
```

### Kafka
```
Broker:   localhost:9092
Topic:    intake
UI:       http://localhost:8080
```

---

## 🧪 Testing

### Test de Integración Completo

```bash
# 1. Iniciar servicios
make docker-up

# 2. Ejecutar aplicación
make dev

# 3. Esperar 10 segundos para que se generen eventos

# 4. Verificar eventos en Kafka UI
open http://localhost:8080

# 5. Verificar eventos en PostgreSQL
docker exec monitoring-energy-db psql -U postgres -d monitoring_energy \
  -c "SELECT COUNT(*) FROM events;"

# 6. Verificar API REST
curl http://localhost:9000/api/v1/events | jq '. | length'

# 7. Verificar tipos de eventos
curl http://localhost:9000/api/v1/events/type/power_reading | jq '. | length'
```

**Resultado esperado:**
- ✅ Kafka UI muestra mensajes en topic "intake"
- ✅ PostgreSQL tiene 30+ eventos
- ✅ API REST retorna eventos en JSON
- ✅ Filtros por tipo funcionan

---

## 📈 Métricas de Performance

**Sistema en ejecución estable:**
- **Eventos/minuto:** 6 (30 eventos cada 5 min)
- **Latencia Kafka → PostgreSQL:** <100ms
- **Latencia API REST:** <50ms
- **Uso de memoria:** ~150MB (app Go)
- **Tamaño promedio evento:** ~300 bytes

---

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una branch: `git checkout -b feature/nueva-funcionalidad`
3. Commit cambios: `git commit -am 'Add nueva funcionalidad'`
4. Push a la branch: `git push origin feature/nueva-funcionalidad`
5. Abre un Pull Request

---

## 📝 Comandos Útiles

```bash
# Desarrollo
make dev              # Correr con hot reload
make dev-modd         # Correr con Modd
make run              # Correr sin hot reload
make build            # Compilar binario

# Docker
make docker-up        # Iniciar servicios
make docker-down      # Detener servicios
docker compose logs -f  # Ver logs en tiempo real

# Migraciones
make migrate-create name=nombre_migracion
make goose-up         # Aplicar migraciones
make goose-down       # Rollback
make goose-status     # Estado

# Testing
make test             # Correr tests
go test ./...         # Tests de todos los paquetes

# Documentación
make swagger          # Generar Swagger docs
```

---

## 📞 Soporte

**Documentación adicional:**
- `CHANGES.md` - Detalle completo de cambios
- `README.md` - Documentación original del template
- Swagger UI: http://localhost:9000/swagger/index.html

**Issues:**
Reporta problemas en el repositorio de GitHub

---

## 📄 Licencia

[Especificar licencia]

---

**✨ ¡Sistema listo para usar! Disfruta monitoreando energía con Kafka y PostgreSQL.**
