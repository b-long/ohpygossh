# ohpygossh

A project to create a new cross-platform SSH wheel for Python.

<img width="617" alt="image" src="docs/assets/banner.png">


## Goals

1. An easy to install, multi-platform SSH wheel
1. Personal learning about Go
1. To explore [`gopy`](https://github.com/go-python/gopys)
1. To explore Python packaging & compiling.  What can we
do with available open source tooling?
1. To explore performance characteristics

## Code of Conduct

1. Respect and give kudos to all the work that's come
before.
    1. In particular, the Herculian efforts of [`paramiko`] and [`asyncssh`].
    1. Please also be very knid to all project dependencies

## Play nicely

It is paramount that we all play nicely.  To that end, please do your very best
not to create churn or spur conversations that may upset other developers,
and/or cause debate without offerinig solutions.

## Build wheel

The following steps should produce a wheel,
located at `dist/ohpygossh-0.0.11-py3-none-any.whl`

```bash
./make_and_validate_script.sh
```

## Naming

This project has to think about components in two different contexts.

To distinguish the two, use the following names:

* The Python module is imported as `ohpygossh`.
* The Golang package is used as `gohpygossh`.

NOTE: Currently, the `pyproject.toml` in this project is named for the
Golang package.  This may need to change as the project matures.

As a Python user:

* Install the library, with `pip install ohpygossh`
* Import functions and data structures, with `from ohpygossh.gohpygossh import ...`


## Reference

Based on:
https://last9.io/blog/using-golang-package-in-python-using-gopy/

Learn Go:
* https://learnxinyminutes.com/docs/go/
* https://gist.github.com/prologic/5f6afe9c1b98016ca278f4d507e65510

## TODO

- [ ] Figure out if generaetd folder (`ls -la ohpygossh/`) should be kept in git.



[`paramiko`]: https://pypi.org/project/paramiko/
[`asyncssh`]: https://pypi.org/project/asyncssh/
