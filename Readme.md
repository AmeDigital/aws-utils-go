# aws-utils-go
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

Quem faz este deploy é o `go get stash.b2w/asp/aws-utils-go.git`. Acontece que este comando irá tentar fazer o download via https, e o nosso stash não suporta https, ele suporta ssh.

O `go get` utiliza o git client para fazer o download, portanto precisamos configurar o git client para usar ssh ao falar com o host *stash.b2w*.

Para configurar o git client, rode `vi ~/.gitconfig` e acrescente no final do arquivo as linhas abaixo:

```
[url "git@stash.b2w:"]
	insteadOf = https://stash.b2w/
```

Feito isto, faça a instalação da lib rodando o `go get stash.b2w/asp/aws-utils-go.git`

#### Importar o aws-utils-go no seu código

Declare o import da lib como no exemplo abaixo:

```golang
package main

import (
    "stash.b2w/asp/aws-utils-go.git/dynamodbutils"
)
...
// save to the "Cities" table an instance of the "City" struct
err := dynamodbutils.PutItem("Cities", city)
...
```

## Como extender o aws-utils-go

Se quiser extender o aws-utils-go o clone do projeto obrigatoriamente tem que ser feito no diretorio `$GOPATH/src/stash.b2w/asp/aws-utils-go.git`.

Isto é porque o próprio codigo do aws-utils-go, quando faz import de um pacote do mesmo projeto, utiliza no importe do pacote o prefixo `stash.b2w/asp/aws-utils-go.git`.

Para fazer o clone, use os comandos:

```shell
mkdir -p $GOPATH/src/stash.b2w/asp/  
cd $GOPATH/src/stash.b2w/asp/  
git clone ssh://git@stash.b2w/asp/aws-utils-go.git aws-utils-go.git  
```

## Gerar uma imagem docker com a lib 'aws-utils-go' embedada e publicar no Nexus B2W

Para que seu código que utilizou 'aws-utils-go' possa ser buildado no bamboo é preciso criar uma imagem docker para golang 
contendo esta lib deployada na GOPATH.  
Este projeto vem com um Dockerfile que cria esta imagem.

O script *build-docker-image.sh* builda a imagem com o nome de *b2wbuild/golang-aws-utils-go:TAG* e publica a mesma no repositório da B2W, *registry.b2w.io*.

Para buildar e publicar uma nova imagem, rode o script como no exemplo:

```shell
./build-docker-image.sh "1.0"
```

Onde "1.0" é uma tag name que irá identificar a versão da lib aws-utils-go que foi empacotada. A medida que esta lib for atualizada, novas tags deverão ser publicadas no repósitório.

O script *build-docker-image.sh* não permite sobrescrever uma tag existente e irá dar erro se já houver uma imagem com a mesma tag no repositório.

## Como buildar o seu código pelo bamboo usando uma imagem docker

#### Criar o 'build plan' usando a imagem gerada

No bamboo, quando for buildar o seu projeto, utilize a imagem docker de GO com a aws-utils-go embedada (a que foi criada no passo acima). 

Veja o exemplo abaixo, o comando usado no "build plan" do bamboo usa a imagem *registry.b2w.io/b2wbuild/golang-aws-utils-go:1.0* para buildar a app:

```shell
run --volume ${bamboo.build.working.directory}/ame-iot-auth:/go/src/ame-iot-auth --workdir /go/src/ame-iot-auth --rm registry.b2w.io/b2wbuild/golang-aws-utils-go:1.0 /bin/bash -c ./device-api/build.sh
```


