#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

bash -n "$ROOT_DIR/perf/run_perf.sh"
bash -n "$ROOT_DIR/perf/api_ping_load.sh"
python3 -m py_compile "$ROOT_DIR/perf/extract_tape_metrics.py"

# Dry-run verifies CLI wiring and output structure without touching services.
"$ROOT_DIR/perf/run_perf.sh" \
  --dry-run \
  --run-id dryrun_check \
  --scenarios invoke,query \
  --n-levels 5,10 \
  --rounds 1 \
  --backend-mode none \
  --collect-docker-stats 0

echo "All script checks passed."

