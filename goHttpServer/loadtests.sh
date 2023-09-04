# Needs 'hey' installed.
# https://github.com/rakyll/hey

# Testing 400 route, which doesn't exist.
hey -n 10000 -m GET http://localhost:8080

# Testing 200 route which responds directly.
hey -n 10000 -m GET http://localhost:8080/bar
