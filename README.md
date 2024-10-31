# Terraform Provider Metabase

This is a Terraform provider for managing Metabase resources.
The provider is based on the [Metabase REST API](https://www.metabase.com/docs/latest/api-documentation) and is currently compatible multiple Metabase versions.

## Usage

The documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/labbs/metabase/latest/docs).

## Metabase Compatibility

| Metabase Version | Supported |
|-----------------:|:---------:|
| 0.50             | ✅        |
| 0.51             | ✅        |

A small is present in this repository to add another version of Metabase to the compatibility list.
If the version is not supported, the provider will return an error when trying to connect to the Metabase API.

## Development

### Folders

- schema_generator: Contains the code to generate the schema from the Metabase API OpenAPI specification.
- metabase: Contains the code used to interact with the Metabase API and manage multiple versions of the API (v0_50, v0_51, ...).
- provider: Contains the code for the Terraform provider.

### Generate the schema

#### Requirements
Start a Metabase instance with a version that you want to support and run the following command:

```bash
docker run -d -p 3000:3000 --name metabase metabase/metabase:v0.50.31
cd schema_generator
go run main.go
```

This will generate the schema for the Metabase version running on `localhost:3000` and the schema in the metabase folder.

#### Update the provider

After generating the schema, you can update the provider to support the new version of Metabase.

- client.go: Integrate the new version of the Metabase client same as the existing versions.
- permissions.go, user.go, ...: Add the new resources and methods to interact with the Metabase API.
