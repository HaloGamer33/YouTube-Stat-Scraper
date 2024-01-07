package main

type CommentCounterJson struct {
    OnResponseReceivedEndpoints []struct {
        ReloadContinuationItemsCommand struct {
            ContinuationItems []struct {
                CommentsHeaderRenderer struct {
                    CountText struct {
                        Runs []struct {
                            Text string `json:"text"`
                        } `json:"runs"`
                    } `json:"countText"`
                } `json:"commentsHeaderRenderer"`
            } `json:"continuationItems"`
        } `json:"reloadContinuationItemsCommand"`
    } `json:"onResponseReceivedEndpoints"`
}

type script45Json struct {
    Contents struct {
        TwoColumnWatchNextResults struct {
            Results struct {
                Results struct {
                    Contents []struct {
                        VideoPrimaryInfoRenderer struct { 
                            VideoActions struct { 
                                MenuRenderer struct { 
                                    TopLevelButtons []struct { 
                                        SegmentedLikeDislikeButtonViewModel struct { 
                                            LikeButtonViewModel struct { 
                                                LikeButtonViewModel struct { 
                                                    ToggleButtonViewModel struct { 
                                                        ToggleButtonViewModel struct { 
                                                            DefaultButtonViewModel struct { 
                                                                ButtonViewModel struct { 
                                                                    // Extracting string that contains the likes
                                                                    AccessibilityText string `json:"accessibilityText"` 
                                                                } `json:"buttonViewModel"` 
                                                            } `json:"defaultButtonViewModel"` 
                                                        } `json:"toggleButtonViewModel"` 
                                                    } `json:"toggleButtonViewModel"` 
                                                } `json:"likeButtonViewModel"` 
                                            } `json:"likeButtonViewModel"` 
                                        } `json:"segmentedLikeDislikeButtonViewModel"` 
                                    } `json:"topLevelButtons"` 
                                } `json:"menuRenderer"` 
                            } `json:"videoActions"` 
                        } `json:"videoPrimaryInfoRenderer"` 
                        ItemSectionRenderer struct {
                            Contents []struct {
                                ContinuationItemRenderer struct {
                                    ContinuationEndpoint struct {
                                        ContinuationCommand struct {
                                            // Extracting continuation token
                                            Token string `json:"token"`
                                        } `json:"continuationCommand"`
                                    } `json:"continuationEndpoint"`
                                } `json:"continuationItemRenderer"`
                                MessageRenderer struct {
                                    Text struct {
                                        Runs []struct {
                                            // Extracting text that says "comments are disabled" if there are no comments
                                            Text string
                                        } `json:"runs"`
                                    } `json:"text"`
                                } `json:"messageRenderer"`
                            } `json:"contents"`
                        } `json:"itemSectionRenderer"`
                    } `json:"contents"`
                } `json:"results"`
            } `json:"results"`
        } `json:"twoColumnWatchNextResults"`
    } `json:"contents"`
}

type ContinuationVideosJson struct {
    OnResponseReceivedActions []struct {
        AppendContinuationItemsAction struct {
            ContinuationItems []struct {
                RichItemRenderer struct {
                    Content struct {
                        VideoRenderer struct {
                            VideoId string `json:"videoId"`
                        } `json:"videoRenderer"`
                    } `json:"content"`
                } `json:"richItemRenderer"`
                ContinuationItemRenderer struct {
                    ContinuationEndpoint struct {
                        ContinuationCommand struct {
                            Token string `json:"token"`
                        } `json:"continuationCommand"`
                    } `json:"continuationEndpoint"`
                } `json:"continuationItemRenderer"`
            } `json:"continuationItems"`
        } `json:"appendContinuationItemsAction"`
    } `json:"onResponseReceivedActions"`
}

type ChannelVideosJson struct {
    Contents struct {
        TwoColumnBrowseResultsRenderer struct {
            Tabs []struct {
                TabRenderer struct {
                    Content struct {
                        RichGridRenderer struct {
                            Contents []struct {
                                RichItemRenderer struct {
                                    Content struct {
                                        VideoRenderer struct {
                                            VideoId string `json:"videoId"`
                                        } `json:"videoRenderer"`
                                    } `json:"content"`
                                } `json:"richItemRenderer"`
                                ContinuationItemRenderer struct {
                                    ContinuationEndpoint struct {
                                        ContinuationCommand struct {
                                            Token string `json:"token"`
                                        } `json:"continuationCommand"`
                                    } `json:"continuationEndpoint"`
                                } `json:"continuationItemRenderer"`
                            } `json:"contents"`
                        } `json:"richGridRenderer"`
                    } `json:"content"`
                } `json:"tabRenderer"`
            } `json:"tabs"`
        } `json:"twoColumnBrowseResultsRenderer"`
    } `json:"contents"`
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


