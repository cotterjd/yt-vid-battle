# yt-vid-battle
Uses google's youtube api to compares youtube videos based on views/likes

Pass in two strings of names of videos
```
$ ./main.linux "chocolate rain" "david after dentist"
```
or 
```
$ vim .env # will need API_KEY
$ go mod tidy
$ go run main.go "chocolate rain" "david after dentist"
```