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
	docker compose down -v
	sqlc generate
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1; do \
		sleep 2; \
	done
	go build -o main_app main.go
	./main_app & echo $$! > app.pid; \
	sleep 2; \
	USER_RESP=$$(curl -s -X POST http://localhost:8080/usuarios -H "Content-Type: application/json" -d '{"nombre": "Justina Test", "email": "justina@example.com"}'); \
	USER_ID=$$(echo $$USER_RESP | grep -o '"id":[0-9]*' | cut -d':' -f2); \
	curl -s http://localhost:8080/usuarios; \
	PRENDA_RESP=$$(curl -s -X POST http://localhost:8080/prendas -H "Content-Type: application/json" -d "{\"usuario_id\": $$USER_ID, \"nombre\": \"Campera Jean\", \"categoria\": \"Abrigos\"}"); \
	PRENDA_ID=$$(echo $$PRENDA_RESP | grep -o '"id":[0-9]*' | cut -d':' -f2); \
	LOOK_RESP=$$(curl -s -X POST http://localhost:8080/looks -H "Content-Type: application/json" -d "{\"usuario_id\": $$USER_ID, \"nombre\": \"Look Urbano\", \"descripcion\": \"Casual\", \"ocasion\": \"Dia\", \"temporada\": \"Primavera\", \"prenda_estrella_id\": $$PRENDA_ID, \"url_imagen\": \"http://example.com/img.jpg\", \"estilo\": \"Casual\"}"); \
	LOOK_ID=$$(echo $$LOOK_RESP | grep -o '"id":[0-9]*' | cut -d':' -f2); \
	curl -s http://localhost:8080/looks; \
	curl -s -i -X POST http://localhost:8080/me-gusta -H "Content-Type: application/json" -d "{\"usuario_id\": $$USER_ID, \"look_id\": $$LOOK_ID}"; \
	curl -s -i -X DELETE http://localhost:8080/me-gusta -H "Content-Type: application/json" -d "{\"usuario_id\": $$USER_ID, \"look_id\": $$LOOK_ID}"; \
	curl -s -i -X DELETE http://localhost:8080/looks/$$LOOK_ID; \
	kill $$(cat app.pid) || true; \
	rm -f main_app app.pid; \
	docker compose down -v