package api

import (
	"testing"

	"github.com/happyreturns/fedex/models"
)

var (
	expectedRequest = models.Request{
		WebAuthenticationDetail: models.WebAuthenticationDetail{
			UserCredential: models.UserCredential{
				Key:      "Key",
				Password: "Password",
			},
		},
		ClientDetail: models.ClientDetail{
			AccountNumber: "Account",
			MeterNumber:   "Meter",
		},
		Version: models.Version{
			ServiceID: "ship",
			Major:     23,
		},
	}
	testAPI = API{
		Key:      "Key",
		Password: "Password",
		Account:  "Account",
		Meter:    "Meter",
	}
)

func TestGroundShipmentNotInternational(t *testing.T) {
	shipment := &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1511 15th Street"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90404",
				CountryCode:         "US",
			},
			FromContact: models.Contact{
				PersonName:  "Joe Customer",
				PhoneNumber: "2045551234",
			},
			ToContact: models.Contact{
				PersonName:  "Returns Department",
				CompanyName: "FedEx",
				PhoneNumber: "9015551234",
			},
			ToAddress: models.Address{
				StreetLines:         []string{"1106 Broadway"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90404",
				CountryCode:         "US",
			},
		},
		NotificationEmail: "NotificationEmail",
		Reference:         "REF",
		Service:           "FEDEX_GROUND",
	}
	envelope, err := testAPI.processShipmentRequest(shipment)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := envelope.Body.(models.ProcessShipmentBody); !ok {
		t.Fatal("should be process shipment body")
	}

	shipmentBody := envelope.Body.(models.ProcessShipmentBody)
	processShipment := shipmentBody.ProcessShipmentRequest
	if processShipment.Request != expectedRequest {
		t.Fatal("Request doesn't match")
	}

	if processShipment.RequestedShipment.DropoffType != "REGULAR_PICKUP" {
		t.Fatal("DropoffType doesn't match")
	}

	if processShipment.RequestedShipment.ServiceType != "FEDEX_GROUND" {
		t.Fatal("ServiceType doesn't match")
	}

	if processShipment.RequestedShipment.PackagingType != "YOUR_PACKAGING" {
		t.Fatal("PackagingType doesn't match")
	}

	if contact := processShipment.RequestedShipment.Shipper.Contact; contact != shipment.FromContact {
		t.Fatal("contact doesn't match")
	}

	if address := processShipment.RequestedShipment.Shipper.Address; len(address.StreetLines) != 1 ||
		address.StreetLines[0] != "1511 15th Street" ||
		address.City != "Santa Monica" ||
		address.StateOrProvinceCode != "CA" ||
		address.PostalCode != "90404" ||
		address.CountryCode != "US" {
		t.Fatal("address doesn't match")
	}

	if contact := processShipment.RequestedShipment.Recipient.Contact; contact != shipment.ToContact {
		t.Fatal("contact doesn't match")
	}

	if address := processShipment.RequestedShipment.Recipient.Address; len(address.StreetLines) != 1 ||
		address.StreetLines[0] != "1106 Broadway" ||
		address.City != "Santa Monica" ||
		address.StateOrProvinceCode != "CA" ||
		address.PostalCode != "90404" ||
		address.CountryCode != "US" {
		t.Fatal("address doesn't match")
	}

	if processShipment.RequestedShipment.ShippingChargesPayment.PaymentType != "SENDER" {
		t.Fatal("PaymentType doesn't match")
	}

	if ssr := processShipment.RequestedShipment.SpecialServicesRequested; ssr.SpecialServiceTypes[0] != "EVENT_NOTIFICATION" ||
		ssr.EventNotificationDetail.AggregationType != "PER_SHIPMENT" ||
		ssr.EventNotificationDetail.EventNotifications[0].Role != "SHIPPER" ||
		len(ssr.EventNotificationDetail.EventNotifications[0].Events) != 5 ||
		ssr.EventNotificationDetail.EventNotifications[0].Events[0] != "ON_DELIVERY" ||
		ssr.EventNotificationDetail.EventNotifications[0].Events[1] != "ON_ESTIMATED_DELIVERY" ||
		ssr.EventNotificationDetail.EventNotifications[0].Events[2] != "ON_EXCEPTION" ||
		ssr.EventNotificationDetail.EventNotifications[0].Events[3] != "ON_SHIPMENT" ||
		ssr.EventNotificationDetail.EventNotifications[0].Events[4] != "ON_TENDER" ||
		ssr.EventNotificationDetail.EventNotifications[0].NotificationDetail.NotificationType != "EMAIL" ||
		ssr.EventNotificationDetail.EventNotifications[0].NotificationDetail.EmailDetail.EmailAddress != "NotificationEmail" ||
		ssr.EventNotificationDetail.EventNotifications[0].NotificationDetail.EmailDetail.Name != "Happy Returns dev team" ||
		ssr.EventNotificationDetail.EventNotifications[0].NotificationDetail.Localization.LanguageCode != "en" ||
		ssr.EventNotificationDetail.EventNotifications[0].FormatSpecification.Type != "HTML" {
		t.Fatal("specialServicesRequested doesn't match")
	}

	// not gonna bother validating every single field
}

func TestGroundShipmentInternational(t *testing.T) {
	commodities := []models.Commodity{
		{
			NumberOfPieces:       1,
			Description:          "Computer Keyboard",
			CountryOfManufacture: "US",
			Weight:               models.Weight{Units: "LB", Value: 10.0},
			Quantity:             1,
			QuantityUnits:        "pcs",
			UnitPrice:            &models.Money{Currency: "USD", Amount: 25.00},
			CustomsValue:         &models.Money{Currency: "USD", Amount: 30.00},
		},
		{
			NumberOfPieces:       1,
			Description:          "Computer Monitor",
			CountryOfManufacture: "US",
			Weight:               models.Weight{Units: "LB", Value: 5.0},
			Quantity:             1,
			QuantityUnits:        "pcs",
			UnitPrice:            &models.Money{Currency: "USD", Amount: 214.42},
			CustomsValue:         &models.Money{Currency: "USD", Amount: 381.12},
		},
	}
	shipment := &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1234 Main Street", "Suite 200"},
				City:                "Winnipeg",
				StateOrProvinceCode: "MB",
				PostalCode:          "R2M4B5",
				CountryCode:         "CA",
			},
			FromContact: models.Contact{
				PersonName:  "Joe Customer",
				PhoneNumber: "2045551234",
			},
			ToContact: models.Contact{
				PersonName:  "Returns Department",
				CompanyName: "FedEx",
				PhoneNumber: "9015551234",
			},
			ToAddress: models.Address{
				StreetLines:         []string{"3610 Hacks Cross Road", "First Floor"},
				City:                "Memphis",
				StateOrProvinceCode: "TN",
				PostalCode:          "38125",
				CountryCode:         "US",
			},
		},
		NotificationEmail: "NotificationEmail",
		Reference:         "REF",
		Service:           "FEDEX_GROUND",
		Commodities:       commodities,
	}
	envelope, err := testAPI.processShipmentRequest(shipment)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := envelope.Body.(models.ProcessShipmentBody); !ok {
		t.Fatal("should be process shipment body")
	}

	shipmentBody := envelope.Body.(models.ProcessShipmentBody)
	processShipment := shipmentBody.ProcessShipmentRequest
	if processShipment.Request != expectedRequest {
		t.Fatal("Request doesn't match")
	}

	if processShipment.RequestedShipment.DropoffType != "BUSINESS_SERVICE_CENTER" {
		t.Fatal("DropoffType doesn't match")
	}

	if processShipment.RequestedShipment.ServiceType != "FEDEX_GROUND" {
		t.Fatal("ServiceType doesn't match")
	}

	if processShipment.RequestedShipment.PackagingType != "YOUR_PACKAGING" {
		t.Fatal("PackagingType doesn't match")
	}

	if contact := processShipment.RequestedShipment.Shipper.Contact; contact != shipment.FromContact {
		t.Fatal("contact doesn't match")
	}

	if address := processShipment.RequestedShipment.Shipper.Address; len(address.StreetLines) != 2 ||
		address.StreetLines[0] != "1234 Main Street" ||
		address.StreetLines[1] != "Suite 200" ||
		address.City != "Winnipeg" ||
		address.StateOrProvinceCode != "MB" ||
		address.PostalCode != "R2M4B5" ||
		address.CountryCode != "CA" {
		t.Fatal("address doesn't match")
	}

	if contact := processShipment.RequestedShipment.Recipient.Contact; contact != shipment.ToContact {
		t.Fatal("contact doesn't match")
	}

	if address := processShipment.RequestedShipment.Recipient.Address; len(address.StreetLines) != 2 ||
		address.StreetLines[0] != "3610 Hacks Cross Road" ||
		address.StreetLines[1] != "First Floor" ||
		address.City != "Memphis" ||
		address.StateOrProvinceCode != "TN" ||
		address.PostalCode != "38125" ||
		address.CountryCode != "US" {
		t.Fatal("address doesn't match")
	}

	if processShipment.RequestedShipment.ShippingChargesPayment.PaymentType != "SENDER" {
		t.Fatal("PaymentType doesn't match")
	}

	if specialServicesRequested := processShipment.RequestedShipment.SpecialServicesRequested; len(specialServicesRequested.SpecialServiceTypes) != 2 ||
		specialServicesRequested.SpecialServiceTypes[0] != "ELECTRONIC_TRADE_DOCUMENTS" ||
		specialServicesRequested.SpecialServiceTypes[1] != "EVENT_NOTIFICATION" ||
		specialServicesRequested.EtdDetail.RequestedDocumentCopies != "COMMERCIAL_INVOICE" {
		t.Fatal("specialServicesRequested doesn't match")
	}

	if ls := processShipment.RequestedShipment.LabelSpecification; ls.LabelFormatType != "COMMON2D" ||
		ls.ImageType != "PDF" ||
		*ls.LabelStockType != "PAPER_4X6" {
		t.Fatal("labelSpecification doesn't match")
	}

	if sds := processShipment.RequestedShipment.ShippingDocumentSpecification; len(sds.ShippingDocumentTypes) != 1 ||
		sds.ShippingDocumentTypes[0] != "COMMERCIAL_INVOICE" ||
		len(sds.CommercialInvoiceDetail) != 1 {
		t.Fatal("ShippingDocumentSpecification doesn't match")
	}

	if cid := processShipment.RequestedShipment.ShippingDocumentSpecification.CommercialInvoiceDetail[0]; cid.Format.ImageType != "PDF" ||
		cid.Format.StockType != "PAPER_LETTER" ||
		len(cid.CustomerImageUsages) != 2 ||
		cid.CustomerImageUsages[0].Type != "LETTER_HEAD" ||
		cid.CustomerImageUsages[0].ID != "IMAGE_1" ||
		cid.CustomerImageUsages[1].Type != "SIGNATURE" ||
		cid.CustomerImageUsages[1].ID != "IMAGE_2" {
		t.Fatal("ShippingDocumentSpecification doesn't match")
	}
}
