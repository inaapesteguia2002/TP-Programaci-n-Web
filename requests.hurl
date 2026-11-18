# 1. CARGA INICIAL DE LA PÁGINA PRINCIPAL (SSR HTML + HTMX)

GET http://localhost:8080/
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "De Rabona"
body contains "htmx.min.js"
body contains "hx-post=\"/ligas\""
body contains "No hay ligas registradas."
body contains "No hay equipos registrados."


# 2. GESTIÓN DE LIGAS (CREACIÓN SIN RECARGA -> FRAGMENTO HTML + OOB)

# [ERROR] Intentar crear liga con datos vacíos
POST http://localhost:8080/ligas
[FormParams]
nombre: Liga Incompleta
pais: 
HTTP 400

# [POST] Crear Liga 1 y verificar que devuelve solo el fragmento HTML (200 OK)
POST http://localhost:8080/ligas
[FormParams]
nombre: Liga Profesional
pais: Argentina
HTTP 200
[Asserts]
body contains "id=\"lista-ligas\""
body contains "<strong>Liga Profesional</strong>"
body contains "<td>Argentina</td>"
body contains "hx-delete=\"/ligas/1\""
body contains "hx-swap-oob=\"true\""
body not contains "<!doctype html>"


# 3. GESTIÓN DE EQUIPOS (CREACIÓN SIN RECARGA -> FRAGMENTO HTML + OOB)

# [ERROR] Intentar crear equipo con nombre vacío
POST http://localhost:8080/equipos
[FormParams]
nombre: 
HTTP 400

# [POST] Crear Equipo 1 (Boca Juniors) y verificar fragmento de lista
POST http://localhost:8080/equipos
[FormParams]
nombre: Boca Juniors
HTTP 200
[Asserts]
body contains "id=\"lista-equipos\""
body contains "<strong>Boca Juniors</strong>"
body contains "hx-delete=\"/equipos/1\""
body contains "hx-swap-oob=\"true\""
body not contains "<!doctype html>"

# [POST] Crear Equipo 2 (River Plate) y verificar que ambos están en el fragmento
POST http://localhost:8080/equipos
[FormParams]
nombre: River Plate
HTTP 200
[Asserts]
body contains "id=\"lista-equipos\""
body contains "<strong>Boca Juniors</strong>"
body contains "<strong>River Plate</strong>"
body contains "hx-delete=\"/equipos/2\""


# 4. INSCRIPCIÓN EN TABLA DE POSICIONES (FRAGMENTO HTML)

# [ERROR] Consultar tabla con liga_id inválido en query param
GET http://localhost:8080/?liga_id=invalido
HTTP 400

# [POST] Inscribir Equipo 1 en la Liga 1 (devuelve fragmento #lista-tabla actualizado)
POST http://localhost:8080/tabla
[FormParams]
liga_id: 1
equipo_id: 1
HTTP 200
[Asserts]
body contains "id=\"lista-tabla\""
body contains "<strong>Boca Juniors</strong>"
body not contains "<!doctype html>"

# [POST] Inscribir Equipo 2 en la Liga 1 (devuelve fragmento con ambos equipos)
POST http://localhost:8080/tabla
[FormParams]
liga_id: 1
equipo_id: 2
HTTP 200
[Asserts]
body contains "id=\"lista-tabla\""
body contains "<strong>Boca Juniors</strong>"
body contains "<strong>River Plate</strong>"


# 5. GESTIÓN DE PARTIDOS Y CONSULTA PARCIAL DE TABLA CON HTMX

# [POST] Registrar un partido finalizado (Boca Juniors 2 - 0 River Plate)
POST http://localhost:8080/partidos
[FormParams]
liga_id: 1
equipo_local_id: 1
equipo_visitante_id: 2
goles_local: 2
goles_visitante: 0
estado: finalizado
HTTP 200
[Asserts]
body contains "id=\"lista-partidos\""
body contains "<strong>Boca Juniors</strong>"
body contains "2 - 0"
body contains "<strong>River Plate</strong>"
body contains "<mark>finalizado</mark>"
body contains "hx-delete=\"/partidos/1\""
body not contains "<!doctype html>"

# [GET HTMX] Consultar la tabla de la Liga 1 enviando HX-Request: true (verifica fragmento y 3 pts)
GET http://localhost:8080/?liga_id=1
HX-Request: true
HTTP 200
[Asserts]
body contains "id=\"lista-tabla\""
body contains "<strong>Boca Juniors</strong>"
body contains "<span class=\"pts-pill\">3</span>"
body not contains "<!doctype html>"


# 6. ELIMINACIÓN SIN RECARGA 

# 1. Eliminar el partido 1
DELETE http://localhost:8080/partidos/1
HTTP 200
[Asserts]
body == ""

# 2. Sacar a los equipos de la tabla de posiciones
DELETE http://localhost:8080/tabla/1/1
HTTP 200
[Asserts]
body == ""

DELETE http://localhost:8080/tabla/1/2
HTTP 200
[Asserts]
body == ""

# 3. Eliminar los equipos
DELETE http://localhost:8080/equipos/1
HTTP 200
[Asserts]
body == ""

DELETE http://localhost:8080/equipos/2
HTTP 200
[Asserts]
body == ""

# 4. Eliminar la liga
DELETE http://localhost:8080/ligas/1
HTTP 200
[Asserts]
body == ""

# 5. Comprobar en el HTML principal que todas las tablas volvieron a quedar vacías
GET http://localhost:8080/
HTTP 200
[Asserts]
body not contains "<strong>Liga Profesional</strong>"
body not contains "<strong>Boca Juniors</strong>"
body not contains "<strong>River Plate</strong>"
body contains "No hay ligas registradas."
body contains "No hay equipos registrados."
body contains "No hay partidos registrados."
