package agenda

import (
	"context"
	"fmt"
	_"embed"
	"strings"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

type CalendarService struct {
	srv        *calendar.Service
	CalendarID string
}

type TasksService struct {
	srv 	*tasks.Service
	TaskID	string
}

const (
	NomeAgendaCalendar 	string = "Horário - UFPI"
	NomeListaTasks 		string = "Atividades - UFPI"
	AgendaPrincipal		string = "primary"
	ListaPrincipal		string = "@default"
	ServicoTasks		string = "Tasks"
	ServicoCalendar		string = "Calendar"
)

//go:embed credentials.json
var credenciais []byte
// Retorna um ponteiro para um serviço do google Tasks
// O serviço inclui o acesso à API e o ID da lista de tarefas
func NewTasksService() (*TasksService, error){
	ctx := context.Background()
	// Antiga forma de usar o sistema via leitura do arquivo
	/*credenciais, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler credenciais: %v", err)
	}*/
	config, err := google.ConfigFromJSON(credenciais, tasks.TasksScope)
	if err != nil {
		return nil, fmt.Errorf("Erro ao configurar o cliente: %v", err)
	}

	config.RedirectURL = "http://localhost:8080"
	client, err := getClient(config, ServicoTasks) 
	if err != nil {
		return nil, fmt.Errorf("Erro ao resgatar o cliente: %v", err)
	}
	
	srv, err := tasks.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar o novo serviço: %v", err)
	}
	taskID := ListaPrincipal

	listas, err := srv.Tasklists.List().Do()
	if err == nil {
		for _, lista := range listas.Items {
			if strings.Contains(lista.Title, NomeListaTasks) {
				taskID = lista.Id
				break
			}
		}
	}

	if taskID == ListaPrincipal {
		novaLista := &tasks.TaskList{
			Title: NomeListaTasks,
		}
		
		list, err := srv.Tasklists.Insert(novaLista).Do()
		if err != nil {
			return nil, fmt.Errorf("Não foi possível criar nova lista: %v", err)
		}
		taskID = list.Id
	}

	return &TasksService{
		srv: srv,
		TaskID: taskID,
	}, nil	
}

// Retorna um ponteiro para um serviço do google calendar
// O serviço inclui o acesso a API e o ID da agenda
func NewCalendarService() (*CalendarService, error) { 
	ctx := context.Background()
	// Maneira antiga de conseguir as credenciais via arquivo
	/*credenciais, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler credenciais: %v", err)
	}*/
	// Passa as credenciais e o escopo de permissões e recebe a config pro client
	config, err := google.ConfigFromJSON(credenciais, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("Erro ao configurar cliente: %v", err)
	}
	client, err := getClient(config, ServicoCalendar)
	if err != nil {
		return nil, fmt.Errorf("Erro ao resgatar cliente: %v", err)
	}

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return nil, fmt.Errorf("Erro ao criar o novo serviço: %v", err)
	}
	// Pega o calendar ID
	calendarID := AgendaPrincipal
	// Procura pela agenda de horários nas agendas do usuário
	agendas, err := srv.CalendarList.List().Do()
	if err == nil {
		for _, userAgenda := range agendas.Items {
			if strings.Contains(userAgenda.Summary, NomeAgendaCalendar){
				calendarID = userAgenda.Id
				break
			}
		}
	}
	// Se não encontrar a agenda, cria uma nova agenda
	if calendarID == AgendaPrincipal {
		novoCalendario := &calendar.Calendar{
			Summary: NomeAgendaCalendar,
			Description: "Horário de aulas da UFPI\n By Sigaa-Calendar",
			TimeZone: models.TimeZone, 
		}

		cal, err := srv.Calendars.Insert(novoCalendario).Do()
		if err != nil {
			return nil, fmt.Errorf("Não foi possível criar uma nova agenda: %v", err)
		}
		calendarID = cal.Id
	}
	return &CalendarService{
		srv: srv,
		CalendarID: calendarID,
	}, nil
}
