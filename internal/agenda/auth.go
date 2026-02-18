package agenda

import (
	"os"
	"fmt"
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"encoding/json"
)

// Resgata um token, salva o token e retorna o cliente
func getClient(config *oauth2.Config) (*http.Client, error) {

	tokFile := "token.json"
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
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("Erro ao requerir token da web: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("Erro ao requerir token da web: %v", err)
	}
	return tok, nil
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
