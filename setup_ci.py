import json
import os
import subprocess
import sys
import re
from distutils.core import Extension

import setuptools
from setuptools.command.build_ext import build_ext

"""
Based on:
    https://github.com/popatam/gopy_build_wheel_example/blob/main/setup_ci.py
"""

def normalize(name):  # https://peps.python.org/pep-0503/#normalized-names
    return re.sub(r"[-_.]+", "-", name).lower()

# PACKAGE_PATH="simple_go_timer"
# PACKAGE_NAME=PACKAGE_PATH.split("/")[-1]

PACKAGE_PATH="gohpygossh"
PACKAGE_NAME="ohpygossh"

if sys.platform == 'darwin':
    # PYTHON_BINARY_PATH is setting explicitly for 310 and 311, see build_wheel.yml
    # on macos PYTHON_BINARY_PATH must be python bin installed from python.org or from brew
    PYTHON_BINARY = os.getenv("PYTHON_BINARY_PATH", sys.executable)
    if PYTHON_BINARY == sys.executable:
        subprocess.check_call([sys.executable, '-m', 'pip', 'install', 'pybindgen'])
else:
    # linux & windows
    PYTHON_BINARY = sys.executable
    subprocess.check_call([sys.executable, '-m', 'pip', 'install', 'pybindgen'])

def _generate_path_with_gopath() -> str:
    go_path = subprocess.check_output(["go", "env", "GOPATH"]).decode("utf-8").strip()
    path_val = f'{os.getenv("PATH")}:{go_path}/bin'
    return path_val


class CustomBuildExt(build_ext):
    def build_extension(self, ext: Extension):
        bin_path = _generate_path_with_gopath()
        go_env = json.loads(subprocess.check_output(["go", "env", "-json"]).decode("utf-8").strip())

        destination = os.path.dirname(os.path.abspath(self.get_ext_fullpath(ext.name))) + f"/{PACKAGE_NAME}"

        subprocess.check_call(
            [
                "gopy",
                "build",
                "-no-make",
                "-dynamic-link=True",
                "-output",
                destination,
                "-vm",
                PYTHON_BINARY,
                *ext.sources,
            ],
            env={"PATH": bin_path, **go_env, "CGO_LDFLAGS_ALLOW": ".*"},
        )

        # dirty hack to avoid "from pkg import pkg", remove if needed
        with open(f"{destination}/__init__.py", "w") as f:
            f.write(f"from .{PACKAGE_NAME} import *")


setuptools.setup(
    name="ohpygossh",
    version="0.0.9",
    author="b-long",
    description="A project to create a new cross-platform SSH wheel for Python.",
    url="https://github.com/b-long/ohpygossh",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: BSD License",
        "Operating System :: OS Independent",
    ],
    include_package_data=True,
    cmdclass={
        "build_ext": CustomBuildExt,
    },
    ext_modules=[
        Extension(PACKAGE_NAME, [PACKAGE_PATH],)
    ],
)