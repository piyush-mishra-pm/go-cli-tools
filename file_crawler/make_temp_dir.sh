#!/bin/bash
# Creates a set of temp files, to test code.
mkdir -p /tmp/testdir/{text,logs}
touch /tmp/testdir/file1.txt
touch /tmp/testdir/text/{text1,text2,text3}.txt
touch /tmp/testdir/logs/{log1,log2,log3}.log
echo "made files"
