#!/bin/bash

# set -x
# set -eou pipefail

# Ensure we aren't in a virtual environment
# deactivate || { echo "Not currently in a virtual environment" ; }

# Cleanup
rm -rf .venv/
rm -rf dist/

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

# For every step below, 'which python' should return '.venv/bin/python'
PATH="$PATH:$HOME/go/bin" gopy build -output=ohpygossh -vm=python3 .

python3 -m pip install --upgrade setuptools wheel

# Build the 'dist/' folder (wheel)
python3 setup.py bdist_wheel

# Prove that the wheel can be installed
pip install dist/ohpygossh-0.0.7-py3-none-any.whl

# Validate functionality
python3 validate_ohpygossh.py
