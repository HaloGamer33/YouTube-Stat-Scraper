package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "strconv"
    "os"
    "encoding/json"
    "strings"
    "runtime"
    "net/http"
    "io"
)

func main() { 
    vidStats := NewVideoStats()
    collector := colly.NewCollector()
    //  ┌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┐
    //  ╎ Callback Functions ╎
    //  └╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┘
//  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    collector.OnRequest(
        func (request *colly.Request) {
            fmt.Println("Visiting", request.URL.String())
        },
    )
    // collector.OnHTML("html",
    //     func (element *colly.HTMLElement) {
    //         string, err := element.DOM.Html()
    //         err = os.WriteFile("output.html", []byte(string), 0644)
    //         if err != nil { panic(err) }
    //     },
    // )
    collector.OnHTML("meta[name=title]",
        func (element *colly.HTMLElement) {
            vidStats.Title = element.Attr("content")
        },
    )
    collector.OnHTML("meta[itemprop=uploadDate]",
        func (element *colly.HTMLElement) {
            // Getting date & hour from json with ISO 8601 format (YYYY-MM-DDTHH:MM:SS±HH:MM)
            // The 'T' is the divider betwen date and hour.
            vidStats.UploadDate = element.Attr("content")[0:10]
            vidStats.UploadHour = element.Attr("content")[11:]
        },
    )
    // Searching and processing <script> elements that contain important json
    var vidJson VideoJson
    var scriptCounter int
    collector.OnHTML("script",
        func (element *colly.HTMLElement) {

            // Writing all scripts for debug
            title := fmt.Sprintf("scripts/%v.txt", scriptCounter)
            err := os.WriteFile(title, []byte(element.Text), 0644)
            if err != nil { panic(err) }

            if scriptCounter == 20 {
                // Removing js from the json
                var indexJsonStart int = strings.Index(element.Text, "{")
                var jsonStr string = element.Text[indexJsonStart:len(element.Text)-1]

                err := json.Unmarshal([]byte(jsonStr), &vidJson)
                if err != nil { panic(err) }
            }
            if scriptCounter == 45 {
                // Removing js from the json
                var indexJsonStart int = strings.Index(element.Text, "{")
                var jsonStr string = element.Text[indexJsonStart:len(element.Text)-1]

                var likesJson LikesJson
                err := json.Unmarshal([]byte(jsonStr), &likesJson)
                if err != nil { panic(err) }

                likesStr := GetLikesStr(likesJson)
                likes := ExtractLikes(likesStr)
                vidStats.Likes = int(likes)
            }
            if scriptCounter == 45 {
                // Removing js from the json
                var indexJsonStart int = strings.Index(element.Text, "{")
                var jsonStr string = element.Text[indexJsonStart:len(element.Text)-1]

                var continueTokenJson ContinueTokenJson
                err := json.Unmarshal([]byte(jsonStr), &continueTokenJson)
                if err != nil { panic(err) }

                fmt.Println(continueTokenJson)

                token := GetContinuationToken(continueTokenJson)
                continuationJson := PageLoadContinuation(token)

                var commentCounterJson CommentCounterJson
                err = json.Unmarshal([]byte(continuationJson), &commentCounterJson)
                noCommentsStr := GetCommentCounter(commentCounterJson)
                noCommentsStr = strings.ReplaceAll(noCommentsStr, ",", "")
                comments, err := strconv.ParseInt(noCommentsStr, 10, 0)
                if err != nil { panic(err) }

                vidStats.Comments = int(comments)
            }
            scriptCounter++
        },
    )
    //  ┌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┐
    //  ╎ User Input Start ╎
    //  └╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┘
    var selection string
    PrintMenu()
    fmt.Scanln(&selection)
    switch selection {
    case "1":
        fmt.Println("Insert the link:")
        fmt.Scanln(&selection)
        ScrapeSingleVideo(selection, &vidStats, &vidJson, collector)
    case "2":
        fmt.Println("Name of the .txt file containing the links:")
        fmt.Scanln(&selection)
        fileContents, err := os.ReadFile(selection)
        stringContents := string(fileContents)
        links := strings.Split(stringContents, "\n")
        if err != nil { panic(err) }
        ScrapeVideos(links, &vidStats, &vidJson, collector, &scriptCounter)
    case "3":
        fmt.Println("Channel link:")
        fmt.Scanln(&selection)

        startOfName := strings.Index(selection, "@")
        channelUsername := selection[startOfName:]
        endOfName := strings.Index(channelUsername, "/")
        channelUsername = channelUsername[:endOfName]
        channelVideosLink := fmt.Sprintf("https://www.youtube.com/%v/videos", channelUsername)
        folderName := fmt.Sprintf("%v - Channel Video Stats", channelUsername)


        links, token := ScrapeChannel(channelVideosLink)

        if token == "" {
            ScrapeVideosIntoFolder(links, folderName, &vidStats, &vidJson, collector, &scriptCounter)
        } else {
            for token != "" {
                continuationVideosJsonStr := PageLoadChannelVideoContinuation(token)
                var continuationVideosJson ContinuationVideosJson
                err := json.Unmarshal([]byte(continuationVideosJsonStr), &continuationVideosJson)
                if err != nil { panic(err) }

                items := continuationVideosJson.OnResponseReceivedActions[0].AppendContinuationItemsAction.ContinuationItems

                for index, item := range items {
                    if index == len(items) - 1 {
                        token = item.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
                        if token != "" { continue }
                    }

                    id := &item.RichItemRenderer.Content.VideoRenderer.VideoId
                    link := fmt.Sprintf("https://www.youtube.com/watch?v=%v", *id)
                    links = append(links, link)
                }
            }

            ScrapeVideosIntoFolder(links, folderName, &vidStats, &vidJson, collector, &scriptCounter)
        }
    }
}

/*
     ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
            ┌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┐
            ╎ Function & Struct Declaration Beggining ╎
            └╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┘
     ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
*/
func ScrapeChannel(channelLink string) ([]string, string) {
    channelScriptCounter := 0
    links := []string{}
    continuationToken := ""

    collector := colly.NewCollector()
    collector.OnHTML("script",
        func (element *colly.HTMLElement) {
            path := fmt.Sprintf("channelScripts/%v", channelScriptCounter)
            os.WriteFile(path, []byte(element.Text), 0644)

            if channelScriptCounter == 36 {
                // Removing js from the json
                var indexJsonStart int = strings.Index(element.Text, "{")
                var jsonStr string = element.Text[indexJsonStart:len(element.Text)-1]

                var channelVideosJson ChannelVideosJson
                err := json.Unmarshal([]byte(jsonStr), &channelVideosJson)
                if err != nil { panic(err) }

                loadedContent := channelVideosJson.Contents.TwoColumnBrowseResultsRenderer.Tabs[1].TabRenderer.Content.RichGridRenderer.Contents
                for index, content := range loadedContent {
                    if index == len(loadedContent) - 1 {
                        continuationToken = content.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
                        continue
                    }
                    links = append(links, content.RichItemRenderer.Content.VideoRenderer.VideoId)
                }

            }
            channelScriptCounter++
        },
    )
    collector.Visit(channelLink)

    var fullLinks []string
    for _, link := range links {
        link = fmt.Sprintf("https://www.youtube.com/watch?v=%v", link)
        fullLinks = append(fullLinks, link)
    }

    return fullLinks, continuationToken
}

func PageLoadChannelVideoContinuation(token string) string {
    link := "https://www.youtube.com/youtubei/v1/browse?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false"
    method := "POST"

    json := `
    {
        "context": {
            "client": {
                "clientName": "WEB",
                "clientVersion": "2.20240101.07.00"
            }
        },
        "continuation": "%v"
    }`

    json = fmt.Sprintf(json, token)
    payload := strings.NewReader(json)

    client := &http.Client {}
    request, err := http.NewRequest(method, link, payload)
    if err != nil { panic(err) }
    request.Header.Add("Content-Type", "application/json")

    response, err := client.Do(request)
    if err != nil { panic(err) }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil { panic(err) }
    return string(body)
}

func PageLoadContinuation(token string) string {
    link := "https://www.youtube.com/youtubei/v1/next?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false"
    method := "POST"

    json := `
    {
        "context": {
            "client": {
                "clientName": "WEB",
                "clientVersion": "2.20240101.07.00"
            }
        },
        "continuation": "%v"
    }`

    json = fmt.Sprintf(json, token)
    payload := strings.NewReader(json)

    client := &http.Client {}
    request, err := http.NewRequest(method, link, payload)
    if err != nil { panic(err) }
    request.Header.Add("Content-Type", "application/json")

    response, err := client.Do(request)
    if err != nil { panic(err) }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil { panic(err) }
    return string(body)
}


func ScrapeVideosIntoFolder(links []string, folderName string, vidStats *VideoStats, vidJson *VideoJson, collector *colly.Collector, scriptCounter *int) {
    for _, link := range links {
        *scriptCounter = 0
        collector.Visit(link)
        vidStats.TransferJson(*vidJson)

        fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Title)
        if OnWindows() == true {
            fileTitle = ReplaceIllegalCharsWindows(fileTitle)
        } else {
            fileTitle = strings.ReplaceAll(fileTitle, "/", "")
        }
        err := os.MkdirAll(folderName, os.ModePerm)
        if err != nil { panic(err) }
        err = os.WriteFile(folderName + "/" + fileTitle, []byte(vidStats.Format()), 0644)
        if err != nil { panic(err) }
    }
}

func ScrapeVideos(links []string, vidStats *VideoStats, vidJson *VideoJson, collector *colly.Collector, scriptCounter *int) {
    for _, link := range links {
        *scriptCounter = 0
        collector.Visit(link)
        vidStats.TransferJson(*vidJson)

        fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Title)
        if OnWindows() == true {
            fileTitle = ReplaceIllegalCharsWindows(fileTitle)
        } else {
            fileTitle = strings.ReplaceAll(fileTitle, "/", "")
        }
        err := os.WriteFile(fileTitle, []byte(vidStats.Format()), 0644)
        if err != nil { panic(err) }
    }
}

func ScrapeSingleVideo(link string, vidStats *VideoStats, vidJson *VideoJson, collector *colly.Collector) {
    collector.Visit(link)
    vidStats.TransferJson(*vidJson)

    fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Title)
    if OnWindows() == true {
        fileTitle = ReplaceIllegalCharsWindows(fileTitle)
    } else {
        fileTitle = strings.ReplaceAll(fileTitle, "/", "")
    }
    err := os.WriteFile(fileTitle, []byte(vidStats.Format()), 0644)
    if err != nil { panic(err) }
}

func PrintMenu() {
    var menu string
    menu += fmt.Sprintf("1. To insert YouTube link.\n")
    menu += fmt.Sprintf("2. To bulk scrape multiple videos. (external .txt file)\n")
    menu += fmt.Sprintf("3. To scrape information & stadistics for all videos of a channel.\n")
    fmt.Println(menu)
}

func ReplaceIllegalCharsWindows(str string) string {
    chars := []string{ "<", ">", ":", "\"", "\\", "/", "|", "?", "*" }
    for _, char := range chars {
        str = strings.ReplaceAll(str, char, "")
    }
    return str
}

func OnWindows() bool {
    if runtime.GOOS == "windows" { return true }
    return false
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

/*
Getting the likes from the JSON string.

Example of how the JSON value looks:
"like this video along with 1,363 other people"
*/
func ExtractLikes(str string) int {
    var likesStr string
    index := FindFirstInt(str)
    likesStr = str[index:]
    final := strings.Index(likesStr, " ")
    likesStr = likesStr[:final]
    likesStr = strings.ReplaceAll(likesStr, ",", "")
    likes, err := strconv.ParseInt(likesStr, 10, 0)
    if err != nil { panic(err) }

    return int(likes)
}

//  ┌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┐
//  ╎ Video Stats Struc & Functions ╎
//  └╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┘

type VideoStats struct {
    Title string
    ViewCount int
    Likes int
    Comments int
    UploadDate string
    UploadHour string
    VideoID  string  
    LengthSeconds int
    Keywords []string 
    ChannelID string 
    Description string 
    IsCrawlable bool 
    AllowRatings bool 
    Author string 
    IsPrivate bool 
    IsLiveContent bool 
}

func NewVideoStats() VideoStats {
    vidStats := VideoStats{
        Title: "",
        ViewCount: 0,
        Likes: 0,
        UploadDate: "",
        UploadHour: "",
        LengthSeconds: 0,
        Description: "",
        Author: "",
        VideoID:  "",
        ChannelID: "",
        Keywords: []string{},
        IsCrawlable: false,
        AllowRatings: false,
        IsPrivate: false,
        IsLiveContent: false, 
    }
    return vidStats
}

// Formating the contents of the VideoStats struct into a human readable format.
func (v VideoStats) Format() string {
    var s string
    var keywords string

    lengthMinutes := SecondsToMinutes(v.LengthSeconds)
    if len(v.Keywords) == 0 {
        keywords = "none"
    } else {
        for _, keyword := range v.Keywords {
            keywords += fmt.Sprintf("%v, ", keyword)
        }
        keywords = keywords[:len(keywords)-2]
    }

    s += fmt.Sprintf("%-20v %v\n", "Title:", v.Title)
    s += fmt.Sprintf("%-20v %v\n", "View Count:", v.ViewCount)
    s += fmt.Sprintf("%-20v %v\n", "Likes:", v.Likes)
    s += fmt.Sprintf("%-20v %v\n", "Comments:", v.Comments)
    s += fmt.Sprintf("%-20v %v\n", "Upload Date:", v.UploadDate)
    s += fmt.Sprintf("%-20v %v\n", "Upload Hour:", v.UploadHour)
    s += fmt.Sprintf("%-20v %v\n", "Length:", lengthMinutes)
    s += fmt.Sprintf("%-20v %v\n", "Length (seconds):", v.LengthSeconds)
    s += fmt.Sprintf("%-20v %v\n", "Author:", v.Author)
    s += fmt.Sprintf("%-20v %v\n", "Video ID:", v.VideoID)
    s += fmt.Sprintf("%-20v %v\n", "Channel ID:", v.ChannelID)
    s += fmt.Sprintf("%-20v %v\n", "Keywords:", keywords)
    s += fmt.Sprintf("%-20v %v\n", "IsCrawlable:", v.IsCrawlable)
    s += fmt.Sprintf("%-20v %v\n", "AllowRatings:", v.AllowRatings)
    s += fmt.Sprintf("%-20v %v\n", "IsPrivate:", v.IsPrivate)
    s += fmt.Sprintf("%-20v %v\n", "IsLiveContent:", v.IsLiveContent)
    s += fmt.Sprintf("Description:\n\n%v", v.Description)
    return s
}

//  ┌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┐
//  ╎ Json Helper Functions ╎
//  └╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┘

// Transfering the information collected with VideoJson to VideoStats, trying to keep everything in a single place.
func (vidStats *VideoStats) TransferJson(json VideoJson) {
    secs, err := strconv.ParseInt(json.VideoDetails.LengthSeconds, 10, 0)
    if err != nil { panic(err) }
    views, err := strconv.ParseInt(json.VideoDetails.ViewCount, 10, 0)
    if err != nil { panic(err) }
    vidStats.ViewCount = int(views)
    vidStats.LengthSeconds = int(secs)
    vidStats.VideoID = json.VideoDetails.VideoID
    vidStats.Keywords = json.VideoDetails.Keywords
    vidStats.ChannelID = json.VideoDetails.ChannelID
    vidStats.Description = json.VideoDetails.Description
    vidStats.IsCrawlable = json.VideoDetails.IsCrawlable
    vidStats.AllowRatings = json.VideoDetails.AllowRatings
    vidStats.Author = json.VideoDetails.Author
    vidStats.IsPrivate = json.VideoDetails.IsPrivate
    vidStats.IsLiveContent = json.VideoDetails.IsLiveContent
}

// Function that gets the number of likes from the json, makes the code look cleaner. (the json is massive as you can see)
func GetLikesStr(likesJson LikesJson) string {
    return likesJson.Contents.TwoColumnWatchNextResults.Results.Results.Contents[0].VideoPrimaryInfoRenderer.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.LikeButtonViewModel.LikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.DefaultButtonViewModel.ButtonViewModel.AccessibilityText
}

func GetContinuationToken(continueTokenJson ContinueTokenJson) string {
    return continueTokenJson.Contents.TwoColumnWatchNextResults.Results.Results.Contents[3].ItemSectionRenderer.Contents[0].ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
}

func GetCommentCounter(commentCounterJson CommentCounterJson) string {
    return commentCounterJson.OnResponseReceivedEndpoints[0].ReloadContinuationItemsCommand.ContinuationItems[0].CommentsHeaderRenderer.CountText.Runs[0].Text
}
