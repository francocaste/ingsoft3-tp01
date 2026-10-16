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
  - `img/20-demo2-test-borde-falla.png`: demo 2, el caso `justo_en_el_tope_(se_permite)` fallando.
  - `img/18-resumen-cobertura.png`: resumen de cobertura publicado en el *Summary* del run.
  - `img/19-artifacts.png`: artefactos `cobertura-backend` y `cobertura-frontend` del run.

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