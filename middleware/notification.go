package middleware


import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)
 // struct  yapısı olarak 
 // notification init 

 //notif token

// eEi1Yy4sTe2ihqZ2DIJTAT:APA91bHSFXqsCTtvE4z6Ea-Eleity9hpAEGBSy3i6HWpl4EqcDLfcqtoZHgrmvgJUkRgtLd_C_EGn4GMKOVVLUkhLvue4SkhiEZ__x8zhcRY8U2sQP9wWaN-CAAVuTTgXgiEeRak2EJ2

func getDecodedFireBaseKey() ([]byte, error) {
	// Here you can implement the logic to retrieve and decode your Firebase key
	// For simplicity, we assume you load it from an environment variable
	key := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	data, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SendPushNotification(deviceTokens []string) error {
	decodedKey, err := getDecodedFireBaseKey()
	if err != nil {
		return err
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		log.Printf("Error initializing Firebase: %s", err)
		return err
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	response, err := fcmClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Congratulations!!",
			Body:  "You have just implemented push notification",
		},
		Tokens: deviceTokens,
	})

	if err != nil {
		return err
	}

	log.Printf("Response success count: %d", response.SuccessCount)
	log.Printf("Response failure count: %d", response.FailureCount)

	return nil
}
