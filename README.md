# FIAP X - API de Gestão de Vídeos

Esta API é o ponto central do sistema FIAP X, responsável pela gestão de utilizadores, autenticação e orquestração do pipeline de processamento de vídeos.

## 🚀 Funcionalidades Principais

* **Autenticação Segura**: Registo e Login de utilizadores com hashing de passwords (BCrypt) e tokens JWT.
* **Pipeline de Vídeo**: Upload de ficheiros diretamente para o Amazon S3 e disparo de mensagens para a fila SQS.
* **Gestão de Histórico**: Listagem do estado de processamento dos vídeos do utilizador.
* **Download Seguro**: Geração de URLs pré-assinadas (Presigned URLs) para download dos frames processados.
* **Documentação Viva**: Interface Swagger integrada para testes de endpoints.

## 🏗️ Arquitetura

A aplicação segue os princípios da **Clean Architecture**, garantindo desacoplamento entre a lógica de negócio e a infraestrutura:

* **Domain/Entity**: Regras essenciais (User, Video).
* **Usecase**: Orquestração das operações de negócio.
* **Infra/Interface Adapters**: Repositórios GORM, Clientes AWS e Handlers Gin Gonic.

## 🛠️ Tecnologias Utilizadas

* **Linguagem**: Go 1.24
* **Framework Web**: Gin Gonic
* **Base de Dados**: PostgreSQL (via GORM)
* **Serviços AWS (SDK v2)**: S3, SQS e Secrets Manager.

## 📦 Configuração e Variáveis de Ambiente

A API prioriza o **AWS Secrets Manager** para credenciais sensíveis, recorrendo a variáveis locais apenas como *fallback* em Desenvolvimento.

| Variável | Descrição | Exemplo |
| --- | --- | --- |
| `PORT` | Porta de escuta da API | `8080` |
| `DB_SECRET_NAME` | Nome do segredo no Secrets Manager | `db-credentials` |
| `AWS_REGION` | Região da infraestrutura | `us-east-1` |
| `AWS_QUEUE_URL` | URL da fila SQS para processamento | `https://sqs...` |
| `JWT_SECRET` | Chave para assinatura dos tokens | `sua_chave_secreta` |

## 🚀 Como Executar

### Via Docker Compose (Ambiente Completo)

Para que o ecossistema de microserviços funcione corretamente via docker-compose, deves garantir os seguintes requisitos:

#### 1. Disponibilidade dos Projetos
O ficheiro docker-compose.yml está configurado para orquestrar múltiplos serviços que residem em repositórios distintos. Para uma subida íntegra, deves:

* Ter as pastas [hackaton-service-api](https://github.com/fiap-11SOAT-2025/hackaton-service-api). e [hackaton-service-worker](https://github.com/fiap-11SOAT-2025/hackaton-service-worker) no mesmo nível de diretório (ou ajustar os caminhos no build: context do compose).

* Garantir que ambos os projetos foram baixados por completo, incluindo os seus respetivos Dockerfile e dependências de código.

#### 2. Configuração do LocalStack (Scripts de Inicialização)
O mapeamento de volumes para o script init-localstack.sh pode variar dependendo do teu sistema operativo e da forma como o Docker está instalado:

* Docker no WSL2 (Recomendado): O mapeamento ./scripts/init-localstack.sh:/etc/localstack/init/ready.d/init-aws.sh funciona nativamente.

* Docker Desktop (Windows nativo): Se não estiveres a usar o WSL, deves garantir que o Docker tem permissões de leitura na pasta do projeto (Shared Drives).

Por fim, na raiz do projeto principal, execute:

```bash
docker-compose up --build
```


### Manualmente (Desenvolvimento)

1. Certifique-se de que o PostgreSQL e o LocalStack (ou AWS) estão acessíveis.
2. Execute a aplicação:

```bash
go run cmd/main.go
```

## 📖 Documentação da API (Swagger)

Com a API a rodar, aceda à documentação interativa:
`http://localhost:8080/swagger-ui/index.html`

## 🧪 Testes

Para garantir a integridade da lógica de autenticação e gestão de vídeos:

```bash
go test -coverprofile=coverage.out ./internal/usecase/...
```
Para gerar em html:

```bash
go tool cover -html=coverage.out
```

