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
full the main function