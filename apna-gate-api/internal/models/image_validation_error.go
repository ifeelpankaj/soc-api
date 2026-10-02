package models

var (
	ErrAvatarUploadRequired       = NewAppError("AVATAR_UPLOAD_REQUIRED", "Use PUT /v1/auth/profile/avatar to upload a profile photo", 400, nil)
	ErrVisitorPhotoUploadRequired = NewAppError("VISITOR_PHOTO_UPLOAD_REQUIRED", "Upload visitor photos through the visitor-entry photo endpoint", 400, nil)
)
