## Requirements

To build a payments gateway through a REST API. It should simulate a payment flow.

### Endpoints

`/authorize` - Responsible for authorising a payment for a customer

Input:

* CardDetails
    * PAN
    * Expiry
    * CVV
    * Amount
        * MinorUnits
        * Currency

Response:

* Unique ID that can be used for all API Calls
* Success or error

`/void` - Cancel the whole transaction without billing the customer. No further action is possible once a transaction is
voided.

Input:

* Authorization ID

`/capture` - Capture money on the customers bank. It can be called multiple times with the amount that is not greater to
the amount authorised in the first call. e.g £10 authorisation can be captured 2 times with a £4 and £6 call.

Input:

* Authorization ID
* Amount
    * MinorUnits
    * Currency

`/refund` - Will refund the money taken from the customer bank account. It can be also called multiple times with the
amount captured. Once a refund has occured a capture cannot be made on the specific transaction.

Input:

* Authorization ID
* Amount
    * MinorUnits
    * Currency

Notes

* Amount and currency available ?? **Check what this means** - Is this the availability on the account?


## Repository
The repository follows a monorepo style approach.
It is structured as follows:

* `/github` - github setup/github actions
* `/build` - generated files such as protobuf definitions
* `/docs` - repo related docs
* `/pkg` - shared libraries for project
* `/proto` - global protobuf definitions
* `/services` - containing any microservices

### Services
All services in `/service` directory should be containerised using Docker.
In order to spin up infrastructure run `make run`.

### Protobuf
In order to generate definitions run `make proto-generate`. This will
generate code for languages specified in `buf.gen.yaml` based on
proto definitions in `/proto`. For extra language support
add in `buf.gen.yaml`


## Improvements
* Exposing the GET `/payment/{id}` endpoint to be able to fetch payment details after completion. Also would be good to expose an API to list payment actions.s 
* Move payment update/processing code out of main flow. This could be done asynchronously to avoid the chance of not
  writing to db. When a payment is created an event is produced and it is processed separately.
* Protobuf
    * CI - Responsible for protobuf generation to ensure no compatibility/versioning issues across machines.
    * `WIRE_JSON` - In order to share the protobuf schemas and avoid duplication I added the WIRE_JSON check. This was
      to avoid writing extra mapping functions. However, by doing this it stops the rpc/internal formats to be able to
      benefit from WIRE changes.
* Database
    * If more time would have written table driven tests to tidy up tests.
    * Storing card information should really be in encrypted storage
      like [Hashicorp Vault](https://www.vaultproject.io/)
* Metrics
    * To add metrics I would look at adding [promhttp](https://github.com/prometheus/client_golang/tree/master/prometheus/promhttp) to be able to instrument the HTTP handler. This would enable
      dashboards to be built to track things like latency and number of requests
* Idempotency
    * Ideally an idempotency mechanism would be implemented to prevent clients making duplicate requests. It would also
      cache results of previous calls.
      https://app.diagrams.net/#G16LSiTc8i5i_N0f7TDqM6yrpQrAbg5AdF
* Testing
  * Improve service layer tests, ran out of time to cover further edge cases
  * e2e tests
    * I would like to write e2e tests to spin up the gateway and call each endpoint validating that they work. This would be done by spinning up via docker-compose and writing BDD styled tests. The Ginkgo library is good for this.
* Production Readiness
  * All committed secrets, setup should be removed and injected in as a separate process.
  * In order to scale services accordingly Kubernetes could be used for each service
    so that they can be scaled independently and horizontally.
  * Monitors setup to track service health such as `/health`, `/metrics` endpoint as well
    as business metrics where alerts can be triggered.
