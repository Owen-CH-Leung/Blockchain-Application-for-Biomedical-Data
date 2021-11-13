package email

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"upload-view/backend/logger"

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
	fmt.Println("Caller Name is ", caller.Name())
	var email_header string
	switch caller.Name() {
	case "upload-view/frontend/download.Register_DownloadPatient":
		email_header = "Here is your Patient record in Blockchain"
	case "upload-view/frontend/download.Register_DownloadAllergy":
		email_header = "Here is your Allergy record in Blockchain"
	case "upload-view/frontend/download.Vaccine_DownloadVaccination":
		email_header = "Here is your Vaccination record in Blockchain"
	case "upload-view/frontend/download.Diet_DownloadBodyMeasurement":
		email_header = "Here is your Body Measurement record in Blockchain"
	case "upload-view/frontend/download.Diet_DownloadMealPlan":
		email_header = "Here is your Meal Plan record in Blockchain"
	case "upload-view/frontend/download.AE_DownloadAssessment":
		email_header = "Here is your Assessment record in Blockchain"
	case "upload-view/frontend/download.ChestPain_DownloadChestPain":
		email_header = "Here is your ChestPain record in Blockchain"
	case "upload-view/frontend/download.PCI_DownloadRiskFactor":
		email_header = "Here is your Risk Factor record in Blockchain"
	case "upload-view/frontend/download.PCI_DownloadProcedure":
		email_header = "Here is your Procedure record in Blockchain"
	default:
		email_header = "This is Default Email Header"
	}
	if len(attachment) > 1 {
		email_header = email_header + ". The Keys are also attached for your reference"
	}
	email_sender := "owen.leung2@gmail.com"

	os.Setenv("AWS_ACCESS_KEY_ID", "XXXXXX")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "XXXXXX")
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
	msg.SetBody("text/plain", "Hello, here attach the key / files for your reference!!")

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
