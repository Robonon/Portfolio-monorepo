#!/bin/bash

set -e

trap 'echo "Error occurred during load test. Exiting."; exit 1' ERR

TOTAL=1000
PARALLEL=50
TMPFILE=$(mktemp)

cleanup() {
  rm -f "$TMPFILE"
}
trap cleanup EXIT

start_time=$(date +%s)

# Start background progress watcher
(
  while true; do
    COMPLETED=$(wc -l < "$TMPFILE")
    elapsed=$(( $(date +%s) - start_time ))
    echo -ne "Completed requests: $COMPLETED/$TOTAL | Elapsed: ${elapsed}s\r"
    [[ $COMPLETED -ge $TOTAL ]] && break
    sleep 0.2
  done
  echo
) &

WATCHER_PID=$!

seq 1 $TOTAL | xargs -P$PARALLEL -I{} bash -c '
  nums="[$((RANDOM%100)), $((RANDOM%100)), $((RANDOM%100)), $((RANDOM%100)), $((RANDOM%100))]"
  curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/calculations/max \
    -H "Content-Type: application/json" \
    -d "{\"numbers\": $nums}" > /dev/null
  echo 1 >> '"$TMPFILE"'
'

wait $WATCHER_PID

end_time=$(date +%s)
duration=$((end_time - start_time))
echo -e "\nLoad test completed successfully in ${duration}s."