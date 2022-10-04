#!/usr/bin/env bash

files_to_lint=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
if [[ -n "$files_to_lint" ]]; then
    echo "Running static analysis..."
    make lint

    status=$?

    if [[ "$status" = 0 ]] ; then
        echo "Static analysis found no problems."
        exit 0
    else
        echo 1>&2 "Static analysis found violations it could not fix."
        echo "Run make fmt to fix issues"
        exit 1
    fi
else
    echo  "Static analysis not needed - No Go file modified"
    exit 0
fi
