-- CRUD PARA LA TABLA: EQUIPOS

-- name: CreateEquipo :one
INSERT INTO equipos (nombre) 
VALUES ($1) 
RETURNING *; --Devuelve el id asiginado del nuevo equipo creado

-- name: GetEquipo :one
SELECT * FROM equipos 
WHERE id = $1;

-- name: ListEquipos :many
SELECT * FROM equipos 
ORDER BY nombre ASC;

-- name: UpdateEquipo :exec
UPDATE equipos 
SET nombre = $2 
WHERE id = $1;

-- name: DeleteEquipo :exec
DELETE FROM equipos 
WHERE id = $1;


-- CRUD PARA LA TABLA PRINCIPAL: PARTIDOS

-- name: CreatePartido :one
INSERT INTO partidos (liga_id, equipo_local_id, equipo_visitante_id, goles_local, goles_visitante, estado)
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *; --Devuelve el id asiginado del nuevo partido creado

-- name: GetPartido :one
SELECT * FROM partidos 
WHERE id = $1;

-- name: ListPartidos :many
-- lista todos los partidos del dia sin importar la liga a la que pertezcan
SELECT * FROM partidos 
ORDER BY id DESC;

-- name: UpdatePartido :exec
UPDATE partidos 
SET goles_local = $2, 
    goles_visitante = $3, 
    estado = $4 
WHERE id = $1;

-- name: DeletePartido :exec
DELETE FROM partidos
WHERE id = $1;


-- CRUD PARA LA TABLA PRINCIPAL: LIGAS

-- name: CreateLiga :one
INSERT INTO ligas (nombre, pais) 
VALUES ($1, $2) 
RETURNING *; --Devuelve el id asiginado de la nueva liga creada

-- name: GetLiga :one
SELECT * FROM ligas
WHERE id = $1;

-- name: ListLigas :many
SELECT * FROM ligas
ORDER BY nombre ASC;

-- name: ListPartidosLiga :many
-- lista los partidos de una liga en especifico, sin importar el estado del partido
SELECT * FROM partidos 
WHERE liga_id = $1;

-- name: DeleteLiga :exec
DELETE FROM ligas
WHERE id = $1;


-- CRUD PARA LA TABLA PRINCIPAL: TABLA DE POSICIONES

-- name: CreateTabla :one
INSERT INTO tabla_posiciones (liga_id, equipo_id)
VALUES ($1, $2) 
RETURNING *; --Devuelve el id de la fila creada en la tabla de posiciones (todos los ids de cada componente y las estadisticas iniciales en 0)

-- name: UpdateTabla :exec
UPDATE tabla_posiciones 
SET puntos = puntos + $3, 
    partidos_jugados = partidos_jugados + 1,
    partidos_ganados = partidos_ganados + $4,
    partidos_empatados = partidos_empatados + $5,
    partidos_perdidos = partidos_perdidos + $6
WHERE liga_id = $1 AND equipo_id = $2;

-- name: GetTabla :many
SELECT e.nombre as equipo, t.puntos, t.partidos_jugados, t.partidos_ganados, t.partidos_empatados, t.partidos_perdidos
FROM tabla_posiciones t
JOIN equipos e ON t.equipo_id = e.id
WHERE t.liga_id = $1
ORDER BY t.puntos DESC;

-- name: DeleteEquipoDeTabla :exec
-- elimina un equipo de la tabla de posiciones de una liga en especifico
DELETE FROM tabla_posiciones
WHERE liga_id = $1 AND equipo_id = $2;
