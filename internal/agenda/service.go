package agenda

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"golang.org/x/oauth2"
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

type GoogleService struct {
	CalendarService		*CalendarService
	TasksService 		*TasksService
}

const (
	NomeAgendaCalendar 	string = "Horário - UFPI"
	NomeListaTasks 		string = "Atividades - UFPI"
	AgendaPrincipal		string = "primary"
	ListaPrincipal		string = "@default"
	ServicoTasks		string = "Tasks"
	ServicoCalendar		string = "Calendar"
)

func NewGoogleService(client *http.Client) (*GoogleService, error){
	// Confguração principal pra pegar permissão
	ctx := context.Background()
	if client == nil {
		return nil, fmt.Errorf("Nenhum cliente do Google encontrado")	
	}
	
	//Criação do serviço tasks e criação da lista se não existir
	srvTasks, err := tasks.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar o serviço do tasks: %v", err)
	}
	taskID := ListaPrincipal

	listas, err := srvTasks.Tasklists.List().Do()
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
		
		list, err := srvTasks.Tasklists.Insert(novaLista).Do()
		if err != nil {
			return nil, fmt.Errorf("Não foi possível criar nova lista: %v", err)
		}
		taskID = list.Id
	}
	
	// Cria o serviço da calendar e cria a agenda se não existir
	srvCalendar, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Erro ao criar o novo serviço: %v", err)
	}
	// Pega o calendar ID
	calendarID := AgendaPrincipal
	// Procura pela agenda de horários nas agendas do usuário
	agendas, err := srvCalendar.CalendarList.List().Do()
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

		cal, err := srvCalendar.Calendars.Insert(novoCalendario).Do() 
		if err != nil {
			return nil, fmt.Errorf("Não foi possível criar uma nova agenda: %v", err)
		}
		calendarID = cal.Id
	}

	TarefasServico := &TasksService{
		srv: srvTasks,
		TaskID: taskID,
	}
	CalendarioServico := &CalendarService{
		srv: srvCalendar,
		CalendarID: calendarID,
	}
	return &GoogleService{
		TasksService: TarefasServico,
		CalendarService: CalendarioServico,
	}, nil

}
// Recebe as credenciais do Google Cloud e devolve um client
func GetClientFromCredentials(clientID string, clientSecret string) (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080", // A rota de retorno
		Scopes: []string{
			tasks.TasksScope,
			calendar.CalendarScope,
		},
		Endpoint: google.Endpoint,
	}
	client, err := getClient(config, fmt.Sprintf("%s+%s",ServicoTasks, ServicoCalendar)) 
	if err != nil {
		return nil, err
	}
	return client, nil
}
