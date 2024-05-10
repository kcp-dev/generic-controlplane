# Generic Control Plane

Generic Control Plane is a Kubernetes API Server flavor that is designed to be used in a standalone mode. It is a lightweight control plane that can be used for development, testing, and learning purposes.

It has all the core components of a Kubernetes control plane without container primitives.

## Usage

```
# Build the control plane
make
# Get all the options
./bin/gcp start options
# Start the control plane
./bin/gcp start
```
