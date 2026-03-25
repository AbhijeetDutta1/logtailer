#!/bin/bash
set -e

# Create a test log file
TEST_LOG="test.log"
echo "Initializing $TEST_LOG..." > "$TEST_LOG"

# Start the tailer in the background, redirect output to results.log
./logtailer "$TEST_LOG" > results.log 2>&1 &
TAILER_PID=$!

# Wait a bit for it to start
sleep 2

# Append more lines
echo "Line 1" >> "$TEST_LOG"
echo "Line 2" >> "$TEST_LOG"
sleep 1
echo "Line 3" >> "$TEST_LOG"

# Wait a bit more for lines to be processed
sleep 2

# Kill the tailer
kill $TAILER_PID || true

# Check results
echo "Results from results.log:"
cat results.log

if grep -q "Line 1" results.log && grep -q "Line 2" results.log && grep -q "Line 3" results.log; then
    echo "Verification SUCCESS"
else
    echo "Verification FAILED"
    exit 1
fi
