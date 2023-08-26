package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	embed "github.com/Clinet/discordgo-embed"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/drexedam/gravatar"
	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Title       string
	Description string
	Author      *gofeed.Person
	URL         string
}

func readFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func ParseMeetups() {
	url := "https://dou.ua/calendar/feed/%D0%BC%D0%B8%D1%82%D0%B0%D0%BF/"
	fp := gofeed.NewParser()
	fp.Client = &http.Client{Timeout: time.Second * 5}

	feed_items := make([]FeedItem, 1)

	for true {
		feed, err := fp.ParseURL(url)
		converter := md.NewConverter("", true, nil)
		converter.AddRules(md.Rule{
			Filter: []string{"img"},
			Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
				// You need to return a pointer to a string (md.String is just a helper function).
				// If you return nil the next function for that html element
				// will be picked. For example you could only convert an element
				// if it has a certain class name and fallback if not.
				content = strings.TrimSpace(content)
				return md.String("")
			},
		})

		if err == nil {
			l.Printf("[INFO] RSS Parser started to running for %s", url)
			items := feed.Items
			//Take only last item.
			item := items[len(items)-1]

			if !strings.Contains(readFile("feed_item.list"), item.Link) {
				// Create a new FeedItem Obj
				feedItem := FeedItem{
					Title:       item.Title,
					Description: item.Description,
					Author:      item.Authors[0],
					URL:         item.Link,
				}
				feed_items = append(feed_items, feedItem)

				description, err := converter.ConvertString(feedItem.Description)

				AvatarURL := gravatar.New(feedItem.Author.Email).
					DefaultURL("https://cdn.icon-icons.com/icons2/2438/PNG/512/boy_avatar_icon_148455.png").
					Size(200).
					Default(gravatar.NotFound).
					Rating(gravatar.Pg).
					AvatarURL()

				genericEmbed := embed.NewEmbed().
					SetTitle(feedItem.Title).
					SetDescription(description).
					SetURL(feedItem.URL).
					SetAuthor(feedItem.Author.Name, AvatarURL).
					SetColor(0x1c1c1c)

				Dg.ChannelMessageSendEmbed("1145004345545474129", genericEmbed.MessageEmbed)

				// Write Link info to the feed_item.list file
				file, err := os.OpenFile("feed_item.list", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					panic(err)
				}

				defer file.Close()

				if _, err := file.WriteString(item.Link + "\n"); err != nil {
					l.Fatal(err)
				}

			} else {
				l.Printf("[WARN] The parser already parsed %s", item.Link)
			}
		}
		feed_items = make([]FeedItem, 1)

		time.Sleep(300 * time.Second)
	}
}
