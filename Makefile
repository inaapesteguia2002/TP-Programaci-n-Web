.PHONY: test

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
	