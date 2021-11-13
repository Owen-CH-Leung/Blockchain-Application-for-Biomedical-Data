package email

import (
	"bytes"
	"context"
	"download-view/backend/logger"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/gomail.v2"
)

func SendEmailViaAWS(recipient string, attachment ...string) {
	var log = logger.BlockchainLog
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		log.Println("Unable to detect caller name when sending email")
		os.Exit(1)
	}
	caller := runtime.FuncForPC(pc)
	log.Println("Caller Name is ", caller.Name())
	var email_header string
	switch caller.Name() {
	case "download-view/frontend/handler.Patient":
		email_header = "Dear Patients, Thank you very much for using our service!"
	case "download-view/frontend/handler.Nurse":
		email_header = "Dear Nurse, Thank you very much for using our service!"
	case "download-view/frontend/handler.Doctor":
		email_header = "Dear Doctor, Thank you very much for using our service!"
	case "download-view/frontend/handler.Insurance":
		email_header = "Dear Insurance Practitioner, Thank you very much for using our service!"
	case "download-view/frontend/handler.Emergency":
		email_header = "Dear Doctor, We understand this is emergency and all historical records are attached for reference."
	case "download-view/frontend/handler.Researcher":
		email_header = "Dear Researcher, Here is the de-identified patient, vaccine & chestpain record for reference."
	default:
		email_header = "This is default email header"
	}

	email_sender := "owen.leung2@gmail.com"

	os.Setenv("AWS_ACCESS_KEY_ID", "XXXXXXX")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "XXXXXXXXX")
	env_config, _ := config.NewEnvConfig()
	awscredentials := env_config.Credentials
	aws_config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: awscredentials,
		}))
	if err != nil {
		log.Println("Creating Config encounters error")
		log.Println(err)
		os.Exit(1)
	}
	ses_client := ses.NewFromConfig(aws_config)
	msg := gomail.NewMessage()
	msg.SetHeader("From", email_sender)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", email_header)
	msg.SetBody("text/plain", "Here attached the file you requested!")

	for _, attn := range attachment {
		msg.Attach(attn)
	}

	var emailRaw bytes.Buffer
	msg.WriteTo(&emailRaw)

	raw_message := types.RawMessage{
		Data: []byte(emailRaw.Bytes()),
	}
	output, err := ses_client.SendRawEmail(context.TODO(),
		&ses.SendRawEmailInput{
			RawMessage:   &raw_message,
			Destinations: []string{recipient},
			Source:       &email_sender,
		},
	)
	log.Println("Output : ", output)
	log.Println("Error : ", err)

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
}

func UnitTest_AWSEmail(recipient string, attachment ...string) {
	SendEmailViaAWS(recipient, attachment...)
}
