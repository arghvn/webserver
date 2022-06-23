Creating a web server with Golang :
Go has HTTP package that contains utilities for quickly creating a web or file server.

The goal of this tutorial is to create a web server that can accept a GET request and serve a response. 
We’ll use the server to serve static files, acting as a file server.
 We’ll then make the web server respond to a POST request coming from a form submission, such as a contact form.

 The file server.go sits at the root of your project
  as does the static folder, which contains two HTML files: index.html and form.html.

  If we were successful in the startup phase and saw the relevant message, we will move on to the next step.

  Starting a web server with GET routes
At this stage, we’ll create a web server that is actually served on port 8080 and can respond to incoming GET requests.
Let’s modify the code in our main() function to start a web server at port 8080. The ListenAndServe method is exported by the http packet we imported during step one. This method allows us to start the web server and specify the port to listen for incoming requests.
Note that the port parameter needs to be passed as a string prepended by colon punctuation. The second parameter accepts a handler to configure the server for HTTP/2. However, this isn’t important for this tutorial, so we can safely pass nil as the second argument.

full the main function.

after this point, the server can start, but it still doesn’t know how to handle requests. We need to pass handlers to the server so it knows how to respond to incoming requests and which requests to accept.

We’ll use the HandleFunc function to add route handlers to the web server. The first argument accepts the path it needs to listen for /hello. Here, you tell the server to listen for any incoming requests for http://localhost:8080/hello. The second argument accepts a function that holds the business logic to correctly respond to the request.

By default, this function accepts a ResponseWriter to send a response back and a Request object that provides more information about the request itself. For example, you can access information about the sent headers, which can be useful for authenticating the request.


s you can see, the handler sends a “Hello!” message as we pass this response to the ResponseWriter.

Now let’s try out this setup. Start the web server with go run server.go and visit http://localhost:8080/hello. If the server responds with "Hello!", you can continue to the next step, where you’ll learn how to add basic security to your Golang web server routes.

it was ok
