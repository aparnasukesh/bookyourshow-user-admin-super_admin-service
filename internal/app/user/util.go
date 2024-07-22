package user

func createUserRoles(userId uint, roleId uint) UserRole {

	return UserRole{
		UserID: userId,
		RoleID: roleId,
	}

}
