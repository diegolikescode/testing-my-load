#!/usr/bin/bash

# docker compose down --remove-orphans
# docker system prune -f

RESULTS_DIR="$(pwd)/load-test/resultados"

mkdir -p $RESULTS_DIR

echo "iniciando e logando execução da API"
docker compose up -d --build
docker compose logs > "docker-compose.logs"
echo "pausa de 6 segundos pro start da API"
sleep 6
echo "iniciando teste"

cd load-test
npx gatling run --simulation rinhabackend

cd ..
echo "teste finalizado"

# Archive the Gatling HTML report alongside contagem-pessoas.
TS="$(date +%Y%m%dT%H%M%S)"
GATLING_RUN_DIR="$(ls -dt load-test/target/gatling/*/ 2>/dev/null | head -1)"
if [ -n "$GATLING_RUN_DIR" ]; then
  GATLING_DEST="$RESULTS_DIR/gatling-${TS}"
  cp -r "$GATLING_RUN_DIR" "$GATLING_DEST"
  echo "relatório gatling salvo em $GATLING_DEST"
else
  echo "WARN: nenhum diretório de relatório gatling encontrado em load-test/target/gatling/"
fi

echo "fazendo request e salvando a contagem de pessoas"
SAVE_CONTAGEM="$RESULTS_DIR/contagem-pessoas-$(date)"
curl -v "http://localhost:9999/contagem-pessoas" > "$SAVE_CONTAGEM"
echo "resultado da contagem em $SAVE_CONTAGEM"
echo "INSERTS: $(cat "$SAVE_CONTAGEM")"
# echo "cleaning up do docker"
# docker compose rm -s -f
# docker compose down --remove-orphans api1 api2 nginx
