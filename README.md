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

## Considerations

* Authorize Endpoint
    * In order to authorize requests the card information is sent across to trigger payments. For now I have chosen to
      only store the first six and last four in the database as a reference. If there was a requirement to persist this
      I would look at tokenizing the card and persisting it in encrypted storage like
      [Hashicorp Vault](https://www.vaultproject.io/)

## Improvements

* Protobuf
    * CI - Responsible for protobuf generation to ensure no compatibility/versioning issues across machines.
    * `WIRE_JSON` - In order to share the protobuf schemas and avoid duplication I added the WIRE_JSON check. This was
      to avoid writing extra mapping functions. However by doing this it stops the rpc/internal formats to be able to
      benefit from WIRE changes.
* Database
    * If more time would have written table driven tests to tidy up tests.

https://app.diagrams.net/#G16LSiTc8i5i_N0f7TDqM6yrpQrAbg5AdF
