[![CI](https://github.com/francocaste/ingsoft3-tp01/actions/workflows/ci.yml/badge.svg)](https://github.com/francocaste/ingsoft3-tp01/actions/workflows/ci.yml)

# Gestión de Gastos Personales

Aplicación web para registrar gastos personales, organizarlos por categoría y controlar un
presupuesto mensual. Proyecto de la materia **Ingeniería de Software 3** (Universidad Católica
de Córdoba), desarrollado a lo largo de siete trabajos prácticos: control de versiones,
contenedores, planificación, integración continua, testing en el pipeline, entrega continua y
pruebas de integración y e2e como compuerta de promoción.

## Funcionalidades

- Alta, listado y baja de gastos (descripción, monto, categoría y fecha).
- Estado de cada gasto: `pendiente` o `pagado` (un gasto pagado no puede volver a pendiente).
- Presupuesto mensual por categoría, con el porcentaje usado y aviso visual al excederlo.
- Resumen de gastos por categoría.

### Reglas de negocio

| Regla | Respuesta de la API |
|---|---|
| Descripción y categoría obligatorias, monto mayor a cero | `400` |
| La categoría debe ser una de las válidas (Comida, Transporte, Ocio, Servicios, Salud, Otros) | `400` |
| La fecha no puede ser futura | `400` |
| El gasto no puede superar el presupuesto mensual de la categoría (el tope exacto se permite) | `422` |
| Un gasto `pagado` no puede volver a `pendiente` | `409` |
| El gasto solicitado no existe | `404` |

## Tecnologías

| Capa | Tecnología |
|---|---|
| Backend | Go 1.22 (`net/http`, `database/sql`, `lib/pq`) |
| Frontend | React 18 + Vite 5, servido por nginx |
| Base de datos | PostgreSQL 16 |
| Contenedores | Docker (multi-stage) y Docker Compose |
| Integración y entrega continua | GitHub Actions (con *environments* `qa` y `production`) |
| Registry de imágenes | GitHub Container Registry (`ghcr.io`) |
| Entornos desplegados | Render (servicios Docker) y Neon (PostgreSQL) |
| Tests | `go test` (backend), Vitest con cobertura v8 (frontend) y Playwright (integración y e2e) |

## Estructura del repositorio

```text
.
├── backend/
│   ├── main.go                  # arranque, rutas y handlers HTTP (capa fina)
│   ├── repo_sql.go              # implementación del repositorio sobre PostgreSQL
│   ├── internal/gastos/         # reglas de negocio, servicio y sus tests
│   ├── scripts/test.sh          # corre los tests y valida el umbral de cobertura
│   └── Dockerfile               # etapas: build, test y final
├── frontend/
│   ├── src/App.jsx              # interfaz
│   ├── src/lib/gastos.js        # lógica pura del frontend (y su test)
│   ├── default.conf.template    # plantilla de nginx: BACKEND_URL y DNS_RESOLVER por variable
│   ├── e2e/                     # Playwright: api.spec.js (integración) y gastos.spec.js (e2e)
│   ├── playwright.config.js     # tiempos pensados para el cold start de QA
│   └── Dockerfile               # etapas: build, test y final
├── .github/workflows/ci.yml     # CI, publicación en ghcr, pruebas contra QA y despliegue a QA y PROD
├── docker-compose.yml
├── .env.example
├── decisiones.md                # decisiones y justificaciones de cada TP
└── img/                         # capturas de evidencia
```

## Ejecución

Requisitos: Docker Desktop (o Docker Engine con el plugin de Compose).

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/francocaste/ingsoft3-tp01
   cd ingsoft3-tp01
   ```

2. Crear el archivo de variables de entorno a partir de la plantilla y definir una contraseña
   para la base de datos:

   ```bash
   cp .env.example .env        # en PowerShell: Copy-Item .env.example .env
   ```

3. Levantar los servicios:

   ```bash
   docker compose up --build
   ```

4. Abrir la aplicación en <http://localhost:3000>. La API queda disponible en
   <http://localhost:8080>.

Para detener los servicios: `docker compose down`. Los datos se guardan en el volumen `db_data`;
para borrarlos también, usar `docker compose down -v`.

## API

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/health` | Estado del servicio y de la base, con el commit que está corriendo |
| `GET` | `/api/gastos` | Lista los gastos |
| `POST` | `/api/gastos` | Crea un gasto |
| `DELETE` | `/api/gastos/{id}` | Elimina un gasto |
| `PATCH` | `/api/gastos/{id}/estado` | Cambia el estado (`pagado`) |
| `GET` | `/api/resumen` | Totales por categoría y porcentaje de presupuesto usado |
| `GET` | `/api/presupuestos` | Lista los presupuestos mensuales |
| `POST` | `/api/presupuestos` | Define o actualiza el presupuesto de una categoría |

## Tests y cobertura

La lógica de negocio está separada de la infraestructura (HTTP y base de datos), de modo que se
prueba con tests unitarios y con dobles en lugar de una base real. Los tests se miden con un
umbral de cobertura del **90 %**; por debajo de ese valor el pipeline falla y el merge queda
bloqueado.

**Backend** (desde `backend/`):

```bash
go test -cover ./internal/...
```

**Frontend** (desde `frontend/`):

```bash
npm install
npm test               # tests
npm run test:ci        # tests con cobertura y umbral
```

**En contenedor**, igual que en el pipeline (cada Dockerfile tiene una etapa `test`):

```bash
docker build --target test -t backend-test:local backend
docker run --rm -v "$PWD/TestResults:/out" backend-test:local
```

El alcance medido es la lógica pura (`backend/internal/...` y `frontend/src/lib/**`). Queda
fuera el código de arranque, los handlers, el SQL y el render de React; los criterios de esta
decisión están en [`decisiones.md`](decisiones.md).

**Integración y e2e (Playwright)**, desde `frontend/`. Son las dos suites que el pipeline corre
contra QA antes de promover a PROD; también corren en local contra el compose:

```bash
npm install
npx playwright install chromium
docker compose up -d --build                                  # desde la raíz del repo
npx playwright test e2e/api.spec.js                           # integración: API + base, sin navegador
npx playwright test e2e/gastos.spec.js                        # e2e: navegador contra el front
```

Sin variables apuntan a `http://localhost:8080` (API) y `http://localhost:3000` (front); en el
pipeline se definen `API_BASE_URL` y `E2E_BASE_URL` con las URLs de QA. Qué cubre cada suite y qué
queda afuera está en la sección 5 del TP7 de [`decisiones.md`](decisiones.md).

## Integración continua

El workflow [`ci.yml`](.github/workflows/ci.yml) se ejecuta en cada Pull Request y en cada push a
`main`, con dos jobs en paralelo: `build-backend` y `build-frontend`. Cada uno:

1. Construye y ejecuta la etapa `test` del Dockerfile (tests + cobertura + umbral).
2. Publica el resumen de cobertura en el *Summary* de la ejecución y sube el reporte como
   artefacto (`cobertura-backend`, `cobertura-frontend`).
3. Construye la imagen final y, **sólo en un push a `main`**, la publica en ghcr
   (`ghcr.io/francocaste/ingsoft3-tp01-backend|frontend:sha-<commit>`), como último paso del
   mismo job que corrió los tests.

Ambos checks son obligatorios para integrar a `main`, que está protegida (se trabaja con una
rama y un Pull Request por cambio, integrados con *squash merge*).

## Entrega continua (TP6 y TP7)

Tras el merge a `main`, el mismo workflow despliega en dos entornos de Render, cada uno con su
propia base en Neon. Desde el TP7 **se construye una sola vez y se promueve esa imagen**, y PROD
queda detrás de las pruebas contra QA:

```text
main ─→ build-backend + build-frontend ─→ deploy-qa (automático + smoke) ─→ integracion ─→ e2e ─→ deploy-prod (aprobación manual + smoke)
```

- **QA** se despliega solo; **PROD** queda en pausa hasta que alguien aprueba el environment
  `production`, y sólo se llega ahí con las dos suites de Playwright en verde. Es *Continuous
  Delivery*, no *Deployment*.
- Los servicios de Render ejecutan la **imagen ya publicada en ghcr**
  (`...-backend|frontend:sha-<commit>`); Render no reconstruye nada. El pipeline le indica cuál
  con el parámetro `imgURL` de cada *deploy hook*, y QA y PROD reciben **la misma imagen**. Los
  hooks viven como secrets de cada environment, nunca en el repositorio.
- No se publica `latest`: la etiqueta de cada imagen es el commit que la originó.
- El *smoke test* reintenta hasta que `/health` informa el commit desplegado y la base y el
  frontend responden.
- `integracion` prueba la API de QA con su base real (sin navegador) y `e2e` prueba el front de QA
  en un navegador. Si cualquiera falla, `deploy-prod` no arranca y no llega a pedir aprobación;
  los reportes de Playwright quedan como artefactos de la corrida (7 días).

| Variable | Dónde | Para qué |
|---|---|---|
| `DATABASE_URL`, `PORT` | backend (entorno) | conexión a la base y puerto |
| `BACKEND_URL`, `DNS_RESOLVER` | frontend (entorno) | a dónde proxea nginx `/api/` (por defecto `http://backend:8080` y `127.0.0.11`) |
| `GIT_SHA` | backend (argumento de build) | commit que informa `/health`; lo fija el CI al construir la imagen |
| `API_BASE_URL`, `E2E_BASE_URL` | Playwright (job del pipeline o local) | API y front sobre los que corren la integración y la e2e |

Las URLs vivas de QA y PROD, las imágenes publicadas y las corridas de referencia están al
principio de [`decisiones.md`](decisiones.md). El plan de rollback (con el tiempo medido) está en la
sección 8 del TP6, y la estrategia de etiquetas y cómo llegar de una release a su imagen, en las
secciones 2 y 3 del TP7.

## Documentación

Las decisiones técnicas, los problemas encontrados y las evidencias de cada trabajo práctico
están en [`decisiones.md`](decisiones.md) y en la carpeta [`img/`](img/).
