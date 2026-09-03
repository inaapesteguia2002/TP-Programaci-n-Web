-- CRUD PARA LA TABLA: EQUIPOS

-- name: CreateEquipo :one
-- Crea un nuevo registro
INSERT INTO equipos (nombre) 
VALUES ($1) 
RETURNING *;

-- name: GetEquipo :one
-- Obtiene un registro especifico por su ID
SELECT * FROM equipos 
WHERE id = $1;

-- name: ListEquipos :many
-- Lista todos los registros
SELECT * FROM equipos 
ORDER BY nombre ASC;

-- name: UpdateEquipo :exec
-- Actualiza un registro existente
UPDATE equipos 
SET nombre = $2 
WHERE id = $1;

-- name: DeleteEquipo :exec
-- Borra un registro por su ID
DELETE FROM equipos 
WHERE id = $1;


-- CRUD PARA LA TABLA PRINCIPAL: PARTIDOS

-- name: CreatePartido :one
-- Crea un partido pasándole todos los datos (goles y estado)
INSERT INTO partidos (liga_id, equipo_local_id, equipo_visitante_id, goles_local, goles_visitante, estado)
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: GetPartido :one
SELECT * FROM partidos 
WHERE id = $1;

-- name: ListPartidos :many
SELECT * FROM partidos 
ORDER BY id DESC;

-- name: UpdatePartido :exec
-- Ideal para cargar el resultado y cambiar el estado manualmente
UPDATE partidos 
SET goles_local = $2, 
    goles_visitante = $3, 
    estado = $4 
WHERE id = $1;

-- name: DeletePartido :exec
DELETE FROM partidos
WHERE id = $1;


-- FUNCIONES ESPECIFICAS DE LA APP 

-- name: CreateLiga :one
INSERT INTO ligas (nombre, pais) 
VALUES ($1, $2) RETURNING *;

-- name: CreatePartidoProgramado :one
-- Alternativa rápida: Crea un partido que aún no se jugó (solo pide liga y equipos)
INSERT INTO partidos (liga_id, equipo_local_id, equipo_visitante_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: ActualizarResultadoPartido :exec
-- Atajo: Solo pasas los goles y automáticamente lo marca como 'Finalizado'
UPDATE partidos 
SET goles_local = $2, goles_visitante = $3, estado = 'Finalizado'
WHERE id = $1;

-- name: ListarPartidosPorLiga :many
-- Para mostrar todos los partidos de una liga específica
SELECT * FROM partidos 
WHERE liga_id = $1;

-- name: InicializarEquipoEnTabla :one
-- Agrega un equipo a la tabla de una liga con 0 puntos
INSERT INTO tabla_posiciones (liga_id, equipo_id)
VALUES ($1, $2) RETURNING *;

-- name: ActualizarEstadisticas :exec
-- Para sumar puntos y partidos cuando termina un encuentro
UPDATE tabla_posiciones 
SET puntos = puntos + $3, 
    partidos_jugados = partidos_jugados + 1,
    partidos_ganados = partidos_ganados + $4,
    partidos_empatados = partidos_empatados + $5,
    partidos_perdidos = partidos_perdidos + $6
WHERE liga_id = $1 AND equipo_id = $2;

-- name: ObtenerTablaDePosiciones :many
-- ¡El corazón de Promiedos! Trae la tabla ordenada por puntos de mayor a menor.
SELECT e.nombre as equipo, t.puntos, t.partidos_jugados, t.partidos_ganados, t.partidos_empatados, t.partidos_perdidos
FROM tabla_posiciones t
JOIN equipos e ON t.equipo_id = e.id
WHERE t.liga_id = $1
ORDER BY t.puntos DESC;
