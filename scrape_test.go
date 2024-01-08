package main

import (
    "testing"
    "fmt"
    "strings"
    "regexp"
)

func TestSingleDownload(t *testing.T) {
    link := "https://www.youtube.com/watch?v=Lwr3-doAgaI"

    collector, vidStats := createScraper()
    collector.Visit(link)

    fileTitle := fmt.Sprintf("%v - Stadistics.txt", vidStats.Title)
    if onWindows() == true {
        fileTitle = replaceIllegalCharsWindows(fileTitle)
    } else {
        fileTitle = strings.ReplaceAll(fileTitle, "/", "")
    }

    regExp := regexp.MustCompile(`^Title:\s*The Last Algorithms Course You'll Need by ThePrimeagen \| Preview\nView Count:\s*[1-9][0-9]*\nLikes:\s*[1-9][0-9]*\nComments:\s*[1-9][0-9]*\nUpload Date:\s*[1-9][0-9]*-[0-9]*-[0-9]*\nUpload Hour:\s*[0-9]*:[0-9]*:[0-9]*-[0-9]*:[0-9]*\nLength:\s*[0-9]*:[0-9]*\nLength \(seconds\):\s*[0-9]*\nAuthor:\s*Frontend Masters\nVideo ID:\s*Lwr3-doAgaI\nChannel ID:\s*UCzumJvwc0KBrdq4jpvOR7RA\nKeywords:\s*#Algorithms, #ThePrimeagen, #FrontendMasters\nIsCrawlable:\s*true\nAllowRatings:\s*true\nIsPrivate:\s*false\nIsLiveContent:\s*false\nDescription:[\S\s]*$`)
    format := vidStats.format()
    match := regExp.MatchString(format)
    if match != true {
        t.Errorf("The RegExp does not match, FAILED AT SINGLE LINK DOWNLOAD")
    }
} 
