#!/bin/bash

# Needs 'hey' installed.
# https://github.com/rakyll/hey

echo "# # # Testing 400 route, which doesn't exist. # # #"
hey -n 10000 -m GET http://localhost:8080
echo "# # # # # # # # # # # # # # # # # # # # # # # # # #"


echo "# # # # Testing 200 route which responds directly. # # #"
hey -n 10000 -m GET http://localhost:8080/bar
echo "# # # # # # # # # # # # # # # # # # # # # # # # # #"

echo "# # # #  Testing Fibonacci calculation (Iterative) of random number in 1:1000_000 # # #"
hey -n 50000 -m GET http://localhost:8080/fibonacci
echo "# # # # # # # # # # # # # # # # # # # # # # # # # #"