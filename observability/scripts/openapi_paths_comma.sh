#!/usr/bin/env sh
# Печатает все path из OpenAPI composition-api для custom-переменной Grafana (через запятую).
# Запуск из корня репозитория: ./med_ml/observability/scripts/openapi_paths_comma.sh
# или из med_ml: ./observability/scripts/openapi_paths_comma.sh

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
SPEC="${ROOT}/composition-api/cmd/service/server.yml"
exec grep -E '^  /[^:]+:' "$SPEC" | sed 's/^  //;s/:$//' | paste -sd, -
