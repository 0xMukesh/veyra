#!/usr/bin/fish

go test -v ./internal/cpu/ \
| grep 'A:' \
> ./output/cpu_test.log
