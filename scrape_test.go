package main

import (
    "testing"
    "regexp"
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

