package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var discord_url string

func main() {
	discord_url = os.Getenv("DISCORD_URL")

	r := httprouter.New()

	r.POST("/post/", myHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	}
}

func myHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bBody, errR := ioutil.ReadAll(r.Body)
	if errR != nil {
		io.WriteString(w, fmt.Sprintf("error in reading body"))
		return
	}
	request := make(map[string]interface{})
	if err := json.Unmarshal(bBody, &request); err != nil {
		io.WriteString(w, fmt.Sprintf("error in unmarshalling json body"))
		return
	}

	if request["from"] == nil {
		io.WriteString(w, fmt.Sprintf("Please tell me who you are. just type in \"from\" json field name / email / twitter / etc"))
		return
	}

	if err := SavePostcard(request); err != nil {
		log.Println("[ERROR]: error in savePostcard" + err.Error())
		io.WriteString(w, fmt.Sprintf("some error in saving your postcard"))
		return
	}
	io.WriteString(w, "Hello friend, I got your post! Have a great day")
}

func SavePostcard(request map[string]interface{}) error {
	f, err := os.OpenFile("./postcards/"+request["from"].(string), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	request["created_date"] = time.Now().UTC().Format(time.RFC3339)

	fileJson, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := f.WriteString(string(fileJson) + "\n"); err != nil {
		log.Println(err)
		return err
	}
	Notify("You got postcard from " + request["from"].(string))

	return nil
}

func Notify(message string) {
	if len(discord_url) == 0 {
		return
	}
	payload := strings.NewReader("{ \"username\": \"" + "uploader" + "\", \"content\": \"" + message + "\" }")
	req, err := http.NewRequest("POST", discord_url, payload)
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	res.Body.Close()
}
