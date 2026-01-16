package twitch

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	twitch2 "github.com/joeyak/go-twitch-eventsub/v3"

	"undertale-tts/internal/charselector"
	"undertale-tts/internal/common"
	"undertale-tts/internal/reproductor"
	"undertale-tts/internal/whitelist"
)

const TWITCH_URL = "https://id.twitch.tv/oauth2/authorize" +
	"?response_type=token" +
	"&client_id=4u4v1h8d2yfvftoqtstu0pley1pooo" +
	"&redirect_uri=http://localhost:9000" +
	"&scope=chat:read+chat:edit+channel:read:redemptions" +
	"&state=c3ab8aa609ea11e793ae92361f002671"

const TWITCH_BROADCASTER_ID = "https://api.twitch.tv/helix/users"

const CLIENT_ID = "4u4v1h8d2yfvftoqtstu0pley1pooo"

const IRC_TWITCH_SERVER = "irc.chat.twitch.tv:6667"

type RedemptionEvent struct {
	UserId               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserId    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Id                   string `json:"id"`
	UserInput            string `json:"user_input"`
	Status               string `json:"status"`
	Reward               Reward `json:"reward"`
	RedeemedAt           string `json:"redeemed_at"`
}

type Reward struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Cost   int    `json:"cost"`
	Prompt string `json:"prompt"`
}

type Broadcaster struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadCasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	CreatedAt       string `json:"created_at"`
}

type Response struct {
	Data []Broadcaster `json:"data"`
}

type RedemptionCallback func(event RedemptionEvent)

var TwitchConnectButton *widget.Button
var TwitchConfWindow fyne.Window
var TwitchRedeemName *widget.Entry
var ConnectToTwitch *widget.Button
var IsTwitchConnected = false
var TwitchActiveRedeemOption *widget.Check
var twitchConfiWindowsIsOpen = false
var TwitchErrorWindow fyne.Window
var Active = false

var IsRedeemOptionActive = false

var CTX context.Context

func GetAuthorization() {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", TWITCH_URL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", TWITCH_URL).Start()
	case "darwin":
		err = exec.Command("open", TWITCH_URL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		//logging.CreateLog("twtich - unsupported platform", err)
		log.Fatal(err)
	}
}

func InitTwitchConfigWindows(app fyne.App) {
	TwitchConfWindow = app.NewWindow("Twitch configuration")

	TwitchConfWindow.SetOnClosed(func() {
		twitchConfiWindowsIsOpen = false
	})

	TwitchRedeemName = widget.NewEntry()
	TwitchRedeemName.SetPlaceHolder("Enter Redeem's name")
	TwitchRedeemName.Disabled()
	TwitchRedeemName.Resize(fyne.NewSize(100, TwitchRedeemName.MinSize().Width))

	form := widget.NewForm(
		widget.NewFormItem("Twitch Redeem's name", TwitchRedeemName),
	)

	radio := widget.NewRadioGroup([]string{"Read all messages in chat", "Use channel points"}, func(value string) {
		if value == "Read all messages in chat" {
			IsRedeemOptionActive = false
			TwitchRedeemName.Disable()
		} else {
			IsRedeemOptionActive = true
			TwitchRedeemName.Enable()
		}
		form.Refresh()
	})

	radio.SetSelected("Read all messages in chat")

	ConnectToTwitch = widget.NewButton("Connect to Twitch", func() {

		if IsRedeemOptionActive && TwitchRedeemName.Text == "" {
			initTwitchErrorWindow(app, "You must enter a Redeem's name")
			TwitchErrorWindow.Show()
			return
		}
		GetAuthorization()
		TwitchConfWindow.Close()
	})

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		ConnectToTwitch,
	)

	TwitchConfWindow.SetContent(
		container.NewVBox(
			radio,
			form,
			centeredButton,
		),
	)

	TwitchConfWindow.Resize(fyne.NewSize(400, 100))

	TwitchConnectButton = widget.NewButton("Connect to Twitch", func() {
		if !twitchConfiWindowsIsOpen {
			twitchConfiWindowsIsOpen = true
			TwitchConfWindow.Show()
		}
	})

}

func SubscribeToChat(token string) {
	broadcasterId, login, err := GetBroadcasterId(token)
	if err != nil {
		//logging.CreateLog("twitch - couldn't get broadcaster info", err)
		log.Fatal("Error retrieving broadcaster id", err)
	}

	if IsRedeemOptionActive {
		subscribeToEvent(token, broadcasterId)

	}

	if !IsRedeemOptionActive {
		conn, err := net.Dial("tcp", IRC_TWITCH_SERVER)
		if err != nil {
			//logging.CreateLog("twitch - couldn't connect to twitch chat", err)
			log.Fatal("Error conectandose a IRC", err)
		}

		fmt.Fprintf(conn, "PASS %s\r\n", "oauth:"+token)
		fmt.Fprintf(conn, "NICK %s\r\n", login)
		fmt.Fprintf(conn, "JOIN #%s\r\n", login)
		reader := bufio.NewReader(conn)
		go func() {
			for {

				select {
				case <-CTX.Done():
					return
				default:
					line, err := reader.ReadString('\n')
					if strings.HasPrefix(line, "PING") {
						fmt.Fprintf(conn, "PONG :tmi.twitch.tv\r\n")
					}
					if err != nil {
						//logging.CreateLog("twitch - couldn't get new message in chat", err)
						log.Fatal(err)
					}
					splitted := strings.Split(line, "#")
					if len(splitted) == 2 {
						user := strings.Split(splitted[0], "!")[0]
						user = user[1:]
						message := strings.Split(splitted[1], ":")
						if len(message) == 2 {
							if Active {
								chatMsg := message[1]
								if common.IsCommandActive && strings.HasPrefix(chatMsg, common.TTSCommand) {
									if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(user) {
										reproductor.Reproduce(chatMsg, charselector.CurrentCharacter)
									} else if !whitelist.IsWhitelistActive {
										reproductor.Reproduce(chatMsg, charselector.CurrentCharacter)
									}
								} else if !common.IsCommandActive {
									if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(user) {
										reproductor.Reproduce(chatMsg, charselector.CurrentCharacter)
									} else if !whitelist.IsWhitelistActive {
										reproductor.Reproduce(chatMsg, charselector.CurrentCharacter)
									}
								}

							} else if strings.Compare(strings.TrimSpace(message[1]), "End of /NAMES list") == 0 && !Active {
								Active = true
							}

						}
					}
				}

			}
		}()

	}
	SetConnected()

}

func GetBroadcasterId(token string) (string, string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", TWITCH_BROADCASTER_ID, nil)
	if err != nil {
		//logging.CreateLog("twitch - couldn't create HTTP request", err)
		log.Fatal("Error creating request", err)
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", CLIENT_ID)

	resp, err := client.Do(req)
	if err != nil {
		//logging.CreateLog("twitch - couldn't make HTTP request", err)
		log.Fatal("Error sending request", err)
		return "", "", err
	}
	defer resp.Body.Close()

	var result Response

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		//logging.CreateLog("twitch - couldn't deserealize response", err)
		log.Fatal("Error decoding JSON response", err)
		return "", "", err
	}

	if len(result.Data) > 0 {
		b := result.Data[0]
		broadcasterid := b.Id
		login := b.Login
		return broadcasterid, login, nil
	}

	return "", "", errors.New("No data in")
}

func subscribeToEvent(accessToken string, userID string) {
	client := twitch2.NewClient()
	client.OnError(func(err error) {
		//logging.CreateLog("twitch - On Error Subscribe to event - ", err)
		log.Fatal("twitch - On Error Subscribe to even", err)
	})
	client.OnWelcome(func(message twitch2.WelcomeMessage) {

		var events []twitch2.EventSubscription

		if IsRedeemOptionActive {
			ev := []twitch2.EventSubscription{
				"channel.channel_points_custom_reward_redemption.add",
			}
			events = append(events, ev...)
		}

		for _, event := range events {
			_, err := twitch2.SubscribeEvent(twitch2.SubscribeRequest{
				SessionID:   message.Payload.Session.ID,
				ClientID:    CLIENT_ID,
				AccessToken: accessToken,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id": userID,
				},
			})
			if err != nil {
				//logging.CreateLog("twitch - Event error - ", err)
				return
			}
		}
	})

	client.OnRawEvent(func(event string, metadata twitch2.MessageMetadata, subscription twitch2.PayloadSubscription) {

		var r RedemptionEvent
		if err := json.Unmarshal([]byte(event), &r); err != nil {
			panic(err)
		}
		if TwitchRedeemName.Text == r.Reward.Title {
			reproductor.Reproduce(r.UserInput, charselector.CurrentCharacter)
		}
	})

	go func() {
		if err := client.Connect(); err != nil {
			//logging.CreateLog("twitch - Could not connect client - ", err)
			log.Fatal("twitch - Could not connect client - ", err)
		}
	}()

	go func() {
		<-CTX.Done()
		client.Close()
	}()
}

func initTwitchErrorWindow(app fyne.App, msg string) {
	TwitchErrorWindow = app.NewWindow("!Error")
	TwitchErrorWindow.SetContent(widget.NewLabel(msg))
}

func SetConnected() {
	fyne.Do(func() {
		TwitchConnectButton.SetText("Connected")
		TwitchConnectButton.Disable()
	})
}
