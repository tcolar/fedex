package models

// ServiceType determines the service type (the FedEx API field) based on the
// service (the HR-defined property) and the source/destination. Needed for
// both rates and shipments. Ideally I don't think this should live in models.
func ServiceType(fromAndTo FromAndTo, service string) string {
	// TODO This is confusing. If the service is marked as "fedex_smart_post" or
	// "fedex_international_economy" (this is done through the CMS), then
	// explicitly set the service type as SmartPost or InternationalEconomy
	// respectively. Otherwise, we deduce the service type based on whether
	// the service was "return", whether the return is international, and where
	// the return is coming from. In the future, we should just not allow using
	// anything other than "fedex_smart_post", "fedex_international_economy",
	// "fedex_ground" and make the CMS user to be explicit. However currently
	// there are many shipping methods that depend on this deduction logic.

	isInternational := fromAndTo.IsInternational()
	shipsOutWithInternationalEconomy := fromAndTo.FromAddress.ShipsOutWithInternationalEconomy()
	switch {
	case service == "fedex_smart_post",
		service == "return" && !isInternational:
		return ServiceTypeSmartPost
	case service == "fedex_international_economy" ||
		(isInternational && shipsOutWithInternationalEconomy):
		return ServiceTypeInternationalEconomy
	default:
		return ServiceTypeFedexGround
	}
}
