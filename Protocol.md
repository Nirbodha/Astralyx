# Protocol
This is the protocol that the client and the server use. The hexadecimal numbers are the IDs to the packets.
## Serverbound
- `0x0`: This is for pinging the server to check if it's up.
- `0x1`: This is for logging in.
- `0x2`: This is for chatting.
- `0x3`: This is for moving.
- `0x4`: This is for interactions. (This might be split into more IDs since this is so broad).
## Clientbound
- `0x0`: This is a response to a ping request.
- `0x1`: This is to verify if a login is a success.
- `0x2`: This is to tell every client in the server that a message has been sent.
- `0x3`: This is to tell every client in the server that someone is moving to a set coordinates.
- `0x4`: This is to tell every client in the server that someone is interacting with something.

This will be edited in the future. 
