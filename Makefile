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
	# 1. Limpieza de contenedores y volúmenes previos
	docker compose down -v --remove-orphans || true

	# 2. Generación de código Go con sqlc
	sqlc generate

	# 3. Levantado contenedor de PostgreSQL
	docker compose up -d postgresql

	# 4. Espera activa a que PostgreSQL acepte conexiones
	@until docker exec postgres_container pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done

	# 5. Compilación y levantamiento de la aplicación en segundo plano
	go build -o tmp/derabona .
	./tmp/derabona & echo $$! > tmp/server.pid
	@sleep 2

	# 6. Prueba de los endpoints utilizando curl[cite: 3]
	@echo "==> Probando endpoints de la API..."
	curl -s -X POST http://localhost:8080/equipos -d '{"nombre":"River Plate"}' -H "Content-Type: application/json" || true
	curl -s -X GET http://localhost:8080/equipos || true
	curl -s -X POST http://localhost:8080/ligas -d '{"nombre":"Liga Profesional", "pais":"Argentina"}' -H "Content-Type: application/json" || true
	curl -s -X GET http://localhost:8080/ligas || true

	# 7. Apagar el servidor de la aplicación en ejecución
	@if [ -f tmp/server.pid ]; then kill $$(cat tmp/server.pid) 2>/dev/null || true; fi

	# 8. Limpieza posterior de contenedores y volúmenes[cite: 3]
	docker compose down -v