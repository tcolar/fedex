fedex
=====

Some Fedex API support for GoLang (ATM just for tracking)

Fedex API's are one of those WDSL SOAP monster documented in a gigantic PDF file, don't we all love those.

I did not bother dealing with all of that here and only created what I needed so far.

I might add more over time but for now it provides:
- Retrieving Tracking info by either:
  Tracking number, PO number, or shipper reference number (~order ID)
  The data is unmarshalled from SOAP into Go structures for more practical usage.

See [fedex_example.go](fedex_example.go) for usage examples

Note that you will need an API key and Password as well as Accont and Meter numbers from Fedex.
See: http://images.fedex.com/ca_english/businesstools/webservices/Web_Services_Guide_ENG.pdf