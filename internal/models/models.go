package models

import "time"


const (
	TimeZone string = "America/Sao_Paulo"
)

type Turma struct {
	Disciplina string
	Local string
	Horario string // Formato Sigaa
	Creditos string // Formato Sigaa, ex: 4 (60h)
	AnoSemestre string // Semestre no sigaa: YYYY.1 ou YYYY.2
	Payload map[string] string // Informações para entrar na página da turma do Sigaa

}

type Tarefa struct {
	Disciplina string
	Titulo string
	Descricao string
	DataVencimento time.Time
	SigaaID string

}

type Aula struct {
	Disciplina string
	Local string
	StartTime time.Time
	EndTime	time.Time
	Recorrencia string

}

type Emprestimo struct {
	Livro string
	Biblioteca string
	Prazo time.Time
}
