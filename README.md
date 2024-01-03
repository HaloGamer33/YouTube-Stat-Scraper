# YouTube-Stat-Scrapper

![Static Badge](https://img.shields.io/badge/Made%20With-Go-blue?style=flat-square&logo=Go&logoColor=white) ![Static Badge](https://img.shields.io/badge/Made%20By-HaloGamer33-white?style=flat-square&label=Made%20by%20&color=%23e12a56) ![GitHub repo size](https://img.shields.io/github/repo-size/HaloGamer33/YouTube-Stat-Scrapper?style=flat-square&label=Size&color=success)

Go application that collects & downloads information about a YouTube video. It uses the Colly library to scrape data from the video's webpage and parses JSON data embedded in the page's HTML.

## Features

- Download description, title, view count, likes, upload date, author, video and channel ID, etc.
- (Coming soon) Download information & stadistics for all videos of a channel.
- (Coming soon) Bulk download information of multiple videos.

## Usage

1. Run the application.
2. When prompted, enter the URL of the YouTube video you want to analyze.
3. The application will print out the video's statistics and write them to a text file named "[Video Title] - Statistics.txt".

## Dependencies

- Go
- Colly

## Installation

1. Clone the repository.
2. Navigate to the project directory.
3. Run `go build` to compile the application.
4. Run the resulting executable.

## Code Overview

The main function initializes a new `VideoStats` object and a new Colly collector. It sets up several callbacks on the collector to handle various events such as visiting a URL, finding HTML elements, and receiving a response.

The `VideoStats` struct holds all the collected data about the video. The `Format` method formats this data into a human-readable string.

The `NewVideoStats` function returns a new `VideoStats` object with default values.

The `SecondsToMinutes` function converts a duration from seconds to a string in the format "minutes:seconds".

The `FindFirstInt` and `ExtractLikes` functions are helper functions used to parse the number of likes from a string.

The `TransferJson` method updates a `VideoStats` object with data extracted from a `VideoJson` object.

The `VideoJson` and `MyJSON` types define the structure of the JSON data embedded in the video's webpage.
