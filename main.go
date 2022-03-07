package main

import (
  "net/http"
  "os"
  "fmt"
  "net/url"
  "io/ioutil"
  "encoding/json"
  "strconv"

  "github.com/joho/godotenv"
)

var ytURL string = "https://youtube.googleapis.com/youtube/v3"

type IDObj struct {
  VideoID string `json:"videoId"`
}
type IDItem struct {
  Id IDObj `json:"id"`
}
type ErrorObj struct{
  Message string `json:"message"`
}
type IDResponse struct{
  Error ErrorObj `json:"error"`
  Items []IDItem `json:"items"`
}
type StatsObj struct{
  ViewCount string `json:"viewCount"`
  LikeCount string `json:"likeCount"`
  FavoriteCount string `json:"favoriteCount"`
}
type StatItem struct {
  Statistics StatsObj `json:"statistics"`
}
type StatsResponse struct{
  Items []StatItem `json:"items"`
}
func main () {
  if err := godotenv.Load(".env"); err != nil {
    print(err)
  }
  if len(os.Args) < 3 {
    fmt.Println("Missing arguments. Please pass in the names of two YouTube videos you want to do battle")
    return
  }

  video1Name := os.Args[1]
  video1ID := getIdByName(video1Name)
  fmt.Println("HEY", video1ID)
  video1Stats := getStats(video1ID)
  var video1Total int = getTotal(video1Stats)
  fmt.Println(video1Name)
  fmt.Println("Views:", formatNumber(video1Stats.ViewCount))
  fmt.Println("Likes:", formatNumber(video1Stats.LikeCount))
  fmt.Println("Total:", formatNumber(strconv.Itoa(video1Total)))
  fmt.Println("")

  video2Name := os.Args[2]
  video2ID := getIdByName(video2Name)
  video2Stats := getStats(video2ID)
  var video2Total int = getTotal(video2Stats)
  fmt.Println(video2Name)
  fmt.Println("Views:", formatNumber(video2Stats.ViewCount))
  fmt.Println("Likes:", formatNumber(video2Stats.LikeCount))
  fmt.Println("Total:", formatNumber(strconv.Itoa(video2Total)))

  if video1Total > video2Total {
    fmt.Println(os.Args[1], " Wins!")
  }
  if video1Total < video2Total {
    fmt.Println(os.Args[2], " Wins!")
  }
  if video1Total == video2Total {
    fmt.Println("Tie!")
  }
}

func formatNumber (num string) string {
  var formattedString string = ""
  for i:=len(num)-1; i>=0; i-- {
    if (i+1) % 3 == 0 {
      if (i !=0) {
        formattedString = formattedString + "," + string(num[i])
      }
    } else {
      formattedString = formattedString + string(num[i])
    }
  }
  return formattedString
}

func getTotal (stats StatsObj) int {
  views, _ := strconv.Atoi(stats.ViewCount)
  likes, _ := strconv.Atoi(stats.LikeCount)
  favs, _ := strconv.Atoi(stats.FavoriteCount)
  return views + likes + favs
}

func getIdByName(searchTerm string) string {
  urlStr := ytURL + "/search?part=id&maxResults=1&q=" + url.QueryEscape(searchTerm) + "&key=" + os.Getenv("API_KEY")
  response, getErr := http.Get(urlStr)
  if getErr != nil {
    print(getErr)
  }
  defer response.Body.Close()

  resultBytes, readErr := ioutil.ReadAll(response.Body)
  if readErr != nil {
    print(readErr)
  }

  var result IDResponse
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		print(err)
	}

  if result.Error.Message != ""  {
    panic(result.Error.Message)
  }

  return result.Items[0].Id.VideoID
}

func getStats(id string) StatsObj {
  urlStr := ytURL + "/videos?part=statistics&id=" + url.QueryEscape(id) + "&key=AIzaSyCivFO1PWBQahwRh9-BGm16iNz0CcvGqRg"
  response, getErr := http.Get(urlStr)
  if getErr != nil {
    print(getErr)
  }
  defer response.Body.Close()

  resultBytes, readErr := ioutil.ReadAll(response.Body)
  if readErr != nil {
    print(readErr)
  }

  var result StatsResponse
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		print(err)
	}

  return result.Items[0].Statistics
}
