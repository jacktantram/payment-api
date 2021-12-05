# Payment Gateway

## Improvements

Currently the service structure is fairly simplistic as the handler request is just being forwarded straight to the
payment processor. However, in the future if there is further business logic, such as routing to various Payment
processors it might make sense to introduce it in the application layer. 
