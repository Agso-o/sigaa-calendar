<h1 align="center">Sigaa Calendar</h1>

Não perca a hora e automatize seu semestre: com o **Sigaa Calendar**, todos os seus prazos de entrega e horarios de aula, que antes estavam presos atrás de uma lenta navegação pelo SIGAA (Sistema Integrado de Gestão de Atividades Acadêmicas), agora podem ir automaticamente para a sua agenda do Google com um simples comando.

Este projeto foi desenvolvido com o objetivo de organizar meus próprios horários e auxiliar estudantes em sua jornada acadêmica, tentando melhorar a experiência de quem usa o sistema da Universidade Federal do Piauí (UFPI). Para isso, o **Sigaa Calendar** conecta sistema antigo do SIGAA ao sistema moderno do Google Calendar, sincronizando horários de aula, datas de entrega de trabalhos/atividades e datas de renovação da biblioteca.

#
## Funcionalidades
- **Sincronização de Horários de aulas**: Salva todos os seus horários de aulas direto do SIGAA para uma nova agenda no Google Calendar como eventos recorrentes que começam do início do período e seguem até o final automaticamente

- **Sincronização de Prazos de entrega**: Salva todos os prazos de entrega dos trabalhos e atividades que ainda não foram entregues no sigaa direto para o Google Tasks. Esses prazos são salvos como tarefas contendo o título, descrição, e data de entrega, e podem ser marcados como concluidos pelo usuário direto pela plataforma do Google Tasks.

- **Sincronização de datas de renovação da Biblioteca**: ---Em desenvolvimento---

## Guia de uso
O SIGAA-Calendar oferece duas formas de interagir com a ferramenta: uma interface gráfica via web, para maior acessibilidade e uma interface de linha de comando através do binário disponível em [releases](https://github.com/Agso-o/sigaa-calendar/releases)
> Por não ser um app verificado pela Google, esta ferramenta possui um limite de usuários. Portanto, antes de conseguir utilizar uma das interfaces abaixo, é necessário ter seu e-mail adicionado na lista de usuários testadores.
> 
> Você pode obter esse acesso entrando em contato comigo, caso haja disponibilidade de vagas, ou utilizando sua própria conta do Google Cloud e compilando o projeto de forma independente.

### **Modo Web**
Sem instalação e de uso rápido, ideal para usuários mobile
1. Acesse: https://sigaa-calendar.onrender.com
2. Clque em ``Continuar com o Google``
3. Escolha a mesma conta que você cadastrou como usuário e selecione todas as permissões necessárias (o programa não irá funcionar caso você não conceda as devidas permissões)
4. Digite seu usuário e senha do sigaa (obs: Atualmente o sistema não identifica senhas digitadas incorretamente, certifique-se de que os dados estão corretos)
5. Selecione oque deseja sincronizar (Horários de aula deverão ser sincronizados somente no primeiro acesso, se não houver mudanças no horário após isso essa sincronização é redundante)
6. Clique em ``Sincronizar agora`` e pronto, seus dados seráo sincronizados no plano de fundo.

### **Modo CLI**
Essa interface tende a ser um pouco menos acessível, porém bem mais poderosa, permitindo automações localmente. Ela exige um pouco mais de configuração, então antes de prosseguir siga as instruções disponíveis em: [instalação](#instalacao)

Após instalar a ferramenta, você poderá ver o guia a seguir digitando o seguinte comando:
```bash
sigaa-calendar --help
```
O programa deverá responder com a seguinte mensagem:
``` bash
----------Sigaa Calendar----------

Uso:
 sigaa-calendar [flags]

Flags disponíveis:
  -aulas
        Sincroniza os horários das aulas com o Google Agenda
  -tarefas
        Sincroniza prazos de entrega de trabalhos com o Google Tasks
Padrão: Sincroniza somente tarefas caso nenhuma flag seja fornecida
```
Exemplos:
```bash
# Sincronizará apenas tarefas
sigaa-calendar --tarefas

# Sincronizará apenas aulas
sigaa-calendar --aulas

# Para sincronizar ambos os dados
sigaa-calendar --tarefas --aulas

```

> Fique tranquilo, seus dados não serão guardados em nenhuma etapa dessa aplicação (seja web ou CLI), o usuário e senha do SIGAA são utilizados somente para o acesso ao sistema, e são descartados imediatamente. Durante o login esses dados ficam apenas na memória RAM e não são armazenados de nenhuma forma.

## Instalação

## Como contribuir
- O desenvolvimento desse projeto faz parte dos meus estudos pessoais, portanto, atualmente este repositório não aceita contribuições de códigos de terceiros.
- Caso queira contribuir, você pode abrir uma [issue](https://github.com/Agso-o/sigaa-calendar/issues) com algo que você gostaria que seja corrigido ou adicionado ao projeto.
- Este projeto foca em organização de prazos e horários, qualquer idéia adicionada as issues deve ser limitada a esse escoṕo

## Avisos legais
- Este software é uma **iniciativa independente** e não possui vínculo oficial com o desenvolvimento do SIGAA.
- O usuário é responsável por verificar a precisão dos dados informados por esse sistema, visto que, qualquer alteração ou instabilidade no portal, poderá alterar os resultados e/ou afetar o funcionamento da ferramenta.
- Este projeto foi inteiramente desenvolvido com foco no sistema da UFPI, o uso dessa ferramenta em outras IFES que utilizam o SIGAA não foi testado e portanto deve ser evitado
