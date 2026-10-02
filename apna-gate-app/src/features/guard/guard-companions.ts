export type CompanionDetail = {
  name: string;
  phoneNumber: string;
};

export type CompanionFieldErrors = {
  name?: string;
  phoneNumber?: string;
};

export function normalizePhone(value: string) {
  return value.replace(/\D/g, "");
}

export function isValidPhone(value: string) {
  const digits = normalizePhone(value);
  return digits.length >= 10 && digits.length <= 15;
}

export function emptyCompanion(): CompanionDetail {
  return { name: "", phoneNumber: "" };
}

export function resizeCompanions(
  current: CompanionDetail[],
  count: number,
) {
  if (count <= 0) {
    return [];
  }
  if (current.length === count) {
    return current;
  }
  if (current.length < count) {
    return [
      ...current,
      ...Array.from({ length: count - current.length }, emptyCompanion),
    ];
  }
  return current.slice(0, count);
}

export function validateCompanionDetails(
  companions: CompanionDetail[],
  count: number,
): CompanionFieldErrors[] {
  const errors: CompanionFieldErrors[] = [];

  for (let index = 0; index < count; index += 1) {
    const companion = companions[index] ?? emptyCompanion();
    const name = companion.name.trim();
    const phone = companion.phoneNumber.trim();
    const fieldErrors: CompanionFieldErrors = {};

    if (!name && !phone) {
      fieldErrors.name = "Enter name or phone";
      fieldErrors.phoneNumber = "Enter name or phone";
    } else if (phone && !isValidPhone(phone)) {
      fieldErrors.phoneNumber = "Enter a valid 10-digit phone number";
    }

    if (fieldErrors.name || fieldErrors.phoneNumber) {
      errors[index] = fieldErrors;
    }
  }

  return errors;
}

export function hasCompanionDetailErrors(errors: CompanionFieldErrors[]) {
  return errors.some((entry) => Boolean(entry?.name || entry?.phoneNumber));
}

export function serializeCompanionDetails(
  companions: CompanionDetail[],
  count: number,
) {
  if (count <= 0) {
    return undefined;
  }

  return companions.slice(0, count).map((companion) => {
    const detail: Record<string, string> = {};
    const name = companion.name.trim();
    const phone = normalizePhone(companion.phoneNumber);

    if (name) {
      detail.full_name = name;
    }
    if (phone) {
      detail.phone_number = phone;
    }

    return detail;
  });
}
