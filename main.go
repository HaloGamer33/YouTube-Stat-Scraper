package main

import (
    "fmt"
    "github.com/gocolly/colly"
    "os"
    "encoding/json"
    "strings"
)

func main() { 
    //  ┌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┐
    //  ╎ User Input Start ╎
    //  └╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┘
    var selection string
    printMenu()
    fmt.Scanln(&selection)
    switch selection {
    case "1":
        fmt.Println("Insert the link:")
        fmt.Scanln(&selection)
        scrapeSingleVideo(selection)
    case "2":
        fmt.Println("Name of the .txt file containing the links (include the .txt):")
        fmt.Scanln(&selection)
        fileContents, err := os.ReadFile(selection)
        stringContents := string(fileContents)
        links := strings.Split(stringContents, "\n")
        if err != nil { panic(err) }
        scrapeVideos(links, "Videos")
    case "3":
        fmt.Println("Channel link:")
        fmt.Scanln(&selection)
        channelLink := &selection
        channelUsername := getUsernameFromLink(channelLink)

        channelVideosLink := fmt.Sprintf("https://www.youtube.com/%v/videos", channelUsername)
        folderName := fmt.Sprintf("%v - Channel Video Stats", channelUsername)

        links, token := scrapeChannelVideosPage(channelVideosLink)

        if token == "" {
            scrapeVideos(links, folderName)
        } else {
            for token != "" {
                var scrapedLinks []string
                scrapedLinks, token = scrapeUntilEndOfPage(token)
                links = append(links, scrapedLinks...)
            }
            scrapeVideos(links, folderName)
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

func scrapeUntilEndOfPage(continuationToken string) ([]string, string) {
    continuationVideosJsonStr := pageLoadChannelVideoContinuation(continuationToken)
    var continuationVideosJson ContinuationVideosJson
    err := json.Unmarshal([]byte(continuationVideosJsonStr), &continuationVideosJson)
    if err != nil { panic(err) }

    items := continuationVideosJson.OnResponseReceivedActions[0].AppendContinuationItemsAction.ContinuationItems

    var links []string
    for index, item := range items {
        if index == len(items) - 1 {
            continuationToken = item.ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
            if continuationToken != "" { continue }
        }

        id := &item.RichItemRenderer.Content.VideoRenderer.VideoId
        link := fmt.Sprintf("https://www.youtube.com/watch?v=%v", *id)
        links = append(links, link)
    }

    return links, continuationToken
}

func scrapeChannelVideosPage(channelLink string) ([]string, string) {
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
                jsonStr := removeJS(element.Text)

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


/*
Getting the likes from the JSON string.

Example of how the JSON value looks:
"like this video along with 1,363 other people"
*/

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

func newVideoStats() VideoStats {
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
func (v VideoStats) format() string {
    var s string
    var keywords string

    lengthMinutes := secondsToMinutes(v.LengthSeconds)
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
