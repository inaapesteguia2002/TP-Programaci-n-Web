package main

import (
	"context"      // Manejo de tiempos límite (timeouts) y cancelación de peticiones// Manejo de tiempos límite (timeouts) y cancelación de peticiones
	"database/sql" // Interfaz genérica para interactuar con bases de datos relacionales
	"fmt"          // Entrada y salida de texto (impresión en consola)
	"log"          // Registro de eventos y manejo de errores críticos (Fatalf)
	"net/http"     // Creación del servidor web y enrutamiento de peticiones HTTP
	"os"           // Interacción con el sistema (lectura de variables de entorno)
	"time"         // Manejo de duraciones y fechas (los 5 segundos del timeout)

	_ "github.com/jackc/pgx/v5/stdlib" // Driver de PostgreSQL (se registra internamente)

	"DeRabona/DataBase" // Repositorio de consultas generado por sqlc
	"DeRabona/Logic"
)

func main() {

	connSrt := os.Getenv("DATABASE_URL")
	if connSrt == "" {
		connSrt = "postgres://postgres:tu_contraseña@localhost:5432/deRabona_db?sslmode=disable"
	}

	db, err := sql.Open("pgx", connSrt)
	if err != nil {
		log.Fatalf("Error al abrir la conexion: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Ls base de datos no esta respodiendo: %v", err)
	}
	fmt.Println("Conexion a la base de datos exitosa.")

	repo := DataBase.New(db)

	app := &Logic.APIHandler{Repo: repo}

	http.HandleFunc("/equipos", app.EquiposHandler)
	http.HandleFunc("/equipos/", app.EquipoHandler)
	http.HandleFunc("/ligas", app.LigasHandler)
	http.HandleFunc("/ligas/", app.LigaHandler)
	http.HandleFunc("/partidos", app.PartidosHandler)
	http.HandleFunc("/partidos/", app.PartidoHandler)
	http.HandleFunc("/tabla", app.TablasHandler)
	http.HandleFunc("/tabla/", app.TablaLigaHandler)
	http.HandleFunc("/tabla-equipo", app.TablaAccionHandler)

	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
