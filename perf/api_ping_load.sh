#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://127.0.0.1:9090"
TOTAL=500
CONCURRENCY=50
OUT_DIR=""

usage() {
  cat <<'EOF'
Usage:
  ./perf/api_ping_load.sh [options]

Options:
  --base-url <url>       API base URL (default: http://127.0.0.1:9090)
  --total <num>          Total requests (default: 500)
  --concurrency <num>    Concurrent workers (default: 50)
  --out-dir <path>       Output directory (default: perf/runs/api_<ts>)
  --help                 Show help

This script benchmarks GET /ping and writes raw latency + summary files.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --base-url)
      BASE_URL="$2"
      shift 2
      ;;
    --total)
      TOTAL="$2"
      shift 2
      ;;
    --concurrency)
      CONCURRENCY="$2"
      shift 2
      ;;
    --out-dir)
      OUT_DIR="$2"
      shift 2
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

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -z "$OUT_DIR" ]]; then
  OUT_DIR="$ROOT_DIR/perf/runs/api_$(date +%Y%m%d_%H%M%S)"
fi
mkdir -p "$OUT_DIR"

RAW_FILE="$OUT_DIR/ping_latencies.csv"
SUMMARY_FILE="$OUT_DIR/ping_summary.txt"

echo "req_id,latency_ms,http_code" >"$RAW_FILE"

export BASE_URL

# Each worker prints: req_id,latency_ms,http_code
seq 1 "$TOTAL" | xargs -I{} -P "$CONCURRENCY" bash -c '
  req_id="$1"
  start_ms=$(date +%s%3N)
  code=$(curl -sS -m 5 -o /dev/null -w "%{http_code}" "$BASE_URL/ping" || echo 000)
  end_ms=$(date +%s%3N)
  latency=$((end_ms - start_ms))
  printf "%s,%s,%s\n" "$req_id" "$latency" "$code"
' _ {} >>"$RAW_FILE"

awk -F',' 'NR>1 {print $2}' "$RAW_FILE" | sort -n >"$OUT_DIR/_latencies_sorted.txt"

count=$(wc -l <"$OUT_DIR/_latencies_sorted.txt")
if [[ "$count" -eq 0 ]]; then
  echo "No samples collected." >"$SUMMARY_FILE"
  cat "$SUMMARY_FILE"
  exit 1
fi

p50_idx=$(( (count * 50 + 99) / 100 ))
p95_idx=$(( (count * 95 + 99) / 100 ))
p99_idx=$(( (count * 99 + 99) / 100 ))

p50=$(sed -n "${p50_idx}p" "$OUT_DIR/_latencies_sorted.txt")
p95=$(sed -n "${p95_idx}p" "$OUT_DIR/_latencies_sorted.txt")
p99=$(sed -n "${p99_idx}p" "$OUT_DIR/_latencies_sorted.txt")

success=$(awk -F',' 'NR>1 && $3==200 {c++} END{print c+0}' "$RAW_FILE")
failed=$((TOTAL - success))
avg=$(awk -F',' 'NR>1 {sum+=$2; c++} END{if(c==0) print "NA"; else printf "%.2f", sum/c}' "$RAW_FILE")

echo "base_url=$BASE_URL" >"$SUMMARY_FILE"
echo "total=$TOTAL" >>"$SUMMARY_FILE"
echo "concurrency=$CONCURRENCY" >>"$SUMMARY_FILE"
echo "success=$success" >>"$SUMMARY_FILE"
echo "failed=$failed" >>"$SUMMARY_FILE"
echo "avg_ms=$avg" >>"$SUMMARY_FILE"
echo "p50_ms=$p50" >>"$SUMMARY_FILE"
echo "p95_ms=$p95" >>"$SUMMARY_FILE"
echo "p99_ms=$p99" >>"$SUMMARY_FILE"
echo "raw=$RAW_FILE" >>"$SUMMARY_FILE"

cat "$SUMMARY_FILE"
rm -f "$OUT_DIR/_latencies_sorted.txt"

