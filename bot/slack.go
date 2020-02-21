package bot

import (
	"fmt"
	"github.com/slack-go/slack"
	"strings"

	"github.com/tempor1s/topofreddit/scraper"
)

/*
   CreateSlackClient sets up the slack RTM (real-time messaging) client library,
   initiating the socket connection and returning the client.
*/
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

/*
   RespondToEvents waits for messages on the Slack client's incomingEvents channel,
   and sends a response when it detects the bot has been tagged in a message with @<botTag>.
*/
func RespondToEvents(slackClient *slack.RTM) {
	for msg := range slackClient.IncomingEvents {
		// Log all events
		fmt.Println("Event Received: ", msg.Type)
		// Switch on the incoming event type
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// The bot's prefix (@topofreddit)
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			// Clean the prefix
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			// Register commands
			echoMessage(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
			sendSubreddits(slackClient, message, ev.Channel)
		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.DesktopNotificationEvent:
			fmt.Printf("Desktop Notification: %v\n", ev)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		default:
		}
	}
}

const helpMessage = "type in `@Reddit Top 5 <command_arg_1> <command_arg_2>` to run a command.\n\nCommands: `top all`\n`echo 'Hello there!'`"

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// echoMessage will just echo anything after the echo keyword.
func echoMessage(slackClient *slack.RTM, message, slackChannel string) {
	splitMessage := strings.Fields(strings.ToLower(message))

	if splitMessage[0] != "echo" {
		return
	}

	slackClient.SendMessage(slackClient.NewOutgoingMessage(strings.Join(splitMessage[1:], " "), slackChannel))
}

func sendSubreddits(slackClient *slack.RTM, message, slackChannel string) {
	splitMessage := strings.Fields(strings.ToLower(message))

	if splitMessage[0] != "top" {
		return
	}

	response := "Please pass in a subreddit name to get the top 5 posts for. Example: `@Reddit Top 5 top all`"
	if len(splitMessage) < 2 {
		slackClient.SendMessage(slackClient.NewOutgoingMessage(response, slackChannel))
	}

	posts := scraper.GetSubreddits(splitMessage[1])
	slackClient.SendMessage(slackClient.NewOutgoingMessage(posts, slackChannel))
}