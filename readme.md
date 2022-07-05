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

Add basic security to routes
It goes without saying that security is important. Let’s explore some basic strategies to enhance the security of your Go web server.

Before we do, we should take a moment to increase the readability of our code. Let’s create the helloHandler function, which holds all the logic related to the /hello request.


This handler uses the Request object to check whether the requested path is correct. This is a very basic example of how you can use the Request object.

If the path is incorrect, the server returns a StatusNotFound error to the user. To write an error to the user, you can use the http.Error method. Notice that the StatusNotFound code corresponds to a 404 error. All status codes can be found in the Golang documentation.

Next, we add a check for verifying the type of the request. If the method doesn’t correspond to GET, the server returns a new error. When both checks pass, the server returns its success response "Hello!".

The last thing we need to do is modify the handleFunc function in our main() function to accept the above helloHandler function.


Next, we’ll start the Go web server with go run server.go. You can test your security by sending a POST request to http://localhost:8080/hello using a tool such as Postman or cURL.

Start a static web server
In this step, we’ll create a simple file server to host static files. This will be a very simple addition to the web server.

To make sure we have content to serve on the web server, let’s modify the index.html file located in the static folder. To keep things simple, just add a heading to the file that says “Static Website.” If you wish, you can add more files or styling files to make your web server look a bit nicer.

To serve the static folder, you’ll have to add two lines of code to server.go. The first line of code creates the file server object using the FileServer function. This function accepts a path in the http.Dir type. Therefore, we have to convert the string path “./static” to an http.Dir path type.

Don’t forget to specify the Handle route, which accepts a path and the fileserver. This function acts in the same way as the HandleFunc function, with some small differences.

It’s time to try out the code. Fire up the server with go run server.go and visit http://localhost:8080/. You should see the “Static Website” header.

Accept a form submission POST request
Lastly, the web server has to respond to a form submission.

Let’s add some content to the form.html file in the static folder. Notice that the form action is sent to /form. This means the POST request from the form will be sent to http://localhost:8080/form. The form itself asks for input for two variables: name and address.

The next step is to create the handler to accept the /form request. The form.html file is already served via the FileServer and can be accessed via http://localhost:8080/form.html.

First, the function has to call ParseForm() to parse the raw query and update r.PostForm and r.Form. This will allow us to access the name and address values via the r.FormValue method.

At the end of the function, we write both values to the ResponseWriter using fmt.Fprintf.

Trying out the form handler
We can test the form by starting the server with go run server.go. When the server starts, visit http://localhost:8080/form.html. You should see two input fields and a submit button.

When you’ve filled out the form, hit the submit button. The server should process your POST request and show you the result on the http://localhost:8080/form response page.
