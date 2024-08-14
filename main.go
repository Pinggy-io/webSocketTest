package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WebSocket Time</title>
</head>
<body>
    <div id="time"></div>

    <script>
        const isSecure = window.location.protocol === "https:";
        const protocol = isSecure ? "wss" : "ws";
        const socket = new WebSocket(protocol + "://" + window.location.host + "/ws");

        socket.onmessage = function(event) {
            document.getElementById("time").innerHTML = "Current Time: " + event.data;
            socket.send("asda");
        };
    </script>
</body>
</html>
`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		err := conn.WriteMessage(websocket.TextMessage, []byte(currentTime))
		if err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(1 * time.Second)
		typ, _, err := conn.ReadMessage()
		fmt.Println("Hello", typ, err)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "p", 8080, "Port to run the server on")
	flag.Parse()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
