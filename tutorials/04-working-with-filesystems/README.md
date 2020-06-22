What if we need to have some sort of persistence for a container?  What if we need to be able to work with the files
a container uses?

Let's start off with something simple.

In this directory, there's a `test.sh` file:

```
#!/bin/bash
echo "Test"
while :
do
    echo $(date) >> /log/output
    sleep 1
done
```

This is a pretty simple script.  It will log the output of `date`, wait 1 second, then do it again.

Let's build our imagee.  For this example, we'll use tag it as `files`

```
docker build . -t files
```

So now we can run the image

```
docker run -d files
```

And if we get a shell on this container, we can `tail -f` the file `/log/output`
```
docker exec -it festive_davinci /bin/bash
root@ae4b30a72468:/# tail -f /log/output
```

And we can see that the file is still being written to.  But what if we want to do some file processing, such as
video transcoding, parsing a file and printing some output, or something along those lines?  Surely there's a better
way to do this than trying to send files over the network?

## Using bind mounts
The easiest way (and the one we're going to go over here) is *bind* mounting.  If you're familiar with bind mounts from
Linux, these work similarly.  If not, these will allow us to make a directory on the host available within the
container.

This time, let's try telling Docker to start the container with a bind mount!

```
docker run -d --mount type=bind,source="$(pwd)"/log,target=/log files
```

There are other types of mounts, but binds are the easiest to figure out.  So that's the `type` we'll use.

The source is the location on the host machine (in this case, it takes the present working directory `$(pwd)`, and 
appends /log).  The target is the directory *inside* the container to "mount" to.  That's it!

And now when we `tail -f` the output, it functions exactly as we think it should, but it's being written to a *local*
file from within the container!