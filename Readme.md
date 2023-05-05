# Server Sent Events (Golang)

This is an example of a server sent event that is handled using go, and I also  created a client as "streamer" and publisher

## Run

Using golang version 1.19

```bash
go run main.go
```

## Usage

1. Open your browser on port 8080
2. On main routes, it should be the sreamer client
3. Access example `http://localhost:8080/post` to access the publisher
4. If someone post a tweet, then it will appear on the "streamer" client in realtime