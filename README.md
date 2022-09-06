# Padawan

Padawan is a PaaS built on top of Kubernetes. It works as follows:

Each user can create a single container from a set Container images that has been defined by administrators. The user can then ssh into the container and run whatever they want. Users can then delete their container and create a new one.

Padawan is handled using a cli tool called `padawan`. This tool can be found in the `padawan-cli` repository.

## Implementation details

Padawan is built to run behind a reverse proxy that handles authentication and passes a X-Fowarded-User header to the Padawan server. This header contains the username of the user that is currently logged in.

To provide ssh functionnality, the pod containing the user's container contains a sidecar ssh container that is responsible for handling ssh connections and establishing a shell inside the user's container using `kubectl exec`.

## Deployment

A helm chart is provided to deploy on a kubernetes cluster. The chart relies on Traefik and Cert-Manager but contains the OAuth2 reverse proxy using the `traefik-forward-auth` project.
