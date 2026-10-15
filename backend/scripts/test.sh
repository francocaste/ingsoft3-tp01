#!/bin/sh
# Corre los tests unitarios del backend, mide cobertura y FALLA (exit 1) si:
#   - algun test falla, o
#   - la cobertura de sentencias queda por debajo del umbral.
#
# Alcance de la cobertura: SOLO ./internal/... (reglas de negocio + servicio).
# Excluido a proposito: main.go (arranque, rutas, CORS), repo_sql.go (SQL) y los
# handlers HTTP. Son codigo de "pegamento" que necesita una base real o un
# servidor; se prueba con tests de integracion (otro tipo de test), no unitarios.
set -eu

MIN="${COVERAGE_MIN:-90}"   # umbral minimo de cobertura de sentencias (%)
OUT="${OUT_DIR:-/out}"      # carpeta donde dejamos los reportes
mkdir -p "$OUT"

status=0
go test -count=1 -v -coverprofile="$OUT/coverage.out" ./internal/... >"$OUT/test.log" 2>&1 || status=$?
cat "$OUT/test.log"

if [ "$status" -ne 0 ]; then
  {
    echo "### Backend: tests"
    echo ""
    echo "FALLARON los tests (exit $status). Ver test.log en el artifact."
  } >"$OUT/resumen.md"
  exit 1
fi

go tool cover -func="$OUT/coverage.out" >"$OUT/cobertura-por-funcion.txt"
go tool cover -html="$OUT/coverage.out" -o "$OUT/coverage.html"
TOTAL=$(awk '/^total:/ { sub("%", "", $3); print $3 }' "$OUT/cobertura-por-funcion.txt")
PASARON=$(grep -c '^--- PASS' "$OUT/test.log" || true)

if awk -v t="$TOTAL" -v m="$MIN" 'BEGIN { exit !(t + 0 >= m + 0) }'; then
  VEREDICTO="OK"; CODE=0
else
  VEREDICTO="POR DEBAJO DEL UMBRAL"; CODE=1
fi

{
  echo "### Backend: tests y cobertura"
  echo ""
  echo "| Metrica | Valor |"
  echo "|---|---|"
  echo "| Tests (funciones) que pasaron | $PASARON |"
  echo "| Cobertura de sentencias (statements) | **${TOTAL}%** |"
  echo "| Umbral minimo | ${MIN}% |"
  echo "| Veredicto | **${VEREDICTO}** |"
} >"$OUT/resumen.md"

echo ""
echo "Cobertura: ${TOTAL}% (umbral ${MIN}%) -> ${VEREDICTO}"
exit "$CODE"