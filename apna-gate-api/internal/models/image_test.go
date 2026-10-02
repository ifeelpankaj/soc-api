package models

import (
	"errors"
	"testing"
)

func TestRawImageURLWritesRejected(t *testing.T) {
	url := "https://external/image.png"
	user := UpdateUserRequest{AvatarURL: &url}
	if !errors.Is(user.Validate(), ErrAvatarUploadRequired) {
		t.Fatal("avatar URL accepted")
	}
	form := VisitorFormRequest{PhotoURL: &url}
	if !errors.Is(form.Validate(true), ErrVisitorPhotoUploadRequired) {
		t.Fatal("visitor URL accepted")
	}
	if !errors.Is(form.ValidateInviteSubmit(), ErrVisitorPhotoUploadRequired) {
		t.Fatal("invite URL accepted")
	}
	update := UpdateGuardVisitorEntryRequest{PhotoURL: &url}
	if !errors.Is(update.Validate(), ErrVisitorPhotoUploadRequired) {
		t.Fatal("update URL accepted")
	}
	blank := " "
	name := "Valid Name"
	update = UpdateGuardVisitorEntryRequest{PhotoURL: &blank, FullName: &name}
	if err := update.Validate(); err != nil || update.PhotoURL != nil {
		t.Fatal("blank URL must not reach database update")
	}
	wrapped := NewAppError("outer", "outer", 400, ErrVisitorPhotoUploadRequired)
	if !errors.Is(wrapped, ErrVisitorPhotoUploadRequired) {
		t.Fatal("lost domain error identity")
	}
}
