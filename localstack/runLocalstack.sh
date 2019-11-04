#!/bin/bash

function die {
    declare MSG="$@"
    echo -e "\n$0: Error: $MSG">&2
    exit 1
}

##########################
# iniciando o localstack
############################

which localstack > /dev/null || die "localstack não está instalado. Instale-o com 'pip install localstack'."

# A VARIAVEL DE AMBIENTE '$SERVICES' DEFINE QUAIS SERVIÇOS DO LOCALSTACK DEVEM SER INICIADOS
# ver: https://github.com/localstack/localstack e https://docs.aws.amazon.com/cli/latest/reference/#available-services
# ex: SERVICES="dynamodb,kinesis,s3,sns,sqs"

[ -z "$SERVICES" ] && die "A variável de ambiente SERVICES não pode ser vazia. Quais serviços vc quer rodar?"

# inicia o localstack em background, 
# direciona o output para um arquivo de log
# bloqueia a thread até que apareça a string "Ready." no log, com um timeout de $TIMEOUT segundos
# retorna uma string com o PID do localstack

OUTPUT_FILE=/tmp/localstack-start.log
FORCE_NONINTERACTIVE=true localstack start > $OUTPUT_FILE 2>&1 &
LOCALSTACK_PID="$!"
TIMEOUT=60
while ! grep -Eq "^Ready\.$" $OUTPUT_FILE; do
    ((TIMEOUT > 0)) || die "Timeout ao iniciar o localstack. \
        Favor tentar inicia-lo na mão com o comando 'localstack start'\n.
        $(cat $OUTPUT_FILE >&2)"

    ps $LOCALSTACK_PID > /dev/null || die "O processo do localstack $LOCALSTACK_PID não está mais rodando.\n \
        $(cat $OUTPUT_FILE >&2)"

    sleep 1

    ((TIMEOUT--))
done
rm $OUTPUT_FILE


[ -z "$LOCALSTACK_PID" ] && die "Bug: variavel LOCALSTACK_PID não pode estar vazia!"
echo -n $LOCALSTACK_PID