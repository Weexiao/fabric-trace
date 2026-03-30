#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NETWORK_DIR="$ROOT_DIR/blockchain/network"
TAPE_DIR="$ROOT_DIR/blockchain/tape"
APP_DIR="$ROOT_DIR/application"
BACKEND_DIR="$ROOT_DIR/application/backend"
RUNS_BASE_DIR="$ROOT_DIR/perf/runs"

RUN_ID="$(date +%Y%m%d_%H%M%S)"
SCENARIOS="invoke,query"
N_LEVELS="100,300,500"
ROUNDS=1
SEED_N=1
DO_SEED=1
START_NETWORK=0
STOP_NETWORK=0
BACKEND_MODE="none" # none | local | docker
BACKEND_URL="http://127.0.0.1:9090"
WAIT_BACKEND_SECONDS=60
COLLECT_DOCKER_STATS=1
DRY_RUN=0

# Stability-oriented defaults for avoiding gRPC GOAWAY too_many_pings.
TAPE_NUM_OF_CONN=4
TAPE_CLIENT_PER_CONN=2
INTER_CASE_SLEEP_SECONDS=3
RETRY_ON_GOAWAY=1
GOAWAY_MAX_RETRIES=2
GOAWAY_RETRY_BACKOFF_SECONDS=15

usage() {
  cat <<'EOF'
Usage:
  ./perf/run_perf.sh [options]

Options:
  --run-id <id>                 Override run id (default: timestamp)
  --scenarios <list>            Comma list: register,invoke,query (default: invoke,query)
  --n-levels <list>             Comma list of -n values (default: 100,300,500)
  --rounds <num>                Repeat count per n-level (default: 1)
  --seed-n <num>                Registration tx count for seeding (default: 1)
  --skip-seed                   Skip seed stage
  --start-network               Run blockchain/network/start.sh first (destructive)
  --stop-network                Run blockchain/network/stop.sh at end
  --backend-mode <mode>         none | local | docker (default: none)
  --backend-url <url>           Backend base URL for healthcheck (default: http://127.0.0.1:9090)
  --wait-backend-seconds <num>  Healthcheck timeout (default: 60)
  --collect-docker-stats <0|1>  Capture docker stats per round (default: 1)
  --tape-num-of-conn <num>      Override tape num_of_conn in runtime configs (default: 4)
  --tape-client-per-conn <num>  Override tape client_per_conn in runtime configs (default: 2)
  --inter-case-sleep-seconds <num>
                                Sleep between rounds/cases to reduce connection bursts (default: 3)
  --retry-on-goaway <0|1>       Retry on gRPC GOAWAY too_many_pings (default: 1)
  --goaway-max-retries <num>    Max retries per case on GOAWAY (default: 2)
  --goaway-retry-backoff-seconds <num>
                                Backoff seconds before GOAWAY retry (default: 15)
  --dry-run                     Print commands only
  --help                        Show this help

Notes:
  - start.sh internally clears old network and chain data via stop.sh.
  - tape command uses: ./tape -c <config> -n <num> run
EOF
}

log() {
  echo "[$(date '+%F %T')] $*"
}

run_cmd() {
  local cmd="$1"
  log "+ $cmd" | tee -a "$RUN_LOG"
  if [[ "$DRY_RUN" -eq 1 ]]; then
    return 0
  fi
  bash -lc "$cmd" >>"$RUN_LOG" 2>&1
}

wait_for_backend() {
  local url="$1"
  local timeout="$2"
  local start
  start=$(date +%s)

  while true; do
    if curl -sS -m 2 "$url/ping" >/dev/null 2>&1; then
      log "Backend is ready: $url/ping" | tee -a "$RUN_LOG"
      return 0
    fi
    if (( $(date +%s) - start >= timeout )); then
      log "Backend healthcheck timeout after ${timeout}s: $url/ping" | tee -a "$RUN_LOG"
      return 1
    fi
    sleep 2
  done
}

parse_csv_to_array() {
  local raw="$1"
  local -n out_ref=$2
  IFS=',' read -r -a out_ref <<<"$raw"
}

config_for_scene() {
  case "$1" in
    register) echo "config_register.yaml" ;;
    invoke) echo "config_invoke.yaml" ;;
    query) echo "config_query.yaml" ;;
    *)
      echo "unsupported scene: $1" >&2
      return 1
      ;;
  esac
}

prepare_tape_config() {
  local scene="$1"
  local source_config
  source_config=$(config_for_scene "$scene")

  local src="$TAPE_DIR/$source_config"
  local dst_dir="$RUN_DIR/meta/runtime_configs"
  mkdir -p "$dst_dir"
  local dst="$dst_dir/${scene}.yaml"

  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "$source_config"
    return 0
  fi

  cp "$src" "$dst"
  sed -E -i "s/^num_of_conn:.*/num_of_conn: ${TAPE_NUM_OF_CONN}/" "$dst"
  sed -E -i "s/^client_per_conn:.*/client_per_conn: ${TAPE_CLIENT_PER_CONN}/" "$dst"
  echo "$dst"
}

collect_docker_stats() {
  local scene="$1"
  local n="$2"
  local round="$3"

  if [[ "$COLLECT_DOCKER_STATS" -ne 1 ]]; then
    return 0
  fi

  local out="$RUN_DIR/raw/docker_stats_${scene}_n${n}_r${round}.csv"
  if [[ "$DRY_RUN" -eq 1 ]]; then
    log "+ docker stats --no-stream > $out" | tee -a "$RUN_LOG"
    return 0
  fi

  if ! command -v docker >/dev/null 2>&1; then
    log "docker command not found, skip docker stats" | tee -a "$RUN_LOG"
    return 0
  fi

  docker stats --no-stream --format "{{.Name}},{{.CPUPerc}},{{.MemUsage}},{{.NetIO}},{{.BlockIO}},{{.PIDs}}" >"$out" || true
}

run_tape_case() {
  local scene="$1"
  local n="$2"
  local round="$3"
  local config
  config=$(prepare_tape_config "$scene")

  local log_file="$RUN_DIR/raw/tape_${scene}_n${n}_r${round}.log"
  local started_at ended_at duration rc

  started_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  local started_s
  started_s=$(date +%s)

  local cmd="cd \"$TAPE_DIR\" && ./tape -c \"$config\" -n \"$n\" run"
  log "+ $cmd" | tee -a "$RUN_LOG"

  local attempt=1
  local max_attempts=1
  if [[ "$RETRY_ON_GOAWAY" -eq 1 ]]; then
    max_attempts=$((GOAWAY_MAX_RETRIES + 1))
  fi

  if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "[dry-run] $cmd" >"$log_file"
    rc=0
  else
    while true; do
      set +e
      bash -lc "$cmd" >"$log_file" 2>&1
      rc=$?
      set -e

      if [[ "$rc" -eq 0 ]]; then
        break
      fi

      if [[ "$attempt" -ge "$max_attempts" ]]; then
        break
      fi

      if grep -Eq "ENHANCE_YOUR_CALM|too_many_pings|GoAway" "$log_file"; then
        log "scene=${scene} n=${n} round=${round} got GOAWAY too_many_pings; retry ${attempt}/${GOAWAY_MAX_RETRIES} after ${GOAWAY_RETRY_BACKOFF_SECONDS}s" | tee -a "$RUN_LOG"
        sleep "$GOAWAY_RETRY_BACKOFF_SECONDS"
        attempt=$((attempt + 1))
        continue
      fi

      break
    done
  fi

  ended_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  duration=$(( $(date +%s) - started_s ))

  echo "${scene},${n},${round},${started_at},${ended_at},${duration},${rc},${log_file}" >>"$SUMMARY_CSV"

  local metrics_csv="$RUN_DIR/summary/tape_metrics.csv"
  if [[ "$DRY_RUN" -eq 1 ]]; then
    if [[ ! -f "$metrics_csv" ]]; then
      echo "scene,n,round,tps,success,fail,p95_ms,p99_ms,input" >"$metrics_csv"
    fi
    echo "${scene},${n},${round},NA,NA,NA,NA,NA,${log_file}" >>"$metrics_csv"
  else
    python3 "$ROOT_DIR/perf/extract_tape_metrics.py" \
      --input "$log_file" \
      --scene "$scene" \
      --n "$n" \
      --round "$round" \
      --output "$metrics_csv" \
      --append >>"$RUN_LOG" 2>&1 || true
  fi

  collect_docker_stats "$scene" "$n" "$round"

  if [[ "$rc" -ne 0 ]]; then
    log "scene=${scene} n=${n} round=${round} failed (rc=${rc})" | tee -a "$RUN_LOG"
    return "$rc"
  fi

  return 0
}

start_backend_if_needed() {
  case "$BACKEND_MODE" in
    none)
      log "Skip backend startup (backend-mode=none)" | tee -a "$RUN_LOG"
      ;;
    local)
      run_cmd "cd \"$BACKEND_DIR\" && nohup go run main.go > \"$RUN_DIR/logs/backend.log\" 2>&1 & echo \\\$! > \"$RUN_DIR/meta/backend.pid\""
      wait_for_backend "$BACKEND_URL" "$WAIT_BACKEND_SECONDS"
      ;;
    docker)
      run_cmd "cd \"$APP_DIR\" && ./start_docker.sh"
      wait_for_backend "$BACKEND_URL" "$WAIT_BACKEND_SECONDS"
      ;;
    *)
      echo "invalid backend-mode: $BACKEND_MODE" >&2
      exit 1
      ;;
  esac
}

stop_backend_if_needed() {
  case "$BACKEND_MODE" in
    local)
      local pid_file="$RUN_DIR/meta/backend.pid"
      if [[ -f "$pid_file" ]]; then
        local pid
        pid=$(cat "$pid_file")
        if [[ -n "$pid" ]]; then
          run_cmd "kill $pid || true"
        fi
      fi
      ;;
    docker)
      run_cmd "cd \"$APP_DIR\" && ./stop_docker.sh"
      ;;
    none)
      ;;
  esac
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --run-id)
      RUN_ID="$2"
      shift 2
      ;;
    --scenarios)
      SCENARIOS="$2"
      shift 2
      ;;
    --n-levels)
      N_LEVELS="$2"
      shift 2
      ;;
    --rounds)
      ROUNDS="$2"
      shift 2
      ;;
    --seed-n)
      SEED_N="$2"
      shift 2
      ;;
    --skip-seed)
      DO_SEED=0
      shift
      ;;
    --start-network)
      START_NETWORK=1
      shift
      ;;
    --stop-network)
      STOP_NETWORK=1
      shift
      ;;
    --backend-mode)
      BACKEND_MODE="$2"
      shift 2
      ;;
    --backend-url)
      BACKEND_URL="$2"
      shift 2
      ;;
    --wait-backend-seconds)
      WAIT_BACKEND_SECONDS="$2"
      shift 2
      ;;
    --collect-docker-stats)
      COLLECT_DOCKER_STATS="$2"
      shift 2
      ;;
    --tape-num-of-conn)
      TAPE_NUM_OF_CONN="$2"
      shift 2
      ;;
    --tape-client-per-conn)
      TAPE_CLIENT_PER_CONN="$2"
      shift 2
      ;;
    --inter-case-sleep-seconds)
      INTER_CASE_SLEEP_SECONDS="$2"
      shift 2
      ;;
    --retry-on-goaway)
      RETRY_ON_GOAWAY="$2"
      shift 2
      ;;
    --goaway-max-retries)
      GOAWAY_MAX_RETRIES="$2"
      shift 2
      ;;
    --goaway-retry-backoff-seconds)
      GOAWAY_RETRY_BACKOFF_SECONDS="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

RUN_DIR="$RUNS_BASE_DIR/$RUN_ID"
mkdir -p "$RUN_DIR/logs" "$RUN_DIR/raw" "$RUN_DIR/summary" "$RUN_DIR/meta"

RUN_LOG="$RUN_DIR/logs/run.log"
SUMMARY_CSV="$RUN_DIR/summary/rounds.csv"

echo "scene,n,round,started_at_utc,ended_at_utc,duration_sec,exit_code,log_file" >"$SUMMARY_CSV"

{
  echo "run_id=$RUN_ID"
  echo "scenarios=$SCENARIOS"
  echo "n_levels=$N_LEVELS"
  echo "rounds=$ROUNDS"
  echo "seed_n=$SEED_N"
  echo "do_seed=$DO_SEED"
  echo "start_network=$START_NETWORK"
  echo "stop_network=$STOP_NETWORK"
  echo "backend_mode=$BACKEND_MODE"
  echo "backend_url=$BACKEND_URL"
  echo "wait_backend_seconds=$WAIT_BACKEND_SECONDS"
  echo "collect_docker_stats=$COLLECT_DOCKER_STATS"
  echo "tape_num_of_conn=$TAPE_NUM_OF_CONN"
  echo "tape_client_per_conn=$TAPE_CLIENT_PER_CONN"
  echo "inter_case_sleep_seconds=$INTER_CASE_SLEEP_SECONDS"
  echo "retry_on_goaway=$RETRY_ON_GOAWAY"
  echo "goaway_max_retries=$GOAWAY_MAX_RETRIES"
  echo "goaway_retry_backoff_seconds=$GOAWAY_RETRY_BACKOFF_SECONDS"
  echo "dry_run=$DRY_RUN"
  echo "started_at_utc=$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
} >"$RUN_DIR/meta/run.env"

log "Run directory: $RUN_DIR" | tee -a "$RUN_LOG"

if [[ "$START_NETWORK" -eq 1 ]]; then
  log "Starting network via blockchain/network/start.sh (this clears old chain data)." | tee -a "$RUN_LOG"
  run_cmd "cd \"$NETWORK_DIR\" && ./start.sh"
fi

start_backend_if_needed

if [[ "$DO_SEED" -eq 1 ]]; then
  log "Seeding users with register scene: n=$SEED_N" | tee -a "$RUN_LOG"
  run_tape_case "register" "$SEED_N" 1 || true
fi

declare -a scene_arr
parse_csv_to_array "$SCENARIOS" scene_arr

declare -a n_arr
parse_csv_to_array "$N_LEVELS" n_arr

for scene in "${scene_arr[@]}"; do
  for n in "${n_arr[@]}"; do
    for ((r=1; r<=ROUNDS; r++)); do
      run_tape_case "$scene" "$n" "$r" || true
      if [[ "$INTER_CASE_SLEEP_SECONDS" -gt 0 ]]; then
        sleep "$INTER_CASE_SLEEP_SECONDS"
      fi
    done
  done
done

stop_backend_if_needed

if [[ "$STOP_NETWORK" -eq 1 ]]; then
  log "Stopping network via blockchain/network/stop.sh" | tee -a "$RUN_LOG"
  run_cmd "cd \"$NETWORK_DIR\" && ./stop.sh"
fi

log "Completed. Summary: $SUMMARY_CSV" | tee -a "$RUN_LOG"

