# Protocol
This is the protocol that the client and the server use. The hexadecimal numbers are the IDs to the packets.
\
\
This is how a data packet should be formatted. Note that this is also how an uncompressed packet looks like in Minecraft servers. Is this a straight rip from [wiki.vg?](https://wiki.vg) Yes.

Field Name | Field Type | Notes
---------- | ---------- | ------
Length     | VarInteger | This just tells the client the size of the incoming packet ID and the data. All of the data doesn't just come in at once.
Packet ID  | VarInteger | Nothing to say here.
Data       | Byte Array | This data is dependent on what the packet ID is about. It will be validated.

## Serverbound
###### Note: The "2" in the total length represents the number of bytes used for the packet IDs and an EOF character.

- `0x0`: This is for pinging the server to check if it's up.

Total Length | Data
------------ | ----
2            | Only an EOF character is to be expected

- `0x1`: This is for logging in.

Total Length | Data
------------ | ----
2 + ?        | Not known yet. When more progress with Google Logins happens, this will be updated.

- `0x2`: This is for chatting.

Total Length     | Data
---------------- | ----
2 + 254 (256) + ?| Byte array of a string. Authentication token size yet to be known.

- `0x3`: This is for moving.

Total Length | Data
------------ | ----
2 + ?        | Map size not yet known.

- `0x4`: This is for interactions. (This might be split into more IDs since this is so broad).

Total Length      | Data
----------------- | ----
2 + 1 (3)  + ? + ?| 1 byte identifier for the interaction. This may be upscaled in the future. Map size and Google Login Tokens not yet known.

## Clientbound
- `0x0`: This is a response to a ping request.

Total Length | Data
------------ | ----
2 + 1 (3)    | A one byte ID to see the status of the server. `0x0` means the server is joinable. `0x1` means the server is up but isn't accepting any connections. `0x2` means the server is up but it is full.

- `0x1`: This is to verify if a login is a success.
- `0x2`: This is to tell every client in the server that a message has been sent.
- `0x3`: This is to tell every client in the server that someone is moving to a set coordinates.
- `0x4`: This is to tell every client in the server that someone is interacting with something.

This will be edited in the future. I'm too tired to finish all of this at once. I'll be working on the server code to receive and read packets in the meanwhile.
