// Logica pura del frontend (sin React, sin DOM, sin red directa).
// Vive aca, separada de App.jsx, para poder probarla con tests unitarios.

export const CATEGORIAS = ['Comida', 'Transporte', 'Ocio', 'Servicios', 'Salud', 'Otros']

// Fecha de hoy en formato YYYY-MM-DD. "ahora" es inyectable para los tests.
export function hoyISO(ahora = new Date()) {
  return ahora.toISOString().slice(0, 10)
}

// El formulario es invalido si falta algun dato o la fecha es futura.
// Con esto se deshabilita el boton en vez de dejar que el usuario descubra
// el error recien despues de enviarlo.
export function formularioValido({ descripcion, monto, fecha }, hoy = hoyISO()) {
  return descripcion.trim() !== '' && Number(monto) > 0 && fecha !== '' && fecha <= hoy
}

export function calcularTotal(gastos) {
  return gastos.reduce((acc, g) => acc + Number(g.monto), 0)
}

// El backend devuelve 422 puntualmente cuando se excede el presupuesto de la
// categoria: se muestra distinto al resto de los errores de validacion.
export function mensajeDeError(status, texto) {
  if (status === 422) {
    return `Presupuesto excedido: ${texto}`
  }
  return texto
}

// Clase CSS del porcentaje de presupuesto usado: rojo si se paso del 100%.
export function claseDePorcentaje(porcentaje) {
  return porcentaje > 100 ? 'excedido' : ''
}

// Envia un gasto nuevo. Recibe la funcion de red (fetch) por parametro:
// en la app es el fetch real; en los tests, un doble (vi.fn()) que no sale a
// la red. Devuelve { ok: true } o { ok: false, error: '<mensaje para la UI>' }.
export async function enviarGasto(fetchFn, datos, hoy = hoyISO()) {
  if (!formularioValido(datos, hoy)) {
    return { ok: false, error: 'Formulario invalido' }
  }
  const res = await fetchFn('/api/gastos', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      descripcion: datos.descripcion,
      monto: parseFloat(datos.monto),
      categoria: datos.categoria,
      fecha: new Date(datos.fecha).toISOString(),
    }),
  })
  if (!res.ok) {
    return { ok: false, error: mensajeDeError(res.status, await res.text()) }
  }
  return { ok: true }
}

