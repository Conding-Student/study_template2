package routers

import (
	admincardinc "chatbot/pkg/admin_cardinc"
	"chatbot/pkg/admin_cardinc/dashboard"
	loanreversal "chatbot/pkg/admin_cardinc/loans_reversal"
	offices "chatbot/pkg/admin_cardinc/office_management"
	adminmlni "chatbot/pkg/admin_mlni"
	administrator "chatbot/pkg/admin_super"
	switches "chatbot/pkg/admin_super/server_switch"
	triviafacts "chatbot/pkg/admin_super/trivia_facts"
	usermanagement "chatbot/pkg/admin_super/user_management"
	"chatbot/pkg/authentication"
	"chatbot/pkg/controllers/healthchecks"
	"chatbot/pkg/eloading"
	"chatbot/pkg/empc"
	"chatbot/pkg/esystem"
	"chatbot/pkg/features"
	"chatbot/pkg/gabaykonek/audittrail"
	"chatbot/pkg/gabaykonek/creditline"
	gabaykonekdashboard "chatbot/pkg/gabaykonek/dashboard"
	loans "chatbot/pkg/gabaykonek/loanapplication"
	"chatbot/pkg/gabaykonek/qslsal"
	"chatbot/pkg/gabaykonek/reports"
	"chatbot/pkg/handler"
	"chatbot/pkg/hcis"
	"chatbot/pkg/loancalc"
	"chatbot/pkg/logs"
	"chatbot/pkg/models/response"
	users "chatbot/pkg/user"
	"chatbot/pkg/utils/go-utils/encryptDecrypt"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/cagabay")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1", handler.ServerSwitchMain)
	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)
	v1Endpoint.Post("/successful/app/launched", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(response.ResponseModel{
			RetCode: "200",
			Message: "Successful! You are now connected to CA-GABAY Server",
			Data:    nil,
		})
	})
	v1Endpoint.Post("/educationalAttainment", handler.AuthMiddleware, handler.EducationalAttainments)
	v1Endpoint.Post("/genders", handler.AuthMiddleware, handler.Genders)
	v1Endpoint.Post("/institutions", handler.AuthMiddleware, handler.Institutions)
	v1Endpoint.Post("/institutionsDateEstablished", handler.AuthMiddleware, handler.InstitutionsAndEstablished)
	v1Endpoint.Post("/maritalStatus", handler.AuthMiddleware, handler.MaritalStats)
	v1Endpoint.Post("/religion", handler.AuthMiddleware, handler.Religions)
	v1Endpoint.Post("/occupation", handler.AuthMiddleware, handler.Occupations)
	v1Endpoint.Post("/legalId", handler.AuthMiddleware, handler.LegalIDs)
	v1Endpoint.Post("/monthlyIncome", handler.AuthMiddleware, handler.MonthlyIncome)
	v1Endpoint.Post("/natureOfBusiness", handler.AuthMiddleware, handler.NatureOfBusiness)

	authEndpoint := v1Endpoint.Group("/auth", handler.AuthMiddleware)
	authEndpoint.Get("/validatetoken", authentication.ValidateSession)

	utilitiesEndpoint := v1Endpoint.Group("/utils")
	utilitiesEndpoint.Post("/encrypt", encryptDecrypt.EncryptHandler)
	utilitiesEndpoint.Post("/decrypt", encryptDecrypt.DecryptHandler)

	featuresEndpoint := v1Endpoint.Group("/features", authentication.ValidateUserToken, handler.AuthMiddleware)
	featuresEndpoint.Post("/featuresAndVersions/:id", features.FeaturesAndVersions)
	featuresEndpoint.Post("/secondaryFeatures/:id", features.SecondaryFeatures)

	hcisEndpoint := v1Endpoint.Group("/hcis", handler.AuthMiddleware)
	hcisEndpoint.Post("/staffInfo/:id", handler.ServerSwitchRegistration, hcis.HcisStaffInfo)

	empcEndpoint := v1Endpoint.Group("/empc", handler.ServerSwitchEMPCSOA, authentication.ValidateUserToken, handler.AuthMiddleware)
	empcEndpoint.Post("/getEmpcData/:id", empc.EmpcStaffInfo)
	empcEndpoint.Post("/getEmpcAmortization/:id", empc.GetEMPCAmortization)
	empcEndpoint.Post("/checkLoanBalance/:id", empc.CheckLoanBalance)
	empcEndpoint.Get("/getEmpcLoanProducts/:id", empc.GetEmpcLoanProducts)
	empcEndpoint.Post("/getStatementOfAccount/:id", empc.ViewStatementOfAccount)

	//old endpoint to create loan
	v1Endpoint.Post("/loan/addLoanApplication", loans.InsertLoanApplication)

	kplusEndpoint := v1Endpoint.Group("/kplus", handler.AuthMiddleware)
	kplusEndpoint.Get("/viewbranches/:id", offices.ViewCardIncBranches)
	kplusEndpoint.Get("/viewCardIncBranches", offices.ViewCardIncBranches)
	kplusEndpoint.Post("/loancreation", loans.InsertLoanApplication)
	kplusEndpoint.Post("/loanCalculatorPlus", loans.LoanCalculatorPlus)
	kplusEndpoint.Post("/rbiCalculatorPlus/:id", loancalc.RBILoanCalculator)

	esystemEndpoint := kplusEndpoint.Group("/esystem")
	esystemEndpoint.Post("/getClientInformation/:id", esystem.GetClientInformationInEsystem)
	esystemEndpoint.Post("/getCoBorrowerInformation/:id", authentication.ValidateUserToken, esystem.GetCoBorrowerInformationInEsystem)
	esystemEndpoint.Post("/getClientLoanDetails", esystem.GetClientLoanDetails)
	esystemEndpoint.Post("getClientCurrentLoans", esystem.GetClientCurrentLoans)
	esystemEndpoint.Post("/getClientSavingsBalance", esystem.GetClientSavingsBalance)
	esystemEndpoint.Post("/getClientCurrentLoansAndSavingsBalance/:id", esystem.GetClientCurrentLoansAndSavingsBalance)

	loanCalculatorEndpoint := v1Endpoint.Group("/loancalculator", handler.ServerSwitchLoanCalculator, authentication.ValidateUserToken)
	loanCalculatorEndpoint.Get("/calcLoanProducts/:id", loans.LoanProducts)
	loanCalculatorEndpoint.Post("/computeloan/:id", loancalc.BankLoanCalculator)

	gabayKonekEndpoint := v1Endpoint.Group("/gabaykonek", handler.AuthMiddleware, handler.ServerSwitchGabayKonek)

	loanEndpoint := gabayKonekEndpoint.Group("/loan")
	loanEndpoint.Get("/loanProductListAndDetails", authentication.ValidateUserToken, loans.LoanProductListAndDetails)
	loanEndpoint.Get("/midasDetails/:id", authentication.ValidateUserToken, loans.ViewMidasDetails)
	loanEndpoint.Get("/loanStatusAndRoles/:id", authentication.ValidateUserToken, loans.LoanStatusAndRoles)
	loanEndpoint.Get("/getGradeLevel/:id", authentication.ValidateUserToken, loans.GetGradeLevel)
	loanEndpoint.Get("/getRelationships/:id", authentication.ValidateUserToken, loans.GetCoBorrowerRelationships)
	loanEndpoint.Post("/loanCalculatorPlus", authentication.ValidateUserToken, loans.LoanCalculatorPlus)
	loanEndpoint.Post("/getLoanApplication/:id", authentication.ValidateUserToken /*,loans.AutoCancellationOfPendingApprovedLoans,*/, loans.GetLoanApplications)
	loanEndpoint.Post("/updateLoanApplication/:id", authentication.ValidateUserToken, loans.UpdateLoanApplication)
	loanEndpoint.Get("/getppiquestionaire/:id", authentication.ValidateUserToken, loans.GetPPIQuestionaire)
	loanEndpoint.Post("/getloansperclient/:id", authentication.ValidateUserToken, loans.GetLoansPerClient)
	loanEndpoint.Get("/loanCategory/:id", loans.GetLoanCategory)     // no params
	loanEndpoint.Get("/getLoanPurposes/:id", loans.GetLoanPurpose)   // no params
	loanEndpoint.Get("/getBusinessTypes/:id", loans.GetBusinessType) // no params
	loanEndpoint.Post("/loanPurposes/:id", loans.LoanPurpose)
	loanEndpoint.Post("/businessTypes/:id", loans.BusinessType)
	loanEndpoint.Post("/getLoanStatus", loans.GetLoanStatus)

	qslsalEndpoint := gabayKonekEndpoint.Group("/qslsal", authentication.ValidateUserToken)
	qslsalEndpoint.Get("/qslsalfields/:id", qslsal.GetFields)

	creditLineEndpoint := gabayKonekEndpoint.Group("/creditline", authentication.ValidateUserToken)
	creditLineEndpoint.Get("/fields/:id", creditline.GetCreditLineFields)
	creditLineEndpoint.Get("/properties/:id", creditline.GetCreditLineProperties)
	creditLineEndpoint.Post("/creation/:id", creditline.CreditLineCreation)
	creditLineEndpoint.Post("/getlist/:id", creditline.GetCreditLineList)
	creditLineEndpoint.Post("/approvedcreditline/:id", creditline.GetApprovedCreditLine)

	dashboardEndpoint := gabayKonekEndpoint.Group("/dashboard", authentication.ValidateUserToken)
	dashboardEndpoint.Get("/dashBoardData/:id", gabaykonekdashboard.LOSDashboard)
	dashboardEndpoint.Post("/loansAndTotals/:id", gabaykonekdashboard.LoansAndTotals)
	dashboardEndpoint.Post("/efficiency/:id", gabaykonekdashboard.GetEfficiency)
	dashboardEndpoint.Post("/staffoffice/:id", gabaykonekdashboard.GetOffices)

	reportsEndpoint := gabayKonekEndpoint.Group("/reports", authentication.ValidateUserToken, handler.AuthMiddleware, handler.ServerSwitchGabayKonek)
	reportsEndpoint.Post("/getLoanApplied/:id", reports.GetAppliedLoanSummary)
	reportsEndpoint.Post("/getLoanRecommended/:id", reports.GetRecommendedLoanSummary)
	reportsEndpoint.Post("/getLoanApproved/:id", reports.GetApprovedLoanSummary)
	reportsEndpoint.Post("/getLoanPending/:id", reports.GetPendingLoanSummary)
	reportsEndpoint.Post("/getLoanReleased/:id", reports.GetReleasedLoanSummary)
	reportsEndpoint.Post("/getLoanCancelled/:id", reports.GetCancelledLoanSummary)

	eloadingEndpoint := v1Endpoint.Group("/eloading", authentication.ValidateUserToken, handler.AuthMiddleware)
	eloadingEndpoint.Post("/eloadrequest/:id", eloading.EloadLoadRequest)

	usersEndpoint := v1Endpoint.Group("/user")
	usersEndpoint.Post("/accountCreation/:id", handler.AuthMiddleware, users.AccountCreation)
	usersEndpoint.Post("/accountLogin/:id", handler.AuthMiddleware, users.AccountLogin)
	usersEndpoint.Post("/login/:id/:deviceid", handler.AuthMiddleware, users.Login)
	usersEndpoint.Post("/viewprofile/:id", authentication.ValidateUserToken, handler.AuthMiddleware, users.ViewProfile)
	usersEndpoint.Post("/updateUsers/:id", handler.AuthMiddleware, users.UpdateUser)
	usersEndpoint.Post("/updatedevice/:id/", handler.AuthMiddleware, users.UpdateDevice)
	usersEndpoint.Post("/changePassword/:id", authentication.ValidateUserToken, handler.AuthMiddleware, users.ChangePassword)
	usersEndpoint.Post("/resetPasswordLink/:id", handler.AuthMiddleware, users.PasswordResetViaLink)
	usersEndpoint.Post("/resetPassword/:id", handler.AuthMiddleware, users.PasswordReset)
	usersEndpoint.Post("/clientInformationForm", handler.AuthMiddleware, handler.ClientInformationForm)
	usersEndpoint.Post("/updateProfilePicture/:id", handler.AuthMiddleware, users.UpdateUserProfilePicture)
	usersEndpoint.Post("/chat/:id", handler.ServerSwitchChatbot, handler.AuthMiddleware, handler.ChatHandler)
	usersEndpoint.Post("/createWishList/:id", users.CreateWishList)
	usersEndpoint.Post("/createpin/:id", authentication.ValidateUserToken, handler.AuthMiddleware, users.CreatePin)
	// usersEndpoint.Post("/deleteAccount")

	usersEndpoint.Post("/newUser/:id/:deviceid", users.AccountCreation)
	usersEndpoint.Post("/logout/:id", users.Logout)

	//--------------------------------------------------------------------------------------------------//

}

func SetupPublicRoutesB(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealthB)
	adminEnpoint := v1Endpoint.Group("/admin", handler.AuthMiddleware)
	adminEnpoint.Post("/login/:id", administrator.AdminLogin)
	adminEnpoint.Post("/logout/:id", users.Logout)

	adminEnpoint.Post("/institutionslist", authentication.ValidateAdminToken, handler.Institutions)
	adminEnpoint.Get("/getallusers/:id", authentication.ValidateSuperAdminToken, usermanagement.GetAllUsers)
	adminEnpoint.Post("accountCreationAD/:id", authentication.ValidateAdminToken, usermanagement.AccountCreationAdmin)
	adminEnpoint.Post("/updateusers/:id", authentication.ValidateAdminToken, usermanagement.UpdateUsers)
	adminEnpoint.Post("syncuserdata/:id", authentication.ValidateAdminToken, usermanagement.SyncUserData)

	// With superadmin token validation
	adminEnpoint.Get("/serverSwitch", authentication.ValidateSuperAdminToken, switches.GetSwitch)
	adminEnpoint.Post("/updateSwitch", authentication.ValidateSuperAdminToken, switches.UpdateSwitch)
	adminEnpoint.Get("/getTrivia", authentication.ValidateSuperAdminToken, triviafacts.GetTrivia)
	adminEnpoint.Post("/updateTrivia", authentication.ValidateSuperAdminToken, triviafacts.UpdateTrivia)
	adminEnpoint.Get("/getArticles", authentication.ValidateSuperAdminToken, triviafacts.GetArticles)
	adminEnpoint.Post("/updateArticles", authentication.ValidateSuperAdminToken, triviafacts.UpdateArticles)
	adminEnpoint.Post("/getLogs", authentication.ValidateSuperAdminToken, logs.GetLogs)
	adminEnpoint.Get("/getWishLists", authentication.ValidateSuperAdminToken, administrator.GetWishList)
	adminEnpoint.Get("/getInstitutionAndClientCount", authentication.ValidateSuperAdminToken, administrator.GetInstiAndClientCount)
	adminEnpoint.Put("/updateUsers/:id", authentication.ValidateSuperAdminToken, users.UpdateUser)

	//old get branches endpoint used in kplus
	adminEnpoint.Post("/viewCardIncBranches", offices.ViewCardIncBranches)

	// CARD, Inc. API's
	cardIncEnpoint := adminEnpoint.Group("/cardinc", authentication.ValidateAdminToken)
	// roles
	cardIncEnpoint.Post("/addRoles/:id", admincardinc.AddUserRole)
	cardIncEnpoint.Post("/viewRoles/:id", admincardinc.ViewUserRole)
	cardIncEnpoint.Delete("/deleteRoles/:id", admincardinc.DeleteuserRoles)
	cardIncEnpoint.Post("/viewApprovingAuthority", admincardinc.ViewApprovingAuthority)
	cardIncEnpoint.Post("/approvingAuthority", admincardinc.AddApprovingAuthority)
	cardIncEnpoint.Delete("/deleteApprovingAuthority", admincardinc.DeleteApprovingAuthority)
	// view offices
	cardIncEnpoint.Get("/clusters/:id", offices.GetClusters)
	cardIncEnpoint.Post("/regions/:id", offices.GetRegions)
	cardIncEnpoint.Post("/branches/:id", offices.GetBranches)
	cardIncEnpoint.Post("/units/:id", offices.GetUnits)
	cardIncEnpoint.Post("/centers/:id", offices.GetCenters)
	cardIncEnpoint.Post("/viewCardIncUnits", offices.ViewCardIncUnits)
	// manage offices
	cardIncEnpoint.Post("/upsertCluster/:id", offices.UpsertCluster)
	cardIncEnpoint.Post("/upsertRegion/:id", offices.UpsertRegion)
	cardIncEnpoint.Post("/upsertBranch/:id", offices.UpsertBranches)
	cardIncEnpoint.Post("/upsertUnit/:id", offices.UpsertUnits)
	cardIncEnpoint.Post("/upsertCenter/:id", offices.UpsertCenters)
	// users
	cardIncEnpoint.Get("/getcardincusers/:id", admincardinc.GetCardIncUsers)
	// dashboard
	cardIncEnpoint.Get("/dashboardData/:id", dashboard.GetDashBoardData)
	// logs
	cardIncEnpoint.Post("/auditlogs/:id", audittrail.AccessLogs)
	// loan
	cardIncEnpoint.Post("/viewdisbursedloans/:id", admincardinc.GetLoanReleased)
	cardIncEnpoint.Post("/viewloans/:id", loanreversal.ViewLoans)
	cardIncEnpoint.Post("/loanreversal/:id", loanreversal.ReverseLoan)
	// ppi
	cardIncEnpoint.Get("/getppiquestionaire/:id", loans.GetPPIQuestionaire)
	// creditline
	cardIncEnpoint.Post("/approvedcreditline/:id", creditline.GetApprovedCreditLine)
	// unknown
	cardIncEnpoint.Post("/getname/:id", offices.GetStaffName)
	cardIncEnpoint.Post("/getStaffByDesignation", offices.GetStaffByDesignation)
	cardIncEnpoint.Post("/getCenterByStaffID", offices.GetCenterByStaffID)
	cardIncEnpoint.Post("/updateCenterTagStaff", offices.UpdateCenterTagStaff)

	// MLNI Tracking API's
	mlniTrackingEnpoint := adminEnpoint.Group("/mlni", authentication.ValidateAdminToken)
	// users
	mlniTrackingEnpoint.Get("/getmlniusers/:id", adminmlni.GetMlniUsers)
	mlniTrackingEnpoint.Post("/updatemlniusers/:id", adminmlni.UpdateMlniUser)
}
