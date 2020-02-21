package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

// Post represents a top post on reddit.
type Post struct {
	Title     string `json:"title"`
	CreatedAt string `json:"createdat"`
	Url       string `json:"url"`
	Upvotes   string `json:"upvotes"`
}

// Posts are a collection of post structs.
type Posts struct {
	Posts []Post `json:"posts"`
}

// scrapeSubreddits will scrape the top 5 posts from any subreddit!
func scrapeSubreddits(subredditName string) Posts {
	// Instantiate default collector
	c := colly.NewCollector()
	var posts Posts

	// On every a element which has href attribute call callback
	c.OnHTML(".rpBJOHq2PR60pnwJlUyP0", func(wrapper *colly.HTMLElement) {
		// Get all the titles
		wrapper.ForEach("div.y8HYJ-y_lTUHkQIc1mdCq._2INHSNB8V5eaWp4P0rY_mE", func(i int, title *colly.HTMLElement) {
			if i < 5 {
				posts.Posts = append(posts.Posts, Post{Title: title.Text})
			}
		})

		// Get URL and created at.
		wrapper.ForEach("div._3AStxql1mQsrZuUIFP9xSg.nU4Je7n-eSXStTBAPMYt8 > a", func(j int, subreddit *colly.HTMLElement) {
			if j < 5 {
				posts.Posts[j].CreatedAt = subreddit.Text
				posts.Posts[j].Url = subreddit.Attr("href")
			}
		})

		wrapper.ForEach("div._23h0-EcaBUorIHC-JZyh6J > div > div", func(k int, upvotes *colly.HTMLElement) {
			if k < 5 {
				posts.Posts[k].Upvotes = upvotes.Text
			}
		})
	})

	// Start scraping a reddit URL.
	if err := c.Visit(fmt.Sprintf("https://www.reddit.com/r/%s/", subredditName)); err != nil {
		log.Fatalf("Error occurred when visiting URL. Error: %v\n", err)
	}

	return posts
}

// GetSubreddits will pull the top 5 posts from any subreddit and format them to look nice in slack/discord.
func GetSubreddits(subreddit string) string {
	posts := scrapeSubreddits(subreddit)

	stringBuilder := fmt.Sprintf("Recent top 5 trending posts for reddit.com/r/%s\n", subreddit)
	stringBuilder += "```"

	for i, post := range posts.Posts {
		stringBuilder += "Title: " + post.Title + "\n"
		stringBuilder += "URL: " + post.Url + "\n"
		stringBuilder += "Created At: " + post.CreatedAt + "\n"
		stringBuilder += "Upvotes: " + post.Upvotes + "\n"

		if i != 4 {
			stringBuilder += "\n"
		}
	}

	stringBuilder += "```"
	return stringBuilder
}