package agenda

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"google.golang.org/api/tasks/v1"
)

var reSigaaID = regexp.MustCompile(`<!--sigaa-id:([a-f0-9]+)-->`)

func normalizarTexto(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

func gerarSigaaID(disciplina, titulo string) string {
	hash := md5.Sum([]byte(normalizarTexto(disciplina) + "|" + normalizarTexto(titulo)))
	return hex.EncodeToString(hash[:])
}

// Cria uma task e salva no google Tasks
func (s *TasksService) SaveTask(tarefa models.Tarefa, tarefasExistentes map[string]string) (error) {
	chave := tarefa.SigaaID
	if chave == "" {
		chave = gerarSigaaID(tarefa.Disciplina, tarefa.Titulo)
	}

	task := &tasks.Task{
		Title: fmt.Sprintf("[%s]: %s", tarefa.Disciplina, tarefa.Titulo),
		Notes: fmt.Sprintf("%s\n%s\n\n<!--sigaa-id:%s-->",
			tarefa.DataVencimento.Format("02/01/2006 às 15:04"),
			tarefa.Descricao,
			chave,
		),
		Due: tarefa.DataVencimento.Format("2006-01-02") + "T00:00:00Z",
	}

	if id, existe := tarefasExistentes[chave]; existe {
		existente, err := s.srv.Tasks.Get(s.TaskID, id).Do()
		if err != nil {
			return fmt.Errorf("Erro ao buscar tarefa existente: %v", err)
		}
		if existente.Due == task.Due && existente.Notes == task.Notes {
			return nil
		}
		
		var dueMudou bool
		ta, errA := time.Parse(time.RFC3339, existente.Due)
		tb, errB := time.Parse(time.RFC3339, task.Due)
		if errA != nil || errB != nil {
			dueMudou = (existente.Due == task.Due)
		}
		dueMudou = ta.Equal(tb)
		notesMudaram := existente.Notes != task.Notes

		if !dueMudou && !notesMudaram {
			return nil
		}
		task.Id = id
		task.Status = existente.Status

		_, err = s.srv.Tasks.Patch(s.TaskID, id, task).Do()
		if err != nil {
			return fmt.Errorf("Erro ao atualizar tarefa: %v", err)
		}
		return nil
	}

	inserted, err := s.srv.Tasks.Insert(s.TaskID, task).Do()
	if err != nil {
		return fmt.Errorf("Erro ao salvar tarefa: %v", err)
	}

	tarefasExistentes[chave] = inserted.Id
	return nil
}

// Recupera as tasks do Google Tasks e retorna em um mapa
// Somente retorna as tasks da lista usada para salvar as tarefas do sigaa
func (s *TasksService) GetTasks() (map[string]string, error){
	mapaTarefas := make(map[string]string)
    pageToken := ""

    for {
        tarefas, err := s.srv.Tasks.List(s.TaskID).ShowCompleted(true).ShowHidden(true).PageToken(pageToken).Do()
        if err != nil {
            return nil, fmt.Errorf("erro ao listar tasks: %v", err)
        }

        for _, tarefa := range tarefas.Items {
			chave := tarefa.Title
			if match := reSigaaID.FindStringSubmatch(tarefa.Notes); len(match) > 1 {
				chave = match[1]
			}
            mapaTarefas[chave] = tarefa.Id
        }
        if tarefas.NextPageToken == "" {
            break
        }
        pageToken = tarefas.NextPageToken
    }

    return mapaTarefas, nil

}

