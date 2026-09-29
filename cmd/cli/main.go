package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/Agso-o/sigaa-calendar/internal/agenda"
	"github.com/Agso-o/sigaa-calendar/internal/wrapper"
	"github.com/joho/godotenv"
)

func main() {
	// Flags de Terminal: Interface CLI
	flag.Usage = func() {
helpMessage := `----------Sigaa Calendar----------
Uso:
 sigaa-calendar [flags]

Flags disponíveis:`
	fmt.Println(helpMessage)
	flag.PrintDefaults()
	}
	flagAulas := flag.Bool("aulas", false, "Sincroniza os horários das aulas com o Google Agenda")
	flagTarefas := flag.Bool("tarefas", false, "Sincroniza prazos de entrega de trabalhos com o Google Tasks")
	flagBiblioteca := flag.Bool("lib", false, "Sincroniza prazos de entrega de livros com o Google Tasks")
	
	flag.Parse()
	
	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	_ = godotenv.Load()
	sigaaUser := os.Getenv("SIGAA_USER")
	sigaaSenha := os.Getenv("SIGAA_SENHA")
	if sigaaSenha == "" || sigaaUser == "" {
		log.Fatal("Não foi possível recuperar credenciais sigaa")
	}
	clientID := os.Getenv("GOOGLE_CLIENTID")
	clientSecret := os.Getenv("GOOGLE_CLIENTSECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("Incapaz de carregar as credenciais do Google")
	}
	client, err := agenda.GetClientFromCredentials(clientID, clientSecret)
	
	err = wrapper.SigaaSync(*flagAulas, *flagTarefas, *flagBiblioteca ,sigaaUser, sigaaSenha, client)
	if err != nil {
		log.Fatal("Não foi possível sincronizar calendários: ", err)
	}
}
