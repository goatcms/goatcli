# GoatCLI
[![Go Report Card](https://goreportcard.com/badge/github.com/goatcms/goatcli)](https://goreportcard.com/report/github.com/goatcms/goatcli)
[![GoDoc](https://godoc.org/github.com/goatcms/goatcli?status.svg)](https://godoc.org/github.com/goatcms/goatcli)

## About
GoatCLI is a project orchestration tool. It can generate/re-generate your code, manage dependencies and sub-projects/modules. Your build scripts are too slow? GoatCLI runs it concurrency.

## Install
```
go get -u github.com/goatcms/goatcli
```

## Commands
Run the single command by **goatcli command_name** or open terminal by **goatcli terminal** power.

### Init/Start
* clone - Clone project (and modules)
* init - Initialize the new empty project
* help - Show a command-line help
* health - Show application health (check system dependencies, environment variables, versions etc)
* terminal - Open internal terminal

### Build
* build - Build goat project in the current directory
* rebuild - Clean build and run a new build
* clean - Clean built files and dependencies
* clean:build - Clean generated files only
* clean:dependencies - Clean dependencies files only

### Data      
* data:add - Add new data set to project

### Dependencies  
* deps:add - Add new static dependency like golang vendor or js node module
* deps:add:go - Add new golang dependency like 'github.com/goatcms/goatcore'
* deps:add:go:import  - Scan project dependency and add it

### Properties
* properties:get - Display a property
* properties:set - Add or update a property with a specified key

### Secrets
* secrets:get - display a secret by key
* secrets:set - add or update a secret

### Scripts - Pipelines
Pipeline (pip) is a code block run concurrency. It is run in a sandbox (like internal terminal, docker image etc).
* pip:clear - clear current pipeline context
* pip:logs - Show execution logs
* pip:run - Run script
* pip:summary - Show execution summary
* pip:wait - Wait for all tasks in context

### Scripts
* scripts:run - Run script by name

### VCS (Version Control System)
* vcs:clean - Clean vcs ignored files
* vcs:generated:list - Show generated files listing
* vcs:ignored:add - Add new vcs ignored file [--path=file path to be ignored]
* vcs:ignored:list - Show ignored files listing
* vcs:ignored:remove - Remove a vcs ignored file [--path=file path]
* vcs:scan - Scan files for changes (and add it to vcs ignored files)

## Arguments
* cwd - set Current Working Directory

## Presentations
* (PL) GoatCLI on Kariera IT Poznan 21-04-2018: https://youtu.be/YX6Ne1Z83l8 and [Slides](https://docs.google.com/presentation/d/1qaqgWtXEjiPy0CljDwsvlryFVut3fm0bQ5WJpbRIXGI/edit#slide=id.p)
* (PL) P.I.W.O (Poznańska Impreza Wolnego Oprogramowania) 28-04-2018: https://www.twitch.tv/videos/264532807 and [Slides](https://docs.google.com/presentation/d/1i4_a_G-ZvvPaZXuyajok4jlfg_lA4b8S_g3dTZv52mw/edit?usp=sharing)

### Run tests
Define the system environment to tests runs.
 * GOATCORE_TEST_DOCKER - Run docker tests if defined and equals to "YES". The tests require docker daemon started.
 * GOATCORE_TEST_TMPDIR - Is base filesystem directory.

On Linux/Unix/macOS add environment variables and run go tests by terminal.
```
export GOATCORE_TEST_SMTP_FROM_ADDRESS="test@email.com"
export GOATCORE_TEST_SMTP_TO_ADDRESS="test@email.com"
export GOATCORE_TEST_SMTP_SERVER="smtp.gmail.com:465"
export GOATCORE_TEST_SMTP_USERNAME="YOUR_TEST_USER"
export GOATCORE_TEST_SMTP_PASSWORD="YOUR_TEST_USER_PASSWORD"
export GOATCORE_TEST_SMTP_IDENTITY=""
export GOATCORE_TEST_DOCKER="YES"
go test ./...
```