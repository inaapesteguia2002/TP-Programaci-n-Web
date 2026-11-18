.PHONY: test test-api run-app

test:
	@set -e; \
	trap 'echo "\nLimpiando contenedores..."; docker compose down -v --remove-orphans || true' EXIT INT TERM; \
	echo "1. Limpieza inicial de contenedores y volúmenes..."; \
	docker compose down -v --remove-orphans || true; \
	echo "2. Generando código con sqlc y templ..."; \
	sqlc generate; \
	templ generate; \
	echo "3. Levantando contenedor de PostgreSQL..."; \
	docker compose up -d postgresql; \
	echo "4. Esperando a que PostgreSQL acepte conexiones..."; \
	until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done; \
	echo "5. Ejecutando tests de Go..."; \
	go test -v ./...

test-api:
	@set -e; \
	trap 'echo "\nLimpiando contenedores..."; docker compose down -v --remove-orphans || true' EXIT INT TERM; \
	echo "1. Limpieza previa de contenedores y volúmenes..."; \
	docker rm -f postgres_container 2>/dev/null || true; \
	docker compose down -v --remove-orphans || true; \
	echo "2. Generando código con sqlc y templ..."; \
	sqlc generate; \
	templ generate; \
	echo "3. Levantando base de datos PostgreSQL..."; \
	docker compose up -d postgresql; \
	echo "4. Esperando a que PostgreSQL acepte conexiones..."; \
	until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done; \
	echo "5. Levantando el servidor web en Go..."; \
	docker compose up -d --build api; \
	echo "6. Esperando a que el servidor HTTP responda en http://localhost:8080..."; \
	ready=0; \
	for i in $$(seq 1 15); do \
		if curl -s http://localhost:8080/ > /dev/null; then \
			echo "¡Servidor conectado y listo!"; \
			ready=1; \
			break; \
		fi; \
		sleep 2; \
	done; \
	if [ "$$ready" -ne 1 ]; then \
		echo "Error: El servidor no respondió a tiempo."; \
		exit 1; \
	fi; \
	echo "7. Ejecutando suite Hurl (recorrido HTML y formularios)..."; \
	hurl --test requests.hurl

run-app:
	@set -e; \
	trap 'echo "\nLimpiando contenedores..."; docker compose down -v --remove-orphans || true' EXIT INT TERM; \
	echo "1. Preparando entorno (limpieza, sqlc, templ y PostgreSQL)..."; \
	docker compose down -v --remove-orphans || true; \
	sqlc generate; \
	templ generate; \
	docker compose up -d postgresql; \
	echo "Esperando a que PostgreSQL acepte conexiones..."; \
	until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done; \
	echo "2. Arrancando la aplicación en http://localhost:8080/ (Ctrl+C para salir y limpiar)"; \
	go run main.go
