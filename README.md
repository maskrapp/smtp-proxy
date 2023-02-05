## Explanation

I'm running both my testing and production environments in the same Kubernetes cluster, and both my mail servers need port 25 to operate. However, the public load balancer can only forward the same port once (this is how networking works).

The TCP proxy sits between the kubernetes pod (or service) and the client that would normally connect to the mailserver directly. The tcp proxy connects to the mailserver pod over a private tailscale network. 

Read https://tailscale.com/kb/1185/kubernetes/?q=Kubernetes for more information about this.

## Installation
TODO