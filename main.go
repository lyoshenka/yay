package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var slackApi *slack.Client

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slackToken := os.Getenv("SLACK")
	if slackToken != "" {
		slackApi = slack.New(slackToken)
	}

	http.HandleFunc("/", handler)
	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

const (
	formInputName = "f"
	slackChannel  = "@grin"
	slackUsername = "yay"
)

func handler(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimLeft(r.URL.String(), "/")

	log.Printf("%s %s\n", r.Method, url)

	switch url {
	case "":
		w.Write([]byte("Yay!"))
		return
	case "favicon.ico": // ignore these
		return
	case "robots.txt": // allow only homepage. don't want bots scraping the rest.
		w.Write([]byte(`User-Agent: *
Allow: /$
Disallow: /`))
		return
	case "favicon.svg":
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg"><text y="32" font-size="32">🙌</text></svg>`)) // emoji favicon
		return
	case "new.min.css":
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(newCSS))
		return
	case "thank-you":
		w.Write([]byte(strings.Replace(layout, "BODY_GOES_HERE", `<h1>Thanks again!</h1>`, 1)))
		return
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "<hostname-unknown>"
	}

	if r.Method == http.MethodPost {
		sendToSlack("%s | *%s* %s\n%s", hostname, url, r.Header.Get("User-Agent"), r.FormValue(formInputName))
		w.Header().Set("Location", "/thank-you")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	sendToSlack("%s | *%s* %s", hostname, url, r.Header.Get("User-Agent"))
	w.Write([]byte(strings.Replace(layout, "BODY_GOES_HERE", `
    <h1>Thanks for your feedback</h1>
	<br>
	<p>Is there anything you'd like to add?</p>
	<form method="POST" action="">
		<p><textarea name="`+formInputName+`" rows=8 style="width: 100%"></textarea></p>
		<p><input type="submit"></p>
	</form>
`, 1)))
}

func sendToSlack(format string, a ...interface{}) {
	if slackApi == nil {
		log.Println("SLACK TOKEN NOT SET")
		return
	}

	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}

	_, _, err := slackApi.PostMessage(slackChannel, slack.MsgOptionText(message, false), slack.MsgOptionUsername(slackUsername))
	if err != nil {
		log.Println("error sending to slack: " + err.Error())
	}
}

var layout = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Yay</title>
	<link rel="icon" href="/favicon.svg" />
	<link rel="stylesheet" href="/new.min.css">
</head>
<body>
	BODY_GOES_HERE
</body>
</html>
`

var newCSS = `/**
 * Minified by jsDelivr using clean-css v4.2.1.
 * Original file: /npm/@exampledev/new.css@1.1.2/new.css
 *
 * Do NOT use SRI with dynamically generated files! More information: https://www.jsdelivr.com/using-sri-with-dynamic-files
 */
:root{--nc-font-sans:'Inter',-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxygen,Ubuntu,Cantarell,'Open Sans','Helvetica Neue',sans-serif,"Apple Color Emoji","Segoe UI Emoji","Segoe UI Symbol";--nc-font-mono:Consolas,monaco,'Ubuntu Mono','Liberation Mono','Courier New',Courier,monospace;--nc-tx-1:#000000;--nc-tx-2:#1A1A1A;--nc-bg-1:#FFFFFF;--nc-bg-2:#F6F8FA;--nc-bg-3:#E5E7EB;--nc-lk-1:#0070F3;--nc-lk-2:#0366D6;--nc-lk-tx:#FFFFFF;--nc-ac-1:#79FFE1;--nc-ac-tx:#0C4047}@media (prefers-color-scheme:dark){:root{--nc-tx-1:#ffffff;--nc-tx-2:#eeeeee;--nc-bg-1:#000000;--nc-bg-2:#111111;--nc-bg-3:#222222;--nc-lk-1:#3291FF;--nc-lk-2:#0070F3;--nc-lk-tx:#FFFFFF;--nc-ac-1:#7928CA;--nc-ac-tx:#FFFFFF}}*{margin:0;padding:0}address,area,article,aside,audio,blockquote,datalist,details,dl,fieldset,figure,form,iframe,img,input,meter,nav,ol,optgroup,option,output,p,pre,progress,ruby,section,table,textarea,ul,video{margin-bottom:1rem}button,html,input,select{font-family:var(--nc-font-sans)}body{margin:0 auto;max-width:750px;padding:2rem;border-radius:6px;overflow-x:hidden;word-break:break-word;overflow-wrap:break-word;background:var(--nc-bg-1);color:var(--nc-tx-2);font-size:1.03rem;line-height:1.5}::selection{background:var(--nc-ac-1);color:var(--nc-ac-tx)}h1,h2,h3,h4,h5,h6{line-height:1;color:var(--nc-tx-1);padding-top:.875rem}h1,h2,h3{color:var(--nc-tx-1);padding-bottom:2px;margin-bottom:8px;border-bottom:1px solid var(--nc-bg-2)}h4,h5,h6{margin-bottom:.3rem}h1{font-size:2.25rem}h2{font-size:1.85rem}h3{font-size:1.55rem}h4{font-size:1.25rem}h5{font-size:1rem}h6{font-size:.875rem}a{color:var(--nc-lk-1)}a:hover{color:var(--nc-lk-2)}abbr:hover{cursor:help}blockquote{padding:1.5rem;background:var(--nc-bg-2);border-left:5px solid var(--nc-bg-3)}abbr{cursor:help}blockquote :last-child{padding-bottom:0;margin-bottom:0}header{background:var(--nc-bg-2);border-bottom:1px solid var(--nc-bg-3);padding:2rem 1.5rem;margin:-2rem calc(0px - (50vw - 50%)) 2rem;padding-left:calc(50vw - 50%);padding-right:calc(50vw - 50%)}header h1,header h2,header h3{padding-bottom:0;border-bottom:0}header>:first-child{margin-top:0;padding-top:0}header>:last-child{margin-bottom:0}a button,button,input[type=button],input[type=reset],input[type=submit]{font-size:1rem;display:inline-block;padding:6px 12px;text-align:center;text-decoration:none;white-space:nowrap;background:var(--nc-lk-1);color:var(--nc-lk-tx);border:0;border-radius:4px;box-sizing:border-box;cursor:pointer;color:var(--nc-lk-tx)}a button[disabled],button[disabled],input[type=button][disabled],input[type=reset][disabled],input[type=submit][disabled]{cursor:default;opacity:.5;cursor:not-allowed}.button:focus,.button:hover,button:focus,button:hover,input[type=button]:focus,input[type=button]:hover,input[type=reset]:focus,input[type=reset]:hover,input[type=submit]:focus,input[type=submit]:hover{background:var(--nc-lk-2)}code,kbd,pre,samp{font-family:var(--nc-font-mono)}code,kbd,pre,samp{background:var(--nc-bg-2);border:1px solid var(--nc-bg-3);border-radius:4px;padding:3px 6px;font-size:.9rem}kbd{border-bottom:3px solid var(--nc-bg-3)}pre{padding:1rem 1.4rem;max-width:100%;overflow:auto}pre code{background:inherit;font-size:inherit;color:inherit;border:0;padding:0;margin:0}code pre{display:inline;background:inherit;font-size:inherit;color:inherit;border:0;padding:0;margin:0}details{padding:.6rem 1rem;background:var(--nc-bg-2);border:1px solid var(--nc-bg-3);border-radius:4px}summary{cursor:pointer;font-weight:700}details[open]{padding-bottom:.75rem}details[open] summary{margin-bottom:6px}details[open]>:last-child{margin-bottom:0}dt{font-weight:700}dd::before{content:'→ '}hr{border:0;border-bottom:1px solid var(--nc-bg-3);margin:1rem auto}fieldset{margin-top:1rem;padding:2rem;border:1px solid var(--nc-bg-3);border-radius:4px}legend{padding:auto .5rem}table{border-collapse:collapse;width:100%}td,th{border:1px solid var(--nc-bg-3);text-align:left;padding:.5rem}th{background:var(--nc-bg-2)}tr:nth-child(even){background:var(--nc-bg-2)}table caption{font-weight:700;margin-bottom:.5rem}textarea{max-width:100%}ol,ul{padding-left:2rem}li{margin-top:.4rem}ol ol,ol ul,ul ol,ul ul{margin-bottom:0}mark{padding:3px 6px;background:var(--nc-ac-1);color:var(--nc-ac-tx)}input,select,textarea{padding:6px 12px;margin-bottom:.5rem;background:var(--nc-bg-2);color:var(--nc-tx-2);border:1px solid var(--nc-bg-3);border-radius:4px;box-shadow:none;box-sizing:border-box}img{max-width:100%}
/*# sourceMappingURL=/sm/4a51164882967d28a74fabce02685c18fa45a529b77514edc75d708f04dd08b9.map */`
