# Server
This is just the code used for the server. Do not use this yet. It's horrid right now.
\
\
Let's specify the protocol. Is this finished? Definitely not!
## Server Protocol
NOTE: The hexadecimal numbers are the IDs of the packets that the server will be receiving and sending.

### Receive
- `0x0`: This is for pinging. This checks if the server is online. Expect to receive 1 byte.
- `0x1`: This is for logging in. Authentication questionable. TODO.
- `0x2`: This is for chatting. If a client wants to send a message, it will send a packet with this ID. Expect to receive 16 bytes. Because username restrictions are not yet known, TODO.
- `0x3`: This is for moving. The client has to send the desired coordinates to the server to move. Positional limits are not known. TODO.
A lot is left to write for the protocol.

### Send
I'll work on this tomorrow. I'm tired. :(
