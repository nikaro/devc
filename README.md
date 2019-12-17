# DevContainer CLI managment tool

## What is a "DevContainer"?

The devcontainer concept have been developped by the authors of Visual Studio Code and its "[Remote - Containers](https://code.visualstudio.com/docs/remote/containers)" extension.

> Work with a sandboxed toolchain or container-based application inside (or mounted into) a container.
– <https://code.visualstudio.com/docs/remote/remote-overview>

This way you don't mess your computer with all the dependencies of all the projects and their programming languages on which you work on.
It can also make it easier for others to start working on your projects, without having to guess what are the required tools to develop, lint, test, build, etc.

## Limitations

It currently (and probably forever if nobody send a pull request for it) does not support devcontainer that uses Dokerfile only. A Docker-Compose file is mandatory.

## Install

There are different methods to install `devc`, ordered by preference.

To install from the devc devcontainer (requires: docker, docker-compose, make):

```
$ git clone https://git.karolak.fr/nicolas/devc.git
$ cd devc
$ docker-composer -p devc_devcontainer -f .devcontainer/docker-compose.yml up -d
$ docker-composer -p devc_devcontainer -f .devcontainer/docker-compose.yml exec app bash
(container) $ make
(container) $ exit
$ sudo make instal
```

To install from sources into `/usr/local/bin/` (requires: golang, make):

```
$ git clone https://git.karolak.fr/nicolas/devc.git
$ cd devc
$ make
$ sudo make install
```

To install `devc` in your `GOPATH` (requires: golang):

```
$ go get -u git.karolak.fr/nicolas/devc
```

## Usage

```
$ devc --help
A CLI tool to manage your devcontainers using Docker-Compose

Usage:
  devc [command]

Available Commands:
  build       Build or rebuild devcontainer services
  completion  Generate completion script
  down        Stop devcontainer services
  exec        Execute a command inside a running container
  help        Help about any command
  ps          List containers
  shell       Execute a shell inside the running devcontainer
  up          Start devcontainer services

Flags:
  -f, --file string           alternate Compose file
  -h, --help                  help for devc
  -p, --project-name string   alternate project name
  -P, --project-path string   specify project path

Use "devc [command] --help" for more information about a command.
```

## Demo

[![asciicast](https://asciinema.org/a/kkM3UIF6YDg8tWjjx1MJgLl6z.svg)](https://asciinema.org/a/kkM3UIF6YDg8tWjjx1MJgLl6z)<Paste>

## Configure (Neo)Vim

With this configuration you can make (Neo)Vim install plugins inside your container (and only inside, not on your host).


`~/.config/nvim/init.vim` with vim-plug as plugin manager:

```vimscript
if &compatible
  set nocompatible
endif

call plug#begin('~/.local/share/nvim/site/plugged')

Plug 'scrooloose/nerdtree'
[...]

" detect wether we are in a docker container or not
function! s:IsDocker()
  if filereadable('/.dockerenv')
    return 1
  endif
  if filereadable('/proc/self/cgroup')
    let l:cgroup = join(readfile('/proc/self/cgroup'), ' ')
    let l:docker = matchstr(l:cgroup, 'docker')
    if l:docker != ""
      return 1
    endif
  endif
endfunction

" install devcontainer plugins if exist and we are in a container
if filereadable('.devcontainer/devcontainer.json') && s:IsDocker()
  let devcontainer = json_decode(readfile('.devcontainer/devcontainer.json'))
    for plugin in get(devcontainer, 'vim-extensions')
      Plug plugin
    endfor
endif

call plug#end()

[...]
```

`.devcontainer/devcontainer.json`:

```json
{
    [...]
    "vim-extensions": [
        "fatih/vim-go"
    ],
    [...]
}
```

And take a look at my [docker-compose.yml](/nicolas/devc/src/branch/master/.devcontainer/docker-compose.yml) and [Dockerfile](/nicolas/devc/src/branch/master/.devcontainer/Dockerfile) (based on <https://hub.docker.com/r/nikaro/debian-dev>) to see how to configure your containers.

## TODO

- add tests
- add an init command to create boilerplate `.devcontainer/`
- auto-start devcontainer on shell command
