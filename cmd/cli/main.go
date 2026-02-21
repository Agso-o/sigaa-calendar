package main

import (
	"fmt"
	"log"
	"os"
	"github.com/Agso-o/sigaa-calendar/internal/agenda"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"github.com/Agso-o/sigaa-calendar/internal/sigaa"
	"github.com/joho/godotenv"
)

func main() {
	// Recebe Credenciais das variaveis de ambiente
	_ = godotenv.Load()
	sigaaUser := os.Getenv("SIGAA_USER")
	sigaaSenha := os.Getenv("SIGAA_SENHA")
	if sigaaSenha == "" || sigaaUser == "" {
		log.Fatal("Não foi possível recuperar credenciais sigaa")
	}

	// Loga no sigaa e recupera as turmas
	sigaaSrv, err:= sigaa.NewSigaaService(sigaaUser, sigaaSenha)
	if err != nil{
		log.Fatal("Não foi possíel inicira serviço sigaa: ", err)
	}
	err = sigaaSrv.LoginSigaa()
	if err != nil {
		log.Fatal("Erro ao logar no sigaa")
	}
	turmas, err := sigaaSrv.GetTurmasAnteriores()
	if err != nil {
		log.Fatal("Não foi possível recuperar turmas anteriores: ", err)
	}
	
	// Cria o serviço do google
	agendaSrv, err := agenda.NewGoogleService()
	if err != nil {
		log.Fatal("Não foi possível iniciar o serviço ", err)
	}
	
	// Pega as tarefas da conta do tasks e guarda em um mapa pra não salvar repetida
	tarefasExistentes, err := agendaSrv.TasksService.GetTasks()
	if err != nil {
		log.Fatal("Não foi possível listar tarefas: ", err)
	}

	for _, turma := range turmas {
		fmt.Printf("Turma: %s [%s]\n", turma.Disciplina, turma.AnoSemestre)
		rec, st, end, _ := sigaa.ParseHorarioAgenda(turma.Horario, turma.AnoSemestre)
		fmt.Printf("\tLocal: %s\n\tCréditos: %s\n\tHorário Sigaa: %s\n\tHorario Parsed: \n\t\t%s\n\t\t%v\n\t\t%v\n", 
			turma.Local, turma.Creditos, turma.Horario, rec, st, end)
		
		// Salva o horário no calendar
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
			fmt.Printf("Não foi possível salvar evento para: %s\n\tErro: %v\n",
				novaAula.Disciplina, err)
				return
		} else {
			fmt.Println("----------Evento de aula salvo com sucesso----------")
		}
		
		// Salva as tarefas no Tasks
		// Pega as tarefas da conta do sigaa
		tarefas, err := sigaaSrv.GetTarefasByTurma(turma)
		if err != nil {
			log.Printf("Não foi possível pegar atividades de %s, %v\n", turma.Disciplina, err)
			return
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
