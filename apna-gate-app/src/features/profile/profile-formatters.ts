const EMPTY = "—";

export function formatPhoneDisplay(phone?: string | null) {
  if (!phone?.trim()) {
    return EMPTY;
  }

  const digits = phone.replace(/\D/g, "");

  if (digits.length === 10) {
    return `+91 ${digits.slice(0, 5)} ${digits.slice(5)}`;
  }

  if (digits.length === 12 && digits.startsWith("91")) {
    const local = digits.slice(2);
    return `+91 ${local.slice(0, 5)} ${local.slice(5)}`;
  }

  if (digits.length === 11 && digits.startsWith("0")) {
    const local = digits.slice(1);
    return `+91 ${local.slice(0, 5)} ${local.slice(5)}`;
  }

  return phone.trim();
}

export function formatDateOfBirth(dob?: string | null) {
  if (!dob?.trim()) {
    return EMPTY;
  }

  const date = new Date(dob);

  if (Number.isNaN(date.getTime())) {
    return dob;
  }

  return new Intl.DateTimeFormat(undefined, {
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(date);
}

export function formatGender(gender?: string | null) {
  if (!gender?.trim()) {
    return EMPTY;
  }

  return gender
    .replace(/_/g, " ")
    .split(" ")
    .filter(Boolean)
    .map((part) => `${part.charAt(0).toUpperCase()}${part.slice(1).toLowerCase()}`)
    .join(" ");
}

const TIMEZONE_LABELS: Record<string, string> = {
  "Asia/Kolkata": "India Standard Time (IST)",
  "Asia/Calcutta": "India Standard Time (IST)",
  UTC: "Coordinated Universal Time (UTC)",
};

export function formatTimezone(timezone?: string | null) {
  if (!timezone?.trim()) {
    return EMPTY;
  }

  return TIMEZONE_LABELS[timezone] ?? timezone.replace(/_/g, " ");
}

export function formatLanguage(language?: string | null) {
  if (!language?.trim()) {
    return EMPTY;
  }

  try {
    const display = new Intl.DisplayNames(undefined, { type: "language" }).of(language);
    return display ? `${display} (${language})` : language.toUpperCase();
  } catch {
    return language.toUpperCase();
  }
}

export function formatDisplayValue(value?: string | null) {
  if (!value?.trim()) {
    return EMPTY;
  }

  return value.trim();
}

export const PROFILE_GENDER_OPTIONS = [
  { label: "Male", value: "male" },
  { label: "Female", value: "female" },
  { label: "Other", value: "other" },
  { label: "Prefer not to say", value: "prefer_not_to_say" },
] as const;

export type ProfileGenderValue = (typeof PROFILE_GENDER_OPTIONS)[number]["value"];

function padDatePart(value: number) {
  return String(value).padStart(2, "0");
}

function toIsoDate(year: number, month: number, day: number) {
  if (month < 1 || month > 12 || day < 1 || day > 31 || year < 1900 || year > 2100) {
    return null;
  }

  const date = new Date(Date.UTC(year, month - 1, day));
  if (
    date.getUTCFullYear() !== year ||
    date.getUTCMonth() !== month - 1 ||
    date.getUTCDate() !== day
  ) {
    return null;
  }

  return `${year}-${padDatePart(month)}-${padDatePart(day)}`;
}

function extractIsoDatePortion(raw: string): string | null {
  const isoDateMatch = raw.match(/^(\d{4}-\d{2}-\d{2})/);
  if (isoDateMatch) {
    return isoDateMatch[1];
  }

  const parsed = new Date(raw);
  if (Number.isNaN(parsed.getTime())) {
    return null;
  }

  return toIsoDate(parsed.getUTCFullYear(), parsed.getUTCMonth() + 1, parsed.getUTCDate());
}

export function formatDateOfBirthForInput(dob?: string | null): string {
  const trimmed = dob?.trim();
  if (!trimmed) {
    return "";
  }

  const iso = extractIsoDatePortion(trimmed);
  return iso ?? trimmed;
}

export function normalizeProfileDateInput(raw?: string | null):
  | { ok: true; value: string }
  | { ok: false; error: string } {
  const trimmed = raw?.trim();
  if (!trimmed) {
    return { ok: true, value: "" };
  }

  const isoDatePrefix = trimmed.match(/^(\d{4}-\d{2}-\d{2})(?:T|$)/);
  if (isoDatePrefix) {
    const [year, month, day] = isoDatePrefix[1].split("-").map(Number);
    const iso = toIsoDate(year, month, day);
    return iso ? { ok: true, value: iso } : { ok: false, error: "Enter a valid date of birth." };
  }

  if (/^\d{4}-\d{2}-\d{2}$/.test(trimmed)) {
    const [year, month, day] = trimmed.split("-").map(Number);
    const iso = toIsoDate(year, month, day);
    return iso ? { ok: true, value: iso } : { ok: false, error: "Enter a valid date of birth." };
  }

  const slashMatch = trimmed.match(/^(\d{4})[/-](\d{1,2})[/-](\d{1,2})$/);
  if (slashMatch) {
    const iso = toIsoDate(Number(slashMatch[1]), Number(slashMatch[2]), Number(slashMatch[3]));
    return iso ? { ok: true, value: iso } : { ok: false, error: "Enter a valid date of birth." };
  }

  const dmyMatch = trimmed.match(/^(\d{1,2})[/-](\d{1,2})[/-](\d{4})$/);
  if (dmyMatch) {
    const iso = toIsoDate(Number(dmyMatch[3]), Number(dmyMatch[2]), Number(dmyMatch[1]));
    return iso ? { ok: true, value: iso } : { ok: false, error: "Enter a valid date of birth." };
  }

  return { ok: false, error: "Use YYYY-MM-DD format for date of birth." };
}
