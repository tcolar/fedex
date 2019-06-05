package fedex

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/happyreturns/fedex/models"
)

var (
	testFedex             Fedex
	prodFedex             Fedex
	laSmartPostFedex      Fedex
	blandonSmartPostFedex Fedex
)

func TestMain(m *testing.M) {
	credData, err := ioutil.ReadFile("creds.json")
	if err != nil {
		panic(err)
	}

	creds := map[string]Fedex{}
	if err := json.Unmarshal(credData, &creds); err != nil {
		panic(err)
	}

	testFedex = creds["test"]
	prodFedex = creds["prod"]
	laSmartPostFedex = creds["laSmartPost"]
	blandonSmartPostFedex = creds["blandonSmartPost"]

	os.Exit(m.Run())
}

func TestTrack(t *testing.T) {
	var (
		reply *models.TrackReply
		err   error
	)

	// Error case - invalid tracking
	_, err = testFedex.TrackByNumber(CarrierCodeExpress, "dkfjdkfj")
	checkErrorMatches(t, err, "make track request and unmarshal: response error: track detail error:")

	// Successful case
	reply, err = testFedex.TrackByNumber(CarrierCodeExpress, "123456789012")
	if err != nil {
		t.Fatal(err)
	}

	// Basic validation
	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "trck" ||
		reply.Notifications[0].Code != "0" ||
		reply.Notifications[0].Message != "Request was successfully processed." ||
		reply.Notifications[0].LocalizedMessage != "Request was successfully processed." ||
		reply.Version.ServiceID != "trck" ||
		reply.Version.Major != 16 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		len(reply.CompletedTrackDetails) != 1 ||
		!reply.CompletedTrackDetails[0].DuplicateWaybill ||
		reply.CompletedTrackDetails[0].MoreData ||
		len(reply.CompletedTrackDetails[0].TrackDetails) < 1 ||
		reply.CompletedTrackDetails[0].TrackDetails[0].OperatingCompanyOrCarrierDescription != "FedEx Express" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumber != "123456789012" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].TrackingNumberUniqueIdentifier != "2458162000~123456789012~FX" ||
		reply.CompletedTrackDetails[0].TrackDetails[0].CarrierCode != "FDXE" {
		t.Fatal("output not correct")
	}

	// Check the timestamps make sense
	if len(reply.CompletedTrackDetails[0].TrackDetails[0].DatesOrTimes) == 0 {
		t.Fatal("should get at least one dateortime")
	}
	if reply.ActualDelivery().IsZero() {
		t.Fatal("actual delivery should be set")
	}
	if reply.EstimatedDelivery().IsZero() {
		t.Fatal("actual delivery should be set")
	}

	// Successful case - ground
	reply, err = prodFedex.TrackByNumber(CarrierCodeGround, "786748004585")
	if err != nil {
		t.Fatal(err)
	}
	shipmentTime := reply.Ship()
	if shipmentTime == nil {
		t.Fatal("should get a shipment time")
	}
	s := *shipmentTime
	if s.Hour() == 0 && s.Minute() == 0 && s.Second() == 0 {
		t.Fatal("shipmentTime didn't set time of day")
	}

	// Successful case - smart post
	reply, err = testFedex.TrackByNumber(CarrierCodeSmartPost, "02396343485320033856")
	if err != nil {
		t.Fatal(err)
	}

	// Make sure smart post returned Events
	if len(reply.Events()) == 0 {
		t.Fatal("reply should have an event")
	}

	// Successful case - smart post
	reply, err = testFedex.TrackByNumber(CarrierCodeSmartPost, "02396343484520070272")
	if err != nil {
		t.Fatal(err)
	}

	// Make sure smart post returned Events
	if len(reply.Events()) == 0 {
		t.Fatal("reply should have an event")
	}
}

func TestRate(t *testing.T) {
	// Error case - invalid request
	_, err := prodFedex.Rate(&models.Rate{})
	checkErrorMatches(t, err, "make rate request and unmarshal: response error: reply got error:")

	// Successful case
	reply, err := prodFedex.Rate(
		&models.Rate{
			FromAndTo: models.FromAndTo{
				FromAddress: models.Address{
					StreetLines:         []string{"1517 Lincoln Blvd"},
					City:                "Santa Monica",
					StateOrProvinceCode: "CA",
					PostalCode:          "90401",
					CountryCode:         "US",
				},
				ToAddress: models.Address{
					StreetLines:         []string{"1106 Broadway"},
					City:                "Santa Monica",
					StateOrProvinceCode: "CA",
					PostalCode:          "90401",
					CountryCode:         "US",
				},
				FromContact: models.Contact{
					PersonName:   "Jenny",
					PhoneNumber:  "213 867 5309",
					EmailAddress: "jenny@jenny.com",
				},
				ToContact: models.Contact{
					CompanyName:  "Some Company",
					PhoneNumber:  "214 867 5309",
					EmailAddress: "somecompany@somecompany.com",
				},
			},
		})
	if err != nil {
		t.Fatal(err)
	}

	// Basic validation
	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "crs" ||
		reply.Notifications[0].Code != "0" ||
		reply.Version.ServiceID != "crs" ||
		reply.Version.Major != 24 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		len(reply.RateReplyDetails) != 1 ||
		reply.RateReplyDetails[0].ServiceType != "FEDEX_GROUND" ||
		reply.RateReplyDetails[0].ServiceDescription.ServiceType != "FEDEX_GROUND" ||
		reply.RateReplyDetails[0].ServiceDescription.Code != "92" ||
		reply.RateReplyDetails[0].ServiceDescription.AstraDescription != "FXG" ||
		reply.RateReplyDetails[0].PackagingType != "YOUR_PACKAGING" ||
		reply.RateReplyDetails[0].DestinationAirportID != "LAX" ||
		reply.RateReplyDetails[0].IneligibleForMoneyBackGuarantee ||
		reply.RateReplyDetails[0].SignatureOption != "SERVICE_DEFAULT" ||
		reply.RateReplyDetails[0].ActualRateType != "PAYOR_ACCOUNT_PACKAGE" ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails) != 2 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[0].EffectiveNetDiscount.Amount == 0.0 ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails[0].RatedPackages) != 1 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[0].RatedPackages[0].PackageRateDetail.NetCharge.Amount == 0.0 ||
		len(reply.RateReplyDetails[0].RatedShipmentDetails[1].RatedPackages) != 1 ||
		reply.RateReplyDetails[0].RatedShipmentDetails[1].RatedPackages[0].PackageRateDetail.NetCharge.Amount == 0.0 {
		t.Fatal("output not correct")
	}
	charge, err := reply.TotalCost()
	if err != nil {
		t.Fatal(err)
	}
	if charge.Currency != "USD" || charge.Amount == 0.00 {
		t.Fatal("totalCost should be non-zero, USD")
	}
}

func TestActual(t *testing.T) {
	t.SkipNow()
	myBytes := []byte(`{"FromAddress":{"StreetLines":["1290 Rue Belvédère S",""],"City":"Sherbrooke","StateOrProvinceCode":"QC","PostalCode":"J1H 4C7","CountryCode":"CA","Residential":false},"ToAddress":{"StreetLines":["1106 Broadway",""],"City":"Santa Monica","StateOrProvinceCode":"CA","PostalCode":"90401","CountryCode":"US","Residential":false},"FromContact":{"PersonName":"Jenny","CompanyName":"Jenny","PhoneNumber":"1 (214) 867-5309","EmailAddress":""},"ToContact":{"PersonName":"Happy Returns","CompanyName":"Happy Returns","PhoneNumber":"424 325 9510","EmailAddress":""},"NotificationEmail":"","Reference":"","Service":"return","Commodities":[{"Name":"","NumberOfPieces":0,"Description":"","CountryOfManufacture":"","Weight":{"Units":"LB","Value":2},"Quantity":0,"QuantityUnits":"","UnitPrice":{"Currency":"USD","Amount":128},"CustomsValue":{"Currency":"USD","Amount":128}}]}`)
	shipment := models.Shipment{}
	if err := json.Unmarshal(myBytes, &shipment); err != nil {
		panic(err)
	}

	testShipInternational(t, prodFedex, &shipment)

}

func TestShipGround(t *testing.T) {
	// Error case - invalid shipment
	_, err := prodFedex.Ship(&models.Shipment{})
	checkErrorMatches(t, err, "api process shipment: make process shipment request and unmarshal: response error: reply got error:")

	// Successful case
	exampleShipment := &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1517 Lincoln Blvd"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			ToAddress: models.Address{
				StreetLines:         []string{"1106 Broadway"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			FromContact: models.Contact{
				PersonName:   "Jenny",
				PhoneNumber:  "213 867 5309",
				EmailAddress: "jenny@jenny.com",
			},
			ToContact: models.Contact{
				CompanyName:  "Some Company",
				PhoneNumber:  "214 867 5309",
				EmailAddress: "somecompany@somecompany.com",
			},
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		Reference:         "My ship ground reference",
		Service:           "default",
	}
	reply, err := prodFedex.Ship(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}

	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "ship" ||
		reply.Notifications[0].Code != "0000" ||
		reply.Notifications[0].Message != "Success" ||
		reply.Notifications[0].LocalizedMessage != "Success" ||
		reply.Version.ServiceID != "ship" ||
		reply.Version.Major != 23 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		reply.JobID == "" ||
		reply.CompletedShipmentDetail.UsDomestic != "true" ||
		reply.CompletedShipmentDetail.CarrierCode != "FDXG" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingIdType != "FEDEX" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingNumber == "" ||
		reply.CompletedShipmentDetail.ServiceTypeDescription != "FXG" ||
		reply.CompletedShipmentDetail.ServiceDescription.ServiceType != "FEDEX_GROUND" ||
		reply.CompletedShipmentDetail.ServiceDescription.Code != "92" ||
		// skip ServiceDescription.Names
		reply.CompletedShipmentDetail.PackagingDescription != "YOUR_PACKAGING" ||
		reply.CompletedShipmentDetail.OperationalDetail.OriginLocationNumber != "901" ||
		reply.CompletedShipmentDetail.OperationalDetail.DestinationLocationNumber != "901" ||
		reply.CompletedShipmentDetail.OperationalDetail.TransitTime != "ONE_DAY" ||
		reply.CompletedShipmentDetail.OperationalDetail.IneligibleForMoneyBackGuarantee != "false" ||
		reply.CompletedShipmentDetail.OperationalDetail.DeliveryEligibilities != "SATURDAY_DELIVERY" ||
		reply.CompletedShipmentDetail.OperationalDetail.ServiceCode != "92" ||
		reply.CompletedShipmentDetail.OperationalDetail.PackagingCode != "01" ||
		reply.CompletedShipmentDetail.ShipmentRating.ActualRateType != "PAYOR_ACCOUNT_PACKAGE" ||
		len(reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails) != 1 ||
		// // // skip most ShipmentRateDetails fields
		reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails[0].RateType != "PAYOR_ACCOUNT_PACKAGE" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds[0].TrackingIdType != "FEDEX" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Type != "OUTBOUND_LABEL" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PNG" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts) != 1 ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image) == 0 {
		fmt.Println(reply.CompletedShipmentDetail.ShipmentRating.ShipmentRateDetails)
		t.Fatal("output not correct")
	}

	// Decode png bytes from base64 data
	pngBytes, err := base64.StdEncoding.DecodeString(string(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image))
	if err != nil {
		t.Fatal(err)
	}

	// Write label as png, and manually check it
	err = ioutil.WriteFile(fmt.Sprintf("output-ground-not-international-%s.png", prodFedex.API.Key), pngBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// it also works with no email
	exampleShipment.NotificationEmail = ""
	reply, err = prodFedex.Ship(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
		fmt.Println(reply)
		t.Fatal("reply should not have failed")
	}
}

func TestShipSmartPost(t *testing.T) {
	// smartpost fail for international
	internationalShipment := &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1234 Main Street", "Suite 200"},
				City:                "Winnipeg",
				StateOrProvinceCode: "MB",
				PostalCode:          "R2M4B5",
				CountryCode:         "CA",
			},
			ToAddress: models.Address{
				StreetLines:         []string{"1106 Broadway"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			FromContact: models.Contact{
				PersonName:   "Jenny",
				PhoneNumber:  "213 867 5309",
				EmailAddress: "jenny@jenny.com",
			},
			ToContact: models.Contact{
				CompanyName:  "normal",
				PhoneNumber:  "214 867 5309",
				EmailAddress: "somecompany@somecompany.com",
			},
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		Reference:         "My ship ground reference",
		Commodities:       []models.Commodity{},
	}
	_, err := laSmartPostFedex.Ship(internationalShipment)
	checkErrorMatches(t, err, "do not ship internationally with smartpost")

	// Successful cases
	testShipSmartPostSuccess(t, laSmartPostFedex)
	testShipSmartPostSuccess(t, blandonSmartPostFedex)
}

func TestShipInternational(t *testing.T) {
	var err error
	fedex := testFedex

	// Successful case
	exampleShipment := &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1234 Main Street", "Suite 200"},
				City:                "Winnipeg",
				StateOrProvinceCode: "MB",
				PostalCode:          "R2M4B5",
				CountryCode:         "CA",
			},
			ToAddress: models.Address{
				StreetLines:         []string{"1106 Broadway"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			FromContact: models.Contact{
				PersonName:   "Jenny",
				PhoneNumber:  "213 867 5309",
				EmailAddress: "jenny@jenny.com",
			},
			ToContact: models.Contact{
				CompanyName:  "normal",
				PhoneNumber:  "214 867 5309",
				EmailAddress: "somecompany@somecompany.com",
			},
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		Reference:         "My ship ground reference",
		Commodities: []models.Commodity{
			{
				NumberOfPieces:       1,
				Description:          "Computer Keyboard",
				Quantity:             1,
				QuantityUnits:        "unit",
				CountryOfManufacture: "US",
				Weight:               models.Weight{Units: "LB", Value: 10.0},
				UnitPrice:            &models.Money{Currency: "USD", Amount: 25.00},
				CustomsValue:         &models.Money{Currency: "USD", Amount: 30.00},
			},
			{
				NumberOfPieces:       1,
				Description:          "Computer Monitor",
				Quantity:             1,
				QuantityUnits:        "unit",
				CountryOfManufacture: "US",
				Weight:               models.Weight{Units: "LB", Value: 5.0},
				UnitPrice:            &models.Money{Currency: "USD", Amount: 214.42},
				CustomsValue:         &models.Money{Currency: "USD", Amount: 381.12},
			},
		},
	}

	fmt.Println(fedex)
	exampleShipment.ToContact.CompanyName = "dev"
	testShipInternational(t, testFedex, exampleShipment)

	exampleShipment.ToContact.CompanyName = "normal"
	testShipInternational(t, fedex, exampleShipment)

	// it also works with no email
	fmt.Println("No email")
	exampleShipment.NotificationEmail = ""
	exampleShipment.ToContact.CompanyName = "no-email"
	testShipInternational(t, fedex, exampleShipment)

	// it also works when an importer name is supplied
	fmt.Println("Has importer name")
	exampleShipment.Importer = "Rothy's"
	exampleShipment.ToContact.CompanyName = "has-importer-name"
	testShipInternational(t, prodFedex, exampleShipment)

	// it also works when an importer name and an importer address are supplied
	fmt.Println("Has importer name and address")
	exampleShipment.Importer = "Rothy's"
	exampleShipment.ImporterAddress = models.Address{
		StreetLines:         []string{"1511 15th Street"},
		City:                "Santa Monica",
		StateOrProvinceCode: "CA",
		PostalCode:          "90404",
		CountryCode:         "US",
	}
	exampleShipment.ToContact.CompanyName = "has-importer-name-and-address"
	testShipInternational(t, prodFedex, exampleShipment)

	// it also works when we supply a letterhead image id
	fmt.Println("Has letterhead image id")
	exampleShipment.LetterheadImageID = "IMAGE_3"
	exampleShipment.ToContact.CompanyName = "has-letterhead-image-id"
	testShipInternational(t, prodFedex, exampleShipment)

	// it also works when commodities > 800
	exampleShipment.Commodities = append(exampleShipment.Commodities,
		models.Commodity{
			NumberOfPieces:       1,
			Description:          "Computer",
			Quantity:             1,
			QuantityUnits:        "unit",
			CountryOfManufacture: "US",
			Weight:               models.Weight{Units: "LB", Value: 50.0},
			UnitPrice:            &models.Money{Currency: "USD", Amount: 1214.42},
			CustomsValue:         &models.Money{Currency: "USD", Amount: 1381.12},
		},
	)
	fmt.Println(fedex)
	exampleShipment.ToContact.CompanyName = "more-commodities"
	testShipInternational(t, prodFedex, exampleShipment)

	// test commodities with no unitprice
	for _, commodity := range exampleShipment.Commodities {
		commodity.UnitPrice = nil
	}
	exampleShipment.ToContact.CompanyName = "commodities-no-unit-price"
	testShipInternational(t, prodFedex, exampleShipment)

	// test commodities with no customsvalue
	for _, commodity := range exampleShipment.Commodities {
		commodity.CustomsValue = nil
	}
	exampleShipment.ToContact.CompanyName = "commodities-no-customs-value"
	testShipInternational(t, prodFedex, exampleShipment)

	// it tries to make a request with no commodities and returns back the error
	// fedex gives us
	// Error case - invalid tracking number
	exampleShipment.Commodities = nil
	exampleShipment.ToContact.CompanyName = "no-commodities"
	_, err = prodFedex.Ship(exampleShipment)
	checkErrorMatches(t, err, "api process shipment: make process shipment request and unmarshal: response error:")
}

func testShipInternational(t *testing.T, f Fedex, shipment *models.Shipment) {
	reply, err := f.Ship(shipment)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "NOTE" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Type != "OUTBOUND_LABEL" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PDF" ||
		reply.CompletedShipmentDetail.ShipmentDocuments[0].Type != "COMMERCIAL_INVOICE" {
		t.Fatal("shipment international output not correct")
	}

	// Save label, manually check it
	data, err := base64.StdEncoding.DecodeString(string(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image))
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("output-international-label-%s-%s.pdf", shipment.ToContact.CompanyName, f.API.Key), data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Save commercial invoice, manually check it
	data, err = base64.StdEncoding.DecodeString(string(reply.CompletedShipmentDetail.ShipmentDocuments[0].Parts[0].Image))
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("output-international-invoice-%s-%s.pdf", shipment.ToContact.CompanyName, f.API.Key), data, 0644)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCreatePickup(t *testing.T) {
	t.SkipNow()
	reply, err := prodFedex.CreatePickup(
		&models.Pickup{
			PickupLocation: models.PickupLocation{
				Address: models.Address{
					StreetLines:         []string{"1517 Lincoln Blvd"},
					City:                "Santa Monica",
					StateOrProvinceCode: "CA",
					PostalCode:          "90401",
					CountryCode:         "US",
				},
				Contact: models.Contact{
					PersonName:   "Jenny",
					PhoneNumber:  "213 867 5309",
					EmailAddress: "jenny@jenny.com",
				}},
			ToAddress: models.Address{
				StreetLines:         []string{"7631 HAskell Ave"},
				City:                "Van Nuys",
				StateOrProvinceCode: "CA",
				PostalCode:          "94106",
				CountryCode:         "US",
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}

	if reply.HighestSeverity != "SUCCESS" ||
		// Basic validation
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "disp" ||
		reply.Notifications[0].Code != "0000" ||
		reply.Notifications[0].Message != "Success" ||
		reply.Version.ServiceID != "disp" ||
		reply.Version.Major != 17 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		reply.PickupConfirmationNumber == "" ||
		reply.PickupConfirmationNumber == "0" {
		t.Fatal("output not correct")
	}
}

func TestSendNotifications(t *testing.T) {
	// Error case - invalid tracking number
	_, err := prodFedex.SendNotifications("123", "dev-notifications@happyreturns.com")
	checkErrorMatches(t, err, "make send notifications request: response error: reply got error: Invalid tracking numbers.")

	// Successful case
	reply, err := prodFedex.SendNotifications("02396343485320152281", "dev-notifications@happyreturns.com")
	if err != nil {
		t.Fatal("error did not match", err)
	}
	if reply.Error() != nil {
		t.Fatal("reply should not have failed")
	}

	// Basic validation
	if reply.HighestSeverity != "SUCCESS" ||
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "trck" ||
		reply.Notifications[0].Code != "0" ||
		reply.Notifications[0].Message != "Request was successfully processed." ||
		reply.Notifications[0].LocalizedMessage != "Request was successfully processed." ||
		reply.DuplicateWaybill ||
		reply.MoreDataAvailable ||
		len(reply.Packages) != 1 ||
		reply.Packages[0].TrackingNumber != "02396343485320152281" ||
		len(reply.Packages[0].TrackingNumberUniqueIdentifiers) != 1 ||
		reply.Packages[0].CarrierCode != "FXSP" ||
		reply.Packages[0].ShipDate != "2019-03-11" ||
		reply.Packages[0].Destination.City != "BLANDON" ||
		reply.Packages[0].Destination.StateOrProvinceCode != "PA" ||
		reply.Packages[0].Destination.CountryCode != "US" ||
		bool(reply.Packages[0].Destination.Residential) ||
		len(reply.Packages[0].RecipientDetails) != 1 ||
		len(reply.Packages[0].RecipientDetails[0].NotificationEventsAvailable) != 1 {
		t.Fatal("output not correct")
	}
}

func testShipSmartPostSuccess(t *testing.T, fedexAccount Fedex) {
	exampleShipment := &models.Shipment{
		FromAndTo: models.FromAndTo{
			FromAddress: models.Address{
				StreetLines:         []string{"1517 Lincoln Blvd"},
				City:                "Santa Monica",
				StateOrProvinceCode: "CA",
				PostalCode:          "90401",
				CountryCode:         "US",
			},
			ToAddress: models.Address{},
			FromContact: models.Contact{
				PersonName:   "Jenny",
				PhoneNumber:  "213 867 5309",
				EmailAddress: "jenny@jenny.com",
			},
			ToContact: models.Contact{
				CompanyName:  "Some Company",
				PhoneNumber:  "214 867 5309",
				EmailAddress: "somecompany@somecompany.com",
			},
		},
		NotificationEmail: "dev-notifications@happyreturns.com",
		Reference:         "My reference",
		Service:           "return",
	}
	reply, err := fedexAccount.Ship(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	// Basic validation
	if reply.Error() != nil {
		fmt.Println(reply)
		t.Fatal("reply should not have failed")
	}
	if reply.HighestSeverity != "SUCCESS" ||
		len(reply.Notifications) != 1 ||
		reply.Notifications[0].Source != "ship" ||
		reply.Notifications[0].Code != "0000" ||
		reply.Notifications[0].Message != "Success" ||
		reply.Notifications[0].LocalizedMessage != "Success" ||
		reply.Version.ServiceID != "ship" ||
		reply.Version.Major != 23 ||
		reply.Version.Intermediate != 0 ||
		reply.Version.Minor != 0 ||
		reply.JobID == "" ||
		reply.CompletedShipmentDetail.UsDomestic != "true" ||
		reply.CompletedShipmentDetail.CarrierCode != "FXSP" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingIdType != "USPS" ||
		reply.CompletedShipmentDetail.MasterTrackingId.TrackingNumber == "" ||
		reply.CompletedShipmentDetail.ServiceTypeDescription != "SMART POST" ||
		reply.CompletedShipmentDetail.ServiceDescription.ServiceType != "SMART_POST" ||
		reply.CompletedShipmentDetail.PackagingDescription != "YOUR_PACKAGING" ||
		reply.CompletedShipmentDetail.OperationalDetail.TransitTime != "TWO_DAYS" ||
		reply.CompletedShipmentDetail.OperationalDetail.IneligibleForMoneyBackGuarantee != "false" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds) != 1 ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.TrackingIds[0].TrackingIdType != "USPS" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Type != "OUTBOUND_LABEL" ||
		reply.CompletedShipmentDetail.CompletedPackageDetails.Label.ImageType != "PNG" ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts) != 1 ||
		len(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image) == 0 {
		t.Fatal("output not correct")
	}

	// Decode pdf bytes from base64 data
	pngBytes, err := base64.StdEncoding.DecodeString(string(reply.CompletedShipmentDetail.CompletedPackageDetails.Label.Parts[0].Image))
	if err != nil {
		t.Fatal(err)
	}

	// Write label as png, and manually check it
	err = ioutil.WriteFile(fmt.Sprintf("output-smart-post-%s.png", fedexAccount.API.Key), pngBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// It also works with no email
	exampleShipment.NotificationEmail = ""
	reply, err = fedexAccount.Ship(exampleShipment)
	if err != nil {
		t.Fatal(err)
	}

	if reply.Error() != nil {
		fmt.Println(reply)
		t.Fatal("reply should not have failed")
	}
}

func checkErrorMatches(t *testing.T, err error, expectedText string) {
	if err == nil || !strings.HasPrefix(err.Error(), expectedText) {
		t.Fatal("error", err, "doesn't match", expectedText)
	}
}
