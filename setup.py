import setuptools

setuptools.setup(
    name="ohpygossh",
    packages=setuptools.find_packages(include=["ohpygossh"]),
    py_modules = ["ohpygossh.gohpygossh"],
    package_data={"ohpygossh": ["*.so"]},
    # Should match 'pyproject.toml' version number
    version="0.0.2",
    include_package_data=True,
)
