# 1. GESTIÓN DE EQUIPOS

# Crear el equipo local (ID 1)
POST http://localhost:8080/equipos
{
  "nombre": "Boca Juniors"
}
HTTP/1.1 201

# Crear el equipo visitante (ID 2)
POST http://localhost:8080/equipos
{
  "nombre": "River Plate"
}
HTTP/1.1 201

# Listar todos los equipos
GET http://localhost:8080/equipos
HTTP/1.1 200


# 2. GESTIÓN DE LIGAS

# Crear una liga (ID 1)
POST http://localhost:8080/ligas
{
  "nombre": "Liga Profesional",
  "pais": "Argentina"
}
HTTP/1.1 201

# Listar ligas
GET http://localhost:8080/ligas
HTTP/1.1 200


# 3. GESTIÓN DE PARTIDOS

# Registrar un partido
POST http://localhost:8080/partidos
{
  "liga_id": 1,
  "equipo_local_id": 1,
  "equipo_visitante_id": 2,
  "goles_local": 2,
  "goles_visitante": 0,
  "estado": "finalizado"
}
HTTP/1.1 201

# Listar partidos
GET http://localhost:8080/partidos
HTTP/1.1 200


# 4. TABLA DE POSICIONES

# Registrar un equipo en la tabla de una liga
POST http://localhost:8080/tabla
{
  "liga_id": 1,
  "equipo_id": 1
}
HTTP/1.1 201

# Consultar la tabla de posiciones de la liga 1
GET http://localhost:8080/tabla/1
HTTP/1.1 200