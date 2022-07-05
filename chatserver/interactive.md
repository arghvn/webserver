Our design has very simple features:

There is a single chat room.
The user can connect to the server.
The user can set his own name.
A user can send a message to the room and the message will be broadcast to all users of the room.
At the current level there is no need to save the messages and the user will only see the messages if he is connected to the room.

Communication between server and client will be done through TCP using a simple string protocol.

The following are implemented in this protocol:

Send command: The client sends a chat message.
Name command: The client sets its name.
Message Order: The server broadcasts the chat message to all members.
A command is a string that starts with the name of the command, has all the parameters, and ends with n\.

For example, to send a "hello" message, the client sends the SEND Hello\n command on the TCP socket, and then the server sends the MESSAGE username Hello\n command to other clients.
