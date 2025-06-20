package main

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	traq "github.com/traPtitech/go-traq"

	"fmt"
)

type Message struct {
	Content  string `json:"content"`
	UserName string `json:"userName"`
}

func main() {
	token := os.Getenv("TRAQ_TOKEN")

	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, token)

	e := echo.New()

	e.GET("/ubugoe/:userId", func(c echo.Context) error {
		userID := c.Param("userId")
		fmt.Println(userID)

		channelList, _, err := client.ChannelApi.GetChannels(auth).Path("gps/times/" + userID).
			Execute()
		if err != nil {
			c.Logger().Error(err)
			return c.String(500, "something wrong")
		}
		if len(channelList.Public) == 0 {
			return c.String(404, "No channel")
		}

		channel := channelList.Public[0]
		channelID := channel.Id

		fmt.Println(channel)
		fmt.Println(channelID)

		messages, _, _ := client.ChannelApi.GetMessages(auth, channelID).
			Since(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)).
			Limit(10).
			Order("asc").
			Execute()

		res := make([]Message, 0, len(messages))
		for _, message := range messages {

			userNameList, _, err := client.UserApi.GetUser(auth, message.UserId).Execute()
			if err != nil {
				c.Logger().Error(err)
				return c.String(500, "something wrong")
			}
			userName := userNameList.Name

			res = append(res, Message{
				Content:  message.Content,
				UserName: userName,
			})
		}

		return c.JSON(200, res)
	})

	e.Start(":8080")

}
