### PostCrossing

"Postcrossing is an online project that allows its members to send and receive postcards from all over the world"


But what if we exchange digital postcards, like json's via POST requests? This project will allow you to host endpoint to receive jsons and store them in your drive.

After receiving POST-request it checks "from" field, because you need to know who you're writing to.
Then it adds "created_date" field and saves that in `./postcards/` directory with a name from "from" field.

Thats it, all other fields are optional, but if you wanna receive a POSTcard from me, you can provide your return address ( callback url ) and write some description


### How to run
1) clone repo and run localy
2) run via docker 
```
docker run -d \
	-p 8080:8080 \
	-v $(shell pwd)/postcards:/postcards/ \
	minmax1996/postcrossing:latest
```
2.1) run via docker with discord notification 
```
docker run -d \
	-p 8080:8080 \
	-v $(shell pwd)/postcards:/postcards/ \
	-e "DISCORD_URL=<PASTE here you discord channel webhook>" \
	minmax1996/postcrossing:latest
```

### SEND-IT
```
curl -X POST \
  http://<hostname>:8080/post/ \
  -d '{
	"from": "MyOldFriend",
	"return_url": "https://example.com",
	"message": "Hello Friend",
	"stamp_url": "https://upload.wikimedia.org/wikipedia/commons/2/25/Rus_Stamp_RND.jpg",
	"picture_url": "https://www.canyousendmeapostcard.com/upload/09-russia-01-POSTCARD_CARTE-POSTALE_Wz4y.jpg"
}'
```
