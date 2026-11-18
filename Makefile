.PHONY: test test-api run-app

test:
	@set -e
	# 1. Registro de limpieza automática desde el inicio 
	trap 'echo "\nLimpiando contenedores..."; docker compose down -v --remove-orphans || true' EXIT INT TERM

	# 2. Limpieza inicial de contenedores y volúmenes previos
	echo "1. Limpieza inicial de contenedores y volúmenes..."
	docker compose down -v --remove-orphans || true

	# 3. Generación de código Go con sqlc y componentes de vista con templ
	echo "2. Generando código con sqlc y templ..."
	sqlc generate
	templ generate

	# 4. Levantado del contenedor de PostgreSQL
	echo "3. Levantando contenedor de PostgreSQL..."
	docker compose up -d postgresql

	# 5. Espera activa a que PostgreSQL acepte conexiones
	echo "4. Esperando a que PostgreSQL acepte conexiones..."
	until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do
		sleep 1
	done

	# 6. Ejecución de tests de Go (al terminar, el trap ejecuta la limpieza final automáticamente)
	echo "5. Ejecutando tests de Go..."
	go test -v ./...

test-api:
	@set -e
	# 1. Registro de limpieza automática desde el inicio (garantiza limpieza ante cualquier fallo)
	trap 'echo "\nLimpiando contenedores..."; docker compose down -v --remove-orphans || true' EXIT INT TERM

	# 2. Limpieza previa de contenedores y volúmenes
	echo "1. Limpieza previa de contenedores y volúmenes..."
	docker rm -f postgres_container 2>/dev/null || true
	docker compose down -v --remove-orphans || true

	# 3. Generación del código Go a partir de las consultas SQL (sqlc) y plantillas HTML (templ)
	echo "2. Generando código con sqlc y templ..."
	sqlc generate
	templ generate

	# 4. Levantamiento EXCLUSIVO de la base de datos primero
	echo "3. Levantando base de datos PostgreSQL..."
	docker compose up -d postgresql

	# 5. Espera activa hasta que PostgreSQL acepte conexiones reales
	echo "4. Esperando a que PostgreSQL acepte conexiones..."
	until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do
		sleep 1
	done

	# 6. Ahora sí, levantamos el servidor web de Go sabiendo que la base está lista
	echo "5. Levantando el servidor web en Go..."
	docker compose up -d --build api

	# 7. Espera inteligente para asegurar que el servidor HTTP esté escuchando (Máx 30s)
	echo "6. Esperando a que el servidor HTTP responda en http://localhost:8080..."
	ready=0
	for i in $$(seq 1 15); do
		if curl -s http://localhost:8080/ > /dev/null; then
			echo "¡Servidor conectado y listo!"
			ready=1
			break
		fi
		sleep 2
	done
	if [ "$$ready" -ne 1 ]; then
		echo "Error: El servidor no respondió a tiempo."
		exit 1
	fi

	# 8. Ejecución de pruebas del recorrido HTML y formularios con Hurl (el trap limpia al finalizar)
	echo "7. Ejecutando suite Hurl (recorrido HTML y formularios)..."
	hurl --test requests.hurl

run-app:
	@set -e
	# 1. Registro de limpieza automática desde el inicio (se activa al apretar Ctrl+C o si falla algo antes)
	trap 'echo "\nLimpiando contenedores..."; docker compose down -v --remove-orphans || true' EXIT INT TERM

	# 2. Tareas previas: limpieza, generación de código (sqlc + templ), levantar base de datos y esperar conexión
	echo "1. Preparando entorno (limpieza, sqlc, templ y PostgreSQL)..."
	docker compose down -v --remove-orphans || true
	sqlc generate
	templ generate
	docker compose up -d postgresql
	echo "Esperando a que PostgreSQL acepte conexiones..."
	until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do
		sleep 1
	done

	# 3. Ejecutar la aplicación (al apretar Ctrl+C se dispara el trap del paso 1 y limpia los contenedores)
	echo "2. Arrancando la aplicación en http://localhost:8080/ (Ctrl+C para salir y limpiar)"
	go run main.go
