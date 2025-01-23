#!/bin/bash

# set -x
# set -eou pipefail

# Ensure we aren't in a virtual environment
# deactivate || { echo "Not currently in a virtual environment" ; }

# Cleanup
rm -rf .venv/

POETRY=$(type -p poetry)
if [[ $? -ne "0" ]]; then
  echo "Unable to find 'poetry', it must be installed to continue"
  exit 1
fi

# Install python deps
"$POETRY" install --no-root

# Fail fast if virtual environment is missing
if ! [ -f .venv/bin/activate ]; then
    echo "Error: Failed to to locate virtual environment"
    echo "Unable to continue"
    exit 1
fi

# Activate shell with 'pybindgen' etc.
# NOTE: Using 'poetry shell' does not work
source .venv/bin/activate

# Prove that the wheel can be installed
pip install dist/ohpygossh-0.0.10-py3-none-any.whl

# Validate functionality
python validate_ohpygossh.py
