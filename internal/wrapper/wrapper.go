package wrapper

import (
	"fmt"
	"log"
	"net/http"
	"github.com/Agso-o/sigaa-calendar/internal/agenda"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"github.com/Agso-o/sigaa-calendar/internal/sigaa"
)

func SigaaSync (flagAulas, flagTarefas bool, sigaaUser, sigaaPass string, googleClient *http.Client) error {
	// Loga no sigaa e recupera as turmas
	sigaaSrv, err:= sigaa.NewSigaaService(sigaaUser, sigaaPass)
	if err != nil{
		return fmt.Errorf("Não foi possíel iniciar serviço sigaa: %v", err)
	}
	err = sigaaSrv.LoginSigaa()
	if err != nil {
		return fmt.Errorf("Erro ao logar no sigaa")
	}

	turmas, err := sigaaSrv.GetTurmas()

	if err != nil {
		return fmt.Errorf("Não foi possível recuperar turmas: %v", err)
	}
	
	// Cria o serviço do google
	agendaSrv, err := agenda.NewGoogleService(googleClient)

	if err != nil {
		return fmt.Errorf("Não foi possível iniciar o serviço do google: %v", err)
	}
	
	var tarefasExistentes map[string]string

	if flagTarefas {
		// Pega as tarefas da conta do tasks e guarda em um mapa pra não salvar repetida
		tarefasExistentes, err = agendaSrv.TasksService.GetTasks()
		if err != nil {
			return fmt.Errorf("Não foi possível listar tarefas: %v", err)
		}
	}

	for _, turma := range turmas {
		fmt.Printf("Turma: %s [%s]\n", turma.Disciplina, turma.AnoSemestre)
		rec, st, end, _ := sigaa.ParseHorarioAgenda(turma.Horario, turma.AnoSemestre)
		fmt.Printf("\tLocal: %s\n\tCréditos: %s\n\tHorário Sigaa: %s\n\tHorario Parsed: \n\t\t%s\n\t\t%v\n\t\t%v\n", 
			turma.Local, turma.Creditos, turma.Horario, rec, st, end)
		
		// Salva o horário no calendar
		if flagAulas {
			// Cria o objeto de aula
			novaAula := &models.Aula{
				Disciplina: turma.Disciplina,
				Local: turma.Local,
				StartTime: st,
				EndTime: end,
				Recorrencia: rec,
			}
			//Tenta salvar o evento e trata o erro
			err := agendaSrv.CalendarService.SaveEventAula(*novaAula)
			if err != nil {
				log.Printf("Não foi possível salvar evento para: %s\n\tErro: %v\n",
					novaAula.Disciplina, err)
			} else {
				log.Println("----------Evento de aula salvo com sucesso----------")
			}
		}	

		// Salva as tarefas no Tasks
		if flagTarefas {
			// Pega as tarefas da conta do sigaa
			tarefas, err := sigaaSrv.GetTarefasByTurma(turma)
			if err != nil {
				log.Printf("Não foi possível pegar atividades de %s, %v\n", turma.Disciplina, err)
				continue
			}
			log.Println("----------Tarefas Coletadas----------")
			// Tenta salvar as tarefas de cada turma
			for _, tarefa := range tarefas {
				err := agendaSrv.TasksService.SaveTask(tarefa, tarefasExistentes)
				if err != nil {
					log.Println("Não foi possível salvar essa tarefa: ", err)
				}
			}
			log.Println("----------Tarefas salvas----------")
		}
	}

	return nil
}
