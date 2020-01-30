package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/happyreturns/fedex/models"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Entry
)

func init() {
	logger = logrus.WithFields(logrus.Fields{
		"app": "fedex",
	})
}

func (a API) makeRequestAndUnmarshalResponse(url string, request *models.Envelope,
	response models.Response) error {
	// Create request body
	reqXML, err := xml.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request xml: %s", err)
	}

	// Post XML
	content, err := postXML(a.FedExURL+url, string(reqXML))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"url":     url,
			"request": string(reqXML),
			"err":     err,
		}).Error("error-posting-xml")
		return fmt.Errorf("post xml: %s", err)
	}

	// Parse response
	if err := xml.Unmarshal(content, response); err != nil {
		logger.WithFields(logrus.Fields{
			"url":      url,
			"request":  string(reqXML),
			"response": string(content),
			"err":      err,
		}).Error("error-parsing-xml")
		return fmt.Errorf("parse xml: %s", err)
	}

	// Check if reply failed (FedEx responds with 200 even though it failed)
	if err := response.Error(); err != nil {
		// Frequently, an error is thrown because a tracking number could not be found.
		// This is business as usual for us, often we are looking up tracking numbers before
		// they are queriable on the shipping service.
		// In this case
		//   --> we DO NOT log an error, it is not an error from the logging perspective
		//   --> this is still considered an error from the code-level perspective,
		//       so we still return the error
		if false == err.Error().Contains("This tracking number cannot be found") {
			logger.WithFields(logrus.Fields{
				"url":      url,
				"request":  string(reqXML),
				"response": string(content),
				"err":      err,
			}).Error("error-response")
		}

		// return the error, even if we didn't log it
		return fmt.Errorf("response error: %s", err)
	}

	return nil
}

// postXML to Fedex API and return response
func postXML(url, xml string) ([]byte, error) {
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all bytes: %s", err)
	}
	return content, nil
}
