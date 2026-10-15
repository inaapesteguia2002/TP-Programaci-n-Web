package Logic

import (
	"context"
	"errors"

	"DeRabona/DataBase"
)

// TablaService maneja las reglas de negocio de la tabla de posiciones
type TablaService struct {
	Repo *DataBase.Queries
}

// RegistrarEquipo inscribe un equipo en una liga
func (s *TablaService) RegistrarEquipo(ctx context.Context, arg DataBase.CreateTablaParams) (DataBase.TablaPosicione, error) {
	if arg.LigaID <= 0 || arg.EquipoID <= 0 {
		return DataBase.TablaPosicione{}, errors.New("el ID de la liga y del equipo deben ser mayores a 0")
	}
	return s.Repo.CreateTabla(ctx, arg)
}

// ObtenerTabla devuelve las posiciones de una liga específica
func (s *TablaService) ObtenerTabla(ctx context.Context, ligaID int32) ([]DataBase.GetTablaPorLigaRow, error) {
	if ligaID <= 0 {
		return nil, errors.New("el ID de la liga es inválido")
	}
	return s.Repo.GetTablaPorLiga(ctx, ligaID)
}

// ActualizarEstadisticas suma puntos y resultados a un equipo
func (s *TablaService) ActualizarEstadisticas(ctx context.Context, arg DataBase.UpdateTablaParams) error {
	if arg.LigaID <= 0 || arg.EquipoID <= 0 {
		return errors.New("liga_id y equipo_id son obligatorios para actualizar")
	}
	return s.Repo.UpdateTabla(ctx, arg)
}

// RemoverEquipo saca a un equipo de la tabla de posiciones
func (s *TablaService) RemoverEquipo(ctx context.Context, arg DataBase.DeleteEquipoDeTablaParams) error {
	return s.Repo.DeleteEquipoDeTabla(ctx, arg)
}
