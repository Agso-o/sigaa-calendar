package sigaa

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"github.com/Agso-o/sigaa-calendar/internal/models"
	"github.com/gocolly/colly/v2"
)

const (
	urlLoginGet  = "https://sigaa.ufpi.br/sigaa/verTelaLogin.do"
	urlLoginPost = "https://sigaa.ufpi.br/sigaa/logar.do?dispatch=logOn"
	urlHome      = "https://sigaa.ufpi.br/sigaa/portais/discente/discente.jsf"
	urlTurma     = "https://sigaa.ufpi.br/sigaa/portais/discente/turmas.jsf"
	urlTarefas   = "https://sigaa.ufpi.br/sigaa/ava/index.jsf"
)

var (
	reForm     = regexp.MustCompile(`document\.getElementById\('([^']+)'\)`)
	reTurma    = regexp.MustCompile(`'idTurma':'(\d+)'`)
	reBotao    = regexp.MustCompile(`{'([^']+)':'[^']+'`)
)

type SigaaService struct {
	Collector *colly.Collector
	User      string
	Pass      string
	ViewState string // O SIGAA precisa disso para cada click
}

// Retorna um ponteiro pra um serviço do sigaa 
// com as informações de login e viewState
func NewSigaaService(user string, pass string) (*SigaaService, error) {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		// colly.Async(), // Pode acelerar o serviço, mas gera riscos de bloqueio pelo sigaa
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	// Delay ao fazer requisições ao sistema
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*sigaa.ufpi.br*",
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	return &SigaaService{
		Collector: c,
		User: user,
		Pass: pass,

	}, nil
}

// Loga no sistema a a partir de um serviço
func (s *SigaaService) LoginSigaa() error{
	err := s.Collector.Visit(urlLoginGet)
	if err != nil {
		return fmt.Errorf("Erro ao acessar página de login: %v", err)
	}
	
	err = s.Collector.Post(urlLoginPost, map[string]string{
		"user.login": s.User,
		"user.senha": s.Pass,
	})
	
	if err != nil {
		return fmt.Errorf("Erro ao tentar logar no sistema: %v", err)
	}

	return nil
}

// Para implementar: Dados de turmas atuais só aparecem após as matriculas serem processadas 
func (s *SigaaService) GetTurmasAtuais() ([]models.Turma, error) {
	return nil, nil	
}

// Após o fim de um semestre, só essa informação vai ser disponível
func (s *SigaaService) GetTurmasAnteriores() ([]models.Turma, error) {
	var turmas []models.Turma
	
	c := s.Collector
	//Limpa o seletor
	c.OnHTMLDetach("table.listagem tr") 
    c.OnHTMLDetach("#turmas-portal span.mais a")
	// Guarda o ViewState da página
	c.OnHTML("input[name='javax.faces.ViewState']", func(e *colly.HTMLElement) {
		s.ViewState = e.Attr("value")
	})

	// Entra na aba de turmas anteriores
	c.OnHTML("#turmas-portal span.mais a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		e.Request.Visit(link)
	})
	
	var anoSemestre string
	// Pega os dados da turma
	c.OnHTML("table.listagem tr", func(e *colly.HTMLElement) {
		disciplina := strings.TrimSpace(e.ChildText("td:nth-child(1)"))
		creditos := strings.TrimSpace(e.ChildText("td:nth-child(4)"))
		horario := strings.TrimSpace(e.ChildText("td:nth-child(5)"))

		if strings.Contains(disciplina, ".1") || strings.Contains(disciplina, ".2") {
			anoSemestre = disciplina
			return
		}

		if disciplina == "" || strings.Contains(disciplina, "Disciplina") {
			return
		}

		// Lógica para pegar os dados do botão de ir pra turma 
		onclickText := e.ChildAttr("td:last-child a", "onclick")

		matchForm := reForm.FindStringSubmatch(onclickText)
		matchTurma := reTurma.FindStringSubmatch(onclickText)
		matchBotao := reBotao.FindStringSubmatch(onclickText)

		if len(matchForm) > 1 && len(matchTurma) > 1 && len(matchBotao) > 1 {
			nomeForm := matchForm[1]  
			idTurma := matchTurma[1]  
			paramBotao := matchBotao[1]

			Payload := map[string]string{
				nomeForm:                	nomeForm,  
				"javax.faces.ViewState":	s.ViewState,
				"idTurma":               	idTurma,
				paramBotao:              	paramBotao, 
			}

			novaTurma := models.Turma{
				Disciplina: disciplina,
				Horario: horario,
				AnoSemestre: anoSemestre,
				Creditos: creditos,
				Payload: Payload,
			}

			turmas = append(turmas, novaTurma)
		} 
	})

	err := c.Visit(urlHome)	
	if err != nil {
		return nil, fmt.Errorf("Erro ao visitar a home: %v", err)
	}
	if len(turmas) == 0 {
		return nil, fmt.Errorf("Nenhuma turma encontrada")
	}
	return turmas, nil
}

func (s *SigaaService) GetTarefasByTurma(turma models.Turma) ([]models.Tarefa, error) {
	var tarefas []models.Tarefa

	c := s.Collector


	// Limpa Seletores
	c.OnHTMLDetach("div#barraEsquerda a")
	c.OnHTMLDetach("td.first[style*='bold']")
	c.OnHTMLDetach("table.listagem tr")
	c.OnHTMLDetach("#turmas-portal span.mais a")
	
	c.SetRequestTimeout(30 * time.Second)

	//Atualiza o ViewState
	c.OnHTML("input[name='javax.faces.ViewState']", func(e *colly.HTMLElement) {
		s.ViewState = e.Attr("value")
	})
	
	// Entra na aba de turmas anteriores
	c.OnHTML("#turmas-portal span.mais a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		e.Request.Visit(link)
	})

	c.OnHTML("div#barraEsquerda a", func(e *colly.HTMLElement) {
		textoBotao := strings.TrimSpace(e.ChildText(".itemMenu"))

		if !strings.Contains(textoBotao, "Tarefas") {
			return
		}

		onclickText := e.Attr("onclick")
		match := reBotao.FindStringSubmatch(onclickText)
			
		if len(match) > 1 {
			c.OnHTMLDetach("div#barraEsquerda a")
			paramBotao := match[1]
			dadosPost := map[string]string{
				"formMenu": "formMenu",
				"javax.faces.ViewState": s.ViewState,
				paramBotao: paramBotao,
			}
			c.Post(urlTarefas, dadosPost)
		}
	})

	c.OnHTML("td.first[style*='bold']", func(e *colly.HTMLElement) {
		titulo := strings.TrimSpace(e.Text)
		descricao := e.DOM.Parent().Next().Find("td.first p").Text()
		descricao = strings.TrimSpace(descricao)
		dataRaw := e.DOM.Parent().Find("td[style*='center']").Text()
		dataLimpa := strings.ReplaceAll(dataRaw, "\n", " ")
        dataLimpa = strings.ReplaceAll(dataLimpa, "\t", "")
        dataLimpa = strings.TrimSpace(dataLimpa)

		endTime, _ := ParseDataEntrega(dataLimpa)

		novaTarefa := &models.Tarefa{
			Disciplina: turma.Disciplina,
			Titulo: titulo,
			Descricao: descricao,
			DataVencimento: endTime,	
		}

		tarefas = append(tarefas, *novaTarefa)
	})

	err := c.Visit(urlHome)
	if err != nil {
		return nil, fmt.Errorf("erro ao resetar navegação: %v", err)
	}

	turma.Payload["javax.faces.ViewState"] = s.ViewState
	
	err = c.Post(urlTurma, turma.Payload)
	if err != nil {
		return nil, fmt.Errorf("Erro ao visitar turma: %v", err)
	}
	return tarefas, nil 	
}
