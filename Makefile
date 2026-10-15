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
	# 1. Limpieza previa de contenedores y volúmenes huérfanos
	@echo "--- 1) Tareas previas ---"
	docker compose down -v --remove-orphans || true

	# 2. Generación del código Go a partir de las consultas SQL con sqlc
	sqlc generate
	
	# 3. Levantamiento de la base de datos y la API mediante Docker Compose (--build)
	docker compose up -d --build
	
	# 4. Espera activa hasta que PostgreSQL acepte conexiones
	@echo "Esperando a que PostgreSQL acepte conexiones..."
	@until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done
	
	# 5. Espera prudencial para que la API termine de inicializarse
	@echo "Dando unos segundos para que la API termine de conectarse..."
	@sleep 3

	# 6. Ejecución de las pruebas funcionales de todos los endpoints utilizando Hurl
	@echo "--- 2) Prueba de los endpoints ---"
	hurl --test requests.hurl

	# 7. Limpieza final de contenedores y volúmenes de Docker
	@echo "--- 3) Tareas posteriores (Limpieza) ---"
	docker compose down -v