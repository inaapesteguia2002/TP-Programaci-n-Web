# 1. CARGA INICIAL DE LA PÁGINA PRINCIPAL (SSR HTML)

GET http://localhost:8080/
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "De Rabona"
body contains "<form"
body contains "No hay ligas registradas."
body contains "No hay equipos registrados."


# 2. GESTIÓN DE LIGAS (FORMULARIO CLÁSICO + PRG)

# [ERROR] Intentar crear liga con datos vacíos
POST http://localhost:8080/ligas
[FormParams]
nombre: Liga Incompleta
pais: 
HTTP 400

# [POST] Crear Liga 1 vía formulario y verificar redirección PRG (303 See Other)
POST http://localhost:8080/ligas
[FormParams]
nombre: Liga Profesional
pais: Argentina
HTTP 303
[Asserts]
header "Location" contains "/"

# [GET] Verificar que la liga aparece en la tabla HTML
GET http://localhost:8080/
HTTP 200
[Asserts]
body contains "<strong>Liga Profesional</strong>"
body contains "<td>Argentina</td>"


# 3. GESTIÓN DE EQUIPOS (FORMULARIO CLÁSICO + PRG)

# [ERROR] Intentar crear equipo con nombre vacío
POST http://localhost:8080/equipos
[FormParams]
nombre: 
HTTP 400

# [POST] Crear Equipo 1 (Boca Juniors)
POST http://localhost:8080/equipos
[FormParams]
nombre: Boca Juniors
HTTP 303
[Asserts]
header "Location" contains "/"

# [POST] Crear Equipo 2 (River Plate)
POST http://localhost:8080/equipos
[FormParams]
nombre: River Plate
HTTP 303
[Asserts]
header "Location" contains "/"

# [GET] Verificar que ambos equipos se renderizan en las filas de la tabla HTML
GET http://localhost:8080/
HTTP 200
[Asserts]
body contains "<strong>Boca Juniors</strong>"
body contains "<strong>River Plate</strong>"


# 4. INSCRIPCIÓN EN TABLA DE POSICIONES (FORMULARIO CLÁSICO + PRG)

# [ERROR] Consultar tabla con liga_id inválido en query param
GET http://localhost:8080/?liga_id=invalido
HTTP 400

# [POST] Inscribir Equipo 1 en la Liga 1
POST http://localhost:8080/tabla
[FormParams]
liga_id: 1
equipo_id: 1
HTTP 303
[Asserts]
header "Location" contains "/"

# [POST] Inscribir Equipo 2 en la Liga 1
POST http://localhost:8080/tabla
[FormParams]
liga_id: 1
equipo_id: 2
HTTP 303
[Asserts]
header "Location" contains "/"

# [GET] Consultar la tabla de posiciones de la Liga 1 en el HTML
GET http://localhost:8080/?liga_id=1
HTTP 200
[Asserts]
body contains "<strong>Boca Juniors</strong>"
body contains "<strong>River Plate</strong>"


# 5. GESTIÓN DE PARTIDOS Y CÁLCULO DE PUNTOS EN TABLA

# [POST] Registrar un partido finalizado (Boca Juniors 2 - 0 River Plate)
POST http://localhost:8080/partidos
[FormParams]
liga_id: 1
equipo_local_id: 1
equipo_visitante_id: 2
goles_local: 2
goles_visitante: 0
estado: finalizado
HTTP 303
[Asserts]
header "Location" contains "/"

# [GET] Verificar que el partido y su estado aparecen en la página principal
GET http://localhost:8080/?liga_id=1
HTTP 200
[Asserts]
body contains "2 - 0"
body contains "<mark>finalizado</mark>"


# 6. ELIMINACIÓN VIA FORMULARIOS POST Y VERIFICACIÓN HTML

# 1. Eliminar el partido creado
POST http://localhost:8080/partidos/eliminar
[FormParams]
id: 1
HTTP 303
[Asserts]
header "Location" contains "/"

# 2. Sacar a los equipos de la tabla de posiciones
POST http://localhost:8080/tabla/eliminar
[FormParams]
liga_id: 1
equipo_id: 1
HTTP 303

POST http://localhost:8080/tabla/eliminar
[FormParams]
liga_id: 1
equipo_id: 2
HTTP 303

# 3. Eliminar los equipos
POST http://localhost:8080/equipos/eliminar
[FormParams]
id: 1
HTTP 303

POST http://localhost:8080/equipos/eliminar
[FormParams]
id: 2
HTTP 303

# 4. Eliminar la liga
POST http://localhost:8080/ligas/eliminar
[FormParams]
id: 1
HTTP 303

# 5. Comprobar en el HTML principal que las tablas volvieron a quedar vacías
GET http://localhost:8080/
HTTP 200
[Asserts]
body not contains "<strong>Liga Profesional</strong>"
body not contains "<strong>Boca Juniors</strong>"
body not contains "<strong>River Plate</strong>"
body contains "No hay ligas registradas."
body contains "No hay equipos registrados."
body contains "No hay partidos registrados."
