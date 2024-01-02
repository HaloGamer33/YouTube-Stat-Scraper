package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "strconv"
    "os"
    "encoding/json"
    "strings"
)

func main() { 
    vidStats := NewVideoStats()

    collector := colly.NewCollector()
    collector.OnHTML(
        "meta[name=title]",
        func (element *colly.HTMLElement) {
            vidStats.Title = element.Attr("content")
        },
    )

    collector.OnHTML(
        "meta[itemprop=interactionCount]",
        func (element *colly.HTMLElement) {
            likes, _ := strconv.ParseInt(element.Attr("content"), 10, 0)
            vidStats.Likes = int(likes)
        },
    )

    collector.OnHTML(
        "meta[itemprop=uploadDate]",
        func (element *colly.HTMLElement) {
            vidStats.UploadDate = element.Attr("content")[0:10]
            vidStats.UploadHour = element.Attr("content")[11:]
        },
    )

    // Searching and processing the 20th <script> element (contains important json)
    var vidJson VideoJson
    var scriptCounter int = 0
    collector.OnHTML(
        "script",
        func (element *colly.HTMLElement) {
            if scriptCounter == 20 {
                // Searching for the index of "{" (marks the start of json)
                var jsonStart int = strings.Index(element.Text, "{")
                // Selecting only the json, (there is a ";" at the end that messes up the json)
                var jsonString string = element.Text[jsonStart:len(element.Text)-1]

                // Json decoding
                err := json.Unmarshal([]byte(jsonString), &vidJson)
                if err != nil {
                    panic(err)
                }
            }
            scriptCounter++
        },
    )

    
    collector.OnRequest(
        func (request *colly.Request) {
            fmt.Println("Visiting ", request.URL.String())
        },
    )

    collector.Visit("https://www.youtube.com/watch?v=CV0Nr90lKdE")

    vidStats.Json = vidJson

    fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Json.VideoDetails.Title)
    os.WriteFile(fileTitle, []byte(vidStats.Format()), 0644)
}

func (v VideoStats) Format() string {
    var s string

    length64, err := strconv.ParseInt(v.Json.VideoDetails.LengthSeconds, 10, 0)
    if err != nil {
        panic(err)
    }

    lengthMin := SecondsToMinutes(int(length64))

    s += fmt.Sprintf("%-20v %v\n", "Title:", v.Title)
    s += fmt.Sprintf("%-20v %v\n", "Likes:", v.Likes)
    s += fmt.Sprintf("%-20v %v\n", "Upload Date:", v.UploadDate)
    s += fmt.Sprintf("%-20v %v\n", "Upload Hour:", v.UploadHour)
    s += fmt.Sprintf("%-20v %v\n", "Length:", lengthMin)
    s += fmt.Sprintf("%-20v %v\n", "Length (seconds):", v.Json.VideoDetails.LengthSeconds)
    s += fmt.Sprintf("Description:\n\n%v", v.Json.VideoDetails.Description)
    return s
}


func NewVideoStats() VideoStats {
    vidStats := VideoStats{
        Title: "",
        Likes: 0,
        UploadDate: "",
        UploadHour: "",
    }

    return vidStats
}


type VideoStats struct {
    Title string
    Likes int
    UploadDate string
    UploadHour string

    Json VideoJson
}

type VideoJson struct {
    VideoDetails struct {
        VideoID  string  `json:"videoId"`
        Title    string  `json:"title"`
        LengthSeconds string `json:"lengthSeconds"`
        Keywords []string `json:"keywords"`
        ChannelID string `json:"channelId"`
        Description string `json:"shortDescription"`
        IsCrawlable bool `json:"isCrawlable"`
        AllowRatings bool `json:"allowRatings"`
        ViewCount string `json:"viewCount"`
        Author string `json:"author"`
        IsPrivate bool `json:"isPrivate"`
        IsLiveContent bool `json:"isLiveContent"`
    } `json:"videoDetails"`
}

func SecondsToMinutes(inSeconds int) string {
    minutes := inSeconds / 60
    seconds := inSeconds % 60
    str := fmt.Sprintf("%v:%v", minutes, seconds)
    return str
}
