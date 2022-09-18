# DevContainer CLI managment tool

`devc` is a simple CLI tool to manage your devcontainers.

## What is a "DevContainer"?

The devcontainer concept come from Visual Studio Code and its "[Remote -
Containers](https://code.visualstudio.com/docs/remote/containers)" extension.

> Work with a sandboxed toolchain or container-based application inside (or mounted into) a container.\
– <https://code.visualstudio.com/docs/remote/containers>

In other words, it lets you install the toolchain of your project in a
container. This way you don't mess your computer with all the dependencies of
all the projects and their programming languages on which you work on. It can
also make it easier for others to start working on your projects, without
having to guess what are the required tools to develop, lint, test, build, etc.

## Install

* From sources

```
go install github.com/nikaro/devc@latest
```

Or:

```
git clone https://github.com/nikaro/devc
cd devc
make
sudo make install
```

* From pre-build binaries and packages

You can get builds for Linux, Windows and macOS, either arm64 or amd64 on the
[Releases](https://github.com/nikaro/devc/releases) page.

* From [brew](https://brew.sh)

```
brew install nikaro/tap/devc
```

* From [AUR](https://aur.archlinux.org/packages/devc-bin/)

```
yay -Syu devc-bin
```

## Usage

```
> devc --help
devc is a devcontainer managment tool

Usage:
  devc [command]

Available Commands:
  build       Build devcontainer
  help        Help about any command
  init        Initialize devcontainer configuration
  list        List devcontainers
  shell       Execute a shell inside devcontainer
  start       Start devcontainer
  stop        Stop devcontainer

Flags:
  -h, --help            help for devc
  -v, --verbose count   enable verbose output

Use "devc [command] --help" for more information about a command.
```

## Demo

[![asciicast](https://asciinema.org/a/521932.svg)](https://asciinema.org/a/521932)

## Configure Neovim

With this snippet you can make Neovim install plugins inside your container (and only inside, not on your host).

`~/.config/nvim/init.lua` with [packer](https://github.com/wbthomason/packer.nvim) as plugin manager:

```lua
-- ensure packer is installed at launch
local ensure_packer = function()
  local install_path = vim.fn.stdpath('data')..'/site/pack/packer/start/packer.nvim'
  if vim.fn.empty(vim.fn.glob(install_path)) > 0 then
    vim.fn.system({'git', 'clone', '--depth', '1', 'https://github.com/wbthomason/packer.nvim', install_path})
    vim.cmd [[packadd packer.nvim]]
    return true
  end
  return false
end

local packer_bootstrap = ensure_packer()

return require('packer').startup(function()

  -- Packer can manage itself
  use 'wbthomason/packer.nvim'

  -- Others plugins
  [...]

  -- Load devcontainer plugins
  if vim.fn.filereadable('.devcontainer/devcontainer.json') == 1 and vim.fn.filereadable('/.dockerenv') == 1 then
    local devcontainer = vim.fn.json_decode(vim.fn.readfile('.devcontainer/devcontainer.json'))
    local customs = devcontainer.customizations or {}
    local devc_customs = customs.devc or {}
    local devc_extensions = devc_customs.extensions or {}
    local devc_settings = devc_customs.settings or {}

    for _, devc_plugin in ipairs(devc_extensions) do
      use(devc_plugin)
    end

    for k, v in pairs(devc_settings) do
      if k == 'vimscript' then
        for _, script in ipairs(v) do
          vim.cmd(script)
        end
      elseif k == 'lua' then
        for _, script in ipairs(v) do
          vim.cmd('lua' .. script)
        end
      end
    end
  end

  if packer_bootstrap then
    require('packer').sync()
  end

end)
```

`.devcontainer/devcontainer.json`:

```json
{
  [...]
  "customizations": {
    "devc": {
      "extensions": [
        "fatih/vim-go"
      ],
      "settings": {
        "vimscript": [
          "let g:go_fmt_command = 'goimports'"
        ]
      }
    }
  },
  [...]
}
```
