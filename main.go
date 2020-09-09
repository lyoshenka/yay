package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lbryio/lbry.go/v2/extras/util"
)

// TODO: add freeform feedback option after initial feedback is in

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slackToken := os.Getenv("SLACK")
	if slackToken != "" {
		util.InitSlack(slackToken, "@grin", "yay")
	}

	http.HandleFunc("/", handler)
	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimLeft(r.URL.String(), "/")
	if url == "" {
		w.Write([]byte("Yay or Nay"))
		return
	}

	log.Println(url)

	ignored := []string{"favicon.ico", "robots.txt"}
	for _, i := range ignored {
		if url == i {
			return
		}
	}

	ua := r.Header.Get("User-Agent")

	err := util.SendToSlack(url + " | " + ua)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Thanks</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@exampledev/new.css@1.1.2/new.min.css">
	<style>
	</style>
</head>
<body>
    <h1>Thanks for your feedback</h1>
</body>
</html>
`))
}
