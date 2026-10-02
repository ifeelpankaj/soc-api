package imagesvc

// ValidateImage shares the bounded JPEG/PNG decoder with other private media modules.
// Callers must cap the byte size before calling it.
func ValidateImage(data []byte) (string, error) { return validateImage(data) }
