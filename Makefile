.PHONY: test test-api

test:
	# 1. Limpieza de contenedores y volúmenes
	docker compose down -v --remove-orphans || true

	# 2. Generación de código Go con sqlc
	sqlc generate

	# 3. Levantado contenedor de PostgreSQL
	docker compose up -d postgresql

	# 4. Espera activa a que PostgreSQL acepte conexiones
	@until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done

	# 5. Ejecución de tests de Go
	go test -v ./...

	# 6. Limpieza posterior
	docker compose down -v

test-api:
	@echo "--- 1) Tareas previas ---"
	# 1. Limpieza previa de contenedores y volúmenes huérfanos
	docker rm -f postgres_container 2>/dev/null || true
	docker compose down -v --remove-orphans || true

	# 2. Generación del código Go a partir de las consultas SQL con sqlc
	sqlc generate
	
	# 3. Levantamiento EXCLUSIVO de la base de datos primero
	docker compose up -d postgresql
	
	# 4. Espera activa hasta que PostgreSQL acepte conexiones reales
	@echo "Esperando a que PostgreSQL acepte conexiones..."
	@until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done
	
	# 5. Ahora sí, levantamos la API sabiendo que la base está lista
	@echo "Levantando la API de Go..."
	docker compose up -d --build api
	
	# 6. Pequeña espera para asegurar que el servidor HTTP esté escuchando
	@sleep 10

	@echo "--- 2) Prueba de los endpoints ---"
	# 7. Ejecución de las pruebas funcionales de todos los endpoints utilizando Hurl
	hurl --test requests.hurl

	@echo "--- 3) Tareas posteriores (Limpieza) ---"
	# 8. Limpieza final de contenedores y volúmenes de Docker
	docker compose down -v