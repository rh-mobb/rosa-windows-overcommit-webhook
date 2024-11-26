# Summary

This webhook is used to validate licensing requirements on ROSA to ensure that windows virtual machines
do not exceed the total capacity of the windows nodes provided.

> **WARN** this was created as a proof-of-concept only.  Please use at your own risk.


## Limitations

There are several known limitations to the webhook:

1. Only accounts for `CREATE` requests and not `UPDATE` requests
2. Logic between `domain.cpu` and `requests/limits` has not yet been determined (see #1)
3. Only `VirtualMachine` and `VirtualMachineInstance` types are validated (there may be more)
4. Depends on node labels via `WEBHOOK_NODE_LABEL_KEY` and `WEBHOOK_NODE_LABEL_VALUES` input.  If nodes are 
missing labels, they will not be used to calculate the total capacity for windows nodes in the cluster.
5. Validation happens prior to scheduling.
6. Test manifests exist in the `manifests/test` directory.

> **WARN** be advised that the test manifests contain passwords in cleartext for testing only.  This in not
> intended to be for production use and was simply used to validate the proof-of-concept.


## Usage

This is simple usage for the webhook.  Please review the `manifests/deploy.yaml` file for accuracy for your
environment.  You will need to clone this repository and change to the cloned repository directory to run 
through these instructions.

1. Webhooks need their own set of certificates in order to properly function.  You can create your own with a simple
script provided in this repository:

```bash
make certs
```

If you do not use the script, be sure to create the following files:

* **CA Certificate** - `tmp/ca.crt` - The CA certificate used to sign the webhook web certificate.
* **Webhook Key** - `tmp/server.key` - The webhook server key.
* **Webhook Cert** - `tmp/server.crt` - The webhook server certificate, signed by the `tmp/ca.crt` file (above).  It 
should be noted that the certificate must be requested with `windows-overcommit-webhook.windows-overcommit-webhook.svc`
as the common name and/or subject alternative name, as this is the name the Kubernetes API expects.  Please see 
the `scripts/gen-certs.sh` for an example.


2. Create the webhook in the ROSA cluster.  This step assumes you have a functioning ROSA cluster and your 
`KUBECONFIG` configured to run commands against the cluster:

```bash
make create
```


3. Once you have installed OpenShift Virtualization (see https://cloud.redhat.com/experts/rosa/ocp-virt/with-fsx/
for a quick start) your requests for `VirtualMachine` and `VirtualMachineInstances` will be successfully validated
by the webhook.


## Cleanup

1. To cleanup the webhook configuration, and the deployment, namespace and certificates, simply run:

```bash
make destroy
```
