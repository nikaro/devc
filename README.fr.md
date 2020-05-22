# Outil de gestion de DevContainer en ligne de commande

## Qu'est-ce qu'un "DevContainer" ?

Le concept de devcontainer a été développé par les auteurs de Vistual Studio Code et de son extension "[Remote - Containers](https://code.visualstudio.com/docs/remote/containers)".

> Work with a sandboxed toolchain or container-based application inside (or mounted into) a container.
– <https://code.visualstudio.com/docs/remote/remote-overview>

Ainsi vous ne pourrissez pas votre ordinateur avec toutes les dépendences de tous vos projets et leur langages de programmation sur lesquels vous travaillez.
Ça peut aussi permettre de rendre plus facile le démarrage du travail sur vos projets, en évitant d'avoir deviner les outils requis pour développer, tester, construire, etc.

## Limitations

Actuellement ça ne supporte pas les devcontainers qui repose un fichier Dockerfile uniquement. L'utilisation d'un fichier Docker-Compose est obligatoire.
Aussi, toutes les paramètres `devcontainer.json` ne sont pas supportés.

## Installation

Il y a différentes méthodes d'installation possible de `devc`, par ordre de préférence.

Sur ArchLinux, depuis [AUR](https://aur.archlinux.org/packages/devc/) :

```
$ yay -Syu devc
```

Pour installer avec le devcontainer de "devc" (requiert : docker, docker-compose, make) :

```
$ git clone https://git.sr.ht/~nka/devc
$ cd devc
$ docker-composer -p devc_devcontainer -f .devcontainer/docker-compose.yml up -d
$ docker-composer -p devc_devcontainer -f .devcontainer/docker-compose.yml exec app bash
(container) $ make
(container) $ exit
$ sudo make install
```

Pour installer depuis les sources (requiert : golang, make) :

```
$ git clone https://git.sr.ht/~nka/devc
$ cd devc
$ make
$ sudo make install
```

Pour installer dans votre `GOPATH` (requiert : golang) :

```
$ go get -u git.sr.ht/~nka/devc
```

Ou depuis GitHub :

```
$ go get -u github.com/nikaro/devc
```

## Utilisation

```
$ devc --help
A CLI tool to manage your devcontainers using Docker-Compose

Usage:
  devc [command]

Available Commands:
  build       Build or rebuild devcontainer services
  completion  Generate completion script
  down        Stop and remove devcontainer containers, networks, images, and volumes
  exec        Execute a command inside a running container
  help        Help about any command
  ps          List containers
  shell       Execute a shell inside the running devcontainer
  start       Start devcontainer services
  stop        Stop devcontainer services
  up          Create and start devcontainer services

Flags:
  -f, --file string           alternate Compose file
  -h, --help                  help for devc
  -p, --project-name string   alternate project name
  -P, --project-path string   specify project path
  -v, --verbose               show the docker-compose command

Use "devc [command] --help" for more information about a command.
```

## Démonstration

[![asciicast](https://asciinema.org/a/kkM3UIF6YDg8tWjjx1MJgLl6z.svg)](https://asciinema.org/a/kkM3UIF6YDg8tWjjx1MJgLl6z)<Paste>

## Configurer (Neo)Vim

Avec cette configuration vous faire en sorte que (Neo)Vim installe les plugins à l'intérieur de votre conteneur (et uniquement à l'intérieur, pas sur votre machine hôte).


`~/.config/nvim/init.vim` avec vim-plug comme gestionnaire de plugins :

```vimscript
if &compatible
  set nocompatible
endif

call plug#begin('~/.local/share/nvim/site/plugged')

Plug 'scrooloose/nerdtree'
[...]

" détecte si on est dans conteneur docker ou pas
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

" si on est dans un conteneur et qu'une configuration existe
if filereadable('.devcontainer/devcontainer.json') && s:IsDocker()
  let devcontainer = json_decode(readfile('.devcontainer/devcontainer.json'))
    " installe les plugins
    for plugin in get(devcontainer, 'vim-extensions', [])
      Plug plugin
    endfor
	" applique les configurations
    for setting in get(devcontainer, 'vim-settings', [])
      execute setting
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
	"vim-settings": [
	  "let g:go_fmt_command = 'goimports'"
	],
    [...]
}
```

Et jetez un coup d'oeil à mon [docker-compose.yml](/nicolas/devc/src/branch/master/.devcontainer/docker-compose.yml) et [Dockerfile](/nicolas/devc/src/branch/master/.devcontainer/Dockerfile) (basé sur <https://hub.docker.com/r/nikaro/debian-dev>) pour voir comment configurer vos conteneurs.
