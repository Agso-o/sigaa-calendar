package agenda

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
)

// Salva um Evento usando o tipo aula, salvando os horarios de aulas
func (s *CalendarService) SaveEventAula(aula models.Aula) error {
	hash := md5.Sum([]byte(aula.Disciplina+aula.Recorrencia))
	id  := hex.EncodeToString(hash[:]) 
	event := &calendar.Event{
		Id: id,
		Summary: aula.Disciplina, 
		Location: aula.Local,
		Description: "Horario de aula da Universidade Federal do Piauí\nBy sigaa-Calendar",
		Start: &calendar.EventDateTime{
			DateTime: aula.StartTime.Format(time.RFC3339),
			TimeZone: models.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: aula.EndTime.Format(time.RFC3339),
			TimeZone: models.TimeZone,
		},
		Recurrence: []string{aula.Recorrencia},

	}
	
	_, err := s.srv.Events.Insert(s.CalendarID, event).Do()
	if err != nil {
		// Se o evento já foi criado, não vai criar de novo
		if googleError, ok := err.(*googleapi.Error); ok && googleError.Code == 409 {
            return nil
        }
		return fmt.Errorf("Erro ao salvar evento: %v", err)
	}
	return nil
}

// Salva as tarefas do sigaa em formato de evento no google Calendar
// Por preferencia de quem não quer usar o tasks
func (s *CalendarService) SaveEventTarefa(tarefa models.Tarefa) error {
	hash := md5.Sum([]byte(tarefa.Titulo+tarefa.DataVencimento.Format(time.RFC3339)))
	id  := hex.EncodeToString(hash[:]) 
	event := &calendar.Event{
		Id: id,
		Summary: tarefa.Titulo,
		Location: "SIGAA",
		Description: tarefa.Descricao,
		Start: &calendar.EventDateTime{
			DateTime: tarefa.DataVencimento.Add(-time.Hour).Format(time.RFC3339),
			TimeZone: models.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: tarefa.DataVencimento.Format(time.RFC3339),
			TimeZone: models.TimeZone,
		},
	}
	_, err := s.srv.Events.Insert(s.CalendarID, event).Do()
	if err != nil {
		// Se o evento já foi criado, não vai criar de novo
		if googleError, ok := err.(*googleapi.Error); ok && googleError.Code == 409 {
            return nil
        }
		return fmt.Errorf("Erro ao salvar evento: %v", err)
	}

	return nil
}

func (s *CalendarService) SaveEventEmprestimo(emp models.Emprestimo) error {
	hash := md5.Sum([]byte(emp.Livro + emp.Prazo.Format(time.RFC3339)))
	id := hex.EncodeToString(hash[:])
	event := &calendar.Event{
		Id:          id,
		Summary:     "Devolver: " + emp.Livro,
		Location:    "Biblioteca UFPI",
		Description: "Prazo de devolução\nBy sigaa-Calendar",
		Start: &calendar.EventDateTime{
			DateTime: emp.Prazo.Add(-time.Hour).Format(time.RFC3339),
			TimeZone: models.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: emp.Prazo.Format(time.RFC3339),
			TimeZone: models.TimeZone,
		},
	}
	_, err := s.srv.Events.Insert(s.CalendarID, event).Do()
	if err != nil {
		if googleError, ok := err.(*googleapi.Error); ok && googleError.Code == 409 {
			return nil
		}
		return fmt.Errorf("Erro ao salvar evento de empréstimo: %v", err)
	}
	return nil
}
