#!/bin/bash
set -euo pipefail

# Mark the workspace as safe for git operations inside the container
git config --global --add safe.directory /github/workspace

# Build pulse gate flags
FLAGS="--no-color"

if [ -n "${INPUT_MAX_CCN:-}" ] && [ "${INPUT_MAX_CCN}" != "0" ]; then
  FLAGS="${FLAGS} --max-ccn ${INPUT_MAX_CCN}"
fi

if [ -n "${INPUT_MAX_COGNITIVE:-}" ] && [ "${INPUT_MAX_COGNITIVE}" != "0" ]; then
  FLAGS="${FLAGS} --max-cognitive ${INPUT_MAX_COGNITIVE}"
fi

if [ -n "${INPUT_MAX_DUPLICATION:-}" ] && [ "${INPUT_MAX_DUPLICATION}" != "0" ]; then
  FLAGS="${FLAGS} --max-duplication ${INPUT_MAX_DUPLICATION}"
fi

if [ -n "${INPUT_MIN_MAINTAINABILITY:-}" ] && [ "${INPUT_MIN_MAINTAINABILITY}" != "0" ]; then
  FLAGS="${FLAGS} --min-maintainability ${INPUT_MIN_MAINTAINABILITY}"
fi

echo "::group::Pulse Quality Gate"
# `pulse gate` exits non-zero when a gate is breached. Under `set -e` a bare
# non-zero command aborts the script immediately, which skipped the
# INPUT_FAIL_ON_GATE handling below entirely — so `fail_on_gate: false` never
# took effect. Capture the status via `|| ...` (exempt from `set -e`) instead.
EXIT_CODE=0
# shellcheck disable=SC2086
pulse gate ${FLAGS} . || EXIT_CODE=$?
echo "::endgroup::"

if [ "${INPUT_FAIL_ON_GATE:-true}" = "false" ]; then
  echo "fail_on_gate=false: reporting gate result without failing the job (pulse exit ${EXIT_CODE})."
  exit 0
fi

exit ${EXIT_CODE}
