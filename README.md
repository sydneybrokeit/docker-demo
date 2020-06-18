# What is Docker
[Docker](https://docker.com) is one of several options for running containers and containerized applications.

Docker is the most common containerization framework, and is used in other technologies like Kubernetes and is commonly
found in other contexts, such Functions as a Service (Faas).

## What's a container?
A container is a self-contained, sandboxed filesystem that contains everything needed to run an application in its own
space.

### So it's a Virtual Machine?
No, but it's not far off.  Virtual Machines run a full copy of the operating system, and have more strictly defined
resource use, while Docker containers have just enough of the operating system to be functional within the context of
the host system's kernel, and resource utilization is more dynamic, similar to running an application.

Docker containers use certain features of the Linux kernel to have applications running in their own space within the
confines of the host kernel.  Processes within containers actually show up within the host's process tree.

### Advantages
Using Docker containers offer us some nice advantages!
* Consistent, distributable way to run and develop complex software
* Immutable application stacks
* Quickly start from a known working point with Docker Hub images like Node, Python, and Golang

### Disadvantages
* Can complicate debugging if you're not careful
* There is some overhead (but less than a VM)
* Windows and MacOS users have higher overhead

# Getting started
To get started, [install docker on your system](00-install)

## Tutorials
1. Running an image