package constants

func actionFailed(action string) string {
	return "Something went wrong while " + action + ". Please try again."
}

func temporarilyFailed(action string) string {
	return "We couldn’t " + action + " right now. Please try again in a bit."
}

const (
	// Auth / access
	ErrUnknownID        = "Access denied. Please try again."
	ErrUserNotFound     = "We couldn’t find this user. Please check and try again."
	ErrUserUnverified   = "Please verify your account before continuing."
	ErrInvalidLogin     = "Email or password is incorrect."
	ErrPasswordMismatch = "Password and confirm password don’t match."
	ErrEmailExists      = "Email already exists. Please log in instead."
	ErrPhoneExists      = "Phone number already exists. Please log in instead."

	// OTP / verification
	ErrAlreadyVerified = "Your account is already verified. Please log in."
	ErrNoActiveOTP     = "No active OTP found. Please request a new one."
	ErrInvalidOTP      = "Invalid or expired OTP. Please request a new one."

	// Soft success
	MsgLastLoginSkipped = "Request successful. We’ll update your login time shortly."

	// Onboarding / registration
	ErrStaffLimitReached = "You’ve reached the staff limit for your current plan. Please upgrade to add more staff."

	// Onboarding / society
	ErrSocietyNotFound = "We couldn’t find this society. Please check and try again."
	ErrSocietyInactive = "This society is not active yet. Please contact support."

	// Onboarding / flat
	ErrFlatNotFound        = "We couldn’t find this flat. Please check and try again."
	ErrFlatNotInSociety    = "This flat doesn’t belong to the selected society."
	ErrFlatNotActive       = "This flat is not available right now. Please choose another."
	ErrFlatAlreadyOccupied = "This flat already has a primary resident."
	ErrFlatSocietyMismatch = "This flat does not belong to the claim society."

	// Onboarding / claim
	ErrDuplicateClaim  = "You already have a pending claim for this flat."
	ErrClaimNotFound   = "We couldn’t find this claim request. Please check and try again."
	ErrClaimNotPending = "This claim is no longer pending."

	// Onboarding / invite
	ErrInviteInvalid      = "This invite token is invalid."
	ErrInviteRevoked      = "This invite has been revoked."
	ErrInviteExpired      = "This invite has expired."
	ErrInviteUsed         = "This invite has already been used."
	ErrInviteInvalidState = "This invite is no longer valid."
	ErrInvalidInviteType  = "Invalid invite type."

	// Onboarding / user assignment
	ErrUserAlreadyAssignedFlat = "You are already assigned to a flat."

	// Shared service validation
	ErrInvalidSocietyID = "Invalid society ID. Please check and try again."
	ErrInvalidPlanID    = "Invalid plan ID. Please check and try again."
	ErrInvalidFlatID    = "Invalid flat ID. Please check and try again."
	ErrInvalidVisitorID = "Invalid visitor ID. Please check and try again."

	// Society
	ErrGenerateSocietyCode   = "We couldn't generate a unique society code right now. Please try again."
	ErrSocietyAlreadyExists  = "A society with this name already exists in this city."
	ErrSocietyAlreadyDeleted = "This society has already been deleted."

	// Flat
	ErrFlatNumberAlreadyExists = "This flat number already exists in this society."
	ErrFlatNotBelongToSociety  = "This flat does not belong to this society."

	// Approval settings
	ErrApprovalSettingsRequestRequired = "Approval settings request is required."
	ErrApprovalSettingsRequired        = "At least one approval setting is required."
	ErrInvalidVisitorType              = "Invalid visitor type. Please check and try again."
	ErrInvalidVisitorTypes             = "All visitor type values must be valid."
	ErrApprovalSettingsTimeout         = "Approval settings operation timed out. Please try again."

	// Visitor
	ErrVisitorRequestRequired   = "Visitor request is required."
	ErrInvalidResidentID        = "Invalid resident ID. Please check and try again."
	ErrInvalidGuardID           = "Invalid guard ID. Please check and try again."
	ErrResidentDoesNotOwnFlat   = "You can only create visitor entries for your own flat."
	ErrVisitorEntryNotFound     = "We couldn't find this visitor entry. Please check and try again."
	ErrVisitorConcurrentUpdate  = "This visitor entry was updated by another action. Please refresh and try again."
	ErrQRCodeUnavailable        = "QR code is not available for this visitor entry."
	ErrInvalidQRCodeToken       = "Invalid QR code token. Please scan a valid visitor QR code." //nolint:gosec // User-facing token error message, not a credential.
	ErrVisitorAlreadyApproved   = "This visitor entry is already approved."
	ErrVisitorAlreadyRejected   = "This visitor entry is already rejected."
	ErrVisitorAlreadyCheckedIn  = "This visitor entry is already checked in."
	ErrVisitorAlreadyCheckedOut = "This visitor entry is already checked out."
	ErrVisitorNotApproved       = "This visitor entry is not approved yet."
	ErrVisitorNotCheckedIn      = "This visitor entry has not checked in yet."
	ErrVisitorWasRejected       = "This visitor entry was rejected."
	ErrVisitorExpired           = "This visitor entry has expired."
	ErrInvalidVisitorStatus     = "Invalid visitor status. Please check and try again."

	// Plan
	ErrPlanNotFound        = "We couldn't find this plan. Please check and try again."
	ErrPlanAlreadyExists   = "A plan with this name already exists."
	ErrInvalidBillingCycle = "Billing cycle must be 'monthly' or 'yearly'."

	// Subscription
	ErrInactivePlan                     = "This plan is not active. Please choose another plan."
	ErrActiveSubscriptionExists         = "This society already has an active subscription. Please change the plan instead."
	ErrNoActiveSubscription             = "No active subscription found for this society."
	ErrSubscriptionAlreadyCancelled     = "This subscription is already cancelled."
	ErrSameSubscriptionPlan             = "This society is already subscribed to this plan."
	ErrCancelledSubscriptionRenewal     = "Cancelled subscriptions cannot be renewed. Please create a new subscription instead."
	ErrSubscriptionInactive             = "This subscription is not active."
	ErrUnknownSubscriptionFeatureFormat = "Unknown subscription feature: %q."
	ErrSubscriptionLimitExceededFormat  = "%s limit of %d reached. Current usage is %d."
)

var (
	// User
	ErrFetchUser  = temporarilyFailed("fetch your profile")
	ErrCreateUser = actionFailed("creating your account")
	ErrUpdateUser = actionFailed("updating your profile")

	// Registration
	ErrCheckEmail   = actionFailed("checking your email")
	ErrCheckPhone   = actionFailed("checking your phone number")
	ErrHashPassword = actionFailed("processing your password")
	ErrGenerateOTP  = actionFailed("generating OTP")
	ErrSendEmail    = temporarilyFailed("send the email")

	// Verification
	ErrFindVerificationUser = actionFailed("verifying your email")
	ErrUpdateVerification   = actionFailed("updating verification status")
	ErrCreateVerification   = actionFailed("creating verification request")
	ErrMarkOTPUsed          = actionFailed("marking OTP as used")

	// Transactions
	ErrRegistrationTx      = actionFailed("completing registration")
	ErrVerificationTx      = actionFailed("completing verification")
	ErrStaffRegistrationTx = actionFailed("completing staff registration")

	// Onboarding / generic fetch
	ErrFetchSociety       = temporarilyFailed("fetch society details")
	ErrFetchFlats         = temporarilyFailed("fetch flats")
	ErrFetchFlat          = temporarilyFailed("fetch flat details")
	ErrFetchClaim         = temporarilyFailed("fetch claim details")
	ErrFetchClaimant      = temporarilyFailed("fetch claimant details")
	ErrFetchInvite        = temporarilyFailed("fetch invite details")
	ErrFetchGuards        = temporarilyFailed("fetch guards")
	ErrFetchPendingClaims = temporarilyFailed("fetch pending claims")
	ErrFetchClaims        = temporarilyFailed("fetch claims")
	ErrFetchClaimHistory  = temporarilyFailed("fetch claim history")
	ErrFetchUpdatedUser   = temporarilyFailed("fetch updated user")
	ErrFetchUpdatedFlat   = temporarilyFailed("fetch updated flat")

	// Onboarding / operations
	ErrCheckFlatOccupancy  = temporarilyFailed("check flat occupancy")
	ErrCheckExistingClaims = temporarilyFailed("check existing claims")
	ErrSubmitClaim         = actionFailed("submitting your claim")
	ErrApproveClaim        = actionFailed("approving the claim")
	ErrRejectClaim         = actionFailed("rejecting the claim")
	ErrRedeemInvite        = actionFailed("redeeming the invite")
	ErrLookupInvite        = actionFailed("looking up the invite")
	ErrAssignSociety       = actionFailed("assigning society")
	ErrAssignFlat          = actionFailed("assigning your flat")
	ErrActivateUser        = actionFailed("activating your account")
	ErrUpdateFlatStatus    = actionFailed("updating flat status")
	ErrRevokeInvites       = actionFailed("revoking old invites")
	ErrIncrementInvite     = actionFailed("redeeming the invite")
	ErrFetchChecklistUser  = temporarilyFailed("fetch checklist user")

	// Society operations
	ErrCreateSociety        = actionFailed("creating society")
	ErrAssignSocietyCreator = actionFailed("assigning society creator")
	ErrListSocieties        = temporarilyFailed("fetch societies")
	ErrUpdateSociety        = actionFailed("updating society")
	ErrActivateSociety      = actionFailed("activating society")
	ErrDeactivateSociety    = actionFailed("deactivating society")
	ErrDeleteSociety        = actionFailed("deleting society")

	// Subscription operations
	ErrFetchPlan                    = temporarilyFailed("fetch plan details")
	ErrFetchNewPlan                 = temporarilyFailed("fetch new plan details")
	ErrCheckExistingSubscription    = temporarilyFailed("check existing subscription")
	ErrCreateSubscription           = actionFailed("creating subscription")
	ErrFetchActiveSubscription      = temporarilyFailed("fetch active subscription")
	ErrListSubscriptions            = temporarilyFailed("fetch subscription history")
	ErrFetchCurrentSubscription     = temporarilyFailed("fetch current subscription")
	ErrChangeSubscriptionPlan       = actionFailed("changing subscription plan")
	ErrCancelSubscription           = actionFailed("cancelling subscription")
	ErrFetchSubscriptionsForRenewal = temporarilyFailed("fetch subscriptions for renewal")
	ErrRenewSubscription            = actionFailed("renewing subscription")
	ErrCheckSubscriptionStatus      = temporarilyFailed("check subscription status")
	ErrFetchUsageSubscription       = temporarilyFailed("fetch subscription usage")

	// Flat operations
	ErrCreateFlat      = actionFailed("creating flat")
	ErrUpdateFlat      = actionFailed("updating flat")
	ErrActivateFlat    = actionFailed("activating flat")
	ErrDeactivateFlat  = actionFailed("deactivating flat")
	ErrDeleteFlat      = actionFailed("deleting flat")
	ErrCheckFlatNumber = temporarilyFailed("check flat number")

	// Approval setting operations
	ErrFetchApprovalSettings = temporarilyFailed("fetch approval settings")
	ErrSaveApprovalSetting   = actionFailed("saving approval setting")
	ErrSaveApprovalSettings  = actionFailed("saving approval settings")
	ErrDeleteApprovalSetting = actionFailed("deleting approval setting")

	// Plan operations
	ErrCreatePlan                     = actionFailed("creating plan")
	ErrFetchPlanDetails               = temporarilyFailed("fetch plan details")
	ErrListActivePlans                = temporarilyFailed("fetch active plans")
	ErrListPlans                      = temporarilyFailed("fetch plans")
	ErrUpdatePlan                     = actionFailed("updating plan")
	ErrActivatePlan                   = actionFailed("activating plan")
	ErrDeactivatePlan                 = actionFailed("deactivating plan")
	ErrFetchActiveSubscriptionForPlan = temporarilyFailed("fetch active subscription")
	ErrFetchStaffCount                = temporarilyFailed("count staff")

	// Visitor operations
	ErrGenerateQRToken      = actionFailed("generating QR token")
	ErrCreateVisitorEntry   = actionFailed("creating visitor entry")
	ErrFetchVisitorEntry    = temporarilyFailed("fetch visitor entry")
	ErrApproveVisitorEntry  = actionFailed("approving visitor entry")
	ErrRejectVisitorEntry   = actionFailed("rejecting visitor entry")
	ErrCheckInVisitorEntry  = actionFailed("checking in visitor")
	ErrCheckOutVisitorEntry = actionFailed("checking out visitor")
	ErrFetchVisitorEntries  = temporarilyFailed("fetch visitor entries")
	ErrReadApprovalSettings = temporarilyFailed("read approval settings")
)
