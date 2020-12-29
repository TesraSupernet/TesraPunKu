package ztools

import (
	uuid "github.com/satori/go.uuid"
)

func CreatUuidv4() string {
	// or error handling
	//u4 := uuid.Must(uuid.NewV4())

	u4 := uuid.NewV4()
	//if err != nil {
	//	return "", err
	//}
	return u4.String()

}

func CreatCancelLineUuid() string {
	var id string
	uuid4 := CreatUuidv4()
	id = CancelTargetStr(uuid4, "-")
	return id
}
