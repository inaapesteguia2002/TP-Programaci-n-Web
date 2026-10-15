package Logic

import (
	"context"
	"database/sql"
	"errors"

	"DeRabona/DataBase"
)

// PartidosService maneja las reglas de negocio para los partidos
type PartidosService struct {
	Repo *DataBase.Queries
}

// CreatePartido valida y guarda un nuevo partido
func (s *PartidosService) CreatePartido(ctx context.Context, arg DataBase.CreatePartidoParams) (DataBase.Partido, error) {
	// Regla de negocio: Un equipo no puede jugar contra sí mismo
	if arg.EquipoLocalID == arg.EquipoVisitanteID {
		return DataBase.Partido{}, errors.New("un equipo no puede jugar contra sí mismo")
	}
	if arg.LigaID <= 0 {
		return DataBase.Partido{}, errors.New("el ID de la liga es inválido")
	}

	return s.Repo.CreatePartido(ctx, arg)
}

// ListPartidos devuelve todos los partidos
func (s *PartidosService) ListPartidos(ctx context.Context) ([]DataBase.Partido, error) {
	return s.Repo.ListPartidos(ctx)
}

// GetPartido obtiene un partido por su ID
func (s *PartidosService) GetPartido(ctx context.Context, id int32) (DataBase.Partido, error) {
	return s.Repo.GetPartido(ctx, id)
}

// UpdatePartido actualiza los datos de un partido (solo goles y estado)
func (s *PartidosService) UpdatePartido(ctx context.Context, arg DataBase.UpdatePartidoParams) error {
	if !arg.Estado.Valid || arg.Estado.String == "" {
		arg.Estado = sql.NullString{String: "pendiente", Valid: true}
	}
	return s.Repo.UpdatePartido(ctx, arg)
}

// DeletePartido elimina un partido por su ID
func (s *PartidosService) DeletePartido(ctx context.Context, id int32) error {
	return s.Repo.DeletePartido(ctx, id)
}
