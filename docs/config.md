<!--
 Copyright 2022 Kevin Ledesma
 
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 
     http://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->

# Config file

The config file contains the configuration for the application. It is a JSON file that is parsed by the CLI.

In the [configs](../configs)'](../configs) folder, you can find some of all the default config files.

## Table of Contents

- [Config file](#config-file)
  - [Table of Contents](#table-of-contents)
  - [Config file structure](#config-file-structure)
    - [Image](#image)
      - [Package manager](#package-manager)
    - [Podman](#podman)
    - [Container](#container)
      - [Binds](#binds)
    - [Commands](#commands)
    - [Packages](#packages)
    - [Development packages](#development-packages)
    - [User packages](#user-packages)
    - [Experiments](#experiments)

<!-- Index ends -->
## Config file structure

The config file contains the following sections:

- `image` - Contains the image configuration, such as the image name and the package manager to use
- `podman` - Contains the podman configuration, such as the podman path
- `container` - Contains the container configuration, such as the container name and the container ports
- `commands` - Contains the user-defined commands to run in the container
- `packages` - Contains the packages to install in the container
- `devpackages` - Contains the development packages to install in the container
- `userpkgs` - Contains the packages that should be installed as a user in the container
- `experiments` - Contains the experimental features to enable

### Image

The `image` section contains the following fields:

- `URI` - Which is the URI of the image to use, for example, `docker.io/library/ubuntu:latest`
- `oncreation` - This is a list of commands to run when the container is just created
- `onfinish` - This is a list of commands to run when the container has finished building
- `pkgmanager` - Contains the configuration for the package manager to use in the container
- `variables` - Contains the environment variables to set in the container

#### Package manager

The `pkgmanager` section uses base strings, where:

- The `{args}` string is replaced with the arguments passed to the package manager
- Anything else inside `{}` is removed if the operation isn't set to auto-install

So, taking that into account, a valid package manager configuration would look like this:

```json
{
    ...
    "pkgmanager": {
        "operations": {
        "install": "apt-get install {-y} {args}",
        "remove": "apt-get remove {-y} {args}",
        "update": "apt-get update",
        "upgrade": "apt-get upgrade -y",
        "search": "apt-cache search {args}",
        "clean": "apt-get clean {-y}"
        },
        ...
    }
    ...
}
```

The only supported commands are:

- `add` - Installs the packages passed as arguments
- `del` - Removes the packages passed as arguments
- `update` - Updates the package manager's database (or package)
- `upgrade` - Upgrades the packages installed in the container
- `search` - Searches for packages in the package manager's database
- `clean` - Cleans the package manager's cache

The package manager also supports adding prefixes or suffixes to package names, for example, in the nix config it is used to add the `nixpkgs.` prefix to the package name.

```json
{
    ...
    "pkgmanager": {
        "operations": {
            "add": "nix-env -iA {args}",
            "del": "nix-env -e {args}",
            "update": "nix-channel --update {args}",
            "upgrade": "nix-env -u {args}",
            "search": "nix-env -qaP {args}",
            "clean": "nix-collect-garbage"
        },
        "modifiers": {
            "add": "nixpkgs.{package}"
        }
        ...
    }
    ...
}
```

The `modifiers` section is optional, and it is used to add prefixes or suffixes to the package name. The key name is the operation name we want to modify, and the value is the modifier to use. The `{package}` string is replaced with the package name.

### Podman

The `podman` section contains the following fields:

- `path` - This is the path to the podman executable (which can also be `docker`)
- `args` - Contains the arguments to pass to the podman executable (for `podman run`)
- `rootless` - Informs the CLI if the podman executable is rootless or not
- `onlybuild` - Creates the container and after finishing doing its thing, it gets deleted
- `privileged` - Runs the container in privileged mode

### Container

The `container` section contains the following fields:

- `name` - Which is the name of the container
- `workdir` - Contains the working directory to use in the container
- `rootuser` - Uses the root user in the container
- `binds` - Contains the binds to mount in the container
- `ports` - Contains the ports to expose in the container
- `volumes` - Contains the volumes to mount in the container
- `sharedfolders` - Contains the shared folders to mount in the container

#### Binds

The `binds` section contains mainly three fields:

- `xorg` - Mounts the X11 socket in the container
- `dev` - Mounts the `/dev` folder in the container (important for GPU support)
- `variables` - Mounts the environment variables in the container

### Commands

The `commands` section defines commands you can run using `develbox run <command>`. For defining commands, it uses a dictionary of key-value pairs, where the key is the command name and the value is the command to run.

The value can be a string or a list of strings, where each string is a command to run. If the value is a list of strings, the commands will be run in order. It also supports some features:

- `!<command key>` - Runs the command defined in the `commands` section (only works if `!` is the first character)
- `$<env var>` - Replaces the environment variable with its value (work anywhere in the string)
- `~/` - Replaces the `~/` character with the home directory of the user (work anywhere in the string)
- `${command}` - Replaces the command with the output of the command (work anywhere in the string)
- `#<command>` - Runs the command as root (only works if `#` is the first character)

An example of a command would be:

```json
{
    ...
    "commands": {
        "test": "echo $HOME",
        "test2": [
            "echo $HOME",
            "#echo $USER"
        ],
        "test3": "echo ${echo $HOME}",
        "test4": "!test",
    }
    ...
}
```

### Packages

The `packages` section contains the packages to install in the container. It uses a list of strings, where each string is a package to install.

An example would be:

```json
{
    ...
    "packages": [
        "git",
        "vim",
        "python3",
        "python3-pip"
    ]
    ...
}
```

### Development packages

The `devpackages` section contains the development packages to install in the container. It uses a list of strings, where each string is a package to install.

Any package here won't be included in the Dockerfile unless the `--dev` flag is passed to the CLI.

An example would be:

```json
{
    ...
    "devpackages": [
        "build-essential",
        "cmake",
        "clang"
    ]
    ...
}
```

### User packages

The `userpkgs` section contains the packages to install as a user in the container. It includes a `packages` and `devpackages` field, which are the same as the `packages` and `devpackages` sections.

Currently, no provided configuration makes use of this section.

An example would be:

```json
{
    ...
    "userpkgs": {
        "packages": [
            "git",
            "vim",
            "python3",
            "python3-pip"
        ],
        "devpackages": [
            "build-essential",
            "cmake",
            "clang"
        ]
    }
    ...
}
```

Packages here will make the package manager inside develbox run the operation as user instead of root. Any packages here won't be included in Dockerfiles.

### Experiments

The `experiments` section contains the experimental features to enable. It uses a key-value dictionary, where the key is the name of the feature and the value is a boolean indicating if the feature should be enabled or not.

The only feature currently supported is `sockets`, enables package installations from inside the container (when running as user).

An example would be:

```json
{
    ...
    "experiments": {
        "sockets": true
    }
    ...
}
```
