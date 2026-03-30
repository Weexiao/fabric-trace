# Performance Scripts (fabric-trace)

This folder provides a minimal, runnable performance test workflow for your current repository layout.

## Files

- `perf/run_perf.sh`: Orchestrates tape-based performance rounds and saves results.
- `perf/extract_tape_metrics.py`: Best-effort parser for TPS/success/fail/latency fields from tape logs.
- `perf/api_ping_load.sh`: Lightweight API benchmark for `GET /ping` (no external load tool required).
- `perf/test_scripts.sh`: Local script checks + dry-run harness.

## Output Layout

Each run is saved under `perf/runs/<run_id>/`:

- `logs/run.log`: execution log
- `raw/tape_*.log`: raw tape output per scene/round
- `raw/docker_stats_*.csv`: docker stats snapshot per round
- `summary/rounds.csv`: round-level status and timing
- `summary/tape_metrics.csv`: parsed metrics (NA when not found)
- `meta/run.env`: run metadata

## Quick Start

1) (Optional) dry-run only:

```bash
cd /home/qikefan/work/fabric-trace
./perf/test_scripts.sh
```

2) Real chaincode pressure test (uses existing network/backend):

```bash
cd /home/qikefan/work/fabric-trace
./perf/run_perf.sh --run-id baseline_1 --scenarios invoke,query --n-levels 100,300,500 --rounds 2
```

3) Full flow including network startup and local backend:

```bash
cd /home/qikefan/work/fabric-trace
./perf/run_perf.sh --run-id fullflow_1 --start-network --backend-mode local --scenarios invoke,query --n-levels 100,300 --rounds 1
```

4) API `/ping` pressure test:

```bash
cd /home/qikefan/work/fabric-trace
./perf/api_ping_load.sh --base-url http://127.0.0.1:9090 --total 2000 --concurrency 100
```

5) If invoke logs show `GOAWAY too_many_pings`, use conservative throttling:

```bash
cd /home/qikefan/work/fabric-trace
./perf/run_perf.sh --run-id stable_low_1 --scenarios invoke,query --n-levels 20,50 --rounds 1 --tape-num-of-conn 1 --tape-client-per-conn 1 --tape-signers 1 --tape-parallel 1 --tape-rate 5 --tape-burst 10 --inter-case-sleep-seconds 8 --retry-on-goaway 1 --goaway-max-retries 3 --goaway-retry-backoff-seconds 30
```

## Notes

- `blockchain/network/start.sh` calls `stop.sh` internally and clears existing chain data.
- `run_perf.sh` seeds users by default via `config_register.yaml`; use `--skip-seed` to disable.
- Tape command format used by scripts is `./tape -c <config> -n <num> run`.
- `run_perf.sh` also supports `--tape-signers`, `--tape-parallel`, `--tape-rate`, and `--tape-burst` to lower gRPC pressure.
- If your local tape build prints different metric keys, `tape_metrics.csv` may show `NA`; raw logs remain intact for manual analysis.

