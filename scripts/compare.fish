#!/usr/bin/fish

awk 'NR==FNR {a[NR]=$0; next}
     {
       if (a[FNR] != $0) {
         printf("Line %d\nNES: %s\nEMU: %s\n\n", FNR, a[FNR], $0)
       }
     }' ./roms/nestest_no_cycle.log ./output/cpu_test.log > ./output/diff.txt
