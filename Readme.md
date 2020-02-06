# aws-utils-go    IGNORE ESSA MENSAGEM EM CAPS, ESSA MENSAGEM NAO SAO OS DROIDS QUE VC ESTA PROCURANDO
## Utilitários com interface simples para acessar serviços da AWS

A aws disponibiliza para a linguagem go a sdk aws-sdk-go, porém sua interface pode ser bem complicada.

Esta lib oferece funções com interfaces simples para executar tarefas comuns.

## Pacotes

* dynamodbutils: oferece interfaces simplificadas para as ações PutItem, GetItem, UpdateItem, PutItemWithConditional, FindOneFromIndex
* s3utils: oferece GetObject, GetObjectAsString, ListObjects, PutObject.
* snsutils: oferece SendMessage, SendMessageWithAttributes.
* sqsutils: oferece SendMessage, ReadMessage, DeleteMessage, GetMessageAttribute
* sessionutils: permite configurar a Session (aws-sdk-go/aws/session) que será utilizada pelos utils para se comunicarem com a AWS.
* localstack (**experimental**): utilitários para iniciar/parar o localstack e seus serviços na máquina local. Está *experimental* ainda e sua interface deve mudar.

## Como importar e utilizar o código

#### Instalar o aws-utils-go no seu ambiente de desenvolvimento

Para ser utilizada, a lib aws-utils-go precisa estar deployada em seu diretorio $GOPATH/src. 

Faça a instalação da lib rodando o `go get github.com/AmeDigital/aws-utils-go`

#### Importar o aws-utils-go no seu código

Declare o import da lib como no exemplo abaixo:

```golang
package main

import (
    "github.com/AmeDigital/aws-utils-go/dynamodbutils"
)
...
// save to the "Cities" table an instance of the "City" struct
err := dynamodbutils.PutItem("Cities", city)
...
```

## Como extender o aws-utils-go

Se quiser extender o aws-utils-go o clone do projeto obrigatoriamente tem que ser feito no diretorio `$GOPATH/src/github.com/AmeDigital/aws-utils-go`.

Isto é porque o próprio codigo do aws-utils-go, quando faz import de um pacote do mesmo projeto, utiliza no importe do pacote o prefixo `github.com/AmeDigital/aws-utils-go`.

Para fazer o clone, use os comandos:

```shell
mkdir -p $GOPATH/src/github.com/AmeDigital && cd $GOPATH/src/github.com/AmeDigital  
git clone ssh://git@github.com/AmeDigital/aws-utils-go aws-utils-go.git  
```
