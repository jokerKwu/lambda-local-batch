package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

type User struct {
	Id             int       `json:"id" bson:"id" validate:"required"`
	Email          string    `json:"email" bson:"email" validate:"required"`
	LastSignin     time.Time `json:"lastSignin" bson:"lastSignin" validate:"required"`
	LastTokenIssue time.Time `json:"lastTokenIssue" bson:"lastTokenIssue" validate:"required"`
	Name           string    `json:"name" bson:"name" validate:"omitempty"`
}

type Users []User

func userPrint(from string, to string, users Users) {
	log.Println("================================")
	log.Printf("%s - %s 휴면 계좌 조회 \n", from, to)
	for i := 0; i < len(users); i++ {
		log.Println("================================")
		log.Printf("유저 이름 : %s\n", users[i].Name)
		log.Printf("유저 이메일 : %s\n", users[i].Email)
		log.Printf("마지막 토큰 발행 날짜 : %s\n", users[i].LastTokenIssue.String())
	}
	log.Println("================================")
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//휴면 날짜 계산
	warningDate := time.Now().AddDate(0, 0, -335)
	warningDateFrom := time.Date(warningDate.Year(), warningDate.Month(), warningDate.Day(), 0, 0, 0, 0, time.UTC)
	warningDate = time.Now().AddDate(0, 0, -334)
	warningDateTo := time.Date(warningDate.Year(), warningDate.Month(), warningDate.Day(), 0, 0, 0, 0, time.UTC)
	//몽고 디비 연결
	client, err := GetClient()
	//디비 휴면계좌 조회
	collection := client.Database("webboard").Collection("user")
	cur, err := collection.Find(context.Background(), bson.M{"lastTokenIssue": bson.M{"$gte": warningDateFrom, "$lte": warningDateTo}})

	if err != nil {
		log.Println("Error on Finding all the documents", err)
	}
	var users Users
	for cur.Next(context.Background()) {
		var user User
		err = cur.Decode(&user)
		if err != nil {
			log.Println("Error on Decoding the document", err)
		}
		users = append(users, user)
	}
	userPrint(warningDateFrom.String(), warningDateTo.String(), users)
	defer client.Disconnect(context.Background())
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("휴면 계좌, %d", len(users)),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
