# cep-weather-api

API em Go que recebe um CEP, descobre a cidade (via [ViaCEP](https://viacep.com.br/))
e retorna a temperatura atual (via [WeatherAPI](https://www.weatherapi.com/)) em
Celsius, Fahrenheit e Kelvin.

## URL de produção (Cloud Run)

> **TODO:** cole aqui a URL ativa após o deploy, ex.:
> `https://cep-weather-api-939917444284.us-central1.run.app`

Exemplo de requisição:

```bash
curl https://cep-weather-api-939917444284.us-central1.run.app/01310100
```

## Contrato da API

`GET /{cep}`

| Cenário | HTTP | Corpo |
|---|---|---|
| Sucesso | `200` | `{ "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65 }` |
| CEP com formato inválido (≠ 8 dígitos) | `422` | `invalid zipcode` |
| CEP válido mas inexistente | `404` | `can not find zipcode` |

Fórmulas de conversão: `F = C * 1.8 + 32` e `K = C + 273.15`.

## Arquitetura

```
cmd/                         # entrypoint + wiring
configs/                     # configuração via variáveis de ambiente
internal/
  entity/                    # domínio: Temperature (regra de conversão), erros
  dto/                       # payload de resposta
  usecase/                   # orquestração: CEP -> cidade -> clima -> conversão
  infra/
    gateway/                 # clientes HTTP (ViaCEP, WeatherAPI)
    web/handlers/            # handler HTTP fino (só valida e delega)
pkg/                         # utilitários (validação de CEP)
```

Fluxo: o handler valida o formato do CEP e delega ao usecase, que consulta o ViaCEP
para obter a cidade e a WeatherAPI para a temperatura, convertendo o resultado nas três
escalas. A lógica de negócio fica no domínio/usecase; o handler apenas traduz para HTTP.

## Configuração

Variáveis de ambiente (veja `.env.example`):

| Variável | Default | Descrição |
|---|---|---|
| `APP_KEY` | — (obrigatória) | Chave da WeatherAPI |
| `WEATHER_API_URL` | `http://api.weatherapi.com/v1/current.json` | Endpoint da WeatherAPI |
| `VIACEP_URL` | `http://viacep.com.br/ws/%s/json/` | Endpoint do ViaCEP (`%s` = CEP) |
| `TIMEOUT` | `5s` | Timeout das chamadas externas |
| `PORT` | `8080` | Porta HTTP |

```bash
cp .env.example .env
# edite .env e preencha APP_KEY com sua chave da WeatherAPI
```

## Rodando localmente

### Via Docker (recomendado)

```bash
docker compose up --build
# ou, sem compose:
docker build -t cep-weather-api .
docker run --rm -p 8080:8080 --env-file .env cep-weather-api
```

### Via Go

```bash
export $(grep -v '^#' .env | xargs)   # carrega as variáveis do .env
go run ./cmd
```

Testando os cenários:

```bash
curl -i localhost:8080/01310100   # 200 + temperaturas
curl -i localhost:8080/123        # 422 invalid zipcode
curl -i localhost:8080/99999999   # 404 can not find zipcode
```

## Testes

```bash
go test ./...            # todos os testes
go test ./... -cover     # com cobertura
```

## Deploy no Google Cloud Run

```bash
# 1. Build & push da imagem (Cloud Build)
gcloud builds submit --tag gcr.io/<PROJECT_ID>/cep-weather-api

# 2. Deploy (a APP_KEY vai como variável de ambiente do serviço, não na imagem)
gcloud run deploy cep-weather-api \
  --image gcr.io/<PROJECT_ID>/cep-weather-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars APP_KEY=<SUA_CHAVE_WEATHERAPI>
```

O Cloud Run injeta automaticamente a variável `PORT`; a aplicação já a respeita.
Após o deploy, copie a URL retornada para a seção [URL de produção](#url-de-produção-cloud-run).
