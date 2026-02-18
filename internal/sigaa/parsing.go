package sigaa

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Agso-o/sigaa-calendar/internal/models"
)

var mapaHorariosInicio = map[string]int{
	// Manhã
	"M1": 6,  // 06:00
	"M2": 7,  // 07:00
	"M3": 8,  // 08:00
	"M4": 9,  // 09:00
	"M5": 10, // 10:00
	"M6": 11, // 11:00
	// Tarde
	"T1": 12, // 12:00
	"T2": 13, // 13:00
	"T3": 14, // 14:00
	"T4": 15, // 15:00
	"T5": 16, // 16:00
	"T6": 17, // 17:00
	// Noite
	"N1": 18, // 18:00
	"N2": 19, // 19:00
	"N3": 20, // 20:00
	"N4": 21, // 21:00
}

var mapaDias = map[rune]string{
	'1': "SU", // Domingo
	'2': "MO", // Segunda
	'3': "TU", // Terça
	'4': "WE", // Quarta
	'5': "TH", // Quinta
	'6': "FR", // Sexta
	'7': "SA", // Sábado
}

// Recebe a string de data extraida do sigaa, no formato: 
// de DD/MM/YYYY ás HH:MM até DD/MM/YYYY ás HH:MM
// E devolve um tipo time no formato RC33399
func ParseDataEntrega(dataRaw string) (time.Time, error) {
	//  de 17/03/2025 às 00h00 a  23/03/2025 às 23h59 
	dataRaw = strings.Join(strings.Split(strings.TrimSpace(dataRaw), " ")[5:], " ")
	dataRaw = strings.TrimSpace(dataRaw)
	zone, err := time.LoadLocation(models.TimeZone)
	if err != nil {
		return time.Now(), err
	}
	t, err := time.ParseInLocation("02/01/2006 às 15h04", dataRaw, zone)
	return t, err
}

// NECESSITA DE MUDANÇAS  - Como eu vou tratar o semestre no novo HTML? qual vai ser o inicio e o fim do periodo?
// Se manter a logica de data fixa para inicio e fim das aulas, colocar como global
func ParseHorarioAgenda(horarioRaw, anoSemestre string) (recorrencia string, startTime, endTime time.Time, err error) {
	horarioRaw = strings.TrimSpace(horarioRaw)

	partesAno := strings.Split(anoSemestre, ".")
    if len(partesAno) < 2 {
        return "", time.Time{}, time.Time{}, fmt.Errorf("ano/semestre inválido ou não encontrado: '%s'", anoSemestre)
    }

    ano := partesAno[0]
    semestre := partesAno[1]
	
	re := regexp.MustCompile(`(\d+)([MTN])(\d+)`)
	matches := re.FindStringSubmatch(horarioRaw)

	if len(matches) < 4 {
		return "", time.Time{},time.Time{}, fmt.Errorf("Formato deconhecido")
	}

	dias := matches[1]
	turno := matches[2]
	horarios := matches[3]

	primeiroHorario := string(horarios[0]) 
	chaveInicio := turno + primeiroHorario

	horarioInicio, ok := mapaHorariosInicio[chaveInicio]
	
	if !ok {
		return "", time.Time{},time.Time{}, fmt.Errorf("Erro ao converter horário: Invalid Maping")

	}

	duracaoHoras := len(horarios)

	loc, _ := time.LoadLocation(models.TimeZone)

	// Data de inicio do primeiro semestre: 10 de março
	dataAux := "10/03/"+ano
	if semestre == "2" {
		// Data de inicio do segundo semestre
		dataAux = "10/08/"+ano	
	}
	inicio, err := time.Parse("02/01/2006", dataAux)

	dataFim := ano + "0715"
	if semestre == "2" {
		//Fim do segundo semestre, 15 de dezembro
		dataFim = ano + "1215" 
	}

	if err != nil {
		return "", time.Time{}, time.Time{}, fmt.Errorf("Não foi possível converter formato de dada: %v", err)
	}

	startTime = time.Date(inicio.Year(), inicio.Month(), inicio.Day(), horarioInicio, 0, 0, 0, loc)
	endTime = startTime.Add(time.Duration(duracaoHoras) * time.Hour)

	var diasRrule []string
	for _, dia := range dias {
		if val, ok := mapaDias[dia]; ok {
			diasRrule = append(diasRrule, val)
		}
	}


	recorrencia = fmt.Sprintf("RRULE:FREQ=WEEKLY;BYDAY=%s;UNTIL=%sT235959Z", strings.Join(diasRrule, ","), dataFim)

	return recorrencia, startTime, endTime, nil

}
