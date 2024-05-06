package test_test

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"test"
	"test/app"
	"test/consumerTestdata/UDM/TestGenAuthData"

	"github.com/stretchr/testify/assert"

	amf_factory "github.com/free5gc/amf/pkg/factory"
	amf_service "github.com/free5gc/amf/pkg/service"
	ausf_factory "github.com/free5gc/ausf/pkg/factory"
	ausf_service "github.com/free5gc/ausf/pkg/service"
	chf_factory "github.com/free5gc/chf/pkg/factory"
	chf_service "github.com/free5gc/chf/pkg/service"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap"
	nrf_factory "github.com/free5gc/nrf/pkg/factory"
	nrf_service "github.com/free5gc/nrf/pkg/service"
	nssf_factory "github.com/free5gc/nssf/pkg/factory"
	nssf_service "github.com/free5gc/nssf/pkg/service"
	"github.com/free5gc/openapi/models"
	pcf_factory "github.com/free5gc/pcf/pkg/factory"
	pcf_service "github.com/free5gc/pcf/pkg/service"
	smf_factory "github.com/free5gc/smf/pkg/factory"
	smf_service "github.com/free5gc/smf/pkg/service"
	udm_factory "github.com/free5gc/udm/pkg/factory"
	udm_service "github.com/free5gc/udm/pkg/service"
	"github.com/free5gc/udm/pkg/suci"
	udr_factory "github.com/free5gc/udr/pkg/factory"
	udr_service "github.com/free5gc/udr/pkg/service"
	"github.com/free5gc/util/mongoapi"
)

var NFs = []app.NetworkFunction{}

const (
	noInit = iota
	initNF
	multiAMF
)

var initFlag int = initNF

func init() {
	var testID string = ""
	var oauth bool

	for _, arg := range os.Args {
		values := strings.Split(arg, "=")
		for k, v := range values {
			if v == "-test.run" {
				testID = values[k+1]
				fmt.Printf("Run %s\n", testID)
			}
		}

		if arg == "noinit" {
			initFlag = noInit
		}

		if arg == "multiAmf" {
			initFlag = multiAMF
		}

		if arg == "oauth" {
			oauth = true
		}
	}

	switch initFlag {
	case multiAMF:
		// default key log path
		if err := os.MkdirAll("./log/", 0775); err != nil {
			fmt.Printf("Make directory %s failed: %+v", "./log/", err)
		}

		if err := nrfConfig(oauth); err != nil {
			fmt.Printf("NRF Config failed: %v\n", err)
		}
		nrfApp, _ := nrf_service.NewApp(nrf_factory.NrfConfig)
		NFs = append(NFs, nrfApp)

		if err := smfConfig(testID); err != nil {
			fmt.Printf("SMF Config failed: %v\n", err)
		}
		smfApp, _ := smf_service.NewApp(smf_factory.SmfConfig)
		NFs = append(NFs, smfApp)

		if err := udrConfig(); err != nil {
			fmt.Printf("UDR Config failed: %v\n", err)
		}
		udrApp, _ := udr_service.NewApp(udr_factory.UdrConfig)
		NFs = append(NFs, udrApp)

		if err := pcfConfig(); err != nil {
			fmt.Printf("PCF Config failed: %v\n", err)
		}
		pcfApp, _ := pcf_service.NewApp(pcf_factory.PcfConfig)
		NFs = append(NFs, pcfApp)

		if err := udmConfig(); err != nil {
			fmt.Printf("UDM Config failed: %v\n", err)
		}
		udmApp, _ := udm_service.NewApp(udm_factory.UdmConfig)
		NFs = append(NFs, udmApp)

		if err := nssfConfig(); err != nil {
			fmt.Printf("NSSF Config failed: %v\n", err)
		}
		nssfApp, _ := nssf_service.NewApp(nssf_factory.NssfConfig)
		NFs = append(NFs, nssfApp)

		if err := ausfConfig(); err != nil {
			fmt.Printf("AUSF Config failed: %v\n", err)
		}
		ausfApp, _ := ausf_service.NewApp(ausf_factory.AusfConfig)
		NFs = append(NFs, ausfApp)

		if err := chfConfig(); err != nil {
			fmt.Printf("CHF Config failed: %v\n", err)
		}
		chfApp, _ := chf_service.NewApp(chf_factory.ChfConfig)
		NFs = append(NFs, chfApp)

		for _, app := range NFs {
			go app.Start("")
			time.Sleep(200 * time.Millisecond)
		}
		if err := mongoapi.SetMongoDB("free5gc", "mongodb://127.0.0.1:27017"); err != nil {
			fmt.Printf("SetMongoDB failed: %v\n", err)
			return
		}
		fmt.Println("MongoDB Set")
	case initNF:
		// default key log path
		if err := os.MkdirAll("./log/", 0775); err != nil {
			fmt.Printf("Make directory %s failed: %+v", "./log/", err)
		}

		if err := nrfConfig(oauth); err != nil {
			fmt.Printf("NRF Config failed: %v\n", err)
		}
		nrfApp, _ := nrf_service.NewApp(nrf_factory.NrfConfig)
		NFs = append(NFs, nrfApp)

		if err := amfConfig(testID); err != nil {
			fmt.Printf("AMF Config failed: %v\n", err)
		}
		amfApp, _ := amf_service.NewApp(amf_factory.AmfConfig)
		NFs = append(NFs, amfApp)

		if err := smfConfig(testID); err != nil {
			fmt.Printf("SMF Config failed: %v\n", err)
		}
		smfApp, _ := smf_service.NewApp(smf_factory.SmfConfig)
		NFs = append(NFs, smfApp)

		if err := udrConfig(); err != nil {
			fmt.Printf("UDR Config failed: %v\n", err)
		}
		udrApp, _ := udr_service.NewApp(udr_factory.UdrConfig)
		NFs = append(NFs, udrApp)

		if err := pcfConfig(); err != nil {
			fmt.Printf("PCF Config failed: %v\n", err)
		}
		pcfApp, _ := pcf_service.NewApp(pcf_factory.PcfConfig)
		NFs = append(NFs, pcfApp)

		if err := udmConfig(); err != nil {
			fmt.Printf("UDM Config failed: %v\n", err)
		}
		udmApp, _ := udm_service.NewApp(udm_factory.UdmConfig)
		NFs = append(NFs, udmApp)

		if err := nssfConfig(); err != nil {
			fmt.Printf("NSSF Config failed: %v\n", err)
		}
		nssfApp, _ := nssf_service.NewApp(nssf_factory.NssfConfig)
		NFs = append(NFs, nssfApp)

		if err := ausfConfig(); err != nil {
			fmt.Printf("AUSF Config failed: %v\n", err)
		}
		ausfApp, _ := ausf_service.NewApp(ausf_factory.AusfConfig)
		NFs = append(NFs, ausfApp)

		if err := chfConfig(); err != nil {
			fmt.Printf("CHF Config failed: %v\n", err)
		}
		chfApp, _ := chf_service.NewApp(chf_factory.ChfConfig)
		NFs = append(NFs, chfApp)

		for _, app := range NFs {
			go app.Start("")
			time.Sleep(200 * time.Millisecond)
		}
	default:
		if err := mongoapi.SetMongoDB("free5gc", "mongodb://127.0.0.1:27017"); err != nil {
			fmt.Printf("SetMongoDB failed: %v\n", err)
			return
		}
		fmt.Println("MongoDB Set")
	}

}

func NfTerminate() {
	if initFlag != noInit {
		nfNums := len(NFs)
		for i := nfNums - 1; i >= 0; i-- {
			NFs[i].Terminate()
		}
	}
}

func TestNGSetup(t *testing.T) {
	var n int
	var sendMsg []byte
	var recvMsg = make([]byte, 2048)

	// RAN connect to AMF
	conn, err := test.ConnectToAmf("127.0.0.1", "127.0.0.1", 38412, 9487)
	assert.Nil(t, err)

	// send NGSetupRequest Msg
	sendMsg, err = test.GetNGSetupRequest([]byte("\x00\x01\x02"), 24, "free5gc")
	assert.Nil(t, err)
	_, err = conn.Write(sendMsg)
	assert.Nil(t, err)

	// receive NGSetupResponse Msg
	n, err = conn.Read(recvMsg)
	assert.Nil(t, err)
	_, err = ngap.Decoder(recvMsg[:n])
	assert.Nil(t, err)

	// close Connection
	conn.Close()
}

func TestCN(t *testing.T) {
	// New UE
	ue := test.NewRanUeContext("imsi-2089300007487", 1, security.AlgCiphering128NEA2, security.AlgIntegrity128NIA2,
		models.AccessType__3_GPP_ACCESS)
	// ue := test.NewRanUeContext("imsi-2089300007487", 1, security.AlgCiphering128NEA0, security.AlgIntegrity128NIA0, models.AccessType__3_GPP_ACCESS)
	ue.AmfUeNgapId = 1
	ue.AuthenticationSubs = test.GetAuthSubscription(TestGenAuthData.MilenageTestSet19.K,
		TestGenAuthData.MilenageTestSet19.OPC,
		TestGenAuthData.MilenageTestSet19.OP)
	// insert UE data to MongoDB

	servingPlmnId := "20893"
	// test.InsertAuthSubscriptionToMongoDB(ue.Supi, ue.AuthenticationSubs)

	test.InsertUeToMongoDB(t, ue, servingPlmnId)

	getData := test.GetAuthSubscriptionFromMongoDB(ue.Supi)
	assert.NotNil(t, getData)
	{
		amData := test.GetAccessAndMobilitySubscriptionData()
		test.InsertAccessAndMobilitySubscriptionDataToMongoDB(ue.Supi, amData, servingPlmnId)
		getData := test.GetAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		smfSelData := test.GetSmfSelectionSubscriptionData()
		test.InsertSmfSelectionSubscriptionDataToMongoDB(ue.Supi, smfSelData, servingPlmnId)
		getData := test.GetSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		smSelData := test.GetSessionManagementSubscriptionData()
		test.InsertSessionManagementSubscriptionDataToMongoDB(ue.Supi, servingPlmnId, smSelData)
		getData := test.GetSessionManagementDataFromMongoDB(ue.Supi, servingPlmnId)
		assert.NotNil(t, getData)
	}
	{
		amPolicyData := test.GetAmPolicyData()
		test.InsertAmPolicyDataToMongoDB(ue.Supi, amPolicyData)
		getData := test.GetAmPolicyDataFromMongoDB(ue.Supi)
		assert.NotNil(t, getData)
	}
	{
		smPolicyData := test.GetSmPolicyData()
		test.InsertSmPolicyDataToMongoDB(ue.Supi, smPolicyData)
		getData := test.GetSmPolicyDataFromMongoDB(ue.Supi)
		assert.NotNil(t, getData)
	}

	defer beforeClose(t, ue, servingPlmnId)

	// subscribe os signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Signal(syscall.SIGUSR1))
	<-c
}

func beforeClose(t *testing.T, ue *test.RanUeContext, servingPlmnId string) {
	// delete test data
	// test.DelAuthSubscriptionToMongoDB(ue.Supi)
	// test.DelAccessAndMobilitySubscriptionDataFromMongoDB(ue.Supi, "20893")
	// test.DelSmfSelectionSubscriptionDataFromMongoDB(ue.Supi, "20893")

	test.DelUeFromMongoDB(t, ue, servingPlmnId)
}

func nrfConfig(oauth bool) error {
	nrf_factory.NrfConfig = &nrf_factory.Config{
		Info: &nrf_factory.Info{
			Version:     "1.0.2",
			Description: "NRF initial test configuration",
		},
		Configuration: &nrf_factory.Configuration{
			MongoDBName: "free5gc",
			MongoDBUrl:  "mongodb://127.0.0.1:27017",
			Sbi: &nrf_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.10",
				BindingIPv4:  "127.0.0.10",
				Port:         8000,
				Cert: &nrf_factory.Cert{
					Pem: "../cert/nrf.pem",
					Key: "../cert/nrf.key",
				},
				RootCert: &nrf_factory.Cert{
					Pem: "../cert/root.pem",
					Key: "../cert/root.key",
				},
				OAuth: oauth,
			},
			DefaultPlmnId: models.PlmnId{
				Mcc: "208",
				Mnc: "93",
			},
			ServiceNameList: []string{
				"nnrf-nfm",
				"nnrf-disc",
			},
		},
		Logger: &nrf_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if _, err := nrf_factory.NrfConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func amfConfig(testID string) error {
	var ngapIpList, integrityOrder, cipheringOrder []string

	if testID == "TestCN" {
		ngapIpList = []string{"10.200.200.1"}
		integrityOrder = []string{"NIA2", "NIA0"}
		cipheringOrder = []string{"NEA2", "NEA0"}
	} else {
		ngapIpList = []string{"127.0.0.1"}
		integrityOrder = []string{"NIA2"}
		cipheringOrder = []string{"NEA0"}
	}

	amf_factory.AmfConfig = &amf_factory.Config{
		Info: &amf_factory.Info{
			Version:     "1.0.9",
			Description: "AMF initial test configuration",
		},
		Configuration: &amf_factory.Configuration{
			AmfName:    "AMF",
			NgapIpList: ngapIpList,
			NgapPort:   38412,
			Sbi: &amf_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.18",
				BindingIPv4:  "127.0.0.18",
				Port:         8000,
				Tls: &amf_factory.Tls{
					Pem: "cert/amf.pem",
					Key: "cert/amf.key",
				},
			},
			ServiceNameList: []string{
				"namf-comm",
				"namf-evts",
				"namf-mt",
				"namf-loc",
				"namf-oam",
			},
			ServedGumaiList: []models.Guami{{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				AmfId: "cafe00",
			}},
			SupportTAIList: []models.Tai{{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				Tac: "000001",
			}},
			PlmnSupportList: []amf_factory.PlmnSupportItem{{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				SNssaiList: []models.Snssai{{
					Sst: 1,
					Sd:  "fedcba",
				}, {
					Sst: 1,
					Sd:  "112233",
				}},
			}},
			SupportDnnList: []string{
				"internet",
			},
			NrfUri:     "http://127.0.0.10:8000",
			NrfCertPem: "../cert/nrf.pem",
			Security: &amf_factory.Security{
				IntegrityOrder: integrityOrder,
				CipheringOrder: cipheringOrder,
			},
			NetworkName: amf_factory.NetworkName{
				Full:  "free5GC",
				Short: "free",
			},
			NgapIE: &amf_factory.NgapIE{
				MobilityRestrictionList: &amf_factory.MobilityRestrictionList{
					Enable: true,
				},
				MaskedIMEISV: &amf_factory.MaskedIMEISV{
					Enable: true,
				},
				RedirectionVoiceFallback: &amf_factory.RedirectionVoiceFallback{
					Enable: false,
				},
			},
			NasIE: &amf_factory.NasIE{
				NetworkFeatureSupport5GS: &amf_factory.NetworkFeatureSupport5GS{
					Enable:  true,
					Length:  1,
					ImsVoPS: 0,
					Emc:     0,
					Emf:     0,
					IwkN26:  0,
					Mpsi:    0,
					EmcN3:   0,
					Mcsi:    0,
				},
			},
			T3502Value:             720,
			T3512Value:             3600,
			Non3gppDeregTimerValue: 3240,
			T3513: amf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    6000000000,
				MaxRetryTimes: 4,
			},
			T3522: amf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    6000000000,
				MaxRetryTimes: 4,
			},
			T3550: amf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    6000000000,
				MaxRetryTimes: 4,
			},
			T3560: amf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    6000000000,
				MaxRetryTimes: 4,
			},
			T3565: amf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    6000000000,
				MaxRetryTimes: 4,
			},
			T3570: amf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    6000000000,
				MaxRetryTimes: 4,
			},
		},
		Logger: &amf_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if _, err := amf_factory.AmfConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func smfConfig(testID string) error {
	var dnaiList []string

	if testID == "TestAFInfluenceOnTrafficRouting" {
		dnaiList = []string{"edge"}
	}

	smf_factory.SmfConfig = &smf_factory.Config{
		Info: &smf_factory.Info{
			Version:     "1.0.7",
			Description: "SMF initial single test configuration",
		},
		Configuration: &smf_factory.Configuration{
			SmfName: "SMF",
			Sbi: &smf_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.2",
				BindingIPv4:  "127.0.0.2",
				Port:         8000,
				Tls: &smf_factory.Tls{
					Pem: "cert/smf.pem",
					Key: "cert/smf.key",
				},
			},
			ServiceNameList: []string{
				"nsmf-pdusession",
				"nsmf-event-exposure",
				"nsmf-oam",
			},
			SNssaiInfo: []*smf_factory.SnssaiInfoItem{{
				SNssai: &models.Snssai{
					Sst: 1,
					Sd:  "fedcba",
				},
				DnnInfos: []*smf_factory.SnssaiDnnInfoItem{{
					Dnn: "internet",
					DNS: &smf_factory.DNS{
						IPv4Addr: "8.8.8.8",
						IPv6Addr: "2001:4860:4860::8888",
					},
				}},
			}, {
				SNssai: &models.Snssai{
					Sst: 1,
					Sd:  "112233",
				},
				DnnInfos: []*smf_factory.SnssaiDnnInfoItem{
					{
						Dnn: "internet",
						DNS: &smf_factory.DNS{
							IPv4Addr: "8.8.8.8",
							IPv6Addr: "2001:4860:4860::8888",
						},
					},
					{
						Dnn: "internet2",
						DNS: &smf_factory.DNS{
							IPv4Addr: "8.8.8.8",
							IPv6Addr: "2001:4860:4860::8888",
						},
					},
				},
			}},
			PFCP: &smf_factory.PFCP{
				NodeID:       "10.200.200.1",
				ExternalAddr: "10.200.200.1",
				ListenAddr:   "10.200.200.1",
			},
			Locality: "area1",
			UserPlaneInformation: smf_factory.UserPlaneInformation{
				UPNodes: map[string]*smf_factory.UPNode{
					"gNB1": {
						Type: "AN",
						ANIP: "192.188.2.3",
					},
					"UPF": {
						Type:   "UPF",
						NodeID: "10.200.200.101",
						Addr:   "10.200.200.101",
						SNssaiInfos: []*smf_factory.SnssaiUpfInfoItem{{
							SNssai: &models.Snssai{
								Sst: 1,
								Sd:  "fedcba",
							},
							DnnUpfInfoList: []*smf_factory.DnnUpfInfoItem{{
								Dnn:      "internet",
								DnaiList: dnaiList,
								Pools: []*smf_factory.UEIPPool{{
									Cidr: "10.60.0.0/16",
								}},
							}},
						}, {
							SNssai: &models.Snssai{
								Sst: 1,
								Sd:  "112233",
							},
							DnnUpfInfoList: []*smf_factory.DnnUpfInfoItem{{
								Dnn: "internet",
								Pools: []*smf_factory.UEIPPool{{
									Cidr: "10.61.0.0/16",
								}},
							}},
						}},
						InterfaceUpfInfoList: []*smf_factory.InterfaceUpfInfoItem{{
							InterfaceType: "N3",
							Endpoints: []string{
								"10.200.200.102",
							},
							NetworkInstances: []string{
								"internet",
							},
						}},
					},
				},
				Links: []*smf_factory.UPLink{{
					A: "gNB1",
					B: "UPF",
				}},
			},
			T3591: &smf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    5 * time.Second,
				MaxRetryTimes: 2,
			},
			T3592: &smf_factory.TimerValue{
				Enable:        true,
				ExpireTime:    5 * time.Second,
				MaxRetryTimes: 2,
			},
			NrfUri:     "http://127.0.0.10:8000",
			NrfCertPem: "../cert/nrf.pem",
		},
		Logger: &smf_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if testID == "TestRequestTwoPDUSessions" {
		smf_factory.SmfConfig.Configuration.UserPlaneInformation.Links =
			append(smf_factory.SmfConfig.Configuration.UserPlaneInformation.Links, &smf_factory.UPLink{
				A: "gNB1",
				B: "UPF2",
			})

		smf_factory.SmfConfig.Configuration.UserPlaneInformation.UPNodes["UPF2"] =
			&smf_factory.UPNode{
				Type:   "UPF",
				NodeID: "10.200.200.102",
				Addr:   "10.200.200.102",
				SNssaiInfos: []*smf_factory.SnssaiUpfInfoItem{{
					SNssai: &models.Snssai{
						Sst: 1,
						Sd:  "112233",
					},
					DnnUpfInfoList: []*smf_factory.DnnUpfInfoItem{{
						Dnn: "internet2",
						Pools: []*smf_factory.UEIPPool{{
							Cidr: "10.62.0.0/16",
						}},
					}},
				}},
				InterfaceUpfInfoList: []*smf_factory.InterfaceUpfInfoItem{{
					InterfaceType: "N3",
					Endpoints: []string{
						"10.200.200.102",
					},
					NetworkInstances: []string{
						"internet2",
					},
				}},
			}
	}

	smfUeRoutingConfig()

	if _, err := smf_factory.SmfConfig.Validate(); err != nil {
		return err
	}

	if _, err := smf_factory.UERoutingConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func smfUeRoutingConfig() {
	smf_factory.UERoutingConfig = &smf_factory.RoutingConfig{
		Info: &smf_factory.Info{
			Version:     "1.0.7",
			Description: "SMF initial test routing information for UE",
		},
		UERoutingInfo: map[string]smf_factory.UERoutingInfo{
			"UE1": {
				Members: []string{
					"imsi-208930000007487",
				},
				Topology: []smf_factory.UPLink{{
					A: "gNB1",
					B: "BranchingUPF",
				}, {
					A: "BranchingUPF",
					B: "AnchorUPF1",
				}},
				SpecificPaths: []smf_factory.SpecificPath{{
					DestinationIP: "10.60.0.103/32",
					Path: []string{
						"BranchingUPF",
						"AnchorUPF2",
					},
				}},
			},
			"UE2": {
				Members: []string{
					"imsi-208930000007486",
				},
				Topology: []smf_factory.UPLink{{
					A: "gNB1",
					B: "BranchingUPF",
				}, {
					A: "BranchingUPF",
					B: "AnchorUPF1",
				}},
				SpecificPaths: []smf_factory.SpecificPath{{
					DestinationIP: "10.0.0.11/32",
					Path: []string{
						"BranchingUPF",
						"AnchorUPF2",
					},
				}},
			},
		},
		RouteProf: map[smf_factory.RouteProfID]smf_factory.RouteProfile{
			"MEC1": {
				ForwardingPolicyID: "10",
			},
		},
		PfdDatas: []*smf_factory.PfdDataForApp{{
			AppID: "edge",
			Pfds: []smf_factory.PfdContent{{
				PfdID: "pfd1",
				FlowDescriptions: []string{
					"permit out ip from 10.60.0.0/24 8080 to 10.60.0.1",
				},
			}},
		}},
	}
}

func udrConfig() error {
	udr_factory.UdrConfig = &udr_factory.Config{
		Info: &udr_factory.Info{
			Version:     "1.0.2",
			Description: "UDR initial test configuration",
		},
		Configuration: &udr_factory.Configuration{
			Sbi: &udr_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.4",
				BindingIPv4:  "127.0.0.4",
				Port:         8000,
				Tls: &udr_factory.Tls{
					Pem: "cert/udr.pem",
					Key: "cert/udr.key",
				},
			},
			Mongodb: &udr_factory.Mongodb{
				Name: "free5gc",
				Url:  "mongodb://localhost:27017",
			},
			NrfUri:     "http://127.0.0.10:8000",
			NrfCertPem: "../cert/nrf.pem",
		},
		Logger: &udr_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if _, err := udr_factory.UdrConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func pcfConfig() error {
	pcf_factory.PcfConfig = &pcf_factory.Config{
		Info: &pcf_factory.Info{
			Version:     "1.0.2",
			Description: "PCF initial test configuration",
		},
		Configuration: &pcf_factory.Configuration{
			PcfName: "PCF",
			Sbi: &pcf_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.7",
				BindingIPv4:  "127.0.0.7",
				Port:         8000,
				Tls: &pcf_factory.Tls{
					Pem: "cert/pcf.pem",
					Key: "cert/pcf.key",
				},
			},
			TimeFormat:      "2019-01-02 15:04:05",
			DefaultBdtRefId: "BdtPolicyId-",
			NrfUri:          "http://127.0.0.10:8000",
			NrfCertPem:      "../cert/nrf.pem",
			ServiceList: []pcf_factory.Service{{
				ServiceName: "npcf-am-policy-control",
			}, {
				ServiceName: "npcf-smpolicycontrol",
				SuppFeat:    "3fff",
			}, {
				ServiceName: "npcf-bdtpolicycontrol",
			}, {
				ServiceName: "npcf-policyauthorization",
				SuppFeat:    "3",
			}, {
				ServiceName: "npcf-eventexposure",
			}, {
				ServiceName: "npcf-ue-policy-control",
			}},
			Mongodb: &pcf_factory.Mongodb{
				Name: "free5gc",
				Url:  "mongodb://localhost:27017",
			},
			Locality: "area1",
		},
		Logger: &pcf_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if _, err := pcf_factory.PcfConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func udmConfig() error {
	udm_factory.UdmConfig = &udm_factory.Config{
		Info: &udm_factory.Info{
			Version:     "1.0.3",
			Description: "UDM initial test configuration",
		},
		Configuration: &udm_factory.Configuration{
			ServiceNameList: []string{
				"nudm-sdm",
				"nudm-uecm",
				"nudm-ueau",
				"nudm-ee",
				"nudm-pp",
			},
			Sbi: &udm_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.3",
				BindingIPv4:  "127.0.0.3",
				Port:         8000,
				Tls: &udm_factory.Tls{
					Pem: "cert/udm.pem",
					Key: "cert/udm.key",
				},
			},
			NrfUri:     "http://127.0.0.10:8000",
			NrfCertPem: "../cert/nrf.pem",
			SuciProfiles: []suci.SuciProfile{
				{
					ProtectionScheme: "1", // Protect Scheme: Profile A
					PrivateKey:       "c53c22208b61860b06c62e5406a7b330c2b577aa5558981510d128247d38bd1d",
					PublicKey:        "5a8d38864820197c3394b92613b20b91633cbd897119273bf8e4a6f4eec0a650",
				},
				{
					ProtectionScheme: "2", // Protect Scheme: Profile B
					PrivateKey:       "F1AB1074477EBCC7F554EA1C5FC368B1616730155E0041AC447D6301975FECDA",
					PublicKey:        "0472DA71976234CE833A6907425867B82E074D44EF907DFB4B3E21C1C2256EBCD15A7DED52FCBB097A4ED250E036C7B9C8C7004C4EEDC4F068CD7BF8D3F900E3B4",
				},
			},
		},
		Logger: &udm_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if _, err := udm_factory.UdmConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func nssfConfig() error {
	var accessType3GPP models.AccessType = models.AccessType__3_GPP_ACCESS

	nssf_factory.NssfConfig = &nssf_factory.Config{
		Info: &nssf_factory.Info{
			Version:     "1.0.2",
			Description: "NSSF initial test configuration",
		},
		Configuration: &nssf_factory.Configuration{
			NssfName: "NSSF",
			Sbi: &nssf_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.31",
				BindingIPv4:  "127.0.0.31",
				Port:         8000,
				Tls: &nssf_factory.Tls{
					Pem: "cert/nssf.pem",
					Key: "cert/nssf.key",
				},
			},
			ServiceNameList: []models.ServiceName{
				"nnssf-nsselection",
				"nnssf-nssaiavailability",
			},
			NrfUri:     "http://127.0.0.10:8000",
			NrfCertPem: "../cert/nrf.pem",
			SupportedPlmnList: []models.PlmnId{{
				Mcc: "208",
				Mnc: "93",
			}},
			SupportedNssaiInPlmnList: []nssf_factory.SupportedNssaiInPlmn{{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				SupportedSnssaiList: []models.Snssai{{
					Sst: 1,
					Sd:  "fedcba",
				}, {
					Sst: 1,
					Sd:  "112233",
				}, {
					Sst: 1,
					Sd:  "000003",
				}, {
					Sst: 2,
					Sd:  "000001",
				}, {
					Sst: 2,
					Sd:  "000002",
				}},
			}},
			NsiList: []nssf_factory.NsiConfig{{
				Snssai: &models.Snssai{
					Sst: 1,
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "10",
				}},
			}, {
				Snssai: &models.Snssai{
					Sst: 1,
					Sd:  "000001",
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "11",
				}},
			}, {
				Snssai: &models.Snssai{
					Sst: 1,
					Sd:  "000002",
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "12",
				}, {
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "12",
				}},
			}, {
				Snssai: &models.Snssai{
					Sst: 1,
					Sd:  "000003",
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "13",
				}},
			}, {
				Snssai: &models.Snssai{
					Sst: 2,
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "20",
				}},
			}, {
				Snssai: &models.Snssai{
					Sst: 2,
					Sd:  "000001",
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "21",
				}},
			}, {
				Snssai: &models.Snssai{
					Sst: 1,
					Sd:  "fedcba",
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "22",
				}},
			}, {
				Snssai: &models.Snssai{
					Sst: 1,
					Sd:  "112233",
				},
				NsiInformationList: []models.NsiInformation{{
					NrfId: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
					NsiId: "23",
				}},
			}},
			AmfSetList: []nssf_factory.AmfSetConfig{{
				AmfSetId: "1",
				AmfList: []string{
					"ffa2e8d7-3275-49c7-8631-6af1df1d9d26",
					"0e8831c3-6286-4689-ab27-1e2161e15cb1",
					"a1fba9ba-2e39-4e22-9c74-f749da571d0d",
				},
				NrfAmfSet: "http://127.0.0.10:8000/nnrf-nfm/v1/nf-instances",
				SupportedNssaiAvailabilityData: []models.SupportedNssaiAvailabilityData{{
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33456",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 1,
						Sd:  "000002",
					}, {
						Sst: 2,
						Sd:  "000001",
					}},
				}, {
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33457",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 1,
						Sd:  "000002",
					}},
				}},
			}, {
				AmfSetId:  "2",
				NrfAmfSet: "http://localhost:8084/nnrf-nfm/v1/nf-instances",
				SupportedNssaiAvailabilityData: []models.SupportedNssaiAvailabilityData{{
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33456",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 1,
						Sd:  "000003",
					}, {
						Sst: 2,
						Sd:  "000001",
					}},
				}, {
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33458",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 2,
					}},
				}},
			}},
			AmfList: []nssf_factory.AmfConfig{{
				NfId: "469de254-2fe5-4ca0-8381-af3f500af77c",
				SupportedNssaiAvailabilityData: []models.SupportedNssaiAvailabilityData{{
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33456",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000002",
					}, {
						Sst: 2,
					}},
				}, {
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33457",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 1,
						Sd:  "000002",
					}},
				}},
			}, {
				NfId: "fbe604a8-27b2-417e-bd7c-8a7be2691f8d",
				SupportedNssaiAvailabilityData: []models.SupportedNssaiAvailabilityData{{
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33458",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 1,
						Sd:  "000003",
					}, {
						Sst: 2,
					}},
				}, {
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33459",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 2,
					}, {
						Sst: 2,
						Sd:  "000001",
					}},
				}},
			}, {
				NfId: "b9e6e2cb-5ce8-4cb6-9173-a266dd9a2f0c",
				SupportedNssaiAvailabilityData: []models.SupportedNssaiAvailabilityData{{
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33456",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 1,
						Sd:  "000002",
					}, {
						Sst: 2,
					}},
				}, {
					Tai: &models.Tai{
						PlmnId: &models.PlmnId{
							Mcc: "466",
							Mnc: "92",
						},
						Tac: "33458",
					},
					SupportedSnssaiList: []models.Snssai{{
						Sst: 1,
					}, {
						Sst: 1,
						Sd:  "000001",
					}, {
						Sst: 2,
					}, {
						Sst: 2,
						Sd:  "000001",
					}},
				}},
			}},
			TaList: []nssf_factory.TaConfig{{
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "92",
					},
					Tac: "33456",
				},
				AccessType: &accessType3GPP,
				SupportedSnssaiList: []models.Snssai{{
					Sst: 1,
				}, {
					Sst: 1,
					Sd:  "000001",
				}, {
					Sst: 1,
					Sd:  "000002",
				}, {
					Sst: 2,
				}},
			}, {
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "92",
					},
					Tac: "33458",
				},
				AccessType: &accessType3GPP,
				SupportedSnssaiList: []models.Snssai{{
					Sst: 1,
				}, {
					Sst: 1,
					Sd:  "000001",
				}, {
					Sst: 1,
					Sd:  "000003",
				}, {
					Sst: 2,
				}},
				RestrictedSnssaiList: []models.RestrictedSnssai{{
					HomePlmnId: &models.PlmnId{
						Mcc: "310",
						Mnc: "560",
					},
					SNssaiList: []models.Snssai{{
						Sst: 1,
						Sd:  "000003",
					}},
				}},
			}, {
				Tai: &models.Tai{
					PlmnId: &models.PlmnId{
						Mcc: "466",
						Mnc: "92",
					},
					Tac: "33459",
				},
				AccessType: &accessType3GPP,
				SupportedSnssaiList: []models.Snssai{{
					Sst: 1,
				}, {
					Sst: 1,
					Sd:  "000001",
				}, {
					Sst: 2,
				}, {
					Sst: 2,
					Sd:  "000001",
				}},
				RestrictedSnssaiList: []models.RestrictedSnssai{{
					HomePlmnId: &models.PlmnId{
						Mcc: "310",
						Mnc: "560",
					},
					SNssaiList: []models.Snssai{{
						Sst: 2,
						Sd:  "000001",
					}},
				}},
			}},
			MappingListFromPlmn: []nssf_factory.MappingFromPlmnConfig{{
				OperatorName: "NTT Docomo",
				HomePlmnId: &models.PlmnId{
					Mcc: "440",
					Mnc: "10",
				},
				MappingOfSnssai: []models.MappingOfSnssai{{
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000001",
					},
					HomeSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "1",
					},
				}, {
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000002",
					},
					HomeSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000003",
					},
				}, {
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000003",
					},
					HomeSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000004",
					},
				}, {
					ServingSnssai: &models.Snssai{
						Sst: 2,
						Sd:  "000001",
					},
					HomeSnssai: &models.Snssai{
						Sst: 2,
						Sd:  "000002",
					},
				}},
			}, {
				OperatorName: "AT&T Mobility",
				HomePlmnId: &models.PlmnId{
					Mcc: "310",
					Mnc: "560",
				},
				MappingOfSnssai: []models.MappingOfSnssai{{
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000001",
					},
					HomeSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000002",
					},
				}, {
					ServingSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000002",
					},
					HomeSnssai: &models.Snssai{
						Sst: 1,
						Sd:  "000003",
					},
				}},
			}},
		},
		Logger: &nssf_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	return nil
}

func ausfConfig() error {
	ausf_factory.AusfConfig = &ausf_factory.Config{
		Info: &ausf_factory.Info{
			Version:     "1.0.3",
			Description: "AUSF initial test configuration",
		},
		Configuration: &ausf_factory.Configuration{
			Sbi: &ausf_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.9",
				BindingIPv4:  "127.0.0.9",
				Port:         8000,
				Tls: &ausf_factory.Tls{
					Pem: "cert/ausf.pem",
					Key: "cert/ausf.key",
				},
			},
			ServiceNameList: []string{
				"nausf-auth",
			},
			NrfUri:     "http://127.0.0.10:8000",
			NrfCertPem: "../cert/nrf.pem",
			PlmnSupportList: []models.PlmnId{{
				Mcc: "208",
				Mnc: "93",
			}, {
				Mcc: "123",
				Mnc: "45",
			}},
			GroupId: "ausfGroup001",
		},
		Logger: &ausf_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if _, err := ausf_factory.AusfConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func chfConfig() error {
	chf_factory.ChfConfig = &chf_factory.Config{
		Info: &chf_factory.Info{
			Version:     "1.0.3",
			Description: "CHF initial test configuration",
		},
		Configuration: &chf_factory.Configuration{
			ChfName: "CHF",
			Sbi: &chf_factory.Sbi{
				Scheme:       "http",
				RegisterIPv4: "127.0.0.113",
				BindingIPv4:  "127.0.0.113",
				Port:         8000,
				Tls: &chf_factory.Tls{
					Pem: "../cert/chf.pem",
					Key: "../cert/chf.key",
				},
			},
			NrfUri:     "http://127.0.0.10:8000",
			NrfCertPem: "../cert/nrf.pem",
			ServiceNameList: []string{
				"nchf-convergedcharging",
			},
			Mongodb: &chf_factory.Mongodb{
				Name: "free5gc",
				Url:  "mongodb://localhost:27017",
			},
			QuotaValidityTime:   10000,
			VolumeLimit:         50000,
			VolumeLimitPDU:      10000,
			VolumeThresholdRate: 0.8,
			Cgf: &chf_factory.Cgf{
				HostIPv4:   "127.0.0.1",
				Port:       2122,
				ListenPort: 2121,
				Tls: &chf_factory.Tls{
					Pem: "../cert/chf.pem",
					Key: "../cert/chf.key",
				},
			},
			AbmfDiameter: &chf_factory.Diameter{
				Protocol: "tcp",
				HostIPv4: "127.0.0.113",
				Port:     3868,
				Tls: &chf_factory.Tls{
					Pem: "../cert/chf.pem",
					Key: "../cert/chf.key",
				},
			},
			RfDiameter: &chf_factory.Diameter{
				Protocol: "tcp",
				HostIPv4: "127.0.0.113",
				Port:     3869,
				Tls: &chf_factory.Tls{
					Pem: "../cert/chf.pem",
					Key: "../cert/chf.key",
				},
			},
		},
		Logger: &chf_factory.Logger{
			Enable:       true,
			Level:        "info",
			ReportCaller: false,
		},
	}

	if _, err := chf_factory.ChfConfig.Validate(); err != nil {
		return err
	}

	return nil
}
