How can we easily share our images so that other people can use them, or so that other tools can use them?  The
basic answer to this is a Docker registry, but the most common registry is [Docker Hub](https://hub.docker.com)

## Creating an account
When you go to Docker Hub, you will be presented with a sign-up form on the right.  Fill that out and you're done!

## Docker Login
To use certain features of Docker Hub from within Docker, you need to be logged in from within Docker itself.

On Linux, this is just a matter of running `docker login` and filling in your credentials.

On Windows or Mac (where Docker for Desktop is more matured), you can use `docker login` or you can right-click the
Docker Desktop icon, and select `Sign In`.  Fill in your credentials, and you're logged in.

## Using Docker Hub from Docker
Pull private repositories using the same `docker pull` command we've been using.  Now that we're authenticated to
Docker Hub from within Docker, we have access to all the things our user account account is tied to.

We also now have the option to push to Docker Hub.  Assuming we have a repository named `hmschreck/testrepo` and we
have write access to this repository, we can push our image to it.

In this example, our local image is named `localimage:latest`.

```
docker tag localimage:latest hmschreck/testrepo:latest
docker push hmschreck/testrepo:latest
```

A lot of text that looks similar to what we see when we pull an image, and we're done. Our image is now pushed to Docker
Hub, and people with access (either the public or people with access to a private repository) can make use of it.

