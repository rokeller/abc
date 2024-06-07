# abc - Azure Blob Commands

A simple CLI to mess with Azure blobs.

## Usage

`abc` is structured using commands and sub-commands. The command hierarchy is
as follows.

* `containers` - do stuff with containers
  * `ls` - list containers on a storage account
  * `mk` - create containers on a storage account
* `blobs` - do stuff with blobs in containers
  * `ls` - list blobs in a container
  * `rm` - remove blobs from a container

You can use the `-h` on any command to get more information:

```bash
abc blobs ls -h
```

## Authentication

Authentication is attempted automatically according to the options listed in
[Azure authentication with the Azure Identity module for Go
](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication?tabs=bash#2-authenticate-with-azure)

Running `abc` is most easy on Azure workloads that have a managed identity
assigned to them (see options 2 and 3 on the page linked above). Using a service
principal with a secret is almost as easy, using the `AZURE_CLIENT_ID`,
`AZURE_TENANT_ID` and `AZURE_CLIENT_SECRET` environment variables (option 1 on
the page linked above).

## Examples

The following examples assume authentication is taken care of.

### List blobs in a container

```bash
abc blobs ls -n mystorageaccount -c mycontainer
```

This will list each blob on a separate line in the output.

### Remove all but the last three blobs in a container

```bash
abc blobs rm -n mystorageaccount -c mycontainer $(abc blobs ls -n mystorageaccount -c mycontainer | head -n-3)
```
