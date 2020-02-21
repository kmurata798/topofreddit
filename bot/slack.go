package bot

import (
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

const helpMessage = "type in '@legacybotboi <command_arg_1> <command_arg_2>'"

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
		fmt.Println("Event Received: ", msg.Type)
		// Switch on the incoming event type
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(ev.Msg.Text, botTagString) {
				continue
			}
			message := strings.Replace(ev.Msg.Text, botTagString, "", -1)

			echoMessage(slackClient, message, ev.Channel)
			sendHelp(slackClient, message, ev.Channel)
			// Any
		default:
		}
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, message, slackChannel string) {
	if strings.ToLower(message) != "help" {
		return
	}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendResponse is NOT unimplemented --- write code in the function body to complete!

func echoMessage(slackClient *slack.RTM, message, slackChannel string) {
	command := strings.ToLower(message)
	println("[RECEIVED] sendResponse:", command)

	//if strings.ToLower(message) != "echo" {
	//	return
	//}
	slackClient.SendMessage(slackClient.NewOutgoingMessage(message, slackChannel))
}
