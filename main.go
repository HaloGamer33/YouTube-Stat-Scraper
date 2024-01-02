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
    collector.OnRequest(
        func (request *colly.Request) {
            fmt.Println("Visiting", request.URL.String())
        },
    )
    collector.OnHTML("meta[name=title]",
        func (element *colly.HTMLElement) {
            vidStats.Title = element.Attr("content")
        },
    )
    collector.OnHTML("meta[itemprop=uploadDate]",
        func (element *colly.HTMLElement) {
            vidStats.UploadDate = element.Attr("content")[0:10]
            vidStats.UploadHour = element.Attr("content")[11:]
        },
    )

    // Searching and processing <script> elements that contain important json
    var vidJson VideoJson
    var scriptCounter int = 0
    // var scripts []colly.HTMLElement
    collector.OnHTML("script",
        func (element *colly.HTMLElement) {
            // scripts = append(scripts, *element)
            if scriptCounter == 20 {
                // Searching for the index of "{" (marks the start of json)
                var jsonStart int = strings.Index(element.Text, "{")
                // Selecting only the json, (there is a ";" at the end that messes up the json)
                var jsonString string = element.Text[jsonStart:len(element.Text)-1]

                // Json decoding
                err := json.Unmarshal([]byte(jsonString), &vidJson)
                if err != nil { panic(err) }
            }
            if scriptCounter == 45 {
                // Searching for the index of "{" (marks the start of json)
                var jsonStart int = strings.Index(element.Text, "{")
                // Selecting only the json, (there is a ";" at the end that messes up the json)
                var jsonString string = element.Text[jsonStart:len(element.Text)-1]

                // Json decoding
                var myData MyJSON
                err := json.Unmarshal([]byte(jsonString), &myData)
                if err != nil { panic(err) }

                string := myData.Contents.TwoColumnWatchNextResults.Results.Results.Contents[0].VideoPrimaryInfoRenderer.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.LikeButtonViewModel.LikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.DefaultButtonViewModel.ButtonViewModel.AccessibilityText
                strLikes := ExtractLikes(string)
                likes, err := strconv.ParseInt(strLikes, 10, 0)
                if err != nil { panic(err) }
                vidStats.Likes = int(likes)
            }
            scriptCounter++
        },
    )
    collector.OnResponse(func(response *colly.Response) {
        // err := os.WriteFile("output.html", response.Body, 0644)
        // if err != nil { panic(err) }
    })

    var link string
    fmt.Println("Insert the video link: ")
    fmt.Scanln(&link)
    collector.Visit(link)

    vidStats.TransferJson(vidJson)

    fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Title)
    // for index, script := range scripts {
    //     title := strconv.Itoa(index)
    //     os.WriteFile(title + ".txt", []byte(script.Text), 0644)
    // }
    os.WriteFile(fileTitle, []byte(vidStats.Format()), 0644)
}

func (v VideoStats) Format() string {
    var s string

    length64, err := strconv.ParseInt(v.LengthSeconds, 10, 0)
    if err != nil { panic(err) }

    lengthMin := SecondsToMinutes(int(length64))

    s += fmt.Sprintf("%-20v %v\n", "Title:", v.Title)
    s += fmt.Sprintf("%-20v %v\n", "Likes:", v.Likes)
    s += fmt.Sprintf("%-20v %v\n", "Upload Date:", v.UploadDate)
    s += fmt.Sprintf("%-20v %v\n", "Upload Hour:", v.UploadHour)
    s += fmt.Sprintf("%-20v %v\n", "Length:", lengthMin)
    s += fmt.Sprintf("%-20v %v\n", "Length (seconds):", v.LengthSeconds)
    s += fmt.Sprintf("Description:\n\n%v", v.Description)
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

func SecondsToMinutes(inSeconds int) string {
    minutes := inSeconds / 60
    seconds := inSeconds % 60
    str := fmt.Sprintf("%v:%v", minutes, seconds)
    return str
}

func FindFirstInt(str string) int {
    ints := []rune{ '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', }
    for index, char := range str {
        for _, int := range ints {
            if char == int {
                return index
            }
        }
    }
    return -1
}

func ExtractLikes(str string) string {
    var likes string
    index := FindFirstInt(str)
    likes = str[index:]
    final := strings.Index(likes, " ")
    likes = likes[:final]
    likes = strings.ReplaceAll(likes, ",", "")

    return likes
}

func (vidStats *VideoStats) TransferJson(json VideoJson) {
    vidStats.VideoID = json.VideoDetails.VideoID
    vidStats.LengthSeconds = json.VideoDetails.LengthSeconds
    vidStats.Keywords = json.VideoDetails.Keywords
    vidStats.ChannelID = json.VideoDetails.ChannelID
    vidStats.Description = json.VideoDetails.Description
    vidStats.IsCrawlable = json.VideoDetails.IsCrawlable
    vidStats.AllowRatings = json.VideoDetails.AllowRatings
    vidStats.ViewCount = json.VideoDetails.ViewCount
    vidStats.Author = json.VideoDetails.Author
    vidStats.IsPrivate = json.VideoDetails.IsPrivate
    vidStats.IsLiveContent = json.VideoDetails.IsLiveContent
}

type VideoStats struct {
    Title string
    Likes int
    UploadDate string
    UploadHour string
    VideoID  string  
    LengthSeconds string 
    Keywords []string 
    ChannelID string 
    Description string 
    IsCrawlable bool 
    AllowRatings bool 
    ViewCount string 
    Author string 
    IsPrivate bool 
    IsLiveContent bool 
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

// not sure if there is a better way to access the json, but this works
type MyJSON struct { Contents struct { TwoColumnWatchNextResults struct { Results struct { Results struct { Contents []struct { VideoPrimaryInfoRenderer struct { VideoActions struct { MenuRenderer struct { TopLevelButtons []struct { SegmentedLikeDislikeButtonViewModel struct { LikeButtonViewModel struct { LikeButtonViewModel struct { ToggleButtonViewModel struct { ToggleButtonViewModel struct { DefaultButtonViewModel struct { ButtonViewModel struct { AccessibilityText string `json:"accessibilityText"` } `json:"buttonViewModel"` } `json:"defaultButtonViewModel"` } `json:"toggleButtonViewModel"` } `json:"toggleButtonViewModel"` } `json:"likeButtonViewModel"` } `json:"likeButtonViewModel"` } `json:"segmentedLikeDislikeButtonViewModel"` } `json:"topLevelButtons"` } `json:"menuRenderer"` } `json:"videoActions"` } `json:"videoPrimaryInfoRenderer"` } `json:"contents"` } `json:"results"` } `json:"results"` } `json:"twoColumnWatchNextResults"` } `json:"contents"` }
