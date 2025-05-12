package sumup

// Transform saumup transaction user to warehouse code
func TransformUserToWarehouseCode(user string) string {
	switch user {
	case "terminal1@maximal-limit.cz":
		return "code:MAXIMAL-LID"
	case "terminal2@maximal-limit.cz":
		return "code:MAXIMAL-DL"
	case "terminal3@maximal-limit.cz":
		return "code:MAXIMAL-BS"
	case "CENTRAL":
		return "code:SKLAD"
	default:
		return ""
	}
}
