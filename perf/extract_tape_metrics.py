#!/usr/bin/env python3
"""Extract best-effort metrics from a tape log into CSV.

The tape output format can vary by build. This parser searches common tokens and
writes NA when a metric is not present.
"""

from __future__ import annotations

import argparse
import csv
import os
import re
from typing import Dict

PATTERNS = {
    "tps": [
        re.compile(r"\bTPS\b\s*[:=]\s*([0-9]+(?:\.[0-9]+)?)", re.IGNORECASE),
        re.compile(r"\bthroughput\b\s*[:=]\s*([0-9]+(?:\.[0-9]+)?)", re.IGNORECASE),
    ],
    "success": [
        re.compile(r"\bsuccess\b\s*[:=]\s*([0-9]+)", re.IGNORECASE),
        re.compile(r"\bsucc(?:ess)?\b\s*[:=]\s*([0-9]+)", re.IGNORECASE),
    ],
    "fail": [
        re.compile(r"\bfail(?:ure)?\b\s*[:=]\s*([0-9]+)", re.IGNORECASE),
        re.compile(r"\berror\b\s*[:=]\s*([0-9]+)", re.IGNORECASE),
    ],
    "p95_ms": [
        re.compile(r"\bp95\b\s*[:=]\s*([0-9]+(?:\.[0-9]+)?)\s*ms", re.IGNORECASE),
        re.compile(r"95(?:th)?\s*percentile\s*[:=]\s*([0-9]+(?:\.[0-9]+)?)\s*ms", re.IGNORECASE),
    ],
    "p99_ms": [
        re.compile(r"\bp99\b\s*[:=]\s*([0-9]+(?:\.[0-9]+)?)\s*ms", re.IGNORECASE),
        re.compile(r"99(?:th)?\s*percentile\s*[:=]\s*([0-9]+(?:\.[0-9]+)?)\s*ms", re.IGNORECASE),
    ],
}


def extract_metrics(content: str) -> Dict[str, str]:
    metrics: Dict[str, str] = {k: "NA" for k in PATTERNS}
    for key, patterns in PATTERNS.items():
        for pattern in patterns:
            m = pattern.search(content)
            if m:
                metrics[key] = m.group(1)
                break
    return metrics


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", required=True, help="tape log file path")
    parser.add_argument("--output", required=True, help="csv output path")
    parser.add_argument("--scene", required=True)
    parser.add_argument("--n", required=True)
    parser.add_argument("--round", required=True)
    parser.add_argument("--append", action="store_true")
    args = parser.parse_args()

    with open(args.input, "r", encoding="utf-8", errors="ignore") as f:
        content = f.read()

    row = {
        "scene": args.scene,
        "n": args.n,
        "round": args.round,
        **extract_metrics(content),
        "input": args.input,
    }

    fieldnames = ["scene", "n", "round", "tps", "success", "fail", "p95_ms", "p99_ms", "input"]
    os.makedirs(os.path.dirname(args.output), exist_ok=True)

    need_header = not os.path.exists(args.output) or not args.append
    mode = "a" if args.append else "w"
    with open(args.output, mode, newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        if need_header:
            writer.writeheader()
        writer.writerow(row)


if __name__ == "__main__":
    main()

