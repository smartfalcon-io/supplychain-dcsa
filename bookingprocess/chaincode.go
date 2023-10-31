package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ShippingContract defines the smart contract
type ShippingContract struct {
	contractapi.Contract
}

// ShippingAsset represents the shipping information
type ShippingAsset struct {
	BookingID                      string `json:"bookingID"`
	Name                           string `json:"name"`
	Address                        string `json:"address"`
	PhoneNumber                    int    `json:"phoneNumber"`
	ReceiptTypeAtOrigin            string `json:"receiptTypeAtOrigin"`
	DeliveryTypeAtDestination      string `json:"deliveryTypeAtDestination"`
	CargoMovementTypeAtOrigin       string `json:"cargoMovementTypeAtOrigin"`
	ServiceContractReference       string `json:"serviceContractReference"`
	CarrierServiceName             string `json:"carrierServiceName"`
	CarrierServiceCode             string `json:"carrierServiceCode"`
	UniversalServiceReference      string `json:"universalServiceReference"`
	CarrierExportVoyageNumber      string `json:"carrierExportVoyageNumber"`
	UniversalExportVoyageReference string `json:"universalExportVoyageReference"`
	DeclaredValueCurrency          string `json:"declaredValueCurrency"`
	IsPartialLoadAllowed           bool   `json:"isPartialLoadAllowed"`
	IsExportDeclarationRequired    bool   `json:"isExportDeclarationRequired"`
	ExportDeclarationReference     string `json:"exportDeclarationReference"`
	IsImportLicenseRequired        bool   `json:"isImportLicenseRequired"`
	ImportLicenseReference         string `json:"importLicenseReference"`
	ContractQuotationReference     string `json:"contractQuotationReference"`
	BookingChannelReference        string `json:"bookingChannelReference"`
	IncoTerms                      string `json:"incoTerms"`
	IsEquipmentSubstitutionAllowed bool   `json:"isEquipmentSubstitutionAllowed"`
}

// CreateAsset creates a new shipping asset
func (s *ShippingContract) CreateAsset(ctx contractapi.TransactionContextInterface, bookingID string, name string, address string, phoneNumber int, receiptTypeAtOrigin string, deliveryTypeAtDestination string, cargoMovementTypeAtOrigin string, serviceContractReference string, carrierServiceName string, carrierServiceCode string, universalServiceReference string, carrierExportVoyageNumber string, universalExportVoyageReference string, declaredValueCurrency string, isPartialLoadAllowed bool, isExportDeclarationRequired bool, exportDeclarationReference string, isImportLicenseRequired bool, importLicenseReference string, contractQuotationReference string, bookingChannelReference string, incoTerms string, isEquipmentSubstitutionAllowed bool) error {
	exists, err := s.AssetExists(ctx, bookingID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset with ID %s already exists", bookingID)
	}

	// Create a new asset
	asset := ShippingAsset{
		BookingID:                      bookingID,
		Name:                           name,
		Address:                        address,
		PhoneNumber:                    phoneNumber,
		ReceiptTypeAtOrigin:            receiptTypeAtOrigin,
		DeliveryTypeAtDestination:      deliveryTypeAtDestination,
		CargoMovementTypeAtOrigin:       cargoMovementTypeAtOrigin,
		ServiceContractReference:       serviceContractReference,
		CarrierServiceName:             carrierServiceName,
		CarrierServiceCode:             carrierServiceCode,
		UniversalServiceReference:      universalServiceReference,
		CarrierExportVoyageNumber:      carrierExportVoyageNumber,
		UniversalExportVoyageReference: universalExportVoyageReference,
		DeclaredValueCurrency:          declaredValueCurrency,
		IsPartialLoadAllowed:           isPartialLoadAllowed,
		IsExportDeclarationRequired:    isExportDeclarationRequired,
		ExportDeclarationReference:     exportDeclarationReference,
		IsImportLicenseRequired:        isImportLicenseRequired,
		ImportLicenseReference:         importLicenseReference,
		ContractQuotationReference:     contractQuotationReference,
		BookingChannelReference:        bookingChannelReference,
		IncoTerms:                      incoTerms,
		IsEquipmentSubstitutionAllowed: isEquipmentSubstitutionAllowed,
	}

	// Convert the asset to JSON
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	// Save the asset to the ledger
	err = ctx.GetStub().PutState(bookingID, assetJSON)
	if err != nil {
		return err
	}

	return nil
}

// ReadAsset retrieves a shipping asset by ID from the ledger
func (s *ShippingContract) ReadAsset(ctx contractapi.TransactionContextInterface, bookingID string) (*ShippingAsset, error) {
	assetJSON, err := ctx.GetStub().GetState(bookingID)
	if err != nil {
		return nil, err
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset with ID %s does not exist", bookingID)
	}

	var asset ShippingAsset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAllAssets retrieves all shipping assets from the ledger
func (s *ShippingContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*ShippingAsset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*ShippingAsset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset ShippingAsset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// AssetExists checks if a shipping asset with a given ID exists in the ledger
func (s *ShippingContract) AssetExists(ctx contractapi.TransactionContextInterface, bookingID string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(bookingID)
	if err != nil {
		return false, err
	}
	return assetJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&ShippingContract{})
	if err != nil {
		fmt.Printf("Error creating shipping chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting shipping chaincode: %s", err.Error())
	}
}
