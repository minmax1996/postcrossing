package main

import (
	"encoding/json"
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

	r.POST("/post", myHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	}
}

func myHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error in reading body"))
		return
	}

	request := make(map[string]interface{})
	if err := json.Unmarshal(bBody, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error in unmarshalling json body"))
		return
	}

	if request["from"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please tell me who you are. just type in \"from\" json field name / email / twitter / etc"))
		return
	}

	if err = savePostcard(request); err != nil {
		log.Println("[ERROR]: error in savePostcard" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("some error in saving your postcard"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello friend, I got your post! Have a great day"))
}

func savePostcard(request map[string]interface{}) error {
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
	notify("You got postcard from " + request["from"].(string))

	return nil
}

func notify(message string) {
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
