package email

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"gopkg.in/gomail.v2"
)

func SendEmailViaAWS(recipient string, attachment ...string) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		log.Println("Unable to detect caller name when sending email")
		os.Exit(1)
	}
	caller := runtime.FuncForPC(pc)
	fmt.Println("Caller Name is ", caller.Name())
	var email_header string
	switch caller.Name() {
	case "ingest/update.Register_DownloadPatient":
		email_header = "Your Patient record in Blockchain has been updated"
	case "ingest/update.Register_DownloadAllergy":
		email_header = "Your Allergy record in Blockchain has been updated"
	case "ingest/update.Vaccine_DownloadVaccination":
		email_header = "Your Vaccination record in Blockchain has been updated"
	case "ingest/update.Diet_DownloadBodyMeasurement":
		email_header = "Your Body Measurement record in Blockchain has been updated"
	case "ingest/update.Diet_DownloadMealPlan":
		email_header = "Your Meal Plan record in Blockchain has been updated"
	case "ingest/update.AE_DownloadAssessment":
		email_header = "Your Assessment record in Blockchain has been updated"
	case "ingest/update.ChestPain_DownloadChestPain":
		email_header = "Your ChestPain record in Blockchain has been updated"
	case "ingest/update.PCI_DownloadRiskFactor":
		email_header = "Your Risk Factor record in Blockchain has been updated"
	case "ingest/update.PCI_DownloadProcedure":
		email_header = "Your Procedure record in Blockchain has been updated"
	default:
		email_header = "This is Default Email Header"
	}
	if len(attachment) > 1 {
		email_header = email_header + ". The Keys are also attached for your reference"
	}
	email_sender := "owen.leung2@gmail.com"

	os.Setenv("AWS_ACCESS_KEY_ID", "xxxxxxxx")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "xxxxxxxxx")
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
	msg.SetBody("text/plain", "Hello, here attach the key!")

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
