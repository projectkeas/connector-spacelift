# Spacelift Connector

The [Spacelift](https://spacelift.io) connector allows for the ingestion of Audit Trail events from Spacelift. To setup the connector, use the following [guide](https://docs.spacelift.io/integrations/audit-trail) and ensure that the URL used for the ingestion ends with `/integrations/spacelift/audit`.

## Endpoints

|Url|Methods|Description|
|---|---|---|
|`/ingest`|POST|Captures a given event into the system (assuming it passes validation and ingestion policies)|
|`/_system/health`|GET|The liveness health check endpoint|
|`/_system/health/ready`|GET|The readiness health check endpoint|

## Configuration

The Spacelift connector looks the following configuration objects within a Kubernetes cluster:

- ConfigMaps:
  - connector-spacelift-cm
- Secrets:
  - connector-spacelift-secret
  - ingestion-secret

### ConfigMap - `connector-spacelift-cm`

|Key|Description|
|---|---|
|ingestion.uri|The path of the ingestion API, including the `/ingest` suffix. Default: `http://keas-ingestion.keas.svc.cluster.local/ingest`|
|log.level|The log level that should be written to the console. Default: `debug`|
|server.port|The port to listen on. It can be useful to change this for local development. Default: `5000`|

### Secret - `connector-spacelift-secret`

_This secret is required. If the secret does not exist, the readiness checks will fail._

|Key|Description|
|---|---|
|spacelift.webhook.token|The secret used to validate that the incoming requests are indeed coming from Spacelift. This must match what's you set on the Spacelift UI|

### Secret - `ingestion-secret`

_This secret is required. If the secret does not exist, the readiness checks will fail. This secret is often setup by the [Ingestion API](https://github.com/projectkeas/ingestion)._

|Key|Description|
|---|---|
|ingestion.auth.token|The API Key that's used for authentication against the ingestion API|
