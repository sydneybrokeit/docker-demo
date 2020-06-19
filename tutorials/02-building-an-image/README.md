We can build our own images to use with Docker, which is where the real power comes from.

First, change into the `build` directory, and you'll see a file name `Dockerfile`.

```
FROM golang:latest
COPY . /go/src/build/
WORKDIR /go/src/build/
RUN go build .
ENTRYPOINT ["/go/src/build/build"]
```

From within this directory, run 
```
docker build .
```

You'll see it building, step by step:
```
Sending build context to Docker daemon  3.072kB
Step 1/5 : FROM golang:latest
 ---> 5fbd6463d24b
Step 2/5 : COPY . /go/src/build/
 ---> 4271966b5ec7
Step 3/5 : WORKDIR /go/src/build/
 ---> Running in c3fd43f3a3a8
Removing intermediate container c3fd43f3a3a8
 ---> 541eeb8def39
Step 4/5 : RUN go build .
 ---> Running in ed25a452a526
Removing intermediate container ed25a452a526
 ---> b15dbafe5a46
Step 5/5 : ENTRYPOINT ["/go/src/build/build"]
 ---> Running in 0b3f4667137d
Removing intermediate container 0b3f4667137d
 ---> 9df7fa42dc44
Successfully built 9df7fa42dc44
```

Notice how it says `intermediate container $...` after each step?  That's because each step is creating a new container
and becoming a new layer.  Understanding this is important for optimizing builds, which we will get into later.

Also note how the last line says `Successfully built $...`.  This is the image ID we can use to run it.  In this case,
we can run
```
docker run 9df7fa42dc44 
```

When you do this, you'll see it print `hello world!`

## What's going on inside this file?
Dockerfiles are a series of tasks to do, using a 
[specific set of commands](https://kapeli.com/cheat_sheets/Dockerfile.docset/Contents/Resources/Documents/index).

```
FROM golang:latest
```
We are going to use the latest version of the `golang` image for our base.

```
COPY . /go/src/build/
```
We're going to copy the contents of `.` into `/go/src/build` -- note that we didn't have `mkdir` the directory, the
`COPY` command handles it automatically.

```
WORKDIR /go/src/build/
```
This is the equivalent of `cd` into the given directory within the image.

```
RUN go build .
```
This will compile everything within the `build` directory we're inside of.

```
ENTRYPOINT ["/go/src/build/build"]
```
The `ENTRYPOINT` directive tells it what command to run when we start the container - in this case, it will run the
`build` binary we compiled in the previoust step.

## Tagging
When we build the image, we get an ID that we can use to run it.  But we don't want to change this every time, nor do we
want to have to run `docker images` to find the most recent one.

Fortunately, we can tag our own images, as well as those from a source like Docker Hub!

All we have to do to tag a build is add `-t name:tag` to the build command:
```
docker build . -t name:tag
```

So let's try this!  Tag the build from this tutorial as "dockerdemo:latest" and run it!

After that, there's just one more bit of housekeeping - we can add tags to a particular image!

```
docker tag $source:tag $destination:tag
```

Try tagging `dockerdemo:latest` as `dockerdemo-build` without re-building it.