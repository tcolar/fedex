package fedex

import (
	"testing"

	"github.com/happyreturns/fedex/models"
	. "github.com/onsi/gomega"
)

func TestFedexRate(t *testing.T) {
	t.Run("heavier-packages-are-more-expensive", func(t *testing.T) {
		// set up test cases
		type testCase struct {
			name    string
			fedex   Fedex
			getRate func() *models.Rate
		}

		// try rates with different services (smartpost, ground) and
		// destinations (international, domestic)
		testCases := []testCase{
			{
				name:    "international-ground",
				fedex:   testFedex,
				getRate: exampleInternationalRate,
			},
			{
				name:    "domestic-ground-test",
				fedex:   testFedex,
				getRate: exampleDomesticRate,
			},
			{
				name:    "domestic-ground-prod",
				fedex:   prodFedex,
				getRate: exampleDomesticRate,
			},
			{
				name:    "domestic-smartpost-la",
				fedex:   laSmartPostFedex,
				getRate: exampleDomesticRate,
			},
			{
				name:    "domestic-smartpost-blandon",
				fedex:   blandonSmartPostFedex,
				getRate: exampleDomesticRate,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				g := NewWithT(t)

				// Get rate with light weight
				lightRequest := testCase.getRate()
				lightReply, err := testCase.fedex.Rate(lightRequest)
				g.Expect(err).NotTo(HaveOccurred())

				// Get rate with much heavier weight (10 times heavier)
				heavyRequest := testCase.getRate()
				for idx := range heavyRequest.Commodities {
					heavyRequest.Commodities[idx].Weight.Value *= 10
				}
				heavyReply, err := testCase.fedex.Rate(heavyRequest)
				g.Expect(err).NotTo(HaveOccurred())

				// Verify rate for the heavier request is much greater (at least 2 times
				// more expensive) than the rate for the light request
				lightCost, err := lightReply.TotalCost()
				g.Expect(err).NotTo(HaveOccurred())
				heavyCost, err := heavyReply.TotalCost()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(lightCost.Amount).To(BeNumerically(">", 0))
				g.Expect(heavyCost.Amount).To(BeNumerically(">", 0))
				g.Expect(heavyCost.Amount).To(BeNumerically(">", lightCost.Amount*5.0))

			})
		}
	})
}

func exampleInternationalRate() *models.Rate {
	rate := exampleDomesticRate()
	rate.FromAndTo.FromAddress = models.Address{
		StreetLines:         []string{"1234 Main Street", "Suite 200"},
		City:                "Winnipeg",
		StateOrProvinceCode: "MB",
		PostalCode:          "R2M4B5",
		CountryCode:         "CA",
	}
	return rate
}

func exampleDomesticRate() *models.Rate {
	return &models.Rate{
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
		},
		Commodities: []models.Commodity{
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
		},
	}
}
