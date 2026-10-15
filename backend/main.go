package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gastos-backend/internal/gastos"

	_ "github.com/lib/pq"
)

var (
	db       *sql.DB
	servicio *gastos.Servicio
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("falta la variable de entorno DATABASE_URL")
	}

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("no se pudo abrir la conexion: %v", err)
	}
	defer db.Close()

	if err := waitForDB(db, 15); err != nil {
		log.Fatalf("la base de datos no respondio a tiempo: %v", err)
	}

	if err := migrate(db); err != nil {
		log.Fatalf("no se pudo crear el schema: %v", err)
	}

	servicio = gastos.NuevoServicio(&repoSQL{db: db}, time.Now)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/gastos", listGastosHandler)
	mux.HandleFunc("POST /api/gastos", createGastoHandler)
	mux.HandleFunc("DELETE /api/gastos/{id}", deleteGastoHandler)
	mux.HandleFunc("PATCH /api/gastos/{id}/estado", updateEstadoHandler)
	mux.HandleFunc("GET /api/resumen", resumenHandler)
	mux.HandleFunc("GET /api/presupuestos", listPresupuestosHandler)
	mux.HandleFunc("POST /api/presupuestos", upsertPresupuestoHandler)

	log.Println("escuchando en :8080")
	if err := http.ListenAndServe(":8080", withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func waitForDB(db *sql.DB, attempts int) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return err
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS gastos (
			id SERIAL PRIMARY KEY,
			descripcion TEXT NOT NULL,
			monto NUMERIC(12,2) NOT NULL,
			categoria TEXT NOT NULL,
			fecha DATE NOT NULL DEFAULT CURRENT_DATE,
			estado TEXT NOT NULL DEFAULT 'pendiente'
		);

		CREATE TABLE IF NOT EXISTS presupuestos (
			categoria TEXT PRIMARY KEY,
			monto_mensual NUMERIC(12,2) NOT NULL
		);
	`)
	return err
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, `{"status":"error"}`, http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func listGastosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, descripcion, monto, categoria, fecha, estado FROM gastos ORDER BY fecha DESC, id DESC`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	lista := []gastos.Gasto{}
	for rows.Next() {
		var g gastos.Gasto
		if err := rows.Scan(&g.ID, &g.Descripcion, &g.Monto, &g.Categoria, &g.Fecha, &g.Estado); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lista = append(lista, g)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, lista)
}

func createGastoHandler(w http.ResponseWriter, r *http.Request) {
	var g gastos.Gasto
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "body invalido", http.StatusBadRequest)
		return
	}

	// Las reglas 1 a 4 viven en gastos.Servicio; el handler solo traduce
	// el error de dominio a un status HTTP.
	creado, err := servicio.CrearGasto(g)
	if err != nil {
		http.Error(w, err.Error(), statusDe(err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, creado)
}

func deleteGastoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}
	res, err := db.Exec(`DELETE FROM gastos WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "no existe ese gasto", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// La regla 5 (transicion de estado) vive en gastos.Servicio.CambiarEstado.
func updateEstadoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}

	var body struct {
		Estado string `json:"estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "body invalido", http.StatusBadRequest)
		return
	}

	if err := servicio.CambiarEstado(id, body.Estado); err != nil {
		http.Error(w, err.Error(), statusDe(err))
		return
	}
	writeJSON(w, map[string]string{"id": idStr, "estado": body.Estado})
}

func resumenHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT categoria, SUM(monto) FROM gastos GROUP BY categoria ORDER BY SUM(monto) DESC`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	resumen := []gastos.ResumenCategoria{}
	for rows.Next() {
		var rc gastos.ResumenCategoria
		if err := rows.Scan(&rc.Categoria, &rc.Total); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		presupuesto, tiene, err := presupuestoDeCategoria(rc.Categoria)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tiene {
			rc.PresupuestoMensual = &presupuesto
			porcentaje := gastos.PorcentajeUsado(rc.Total, presupuesto)
			rc.PorcentajeUsado = &porcentaje
		}

		resumen = append(resumen, rc)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, resumen)
}

func listPresupuestosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT categoria, monto_mensual FROM presupuestos ORDER BY categoria`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	presupuestos := []gastos.Presupuesto{}
	for rows.Next() {
		var p gastos.Presupuesto
		if err := rows.Scan(&p.Categoria, &p.MontoMensual); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		presupuestos = append(presupuestos, p)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, presupuestos)
}

func upsertPresupuestoHandler(w http.ResponseWriter, r *http.Request) {
	var p gastos.Presupuesto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "body invalido", http.StatusBadRequest)
		return
	}
	if !gastos.CategoriasValidas[p.Categoria] {
		http.Error(w, "categoria invalida", http.StatusBadRequest)
		return
	}
	if p.MontoMensual <= 0 {
		http.Error(w, "monto_mensual debe ser mayor a 0", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
		INSERT INTO presupuestos (categoria, monto_mensual) VALUES ($1, $2)
		ON CONFLICT (categoria) DO UPDATE SET monto_mensual = EXCLUDED.monto_mensual
	`, p.Categoria, p.MontoMensual)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p)
}

// --- Helpers ---

func presupuestoDeCategoria(categoria string) (monto float64, tiene bool, err error) {
	err = db.QueryRow(`SELECT monto_mensual FROM presupuestos WHERE categoria = $1`, categoria).Scan(&monto)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return monto, true, nil
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// statusDe traduce un error de dominio al status HTTP que corresponde.
// Cualquier error que no sea de dominio (ej: la base caida) es un 500.
func statusDe(err error) int {
	switch {
	case errors.Is(err, gastos.ErrCamposObligatorios),
		errors.Is(err, gastos.ErrCategoriaInvalida),
		errors.Is(err, gastos.ErrFechaFutura),
		errors.Is(err, gastos.ErrEstadoInvalido):
		return http.StatusBadRequest
	case errors.Is(err, gastos.ErrPresupuestoExcedido):
		return http.StatusUnprocessableEntity
	case errors.Is(err, gastos.ErrTransicionInvalida):
		return http.StatusConflict
	case errors.Is(err, gastos.ErrNoEncontrado):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
