# YouTube-Stat-Scraper

![Static Badge](https://img.shields.io/badge/Made%20With-Go-blue?style=flat-square&logo=Go&logoColor=white) ![Static Badge](https://img.shields.io/badge/Made%20By-HaloGamer33-white?style=flat-square&label=Made%20by%20&color=%23e12a56) ![GitHub repo size](https://img.shields.io/github/repo-size/HaloGamer33/YouTube-Stat-Scrapper?style=flat-square&label=Size&color=success) ![GitHub release (with filter)](https://img.shields.io/github/v/release/HaloGamer33/YouTube-Stat-Scraper?style=flat-square&label=Release&color=%23ed9b37)

Go application that collects & downloads information about a YouTube video. It uses the Colly library to scrape data from the video's webpage and parses JSON data embedded in the page's HTML.

## Features

- Download description, title, view count, likes, upload date, author, video and channel ID, etc.
- Download information & stadistics for all videos of a channel.
- Bulk download information of multiple videos.

## Usage

1. Run the application.
2. When prompted, enter the URL of the YouTube video you want to analyze.
3. The application will write the video's statistics to a text file named "[Video Title] - Statistics.txt".

## Installation

1. Clone the repository.
2. Navigate to the project directory.
3. Run `go build` to compile the application.
4. Run the resulting executable.

## Dependencies

- Go
- Colly
