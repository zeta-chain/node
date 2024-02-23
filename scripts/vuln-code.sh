#!/bin/bash

# Simple script with intentional issues for shellcheck testing

# Unquoted variable
unquoted_variable=Hello

# Missing shebang
echo "This script doesn't have a shebang"

# Unused variable
unused_variable="This variable is unused"

# Incorrect if statement syntax
if [ $unquoted_variable == "Hello" ]; then
  echo "Condition is true"
fi
