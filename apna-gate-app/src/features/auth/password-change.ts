import { getApiErrorCode, getApiMessage } from "./api-error";

export type PasswordField = "current" | "next" | "confirm";
export type PasswordValues = Record<PasswordField, string>;
export type PasswordErrors = Partial<Record<PasswordField, string>>;

export function utf8Length(value: string) {
  return Array.from(value).reduce((bytes, char) => {
    const code = char.codePointAt(0)!;
    return (
      bytes + (code <= 0x7f ? 1 : code <= 0x7ff ? 2 : code <= 0xffff ? 3 : 4)
    );
  }, 0);
}

export function validatePasswordChange(values: PasswordValues): PasswordErrors {
  const errors: PasswordErrors = {};
  if (!values.current) errors.current = "Enter your current password.";
  if (!values.next) errors.next = "Enter a new password.";
  else if (Array.from(values.next).length < 8)
    errors.next = "Use at least 8 characters.";
  else if (utf8Length(values.next) > 72)
    errors.next =
      "Use at most 72 bytes. Some characters use more than one byte.";
  else if (values.next === values.current)
    errors.next = "Choose a password different from your current password.";
  if (!values.confirm) errors.confirm = "Confirm your new password.";
  else if (values.next !== values.confirm)
    errors.confirm = "Passwords do not match.";
  return errors;
}

export function passwordChangeError(error: unknown): {
  message: string;
  field?: PasswordField;
  uncertain?: boolean;
  reset?: boolean;
} {
  const code = getApiErrorCode(error);
  if (code === "CURRENT_PASSWORD_INCORRECT" || code === "INVALID_CREDENTIALS") {
    return { field: "current", message: "Current password is incorrect." };
  }
  if (code === "PASSWORD_NOT_SET")
    return {
      field: "current",
      reset: true,
      message: "No password is set. Use password reset to create one.",
    };
  if (code === "PASSWORD_TOO_LONG")
    return {
      field: "next",
      message: "Use at most 72 bytes. Some characters use more than one byte.",
    };
  if (code === "PASSWORD_REUSE")
    return {
      field: "next",
      message: "Choose a password different from your current password.",
    };
  if (code === "PASSWORD_MISMATCH")
    return { field: "confirm", message: "Passwords do not match." };
  if (
    error &&
    typeof error === "object" &&
    "status" in error &&
    ["TIMEOUT_ERROR", "FETCH_ERROR"].includes(String(error.status))
  ) {
    return {
      uncertain: true,
      message:
        "We couldn't confirm whether your password changed. Try signing in with your new password, or reset your password.",
    };
  }
  return {
    message: getApiMessage(
      error,
      "Could not change password. Please try again.",
    ),
  };
}
