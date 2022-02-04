package main

import {
  "fmt"
  "github.com/aws/aws-sdk-go"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
}


func requestHandler(req Request) (Response, error) {
    session := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    svc := dynamodb.New(session)

    type Item struct {
        StartDateTime   string
        EndDateTime     string
        Message         string
    }

    params := &dynamodb.ScanInput{
        TableName: process.env.tableName,
        Limit:100,
        FilterExpression: ":dateNow between StartDateTime and EndDateTime",
        ExpressionAttributeValues: {
          ":dateNow": Date.now()
        }
    }

    result, err := svc.Scan(params)
    if err != nil {
        log.Fatalf("Query API call failed: %s", err)
    }

    if result.Item == nil {
        return validator(true, "False")
    } else {
        item := Item{}
        err = dynamodbattribute.UnmarshalMap(result.Item, &item)
        if err != nil {
            panic(fmt.Sprintf("Failed to unmarshal record, %v", err))
        }
        return validator(true, "True", item.Message)
    }
}

func validator(isSuccess bool, holidayFound string, holidayMessage) (Response, error) {
    if isSuccess {
        return {
            holidayFound: holidayFound,
            holidayMessage: holidayMessage,
            lambdaResult: "Success"
        }
    } else {
        return { lambdaResult: "Error" }
    }
}

func main ()  {
    lambda.Start(requestHandler)
}
