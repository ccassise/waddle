# waddle
Chat server written in Go.

Idea taken from https://web.archive.org/web/20111012115624/http://itasoftware.com/careers/work-at-ita/hiring-puzzles.html

https://user-images.githubusercontent.com/58533624/135701220-734877e1-5b7a-4c2c-9549-fe09ff4a02c3.mov

## Quick Start
```
waddle [port]
```
#### Prerequisites
- Go compiler.
- (Optional) Some program that sends data over a TCP socket. Any examples here will use `netcat` since it comes preinstalled on MacOS, most Linux distributions, and is available on Windows.

#### Build & Run
1. Clone or download the repo.
2. `cd` to its directory.
3. Compile and run with Go: `go run cmd/waddle/waddle.go [port]`. NOTE: Port is required.

This will start the server and listen on the given port.
To connect to the server:
```
nc localhost [port]
```
Where port is the same port you used on the compile and run command. You should now be able to send and receive messages to/from the server.

## Protocol
```
<CRLF> indicates the bytes "\r\n".

LOGIN <username><CRLF>                                    - Login as given username.
JOIN #<chatroom><CRLF>                                    - Create or join a chatroom. Chatrooms begin with '#'.
PART #<chatroom><CRLF>                                    - Leave a chatroom. A user is able to join multiple chatrooms at once.
MSG #<chatroom> <message-text><CRLF>                      - Send a message to all users in a chatroom.
MSG <username> <message-text><CRLF>                       - Send a message directly to user.
LOGOUT<CRLF>                                              - Log off and close connection to server.
  
Server responses:
OK<CRLF>                                                  - Indicates command was accepted.
ERROR <reason><CRLF>                                      - Indicates an error has occured.
GOTROOMMSG <sender> #<chatroom> <message-text><CRLF>      - When a message was sent to the room the user is in.
GOTUSERMSG <sender> <message-text><CRLF>                  - When a message was sent directy to the user.
```
