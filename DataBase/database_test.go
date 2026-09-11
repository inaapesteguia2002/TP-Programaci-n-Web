package DataBase

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // importa el driver de Postgres
)

func TestIntegracionBD(t *testing.T) {
	ctx := context.Background()

	// 1.Configuramos la conexion a la base de datos local (la que levantará el Makefile)
	connStr := "postgres://postgres:tu_contraseña@localhost:5432/deRabona_db?sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al abrir la conexión: %v", err)
	}
	defer db.Close()

	// Verificamos que la base de datos responde
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("La base de datos no está respondiendo: %v", err)
	}

	queries := New(db)

	// 2.Crear una Liga
	liga, err := queries.CreateLiga(ctx, CreateLigaParams{
		Nombre: "Liga Profesional",
		Pais:   "Argentina",
	})
	if err != nil {
		t.Fatalf("Fallo al crear la liga: %v", err)
	}
	t.Log("Liga creada con éxito. ID:", liga.ID)

	// 3.Crear un Equipo
	equipo, err := queries.CreateEquipo(ctx, "Sarmiento")
	if err != nil {
		t.Fatalf("Fallo al crear el equipo: %v", err)
	}
	t.Log("Equipo creado con éxito. ID:", equipo.ID)

	// 4.Agregar equipo a la tabla de posiciones
	tabla, err := queries.CreateTabla(ctx, CreateTablaParams{
		LigaID:   liga.ID,
		EquipoID: equipo.ID,
	})
	if err != nil {
		t.Fatalf("Fallo al agregar a la tabla de posiciones: %v", err)
	}
	t.Logf("Equipo ingresado a la tabla de posiciones con 0 puntos.", tabla.Puntos)
}
