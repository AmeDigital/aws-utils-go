# aws-utils-go
## Utilitários com interface simples para acessar serviços da AWS

A aws disponibiliza para a linguagem go a sdk aws-sdk-go, porém sua interface pode ser bem complicada.

Esta lib oferece funções com interfaces simples para executar tarefas comuns.

## Pacotes

* dynamodbutils: oferece interfaces simplificadas para as ações PutItem, GetItem, UpdateItem
* s3utils: oferece GetObject, GetObjectAsString, ListObjects
* snsutils: oferece SendMessage, SendMessageWithAttributes
* localstack: utilitários para iniciar/parar o localstack e seus serviços na máquina local. Está *experimental* ainda e sua interface deve mudar.

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


// save to the "Cities" table an instance of a "City" struct

city := City{
    Id:         1,
    Name:       "New York",
}

err := PutItem("Cities", city)
if err !=nil {
    panic(err)
}
```

É necessario rodar o "go get" para fazer download do aws-utils-go para a sua máquina de desenvolvimento.

## Como buildar seu código no bamboo

### Gerar uma imagem docker com a lib embedada

Para que seu código possa buildar no bamboo é preciso criar uma imagem golang que tenha esta lib deployada lá dentro.

O script *build-docker-image.sh* builda a imagem e publica a mesma para o repositório da B2W.

Para utilizar o script rode como o exemplo abaixo:

```bash
    ./build-docker-image.sh "1.0"
```

Onde "1.0" é uma tag name que irá identificar a versão da lib aws-utils-go que foi empacotada.

O script *build-docker-image.sh* não permite sobrescrever uma tag existente e irá dar erro se já houver uma imagem com a mesma tag no repositório.

### Criar o 'build plan' usando a imagem gerada

No bamboo, quando for buildar o seu projeto, utilize a imagem gerada como no exemplo abaixo:

```bash
run --volume ${bamboo.build.working.directory}/ame-iot-auth:/go/src/ame-iot-auth --workdir /go/src/ame-iot-auth --rm registry.b2w.io/b2wbuild/golang-aws-utils-go:1.0 /bin/bash -c ./device-api/build.sh
```


