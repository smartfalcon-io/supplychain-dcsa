package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type InvoiceContract struct {
	contractapi.Contract
}

type Location struct {
	LocationName             string `json:"locationName"`
	UNLocationCode           string `json:"UNLocationCode"`
	FacilityCode             string `json:"facilityCode,omitempty"`
	FacilityCodeListProvider string `json:"facilityCodeListProvider,omitempty"`
}

type PartyContactDetails struct {
	Name  string `json:"name"`
	Phone int    `json:"phone"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

type Address struct {
	Name         string `json:"name"`
	Street       string `json:"street"`
	StreetNumber int    `json:"streetNumber"`
	Floor        string `json:"floor,omitempty"`
	PostCode     int    `json:"postCode"`
	City         string `json:"city"`
	StateRegion  string `json:"stateRegion,omitempty"`
	Country      string `json:"country"`
}

type IdentifyingCode struct {
	DCSAResponsibleAgencyCode string `json:"DCSAResponsibleAgencyCode"`
	PartyCode                 string `json:"partyCode"`
	CodeListName              string `json:"codeListName"`
}

type Party struct {
	PartyName           string                `json:"partyName"`
	TaxReference1       string                `json:"taxReference1"`
	TaxReference2       string                `json:"taxReference2"`
	PublicKey           string                `json:"publicKey"`
	Address             Address               `json:"address"`
	PartyContactDetails []PartyContactDetails `json:"partyContactDetails"`
	IdentifyingCodes    []IdentifyingCode     `json:"identifyingCodes"`
}

type DocumentParty struct {
	Party            Party    `json:"party"`
	PartyFunction    string   `json:"partyFunction"`
	DisplayedAddress []string `json:"displayedAddress"`
	IsToBeNotified   bool     `json:"isToBeNotified"`
}

type Reference struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type RequestedEquipment struct {
	ISOEquipmentCode                string   `json:"ISOEquipmentCode"`
	TareWeight                      float64  `json:"tareWeight"`
	TareWeightUnit                  string   `json:"tareWeightUnit"`
	Units                           int      `json:"units"`
	EquipmentReferences             []string `json:"equipmentReferences"`
	IsShipperOwned                  bool     `json:"isShipperOwned"`
	CommodityRequestedEquipmentLink string   `json:"commodityRequestedEquipmentLink"`
}

type ShipmentLocation struct {
	Location                 Location `json:"location"`
	ShipmentLocationTypeCode string   `json:"shipmentLocationTypeCode"`
	EventDateTime            string   `json:"eventDateTime"`
}

type Commodity struct {
	CommodityType                   string  `json:"commodityType"`
	HSCode                          string  `json:"HSCode"`
	CargoGrossWeight                float64 `json:"cargoGrossWeight"`
	CargoGrossWeightUnit            string  `json:"cargoGrossWeightUnit"`
	CargoGrossVolume                float64 `json:"cargoGrossVolume"`
	CargoGrossVolumeUnit            string  `json:"cargoGrossVolumeUnit"`
	NumberOfPackages                int     `json:"numberOfPackages"`
	ExportLicenseIssueDate          string  `json:"exportLicenseIssueDate"`
	ExportLicenseExpiryDate         string  `json:"exportLicenseExpiryDate"`
	CommodityRequestedEquipmentLink string  `json:"commodityRequestedEquipmentLink"` // Updated field name
}

type ValueAddedService struct {
	ValueAddedServiceCode string `json:"valueAddedServiceCode"`
}

type Invoice struct {
	BookingID                                 string               `json:"bookingid"`
	ReceiptTypeAtOrigin                       string               `json:"receiptTypeAtOrigin"`
	DeliveryTypeAtDestination                 string               `json:"deliveryTypeAtDestination"`
	CargoMovementTypeAtOrigin                 string               `json:"cargoMovementTypeAtOrigin"`
	CargoMovementTypeAtDestination            string               `json:"cargoMovementTypeAtDestination"`
	ServiceContractReference                  string               `json:"serviceContractReference"`
	VesselName                                string               `json:"vesselName"`
	CarrierServiceName                        string               `json:"carrierServiceName"`
	CarrierServiceCode                        string               `json:"carrierServiceCode"`
	UniversalServiceReference                 string               `json:"universalServiceReference"`
	CarrierExportVoyageNumber                 string               `json:"carrierExportVoyageNumber"`
	UniversalExportVoyageReference            string               `json:"universalExportVoyageReference"`
	DeclaredValue                             float64              `json:"declaredValue"`
	DeclaredValueCurrency                     string               `json:"declaredValueCurrency"`
	PaymentTermCode                           string               `json:"paymentTermCode"`
	IsPartialLoadAllowed                      bool                 `json:"isPartialLoadAllowed"`
	IsExportDeclarationRequired               bool                 `json:"isExportDeclarationRequired"`
	ExportDeclarationReference                string               `json:"exportDeclarationReference"`
	IsImportLicenseRequired                   bool                 `json:"isImportLicenseRequired"`
	ImportLicenseReference                    string               `json:"importLicenseReference"`
	IsCustomsFilingSubmissionByShipper        bool                 `json:"isCustomsFilingSubmissionByShipper"`
	ContractQuotationReference                string               `json:"contractQuotationReference"`
	ExpectedDepartureDate                     string               `json:"expectedDepartureDate"`
	ExpectedArrivalAtPlaceOfDeliveryStartDate string               `json:"expectedArrivalAtPlaceOfDeliveryStartDate"`
	ExpectedArrivalAtPlaceOfDeliveryEndDate   string               `json:"expectedArrivalAtPlaceOfDeliveryEndDate"`
	TransportDocumentTypeCode                 string               `json:"transportDocumentTypeCode"`
	TransportDocumentReference                string               `json:"transportDocumentReference"`
	BookingChannelReference                   string               `json:"bookingChannelReference"`
	IncoTerms                                 string               `json:"incoTerms"`
	CommunicationChannelCode                  string               `json:"communicationChannelCode"`
	IsEquipmentSubstitutionAllowed            bool                 `json:"isEquipmentSubstitutionAllowed"`
	VesselIMONumber                           string               `json:"vesselIMONumber"`
	PreCarriageModeOfTransportCode            string               `json:"preCarriageModeOfTransportCode"`
	InvoicePayableAt                          Location             `json:"invoicePayableAt"`
	PlaceOfBLIssue                            Location             `json:"placeOfBLIssue"`
	Commodities                               []Commodity          `json:"commodities"`
	ValueAddedServices                        []ValueAddedService  `json:"valueAddedServices"`
	References                                []Reference          `json:"references"`
	RequestedEquipments                       []RequestedEquipment `json:"requestedEquipments"`
	DocumentParties                           []DocumentParty      `json:"documentParties"`
	ShipmentLocations                         []ShipmentLocation   `json:"shipmentLocations"`
}

func (ic *InvoiceContract) CreateInvoice(ctx contractapi.TransactionContextInterface, invoiceJSON string) error {
	var invoice Invoice
	err := json.Unmarshal([]byte(invoiceJSON), &invoice)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invoice JSON: %v", err)
	}
	err = ctx.GetStub().PutState("invoiceKey", []byte(invoiceJSON))
	if err != nil {
		return fmt.Errorf("failed to put state for invoice: %v", err)
	}

	return nil
}

func main() {
	contract := new(InvoiceContract)

	cc, err := contractapi.NewChaincode(contract)
	if err != nil {
		fmt.Printf("Error creating InvoiceChaincode: %s", err)
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting InvoiceChaincode: %s", err)
	}
}

//testing
