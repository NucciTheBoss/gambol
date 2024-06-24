# gambol

A simple workflow engine for distributed applications 🎬 🌌

gambol enables you to write straightforward workflows that can be run on
no-frills cloud environments. Rather than require you to be an expert in
cloud infrastructure, gambol allows you focus on what matters: __getting
work done__. You can write end-to-end tests for your distributed application,
test packages in isolated environments, create comphrensive testing suites
that can run both locally and in the cloud, and more. Say goodbye to spending
half your day just setting up your testing infrastructure with shell script soup!

To get started with gambol, check out the getting started how-to below 👇

## Getting started ✨

gambol is inspired by [cleantest](https://github.com/NucciTheBoss/cleantest).
cleantest originally started out with the goal of taking Python functions
that perform some destructive action like modifying the underlying
filesystem, and executing it within an isolated test environment. cleantest
worked great at first, however, there were some challenges that emerged as
the test framework matured.

First, it's really hard to support Python applications that must support
environments with wide interpreter version stratification. The Python API is a
moving target which made it difficult to work with cleantest's _injectable_
mechanism where it would copy over the body of a function and execute within an
isolated environment. Second, cleantest had a bit of an identity crisis where it
became both a testing framework and a workflow orchestrator. You could write tests
Python tests the traditional way like you can with `unittest` or `pytest`, but you
could also orchestrate things that had absolutely nothing to do with testing Python
modules like set up a mini-supercomputer or end-to-end tests for a separate
application. Difficult to sell people on your using your framework when you've
effectively railroaded two different paradigms together!

gambol is a breakout of that second half of cleantest. It's a workflow orchestrator
for distributed applications. It can request isolated environment dynamically, test
a wide variety of applications, and has straightforward syntax for creating
workflows.

### Sounds great, right? Let's get to it then!

You first need to setup a cloud that will provide the instances that gambol will use
to run workflows. For this how-to, we can use [LXD](https://ubuntu.com/lxd) as our
backing cloud, and will use snap to install gambol.

Run the following commands to install LXD on your system:

```shell
snap install lxd
lxd init --minimal
```

And now use the following command to install gambol:

```shell
# [WIP]: Snap is not published yet
snap install gambol
```

Now you're ready to start using gambol.

### But first, some explanations!

Below are some key concepts in gambol that you should be familiar with to better
understand how it works:

#### Playthrough

A playthrough is a YAML file that defines the workflow that gambol to will run.

#### Provider

A provider is a cloud that provides the containers and/or virtual machines that
gambol will use run to run defined playthrough in. Currently, LXD is the only
supported option.

#### Act

An Act is a sequence of steps that you want to execute within an instance requested
from the configured provider. Each Act corresponds to a single instance, or an Act
can correspond an instance that has been requested for a previous Act.

#### Scene

A Scene is a step within an Act. Scenes wrap executable blocks, and will report
whether or not the given block completes successfully within an Act instance.

### Run your first playthrough

Using your favorite text editor, create the file _playthrough.yaml_, and enter the
following document. This playthrough will create the Act instance `act-1`, and will
run the scenes `Install cowsay` and `Say hello word`:

```yaml
name: "my first playthrough"
provider:
  lxd:
acts:
  act-1:
    name: "Say hello to the world with gambol!"
    run-on: noble
    scenes:
      - name: "Install cowsay"
        run: |
          export DEBIAN_FRONTEND=noninteractive
          apt-get -y install cowsay
      - name: "Say hello world"
        run: |
          cowsay hello world
```

After you have created the file _playthrough.yaml_, use the following command to run
the playthrough with gambol:

```shell
gambol run playthrough.yaml
```

Congratulations! You have run your first playthrough using gambol 🎉

## Where to next? 🤔

gambol can do a lot more than just make an ASCII cow say hello within a system
container. I'm still working on more comprehensive documentation for gambol, but in
the meantime, you check out the [e2e tests](./test/e2e/) for more advanced use cases
of gambol/

## Development 🛠️

Right now, gambol is very much a hacking and wacking endevour, so there's no
formal contribution process quite yet. However, if you are interested in
contributing to gambol, please ensure the all e2e tests are passing before
opening a pull request:

```shell
make e2e
```

All contributions must be licensed under the [AGPLv3 license](./LICENSE).

## Project and community 🤝

If you're interested in discussing development of gambol, reach out on
[Mastodon](https://mast.hpc.social/@nuccitheboss), or feel free to start a
new [Discussion](https://github.com/NucciTheBoss/gambol/discussions) thread
on GitHub.

## License 📋

gambol is licensed under the GNU Affero General Public License, version 3.
Please see the [AGPLv3 LICENSE](./LICENSE) file for further details.

