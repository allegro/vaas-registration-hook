# VaaS Registration Hook

[![Build Status](https://travis-ci.org/allegro/vaas-registration-hook.svg?branch=master)](https://travis-ci.org/allegro/vaas-registration-hook)
[![Go Report Card](https://goreportcard.com/badge/github.com/allegro/vaas-registration-hook)](https://goreportcard.com/report/github.com/allegro/vaas-registration-hook)
[![Codecov](https://codecov.io/gh/allegro/vaas-registration-hook/branch/master/graph/badge.svg)](https://codecov.io/gh/allegro/vaas-registration-hook)
[![GoDoc](https://godoc.org/github.com/allegro/vaas-registration-hook?status.svg)](https://godoc.org/github.com/allegro/vaas-registration-hook)


[VaaS][1] integration based on a hook mechanic.
An app is usually registered with VaaS once it becomes healthy and deregistered before termination.
Task’s desired address (`-aaddress`) and port (`--port`) will be registered under director 
provided by `--director`.
If task has defined weight it can be provided with `--weight`
If task is a canary instance (has `--canary` switch) backend is tagged
as a canary.

## Requirements

To run executor tests locally you need following tools installed:

* [Go 1.9+][2]
* [Make][3]
* Linux or Docker


## Debug mode

vaas-registration hook offers a debug mode that provide extended logging and capabilities during
runtime. To enable debug mode add `--debug` flag to the command or set `VAAS_HOOK_DEBUG` 
environment variable to `true`.

## Development

To build this project, just execute the following command in project root folder:

```
$ make
```

It will run tests and create a binary and a ZIP package for release purposes.


## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for more details and code of conduct. 

## License

Mesos Executor is distributed under the [Apache 2.0 License](LICENSE).


[1]: https://github.com/allegro/vaas
[2]: https://golang.org/dl/
[3]: https://www.gnu.org/software/make/
[6]: https://brandur.org/logfmt
