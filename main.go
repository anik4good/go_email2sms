package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	Configuration "github.com/anik4good/go_email2sms/config"
	"github.com/anik4good/go_email2sms/models"
	Route "github.com/anik4good/go_email2sms/routes"

	"gitlab.com/hartsfield/gmailAPI"
	"gitlab.com/hartsfield/inboxer"
	gmail "google.golang.org/api/gmail/v1"
)

var database *sql.DB
var logger *log.Logger

func main() {
	fmt.Println("hello world")
	logger = Configuration.InitLogger()
	//	Configuration.InitDatabaseMysql()
	Configuration.InitDatabaseMysqlGorm()

	app := fiber.New()
	app.Use(requestid.New())
	Route.SetUpRoutes(app)
	app.Use(cors.New())
	// Default middleware config

	// app.Use(func(c *fiber.Ctx) error {
	// 	return c.SendStatus(404) // => 404 "Not Found"
	// })

	// go log.Fatal(app.Listen(":3000"))
	go app.Listen("127.0.0.1:3000")
	time.Sleep(time.Second)
	//util.UserSeed()
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

	for {
		// sendToQueue()
		// CheckProfile()

		msgs, err := inboxer.Query(srv, "label:UNREAD after:2022/01/01 before:2023/12/31")

		if err != nil {
			fmt.Println(err)
		}

		if int(CheckUnread()) > 0 {
			fmt.Printf("You have %d unread emails.", CheckUnread())
			fmt.Println("")
			// Range over the messages
			for _, msg := range msgs {

				fmt.Println(msg.Id)
				time, err := inboxer.ReceivedTime(msg.InternalDate)
				if err != nil {
					fmt.Println(err)
				}

				md := inboxer.GetPartialMetadata(msg)
				body, err := inboxer.GetBody(msg, "text/plain")
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println("From: ", md.From)
				fmt.Println("Date: ", time)
				fmt.Println(body)
				SendToTrash(msg.id)
			}

			// MarkAsRead()

			// queue := []models.Queue{}
			// result := Configuration.GormDBConn.Where("status = ?", 0).Limit(1).Last(&queue)
			// if result.RowsAffected == 0 {
			// 	fmt.Println("No data found on Queue table")
			// }

			// //log.Println(int(result.RowsAffected))

			// if int(result.RowsAffected) > 0 {
			// 	for _, q := range queue {
			// 		//go processEmail(queue)
			// 		//logger.Println("Status change for id ", q.ID)
			// 		fmt.Println("Queued Data: ", q.Name)
			// 		q.Status = 2
			// 		Configuration.GormDBConn.Updates(&q)
			// 	}
		} else {
			fmt.Printf("You have 0 unread emails.")
			fmt.Println("")
		}

		time.Sleep(5 * time.Second)
		// fmt.Printf(">")
	}

}

// func processEmail(queue models.Queue) {
// 	err := sendEmail(queue)
// 	if err != nil {
// 		logger.Println(err)
// 		return
// 	}

// 	 changeStatusToSuccess(queue.ID)

// }

// func sendEmail(queue models.Queue) error {
// 	//	logger.Println("Sending sms to", queuedEmail.Email)
// 	fmt.Println("Sending Email to", queue.Email)

// 	//	send_email(queuedEmail)

// 	return nil
// }

//func checkForNewRecords() *gorm.DB {
//	users := []models.User{}
//	// rows, err := database.Query("select id, name, email from queues WHERE status = 0 LIMIT 10")
//	// if err != nil {
//	// 	logger.Println("Error on new records checking ..", err)
//	// }
//	result:= Configuration.GormDBConn.Where("status = ?", 0).Limit(1).Last(&users)
//
//	return result
//}

func sendToQueue() {
	// rows, err := database.Query("select id, name, email from users WHERE status = 0 LIMIT 10")
	// queue := new(models.Queue)
	// users := new(models.User)

	users := []models.User{}
	var queue models.Queue
	// result := Configuration.GormDBConn.Raw("select id, name, email from users WHERE status = 0 LIMIT 100").Scan(&queue)
	result := Configuration.GormDBConn.Where("status = ?", 0).Last(&users)
	//log.Println(int(result.RowsAffected))
	if result.Error != nil {
		//fmt.Println("No data found on User table")
	}

	// user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	//Configuration.GormDBConn.Create(queues)

	// if err != nil {
	// 	logger.Println("Error on new records checking ..", err)
	// }
	for _, rows := range users {
		queue.ID = rows.ID
		queue.Name = rows.Name
		queue.Email = rows.Email
		queue.Status = rows.Status

		Configuration.GormDBConn.Create(&queue)

		rows.Status = 2

		Configuration.GormDBConn.Updates(&rows)
		// return result.RowsAffected
		//log.Println(result)
	}
}

// func changeStatusToPending(id int) {
// 	_, err := database.Exec("UPDATE queues SET status = ? WHERE id = ?", 2, id)
// 	if err != nil {
// 		logger.Println("Error updating status of "+string(rune(id))+" in users", (err), " to pending")
// 		return
// 	}
// }

// func changeStatusToSuccess(id int) {
// 	_, err := database.Exec("UPDATE queues SET status = ? WHERE id = ?", 3, id)
// 	if err != nil {
// 		logger.Println("Error updating status of "+string(rune(id))+" in users", (err), " to pending")
// 		return
// 	}
// }

func CheckUnread() int64 {
	// Connect to the gmail API service.
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

	// num will be -1 on err
	num, err := inboxer.CheckForUnread(srv)

	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("You have %d unread emails.", int(num))
	return num

}

func CheckProfile() {

	url := "https://www.googleapis.com/gmail/v1/users/anik4nobody@gmail.com/profile"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer ya29.A0ARrdaM_Fe29KPaKQM6cqr9f6PBzdf8pKp2RYf68BrExYgTzPNUOyAhkYCg1vqMEQ2EHKtkJ9tCUmm5qCKqwmHkpW2EYq_FDEpVaHCfe-FB7cGCoEF2oOlEnD8frVFyOfokCmmsC01lRIpgoW4yj7UpA3GO0F")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func MarkAsRead() {
	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

	msgs, err := inboxer.Query(srv, "label:UNREAD after:2022/01/01 before:2023/12/31")

	if err != nil {
		fmt.Println(err)
	}

	req := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
		AddLabelIds:    []string{"TRASH"}}

	// Range over the messages
	for _, msg := range msgs {
		msg, err := inboxer.MarkAs(srv, msg, req)
		if err != nil {
			fmt.Println(err)
		}
		body, err := inboxer.GetBody(msg, "text/plain")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(body + "is marked as read")

	}

}

func SendToTrash(id string) {

	ctx := context.Background()
	srv := gmailAPI.ConnectToService(ctx, gmail.GmailComposeScope)

	msg, err := srv.Users.Messages.Get("me", id).Do()
	if err != nil {
		fmt.Println(err)
	}

}
