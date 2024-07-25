package superadmin

func buildAdminRequestsList(admiLists []AdminStatus) []AdminRequests {
	adminRequest := []AdminRequests{}

	for _, admin := range admiLists {
		adminReq := AdminRequests{
			Email: admin.Email,
		}
		adminRequest = append(adminRequest, adminReq)
	}
	return adminRequest
}
