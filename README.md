# ohpygossh

A project to create a new cross-platform SSH wheel for Python.

<img width="617" alt="image" src="docs/assets/banner.png">

## Goals

1. An easy to install, multi-platform SSH wheel
1. To explore [`gopy`](https://github.com/go-python/gopy)

## Code of Conduct

1. Respect and give kudos to all the work that's come
before.
    1. In particular, the Herculean efforts of [`gopy`], [`paramiko`] and [`asyncssh`].
    1. Please also be kind to all other open-source software projects you find.

It is paramount that we all play nicely.  To that end, please do your very best
not to create churn or spur conversations that may upset other developers,
and/or cause debate without offering solutions.


## Development

⚠️ This section needs to be rewritten ⚠️

### Build wheel

The following steps should produce a wheel,
located at `dist/ohpygossh-<version>-py3-none-any.whl`

```bash
./make_and_validate_script.sh
```

### Naming

This project has to think about naming in two different contexts:

* The Python library, imported by end-users: `ohpygossh`.
* The Golang package, also used in the `pyproject.toml`: `gohpygossh`.

Briefly:

```python
# Install the library:
#   pip install ohpygossh

# Import functions and data structures
from ohpygossh.gohpygossh import ...
```

## Reference

Based on:
* [`gopy`]
* https://github.com/b-long/cookiecutter-gopy

[`gopy`]: https://github.com/go-python/gopy
[`paramiko`]: https://pypi.org/project/paramiko/
[`asyncssh`]: https://pypi.org/project/asyncssh/
