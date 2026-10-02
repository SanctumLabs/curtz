#!/bin/bash
# Elastic setup job. Idempotent: safe to run on every `up`.
#   ELK_MODE=ha   first generates a CA and the node certificates into the shared certs volume
#   both modes    then wait for Elasticsearch and provision users, roles, the ILM policy and the index template
# Env: ELK_MODE (ha|single), ES_URL, ELASTIC_PASSWORD, KIBANA_SYSTEM_PASSWORD, LOGSTASH_WRITER_PASSWORD,
#      GRAFANA_READER_PASSWORD, METRICS_READER_PASSWORD
set -euo pipefail

ES="${ES_URL:?}"

if [ "${ELK_MODE:?}" = ha ]; then
  CERTS=/usr/share/elasticsearch/config/certs
  if [ ! -f "$CERTS/ca/ca.crt" ]; then
    echo "generating the certificate authority"
    bin/elasticsearch-certutil ca --silent --pem --out "$CERTS/ca.zip"
    unzip -q -o "$CERTS/ca.zip" -d "$CERTS"
  fi
  if [ ! -f "$CERTS/es-3/es-3.crt" ]; then
    echo "generating the node certificates"
    cat >"$CERTS/instances.yml" <<'YAML'
instances:
  - name: es-1
    dns: [es-1, localhost]
    ip: [127.0.0.1]
  - name: es-2
    dns: [es-2, localhost]
    ip: [127.0.0.1]
  - name: es-3
    dns: [es-3, localhost]
    ip: [127.0.0.1]
YAML
    bin/elasticsearch-certutil cert --silent --pem --in "$CERTS/instances.yml" \
      --ca-cert "$CERTS/ca/ca.crt" --ca-key "$CERTS/ca/ca.key" --out "$CERTS/certs.zip"
    unzip -q -o "$CERTS/certs.zip" -d "$CERTS"
  fi
  # Development only: the keys are world-readable so the unprivileged elasticsearch user can read them.
  chown -R root:root "$CERTS"
  find "$CERTS" -type d -exec chmod 755 {} +
  find "$CERTS" -type f -exec chmod 644 {} +
  touch "$CERTS/.ready"
fi

es() { curl -sS --fail-with-body -u "elastic:${ELASTIC_PASSWORD}" -H 'Content-Type: application/json' "$@"; }

echo "waiting for Elasticsearch at $ES"
until [ "$(curl -s -o /dev/null -w '%{http_code}' -u "elastic:${ELASTIC_PASSWORD}" "$ES/_cluster/health?wait_for_status=yellow&timeout=5s")" = 200 ]; do
  sleep 3
done

REPLICAS=0
[ "$ELK_MODE" = ha ] && REPLICAS=1

echo "provisioning roles and users"
es -X POST "$ES/_security/user/kibana_system/_password" -d "{\"password\":\"${KIBANA_SYSTEM_PASSWORD}\"}" >/dev/null

es -X PUT "$ES/_security/role/logstash_writer" -d '{
  "cluster": ["monitor"],
  "indices": [{"names": ["logs-curtz-*"], "privileges": ["create_doc", "auto_configure", "view_index_metadata"]}]
}' >/dev/null
es -X PUT "$ES/_security/role/grafana_reader" -d '{
  "indices": [{"names": ["logs-curtz-*"], "privileges": ["read", "view_index_metadata"]}]
}' >/dev/null
es -X PUT "$ES/_security/role/metrics_reader" -d '{
  "cluster": ["monitor"],
  "indices": [{"names": ["*"], "privileges": ["monitor"]}]
}' >/dev/null

es -X PUT "$ES/_security/user/logstash_writer" -d "{\"password\":\"${LOGSTASH_WRITER_PASSWORD}\",\"roles\":[\"logstash_writer\"]}" >/dev/null
es -X PUT "$ES/_security/user/grafana_reader" -d "{\"password\":\"${GRAFANA_READER_PASSWORD}\",\"roles\":[\"grafana_reader\"]}" >/dev/null
es -X PUT "$ES/_security/user/metrics_reader" -d "{\"password\":\"${METRICS_READER_PASSWORD}\",\"roles\":[\"metrics_reader\"]}" >/dev/null

echo "installing the ILM policy and the index template"
es -X PUT "$ES/_ilm/policy/curtz-logs" -d '{
  "policy": {"phases": {
    "hot": {"actions": {"rollover": {"max_age": "1d", "max_primary_shard_size": "5gb"}}},
    "delete": {"min_age": "7d", "actions": {"delete": {}}}
  }}
}' >/dev/null

es -X PUT "$ES/_index_template/logs-curtz" -d "{
  \"index_patterns\": [\"logs-curtz-*\"],
  \"data_stream\": {},
  \"priority\": 500,
  \"template\": {
    \"settings\": {\"index.lifecycle.name\": \"curtz-logs\", \"index.number_of_replicas\": ${REPLICAS}},
    \"mappings\": {\"properties\": {
      \"@timestamp\": {\"type\": \"date\"},
      \"message\": {\"type\": \"text\"},
      \"service\": {\"properties\": {\"name\": {\"type\": \"keyword\"}}},
      \"trace\": {\"properties\": {\"id\": {\"type\": \"keyword\"}}},
      \"span\": {\"properties\": {\"id\": {\"type\": \"keyword\"}}},
      \"log\": {\"properties\": {\"level\": {\"type\": \"keyword\"}}}
    }}
  }
}" >/dev/null

echo "elk setup complete"
