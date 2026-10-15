package service

import (
	"errors"
	"time"

	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
)

type CitaService struct {
	citaRepo     *repository.CitaRepository
	catalogoRepo *repository.CatalogoRepository
}

func NewCitaService(citaRepo *repository.CitaRepository, catalogoRepo *repository.CatalogoRepository) *CitaService {
	return &CitaService{citaRepo: citaRepo, catalogoRepo: catalogoRepo}
}

func (s *CitaService) AgendarCita(clienteID int, req domain.CrearCitaRequest) (*domain.Cita, error) {
	// 1. Obtener los servicios elegidos y sumar duración y precio total
	servicios, err := s.citaRepo.ObtenerServiciosPorIDs(req.ServiciosIDs)
	if err != nil || len(servicios) == 0 {
		return nil, errors.New("debes seleccionar al menos un servicio válido")
	}

	var duracionTotal int
	var montoTotal float64
	var detalles []domain.CitaServicio

	for _, srv := range servicios {
		duracionTotal += srv.DuracionMinutos
		montoTotal += srv.Precio
		detalles = append(detalles, domain.CitaServicio{
			ServicioID:      srv.ID,
			PrecioAplicado:  srv.Precio,
			DuracionMinutos: srv.DuracionMinutos,
			ComisionBarbero: srv.Precio * 0.50, // 50% del coste del servicio para el barbero
		})
	}

	// 2. Calcular la HoraFin sumando los minutos totales a la HoraInicio
	inicioTime, err := time.Parse("15:04", req.HoraInicio)
	if err != nil {
		return nil, errors.New("formato de hora_inicio inválido (usa HH:MM, ej. 15:00)")
	}
	finTime := inicioTime.Add(time.Duration(duracionTotal) * time.Minute)
	horaInicioStr := inicioTime.Format("15:04:05")
	horaFinStr := finTime.Format("15:04:05")

	// 3. Validar contra el horario de apertura, cierre y descanso general de la barbería
	config, err := s.catalogoRepo.ObtenerConfiguracion()
	if err != nil {
		return nil, errors.New("error al verificar horarios de la barbería")
	}

	if horaInicioStr < config.HoraApertura || horaFinStr > config.HoraCierre {
		return nil, errors.New("el horario solicitado está fuera del horario de atención de la barbería")
	}

	// Verificar si se cruza con la hora de descanso general (ej. 13:00 a 14:30)
	if horaInicioStr < config.HoraFinDescanso && horaFinStr > config.HoraInicioDescanso {
		return nil, errors.New("el horario choca con la hora de descanso de los barberos (13:00 a 14:30)")
	}

	// 4. Modelo Híbrido: Asignar barbero específico o buscar uno libre automáticamente
	barberoElegidoID := req.BarberoID
	barberos, err := s.citaRepo.ObtenerBarberosActivos()
	if err != nil || len(barberos) == 0 {
		return nil, errors.New("no hay barberos disponibles en el sistema")
	}

	if barberoElegidoID > 0 {
		// El cliente eligió un barbero específico: verificamos si está libre
		cruce, err := s.citaRepo.BarberoTieneCruce(barberoElegidoID, req.Fecha, horaInicioStr, horaFinStr)
		if err != nil || cruce {
			return nil, errors.New("el barbero seleccionado ya tiene una cita en ese horario")
		}
	} else {
		// El cliente eligió "Cualquier barbero disponible" (barbero_id = 0)
		for _, b := range barberos {
			cruce, err := s.citaRepo.BarberoTieneCruce(b.ID, req.Fecha, horaInicioStr, horaFinStr)
			if err == nil && !cruce {
				barberoElegidoID = b.ID
				break
			}
		}
		if barberoElegidoID == 0 {
			return nil, errors.New("todos los barberos están ocupados en ese horario, intenta otra hora")
		}
	}

	// 5. Calcular garantía del 50% para asegurar la cita
	montoGarantia := montoTotal * (config.PorcentajeGarantiaReserva / 100.0)

	nuevaCita := &domain.Cita{
		ClienteID:           clienteID,
		BarberoID:           barberoElegidoID,
		Fecha:               req.Fecha,
		HoraInicio:          horaInicioStr,
		HoraFin:             horaFinStr,
		Origen:              "web",
		MontoTotal:          montoTotal,
		MontoGarantiaPagado: montoGarantia,
		MetodoPagoGarantia:  req.MetodoPagoGarantia,
		Estado:              "confirmada",
	}

	if err := s.citaRepo.GuardarCitaCompleta(nuevaCita, detalles); err != nil {
		return nil, errors.New("no se pudo registrar la cita en la base de datos")
	}

	return nuevaCita, nil
}
