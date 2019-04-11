# VaaS Registration Hook

[![Build Status](https://travis-ci.org/allegro/vaas-registration-hook.svg?branch=master)](https://travis-ci.org/allegro/vaas-registration-hook)
[![Go Report Card](https://goreportcard.com/badge/github.com/allegro/vaas-registration-hook)](https://goreportcard.com/report/github.com/allegro/vaas-registration-hook)
[![Codecov](https://codecov.io/gh/allegro/vaas-registration-hook/branch/master/graph/badge.svg)](https://codecov.io/gh/allegro/vaas-registration-hook)
[![GoDoc](https://godoc.org/github.com/allegro/vaas-registration-hook?status.svg)](https://godoc.org/github.com/allegro/vaas-registration-hook)


[VaaS][1] integration based on a hook mechanic.
An app is usually registered with VaaS once it becomes healthy and deregistered before termination.

## Usage
### CLI
Task’s desired address (`--addr`) and port (`--port, -p`) will be (de)registered under 
a director provided by `--director`. A VaaS API url needs to be provided (`--vaas-url` or `VAAS_URL`) 
along with an API user (`--user, -u`) and secret key (`--key, -k`). 
If task needs a defined weight it can be provided with `--weight` at registration.
Registered backend can be tagged as a canary using `--canary`. 

Examples:
```bash
export VAAS_URL="http://vaas.example.com/api"
export VAAS_USER="admin"
export VAAS_KEY="secret-key"
vaas-hook --debug --addr=192.168.0.10 --port 80 --director=hook-test register cli --weight 1 --dc dc1
vaas-hook --debug --addr=192.168.0.10 --port 80 --director=hook-test deregister cli
```

### Kubernetes
This hook can also read a Kubernetes environment and access annotations via it's Pod API.
All the available annotations can be viewed in [k8s/pod.go](k8s/pod.go).
Action names and debug flag still needs to be provided via command line.
A working example can be found in [examples/service-with-lifecycle.yaml](examples/service-with-lifecycle.yaml)

Examples:
```bash
vaas-hook --debug register k8s
vaas-hook --debug deregister k8s
```

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
