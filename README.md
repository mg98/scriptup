<h1 align="center">scriptup ⬆️</h1>
<p align="center">
    A migration tool for shell script executions.
</p>

<p align="center">
  <a href="https://github.com/mg98/scriptup/actions/workflows/test.yml">
    <img src="https://github.com/mg98/scriptup/actions/workflows/test.yml/badge.svg">
  </a>
  <a href="https://pkg.go.dev/github.com/mg98/scriptup">
    <img src="http://img.shields.io/badge/godoc-reference-blue.svg"/>
  </a>
  <a href="https://codecov.io/gh/mg98/scriptup">
    <img src="https://codecov.io/gh/mg98/scriptup/branch/main/graph/badge.svg?token=R3OYXX1HC7">
  </a>
  <a href="https://goreportcard.com/report/github.com/mg98/scriptup">
    <img src="https://goreportcard.com/badge/github.com/mg98/scriptup">
  </a>
  <a href="./LICENSE">
    <img src="https://img.shields.io/github/license/mg98/scriptup">
  </a>
</p>

<hr>

_scriptup_ is a framework- and language-agnostic tool that enables software projects to maintain code-based migration scripts.
This is very similar to database migration tools, just that it is based on shell instead of SQL executions.
While traditional SQL migration frameworks aid deploying database changes to production systems and to other developer machines,
_scriptup_ is able to perform state migrations using general-purpose code!
Example use cases include complex data alterations in the database (e.g. encryption) or the installation of IDE extensions
or even git hooks on developer machines.


## Install and Setup

The tool is run as a standalone binary.
You can find the appropriate executable for your OS on the [Releases](https://github.com/mg98/scriptup/releases) page.
However, macOS users are very welcome to install it via the package manager [brew](https://brew.sh).

```sh
brew tap mg98/homebrew-tap
brew install scriptup

# Verify your installation
scriptup -v
```

To setup _scriptup_ in your project, you have to create its configuration file in your project's root directory.
Take [scriptup.yml](./scriptup.yml) as a template and adjust the values as needed.

## Usage

Create a new migration, e.g.

```sh
scriptup new add-git-hook
```

This will create a new file in your configured migration folder, something like `20220626135412_add-git-hook.sh`.
Edit the file to contain the script that you want to be executed when running this migration.

As you will see from the template, the file is subdivided into two sections.
Everything after `### migrate up ###` and before (an optional) `### migrate down ###` will be run on `scriptup up`.
Everything after `### migrate down ###` will be executed on `scriptup down` and is supposed to undo the changes 
performed through the script in the _up_-section (your responsibility though).

### CLI Reference

```
COMMANDS:
   new, n     Generate a new migration file
   up, u      Execute recent scripts that have not been migrated yet
   down, d    Undo recently performed migrations
   status, s  Get status about open migrations and which was run last
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --env value, -e value  specify which configuration to use (default: "dev")
   --help, -h             show help (default: false)
   --version, -v          print the version (default: false)
```