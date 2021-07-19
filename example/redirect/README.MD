# Redirect Example

Rony is capable of handling redirection server side or client side. This example demonstrates
this feature.

To bring up 3 instance of the server you can use `goreman`[https://github.com/mattn/goreman].

1. Build the executable: `go install ./cmd/cli-redirect`

2. Run services with goreman or bring three (or more) instances manually: `goreman start`
> Please mind that there is --seed flag in Procfile, you must change it according to your environment
> otherwise your cluster would not work properly.

3. Open browser: `http://localhost/info/server-redirect/3/sometext`
> This demonstrates the server side redirection. You are asking request from replica-set 3 but
> sending your request to node1 which is replica-set 1. But it internally redirects your request
> to node3 and sends its result back to you


4. Open browser: `http://localhost/info/client-redirect/3/sometext`
> This demonstrates the client side redirection. You will see that server redirects you to new
> url by sending HTTP 307 and the address of the node3
> 