## Enlaces útiles por TP

Acá van los enlaces que no se encuentran navegando el repositorio (paquetes publicados, entornos
vivos y corridas de referencia), agrupados por trabajo práctico. Cada TP agrega su bloque.

### Enlaces de este TP (TP6)

**Paquetes públicos en el registry** (`docker pull` sin credenciales, etiquetados con el commit
del merge del PR #47, `d3e03dd1bdaabb56d4f805f3387ffeb1caa0be16`):

- Backend: [`ghcr.io/francocaste/ingsoft3-tp01-backend`](https://github.com/users/francocaste/packages/container/package/ingsoft3-tp01-backend)

  ```bash
  docker pull ghcr.io/francocaste/ingsoft3-tp01-backend:sha-d3e03dd1bdaabb56d4f805f3387ffeb1caa0be16
  ```

- Frontend: [`ghcr.io/francocaste/ingsoft3-tp01-frontend`](https://github.com/users/francocaste/packages/container/package/ingsoft3-tp01-frontend)

  ```bash
  docker pull ghcr.io/francocaste/ingsoft3-tp01-frontend:sha-d3e03dd1bdaabb56d4f805f3387ffeb1caa0be16
  ```

**Los dos enlaces de la cadena de publicación (sección 1 del TP6):**

- Job de un Pull Request con el paso «Entrar al registry» **salteado** (verifica pero no publica):
  [job `build-backend` del PR #47 (corrida #70)](https://github.com/francocaste/ingsoft3-tp01/actions/runs/39274573037/job/118135410334)
- Lista de pasos de una corrida de `main`, donde «Entrar al registry» corre y la imagen se publica
  como último paso del mismo job que testeó:
  [job `build-backend` de la corrida #71, push a `main`](https://github.com/francocaste/ingsoft3-tp01/actions/runs/39274801397/job/118136208159)

**Entornos vivos hasta la defensa** (Render free: si estuvieron inactivos, la primera carga puede
tardar 50 s o más en responder, ver la sección 5 del TP6):

| | Frontend | API |
|---|---|---|
| QA | <https://gastos-front-qa-ez08.onrender.com> | <https://gastos-api-qa-r6ig.onrender.com/health> |
| PROD | <https://gastos-front-prod-96km.onrender.com> | <https://gastos-api-prod-fvsa.onrender.com/health> |

**Corridas citadas en el TP6** (el resto se encuentra navegando *Actions* y *Deployments*):

- Rechazo del gate, con su motivo (sección 4): [corrida #67](https://github.com/francocaste/ingsoft3-tp01/actions/runs/38921626680)
- Aprobación y llegada a PROD: [corrida #69](https://github.com/francocaste/ingsoft3-tp01/actions/runs/38933270155)
- Smoke que espera el commit, QA y PROD (sección 7): [corrida #71](https://github.com/francocaste/ingsoft3-tp01/actions/runs/39274801397), PR [#47](https://github.com/francocaste/ingsoft3-tp01/pull/47)
- Release: [`v6.0.0`](https://github.com/francocaste/ingsoft3-tp01/releases/tag/v6.0.0)

---

# Decisiones — TP1

## 1. Por qué Git no pudo resolver el conflicto solo

El conflicto se fabricó a propósito siguiendo la guía: dos ramas (`feature/titulo-a` y
`feature/titulo-b`) nacieron ambas desde el mismo commit de `main`, y cada una modificó
**la misma línea** del `README.md` (la primera línea, el título del proyecto) con contenido
distinto ("versión A" vs "versión B").

Git resuelve merges automáticamente cuando los cambios tocan **partes distintas** del archivo,
comparando cada rama contra el ancestro común. Pero cuando dos ramas cambian la misma línea, Git
no tiene forma de saber cuál de las dos versiones es "la correcta" — ambas son ediciones válidas
sobre el mismo punto del archivo, y elegir entre ellas es una decisión de contenido, no algo que
se pueda inferir del historial. Por eso Git se detiene y delega la decisión a una persona,
marcando el archivo con `<<<<<<<`, `=======` y `>>>>>>>`.

Para que este conflicto nunca hubiera aparecido, hubiera bastado con que las dos ramas no
tocaran la misma línea (por ejemplo, si una editaba el título y la otra agregaba una sección
distinta), o con integrar la rama B **después** de traer los cambios de A a `main` antes de
empezar a editar desde ahí — es decir, con ramas más cortas e integración más frecuente, que es
justamente lo que reduce la probabilidad y el tamaño de los conflictos en un equipo real.

## 2. Problemas encontrados y cómo los solucioné

No tuve problemas con la configuración de "Require approvals" (quedó en cero aprobaciones
como pide la guía, para poder mergear mis propios PRs al ser un trabajo individual). El resto
de los pasos —crear el repo, proteger `main`, hacer los PRs, provocar y resolver el conflicto,
y publicar el tag y la release— los hice todos de forma manual desde la web y la consola,
siguiendo la guía paso a paso, sin inconvenientes que resolver.

## 3. Declaración de uso de IA

No usé IA para realizar ninguno de los pasos técnicos del TP (creación del repositorio,
configuración de las protecciones de rama, Pull Requests, resolución del conflicto, tag o
release): todo el trabajo en el repositorio fue manual.

Sí usé IA (Claude) puntualmente al final, para redactar este archivo (`decisiones.md`) y
`evidencias.md` a partir de las capturas ya tomadas y del propio enunciado del TP. La
verificación consistió en revisar que lo escrito describiera con precisión lo que efectivamente
hice y lo que muestran mis capturas (por ejemplo, que la explicación del conflicto coincida con
las ramas y el archivo que realmente usé), y en poder explicar cada punto sin depender del texto
generado, ya que la defensa oral se responde con mis propias palabras.

## TP2 — Contenedores

### 1. Qué app elegí y por qué

Elegí armar una aplicación de **gestión de gastos personales** desde cero, con backend en **Go**
(librería estándar `net/http` + `database/sql`, sin framework) y frontend en **React + Vite**,
sobre una base de datos **PostgreSQL**.

Contra los criterios de la guía:
- **Buildea y corre localmente sin magia**: la probé end-to-end antes de dar por terminado el TP
  (carga, listado, borrado y resumen de gastos funcionando).
- **Tamaño acotado**: es un CRUD con una sola entidad (gastos) y una pantalla — alcanza para
  demostrar los tres requisitos (backend + frontend + base de datos) sin fricción de más.
- **La entiendo por completo**: la escribí yo mismo (con asistencia de IA, ver declaración más
  abajo), así que puedo modificarla y explicar cada parte en la defensa oral.
- Elegí Go en vez de un stack más "tradicional" de la materia porque me interesaba probar un
  lenguaje que compila a un binario único — eso simplifica notablemente el Dockerfile del backend
  (no hace falta instalar ningún runtime en la imagen final, alcanza con copiar el binario sobre
  una base mínima como `alpine`).

### 2. Decisiones de contenerización

- **Imágenes base**: `golang:1.22-alpine` para compilar el backend, `alpine:3.20` para ejecutarlo
  en producción; `node:22-alpine` para compilar el frontend, `nginx:alpine` para servirlo.
- **Estructura multi-stage**: en las dos imágenes, la etapa de build (con el compilador/SDK) se
  descarta por completo en la imagen final — solo viaja el binario de Go compilado, o los
  estáticos que generó `vite build`.
- **Qué persiste y qué no**: los datos de PostgreSQL viven en el volumen nombrado `db_data`,
  montado en `/var/lib/postgresql/data`. Los contenedores de `backend` y `frontend` son
  completamente efímeros — no guardan ningún estado propio, así que se pueden recrear sin
  pérdida de información.
- **Comunicación entre servicios**: el backend se conecta a la base por el nombre de servicio
  `db` (`DATABASE_URL=postgres://postgres:${DB_PASSWORD}@db:5432/app`), y el frontend habla con
  el backend por ruta relativa (`/api/...`), que nginx reenvía internamente a `http://backend:8080`
  — así el navegador nunca necesita conocer el nombre `backend` (que no resolvería fuera de la
  red de compose).
- **Secretos**: la contraseña de la base viaja por la variable `DB_PASSWORD`, definida en un
  `.env` que no se commitea (con `.env.example` como plantilla commiteada).
- **Registry**: publiqué las dos imágenes en GitHub Container Registry (ghcr.io), con tag
  `v0.1.0`, y las hice públicas para que `docker-compose.registry.yml` pueda bajarlas sin
  credenciales.

### 3. Problemas encontrados y cómo los resolví

- **`go mod download` no generaba un `go.sum` completo**: el build fallaba con
  `missing go.sum entry for module providing package github.com/lib/pq`. Lo resolví copiando
  todo el código fuente antes de correr `go mod tidy` (en vez de `go mod download` con solo el
  `go.mod`), porque `tidy` sí escanea los imports reales y arma el `go.sum` correcto.
- **`npm ci` fallaba** porque no tenía un `package-lock.json` commiteado (lo esperable con
  `npm ci`, que exige un lockfile ya generado). Lo resolví usando `npm install` en el Dockerfile,
  que genera el lockfile sobre la marcha si no existe.
- **El primer `docker push` del backend falló** con `error from registry: unknown` después de
  subir todas las capas — fue un error transitorio del lado de GitHub (probablemente por ser la
  primera vez que se creaba ese paquete en mi cuenta). Reintentar el mismo comando lo resolvió sin
  cambiar nada.

### 4. Declaración de uso de IA

Usé IA (Claude) de forma extensa en este TP, a diferencia del TP1: le pedí que me ayudara a
elegir el stack (Go + React) y que armara la aplicación completa de gestión de gastos —
backend, frontend, Dockerfiles multi-stage, `docker-compose.yml`, `docker-compose.registry.yml`
y el `nginx.conf`. Después corrí yo mismo cada comando de build, arranque y prueba, y cuando
aparecieron errores (los tres del punto anterior) se los reporté a la IA, que me indicó la causa
y el cambio puntual a hacer en el Dockerfile correspondiente.

Verifiqué lo generado ejecutándolo yo mismo en mi máquina en cada paso — no di nada por
funcionando hasta ver la salida real de mis propios comandos (`docker compose ps`, `curl`,
la interfaz en el navegador, los `docker images` de la comparación de tamaños, y el `pull`
deslogueado para confirmar que las imágenes eran realmente públicas). Entiendo la estructura de
cada Dockerfile y del compose lo suficiente como para explicarla en la defensa oral: qué hace
cada instrucción, por qué el orden importa para el cache, por qué el volumen es necesario, y por
qué el frontend habla con el backend por ruta relativa y no por URL absoluta.

## TP3 — Planificación y trazabilidad

### 1. Duración del sprint

Elegí sprints de **1 semana**, porque coincide con el ritmo de entrega de la cátedra (un TP
por semana) y permite cerrar un ciclo completo de planificación → ejecución → revisión en cada
entrega, alineando el objetivo del sprint con el objetivo de cada práctico. Un sprint más largo
(por ejemplo, un mes) no me daría feedback a tiempo respecto del calendario real de la materia.

### 2. Límite de trabajo en progreso

Configuré el límite de la columna *In Progress* en **2**, siguiendo la regla de arranque de la
guía: cantidad de personas trabajando (1, en mi caso, porque el TP es individual) más uno. El
"+1" es la válvula para cuando algo queda esperando (por ejemplo, una revisión propia o un
bloqueo externo) y necesito poder avanzar en otra cosa sin romper el límite. Si en la práctica
nunca llego a usar las 2 posiciones, es señal de que podría bajarlo a 1; si lo alcanzo todo el
tiempo y me frena, sería señal de subirlo.

### 3. Diagnóstico de la historia mal escrita

La historia de ejemplo *"Como desarrollador quiero crear la tabla usuarios"* está mal escrita
porque es una **tarea técnica disfrazada de historia**: no describe una capacidad observable
por un usuario, sino un paso de implementación interno. Le falta además el "para qué" (el
beneficio) — "crear una tabla" no es algo que nadie "quiera" en sí mismo, es un medio para
lograr otra cosa. Se puede reconocer por el anti-patrón que menciona la guía: si el "quiero"
describe una pieza técnica en vez de una capacidad de negocio, es una tarea.

Cómo la reescribiría: la elevaría a una historia real, por ejemplo *"Como usuario quiero
registrarme e iniciar sesión para poder acceder a mis datos personales"*, y "crear la tabla
usuarios" pasaría a ser una de las **tareas técnicas** dentro de esa historia (junto con, por
ejemplo, "implementar el endpoint de registro" o "validar el formato del email").

### 4. Problemas encontrados y cómo los resolví

- **Los comandos `gh` con `\` para continuar en varias líneas fallaban en PowerShell**
  (`ParserError: Missing expression after unary operator '--'`). PowerShell no interpreta la
  barra invertida como continuación de línea (eso es sintaxis de bash/zsh). Lo resolví escribiendo
  cada comando `gh issue create` en una sola línea, con comillas dobles en vez de simples.
- **El primer `gh project item-add` falló** porque copié el número de proyecto entre signos
  `< >` literales (`<2>`), que PowerShell interpretó como redirección de archivo. Lo resolví
  usando el número solo, sin los signos.
- **El PR que debía cerrar la tarea "Escribir el workflow de build y tests" (#9) se mergeó con
  `Closes #12`** — un número de issue que no existe en mi repositorio (probablemente copiado de
  un ejemplo sin ajustar). Como el merge ya había ocurrido, el auto-cierre no se disparó
  retroactivamente. Lo resolví abriendo un segundo PR chico (un comentario adicional en el
  `ci.yml`) con la descripción corregida a `Closes #9`, que sí cerró la tarea correcta al
  mergearse.

### 5. Declaración de uso de IA

Usé IA (Claude) para traducir la guía del TP a una secuencia concreta de comandos `gh` y pasos
de la interfaz web adaptados a mi repositorio, y para diagnosticar los tres problemas del punto
anterior a medida que aparecían (en particular, para identificar por qué `Closes #12` no había
cerrado ningún issue). Verifiqué cada paso ejecutándolo yo mismo y revisando el resultado real
en GitHub: el estado de cada issue con `gh issue list`, la jerarquía de sub-issues abriendo cada
uno en la web, el tablero y el límite de trabajo en progreso, y finalmente que la tarea #9
apareciera cerrada y enlazada al PR correcto después de corregir el error. Entiendo la
diferencia entre historia, tarea y épica, por qué la trazabilidad depende de que `Closes #N`
esté en la descripción del PR y no en un comentario posterior, y por qué elegí la duración de
sprint y el límite de trabajo en progreso que elegí — puedo explicar cada uno de estos puntos
en la defensa oral sin depender del texto generado.

## TP4 — CI: Pipelines as Code

### 1. Estructura del pipeline

Elegí **2 jobs en paralelo** (`build-backend` y `build-frontend`), uno por cada Dockerfile que
ya tenía del TP2. Los separé porque son unidades de build independientes entre sí — no hay
ninguna dependencia de orden entre construir la imagen del backend y la del frontend, así que
correrlos en paralelo (cada uno en su propio runner) aprovecha mejor el tiempo que hacerlo en
serie. Si mi app tuviera un solo Dockerfile, un solo job hubiera sido igual de válido.

### 2. Qué cachea mi pipeline

Lo que se cachea son las **capas de la imagen Docker** de cada job, guardadas en el cache de
GitHub Actions (`type=gha`) con un `scope` distinto por job (`backend` y `frontend`) para que no
se pisen entre sí. En mi caso, las capas que más se reutilizan son las de las etapas base
(`FROM golang:1.22-alpine`, `FROM alpine:3.20`) y, cuando no cambian las dependencias, la etapa
de `go mod tidy` / `npm install`. La etapa que casi siempre se rehace es la que compila el
código (`go build`, `npm run build`), porque cambia en cada commit.

Si el cache desaparece (la plataforma lo puede desalojar en cualquier momento, o tiene límite de
tamaño), el pipeline **tiene que funcionar igual**, solo que construyendo todo desde cero — más
lento, pero no roto. Si mi build fallara sin cache, no tendría un cache: tendría una dependencia
escondida que dependía de un estado previo, y eso sería un bug del Dockerfile, no del pipeline.

### 3. Por qué el pipeline construye con mi Dockerfile en vez de compilar por su cuenta

Si el workflow tuviera pasos propios de `go build` o `npm run build` en vez de invocar mi
Dockerfile, tendría **dos definiciones de build**: una para verificar en CI y otra para lo que
después efectivamente se despliega (la imagen). Esas dos definiciones divergen tarde o temprano
— por ejemplo, si cambio una variable de entorno o una versión de imagen base en el Dockerfile y
me olvido de replicarlo en el workflow. Construyendo siempre con el mismo Dockerfile, lo que el
pipeline verifica es exactamente lo mismo que se va a ejecutar en producción, no una aproximación.

### 4. Problemas encontrados y cómo los resolví

- **El primer intento de romper el build a propósito no funcionaba**: agregué un import
  inexistente en Go, pero al guardar el archivo en VS Code, la extensión de Go (con
  `goimports`) lo eliminaba automáticamente por no estar en uso, así que el build seguía
  pasando sin que me diera cuenta de por qué. Lo resolví agregando también una línea que
  **usara** ese import (una llamada a una función del paquete inexistente) — al estar "en uso",
  el formateador ya no lo borraba, y ahí sí el build falló como se buscaba.
- **Al confirmar la configuración del gate con `gh api` y un filtro `--jq` con comillas
  anidadas, PowerShell interpretaba mal el comando** (`accepts 1 arg(s), received 8`). Lo
  resolví pidiendo el JSON completo sin filtro y leyendo a simple vista los campos que
  necesitaba verificar (`strict`, `contexts`, `enforce_admins`).
- **Al armar el commit del fix, lo hice parado en la rama equivocada** (`docs/muestra-del-freno`
  en vez de `feature/demo-gate`), así que el cambio no se aplicó donde correspondía. Lo detecté
  porque el mensaje de Git mostraba "On branch docs/muestra-del-freno", y lo resolví cambiando a
  la rama correcta antes de repetir el commit.

### 5. Declaración de uso de IA

Usé IA (Claude) para traducir la guía en un workflow YAML concreto para mi stack (Go + React +
Docker), y para diagnosticar los tres problemas del punto anterior a medida que aparecían — en
particular, para entender por qué el formateador de Go estaba neutralizando mi intento de romper
el build a propósito. Verifiqué cada paso ejecutándolo yo mismo: confirmé el JSON real de la
protección de rama con `gh api`, vi con mis propios ojos el PR pasar de rojo a bloqueado a verde
a mergeado, y probé que el build efectivamente fallaba en mi máquina (`docker build --no-cache`)
antes de subirlo al pipeline. Entiendo la diferencia entre `push` y `pull_request` como
triggers, por qué el cache necesita `scope` separado por job, y por qué el pipeline usa el
Dockerfile en vez de compilar por su cuenta — puedo explicar cada uno de estos puntos en la
defensa oral sin depender del texto generado.


## TP5 — Testing en el pipeline

### 1. Qué tests elegí y por qué

Probé las **5 reglas de negocio** de la app más el cálculo del porcentaje de presupuesto, porque
son lo que más confianza cuesta si se rompe: monto positivo y campos obligatorios, categoría
válida, fecha no futura, tope del presupuesto mensual y transición de estado (`pagado` no vuelve
a `pendiente`). Tests de backend: 17 funciones de test en `backend/internal/gastos` (más sus
subtests). Tests de frontend: 23 casos en `frontend/src/lib/gastos.test.js`.

- **Parametrizados** (tabla en Go, `it.each` en vitest): una misma lógica con muchos datos. Los
  usé sobre todo para los **bordes**: gastar *exactamente* el tope del presupuesto se permite y
  pasarse por un centavo no; un gasto de "hoy a última hora" es válido y uno de mañana no.
- **Caso de error**: que el repositorio falle (base caída) y el servicio propague el error en
  vez de tragárselo; en el frontend, que `fetch` rechace y el error se propague.
- **Con doble (mock)**: ver sección 6.

No probé `App.jsx` (render) ni los handlers HTTP ni el SQL: son tests de otra capa de la
pirámide (componentes / integración) y los dejo fuera a propósito (ver sección 3).

### 2. Umbral de cobertura: número, métrica y por qué

| | Métrica | Umbral | Valor actual |
|---|---|---|---|
| Backend (Go) | sentencias (statements) | **90 %** | 100 % |
| Frontend (vitest + v8) | líneas **y** ramas | **90 % / 90 %** | 100 % / 100 % |

Go sólo mide cobertura de **sentencias** (no tiene cobertura de ramas), por eso ahí uso esa
métrica; en el frontend uso líneas y ramas porque vitest reporta las dos.

Elegí **90 %** y no otro número por tres razones: (1) el alcance medido es sólo la lógica pura,
donde 100 % es alcanzable, así que 90 % deja margen para una rama defensiva sin test sin frenar
un PR por una sola línea; (2) es lo bastante alto como para que código nuevo sin tests **sí**
frene: lo comprobé (sección 5); (3) con 60 % u 80 % una función nueva de unas diez sentencias
pasaría sin que nadie la pruebe. No puse 100 % porque empuja a escribir tests sólo para tocar
líneas (ley de Goodhart: cuando la métrica se vuelve la meta, deja de medir lo que importaba) y
a frenar por código inalcanzable.

El umbral está en el repo, no en una configuración externa: `COVERAGE_MIN` en
`backend/scripts/test.sh` y `coverage.thresholds` en `frontend/vite.config.js`. Cambiarlo
requiere un PR, que pasa por revisión y por el mismo gate.

### 3. Qué excluí de la cobertura y por qué

Mido **sólo** `backend/internal/...` y `frontend/src/lib/**`. Quedan afuera:

- `backend/main.go` (arranque, rutas, CORS, `migrate`, `waitForDB`) y los handlers HTTP: es
  código de arranque y "pegamento" que necesita una base real o un servidor levantado.
- `backend/repo_sql.go`: es SQL; probarlo de verdad requiere Postgres (test de integración).
- `frontend/src/main.jsx` y `frontend/src/App.jsx`: arranque de React y render; se prueban con
  tests de componentes o E2E, no unitarios.

Para que la exclusión no sea trampa, **moví la lógica fuera de esos archivos** (a
`internal/gastos` y `src/lib`) y dejé en ellos sólo código de pegamento. Excluir un archivo con
reglas de negocio adentro sería esconder lo que sí importa.

### 4. Cómo corre en el pipeline

- Los tests corren **dentro de los jobs que ya existían** (`build-backend` y `build-frontend`),
  sin jobs nuevos: así los nombres de los checks obligatorios de `main` no cambian y la
  protección de rama sigue funcionando.
- Cada Dockerfile tiene una etapa `test` (`FROM build AS test` + `ENTRYPOINT`). El pipeline la
  construye con `target: test` y `load: true` (para poder hacer `docker run`), y la ejecuta. Si
  falla un test o la cobertura queda por debajo del umbral, el contenedor sale con código 1 y el
  job queda en rojo. La etapa final (nginx / alpine) no la incluye.
- Los reportes se publican con `if: ${{ !cancelled() }}`: aunque los tests fallen, el resumen
  (`$GITHUB_STEP_SUMMARY`) y el artefacto (`cobertura-backend`, `cobertura-frontend`) se suben,
  que es cuando más sirven. Si los tests fallan, la imagen final ni se construye.
- Para el frontend, el pipeline lee `coverage-summary.json` con `jq` y **falla si
  `total.lines.total` es 0**: sin esa guarda, un `include` mal escrito mediría 0 de 0 líneas y
  el umbral pasaría sin haber medido nada.
- Backend: la cobertura la calcula `scripts/test.sh` (`go test -coverprofile`), que además genera
  `coverage.html` y `cobertura-por-funcion.txt`.

### 5. Cobertura no es calidad, y los PRs de demostración

**Ejemplo de que 100 % no alcanza (PR de demo #2).** La línea
`if totalDelMes+monto > presupuesto` queda 100 % cubierta con sólo dos casos ("dentro del
presupuesto" y "se pasa por mucho"). Cambié `>` por `>=`: la cobertura **siguió en 100 %**, pero
el caso borde "justo en el tope (se permite)" falló y el pipeline quedó en rojo con el merge
bloqueado. Un test sin ningún `assert` que sólo llame a `CrearGasto` daría la misma cobertura sin
verificar nada. La cobertura dice qué código se **ejecutó**, no que esté **verificado**.

**PR bloqueado por cobertura (PR de demo #1).** Agregué `limites.go` (dos funciones nuevas, sin
tests). Compila y todos los tests existentes pasan, pero la cobertura del backend cayó a
**82,4 %** (< 90 %) y el check `build-backend` quedó en rojo, con *Squash and merge* bloqueado.
Después agregué `limites_test.go`, la cobertura volvió a 100 % y el mismo PR pasó a verde y se
mergeó.

En el frontend probé lo mismo en local con una función nueva sin test: **líneas 84,28 %**
(< 90 %), y vitest 3.2.7 (con `@vitest/coverage-v8` 3.2.7) salió con
`ERROR: Coverage for lines (84.28%) does not meet global threshold (90%)` y código de salida 1.
Dato que me llamó la atención: **las ramas seguían en 100 %**, porque v8 no cuenta las ramas de
funciones que nunca se ejecutaron. Por eso el umbral exige líneas **y** ramas y no confío sólo
en ramas.

- PR de demo #1 (rojo → verde, mergeado): `https://github.com/francocaste/ingsoft3-tp01/pull/38`
- Corrida en rojo (82,4 %): `https://github.com/francocaste/ingsoft3-tp01/actions/runs/38836552181`
- Corrida en verde tras agregar los tests: `https://github.com/francocaste/ingsoft3-tp01/actions/runs/38837448964`
- PR de demo #2 (`demo/rompe-regla`, **abierto y en rojo hasta la defensa**):
  `https://github.com/francocaste/ingsoft3-tp01/pull/39` — run: `https://github.com/francocaste/ingsoft3-tp01/actions/runs/38838021387/job/116771406938`
- Capturas (carpeta `img/`):
  - `img/13-demo-rojo-cobertura-82.png`: demo 1, log en rojo con la cobertura en 82,4 % (< 90 %).
  - `img/14-demo1-pr-bloqueado-rojo.png`: demo 1, PR con `build-backend` en rojo y el merge bloqueado.
  - `img/17-demo1-verde.png`: demo 1, el mismo PR en verde (run rojo #52, run verde #53) después de agregar los tests.
  - `img/15-demo2-pr-bloqueado-rojo.png`: demo 2, PR con el merge bloqueado.
  - `img/16-demo2-log-ci-falla-cobertura-100.png`: demo 2, el log del pipeline fallando con la cobertura en 100 %.

### 6. Qué tuve que refactorizar para poder usar un mock

El código original **no se podía testear sin una base de datos**: había una variable global
`db *sql.DB` y las reglas estaban escritas adentro de los handlers HTTP, mezcladas con SQL.
Para poder probar las reglas hice tres cosas, sin cambiar el comportamiento ni los códigos de
estado:

1. Saqué las reglas a **funciones puras** (`ValidarCampos`, `ValidarCategoria`, `ValidarFecha`,
   `ValidarPresupuesto`, `ValidarTransicion`, `PorcentajeUsado`) en `internal/gastos`.
2. Definí una **interfaz `Repositorio`** (lo que el servicio necesita de la base) y un
   `Servicio` que la recibe por **inyección de dependencias**; en producción la implementa
   Postgres (`repo_sql.go`) y en los tests un doble (`repoFalso`). También inyecté el "ahora"
   (`func() time.Time`) para que los tests de fecha no dependan del reloj.
3. Los handlers quedaron finos: decodifican, llaman al servicio y traducen el error de dominio a
   un status HTTP (`statusDe`). En el frontend pasé `fetch` por parámetro a `enviarGasto`.

Mi doble es **stub y spy a la vez**: *stub* porque devuelve datos que yo fijo (presupuesto,
total del mes), *spy* porque cuenta las llamadas, y eso me permite verificar **comportamiento**:
"si la categoría es inválida, no se toca la base" (0 llamadas) o "se inserta exactamente una
vez". En el frontend, `vi.fn()` hace lo mismo con `fetch` ("si el formulario es inválido, no se
llama al backend").

### 7. Ejercicio: un camino sin cubrir

Con 100 % en el alcance medido, el camino sin cubrir está en lo que excluí: la función
`statusDe` de `main.go`, que traduce cada error de dominio a su código HTTP (400 / 404 / 409 /
422 / 500).

1. **Qué comportamiento queda sin verificar:** que cada error de dominio salga con *su* código
   (por ejemplo `ErrTransicionInvalida` → 409, `ErrPresupuestoExcedido` → 422) y que un error que
   no es de dominio salga como 500.
2. **Qué bug se podría esconder:** alguien cambia el 409 por 400 (o se olvida de agregar un
   error nuevo al `switch`, que caería en 500). Todos los tests seguirían en verde y la
   cobertura en 100 %, pero el frontend recibiría un código distinto.
3. **Qué test lo cubriría:** un test de tabla sobre `statusDe` con cada error (incluso envuelto
   con `fmt.Errorf("...: %w", err)`) y su código esperado, o un test de integración con
   `httptest` que dispare cada caso por HTTP. No lo escribí en este TP; queda como deuda
   conocida y por eso mi 100 % es sobre la lógica y no sobre "toda la app".

### 8. Problemas encontrados y cómo los resolví

- **El código no era testeable** (global `db` y reglas dentro de los handlers): lo resolví con
  el refactor de la sección 6.
- **v8 no cuenta las ramas de funciones nunca ejecutadas**: se vio al probar el bloqueo del
  frontend (ramas 100 % con líneas 84 %); por eso el umbral exige líneas y ramas.
- **Una cobertura que "pasa" midiendo nada**: si `coverage.include` apuntara a una ruta
  inexistente, vitest mediría 0 líneas y no se quejaría. Agregué la guarda con `jq` en el
  pipeline (sección 4).
- **`scripts/test.sh: set: line 10: illegal option -` dentro del contenedor**: en Windows el
  script se bajó con saltos de línea CRLF y Alpine leyó `pipefail\r`. Lo resolví con un
  `.gitattributes` (`*.sh text eol=lf`, `Dockerfile text eol=lf`); el repo ya tenía LF, así que
  el pipeline nunca estuvo afectado, sólo mi copia local.
- **`${COVERAGE_DIR:-coverage}` en un script de npm no funciona en Windows** (usa `cmd`, no
  `sh`): moví esa lógica a `process.env.COVERAGE_DIR` en `vite.config.js`, que anda igual en
  ambos lados.
- **Docker Desktop: `failed to prepare extraction snapshot ... parent snapshot ... does not exist`**
  al exportar la imagen: era un layer de la caché local roto, no el código. Se resolvió
  reconstruyendo con `docker build --no-cache`.
- **`pq: column "estado" of relation "gastos" does not exist`** al probar la app: el volumen de
  Postgres era del TP2 y el `CREATE TABLE IF NOT EXISTS` no altera una tabla que ya existe. Lo
  resolví con `docker compose down -v`. Deuda: faltan migraciones reales.
- **Deuda conocida de zona horaria:** `hoyISO()` usa `toISOString()` (UTC), así que después de
  las 21:00 en Argentina devuelve el día siguiente. Lo vi con el formulario ofreciendo el 16/10
  cuando para mí era 15/10. Viene del código original y el backend usa el mismo criterio, así
  que no lo cambié en este TP.
- **`npm audit` reporta 6 vulnerabilidades** en dependencias de desarrollo (apareció al instalar
  vitest; no revisé una por una de qué paquete vienen). No entran en la imagen final (que sólo
  lleva nginx con el build) y no usé `npm audit fix --force` para no romper versiones.

### 9. Declaración de uso de IA

Usé IA (Claude) para analizar mi código y proponer el refactor que lo hace testeable, para
escribir los tests y para armar la etapa `test` de los Dockerfiles y los pasos del workflow. Fui
armando cada etapa paso a paso, en un PR por etapa, y verifiqué cada cosa corriéndola: los tests
y el cálculo de cobertura, el bloqueo por umbral (con `COVERAGE_MIN=101` y agregando código sin
tests, en backend y en frontend) y el borde `>` vs `>=`. Puedo explicar sin el texto generado
qué es un mock frente a un stub, por qué línea y rama no son lo mismo, por qué elegí 90 % y qué
quedó afuera de la cobertura. 
Los tests, el refactor y los pasos del workflow me los propuso Claude; yo los fui escribiendo y
corriendo etapa por etapa, y entiendo qué hace cada uno. Detecté y corregí por mi cuenta errores
que aparecieron al aplicarlos: una línea sobrante en `limites.go` que rompía la compilación, el
nombre mal escrito de `.gitattributes` (`.gitatributes`, que git ignoraba), y el CRLF de
`test.sh` que sólo se manifestó en mi máquina con Docker. Claude me advirtió antes de que lo
usara que `${COVERAGE_DIR:-coverage}` no funcionaría en `cmd` de Windows, y lo moví a la config.

---

## TP6 — Entrega continua: entornos, aprobaciones y patrones de despliegue

### 1. Por qué el artefacto se publica sólo con la verificación en verde

Una imagen en el registry tiene que querer decir una cosa: «este commit pasó los tests y el
umbral de cobertura». No hay un control que lo garantice solo; son **tres** encadenados:

1. **Nada entra a `main` sin verde**: la rama está protegida y exige los checks `build-backend`
   y `build-frontend` (TP1 y TP4).
2. **Sólo `main` publica**: el `push` de la imagen está condicionado a
   `github.event_name == 'push' && github.ref == 'refs/heads/main'`, y el paso «Entrar al
   registry» sólo corre en `push`. En un Pull Request ese paso aparece salteado.
3. **Publicar es el último paso del mismo job que corre los tests**: en el workflow, «Construir y
   publicar» va después de «Correr tests y verificar cobertura». Si un paso anterior falla, el job
   se corta y no se llega a publicar. No hace falta un `if: tests passed`: el orden **es** la
   condición.

Si se publicara igual, el registry pasaría a ser un depósito de «lo que se construyó», sin
información sobre si funciona, y cualquier `docker pull` sería una apuesta. Lo bueno de estos tres
eslabones es que ninguno requiere disciplina: aunque yo me olvide, el pipeline no deja otra.

**Lo que esta cadena no garantiza.** Un `docker push` hecho a mano desde mi máquina, con un token
con permiso de escritura, entra al registry igual y con cualquier etiqueta: la cadena protege el
camino del pipeline, no el registry. Tampoco la etiqueta `sha-<commit>` es inmutable: se puede
volver a publicar. Lo que sí identificaría una imagen sin ambigüedad es su **digest**
(`@sha256:...`), que depende del contenido; desplegar por digest, y no por etiqueta, es lo que haría
falta para que «lo que corre» sea exactamente «lo que se publicó» (esto vuelve en la sección 6).

El job de publicación tiene permisos mínimos (`contents: read` y `packages: write`) y usa el
`GITHUB_TOKEN` de cada corrida, así que no agregué ningún secret nuevo para publicar. Los nombres
de imagen van en minúscula (`ghcr.io/francocaste/ingsoft3-tp01-backend|frontend`) porque ghcr no
acepta mayúsculas.

### 2. Continuous Delivery o Continuous Deployment

- **CI:** cada cambio se integra a `main` y se verifica automáticamente (TP4 y TP5).
- **Continuous Delivery:** cada commit que pasa la verificación queda **listo para producción**
  y el paso final lo decide una persona.
- **Continuous Deployment:** además, ese paso final es automático: todo commit verde llega a
  producción sin intervención.

Implementé **Continuous Delivery**: QA se despliega solo después del merge, y PROD queda detrás de
una aprobación manual (environment `production` con *required reviewer*). Corresponde a mi
contexto porque la app guarda datos de personas y no tengo hoy nada que me avise de un problema
después de desplegar: no hay métricas de errores, ni alertas, ni logs centralizados, sólo un
smoke test que mira tres cosas. Sin esa señal, automatizar la última puerta sería decidir a ciegas.

Para llegar a Continuous Deployment me faltaría: observabilidad (tasa de errores y latencia por
versión), una forma de desplegar de a poco (canary o blue-green) y **rollback automático** cuando
esas métricas empeoran, además de tests de extremo a extremo que reemplacen lo que hoy mira el
aprobador. Con eso sí lo querría para esta app, porque la aprobación humana sin información no
agrega seguridad, agrega demora.

### 3. Diseño de la cadena y alcance de los secrets

```text
PR  → build-backend + build-frontend (verifican, NO publican, NO despliegan)
main → build-backend + build-frontend (verifican y publican en ghcr)
        └─ needs: ambos ─→ deploy-qa  (if: main · environment: qa)
                            └─ needs: deploy-qa ─→ deploy-prod (environment: production · espera aprobación)
```

- **`needs`**: `deploy-qa` espera a los dos builds, así que sin CI verde no hay deploy.
  `deploy-prod` espera a `deploy-qa`: a PROD sólo se llega después de que QA haya respondido bien.
- **`if: github.ref == 'refs/heads/main'`** va sólo en `deploy-qa`; `deploy-prod` lo hereda por
  `needs`, porque si QA se saltea, PROD también. En un PR, `deploy-qa` figura como salteado: es lo
  esperado (me confundió la primera vez que lo vi).
- **Environments** (`qa` y `production`): le dan a cada entorno sus **secrets propios**, sus
  **reglas de protección** (en `production`, un revisor obligatorio) y un **historial de
  despliegues** (pestaña *Deployments*, con el commit de cada uno). Los secrets de PROD
  (`RENDER_HOOK_API_PROD`, `RENDER_HOOK_FRONT_PROD`) viven **dentro** del environment
  `production`: un job sólo los puede leer si pasó la aprobación. Si estuvieran en el repo, los
  podría leer cualquier job de cualquier rama del repositorio, incluida la de un PR, y un workflow modificado en
  una rama podría desplegar a PROD sin que nadie lo apruebe.
- **`&ref=$GITHUB_SHA` en el deploy hook**: sin eso, el hook despliega la **punta de la rama**,
  no el commit que acabo de verificar. Entre dos merges seguidos, la corrida de A terminaría
  desplegando B. En PROD pasa algo peor: entre que la corrida queda esperando y alguien aprueba
  pueden pasar horas y `main` se mueve; sin `ref`, aprobaría el commit A y subiría lo último que
  haya. Con `ref`, lo que se aprueba es lo que se despliega. Lo verifiqué mirando en Render qué
  commit tenía cada deploy después de cada corrida (columna del commit en *Deploys*).
- **`concurrency: deploy-prod` con `cancel-in-progress: false`**: dos despliegues a PROD nunca se
  pisan, y uno en curso no se cancela a la mitad. Ojo: no reordena las aprobaciones pendientes.
- **Auto-Deploy en `Off`** en los servicios de Render: los despliegues los dispara el pipeline
  con el hook, y la columna *Trigger* de *Deploys* dice `Deploy Hook` en todos. Si estuviera
  prendido, el proveedor desplegaría solo y el pipeline sería decorado.

### 4. Qué mira el aprobador antes de aprobar

No es un trámite. Antes de apretar *Approve* reviso:

1. **Que sea el commit que creo**: el commit de la corrida coincide con el del merge que acabo de
   hacer (`git rev-parse --short HEAD`), y no hay otra corrida más nueva esperando. Si hay una más
   nueva, apruebo la nueva y rechazo la vieja: aprobar la vieja haría **retroceder** PROD.
2. **Que QA haya quedado bien con el smoke nuevo**: en el log tiene que verse el `/health` con el
   commit de esta corrida, no sólo «verde».
3. **Qué cambia el commit**: abro el diff del PR. Si toca `migrate()` o el esquema, un rollback de
   código no deshace eso (sección 8), y no apruebo sin un plan para los datos.
4. **Lo que cambió, funcionando en QA**: abro la URL de QA y pruebo a mano lo modificado.

**Rechazo documentado (corrida #67):** el commit era 9b1584f, el PR #45, que sólo agregó el job deploy-prod al workflow: no cambió nada de backend/ ni de frontend/ (lo comprobé con git show --stat). No había ninguna versión nueva de la app para llevar a PROD, así que aprobar sólo habría redesplegado lo mismo. Rechacé el deploy.. Lo que pasó con
esa corrida: el job `deploy-prod` falló al ser rechazado, con `deploy-qa` ya desplegado; el
rechazo **no** saca nada de `main` ni revierte QA, sólo impide que ese commit llegue a PROD.
**Aprobación (corrida #69):** QA verde, *Waiting for review*, aprobación, PROD desplegado y smoke
verde.

**Lo que mi aprobador no puede ver:** el impacto real sobre usuarios (no hay métricas de errores ni
de latencia), el estado de la base ni qué pasaría con un esquema distinto, los logs de Render en
caliente ni el comportamiento bajo carga. Para decidir mejor me faltarían un panel de errores y
latencia por versión, alertas, y un ambiente de QA con datos parecidos a los reales.

### 5. La letra chica del free tier y cómo la maneja el pipeline

- **Render free:** la instancia se apaga tras unos 15 minutos sin tráfico y la primera petición
  después puede demorar 50 segundos o más (es el cartel que muestra el propio panel). Hay un tope de
  horas de ejecución compartido por workspace (750 h/mes) y de minutos de build por mes (500)
  según la documentación de Render al momento de escribir; con cuatro servicios y deploys
  frecuentes hay que cuidarlos. Render no tiene región en São Paulo: usé Virginia.
- **Neon free:** suspende el cómputo tras unos 5 minutos sin uso. Creé el proyecto primero en São
  Paulo y lo rehíce en N. Virginia para que la base y Render estén en la misma región (menos
  latencia, y las dos bases —QA y PROD— en el mismo lugar).
- **Cómo lo maneja el pipeline:** el smoke test **reintenta** hasta 30 veces con 20 s de espera
  (10 minutos) y cada `curl` tiene `--max-time 10`, porque un entorno dormido tarda en
  responder y no es un error. Aún así, un deploy «terminado» puede no estar respondiendo: por eso
  el smoke mira la app y no el estado del deploy.
- **Efecto colateral que encontré:** después de un rato sin uso, `/health` y `/api/gastos`
  respondían `{"status":"error"}`: Neon había suspendido el cómputo y el pool de conexiones de Go
  mantenía una conexión muerta. Lo resolví con `db.SetConnMaxIdleTime(time.Minute)`, que descarta
  conexiones ociosas antes de que Neon las corte.
- **Las bases son distintas de verdad:** QA y PROD usan dos bases (`app_qa` y `app_prod`) dentro del
  mismo proyecto de Neon, cada una con su propia cadena de conexión. Lo comprobé creando una fila
  `SOY PROD` sólo en la base de PROD y viendo que `/api/gastos` de PROD la devuelve y el de QA no.

### 6. Qué garantía pierdo porque Render reconstruye desde el repo

El pipeline publica en ghcr una imagen verificada, pero **Render no la usa**: construye su propia
imagen desde el repo, con el mismo commit (`&ref=`) y el mismo Dockerfile. Es **el mismo código,
no la misma imagen**. Lo que se pierde:

- **Lo que probé no es literalmente lo que corre.** La imagen de ghcr y la de Render tienen
  digests distintos. Entre los dos builds pueden cambiar cosas que el Dockerfile no fija: las
  imágenes base están por etiqueta móvil (`nginx:alpine`, `node:22-alpine`, `golang:1.22-alpine`)
  y el backend ejecuta `go mod tidy` al construir, que resuelve dependencias en ese momento y no
  desde un `go.sum` versionado.
- **El tiempo de despliegue incluye un build entero.** Cada despliegue paga un build en Render
  (y consume sus minutos de build del plan) además del que ya hizo el CI.
- **El commit que informa `/health` viene de dos lugares:** en la imagen de ghcr, del argumento de
  build `GIT_SHA`; en Render, de la variable `RENDER_GIT_COMMIT` que el proveedor inyecta.

La corrección es el TP7: que Render (o cualquier plataforma) **baje la imagen ya publicada, por
digest**, y no la reconstruya. Lo bueno es que el Dockerfile, el workflow y el contrato de
variables no cambian.

### 7. Qué prueba mi smoke test y qué no

El smoke contesta tres cosas: que la API responde (`/health`), que la base responde
(`/api/gastos` toca la base) y que el front se sirve (`/`). Además, ahora **compara el commit** que
informa `/health` con `$GITHUB_SHA` y reintenta hasta que coincide.

Por qué lo agregué: la primera versión sólo verificaba que respondiera, y **no decía qué versión
corría**. En una medición de PROD, el job terminó en 35 s mientras que el build del front en
Render tardó 58.4 s: el verde podía ser contra la versión vieja, porque Render mantiene viva la
anterior hasta que la nueva está lista. Con el cambio de este práctico, el log de la corrida #71
lo muestra: en QA el intento 1 fue `sin respuesta`, el intento 2 devolvió `{"status":"ok"}` **sin**
`commit` (la versión vieja, que con el smoke anterior ya habría dado verde) y el intento 3 devolvió
el commit correcto y recién ahí cerró. En PROD pasó lo mismo.
(Capturas: `img/18-tp6-smoke-qa-espera-commit.png`, `img/19-tp6-smoke-prod-espera-commit.png`.)

**Lo que sigue sin probar:** (a) la versión del **front**: es nginx estático y no informa nada,
así que puede quedar atrás de la API sin que el smoke lo note; (b) que la **imagen** sea la
verificada (sección 6); (c) que la app funcione bien en lo funcional: el smoke no crea ni borra
gastos; (d) los datos.

### 8. Patrón de despliegue para producción real y plan de rollback

**Patrón que elegiría: blue-green.** Mantener dos copias completas, la actual (blue) y la nueva
(green); se prueba green con el smoke **antes** de mandarle tráfico, y el cambio es mover el
balanceador. 

- *Por qué ese:* mi app es chica, con pocos usuarios, y lo que más valor me da es el **rollback
  inmediato** (volver el tráfico a blue, sin reconstruir nada) y probar la versión nueva con datos
  reales antes de exponerla. Un **canary** manda un porcentaje chico del tráfico a la versión
  nueva y mira las métricas: con pocos usuarios, el 5 % son un puñado de peticiones y no alcanzan
  para sacar una conclusión estadística.
- *Costo:* el doble de capacidad mientras conviven las dos copias, y la **base compartida**: las
  migraciones tienen que ser compatibles hacia atrás (expandir y después contraer), porque blue
  sigue corriendo contra el mismo esquema.
- *Riesgo:* el cambio de tráfico es todo o nada, así que un defecto que el smoke no ve llega a
  todos los usuarios de golpe.
- *Qué me falta hoy:* un balanceador que permita mover el tráfico entre dos copias (no lo
  tengo en el free tier de Render), métricas de errores y latencia por versión, alertas y logs centralizados. Esa observabilidad
  también es lo que me faltaría para hacer **canary en serio**: dividir el tráfico por peso, medir
  tasa de error y latencia por versión y decidir con eso.
- *Lo que ya se parece:* por lo que muestra el log del smoke (intento 2, sección 7), Render
  mantiene viva la versión anterior hasta que la nueva responde, y mi smoke hace de «prueba de
  green». Es un rolling básico.

**Plan de rollback actual** (si el deploy aprobado sale mal):

1. Detectar: smoke rojo, `/health` con error, o un problema reportado.
2. Identificar el **último commit bueno** en GitHub → *Deployments* → `production` (historial de
   qué commit estuvo desplegado y desde cuándo).
3. Llamar a los dos hooks de PROD con ese commit:
   `curl -fsS "$HOOK_API&ref=<commit-bueno>"` y lo mismo con el del front (los hooks salen del
   archivo privado, no del repo).
4. Verificar que `/health` devuelva la versión buena y que PROD responda.
5. **Revertir en `main`** (`git revert` por PR): si no, el próximo merge volvería a desplegar el
   código malo. El rollback con hook es temporal.

**Medido de verdad (no estimado)**, volviendo de `d3e03dd` al commit anterior `b73f624` y de nuevo
a `d3e03dd`, con los deploy hooks, en PROD:

| | Rollback a `b73f624` | Vuelta a `d3e03dd` |
|---|---|---|
| API, duración del deploy en Render | 39.7 s | 40.2 s |
| Front, duración del deploy en Render | 24.3 s | 24.0 s |
| API, desde el disparo hasta que `/health` responde la versión correcta | 45 s | 43 s |

El tiempo de rollback es el del servicio que más tarda: **39.7 s en Render, 45 s vistos desde
afuera**. El rastro está en *Deploys* de `gastos-api-prod` y `gastos-front-prod`
(`img/20-tp6-rollback-api-prod.png`, `img/21-tp6-rollback-front-prod.png`).
Es un **mejor caso**: `b73f624` ya se había construido en Render y reutilizó capas en caché; un
commit nunca construido, o con el caché vencido, tardaría más. Además se hizo **por fuera del
pipeline**, así que *Deployments* de GitHub sigue mostrando `d3e03dd` como activo aunque PROD haya
estado un rato en otro commit.

**Lo que un rollback de código no deshace: los datos.** Las filas escritas por la versión mala
quedan. `migrate()` sólo crea las tablas si no existen y no hay migraciones reversibles (deuda
anotada en el TP5): hoy es seguro porque no cambié el esquema, pero si un commit lo cambiara, volver
atrás no alcanzaría. Haría falta migraciones versionadas con *down*, el patrón expandir/contraer y
respaldos de la base antes de desplegar cambios de esquema. La métrica DORA que ejercito es el
**tiempo de recuperación** (*time to restore service*); el pipeline mismo mide la frecuencia de
despliegue y el tiempo de entrega de un cambio.

### 9. Dos entornos: qué va por variable y qué va en la imagen

- **Por variable (se lee al arrancar):** en el front, `BACKEND_URL` y `DNS_RESOLVER`, que nginx
  sustituye en la plantilla `default.conf.template` al iniciar (en Compose valen
  `http://backend:8080` y `127.0.0.11`; en Render, la URL https de la API y `8.8.8.8`). En el
  backend, `DATABASE_URL` y `PORT` (y `GIT_SHA`, que entra como argumento de build).
- **Dentro de la imagen:** el build estático de React (que llama a rutas relativas `/api/...`), la
  plantilla de nginx y el binario de Go. **La dirección del backend no está en la imagen del
  front**: la misma imagen sirve en QA y en PROD y cada una habla con su API.
- **Dos bases distintas** (una por entorno) y **cuatro servicios** de Render (API y front de QA y de
  PROD), todos con *Auto-Deploy* apagado: los despliega el pipeline.

**Si Render desaparece mañana:** sobreviven sin cambios los Dockerfiles, el workflow de CI, la
publicación en ghcr, el contrato de variables y el smoke. Migran la definición de los servicios, los
hooks y las URLs: cambia el «dónde» y no el «qué», porque el despliegue depende de variables y no de
nada de Render adentro de la imagen.

### 10. Problemas encontrados y cómo los resolví

- **Pegué la URL pública del servicio como valor del secret**, en lugar del *Deploy Hook*
  (`https://api.render.com/deploy/srv-...?key=...`, de *Settings → Deploy Hook*): el pipeline habría
  llamado a la página pública en vez de disparar un deploy. Lo corregí antes de usarlo, con el
  valor correcto en cada secret.
- **La cadena de conexión de Neon no conectaba:** tenía `-pooler` en el host, `channel_binding=require`
  (que `lib/pq` no soporta) y, en un intento, faltaba el `@`. La correcta usa el host directo, `?sslmode=require`
  y se copia de Neon con *pooling* apagado.
- **Una contraseña de la base quedó visible en una captura**: la reseteé y desde entonces oculto los
  valores de las variables en las capturas y no pego hooks ni cadenas de conexión.
- **São Paulo en Neon y Virginia en Render:** distinta región entre base y servicio. Rehíce el
  proyecto de Neon en N. Virginia.
- **Faltaba un servicio en Render:** el cuarto estaba dentro de otro proyecto; no se había
  perdido.
- **`{"status":"error"}` tras un rato inactivo:** conexión muerta tras la suspensión de Neon (sección 5).
- **GitHub tardó unos 10 minutos en crear la corrida del push a `main`**: no era un error de mi
  configuración.
- **Un `deploy-qa` salteado en un PR:** es lo esperado (`if` sólo en `main`).
- **El smoke daba verde contra la versión vieja** (sección 7).
- **Un error mío al probar el cambio de `/health`:** puse `ARG GIT_SHA=desconocido` en el Dockerfile,
  lo que dejaba la variable siempre definida y el fallback a `RENDER_GIT_COMMIT` nunca se
  ejecutaba (en Render habría pasado lo mismo). Lo vi probando las tres ramas del handler en local
  con `docker compose` antes de abrir el PR; el default del `ARG` ahora es vacío.
- **Dejé PROD en el commit anterior unas dos horas después de probar el rollback**, antes de volver
  al actual. Lo noté al mirar las horas de *Deploys* y lo anoto como el riesgo de hacer un rollback
  por fuera del pipeline: nada obliga a deshacerlo, y *Deployments* de GitHub no lo registra.

### 11. Declaración de uso de IA

Usé IA (Claude) para armar el TP **paso a paso**, un PR por etapa: me propuso el diseño de la
cadena y los workflows de QA y PROD (`needs`, `if`, environments, `concurrency`, el `&ref=` en los
hooks), la plantilla de nginx con variables, el smoke con reintentos y el cambio de `/health` que
informa el commit, y me ayudó a diagnosticar los problemas de la sección 10. No me limité a
copiar: cada paso lo ejecuté yo en mi máquina, en Render, en Neon y en GitHub, y verifiqué con mis
propios ojos el resultado: las imágenes públicas (`docker pull` sin credenciales), QA desplegándose
sola con el commit correcto, el *Waiting for review*, un rechazo y una aprobación, la separación de las
bases con una fila de prueba, y el rollback cronometrado con los hooks y *Deploys* de Render.

Claude se equivocó en dos cosas que aparecieron al medir y probar: afirmó que un deploy de QA de
34 s implicaba que el smoke había corrido contra la versión vieja, y en realidad el build de Render
estaba en caché; y su primera versión del Dockerfile, con `ARG GIT_SHA=desconocido`, anulaba el
fallback a `RENDER_GIT_COMMIT`. Puedo explicar sin el texto generado qué pasa exactamente entre el merge y PROD:
los dos builds publican la imagen en ghcr; `deploy-qa` dispara los hooks con el commit verificado y
espera con el smoke; `deploy-prod` se pausa en el environment hasta que apruebo, y recién ahí usa los
secrets de PROD; y por qué la aprobación sirve sólo si miro lo que dice la sección 4.
