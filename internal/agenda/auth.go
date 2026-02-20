package agenda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/oauth2"
)

func getTokenPath(tipoServico string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "sigaa-calendar")

	_ = os.MkdirAll(appDir, 0755)

	return filepath.Join(appDir, fmt.Sprintf("%sToken.json", tipoServico))

}

// Resgata um token, salva o token e retorna o cliente
func getClient(config *oauth2.Config, tipoServico string) (*http.Client, error) {
	tokFile := getTokenPath(tipoServico)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, fmt.Errorf("Erro ao requisitar token da web: %v", err)
		}
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok), nil
}

// Faz a requisição do token na web e retorna o token resgatado
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error){
	canalCodigo := make(chan string)
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	server := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		canalCodigo <- code
		http.Redirect(w, r, "https://calendar.google.com/calendar/", http.StatusSeeOther)
	})

	
	go func() {
		server.ListenAndServe()	
	}()

	err := openBrowser(authURL)
	if err != nil {
		return nil, err
	}
	authCode := <- canalCodigo

	go func() {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()


	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("Erro ao requerir token da web: %v", err)
	}
	return tok, nil
}
// Função auxiliar para abrir o navegador
func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Plataforma não suportada")
	}
	return err
}

// Resgata um token de um arquivo local
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Salva um token a partir de um caminho de arquivo
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Não foi possível salvar o arquivo do token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}
