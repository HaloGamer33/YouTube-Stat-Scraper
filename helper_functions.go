package main

import (
    "strings"
    "encoding/json"
    "strconv"
    // "fmt"
    "github.com/gocolly/colly"
    "fmt"
    "net/http"
    "io"
    "os"
    "runtime"
)

func getUsernameFromLink(link *string) string {
    startOfName := strings.Index(*link, "@")
    userName := (*link)[startOfName:]
    endOfName := strings.Index(userName, "/")
    if endOfName == -1 {
        return userName
    }
    userName = userName[:endOfName]
    return userName
}

func processScript20(scriptContents *string, vidStats *VideoStats) {
    jsonStr := removeJS(*scriptContents)

    // Getting number of likes

    var vidJson VideoJson
    err := json.Unmarshal([]byte(jsonStr), &vidJson)
    if err != nil { panic(err) }

    // Transfering the information collected with VideoJson to VideoStats, trying to keep everything in a single place.
    secs, err := strconv.ParseInt(vidJson.VideoDetails.LengthSeconds, 10, 0)
    if err != nil { panic(err) }
    views, err := strconv.ParseInt(vidJson.VideoDetails.ViewCount, 10, 0)
    if err != nil { panic(err) }
    vidStats.ViewCount = int(views)
    vidStats.LengthSeconds = int(secs)
    vidStats.VideoID = vidJson.VideoDetails.VideoID
    vidStats.Keywords = vidJson.VideoDetails.Keywords
    vidStats.ChannelID = vidJson.VideoDetails.ChannelID
    vidStats.Description = vidJson.VideoDetails.Description
    vidStats.IsCrawlable = vidJson.VideoDetails.IsCrawlable
    vidStats.AllowRatings = vidJson.VideoDetails.AllowRatings
    vidStats.Author = vidJson.VideoDetails.Author
    vidStats.IsPrivate = vidJson.VideoDetails.IsPrivate
    vidStats.IsLiveContent = vidJson.VideoDetails.IsLiveContent
}

func processScript45(scriptContents *string, vidStats *VideoStats) {
    jsonStr := removeJS(*scriptContents)

    var script45 script45Json
    err := json.Unmarshal([]byte(jsonStr), &script45)
    if err != nil { panic(err) }

    // Getting number of likes
    accessibilityText := script45.Contents.TwoColumnWatchNextResults.Results.Results.Contents[0].VideoPrimaryInfoRenderer.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.LikeButtonViewModel.LikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.DefaultButtonViewModel.ButtonViewModel.AccessibilityText

    var likesStr string
    var likes int
    index := findFirstInt(accessibilityText)
    if index == -1 {
        likes = 1
    } else {
        likesStr = accessibilityText[index:]
        final := strings.Index(likesStr, " ")
        likesStr = likesStr[:final]
        likesStr = strings.ReplaceAll(likesStr, ",", "")
        likes64, err := strconv.ParseInt(likesStr, 10, 0)
        if err != nil { panic(err) }
        likes = int(likes64)
    }

    vidStats.Likes = int(likes)

    // Getting number of comments
    
    var comments int

    if (len(script45.Contents.TwoColumnWatchNextResults.Results.Results.Contents[2].ItemSectionRenderer.Contents[0].MessageRenderer.Text.Runs) != 0 ) {
        // There is a text displaying "The comments are disabled" thus, no comments.
        comments = 0
    } else {
        token := script45.Contents.TwoColumnWatchNextResults.Results.Results.Contents[3].ItemSectionRenderer.Contents[0].ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
        continuationJson := pageLoadContinuation(token)

        var commentCounterJson CommentCounterJson
        err = json.Unmarshal([]byte(continuationJson), &commentCounterJson)
        noCommentsStr := commentCounterJson.OnResponseReceivedEndpoints[0].ReloadContinuationItemsCommand.ContinuationItems[0].CommentsHeaderRenderer.CountText.Runs[0].Text
        noCommentsStr = strings.ReplaceAll(noCommentsStr, ",", "")

        comments64, err := strconv.ParseInt(noCommentsStr, 10, 0)
        if err != nil { panic(err) }

        comments = int(comments64)
    }

    vidStats.Comments = int(comments)
}

func removeJS(s string) string {
    // Removing js from the json
    var indexJsonStart int = strings.Index(s, "{")
    return s[indexJsonStart:len(s)-1]
}


func createScraper() (*colly.Collector, *VideoStats) {
    collector := colly.NewCollector()
    var vidStats VideoStats
    var scriptCounter int
    collector.OnRequest(
        func (request *colly.Request) {
            scriptCounter = 0
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
    collector.OnHTML("script",
        func (element *colly.HTMLElement) {
            // // Writing all scripts for debugging
            // title := fmt.Sprintf("scripts/%v.txt", scriptCounter)
            // err := os.WriteFile(title, []byte(element.Text), 0644)
            // if err != nil { panic(err) }

            if scriptCounter == 20 {
                scriptContents := &element.Text
                processScript20(scriptContents, &vidStats)
            }
            if scriptCounter == 45 {
                scriptContents := &element.Text
                processScript45(scriptContents, &vidStats)
            }
            scriptCounter++
        },
    )
    return collector, &vidStats
}

func pageLoadChannelVideoContinuation(token string) string {
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

func pageLoadContinuation(token string) string {
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


func scrapeVideos(links []string, folderName string) {
    err := os.MkdirAll(folderName, os.ModePerm)
    if err != nil { panic(err) }

    for _, link := range links {
        collector, vidStats := createScraper()

        collector.Visit(link)

        fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Title)
        if onWindows() == true {
            fileTitle = replaceIllegalCharsWindows(fileTitle)
        } else {
            fileTitle = strings.ReplaceAll(fileTitle, "/", "")
        }
        err = os.WriteFile(folderName + "/" + fileTitle, []byte(vidStats.format()), 0644)
        if err != nil { panic(err) }
    }
}

func scrapeSingleVideo(link string) {
    collector, vidStats := createScraper()
    collector.Visit(link)

    fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Title)
    if onWindows() == true {
        fileTitle = replaceIllegalCharsWindows(fileTitle)
    } else {
        fileTitle = strings.ReplaceAll(fileTitle, "/", "")
    }
    err := os.WriteFile(fileTitle, []byte(vidStats.format()), 0644)
    if err != nil { panic(err) }
}

func printMenu() {
    var menu string
    menu += fmt.Sprintf("1. To insert YouTube link.\n")
    menu += fmt.Sprintf("2. To bulk scrape multiple videos. (external .txt file)\n")
    menu += fmt.Sprintf("3. To scrape information & stadistics for all videos of a channel.\n")
    fmt.Println(menu)
}

func replaceIllegalCharsWindows(str string) string {
    chars := []string{ "<", ">", ":", "\"", "\\", "/", "|", "?", "*" }
    for _, char := range chars {
        str = strings.ReplaceAll(str, char, "")
    }
    return str
}

func onWindows() bool {
    if runtime.GOOS == "windows" { return true }
    return false
}

func secondsToMinutes(inSeconds int) string {
    minutes := inSeconds / 60
    seconds := inSeconds % 60
    str := fmt.Sprintf("%v:%v", minutes, seconds)
    return str
}

func findFirstInt(str string) int {
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
