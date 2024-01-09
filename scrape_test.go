package main

import (
    "testing"
    "regexp"
    "fmt"
)

const singleDownloadRegex = `^Title:\s*.+
View Count:\s*[1-9][0-9]*
Likes:\s*[1-9][0-9]*
Comments:\s*[1-9][0-9]*
Upload Date:\s*[1-9][0-9]*-[0-1][0-9]*-[0-3][0-9]*
Upload Hour:\s*[0-9]*:[0-9]*:[0-9]*-[0-9]*:[0-9]*
Length:\s*[0-9]*:[0-9]*
Length \(seconds\):\s*[1-9][0-9]*
Author:\s*.+
Video ID:\s*.+
Channel ID:\s*.+
Keywords:\s*.+
IsCrawlable:\s*(true|false)
AllowRatings:\s*(true|false)
IsPrivate:\s*(true|false)
IsLiveContent:\s*(true|false)
Description:[\S\s]*$`

const commentsDisabledRegex = `^Title:\s*.+
View Count:\s*[1-9][0-9]*
Likes:\s*[1-9][0-9]*
Comments:\s*0
Upload Date:\s*[1-9][0-9]*-[0-1][0-9]*-[0-3][0-9]*
Upload Hour:\s*[0-9]*:[0-9]*:[0-9]*-[0-9]*:[0-9]*
Length:\s*[0-9]*:[0-9]*
Length \(seconds\):\s*[1-9][0-9]*
Author:\s*.+
Video ID:\s*.+
Channel ID:\s*.+
Keywords:\s*.+
IsCrawlable:\s*(true|false)
AllowRatings:\s*(true|false)
IsPrivate:\s*(true|false)
IsLiveContent:\s*(true|false)
Description:[\S\s]*$`

func TestSingleDownload(t *testing.T) {
    link := "https://www.youtube.com/watch?v=Lwr3-doAgaI"

    collector, vidStats := createScraper()
    collector.Visit(link)

    regExp := regexp.MustCompile(singleDownloadRegex)
    format := vidStats.format()
    match := regExp.MatchString(format)
    if match != true {
        t.Errorf("The RegExp does not match, FAILED AT SINGLE LINK DOWNLOAD")
    }
} 

func TestSingleDownloadNoComments(t *testing.T) {
    link := "https://www.youtube.com/watch?v=bOYl_FjqG_0"

    collector, vidStats := createScraper()
    collector.Visit(link)

    regExp := regexp.MustCompile(commentsDisabledRegex)
    format := vidStats.format()
    match := regExp.MatchString(format)
    if match != true {
        t.Errorf("The RegExp does not match, FAILED AT DOWNLOADING VIDEO WITH COMMENTS DISABLED")
    }
}

func TestChannelScrape(t *testing.T) {
    channelLink := "https://www.youtube.com/@Bordercollie_5"
    // https://www.youtube.com/@Bordercollie_5
    channelUsername := getUsernameFromLink(&channelLink)

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
    //
    // file, err := os.Open("./@Bordercollie_5 - Channel Video Stats")
    // if err != nil { panic(err) }
    //
    // defer file.Close()
    //
    // names, err := file.Readdirnames(0)
    // if err != nil { panic(err) }
    //
    // for _, v := range names {
    //     fmt.Println(v)
    // }
}
