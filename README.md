# Develbox

A tool to create dev containers using Podman and Go. In other words, like JetPack's Devbox but worse.

## Why this exists

After reading about [JetPack's Devbox](https://github.com/jetpack-io/devbox) on HackerNews, I wanted to use it to prevent [supply-chain attacks](https://www.bleepingcomputer.com/news/security/npm-supply-chain-attack-impacts-hundreds-of-websites-and-apps/) when installing packages from NPM or PyPi.

The issue was that Nix isn't available on Fedora Silverblue, as it currently doesn't comply with the [FHS](https://en.wikipedia.org/wiki/Filesystem_Hierarchy_Standard) (_Nix directory is located on `/nix`, but `/` is read-only on [Silverblue](https://docs.fedoraproject.org/en-US/fedora-silverblue/technical-information/#filesystem-layout)_) and has issues with [SELinux](https://github.com/NixOS/nix/issues/2374). There are other ways to install it, like:

- [Matthewpi's guide](https://gist.github.com/matthewpi/08c3d652e7879e4c4c30bead7021ff73) to manually installing Nix
- Changing to another distro, like [NixOS](https://nixos.org/) directly
- [Yajo's fork](https://github.com/moduon/nix-installers/tree/rpm-ostree) (still hasn't been [merged](https://github.com/nix-community/nix-installers/pull/8) into [upstream](https://github.com/nix-community/nix-installers))
- [Nix User Chroot](https://github.com/nix-community/nix-user-chroot) (I became aware of this one too late into the project)
- [Nix Portable](https://github.com/DavHau/nix-portable) (Just discovered this one too)

> _I could instead have used developments containers with VSCode now that I think about it._

So instead of choosing all of those good options, I decided to try to make a script that does something similar to [Devbox](https://github.com/jetpack-io/devbox) and [Toolbox](github.com/containers/toolbox) but worse just to learn Golang.

## What does that mean?

Mainly, it means that the code quality here probably is not the best. There are a ton of things that should be improved or made more secure.

I would fix them but my current knowledge of coding, containers, and Linux isn't enough to make this a better product.

**TL;DR: It's best to not use this in your projects, as it is made for personal use.**

## Getting Started

If you are okay with all that, you can install this script by doing the following. First, we should visit the prerequisites.

### Prerequisites

This project requires you to have installed

- [Go](https://go.dev/)
- [Podman](https://podman.io/) or [Docker](https://www.docker.com/), ideally use Podman.

You can get them from your package manager usually with as `go` and `podman`.

### Installing

You can get the project using `go install` like this.

```bash
go install github.com/kadmuffin/develbox@latest
```

### Usage

> It's recommended that you add `.develbox/home` to your`.gitignore` file.

#### Creating the container

You can now proceed to create the container with the following command.

```bash
develbox create
```

Configs file will be located at `.develbox/config.json`

#### Opening the shell

After that you can enter the container using:

```bash
develbox enter
```

And stop the container with:

```bash
develbox stop
```

You can delete it using `develbox trash` too.

#### Managing packages

To add a package to the container we can run `develbox add`, for example, if we wish to add `nano` to the container:

```bash
develbox add nano
```

Now, if we want to delete the package, we use the `develbox del` command:

```bash
develbox del nano
```

## Contributing

If you wish to contribute to this small repo, you are welcome to submit your pull request. Take into account that I'm a total noob at this, so explanations and patience are appreciated!

## License

This project is under the [Apache 2.0 License](https://github.com/kadmuffin/develbox/blob/main/LICENSE).

## Acknowledgments

- **Billie Thompson** - _Provided README Template_ -
    [PurpleBooth](https://github.com/PurpleBooth)
- **Jetpack's Devbox** - Inspiration - [Devbox](https://github.com/jetpack-io/devbox)
- **Toolbox** - Used as reference for some things (for example, how to load /dev) - [Toolbox](https://github.com/containers/toolbox)
- **Martin Viereck** - Helpful Wiki - [x11docker](https://github.com/mviereck/x11docker)
- And even more projects in [CREDITS](https://github.com/kadmuffin/develbox/blob/main/CREDITS).
