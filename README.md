# <img alt="Logo" width="80px" src="./contrib/logo/blue-green.png" style="vertical-align: middle;" /> Generic Control Plane

## Overview

Generic Control plane (gcp) is a Kubernetes control plane focusing on:

- A **control plane** for Kubernetes-native APIs that can be used **without Kubernetes**
- Enabling API service providers to **offer APIs centrally**

gcp can be a building block for SaaS service providers who need a robust API base platform out-of-the-box.
The goal is to be useful to cloud providers as well as enterprise IT departments offering APIs within their company.

## Documentation

To get started with gcp:

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

## Contributing

We ❤️ our contributors! If you're interested in helping us out, please check out [contributing to gcp](CONTRIBUTING.md).

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
