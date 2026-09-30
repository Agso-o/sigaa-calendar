<h1 align="center">Sigaa Calendar</h1>

Este projeto foi desenvolvido com o objetivo de organizar meus próprios horários e auxiliar estudantes em sua jornada acadêmica, tentando melhorar a experiência de quem usa o sistema da Universidade Federal do Piauí (UFPI). Para isso, o **Sigaa Calendar** conecta sistema antigo do SIGAA ao sistema moderno do Google Calendar, sincronizando horários de aula, datas de entrega de trabalhos/atividades e datas de renovação da biblioteca.

#
## Status do Projeto
O projeto foi lançado oficialmente, mas mantém o foco em correção de bugs, melhorias, e na implementação da funcionalidade de sincronizar os prazos de renovação de livros da biblioteca
## Funcionalidades
- **Sincronização de Horários de aulas**: Salva todos os seus horários de aulas direto do SIGAA para uma nova agenda no Google Calendar como eventos recorrentes que começam do início do período e seguem até o final automaticamente

- **Sincronização de Prazos de entrega**: Salva todos os prazos de entrega dos trabalhos e atividades que ainda não foram entregues no sigaa direto para o Google Tasks. Esses prazos são salvos como tarefas contendo o título, descrição, e data de entrega, e podem ser marcados como concluidos pelo usuário direto pela plataforma do Google Tasks.

- **Sincronização de datas de renovação da Biblioteca**: Salva a data de renovação dos emprestimo de livros da biblioteca como eventos no Google Calendar na agenda criada pelo programa

## Guia de uso
O SIGAA-Calendar oferece duas formas de interagir com a ferramenta: uma interface gráfica via web, para maior acessibilidade e uma interface de linha de comando através do programa compilado.
> Por não ser um app verificado pela Google, esta ferramenta possui um limite de usuários. Portanto, antes de conseguir utilizar uma das interfaces abaixo, é necessário ter seu e-mail adicionado na lista de usuários testadores.
> 
> Você pode obter esse acesso entrando em contato comigo, caso haja disponibilidade de vagas, ou utilizando sua própria conta do Google Cloud e compilando o projeto de forma independente.

### Modo Web
Sem instalação e de uso rápido, ideal para usuários mobile
1. Acesse: https://sigaa-calendar.onrender.com
2. Clque em ``Continuar com o Google``
3. Escolha a mesma conta que você cadastrou como usuário e selecione todas as permissões necessárias (o programa não irá funcionar caso você não conceda as devidas permissões)
4. Digite seu usuário e senha do sigaa (obs: Atualmente o sistema não identifica senhas digitadas incorretamente, certifique-se de que os dados estão corretos)
5. Selecione oque deseja sincronizar (Horários de aula deverão ser sincronizados somente no primeiro acesso, se não houver mudanças no horário após isso essa sincronização é redundante)
6. Clique em ``Sincronizar agora`` e pronto, seus dados seráo sincronizados no plano de fundo.

### Modo CLI
Essa interface tende a ser um pouco menos acessível, porém bem mais poderosa, permitindo automações localmente. Ela exige um pouco mais de configuração, então antes de prosseguir siga as instruções disponíveis em: **[instalação](#instalação)**

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
  -lib
    	Sincroniza prazos de entrega de livros com o Google Tasks
  -tarefas
    	Sincroniza prazos de entrega de trabalhos com o Google Tasks
```
Exemplos de uso:
```bash
# Sincronizará apenas tarefas
sigaa-calendar --tarefas

# Sincronizará apenas aulas
sigaa-calendar --aulas

# sincronizará apenas as datas da biblioteca
sigaa-calendar --lib

# Para sincronizar ambos os dados
sigaa-calendar --tarefas --aulas --lib

```

> Fique tranquilo, seus dados não serão guardados em nenhuma etapa dessa aplicação (seja web ou CLI), o usuário e senha do SIGAA são utilizados somente para o acesso ao sistema, e são descartados imediatamente. Durante o login esses dados ficam apenas na memória RAM e não são armazenados de nenhuma forma.

## Instalação
O Sigaa Calendar utiliza a API do Google para sincronizar seus dados. Como o binário roda localmente na sua máquina (Modo CLI), você precisa fornecer suas próprias credenciais do Google Cloud e do SIGAA para que o programa funcione.
### 1. Credenciais do Google Cloud (Acesso à Agenda e Tasks)
Para que a ferramenta consiga criar eventos e tarefas na sua conta Google, você precisa **[criar um projeto no Google Cloud](https://console.cloud.google.com/projectcreate )** e gerar as credenciais de ID do cliente OAuth: Client ID e Client Secret

Siga o **[Tutorial Oficial do Google: Criar credenciais de acesso](https://developers.google.com/workspace/guides/create-credentials?hl=pt-br#oauth-client-id)** para gerar sua chave.

#### Pontos de atenção ao gerar as credenciais:
* Ative a **[Google Calendar API](https://console.cloud.google.com/apis/library/calendar-json.googleapis.com)** e a **[Google Tasks API](https://console.cloud.google.com/apis/library/tasks.googleapis.com)** no seu projeto.
* Na aba de publico alvo, selecione externo e adicione o seu próprio e-mail na lista de "Usuários de teste".
* Na hora de criar o cliente OAuth, escolha o tipo **App para computador** (Desktop app).

Após criar o seu cliente OAuth, copie e guarde para a próxima etapa as credenciais: ID do cliente e Chave secreta do cliente

### 2. Variáveis de ambiente
Para manter sua segurança e não deixar senhas e credenciais salvas em texto puro dentro do código, o programa lê seus dados do SIGAA e do Google Cloud diretamente das variáveis de ambiente do seu sistema operacional.

#### No Linux ou macOS
Configure as variáveis no seu terminal (você pode exportar temporariamente ou adicionar essas linhas ao seu `~/.bashrc` ou `~/.zshrc` para ficar permanente):
```Bash
export SIGAA_USER="seu_usuario"
export SIGAA_SENHA="sua_senha_do_sigaa"
export GOOGLE_CLIENTID="seu_ID_do_cliente"
export GOOGLE_CLIENTSECRET="sua_chave_secreta_do_cliente"
```
#### No Windows 
Temporáio:
```PowerShell
$env:SIGAA_USER="seu_usuario" 
$env:SIGAA_SENHA="sua_senha_do_sigaa" 
$env:GOOGLE_CLIENTID="seu_ID_do_cliente" 
$env:GOOGLE_CLIENTSECRET="sua_chave_secreta_do_cliente"
```
Permanente:
```PowerShell
setx SIGAA_USER "seu_usuario"
setx SIGAA_SENHA "sua_senha_do_sigaa"
setx GOOGLE_CLIENTID "seu_ID_do_cliente"
setx GOOGLE_CLIENTSECRET "sua_chave_secreta_do_cliente"
```
### 3. Baixando o executável (Sem instalar o Go)
Se você não é desenvolvedor ou não quer instalar e configurar o ambiente Go na sua máquina, a forma mais rápida de utilizar a ferramenta é baixando a versão pronta:

1. Acesse a página de **[Releases](https://github.com/Agso-o/sigaa-calendar/releases)** do repositório.
2. Baixe o arquivo compactado correspondente ao seu sistema operacional (Windows, Linux ou macOS).
3. Extraia o arquivo baixado. Você terá o executável `sigaa-calendar` (ou `sigaa-calendar.exe` no Windows) pronto para uso.

#### Executando
Se você optou por baixar o arquivo pronto, abra o terminal na pasta onde você extraiu o arquivo e execute o programa diretamente (você pode ignorar as instruções de compilação do passo a seguir).

No Linux ou macOS:
```bash
chmod +x sigaa-calendar
./sigaa-calendar --tarefas --aulas
```
No Windows:
```
.\sigaa-calendar.exe --tarefas --aulas
```

### 4. Compilando e executando o programa
O projeto foi desenvolvido na linguagem Go. Portanto, é necessário ter o Go instalado e configurado corretamente no sistema para compilar e executar o programa.

Agora que está tudo configurado você pode baixar e compilar o projeto normalmente:

Para Linux ou macOS:
```Bash
git clone https://github.com/Agso-o/sigaa-calendar
cd sigaa-calendar
go build -o sigaa-calendar ./cmd/cli/
```
Para Windows:
```DOS
git clone https://github.com/Agso-o/sigaa-calendar
cd sigaa-calendar
go build -o sigaa-calendar.exe .\cmd\cli\
```
#### Executando 

No Linux ou macOS:
```bash
chmod +x sigaa-calendar
./sigaa-calendar --tarefas --aulas
```
No Windows:
```DOS
.\sigaa-calendar.exe --tarefas --aulas
```

Para instalar o programa no sistema Linux ou macOS e poder utilizá-lo fora da raiz do projeto, mova o arquivo para um diretório no $PATH:
```Bash
mv sigaa-calendar ~/.local/bin
```
(Certifique-se de que esse diretório está no path com: `echo $PATH`)

Mais informações de uso em [Como utilizar o modo CLI](#modo-cli)

Na primeira vez que você rodar o comando, a primeira ação do programa será abrir o seu navegador padrão automaticamente, solicitando que você faça login na sua conta do Google e autorize o aplicativo. Após aceitar, um arquivo de token (ex: Tasks+CalendarToken.json) será gerado na pasta de configurações do seu sistema (no Linux, o padrão é ~/.config/sigaa-calendar/).
Com o token salvo, as próximas execuções acontecerão via terminal, sem autorização extra no navegador.

## Como contribuir
- Caso queira contribuir, você pode abrir uma [issue](https://github.com/Agso-o/sigaa-calendar/issues) com algo que você gostaria que seja corrigido ou adicionado ao projeto.
- Este projeto foca em organização de prazos e horários, qualquer idéia adicionada as issues deve ser limitada a esse escoṕo

## Avisos legais
- Este software é uma **iniciativa independente** e não possui vínculo oficial com o desenvolvimento do SIGAA ou com a UFPI.
- O usuário é responsável por verificar a precisão dos dados informados por esse sistema, visto que, qualquer alteração ou instabilidade no portal, poderá alterar os resultados e/ou afetar o funcionamento da ferramenta.
- Este projeto foi inteiramente desenvolvido com foco no sistema da UFPI, o uso dessa ferramenta em outras IFES que utilizam o SIGAA não foi testado e portanto deve ser evitado
