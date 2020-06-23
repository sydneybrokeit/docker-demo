Remember how I said that Docker containers are made up of layers, and that layers are just tarballs with changes 
tracked between steps?  Every one of those tarballs has a size, and if we're not careful, our containers can balloon,
and start keeping a bunch of junk that we don't need -- or, if we're not careful, we can accidentally expand our attack
surface.  We can also greatly improve build times if we're smart.

## Selecting a proper base image
The first place we can optimize our images is actually in the first line of our Dockerfile - the `FROM` line.

Depending on what you're doing, you have a few options.

### Selecting a stripped-down distribution base
One possible starting point is to select something based on a stripped down distribution.  If you look in Docker Hub,
you can find a lot of stripped down images based on niche distributions like Alpine Linux.  These are not without
trade-offs, however.  Let's look at the trade-offs for Alpine, specifically, because it's one of the more common ones.
* Alpine uses MUSL instead of GCC, which means that some things (specifically, those that use C or C++, or interop) will
not work.  This is the biggest drawback.
* Alpine is small, but most of the base images based on it don't have much extra functionality in them.  Things like
git or other necessary tools for doing the build may need to be included as part of the Dockerfile.
* Certain things in Alpine don't work quite as we'd expect (mostly due to system level chagnes like MUSL vs GCC)

### Selecting a `slim` version
Depending on what image you're using, some projects include a `slim` version of their base images.  These are going to
be slightly larger than using a full low-overhead distribution, but are going to be slightly easier to work with
overall.  For example, the official `node` base images have `slim` options.
* These images are based off a more mainstream distribution, using GCC, so most interop will work without issue
* These images tend to have more of the nice things that we take for granted like `git`, `nslookup`, etc
* We tend to be more familiar with these distributions
* General package availability is much better, so the packages required for certain things are easier to get or install.

### Selecting an official base image
Many large projects, like NodeJS, Python, and so on, have official base images you can use.  These have a few niceties:
* These are the official images and will generally follow best practices, not only for the framework/language in use,
but also for creating Docker images (such as cleaning up after installations)
* These images are guaranteed to work with the specified version of the project.

**Note** These are not all mutually exclusive.  Sometimes the official base image can be based on something like Alpine!
The important thing is to understand the concepts and trade-offs.

## Optimizing step order
Every step in a Dockerfile creates a new layer, as we've noted.  By default, when there are no changes to be made to a
layer, Docker will simply use the existing layer, which is significantly faster (and doesn't require an update when you
`push`).  The most important thing to note here for optimization purposes is that if a previous layer changed, _every
subsequent layer must be changed as well_.

### Example
Let's say that we have the following excerpt from a Dockerfile:
```
RUN apt install git
RUN git clone repo
RUN apt install go
RUN go build .
```

We can see what this will turn into from the steps alone - install git, clone the repo, install go, and build the repo.
If anything in the repo has changed, everything afterward, including installing go, will also create a fresh layer.

```
RUN apt install git
RUN git clone repo # this will change
RUN apt install go # this will cahnge
RUN go build . # this will change
```

Let's organize this a little bit more cleanly, so that we have the things least likely to change moved up.

```
RUN apt install git
RUN apt install go
RUN git clone repo
RUN go build .
```

In this case, unless there's an update to git or go, only the git clone and the build will change -- i.e., _only the
parts that we actually think of as changing will have changed in the final image_.  **In order to avoid cascading layer
changes, put the things least likely to change toward the top where possible**.

## Optimizing the number of steps
There is a small cost associated with having layers, so there's a sweet spot of optimization we can do for that, as
well.

Let's use the example from before:
```
RUN apt install git
RUN apt install go
RUN git clone repo
RUN go build .
```

We can easily combine the installation of go and git, so that they become a single layer:
```
RUN apt install git go
RUN git clone repo
RUN go build .
```

This removes a layer just by combining a step.  Easy!

## Cleaning up after yourself
Some steps will include a cache step or intermediate files.  These can (and should!) be removed.

```
RUN apt update && apt install go
RUN apt-get clean
RUN apt-get autoclean
```

The clean and autoclean steps will get rid of some of the downloaded packages, which is all well and good, but we have
one major problem here: each step here is going to be its own layer.  In this case, let's clean it up like this:

```
RUN apt update && apt install go && apt-get clean && apt-get autoclean
```

Now this creates a single step, but more importantly, it means that the step never includes the cached files, saving
us time and space.

**IMPORTANT** - the same reason we combine steps this way is exactly why you should never include sensitive data in the image
at any point, *even if it is deleted at a later step*.  Layers are just files.  They can still be read - quite easily,
in fact!

## Multistage builds
A final optimization we can do, and one that's great for compiled languages like Go, is multi-stage builds.  We can
create a base image for compiling, then copy the resulting binary from the base to a new image, which will become the
image we are going to use or publish.  This is also the only real exception to including sensitive data.

```
FROM golang:alpine as build-env
LABEL maintainer="harold.schreckengost@smartbear.com"

ENV GOBIN /go/bin

RUN mkdir /go/src/app && \
    apk --no-cache add git curl openssh-client && \
    curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

ADD . /go/src/app
WORKDIR /go/src/app
RUN go build -o /vnc-recorder-server .

FROM jrottenberg/ffmpeg:4.0-alpine
COPY --from=build-env /vnc-recorder-server /
ENTRYPOINT ["/vnc-recorder-server"]
EXPOSE 8080
CMD [""]
```

In the above example (taken from real work I've done), inside of the `golang:alpine` container we're nameing
`build-env` within this environment, we're compiling a binary.  After this, we are creating the final base image
out of an image that includes `ffmpeg` based on Alpine (this helps avoid having to build or configure it ourselves).

The real meat of this is the first line, `FROM golang:alpine as build-env`, which specifies our first base image with
a given name, and then inside of the second `FROM` the `COPY` line contains `--from=build-env`.  We are telling Docker
to use this as the source rather than our host machine.

There are more steps we can take to optimize things, but these are good first steps, and will hopefully lay a good
foundation.

