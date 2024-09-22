# Generic Control Plane

# IMPORTANT:

Generic Control Plane is a new project under KCP umbrella and is not yet ready for production use. We are actively working on the project and welcome contributions from the community. If you are interested in contributing, please see our [contributing guide](CONTRIBUTING.md).

## Overview

Generic Control plane is a Kubernetes based control plane focusing on:

- A **control plane** for Kubernetes-native APIs that can be used **without Kubernetes**
- Enabling API service providers to **offer APIs centrally**

gcp can be a building block for SaaS service providers who need a robust API base platform out-of-the-box.
The goal is to be useful to cloud providers as well as enterprise IT departments offering APIs within their company.

## Documentation

To get started with generic control plane:

```bash
# Clone the repository
git clone https://github.com/kcp-dev/generic-controlplane.git

# Build the project
cd generic-controlplane && make build

# Start standalone gcp
./bin/gcp start

# Access the gcp API
export KUBECONFIG=.gcp/admin.kubeconfig

# Check resources in the gcp API
kubectl api-resources
```

## Batteries

Example server contains a simple implementation of batteries that can be used to extend the gcp API.

Batteries:
- `leases` - a Kubernetes lease resources from `coordination.k8s.io`
- `authentication` - Kubernetes native authentication using `authentication.k8s.io`
- `authorization` - Kubernetes native authorization using `authorization.k8s.io`
- `admission` - Kubernetes native admission using `admissionregistration.k8s.io`
- `flowcontrol` - Kubernetes native flow control using `flowcontrol.apiserver.k8s.io`
- `apiservices` - Kubernetes native API services using `apiregistration.k8s.io`


When starting server without any flags, local in-memory etcd will be used and batteries will be disabled by default.

Important: In the long run, we plan to move existing apis into batteries on its own, and make default server to be a simple server without any resources.

To start the server with batteries enabled, use the following flags:

```bash
./bin/gcp start --batteries=lease,authentication,authorization,admission,flowcontrol
```


## Contributing

We ❤️ our contributors! If you're interested in helping us out, please check out [contributing to Generic Control Plane](CONTRIBUTING.md).

This community has a [Code of Conduct](./code-of-conduct.md). Please make sure to follow it.

## Getting in touch

There are several ways to communicate with us:

- The [`#kcp-dev` channel](https://app.slack.com/client/T09NY5SBT/C021U8WSAFK) in the [Kubernetes Slack workspace](https://slack.k8s.io).
- Our mailing lists:
    - [kcp-dev](https://groups.google.com/g/kcp-dev) for development discussions.
    - [kcp-users](https://groups.google.com/g/kcp-users) for discussions among users and potential users.
- By joining the kcp-dev mailing list, you should receive an invite to our bi-weekly community meetings.
- See recordings of past community meetings on [YouTube](https://www.youtube.com/channel/UCfP_yS5uYix0ppSbm2ltS5Q).
- The next community meeting dates are available via our [CNCF community group](https://community.cncf.io/kcp/).
- Check the [community meeting notes document](https://docs.google.com/document/d/1PrEhbmq1WfxFv1fTikDBZzXEIJkUWVHdqDFxaY1Ply4) for future and past meeting agendas.
- Browse the [shared Google Drive](https://drive.google.com/drive/folders/1FN7AZ_Q1CQor6eK0gpuKwdGFNwYI517M?usp=sharing) to share design docs, notes, etc.
    - Members of the kcp-dev mailing list can view this drive.

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkcp-dev%2Fkcp.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkcp-dev%2Fkcp?ref=badge_large)
