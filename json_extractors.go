package main

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
type LikesJson struct { Contents struct { TwoColumnWatchNextResults struct { Results struct { Results struct { Contents []struct { VideoPrimaryInfoRenderer struct { VideoActions struct { MenuRenderer struct { TopLevelButtons []struct { SegmentedLikeDislikeButtonViewModel struct { LikeButtonViewModel struct { LikeButtonViewModel struct { ToggleButtonViewModel struct { ToggleButtonViewModel struct { DefaultButtonViewModel struct { ButtonViewModel struct { AccessibilityText string `json:"accessibilityText"` } `json:"buttonViewModel"` } `json:"defaultButtonViewModel"` } `json:"toggleButtonViewModel"` } `json:"toggleButtonViewModel"` } `json:"likeButtonViewModel"` } `json:"likeButtonViewModel"` } `json:"segmentedLikeDislikeButtonViewModel"` } `json:"topLevelButtons"` } `json:"menuRenderer"` } `json:"videoActions"` } `json:"videoPrimaryInfoRenderer"` } `json:"contents"` } `json:"results"` } `json:"results"` } `json:"twoColumnWatchNextResults"` } `json:"contents"` }

