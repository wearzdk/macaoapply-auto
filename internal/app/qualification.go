package app

import (
	"encoding/json"
	"fmt"
	"log"
	"macaoapply-auto/internal/client"

	"github.com/golang-jwt/jwt/v5"
)

type FormInstance struct {
	CreateUserID                string `json:"createUserId"`
	CreateUserName              string `json:"createUserName"`
	CreateTimestamp             int    `json:"createTimestamp"`
	UpdateUserID                string `json:"updateUserId"`
	UpdateUserName              string `json:"updateUserName"`
	UpdateTimestamp             int    `json:"updateTimestamp"`
	ID                          string `json:"id"`
	FormInstanceID              string `json:"formInstanceId"`
	FormID                      string `json:"formId"`
	DataID                      string `json:"dataId"`
	FormInstCode                string `json:"formInstCode"`
	ProcessInstanceID           string `json:"processInstanceId"`
	ApplyTime                   int    `json:"applyTime"`
	ReceiveTime                 int    `json:"receiveTime"`
	ApplyBy                     string `json:"applyBy"`
	ApplyByName                 string `json:"applyByName"`
	Type                        string `json:"type"`
	AccountType                 string `json:"accountType"`
	Status                      string `json:"status"`
	PlateNumber                 string `json:"plateNumber"`
	Temporaryplate              string `json:"temporaryplate"`
	ValidityPeriod              any    `json:"validityPeriod"`
	QualificationNumber         string `json:"qualificationNumber"`
	CarOwnerID                  string `json:"carOwnerId"`
	FlowType                    string `json:"flowType"`
	FlowTypeStatus              string `json:"flowTypeStatus"`
	TaskDefinitionKey           string `json:"taskDefinitionKey"`
	CompanyCode                 string `json:"companyCode"`
	MacaoCustomsCheck           string `json:"macaoCustomsCheck"`
	MacaoCustomsHandleDate      any    `json:"macaoCustomsHandleDate"`
	IsOwnDriverCard             string `json:"isOwnDriverCard"`
	DriverCancelState           string `json:"driverCancelState"`
	IsPrintFailure              string `json:"isPrintFailure"`
	CmcImage                    string `json:"cmcImage"`
	IsOldData                   string `json:"isOldData"`
	IsTakeTemporary             string `json:"isTakeTemporary"`
	TakeTemporaryHandleDate     any    `json:"takeTemporaryHandleDate"`
	CorrectState                string `json:"correctState"`
	TakeDriverBookStatus        string `json:"takeDriverBookStatus"`
	PrintDriverBookDate         any    `json:"printDriverBookDate"`
	IssueDriverBookDate         any    `json:"issueDriverBookDate"`
	PayTemporaryStatus          string `json:"payTemporaryStatus"`
	PayTemporaryDate            any    `json:"payTemporaryDate"`
	MacaoCustomsHandleByName    string `json:"macaoCustomsHandleByName"`
	TakeTemporaryHandleByName   string `json:"takeTemporaryHandleByName"`
	IsReturnCarInspection       string `json:"isReturnCarInspection"`
	DownloadMacaoDataTime       int    `json:"downloadMacaoDataTime"`
	CiqDealByName               string `json:"ciqDealByName"`
	BorderDealByName            string `json:"borderDealByName"`
	PrintTemporaryReason        string `json:"printTemporaryReason"`
	IsVehicleValueChange        string `json:"isVehicleValueChange"`
	PoliceAuditState            string `json:"policeAuditState"`
	CmcAuditState               string `json:"cmcAuditState"`
	RemoveCardState             string `json:"removeCardState"`
	RemoveCardDate              any    `json:"removeCardDate"`
	VehicleCancelReason         string `json:"vehicleCancelReason"`
	VehicleCancelApplyType      string `json:"vehicleCancelApplyType"`
	VehicleCancelApplyDate      any    `json:"vehicleCancelApplyDate"`
	VehicleCancelDriverName     string `json:"vehicleCancelDriverName"`
	VehicleCancelDriverIDCardNo string `json:"vehicleCancelDriverIdCardNo"`
	ChangeDriverInfoRemark      string `json:"changeDriverInfoRemark"`
	ReissueVehicleReason        string `json:"reissueVehicleReason"`
	ReissueDriverType           string `json:"reissueDriverType"`
	ReissueDriverReason         string `json:"reissueDriverReason"`
	ReissueDriverLocation       string `json:"reissueDriverLocation"`
	InsuranceState              string `json:"insuranceState"`
	QualificationID             string `json:"qualificationId"`
	MsgState                    string `json:"msgState"`
	ReceiveState                string `json:"receiveState"`
	IsEditCarApproveDate        string `json:"isEditCarApproveDate"`
	Pk                          string `json:"pk"`
	TypeName                    any    `json:"typeName"`
	AccountTypeName             any    `json:"accountTypeName"`
	StatusName                  any    `json:"statusName"`
	TaskID                      any    `json:"taskId"`
	FlowTypeName                string `json:"flowTypeName"`
	VehicleID                   any    `json:"vehicleId"`
}
type FormInstanceList struct {
	FormInstance FormInstance `json:"formInstance"`
}
type GetPassQualificationResp struct {
	PassAppointmentAdvance int                `json:"passAppointmentAdvance"`
	FormInstanceList       []FormInstanceList `json:"formInstanceList"`
}

func getPassQualification(plateNumber string) (FormInstance, error) {
	resp, err := client.RequestWithCache("POST", "before/sys/appointment/getPassQualification", jwt.MapClaims{})
	if err != nil {
		log.Println("获取预约资格失败：" + err.Error())
		return FormInstance{}, err
	}
	var getPassQualificationResp GetPassQualificationResp
	err = json.Unmarshal([]byte(resp), &getPassQualificationResp)
	if err != nil {
		log.Println("解析预约资格失败：" + err.Error())
		return FormInstance{}, err
	}
	for _, formInstanceList := range getPassQualificationResp.FormInstanceList {
		if formInstanceList.FormInstance.PlateNumber == plateNumber {
			return formInstanceList.FormInstance, nil
		}
	}
	return FormInstance{}, fmt.Errorf("未找到车牌号为 %s 的预约资格", plateNumber)
}
