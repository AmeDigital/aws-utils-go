#!/bin/bash

function die {
    declare MSG="$@"
    echo "$0: Error: $MSG">&2
    exit 1
}

# inicia o localstack em background, 
# direciona o output para um arquivo de log
# bloqueia a thread até que apareça a string "Ready." no log, com um timeout de $TIMEOUT segundos
# retorna uma string com o PID do localstack
function start_localstack {        
    declare OUTPUT_FILE=/tmp/localstack-start.log
    localstack start > $OUTPUT_FILE 2>&1 &
    declare LOCALSTACK_PID="$!"
    declare TIMEOUT=15
    while ! grep -Eq "^Ready\.$" $OUTPUT_FILE; do
        ((TIMEOUT > 0)) || die "Timeout ao iniciar o localstack. \
            Favor tentar inicia-lo na mão com o comando 'localstack start'.
            $(cat $OUTPUT_FILE)"
        sleep 1
        ((TIMEOUT--))
    done
    rm $OUTPUT_FILE
    echo $LOCALSTACK_PID
}

##########################
# iniciando o localstack
############################

which localstack > /dev/null || die "localstack não está instalado. Instale-o com 'pip install localstack'."

# A VARIAVEL DE AMBIENTE '$SERVICES' DEFINE QUAIS SERVIÇOS DO LOCALSTACK DEVEM SER INICIADOS
# ver: https://github.com/localstack/localstack e https://docs.aws.amazon.com/cli/latest/reference/#available-services
# ex: SERVICES="dynamodb,kinesis,s3,sns,sqs"

[ -z "$SERVICES" ] && die "A variável de ambiente SERVICES não pode ser vazia. Quais serviços vc quer rodar?"

# Iniciando o localstack... 
LOCALSTACK_PID=$(start_localstack)
[ -z $LOCALSTACK_PID ] && die "Bug: variavel LOCALSTACK_PID não pode estar vazia!"
echo -n $LOCALSTACK_PID