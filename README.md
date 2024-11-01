# Go Context examples

1. Basic examples how to use Go Context package.

   - [WithCancel](cmd/withcancel/main.go): Example how you can create a context and cancel it manually.
   - [WithDeadline](cmd/withdeadline/main.go): Example how you can give a deadline to a context.
   - [WithTimeout](cmd/withtimeout/main.go): Example how you can set a timeout for a context.
   - [WithValue](cmd/withvalue/main.go): Example how you can carry values in a context.

2. Complex example.

   In this example we have two microservices. The [shop](cmd/shop/main.go) service has an API endpoint where we can get the user's product list.
   The product list is provided by the [products](cmd/products/main.go) service via HTTP API endpoint. This endpoint handler simulates to gets
   the products from a DB, what is usually takes time. Both services handles syscall signals to do a gracefull shutdown.
   To able to use the shop's API you need to authenticate and the authentication middleware provides the username and it is used in the context to call the products' API.
   It also makes sure that the API call has a timeout via context.
