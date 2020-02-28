package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
import pb "github.com/hyperledger/fabric/protos/peer"

type recordInfo struct {
	Identity string   `json:"identity"`
	SickName string   `json:"sickName"`
	DrugName []string `json:"drugName"`
}

type resultData struct {
	RecordInfos []recordInfo `json:"recordInfos"`
}

func (r *recordInfo) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (r *recordInfo) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "save" {
		return r.saveRecord(stub, args)
	} else if funcName == "query" {
		return r.queryRecord(stub, args)
	} else if funcName == "queryHistory" {
		return r.queryHistoryRecord(stub, args)
	} else {
		return shim.Error("no such function")
	}
}

func (r *recordInfo) saveRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("except two args")
	} else {
		err := stub.PutState(args[0], []byte(args[1]))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	}

}

func (r *recordInfo) queryRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("except one arg")
	} else {
		value, err := stub.GetState(args[0])
		if err != nil {
			shim.Error("no data found")
		}
		return shim.Success(value)
	}
}

func (r *recordInfo) queryHistoryRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("except one arg")
	} else {
		//queryParam :="{\"selector\":{\"identity\":\""+id+"\"}}"
		recordinfos := []recordInfo{}
		it, err := stub.GetHistoryForKey(args[0])
		if err != nil {
			return shim.Error("no data found")
		} else {
			fmt.Println("is data exits?", it.HasNext())
			for it.HasNext() {
				keym, err := it.Next()
				if err != nil {
					return shim.Error("data get error")
				}
				value := keym.Value
				fmt.Println("get value is", string(value))
				recordinfo := new(recordInfo)
				json.Unmarshal(value, &recordinfo)
				fmt.Println("recordinfo is ", recordinfo)
				recordinfos = append(recordinfos, *recordinfo)
				fmt.Println("recordinfos is ", recordinfos)
			}
			resultdata := new(resultData)
			resultdata.RecordInfos = recordinfos
			fmt.Println("resultdata is ", resultdata)
			value, err := json.Marshal(resultdata)
			if err != nil {
				shim.Error(err.Error())
			}
			return shim.Success(value)
		}
	}
}

func main() {
	err := shim.Start(new(recordInfo))
	if err != nil {
		fmt.Println("emr recordInfo chaincode start error")
	}
}
