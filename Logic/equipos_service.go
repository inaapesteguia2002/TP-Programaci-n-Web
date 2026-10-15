package Logic

import (
	"DeRabona/DataBase"
	"context"
	"errors"
)

// EquiposService maneja la lógica de negocio para los equipos
type EquiposService struct {
	Repo *DataBase.Queries
}

// CreateEquipo valida y guarda un nuevo equipo
func (s *EquiposService) CreateEquipo(ctx context.Context, nombre string) (DataBase.Equipo, error) {
	if nombre == "" {
		return DataBase.Equipo{}, errors.New("el nombre del equipo no puede estar vacío")
	}
	return s.Repo.CreateEquipo(ctx, nombre)
}

// ListEquipos devuelve todos los equipos
func (s *EquiposService) ListEquipos(ctx context.Context) ([]DataBase.Equipo, error) {
	return s.Repo.ListEquipos(ctx)
}

// GetEquipo obtiene un equipo por su ID
func (s *EquiposService) GetEquipo(ctx context.Context, id int32) (DataBase.Equipo, error) {
	return s.Repo.GetEquipo(ctx, id)
}

// UpdateEquipo actualiza los datos de un equipo existente
func (s *EquiposService) UpdateEquipo(ctx context.Context, arg DataBase.UpdateEquipoParams) error {
	if arg.Nombre == "" {
		return errors.New("el nuevo nombre no puede estar vacío")
	}
	return s.Repo.UpdateEquipo(ctx, arg)
}

// DeleteEquipo elimina un equipo por su ID
func (s *EquiposService) DeleteEquipo(ctx context.Context, id int32) error {
	return s.Repo.DeleteEquipo(ctx, id)
}
