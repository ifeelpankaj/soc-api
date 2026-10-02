import type {
  ModelsShortLinkResponse,
  ModelsVisitorPurpose,
} from "@/lib/api/generated-api";
import type { ModelsFlatMemberInviteRole } from "@/lib/api/resident-api-extensions";

export type InviteShareContent = {
  title: string;
  message: string;
  url: string;
};

export type InviteShareResult =
  | { available: true; content: InviteShareContent }
  | { available: false; reason: "short_link_unavailable" };

export type VisitorInviteShareInput = {
  expectedAt?: string;
  expiresAt?: string;
  flatLabel?: string;
  purpose?: ModelsVisitorPurpose | string;
  shortLink?: ModelsShortLinkResponse | null;
  societyName?: string | null;
};

export type MemberInviteShareInput = {
  expiresAt?: string;
  flatLabel?: string;
  fullName?: string | null;
  role?: ModelsFlatMemberInviteRole | string;
  shortLink?: ModelsShortLinkResponse | null;
  societyName?: string | null;
};

export function buildVisitorInviteShareContent(
  invite: VisitorInviteShareInput,
): InviteShareResult {
  const url = invite.shortLink?.url?.trim();
  if (!url) {
    return { available: false, reason: "short_link_unavailable" };
  }

  const societyName = cleanText(invite.societyName) || "your society";
  const purposeLabel = titleize(invite.purpose ?? "guest");
  const lines = compactLines([
    `You're invited to visit ${societyName} 🏠`,
    "",
    `Purpose: ${purposeLabel}`,
    invite.flatLabel ? `Flat: ${invite.flatLabel}` : null,
    invite.expectedAt ? `Expected: ${formatShareDateTime(invite.expectedAt)}` : null,
    "",
    "Complete your visitor details here:",
    url,
    "",
    invite.expiresAt
      ? `This invite expires on ${formatShareDateTime(invite.expiresAt)}.`
      : null,
    "",
    "— Apna Gate",
  ]);

  return {
    available: true,
    content: {
      message: lines.join("\n"),
      title: "Visitor invite",
      url,
    },
  };
}

export function buildMemberInviteShareContent(
  invite: MemberInviteShareInput,
): InviteShareResult {
  const url = invite.shortLink?.url?.trim();
  if (!url) {
    return { available: false, reason: "short_link_unavailable" };
  }

  const name = cleanText(invite.fullName) || "Hi";
  const societyName = cleanText(invite.societyName) || "your society";
  const roleLabel = titleize(invite.role ?? "member");
  const lines = compactLines([
    `${name}, you've been invited to join ${societyName} 🏠`,
    "",
    invite.flatLabel ? `Flat: ${invite.flatLabel}` : null,
    `Role: ${roleLabel}`,
    "",
    "Accept your invitation:",
    url,
    "",
    invite.expiresAt
      ? `This invite expires on ${formatShareDateTime(invite.expiresAt)}.`
      : null,
    "",
    "— Apna Gate",
  ]);

  return {
    available: true,
    content: {
      message: lines.join("\n"),
      title: "Member invite",
      url,
    },
  };
}

function titleize(value: string) {
  return value
    .replace(/[_-]+/g, " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function cleanText(value?: string | null) {
  return value?.trim() ?? "";
}

function compactLines(lines: (string | null | undefined)[]) {
  const compacted: string[] = [];

  for (const line of lines) {
    if (line === null || line === undefined) {
      continue;
    }
    if (line === "" && compacted[compacted.length - 1] === "") {
      continue;
    }
    compacted.push(line);
  }

  while (compacted[0] === "") {
    compacted.shift();
  }
  while (compacted[compacted.length - 1] === "") {
    compacted.pop();
  }

  return compacted;
}

function formatShareDateTime(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString(undefined, {
    day: "2-digit",
    hour: "numeric",
    minute: "2-digit",
    month: "short",
  });
}
