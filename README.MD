# Visão geral

- "Serviço de banco de dados NoSQL rápido e flexível para qualquer escala"
- Totalmente gerenciado pela Amazon
- Fácil de escalar
- Fácil de usar em qualquer linguagem (SDKs Amazon: Go, Java, Python, Node.js)
- Suporte a transações (Até 25 itens / requisição)

# Configuração DynamoDB Local

Baixar e rodar a versão local que a Amazon disponibiliza:

    docker pull amazon/dynamodb-local
    docker run -p 8000:8000 amazon/dynamodb-local

## Instalando o dynamodb admin

    npm install -g dynamodb-admin
    dynamodb-admin

### Resultado:

- DynamoDB respondendo na porta 8000
- DynamoDB Admin respondendo em http://localhost:8001 <br>
- [DynamoLocalAdmin](http://localhost:8001)

## Up and running!

# Como conectar
### Código !!! <br>
<code>internal/connection/connection.go</code>

#Create table

    Cria uma tabela no DynamoDB
    Devem ser especificados os campos de chave. 
    
Os modelos de criação de chave são: <br>

    Hash: Apenas uma chave, esta é única na tabela (igual uma primmary key de um banco relacional)
    Hash+Range: Uma chave de partição e uma de ordenação

## Chaves

### Partição
    Chaves de partição (HASH) devem sempre ser buscadas usando o operador =.
    Elas recebem esse nome por causa do algoritimo interno que o DynamoDB usa para
    distribuir as chaves em partições.
    É por este motivo que para fazer queries nas chaves de partição, temos que sempre usar o =
    
### Range
    Chaves de ordenação são chaves opcionais.
    Itens devem ser únicos na combinação Partition+Range.
    Os itens com a mesma Partition são armazenados em partições próximas, 
    ordenados pelo valor da chave de Ordenação

### PK + SK
    Um nome comum para se dar aos campos que vão compor
    as chaves do dynamo é (PK) = Partition (SK) = Sort


# Importante
    Ao modelar a informação, sempre ter em mente que para fazer queries 
    por chave ou índice, a chave de partição é 

# OBRIGATÓRIA!!!
###### (ou algum índice que se tenha a partition key em mãos :D) 

#Exemplo de modelagem

- Definir padrões de acesso 
- Como quero acessar meu dado?
- Quais índices serão necessários para atender esses padrões de acesso
- Reutilização de índices entre entidades (GSI overloading)
- Evitar Hot Partition
- Evitar criar índices demais
- Definir projeções (KEYS_ONLY, INCLUDE, ALL)
- Via de regra, Scan não deve ser uma opção para operações transacionais. 
  - Usar com muito cuidado pois gera custos...
  - Que serão extremamente elevados no tamanho das tabelas do Meli
  - É uma opção caso seja necessária fazer alguma migração / alteração de dados em massa    
  
  
#### Modelagem simples somente com uma PK + 1 indice:

|          | PK                | gsi1pk            | gsi1sk           |
|----------|-------------------|-------------------|------------------|
| Order    | ORDER#{ORDERID}   | CUSTOMER#{CUSTID} | ORDERDATE#{date} |
| Customer | CUSTOMER#{CUSTID} |                   |                  |

<br>

#### Modelagem clearing !


# Código !!! <br>
