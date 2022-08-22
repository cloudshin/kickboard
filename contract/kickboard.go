package main

import (
	"encoding/json"
	"fmt"
	"encoding/base64"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Kickboard struct{
	KickID string `json:"kickid"`
	Bat string `json:"bat"`
	State string `json:"state"`
	History []List `json:"history"`
}
type List struct{
	UserID string  `json:"userid"`
	StartTime string `json:"starttime"`
	FinishTime string `json:"finishtime"`
	StartLocation string `json:"startlocation"`
	FinishLocation string `json:"Finishlocation"`
}
 
func (s *SmartContract) RegisterKickboard(ctx contractapi.TransactionContextInterface, KickID string, StartTime string, StartLocation string) error {

	// GetState()로 킥보드 등록 여부 확인

	var list = List{StartTime: StartTime, StartLocation: StartLocation}
	var kickboard = Kickboard{KickID: KickID, State: "1"}

	kickboard.History=append(kickboard.History,list)

	KickboardAsBytes, _ := json.Marshal(kickboard)	

	return ctx.GetStub().PutState(KickID, KickboardAsBytes)
}

func (s *SmartContract) QueryKickboard(ctx contractapi.TransactionContextInterface, KickID string) (string, error) {

	KickboardAsBytes, err := ctx.GetStub().GetState(KickID)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if KickboardAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", KickID)
	}
	
	return string(KickboardAsBytes[:]), nil	
}

func (s *SmartContract) DiscardKickboard(ctx contractapi.TransactionContextInterface, KickID string) error {
	
	// getState User 
	KickboardAsBytes, err := ctx.GetStub().GetState(KickID)	

	if err != nil{
		return err
	} else if KickboardAsBytes == nil{ // no State! error
		return fmt.Errorf("\"Error\":\"Kickboard does not exist: "+ KickID+"\"")
	}
	// state ok
	kickboard := Kickboard{}
	err = json.Unmarshal(KickboardAsBytes, &kickboard)
	if err != nil {
		return err
	}

	kickboard.State = "3"
	// update to User World state
	KickboardAsBytes, err = json.Marshal(kickboard);
	if err != nil {
		return fmt.Errorf("failed to Marshaling: %v", err)
	}	
	err = ctx.GetStub().PutState(KickID, KickboardAsBytes)
	if err != nil {
		return fmt.Errorf("failed to Discarding: %v", err)
	}	
	return nil
}

func (s *SmartContract) EnrollData(ctx contractapi.TransactionContextInterface, KickID string, StartTime string, StartLocation string, Bat string) error {
	
	// getState User 
	KickboardAsBytes, err := ctx.GetStub().GetState(KickID)	

	if err != nil{
		return err
	} else if KickboardAsBytes == nil{ // no State! error
		return fmt.Errorf("\"Error\":\"Kickboard does not exist: "+ KickID+"\"")
	}
	// state ok
	kickboard := Kickboard{}
	err = json.Unmarshal(KickboardAsBytes, &kickboard)
	if err != nil {
		return err
	}

	// create rate structure
	var List = List{StartTime: StartTime, StartLocation: StartLocation}

	kickboard.History=append(kickboard.History,List)

	kickboard.Bat = Bat
	// update to User World state
	KickboardAsBytes, err = json.Marshal(kickboard);
	if err != nil {
		return fmt.Errorf("failed to Marshaling: %v", err)
	}	
	err = ctx.GetStub().PutState(KickID, KickboardAsBytes)
	if err != nil {
		return fmt.Errorf("failed to Enrolling Data: %v", err)
	}	
	return nil
}

func (s *SmartContract) UseKickboard(ctx contractapi.TransactionContextInterface, KickID string, StartTime string, StartLocation string) error {
	
	// getState User 
	KickboardAsBytes, err := ctx.GetStub().GetState(KickID)	
	
	if err != nil{
		return err
	} else if KickboardAsBytes == nil{ // no State! error
		return fmt.Errorf("\"Error\":\"Kickboard does not exist: "+ KickID+"\"")
	}
	// state ok
	kickboard := Kickboard{}
	err = json.Unmarshal(KickboardAsBytes, &kickboard)
	if err != nil {
		return err
	}

	if kickboard.State == "2" {
		return fmt.Errorf("\"Error\":\"Kickboard 사용중: "+ KickID+"\"")
	} else if kickboard.State == "3" {
		return fmt.Errorf("\"Error\":\"Kickboard 폐기: "+ KickID+"\"")
	}

	UserID, err := submittingClientIdentity(ctx)
	if err != nil {
		return err
	}


	// create rate structure
	var List = List{UserID: UserID, StartTime: StartTime, StartLocation: StartLocation}

	kickboard.History=append(kickboard.History,List)

	kickboard.State = "2"
	// update to User World state
	KickboardAsBytes, err = json.Marshal(kickboard);
	if err != nil {
		return fmt.Errorf("failed to Marshaling: %v", err)
	}	
	err = ctx.GetStub().PutState(KickID, KickboardAsBytes)
	if err != nil {
		return fmt.Errorf("failed to using: %v", err)
	}	
	return nil
}

func (s *SmartContract) FinishKickboard(ctx contractapi.TransactionContextInterface, KickID string, FinishTime string, FinishLocation string) error {
	
	// getState User 
	KickboardAsBytes, err := ctx.GetStub().GetState(KickID)	
	
	if err != nil{
		return err
	} else if KickboardAsBytes == nil{ // no State! error
		return fmt.Errorf("\"Error\":\"Kickboard does not exist: "+ KickID+"\"")
	}
	// state ok
	kickboard := Kickboard{}
	err = json.Unmarshal(KickboardAsBytes, &kickboard)
	if err != nil {
		return err
	}

	if kickboard.State == "1" {
		return fmt.Errorf("\"Error\":\"Kickboard 대기중: "+ KickID+"\"")
	} else if kickboard.State == "3" {
		return fmt.Errorf("\"Error\":\"Kickboard 폐기: "+ KickID+"\"")
	}

	UserID, err := submittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	// create rate structure
	var List = List{UserID: UserID, FinishTime: FinishTime, FinishLocation: FinishLocation}

	kickboard.History=append(kickboard.History,List)

	kickboard.State = "1"
	// update to User World state
	KickboardAsBytes, err = json.Marshal(kickboard);
	if err != nil {
		return fmt.Errorf("failed to Marshaling: %v", err)
	}	
	err = ctx.GetStub().PutState(KickID, KickboardAsBytes)
	if err != nil {
		return fmt.Errorf("failed to Finish: %v", err)
	}	
	return nil
}

func submittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}

	DecodeID := string(decodeID)
	user := strings.Split(DecodeID, ",")

	getuser := user[0]
	UserID := getuser[9:]

	return UserID, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create teamate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting teamate chaincode: %s", err.Error())
	}
}
