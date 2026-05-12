#!/bin/bash
# Generate a detailed Markdown coverage report from coverage.out

COVERAGE_FILE=$1
MODULE_NAME="digital-innovation/gostrategy"

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "No coverage data available."
    exit 0
fi

# Parse coverage.out, group by package, and calculate per-file coverage
awk -v mod="$MODULE_NAME" '
  BEGIN { FS="[ :]" }
  $1 == "mode" { next }
  {
    full_path = $1
    sub("^" mod "/", "", full_path)
    
    if (full_path ~ /\//) {
      pkg = full_path; sub(/\/[^\/]+$/, "", pkg)
      file = full_path; sub(/.*\//, "", file)
    } else {
      pkg = "(root)"; file = full_path
    }
    
    key = pkg "|" file
    f_stmts[key] += $3
    if ($4 > 0) { f_hits[key] += $3 }
  }
  END {
    for (k in f_stmts) {
      split(k, a, "|")
      cov = f_stmts[k] > 0 ? (f_hits[k]/f_stmts[k])*100 : 0
      printf "%s|%s|%.1f%%\n", a[1], a[2], cov
    }
  }
' "$COVERAGE_FILE" | sort -t'|' -k1,1 -k2,2 | awk -F'|' '
  BEGIN { 
    print "| Package | File | Coverage |"
    print "| :--- | :--- | :--- |" 
  }
  {
    pkg = $1; file = $2; cov = $3
    if (pkg == last_pkg) {
      printf "| | `%s` | %s |\n", file, cov
    } else {
      printf "| **%s** | `%s` | %s |\n", pkg, file, cov
      last_pkg = pkg
    }
  }
'
