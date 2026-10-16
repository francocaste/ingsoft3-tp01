import { describe, it, expect, vi } from 'vitest'
import {
  calcularTotal,
  claseDePorcentaje,
  enviarGasto,
  formularioValido,
  hoyISO,
  mensajeDeError,
} from './gastos'

// Todos los tests siguen Arrange / Act / Assert.
const HOY = '2026-10-09'

describe('hoyISO', () => {
  it('devuelve la fecha en formato YYYY-MM-DD', () => {
    const ahora = new Date('2026-10-09T15:30:00Z')
    expect(hoyISO(ahora)).toBe('2026-10-09')
  })

  it('sin argumentos usa la fecha actual', () => {
    expect(hoyISO()).toMatch(/^\d{4}-\d{2}-\d{2}$/)
  })
})

describe('formularioValido (parametrizado)', () => {
  it.each([
    ['todo completo', { descripcion: 'Almuerzo', monto: '1500', fecha: HOY }, true],
    ['fecha pasada', { descripcion: 'Cafe', monto: '300', fecha: '2026-01-01' }, true],
    ['descripcion vacia', { descripcion: '', monto: '1500', fecha: HOY }, false],
    ['descripcion solo espacios', { descripcion: '   ', monto: '1500', fecha: HOY }, false],
    ['monto cero', { descripcion: 'Cafe', monto: '0', fecha: HOY }, false],
    ['monto negativo', { descripcion: 'Cafe', monto: '-5', fecha: HOY }, false],
    ['monto vacio', { descripcion: 'Cafe', monto: '', fecha: HOY }, false],
    ['sin fecha', { descripcion: 'Cafe', monto: '10', fecha: '' }, false],
    ['fecha futura', { descripcion: 'Cafe', monto: '10', fecha: '2026-10-10' }, false],
  ])('%s -> %s', (_nombre, datos, esperado) => {
    expect(formularioValido(datos, HOY)).toBe(esperado)
  })
})

describe('calcularTotal', () => {
  it('suma los montos (aunque vengan como string)', () => {
    const gastos = [{ monto: 100 }, { monto: '50.5' }, { monto: 20 }]
    expect(calcularTotal(gastos)).toBe(170.5)
  })

  it('con lista vacia devuelve 0', () => {
    expect(calcularTotal([])).toBe(0)
  })
})

describe('mensajeDeError', () => {
  it('el 422 se muestra como presupuesto excedido', () => {
    expect(mensajeDeError(422, 'tope superado')).toBe('Presupuesto excedido: tope superado')
  })

  it('cualquier otro status devuelve el texto tal cual', () => {
    expect(mensajeDeError(400, 'categoria invalida')).toBe('categoria invalida')
  })
})

describe('claseDePorcentaje', () => {
  it.each([
    [50, ''],
    [100, ''],
    [100.5, 'excedido'],
    [180, 'excedido'],
  ])('%s%% -> "%s"', (porcentaje, clase) => {
    expect(claseDePorcentaje(porcentaje)).toBe(clase)
  })
})

describe('enviarGasto (con fetch simulado)', () => {
  const datos = { descripcion: 'Almuerzo', monto: '1500', categoria: 'Comida', fecha: HOY }

  it('manda un POST a /api/gastos con el body correcto', async () => {
    // Arrange: un doble de fetch que responde OK, sin salir a la red.
    const fetchFalso = vi.fn().mockResolvedValue({ ok: true })
    // Act
    const resultado = await enviarGasto(fetchFalso, datos, HOY)
    // Assert: verificamos el COMPORTAMIENTO (como se llamo), no solo el resultado.
    expect(resultado).toEqual({ ok: true })
    expect(fetchFalso).toHaveBeenCalledTimes(1)
    const [url, opciones] = fetchFalso.mock.calls[0]
    expect(url).toBe('/api/gastos')
    expect(opciones.method).toBe('POST')
    expect(JSON.parse(opciones.body)).toMatchObject({
      descripcion: 'Almuerzo',
      monto: 1500,
      categoria: 'Comida',
    })
  })

  it('si el backend responde 422, devuelve el error de presupuesto excedido', async () => {
    const fetchFalso = vi.fn().mockResolvedValue({
      ok: false,
      status: 422,
      text: async () => 'el gasto excede el presupuesto mensual de la categoria',
    })

    const resultado = await enviarGasto(fetchFalso, datos, HOY)

    expect(resultado.ok).toBe(false)
    expect(resultado.error).toBe(
      'Presupuesto excedido: el gasto excede el presupuesto mensual de la categoria',
    )
  })

  it('si el formulario es invalido, NO llama al backend', async () => {
    const fetchFalso = vi.fn()

    const resultado = await enviarGasto(fetchFalso, { ...datos, monto: '0' }, HOY)

    expect(resultado).toEqual({ ok: false, error: 'Formulario invalido' })
    expect(fetchFalso).not.toHaveBeenCalled()
  })

  it('si la red falla, el error se propaga', async () => {
    const fetchFalso = vi.fn().mockRejectedValue(new Error('sin conexion'))

    await expect(enviarGasto(fetchFalso, datos, HOY)).rejects.toThrow('sin conexion')
  })
})