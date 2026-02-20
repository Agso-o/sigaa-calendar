package agenda

import (
	"fmt"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"google.golang.org/api/tasks/v1"
)
// Cria uma task e salva no google Tasks
func (s *TasksService) SaveTask(tarefa models.Tarefa, tarefasExistentes map[string]string) (error) {
	task := &tasks.Task{
		Title: fmt.Sprintf("[%s]: %s", tarefa.Disciplina, tarefa.Titulo),
		Notes: fmt.Sprintf("%s\n%s", tarefa.DataVencimento.Format("02/01/2006 às 15:04"), tarefa.Descricao),
		Due: tarefa.DataVencimento.Format("2006-01-02")+"T00:00:00Z",
		
	}
	// Aqui tem que implementar a verificação se mudou a data de entrega
	if _, existe := tarefasExistentes[task.Title]; existe {
		return nil
	}
	_, err := s.srv.Tasks.Insert(s.TaskID, task).Do()
	if err != nil {
		return fmt.Errorf("Erro ao salvar evento: %v", err)
	}

	tarefasExistentes[task.Title] = "ID"
	return nil
}

// Recupera as tasks do Google Tasks e retorna em um mapa
// Somente retorna as tasks da lista usada para salvar as tarefas do sigaa
func (s *TasksService) GetTasks() (map[string]string, error){
	mapaTarefas := make(map[string]string)
    pageToken := "" // Começa vazio

    for {
		// Como o sistema já verifica se a tarefa foi completada no sigaa, nao precisa pegar as completadas aqui
        tarefas, err := s.srv.Tasks.List(s.TaskID).ShowCompleted(false).ShowHidden(true).PageToken(pageToken).Do()
        if err != nil {
            return nil, fmt.Errorf("erro ao listar tasks: %v", err)
        }

        for _, tarefa := range tarefas.Items {
            mapaTarefas[tarefa.Title] = tarefa.Id
        }
        if tarefas.NextPageToken == "" {
            break
        }
        pageToken = tarefas.NextPageToken
    }

    return mapaTarefas, nil

}

