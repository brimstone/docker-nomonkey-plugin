Docker No monkey Business Plugin
================================
_In order to use this plugin you need to be running at least Docker 1.10 which
has support for authorization plugins._

Sometimes, you want to let people have access to your docker socket that you
don't really trust. This plugin helps cut down on the monkey business they
could cause.

This plugin solves this issue by disallowing starting a container with options
that could allow compromise of the host.
In particular, the plugin will block `docker run` with:

- `-v` host volumes bound, some directories can be whitelisted.
- `--capadd` adding additional capabilities to a container
- `--device` adding host devices to the container
- `--privileged` giving all capabilities to a container

Building
--------
This is a [projectbuilder](https://github.com/brimstone/projectbuilder) enabled project.

If you don't have projectbuilder configured:

1. Properly configure your `GOPATH`
2. `go get && go build` should be all you need

Installing
----------
1. Copy this to somewhere on the host.
2. Configure your process supervisor to start it before or shortly after
docker-engine.
3. Configure docker-engine to start with `--authorization-plugin=docker-novolume-plugin`

If this program stops running, docker-engine will not start any new containers.

License
-
MIT
