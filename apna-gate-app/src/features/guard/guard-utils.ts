export * from "@/features/visitors/visitor-utils";

export const GUARD_GATE_NAME = "Main Gate";

export function extractQrToken(data: string) {
  const trimmed = data.trim();

  if (!trimmed) {
    return "";
  }

  try {
    const url = new URL(trimmed);
    const token =
      url.searchParams.get("token") ||
      url.searchParams.get("qr_token") ||
      url.searchParams.get("visitor_token");

    if (token) {
      return token;
    }

    const lastSegment = url.pathname.split("/").filter(Boolean).at(-1);
    return lastSegment || trimmed;
  } catch {
    return trimmed;
  }
}

export const DELIVERY_PARTNERS = [
  "Zomato",
  "Swiggy",
  "Blinkit",
  "Amazon",
] as const;

export const DELIVERY_PARTNER_OTHER_LABEL = "Not listed";
