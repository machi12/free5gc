package producer

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/free5gc/amf/internal/context"
	gmm_common "github.com/free5gc/amf/internal/gmm/common"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/internal/nas/nas_security"
	"github.com/free5gc/amf/internal/sbi/consumer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/machi12/nas/security"
	"github.com/machi12/openapi/models"
)

// TS 29.518 5.2.2.2.3
func HandleCreateUEContextRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.CommLog.Infof("Handle Create UE Context Request")

	createUeContextRequest := request.Body.(models.CreateUeContextRequest)
	ueContextID := request.Params["ueContextId"]

	createUeContextResponse, ueContextCreateError := CreateUEContextProcedure(ueContextID, createUeContextRequest)
	if ueContextCreateError != nil {
		return httpwrapper.NewResponse(int(ueContextCreateError.Error.Status), nil, ueContextCreateError)
	} else {
		return httpwrapper.NewResponse(http.StatusCreated, nil, createUeContextResponse)
	}
}

func CreateUEContextProcedure(ueContextID string, createUeContextRequest models.CreateUeContextRequest) (
	*models.CreateUeContextResponse, *models.UeContextCreateError,
) {
	amfSelf := context.GetSelf()
	ueContextCreateData := createUeContextRequest.JsonData

	if ueContextCreateData.UeContext == nil || ueContextCreateData.TargetId == nil ||
		ueContextCreateData.PduSessionList == nil || ueContextCreateData.SourceToTargetData == nil ||
		ueContextCreateData.N2NotifyUri == "" {
		ueContextCreateError := &models.UeContextCreateError{
			Error: &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "HANDOVER_FAILURE",
			},
		}
		return nil, ueContextCreateError
	}
	// create the UE context in target amf
	ue := amfSelf.NewAmfUe(ueContextID)
	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	// amfSelf.AmfRanSetByRanId(*ueContextCreateData.TargetId.RanNodeId)
	// ue.N1N2Message[ueContextId] = &context.N1N2Message{}
	// ue.N1N2Message[ueContextId].Request.JsonData = &models.N1N2MessageTransferReqData{
	// 	N2InfoContainer: &models.N2InfoContainer{
	// 		SmInfo: &models.N2SmInformation{
	// 			N2InfoContent: ueContextCreateData.SourceToTargetData,
	// 		},
	// 	},
	// }
	ue.HandoverNotifyUri = ueContextCreateData.N2NotifyUri

	amfSelf.AmfRanFindByRanID(*ueContextCreateData.TargetId.RanNodeId)
	supportedTAI := context.NewSupportedTAI()
	supportedTAI.Tai.Tac = ueContextCreateData.TargetId.Tai.Tac
	supportedTAI.Tai.PlmnId = ueContextCreateData.TargetId.Tai.PlmnId
	// ue.N1N2MessageSubscribeInfo[ueContextID] = &models.UeN1N2InfoSubscriptionCreateData{
	// 	N2NotifyCallbackUri: ueContextCreateData.N2NotifyUri,
	// }
	ue.UnauthenticatedSupi = ueContextCreateData.UeContext.SupiUnauthInd
	// should be smInfo list

	//for _, smInfo := range ueContextCreateData.PduSessionList {
	//if smInfo.N2InfoContent.NgapIeType == "NgapIeType_HANDOVER_REQUIRED" {
	// ue.N1N2Message[amfSelf.Uri].Request.JsonData.N2InfoContainer.SmInfo = &smInfo
	//}
	//}

	ue.RoutingIndicator = ueContextCreateData.UeContext.RoutingIndicator

	// optional
	ue.UdmGroupId = ueContextCreateData.UeContext.UdmGroupId
	ue.AusfGroupId = ueContextCreateData.UeContext.AusfGroupId
	// ueContextCreateData.UeContext.HpcfId
	ue.RatType = ueContextCreateData.UeContext.RestrictedRatList[0] // minItem = -1
	// ueContextCreateData.UeContext.ForbiddenAreaList
	// ueContextCreateData.UeContext.ServiceAreaRestriction
	// ueContextCreateData.UeContext.RestrictedCoreNwTypeList

	// it's not in 5.2.2.1.1 step 2a, so don't support
	// ue.Gpsi = ueContextCreateData.UeContext.GpsiList
	// ue.Pei = ueContextCreateData.UeContext.Pei
	// ueContextCreateData.UeContext.GroupList
	// ueContextCreateData.UeContext.DrxParameter
	// ueContextCreateData.UeContext.SubRfsp
	// ueContextCreateData.UeContext.UsedRfsp
	// ue.UEAMBR = ueContextCreateData.UeContext.SubUeAmbr
	// ueContextCreateData.UeContext.SmsSupport
	// ueContextCreateData.UeContext.SmsfId
	// ueContextCreateData.UeContext.SeafData
	// ueContextCreateData.UeContext.Var5gMmCapability
	// ueContextCreateData.UeContext.PcfId
	// ueContextCreateData.UeContext.PcfAmPolicyUri
	// ueContextCreateData.UeContext.AmPolicyReqTriggerList
	// ueContextCreateData.UeContext.EventSubscriptionList
	// ueContextCreateData.UeContext.MmContextList
	// ue.CurPduSession.PduSessionId = ueContextCreateData.UeContext.SessionContextList.
	// ue.TraceData = ueContextCreateData.UeContext.TraceData
	createUeContextResponse := new(models.CreateUeContextResponse)
	createUeContextResponse.JsonData = &models.UeContextCreatedData{
		UeContext: &models.UeContext{
			Supi: ueContextCreateData.UeContext.Supi,
		},
	}

	// response.JsonData.TargetToSourceData =
	// ue.N1N2Message[ueContextId].Request.JsonData.N2InfoContainer.SmInfo.N2InfoContent
	createUeContextResponse.JsonData.PduSessionList = ueContextCreateData.PduSessionList
	createUeContextResponse.JsonData.PcfReselectedInd = false
	// TODO: When  Target AMF selects a nw PCF for AM policy, set the flag to true.

	//response.UeContext = ueContextCreateData.UeContext
	//response.TargetToSourceData = ue.N1N2Message[amfSelf.Uri].Request.JsonData.N2InfoContainer.SmInfo.N2InfoContent
	//response.PduSessionList = ueContextCreateData.PduSessionList
	//response.PcfReselectedInd = false // TODO:When  Target AMF selects a nw PCF for AM policy, set the flag to true.
	//

	// return httpwrapper.NewResponse(http.StatusCreated, nil, createUeContextResponse)
	return createUeContextResponse, nil
}

// TS 29.518 5.2.2.2.4
func HandleReleaseUEContextRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.CommLog.Info("Handle Release UE Context Request")

	ueContextRelease := request.Body.(models.UeContextRelease)
	ueContextID := request.Params["ueContextId"]

	problemDetails := ReleaseUEContextProcedure(ueContextID, ueContextRelease)
	if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func ReleaseUEContextProcedure(ueContextID string, ueContextRelease models.UeContextRelease) *models.ProblemDetails {
	amfSelf := context.GetSelf()

	// TODO: UE is emergency registered and the SUPI is not authenticated
	if ueContextRelease.Supi != "" {
		logger.GmmLog.Warnf("Emergency registered UE is not supported.")
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "UNSPECIFIED",
		}
		return problemDetails
	}

	if ueContextRelease.NgapCause == nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
		}
		return problemDetails
	}

	logger.CommLog.Debugf("Release UE Context NGAP cause: %+v", ueContextRelease.NgapCause)

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	// TODO: TS 23.502 4.11.1.2.3.4
	// If the target CN node is AMF, the AMF invokes the
	// "Nsmf_PDUSession_UpdateSMContext request (SUPI, Relocation Cancel Indication) toward the SMF. Based
	// on the Relocation Cancel Indication, the target CN node deletes the session resources established during
	// handover preparation phase in SMF and UPF.

	gmm_common.RemoveAmfUe(ue, false)

	return nil
}

func HandleMobiRegUe(ue *context.AmfUe, ueContextTransferRspData *models.UeContextTransferRspData,
	ueContextTransferResponse *models.UeContextTransferResponse,
) {
	ueContextTransferRspData.UeRadioCapability = &models.N2InfoContent{
		NgapMessageType: 0,
		NgapIeType:      models.NgapIeType_UE_RADIO_CAPABILITY,
		NgapData: &models.RefToBinaryData{
			ContentId: "n2Info",
		},
	}
	b := []byte(ue.UeRadioCapability)
	copy(ueContextTransferResponse.BinaryDataN2Information, b)
}

// TS 29.518 5.2.2.2.1
func HandleUEContextTransferRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.CommLog.Info("Handle UE Context Transfer Request")

	ueContextTransferRequest := request.Body.(models.UeContextTransferRequest)
	ueContextID := request.Params["ueContextId"]

	ueContextTransferResponse, problemDetails := UEContextTransferProcedure(ueContextID, ueContextTransferRequest)
	if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return httpwrapper.NewResponse(http.StatusOK, nil, ueContextTransferResponse)
	}
}

func UEContextTransferProcedure(ueContextID string, ueContextTransferRequest models.UeContextTransferRequest) (
	*models.UeContextTransferResponse, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	if ueContextTransferRequest.JsonData == nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
		}
		return nil, problemDetails
	}

	UeContextTransferReqData := ueContextTransferRequest.JsonData

	if UeContextTransferReqData.AccessType == "" || UeContextTransferReqData.Reason == "" {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
		}
		return nil, problemDetails
	}

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	ueContextTransferResponse := &models.UeContextTransferResponse{
		JsonData: new(models.UeContextTransferRspData),
	}
	ueContextTransferRspData := ueContextTransferResponse.JsonData

	//if ue.GetAnType() != UeContextTransferReqData.AccessType {
	//for _, tai := range ue.RegistrationArea[ue.GetAnType()] {
	//if UeContextTransferReqData.PlmnId == tai.PlmnId {
	// TODO : generate N2 signalling
	//}
	//}
	//}

	switch UeContextTransferReqData.Reason {
	case models.TransferReason_INIT_REG:
		_, integrityProtected, err := nas_security.Decode(ue, UeContextTransferReqData.AccessType,
			ueContextTransferRequest.BinaryDataN1Message, true)
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INTEGRITY_CHECK_FAIL",
			}
			ue.NASLog.Errorln(err)
			return nil, problemDetails
		}
		if integrityProtected {
			ueContextTransferRspData.UeContext = buildUEContextModel(ue, UeContextTransferReqData.Reason)
		} else {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INTEGRITY_CHECK_FAIL",
			}
			return nil, problemDetails
		}
		// TODO: handle condition of TS 29.518 5.2.2.2.1.1 step 2a case b
	case models.TransferReason_MOBI_REG:
		_, integrityProtected, err := nas_security.Decode(ue, UeContextTransferReqData.AccessType,
			ueContextTransferRequest.BinaryDataN1Message, false)
		if err != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INTEGRITY_CHECK_FAIL",
			}
			ue.NASLog.Errorln(err)
			return nil, problemDetails
		}
		if integrityProtected {
			ueContextTransferRspData.UeContext = buildUEContextModel(ue, UeContextTransferReqData.Reason)
		} else {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INTEGRITY_CHECK_FAIL",
			}
			return nil, problemDetails
		}
		HandleMobiRegUe(ue, ueContextTransferRspData, ueContextTransferResponse)

	case models.TransferReason_MOBI_REG_UE_VALIDATED:
		ueContextTransferRspData.UeContext = buildUEContextModel(ue, UeContextTransferReqData.Reason)
		HandleMobiRegUe(ue, ueContextTransferRspData, ueContextTransferResponse)

	default:
		logger.ProducerLog.Warnf("Invalid Transfer Reason: %+v", UeContextTransferReqData.Reason)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "MANDATORY_IE_INCORRECT",
			InvalidParams: []models.InvalidParam{
				{
					Param: "reason",
				},
			},
		}
		return nil, problemDetails
	}
	return ueContextTransferResponse, nil
}

func buildUEContextModel(ue *context.AmfUe, Reason models.TransferReason) *models.UeContext {
	ueContext := new(models.UeContext)
	ueContext.Supi = ue.Supi
	ueContext.SupiUnauthInd = ue.UnauthenticatedSupi
	if Reason == models.TransferReason_INIT_REG || Reason == models.TransferReason_MOBI_REG {
		var mmContext models.MmContext
		mmContext.AccessType = models.AccessType__3_GPP_ACCESS
		NasSecurityMode := new(models.NasSecurityMode)
		switch ue.IntegrityAlg {
		case security.AlgIntegrity128NIA0:
			NasSecurityMode.IntegrityAlgorithm = models.IntegrityAlgorithm_NIA0
		case security.AlgIntegrity128NIA1:
			NasSecurityMode.IntegrityAlgorithm = models.IntegrityAlgorithm_NIA1
		case security.AlgIntegrity128NIA2:
			NasSecurityMode.IntegrityAlgorithm = models.IntegrityAlgorithm_NIA2
		case security.AlgIntegrity128NIA3:
			NasSecurityMode.IntegrityAlgorithm = models.IntegrityAlgorithm_NIA3
		}
		switch ue.CipheringAlg {
		case security.AlgCiphering128NEA0:
			NasSecurityMode.CipheringAlgorithm = models.CipheringAlgorithm_NEA0
		case security.AlgCiphering128NEA1:
			NasSecurityMode.CipheringAlgorithm = models.CipheringAlgorithm_NEA1
		case security.AlgCiphering128NEA2:
			NasSecurityMode.CipheringAlgorithm = models.CipheringAlgorithm_NEA2
		case security.AlgCiphering128NEA3:
			NasSecurityMode.CipheringAlgorithm = models.CipheringAlgorithm_NEA3
		}
		NgKsi := new(models.NgKsi)
		NgKsi.Ksi = ue.NgKsi.Ksi
		NgKsi.Tsc = ue.NgKsi.Tsc
		KeyAmf := new(models.KeyAmf)
		KeyAmf.KeyType = models.KeyAmfType_KAMF
		KeyAmf.KeyVal = ue.Kamf
		SeafData := new(models.SeafData)
		SeafData.NgKsi = NgKsi
		SeafData.KeyAmf = KeyAmf
		if ue.NH != nil {
			SeafData.Nh = hex.EncodeToString(ue.NH)
		}
		SeafData.Ncc = int32(ue.NCC)
		SeafData.KeyAmfChangeInd = false
		SeafData.KeyAmfHDerivationInd = false
		ueContext.SeafData = SeafData
		mmContext.NasSecurityMode = NasSecurityMode
		if ue.UESecurityCapability.Buffer != nil {
			mmContext.UeSecurityCapability = base64.StdEncoding.EncodeToString(ue.UESecurityCapability.Buffer)
		}
		mmContext.NasDownlinkCount = int32(ue.DLCount.Get())
		mmContext.NasUplinkCount = int32(ue.ULCount.Get())
		if ue.AllowedNssai[models.AccessType__3_GPP_ACCESS] != nil {
			for _, allowedSnssai := range ue.AllowedNssai[models.AccessType__3_GPP_ACCESS] {
				mmContext.AllowedNssai = append(mmContext.AllowedNssai, *(allowedSnssai.AllowedSnssai))
			}
		}
		ueContext.MmContextList = append(ueContext.MmContextList, mmContext)
	}
	if Reason == models.TransferReason_MOBI_REG_UE_VALIDATED || Reason == models.TransferReason_MOBI_REG {
		sessionContextList := &ueContext.SessionContextList
		ue.SmContextList.Range(func(key, value interface{}) bool {
			smContext := value.(*context.SmContext)
			snssai := smContext.Snssai()
			pduSessionContext := models.PduSessionContext{
				PduSessionId: smContext.PduSessionID(),
				SmContextRef: smContext.SmContextRef(),
				SNssai:       &snssai,
				Dnn:          smContext.Dnn(),
				AccessType:   smContext.AccessType(),
				HsmfId:       smContext.HSmfID(),
				VsmfId:       smContext.VSmfID(),
				NsInstance:   smContext.NsInstance(),
			}
			*sessionContextList = append(*sessionContextList, pduSessionContext)
			return true
		})
	}
	if ue.Gpsi != "" {
		ueContext.GpsiList = append(ueContext.GpsiList, ue.Gpsi)
	}

	if ue.Pei != "" {
		ueContext.Pei = ue.Pei
	}

	if ue.UdmGroupId != "" {
		ueContext.UdmGroupId = ue.UdmGroupId
	}

	if ue.AusfGroupId != "" {
		ueContext.AusfGroupId = ue.AusfGroupId
	}

	if ue.RoutingIndicator != "" {
		ueContext.RoutingIndicator = ue.RoutingIndicator
	}

	if ue.AccessAndMobilitySubscriptionData != nil {
		if ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr != nil {
			ueContext.SubUeAmbr = &models.Ambr{
				Uplink:   ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Uplink,
				Downlink: ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Downlink,
			}
		}
		if ue.AccessAndMobilitySubscriptionData.RfspIndex != 0 {
			ueContext.SubRfsp = ue.AccessAndMobilitySubscriptionData.RfspIndex
		}
	}

	if ue.PcfId != "" {
		ueContext.PcfId = ue.PcfId
	}

	if ue.AmPolicyUri != "" {
		ueContext.PcfAmPolicyUri = ue.AmPolicyUri
	}

	if ue.AmPolicyAssociation != nil {
		if len(ue.AmPolicyAssociation.Triggers) > 0 {
			ueContext.AmPolicyReqTriggerList = buildAmPolicyReqTriggers(ue.AmPolicyAssociation.Triggers)
		}
	}

	for _, eventSub := range ue.EventSubscriptionsInfo {
		if eventSub.EventSubscription != nil {
			ueContext.EventSubscriptionList = append(ueContext.EventSubscriptionList, *eventSub.EventSubscription)
		}
	}

	if ue.TraceData != nil {
		ueContext.TraceData = ue.TraceData
	}
	return ueContext
}

func buildAmPolicyReqTriggers(triggers []models.RequestTrigger) (amPolicyReqTriggers []models.AmPolicyReqTrigger) {
	for _, trigger := range triggers {
		switch trigger {
		case models.RequestTrigger_LOC_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_LOCATION_CHANGE)
		case models.RequestTrigger_PRA_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_PRA_CHANGE)
		case models.RequestTrigger_SERV_AREA_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_SARI_CHANGE)
		case models.RequestTrigger_RFSP_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_RFSP_INDEX_CHANGE)
		}
	}
	return
}

// TS 29.518 5.2.2.6
func HandleAssignEbiDataRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.CommLog.Info("Handle Assign Ebi Data Request")

	assignEbiData := request.Body.(models.AssignEbiData)
	ueContextID := request.Params["ueContextId"]

	assignedEbiData, assignEbiError, problemDetails := AssignEbiDataProcedure(ueContextID, assignEbiData)
	if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else if assignEbiError != nil {
		return httpwrapper.NewResponse(int(assignEbiError.Error.Status), nil, assignEbiError)
	} else {
		return httpwrapper.NewResponse(http.StatusOK, nil, assignedEbiData)
	}
}

func AssignEbiDataProcedure(ueContextID string, assignEbiData models.AssignEbiData) (
	*models.AssignedEbiData, *models.AssignEbiError, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, nil, problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	// TODO: AssignEbiError not used, check it!
	if _, ok := ue.SmContextFindByPDUSessionID(assignEbiData.PduSessionId); ok {
		var assignedEbiData *models.AssignedEbiData
		assignedEbiData.PduSessionId = assignEbiData.PduSessionId
		return assignedEbiData, nil, nil
	}
	logger.ProducerLog.Errorf("SmContext[PDU Session ID:%d] not found", assignEbiData.PduSessionId)
	return nil, nil, nil
}

// TS 29.518 5.2.2.2.2
func HandleRegistrationStatusUpdateRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.CommLog.Info("Handle Registration Status Update Request")

	ueRegStatusUpdateReqData := request.Body.(models.UeRegStatusUpdateReqData)
	ueContextID := request.Params["ueContextId"]

	ueRegStatusUpdateRspData, problemDetails := RegistrationStatusUpdateProcedure(ueContextID, ueRegStatusUpdateReqData)
	if problemDetails != nil {
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return httpwrapper.NewResponse(http.StatusOK, nil, ueRegStatusUpdateRspData)
	}
}

func RegistrationStatusUpdateProcedure(ueContextID string, ueRegStatusUpdateReqData models.UeRegStatusUpdateReqData) (
	*models.UeRegStatusUpdateRspData, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	// ueContextID must be a 5g GUTI (TS 29.518 6.1.3.2.4.5.1)
	if !strings.HasPrefix(ueContextID, "5g-guti") {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "UNSPECIFIED",
		}
		return nil, problemDetails
	}

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		logger.CtxLog.Warnf("AmfUe Context[%s] not found", ueContextID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}

	ue.Lock.Lock()
	defer ue.Lock.Unlock()

	ueRegStatusUpdateRspData := new(models.UeRegStatusUpdateRspData)

	if ueRegStatusUpdateReqData.TransferStatus == models.UeContextTransferStatus_TRANSFERRED {
		// remove the individual ueContext resource and release any PDU session(s)
		for _, pduSessionId := range ueRegStatusUpdateReqData.ToReleaseSessionList {
			cause := models.Cause_REL_DUE_TO_SLICE_NOT_AVAILABLE
			causeAll := &context.CauseAll{
				Cause: &cause,
			}
			smContext, ok := ue.SmContextFindByPDUSessionID(pduSessionId)
			if !ok {
				ue.ProducerLog.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionId)
				continue
			}
			problem, err := consumer.SendReleaseSmContextRequest(ue, smContext, causeAll, "", nil)
			if problem != nil {
				logger.GmmLog.Errorf("Release SmContext[pduSessionId: %d] Failed Problem[%+v]", pduSessionId, problem)
			} else if err != nil {
				logger.GmmLog.Errorf("Release SmContext[pduSessionId: %d] Error[%v]", pduSessionId, err.Error())
			}
		}

		if ueRegStatusUpdateReqData.PcfReselectedInd {
			problem, err := consumer.AMPolicyControlDelete(ue)
			if problem != nil {
				logger.GmmLog.Errorf("AM Policy Control Delete Failed Problem[%+v]", problem)
			} else if err != nil {
				logger.GmmLog.Errorf("AM Policy Control Delete Error[%v]", err.Error())
			}
		}
		// TODO: Currently only consider the 3GPP access type
		if !ue.UeCmRegistered[models.AccessType__3_GPP_ACCESS] {
			gmm_common.RemoveAmfUe(ue, false)
		}
	} else {
		// NOT_TRANSFERRED
		logger.CommLog.Debug("[AMF] RegistrationStatusUpdate: NOT_TRANSFERRED")
	}

	ueRegStatusUpdateRspData.RegStatusTransferComplete = true
	return ueRegStatusUpdateRspData, nil
}
