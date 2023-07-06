# Using the Ubuntu extension with Gardener as end-user

The [`core.gardener.cloud/v1beta1.Shoot` resource](https://github.com/gardener/gardener/blob/master/example/90-shoot.yaml) declares a few fields that must be considered when this OS extension is used.

In this document we describe how this configuration looks like and under which circumstances your attention may be required.

## AWS VPC settings for Ubuntu workers

Gardener allows you to create Ubuntu based worker nodes by:
1. Using a Gardener managed VPC
2. Reusing a VPC that already exists (VPC `id` specified in [InfrastructureConfig](https://github.com/gardener/gardener-extension-provider-aws/blob/master/docs/usage/usage.md#infrastructureconfig)]

If the second option applies to your use-case please make sure that your VPC has enabled **DNS Support**. Otherwise Ubuntu based nodes aren't able to join or operate in your cluster properly.

**[DNS](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-dns.html)** settings (required):

- `enableDnsHostnames`: true
- `enableDnsSupport`: true

