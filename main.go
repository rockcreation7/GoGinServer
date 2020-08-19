package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
)

// SMS fmt
type SMS struct {
	Phone string `form:"phone" binding:"required"`
	Msg   string `form:"msg" binding:"required"`
	Sub   string `form:"sub" binding:"required"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.GET("/color", func(c *gin.Context) {
		c.File("color.html")
	})

	r.GET("/smsform", func(c *gin.Context) {
		c.File("smsform.html")
	})

	r.POST("/", formHandler)

	r.StaticFS("/upload", http.Dir("upload"))

	r.POST("/upload", func(c *gin.Context) {
		img, err := imageupload.Process(c.Request, "file")

		if err != nil {
			panic(err)
		}

		err = img.Save("upload/" + img.Filename)

		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{
			"message": "uploadSuccess",
		})
	})

	r.POST("/smsform", func(c *gin.Context) {
		var sms SMS
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&sms); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := SendSMS(sms.Phone, sms.Msg, sms.Sub)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{"status": "SMS sent"})
	})

	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8000"
	}

	r.Run(":" + Port)
}

type myForm struct {
	Colors []string `form:"colors[]"`
}

func formHandler(c *gin.Context) {
	var fakeForm myForm
	c.Bind(&fakeForm)
	c.JSON(200, gin.H{"color": fakeForm.Colors})
}

// SendSMS unit
func SendSMS(phoneNumber string, message string, subject string) error {

	fmt.Println(subject)
	AwsRegion := "us-east-1"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion),
	},
	)

	svc := sns.New(sess)

	// Pass the phone number and message.
	params := &sns.PublishInput{
		PhoneNumber: aws.String(phoneNumber),
		Message:     aws.String(message),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"AWS.SNS.SMS.SenderID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(subject),
			},
		},
	}

	// sends a text message (SMS message) directly to a phone number.
	resp, err := svc.Publish(params)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	fmt.Println(resp) // print the response data.
	return nil
}
