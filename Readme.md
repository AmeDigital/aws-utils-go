# aws-utils-go
## Utilitários com interface simples para acessar serviços da AWS

A aws disponibiliza para a linguagem go a sdk aws-sdk-go, porém sua interface pode ser bem complicada.

Esta lib oferece funções com interfaces simples para executar tarefas comuns.

## Pacotes

* dynamodbutils: oferece interfaces simplificadas para as ações PutItem, GetItem, UpdateItem.
* s3utils: oferece GetObject, GetObjectAsString, ListObjects.
* snsutils: oferece SendMessage, SendMessageWithAttributes.
* sessionutils: permite configurar a Session (aws-sdk-go/aws/session) que será utilizada pelos utils para se comunicarem com a AWS.
* localstack (**experimental**): utilitários para iniciar/parar o localstack e seus serviços na máquina local. Está *experimental* ainda e sua interface deve mudar.

## Como importar e utilizar o código

Esta lib está publicada no stash. Para importá-la e utilizá-la em seu código siga o exemplo abaixo:

```golang
package main

import (
    "stash.b2w/asp/aws-utils-go.git/dynamodbutils"
)

type City struct {
    Id         int
    Name       string
}
city := City{
    Id:         1,
    Name:       "New York",
}

// save to the "Cities" table an instance of the "City" struct
err := PutItem("Cities", city)

if err !=nil {
    panic(err)
}
```

É necessario rodar o "go get" para fazer download do aws-utils-go para a sua máquina de desenvolvimento.

## Como gerar uma imagem docker com a lib 'aws-utils-go' embedada

Para que seu código que utilizou 'aws-utils-go' possa ser buildado no bamboo é preciso criar uma imagem docker para golang 
contendo esta lib deployada na GOPATH.  
Este projeto vem com um Dockerfile que cria esta imagem.

O script *build-docker-image.sh* builda a imagem com o nome de *b2wbuild/golang-aws-utils-go:TAG* e publica a mesma no repositório da B2W, *registry.b2w.io*.

Para buildar e publicar uma nova imagem, rode o script como no exemplo:

```bash
    ./build-docker-image.sh "1.0"
```

Onde "1.0" é uma tag name que irá identificar a versão da lib aws-utils-go que foi empacotada. A medida que esta lib for atualizada, novas tags deverão ser publicadas no repósitório.

O script *build-docker-image.sh* não permite sobrescrever uma tag existente e irá dar erro se já houver uma imagem com a mesma tag no repositório.

## Como buildar o seu código pelo bamboo usando uma imagem docker

### Criar o 'build plan' usando a imagem gerada

No bamboo, quando for buildar o seu projeto, utilize a imagem gerada como no exemplo abaixo:

```bash
run --volume ${bamboo.build.working.directory}/ame-iot-auth:/go/src/ame-iot-auth --workdir /go/src/ame-iot-auth --rm registry.b2w.io/b2wbuild/golang-aws-utils-go:1.0 /bin/bash -c ./device-api/build.sh
```


