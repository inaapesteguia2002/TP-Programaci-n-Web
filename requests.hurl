# 1. GESTIÓN DE EQUIPOS

# ERROR: Intentar crear equipo con nombre vacío (debe dar 400 Bad Request)
POST http://localhost:8080/equipos
{
  "nombre": ""
}
HTTP/1.1 400

# POST: Crear el equipo local (ID 1)
POST http://localhost:8080/equipos
{
  "nombre": "Boca"
}
HTTP/1.1 201
[Asserts]
jsonpath "$.ID" == 1
jsonpath "$.Nombre" == "Boca"

# PUT: Actualizar el equipo 1 y comprobar que devuelve el objeto actualizado
PUT http://localhost:8080/equipos/1
{
  "nombre": "Boca Juniors"
}
HTTP/1.1 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.nombre" == "Boca Juniors"

# POST: Crear el equipo visitante (ID 2)
POST http://localhost:8080/equipos
{
  "nombre": "River Plate"
}
HTTP/1.1 201
[Asserts]
jsonpath "$.ID" == 2
jsonpath "$.Nombre" == "River Plate"

# GET: Obtener equipo por ID
GET http://localhost:8080/equipos/1
HTTP/1.1 200
[Asserts]
jsonpath "$.Nombre" == "Boca Juniors"

# GET: Listar todos los equipos
GET http://localhost:8080/equipos
HTTP/1.1 200
[Asserts]
jsonpath "$" count == 2


# 2. GESTIÓN DE LIGAS

# ERROR: Intentar crear liga sin país (400 Bad Request)
POST http://localhost:8080/ligas
{
  "nombre": "Liga Incompleta",
  "pais": ""
}
HTTP/1.1 400

# POST: Crear una liga (ID 1)
POST http://localhost:8080/ligas
{
  "nombre": "Primera División",
  "pais": "Argentina"
}
HTTP/1.1 201
[Asserts]
jsonpath "$.ID" == 1
jsonpath "$.Nombre" == "Primera División"

# PUT: Actualizar liga 1 y comprobar objeto devuelto
PUT http://localhost:8080/ligas/1
{
  "nombre": "Liga Profesional",
  "pais": "Argentina"
}
HTTP/1.1 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.nombre" == "Liga Profesional"
jsonpath "$.pais" == "Argentina"

# GET: Listar ligas
GET http://localhost:8080/ligas
HTTP/1.1 200
[Asserts]
jsonpath "$" count == 1


# 3. GESTIÓN DE PARTIDOS

# ERROR: Intentar crear partido con IDs inválidos (400 Bad Request)
POST http://localhost:8080/partidos
{
  "liga_id": 0,
  "equipo_local_id": 1,
  "equipo_visitante_id": 2
}
HTTP/1.1 400

# POST: Registrar un partido (ID 1)
POST http://localhost:8080/partidos
{
  "liga_id": 1,
  "equipo_local_id": 1,
  "equipo_visitante_id": 2,
  "goles_local": 0,
  "goles_visitante": 0,
  "estado": "pendiente"
}
HTTP/1.1 201
[Asserts]
jsonpath "$.ID" == 1

# PUT: Actualizar resultado del partido y comprobar objeto devuelto
PUT http://localhost:8080/partidos/1
{
  "goles_local": 2,
  "goles_visitante": 0,
  "estado": "finalizado"
}
HTTP/1.1 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.goles_local" == 2
jsonpath "$.goles_visitante" == 0
jsonpath "$.estado" == "finalizado"

# GET: Listar partidos
GET http://localhost:8080/partidos
HTTP/1.1 200
[Asserts]
jsonpath "$" count == 1


# 4. TABLA DE POSICIONES

# [ERROR] Intentar registrar en tabla con IDs inválidos (400 Bad Request)
POST http://localhost:8080/tabla
{
  "liga_id": 0,
  "equipo_id": 1
}
HTTP/1.1 400

# POST: Registrar equipo 1 en la tabla de la liga 1
POST http://localhost:8080/tabla
{
  "liga_id": 1,
  "equipo_id": 1
}
HTTP/1.1 201
[Asserts]
jsonpath "$.LigaID" == 1
jsonpath "$.EquipoID" == 1

# PUT: Actualizar estadísticas en la tabla y comprobar objeto devuelto
PUT http://localhost:8080/tabla-equipo
{
  "liga_id": 1,
  "equipo_id": 1,
  "puntos": 3,
  "partidos_ganados": 1,
  "partidos_empatados": 0,
  "partidos_perdidos": 0
}
HTTP/1.1 200
[Asserts]
jsonpath "$.liga_id" == 1
jsonpath "$.equipo_id" == 1
jsonpath "$.puntos" == 3
jsonpath "$.partidos_ganados" == 1

# GET: Consultar la tabla de posiciones de la liga 1
GET http://localhost:8080/tabla/1
HTTP/1.1 200
[Asserts]
jsonpath "$" count == 1


# 5. ELIMINACIÓN (DELETE) Y COMPROBACIÓN ERROR 404

# 1. Eliminar equipo de la tabla de posiciones
DELETE http://localhost:8080/tabla-equipo
{
  "liga_id": 1,
  "equipo_id": 1
}
HTTP/1.1 204

# 2. Eliminar el partido 1
DELETE http://localhost:8080/partidos/1
HTTP/1.1 204

# Comprobar que el partido ya no existe (404 Not Found)
GET http://localhost:8080/partidos/1
HTTP/1.1 404

# 3. Eliminar la liga 1
DELETE http://localhost:8080/ligas/1
HTTP/1.1 204

# Comprobar que la liga ya no existe (404 Not Found)
GET http://localhost:8080/ligas/1
HTTP/1.1 404

# 4. Eliminar el equipo 1
DELETE http://localhost:8080/equipos/1
HTTP/1.1 204

# Comprobar que el equipo ya no existe (404 Not Found)
GET http://localhost:8080/equipos/1
HTTP/1.1 404
HTTP/1.1 200
