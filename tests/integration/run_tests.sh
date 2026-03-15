#!/bin/bash
# ATLAS Compliance Test Runner for Linux/Mac

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "========================================"
echo "ATLAS Compliance Test Suite"
echo "========================================"
echo ""

# Parse arguments
VERBOSE=""
COVERAGE=""
RACE=""
TEST_TARGET="Atlas"

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=1
            shift
            ;;
        -coverage)
            COVERAGE=1
            shift
            ;;
        -race)
            RACE=1
            shift
            ;;
        -run)
            TEST_TARGET="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

echo "Running ATLAS compliance tests..."
echo ""

# Build test flags
GO_TEST_FLAGS="-v"

if [[ -n "$VERBOSE" ]]; then
    GO_TEST_FLAGS="$GO_TEST_FLAGS -count=1"
fi

if [[ -n "$COVERAGE" ]]; then
    GO_TEST_FLAGS="$GO_TEST_FLAGS -coverprofile=coverage.out -covermode=atomic"
fi

if [[ -n "$RACE" ]]; then
    GO_TEST_FLAGS="$GO_TEST_FLAGS -race"
fi

GO_TEST_FLAGS="$GO_TEST_FLAGS -run $TEST_TARGET ./tests/integration/..."

# Run tests
go test $GO_TEST_FLAGS
TEST_RESULT=$?

echo ""
if [[ $TEST_RESULT -eq 0 ]]; then
    echo "========================================"
    echo "All tests passed!"
    echo "========================================"
else
    echo "========================================"
    echo "Tests failed with error code: $TEST_RESULT"
    echo "========================================"
fi

if [[ -n "$COVERAGE" ]] && [[ -f coverage.out ]]; then
    echo ""
    echo "Coverage Report:"
    go tool cover -func=coverage.out
fi

exit $TEST_RESULT
