#!/usr/bin/env bash
# Observer: runs alongside ./stress.sh collecting independent per-layer signals.
#
# Captures, every POLL_INTERVAL seconds, into logs/observe-<TS>.tsv:
#   - docker stats (per-container CPU / mem / net / block)
#   - pg_stat_activity state histogram
#   - pprof goroutine count (whichever API container bound :6060 first)
#
# Also captures a single CPU pprof snapshot PPROF_DELAY seconds after start
# (duration PPROF_SECONDS) to logs/cpu-<TS>.pb.gz.
#
# Usage:
#   ./scripts/observe.sh                       # defaults: poll 5s, pprof at 30s for 30s
#   POLL_INTERVAL=2 PPROF_DELAY=60 PPROF_SECONDS=30 ./scripts/observe.sh
#
# Stop with Ctrl-C; the observed data is flushed on every line.

set -u

POLL_INTERVAL="${POLL_INTERVAL:-5}"
PPROF_DELAY="${PPROF_DELAY:-30}"
PPROF_SECONDS="${PPROF_SECONDS:-30}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$ROOT_DIR/logs"
mkdir -p "$LOG_DIR"

TS="$(date +%Y%m%dT%H%M%S)"
OBS_FILE="$LOG_DIR/observe-${TS}.tsv"
PPROF_FILE="$LOG_DIR/cpu-${TS}.pb.gz"

PG_CONTAINER="db"
PG_USER="${POSTGRES_USER:-bololo}"
PG_DB="${POSTGRES_DB:-shakeit}"

# header
{
  printf 'ts\t'
  printf 'api1_cpu\tapi1_mem\tapi2_cpu\tapi2_mem\tdb_cpu\tdb_mem\tnginx_cpu\tnginx_mem\t'
  printf 'pg_active\tpg_idle\tpg_idle_in_txn\tpg_other\t'
  printf 'goroutines\n'
} >> "$OBS_FILE"

echo "observer: writing to $OBS_FILE"
echo "observer: cpu pprof will be captured ${PPROF_DELAY}s after start (${PPROF_SECONDS}s duration) -> $PPROF_FILE"

# Schedule the one-shot pprof capture in the background.
(
  sleep "$PPROF_DELAY"
  echo "observer: capturing ${PPROF_SECONDS}s CPU pprof..."
  curl -s --max-time $((PPROF_SECONDS + 10)) \
    -o "$PPROF_FILE" \
    "http://localhost:6060/debug/pprof/profile?seconds=${PPROF_SECONDS}"
  if [ -s "$PPROF_FILE" ]; then
    echo "observer: pprof saved -> $PPROF_FILE (analyze with: go tool pprof $PPROF_FILE)"
  else
    echo "observer: WARN pprof empty or curl failed (is any api container bound to :6060?)"
  fi
) &

cleanup() {
  trap '' INT TERM
  echo "observer: stopping (waiting for any in-flight pprof)..."
  wait 2>/dev/null
  echo "observer: done. data in $OBS_FILE"
  trap - INT TERM
  exit 0
}
trap cleanup INT TERM

echo "observer: polling every ${POLL_INTERVAL}s. Ctrl-C to stop."

while :; do
  NOW="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

  # ---- docker stats (one-shot, parseable) ----
  STATS="$(docker stats --no-stream --format \
    '{{.CPUPerc}}\t{{.MemUsage}}' api1 api2 db nginx 2>/dev/null)"
  # Re-split into fields; MemUsage looks like "12.34MiB / 512MiB"; keep the used part only.
  STATS_CLEAN="$(printf '%s' "$STATS" | awk -F'\t' '
    {
      cpu = $1; sub(/%/, "", cpu);
      mem = $2; sub(/[[:space:]].*/, "", mem);
      printf "%s\t%s", cpu, mem;
      if (NR < 4) printf "\t";
    }
  ')"

  # ---- pg_stat_activity states ----
  PG="$(docker exec "$PG_CONTAINER" \
    psql -U "$PG_USER" -d "$PG_DB" -t -A -F= -c \
    "SELECT state, count(*) FROM pg_stat_activity GROUP BY state ORDER BY state;" 2>/dev/null)"
  # Reshape to fixed columns: active, idle, idle in transaction, other.
  pg_active=0; pg_idle=0; pg_idle_txn=0; pg_other=0
  while IFS='=' read -r st cnt; do
    [ -z "${st:-}" ] && continue
    case "$st" in
      active)                  pg_active="$cnt" ;;
      idle)                    pg_idle="$cnt" ;;
      "idle in transaction")   pg_idle_txn="$cnt" ;;
      *)                       pg_other="$((pg_other + cnt))" ;;
    esac
  done <<< "$PG"

  # ---- goroutine count (pprof) ----
  GORO="$(curl -s --max-time 2 'http://localhost:6060/debug/pprof/goroutine?debug=1' \
    | head -1 | sed -E 's/goroutine profile: total ([0-9]+).*/\1/')"
  [ -z "$GORO" ] && GORO="-"

  printf '%s\t%s\t%s\t%s\t%s\t%s\n' \
    "$NOW" "$STATS_CLEAN" "$pg_active" "$pg_idle" "$pg_idle_txn" "$pg_other" "$GORO" \
    >> "$OBS_FILE"

  # also echo a compact line to stdout for live watching
  printf '%s  api1=%s%%  db=%s%%  nginx=%s%%  pg_active=%s  goroutines=%s\n' \
    "$NOW" "$(printf '%s' "$STATS_CLEAN" | cut -f1)" \
    "$(printf '%s' "$STATS_CLEAN" | cut -f5)" \
    "$(printf '%s' "$STATS_CLEAN" | cut -f7)" \
    "$pg_active" "$GORO"

  sleep "$POLL_INTERVAL"
done
