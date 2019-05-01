# Docker

The Dockerfile in this repo is build as part of our CI in
Cloudbuild. However, it can also be build locally.

## 0. gcloud

The Dockerfile depends, during build, on a base image that is stored
in our Google Cloud image repositories. So, you need the `gcloud`
command installed (https://cloud.google.com/sdk/install), and then you
need to do a `gcloud init` and `gcloud auth configure-docker`

## 1. Netrc

During build, docker will download, using Go, various sources that we
require. Some of these are in private github repositories, which
require authentication.

Often in development, you authenticate with github via git, which may
be using your ssh keys. However, we're not sending ssh keys into
Docker, plus Go itself speaks https to github, not git, hence this
mechanism. You may already have a `~/.netrc` file because sometimes
github seems to impose rate limits unless you have one.

If you don't already have a `~/.netrc` file then:
* In a browser, go to https://github.com/settings/tokens and click "Generate new token"
* The token needs access to private repos, so you must check the top most box: "repo"
* Give a description, then scroll to the bottom, and click "Generate token"
* A magic string will appear. You now need to create a `~/.netrc` file with the content:

      machine github.com login $MAGIC_STRING_HERE

* You probably want to `chmod 600 ~/.netrc` or similar

If you do already have a `~/.netrc` file then you can still go to
https://github.com/settings/tokens, select the token, and update its
permissions if it does not have sufficient permissions.

## 2. Docker

Install it.

## 3. Build the image

    antha-lang/antha$ docker build -t local-antha --build-arg NETRC="$(cat ~/.netrc)" .

All the sources necessary should be downloaded and compiled, and the
antha-lang/antha tests will be run. At the end, you should have a
message like:

    Step 13/13 : ONBUILD ADD . /elements
     ---> Running in cfccbe95a60c
    Removing intermediate container cfccbe95a60c
     ---> 2ce136c2261c
    Successfully built 2ce136c2261c

## 4. Run the image

With the image built, you can now run it with:

    docker run -it --rm local-antha

`-i` says interactive, `-t` says you want to interact with a terminal, and
`--rm` says delete the running container once you exit. I.e. once you exit the
shell that you're given, (either `Ctl-D` or `exit`), the running container will
be stopped and deleted and any changes lost. So you're always restarting from a
known good state.

You should be presented with a shell from which you can explore. Antha
and its commands should be installed and available in PATH as normal.

For security the Docker image does not retain netrc.  This may affect you 
running things locally so you may want to mount this file from your host 
machine, eg by calling `docker run` with  `-v ~/.netrc:/root/.netrc`.
