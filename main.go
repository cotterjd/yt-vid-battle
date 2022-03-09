package main

import (
  "net/http"
  "os"
  "fmt"
  "net/url"
  "io/ioutil"
  "encoding/json"
  "strconv"
  "time"

  "github.com/joho/godotenv"
  "github.com/floscodes/golang-thousands"
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
type SnippetObj struct{
  PublishedAt string `json:"publishedAt"`
}
type StatItem struct {
  Statistics StatsObj `json:"statistics"`
  Snippet SnippetObj `json:"snippet"`
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

  vid1Name := os.Args[1]
  vid1Score := getAndShowStats(vid1Name)
  vid2Name := os.Args[2]
  vid2Score := getAndShowStats(vid2Name)

  if vid1Score > vid2Score {
    fmt.Println(vid1Name, " Wins!")
  }
  if vid1Score < vid2Score {
    fmt.Println(vid2Name, " Wins!")
  }
  if vid1Score == vid2Score {
    fmt.Println("Tie!")
  }

}

func getAndShowStats(vidName string) int {
  videoID := getIdByName(vidName)
  var videoInfo StatItem = getStats(videoID)
  var videoTotal int = getTotal(videoInfo)
  videoScore, yearsPublished := getScore(videoTotal, videoInfo)
  displayStats(vidName, videoTotal, videoScore, yearsPublished, videoInfo)

  return videoScore
}

func displayStats (name string, total int, score int, yearsPublished int64, info StatItem) {
  fmt.Println(name)
  views, _ := strconv.Atoi(info.Statistics.ViewCount)
  fmt.Println("Views:", formatNum(views))
  likes, _ := strconv.Atoi(info.Statistics.LikeCount)
  fmt.Println("Likes:", formatNum(likes))
  fmt.Println("Total:", formatNum(total))
  fmt.Println("Years Published:", yearsPublished)
  fmt.Println("Normalized Score:", formatNum(score))
  fmt.Println("")
}

func formatNum(num int) string {
  f, _ := thousands.Separate(num, "en")
  return f
}

func getTotal (item StatItem) int {
  views, _ := strconv.Atoi(item.Statistics.ViewCount)
  likes, _ := strconv.Atoi(item.Statistics.LikeCount)
  favs, _ := strconv.Atoi(item.Statistics.FavoriteCount)
  return views + likes + favs
}

func getScore (total int, item StatItem) (int, int64) {
  now := time.Now()
  nowStamp := now.Unix()
  // normalize for amount of time published
  divider := getYearsOld(nowStamp, item.Snippet.PublishedAt)
  if divider == 0 {
    return total, divider
  } else {
    return int(int64(total) / divider), divider
  }
}

func getYearsOld (nowStamp int64, dateStr string) int64 {
  date, _ := time.Parse(time.RFC3339, dateStr)
  stamp := date.Unix()
  return (nowStamp - stamp) / 365 / 24 / 60 / 60
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

func getStats(id string) StatItem {
  urlStr := ytURL + "/videos?part=statistics&part=snippet&id=" + url.QueryEscape(id) + "&key=AIzaSyCivFO1PWBQahwRh9-BGm16iNz0CcvGqRg"
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

  return result.Items[0]
}
