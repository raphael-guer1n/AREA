#!/bin/bash

# Test summary script - parses Go test output and displays clean summary

# Colors
RESET='\033[0m'
BOLD='\033[1m'
GREEN='\033[32m'
RED='\033[31m'
YELLOW='\033[33m'
CYAN='\033[36m'
DIM='\033[2m'

# Temporary file to store results
tmpfile=$(mktemp)
trap "rm -f $tmpfile" EXIT

current_service=""
pass_count=0
fail_count=0
skip_count=0
coverage="N/A"
max_coverage=0
has_coverage=false
result=""

# Strip ANSI codes function
strip_ansi() {
    echo "$1" | sed 's/\x1b\[[0-9;]*m//g'
}

# Process stdin line by line
while IFS= read -r line; do
    # Strip ANSI color codes
    clean_line=$(strip_ansi "$line")

    # Detect service header (with or without ANSI codes)
    if [[ $clean_line =~ ^==\>\ (.+)\ make\ test_run ]]; then
        # Save previous service if exists
        if [[ -n "$current_service" ]]; then
            # Use the maximum coverage seen for this service
            if [[ $has_coverage == true ]]; then
                coverage="${max_coverage}%"
            fi
            echo "$current_service|$result|$pass_count|$fail_count|$skip_count|$coverage" >> "$tmpfile"
        fi

        # Reset for new service
        current_service="${BASH_REMATCH[1]}"
        pass_count=0
        fail_count=0
        skip_count=0
        coverage="N/A"
        max_coverage=0
        has_coverage=false
        result=""
    fi

    # Count test results
    if [[ $clean_line =~ ^---\ PASS: ]]; then
        ((pass_count++))
    elif [[ $clean_line =~ ^---\ FAIL: ]]; then
        ((fail_count++))
    elif [[ $clean_line =~ ^---\ SKIP: ]]; then
        ((skip_count++))
    fi

    # Extract coverage from summary line - keep the highest coverage value
    if [[ $clean_line =~ coverage:\ ([0-9.]+)%\ of\ statements ]]; then
        cov_value="${BASH_REMATCH[1]}"
        # Compare as floating point
        if (( $(echo "$cov_value > $max_coverage" | bc -l) )); then
            max_coverage="$cov_value"
            has_coverage=true
        fi
    fi

    # Detect overall result
    if [[ $clean_line =~ ^PASS$ ]]; then
        result="PASS"
    elif [[ $clean_line =~ ^FAIL$ ]]; then
        result="FAIL"
    fi
done

# Save last service
if [[ -n "$current_service" ]]; then
    if [[ $has_coverage == true ]]; then
        coverage="${max_coverage}%"
    fi
    echo "$current_service|$result|$pass_count|$fail_count|$skip_count|$coverage" >> "$tmpfile"
fi

# Print summary header
echo ""
echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════${RESET}"
echo -e "${BOLD}${CYAN}                          TEST SUMMARY${RESET}"
echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════${RESET}"
echo ""
printf "${BOLD}%-30s %-10s %-25s %-12s${RESET}\n" "SERVICE" "STATUS" "TESTS" "COVERAGE"
echo -e "${DIM}───────────────────────────────────────────────────────────────────────${RESET}"

# Read and display results
if [[ -s "$tmpfile" ]]; then
    while IFS='|' read -r service result pass fail skip coverage; do
        total=$((pass + fail + skip))

        # Determine status color and text
        if [[ $result == "PASS" ]]; then
            status_color="${GREEN}"
            status_text="✓ PASS"
        else
            status_color="${RED}"
            status_text="✗ FAIL"
        fi

        # Format test counts
        test_info="${pass} OK"
        if [[ $fail -gt 0 ]]; then
            test_info="${test_info}, ${fail} KO"
        fi
        if [[ $skip -gt 0 ]]; then
            test_info="${test_info}, ${skip} SKIP"
        fi
        test_info="${test_info} / ${total}"

        # Print row
        printf "%-30s ${status_color}%-10s${RESET} %-25s %-12s\n" \
            "$service" \
            "$status_text" \
            "$test_info" \
            "$coverage"
    done < <(sort "$tmpfile")
else
    echo "No test results found"
fi

echo -e "${BOLD}═══════════════════════════════════════════════════════════════════════${RESET}"
echo ""
