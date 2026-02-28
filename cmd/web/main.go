package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Agso-o/sigaa-calendar/internal/wrapper"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

func init() {
	_ = godotenv.Load()
	clientID := os.Getenv("CLIENTID")
	clientSecret := os.Getenv("CLIENTSECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("Incapaz de ler as variáveis de ambiente")
	}
	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://sigaa-calendar.onrender.com/callback", // A rota de retorno
		Scopes: []string{
			"https://www.googleapis.com/auth/tasks",
			"https://www.googleapis.com/auth/calendar",
			"https://www.googleapis.com/auth/userinfo.email",

		},
		Endpoint: google.Endpoint,
	}
}
	
func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("super-secreto"))
	r.Use(sessions.Sessions("sigaa_session", store))

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login_google.html", nil)
	})

	r.GET("/login/google", func(ctx *gin.Context) {
		url := googleOAuthConfig.AuthCodeURL("state-seguro")
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/callback", func(ctx *gin.Context) {
		code := ctx.Query("code")
		c := context.Background()
		
		// Troca o código pelo Token real
		token, err := googleOAuthConfig.Exchange(c, code)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Erro no Google")
			return
		}

		// Transforma o Token em string (JSON) para guardar no Cookie
		tokenJSON, _ := json.Marshal(token)
		
		// Salva o token na Sessão do usuário
		session := sessions.Default(ctx)
		session.Set("google_token", string(tokenJSON))
		session.Save()

		// REDIRECIONA PARA A TELA DO SIGAA!
		ctx.Redirect(http.StatusFound, "/sigaa")
	})

	r.GET("/sigaa", func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("google_token") == nil {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "index.html", nil)
	})

	
	r.POST("/sincronizar", func(ctx *gin.Context) {
		usuario := ctx.PostForm("usuario")
		senha := ctx.PostForm("senha")

		flagTarefas := ctx.PostForm("tarefas") == "true"
		flagAulas := ctx.PostForm("aulas") == "true"	

		// 2. Resgata o Token do Google que estava guardado no Cookie
		session := sessions.Default(ctx)
		tokenString := session.Get("google_token")
		if tokenString == nil {
			ctx.String(http.StatusUnauthorized, "Sessão expirada. Volte ao início.")
			return
		}

		var token oauth2.Token
		json.Unmarshal([]byte(tokenString.(string)), &token)
		
		c := context.Background()
		googleClient := googleOAuthConfig.Client(c, &token)

		// Sincroniza o calendário com o sigaa no plano de fundo pro usuário não ter que esperar
		// Essa função pode demorar até 2min 
		go func(user, pass string, client *http.Client) {
			errSync := wrapper.SigaaSync(flagAulas, flagTarefas, user, pass, client)
			if errSync != nil {
				fmt.Printf("Falha no background de %s: %v\n", user, errSync)
			} else {
				fmt.Printf("Sincronização de %s concluída com sucesso!\n", user)
			}
		}(usuario, senha, googleClient)
		
		var emailDoUsuario string
		resp, err := googleClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err == nil {
			defer resp.Body.Close()
			var userInfo struct {
				Email string `json:"email"`
			}
			json.NewDecoder(resp.Body).Decode(&userInfo)
			emailDoUsuario = userInfo.Email
		}

		session.Clear()
		session.Save()

		urlDestino := "https://calendar.google.com/calendar/r"
		
		if emailDoUsuario != "" {
			urlDestino = fmt.Sprintf("https://calendar.google.com/calendar/r?authuser=%s", emailDoUsuario)
		}

		ctx.Redirect(http.StatusSeeOther, urlDestino)
	})

	r.Run()
}

