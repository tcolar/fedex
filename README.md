fedex
=====

Some Fedex API support for GoLang (Bare minimum for tracking)

Fedex API's are one of those WDSL SOAP monster documented ina gigantic PDF file, don't we all love those.

I did not bother dealing with that here and only create a few select custom XML packets to get the data I neeed.

I might add more over time but for now it provides:
- Retrieving Tracking info by either:
  Tracking number, PO number, or shipper reference number (~order ID)

Note that you will need an API key and Password as well as Accont and Meter numbers from Fedex.

See [fedex_example.go](fedex_example.go) for uusage examples

Get more info here:
http://images.fedex.com/ca_english/businesstools/webservices/Web_Services_Guide_ENG.pdf