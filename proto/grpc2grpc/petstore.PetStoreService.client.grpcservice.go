// This file registers with grpc service. This file was auto-generated by mashling at
// 2018-12-05 14:00:51.859472535 -0700 MST m=+0.048388335
package grpc2grpc

import (
	"context"

	"encoding/json"

	"github.com/project-flogo/grpc/support"

	"errors"

	"io"
	"strings"

	"log"

	"github.com/imdario/mergo"

	servInfo "github.com/project-flogo/grpc/activity"
	"google.golang.org/grpc"
)

type clientServicepetstorePetStoreServiceclient struct {
	serviceInfo *servInfo.ServiceInfo
}

var serviceInfopetstorePetStoreServiceclient = &servInfo.ServiceInfo{
	ProtoName:   "petstore",
	ServiceName: "PetStoreService",
}

func init() {
	servInfo.ClientServiceRegistery.RegisterClientService(&clientServicepetstorePetStoreServiceclient{serviceInfo: serviceInfopetstorePetStoreServiceclient})
}

//GetRegisteredClientService returns client implimentaion stub with grpc connection
func (cs *clientServicepetstorePetStoreServiceclient) GetRegisteredClientService(gCC *grpc.ClientConn) interface{} {
	return NewPetStoreServiceClient(gCC)
}

func (cs *clientServicepetstorePetStoreServiceclient) ServiceInfo() *servInfo.ServiceInfo {
	return cs.serviceInfo
}

func (cs *clientServicepetstorePetStoreServiceclient) InvokeMethod(reqArr map[string]interface{}) map[string]interface{} {

	clientObject := reqArr["ClientObject"].(PetStoreServiceClient)
	methodName := reqArr["MethodName"].(string)

	switch methodName {
	case "PetById":
		return PetById(clientObject, reqArr)
	case "UserByName":
		return UserByName(clientObject, reqArr)
	case "ListUsers":
		return ListUsers(clientObject, reqArr)
	case "StoreUsers":
		return StoreUsers(clientObject, reqArr)
	case "BulkUsers":
		return BulkUsers(clientObject, reqArr)
	}

	resMap := make(map[string]interface{}, 2)
	resMap["Response"] = []byte("null")
	resMap["Error"] = errors.New("Method not Available: " + methodName)
	return resMap
}
func PetById(client PetStoreServiceClient, values interface{}) map[string]interface{} {
	req := &PetByIdRequest{}
	support.AssignStructValues(req, values)
	res, err := client.PetById(context.Background(), req)
	b, errMarshl := json.Marshal(res)
	if errMarshl != nil {
		log.Println("Error: ", errMarshl)
		return nil
	}

	resMap := make(map[string]interface{}, 2)
	resMap["Response"] = b
	resMap["Error"] = err
	return resMap
}
func UserByName(client PetStoreServiceClient, values interface{}) map[string]interface{} {
	req := &UserByNameRequest{}
	support.AssignStructValues(req, values)
	res, err := client.UserByName(context.Background(), req)
	b, errMarshl := json.Marshal(res)
	if errMarshl != nil {
		log.Println("Error: ", errMarshl)
		return nil
	}

	resMap := make(map[string]interface{}, 2)
	resMap["Response"] = b
	resMap["Error"] = err
	return resMap
}

func ListUsers(client PetStoreServiceClient, reqArr map[string]interface{}) map[string]interface{} {
	resMap := make(map[string]interface{}, 1)

	if reqArr["Mode"] != nil {
		mode := reqArr["Mode"].(string)
		if strings.Compare(mode, "rest-to-grpc") == 0 {
			resMap["Error"] = errors.New("streaming operation is not allowed in rest to grpc case")
			return resMap
		}
	}

	req := &EmptyReq{}
	reqData := reqArr["reqdata"].(*EmptyReq)
	if err := mergo.Merge(req, reqData, mergo.WithOverride); err != nil {
		resMap["Error"] = errors.New("unable to merge reqData values")
		return resMap
	}

	sReq := reqArr["strmReq"].(PetStoreService_ListUsersServer)

	stream, err := client.ListUsers(context.Background(), req)
	if err != nil {
		log.Println("erorr while getting stream object for ListUsers:", err)
		resMap["Error"] = err
		return resMap
	}
	for {
		obj, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("erorr occured in ListUsers Recv():", err)
			resMap["Error"] = err
			return resMap
		}
		err = sReq.Send(obj)
		if err != nil {
			log.Println("error occured in ListUsers Send():", err)
			resMap["Error"] = err
			return resMap
		}
	}
	resMap["Error"] = nil
	return resMap
}

func StoreUsers(client PetStoreServiceClient, reqArr map[string]interface{}) map[string]interface{} {
	resMap := make(map[string]interface{}, 1)

	if reqArr["Mode"] != nil {
		mode := reqArr["Mode"].(string)
		if strings.Compare(mode, "rest-to-grpc") == 0 {
			resMap["Error"] = errors.New("streaming operation is not allowed in rest to grpc case")
			return resMap
		}
	}

	stream, err := client.StoreUsers(context.Background())
	if err != nil {
		log.Println("erorr while getting stream object for StoreUsers:", err)
		resMap["Error"] = err
		return resMap
	}

	cReq := reqArr["strmReq"].(PetStoreService_StoreUsersServer)

	for {
		dataObj, err := cReq.Recv()
		if err == io.EOF {
			obj, err := stream.CloseAndRecv()
			if err != nil {
				log.Println("erorr occured in StoreUsers CloseAndRecv():", err)
				resMap["Error"] = err
				return resMap
			}
			resMap["Error"] = cReq.SendAndClose(obj)
			return resMap
		}
		if err != nil {
			log.Println("error occured in StoreUsers Recv():", err)
			resMap["Error"] = err
			return resMap
		}

		if err := stream.Send(dataObj); err != nil {
			log.Println("error while sending dataObj with client stream:", err)
			resMap["Error"] = err
			return resMap
		}
	}
}

func BulkUsers(client PetStoreServiceClient, reqArr map[string]interface{}) map[string]interface{} {
	resMap := make(map[string]interface{}, 1)

	if reqArr["Mode"] != nil {
		mode := reqArr["Mode"].(string)
		if strings.Compare(mode, "rest-to-grpc") == 0 {
			resMap["Error"] = errors.New("streaming operation is not allowed in rest to grpc case")
			return resMap
		}
	}

	bReq := reqArr["strmReq"].(PetStoreService_BulkUsersServer)

	stream, err := client.BulkUsers(context.Background())
	if err != nil {
		log.Println("error while getting stream object for BulkUsers:", err)
		resMap["Error"] = err
		return resMap
	}

	waits := make(chan struct{})
	go func() {
		for {
			obj, err := bReq.Recv()
			if err == io.EOF {
				resMap["Error"] = nil
				stream.CloseSend()
				close(waits)
				return
			}
			if err != nil {
				log.Println("error occured in BulkUsers bidi Recv():", err)
				resMap["Error"] = err
				close(waits)
				return
			}
			if err := stream.Send(obj); err != nil {
				log.Println("error while sending obj with stream:", err)
				resMap["Error"] = err
				close(waits)
				return
			}
		}
	}()

	waitc := make(chan struct{})
	go func() {
		for {
			obj, err := stream.Recv()
			if err == io.EOF {
				resMap["Error"] = nil
				close(waitc)
				return
			}
			if err != nil {
				log.Println("erorr occured in BulkUsers stream Recv():", err)
				resMap["Error"] = err
				close(waitc)
				return
			}
			if sdErr := bReq.Send(obj); sdErr != nil {
				log.Println("error while sending obj with bidi Send():", sdErr)
				resMap["Error"] = sdErr
				close(waitc)
				return
			}
		}
	}()
	<-waitc
	<-waits
	return resMap
}
