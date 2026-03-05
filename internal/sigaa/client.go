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

func (s *SigaaService) GetTurmas() ([]models.Turma, error) {
	var turmas []models.Turma
	var anoSemestre string 

	c := s.Collector.Clone()
	c.OnHTML("input[name='javax.faces.ViewState']", func(e *colly.HTMLElement) {
		s.ViewState = e.Attr("value")
	})
	
	c.OnHTML("p.periodo-atual strong", func(e *colly.HTMLElement) {
		anoSemestre = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div#turmas-portal table tbody tr", func(e *colly.HTMLElement) {
		disciplina := strings.TrimSpace(e.ChildText("td.descricao a"))
		local := strings.TrimSpace(e.ChildText("td.info[style*='text-align:left'] "))
		horario := strings.TrimSpace(e.ChildText("td.info[style*='text-align:center']"))

		form := e.DOM.Find("form")
		nomeForm := form.AttrOr("id", "")
		idTurma := form.Find("input[name='idTurma']").AttrOr("value", "")
		paramBotao := form.Find("a[id$=':turmaVirtual']").AttrOr("id", "")

		if nomeForm == "" || idTurma == "" || paramBotao == "" {
    		return 
		}

		payload := map[string]string{
			nomeForm: nomeForm,
			"javax.faces.ViewState": s.ViewState,
			"idTurma": idTurma,
			paramBotao: paramBotao,
		}
		
		novaTurma := models.Turma{
			Disciplina: disciplina,
			Local: local,
			Horario: horario,
			AnoSemestre: anoSemestre,
			Payload: payload,
		}
		turmas = append(turmas, novaTurma)
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

	c := s.Collector.Clone()

	//Atualiza o ViewState
	c.OnHTML("input[name='javax.faces.ViewState']", func(e *colly.HTMLElement) {
		s.ViewState = e.Attr("value")
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
		// No sigaa, quando tem o "Visualizar tarefa enviada", ela já foi enviada
		tarefaEnviada := e.DOM.Parent().Find("a[title*='Visualizar']")
		// Se a tarefa já foi enviada, ela é ignorada aqui
		if tarefaEnviada.Length() > 0 {
			return
		}

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
	
	err = c.Post(urlHome, turma.Payload)
	if err != nil {
		return nil, fmt.Errorf("Erro ao visitar turma: %v", err)
	}
	return tarefas, nil 	
}
