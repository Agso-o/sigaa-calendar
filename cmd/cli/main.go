package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/Agso-o/sigaa-calendar/internal/wrapper"
	"github.com/joho/godotenv"
)

func main() {
	// Flags de Terminal: Interface CLI
	flag.Usage = func() {
		fmt.Printf("----------Sigaa Calendar----------\n\n")
		fmt.Printf("Uso:\n")
		fmt.Printf(" sigaa-calendar [flags]\n\n")
		fmt.Println("Flags disponíveis:")
		flag.PrintDefaults()
		fmt.Println("Padrão: Sincroniza somente tarefas caso nenhuma flag seja fornecida")
	}
	flagAulas := flag.Bool("aulas", false, "Sincroniza os horários das aulas com o Google Agenda")
	flagTarefas := flag.Bool("tarefas", false, "Sincroniza prazos de entrega de trabalhos com o Google Tasks")
	
	flag.Parse()

	if !*flagTarefas && !*flagAulas {
		*flagTarefas = true
	}
	
	// Recebe Credenciais das variaveis de ambiente
	_ = godotenv.Load()
	sigaaUser := os.Getenv("SIGAA_USER")
	sigaaSenha := os.Getenv("SIGAA_SENHA")
	if sigaaSenha == "" || sigaaUser == "" {
		log.Fatal("Não foi possível recuperar credenciais sigaa")
	}
	err := wrapper.SigaaSync(*flagAulas, *flagTarefas, sigaaUser, sigaaSenha, nil)
	if err != nil {
		log.Fatal("Não foi possível sincronizar calendários: ", err)
	}
	
}
