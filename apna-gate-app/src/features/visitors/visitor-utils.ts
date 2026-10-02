import type {
  ModelsVisitorEntry,
  ModelsVisitorPendingEntry,
  ModelsVisitorPurpose,
  ModelsVisitorStatus,
} from "@/lib/api/generated-api";

type VisitorEntryWithFlatFields = (ModelsVisitorEntry | ModelsVisitorPendingEntry) & {
  email?: string;
  full_name?: string;
  phone_number?: string;
};
type FlatVisitorField = "email" | "full_name" | "phone_number";

export type NormalizedCompanion = {
  name: string;
  phoneNumber: string;
};

export const visitorPurposes: ModelsVisitorPurpose[] = [
  "guest",
  "delivery",
  "cab",
  "service",
  "maintenance",
  "staff",
  "other",
];

export const residentInvitePurposes = visitorPurposes.filter(
  (purpose) => purpose !== "service" && purpose !== "staff",
);

export function titleize(value?: string | null) {
  if (!value) {
    return "-";
  }

  return value
    .replace(/_/g, " ")
    .split(" ")
    .filter(Boolean)
    .map((part) => `${part.charAt(0).toUpperCase()}${part.slice(1)}`)
    .join(" ");
}

function presentString(value?: string | null) {
  const trimmed = value?.trim();
  return trimmed ? trimmed : undefined;
}

function getFlatEntryField(
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry | undefined,
  field: FlatVisitorField,
) {
  return presentString((entry as VisitorEntryWithFlatFields | undefined)?.[field]);
}

function getStringFromRecord(
  record: Record<string, unknown>,
  field: string,
) {
  const value = record[field];
  return typeof value === "string" ? presentString(value) : undefined;
}

export function getVisitorName(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
) {
  return (
    presentString(entry?.visitor?.full_name) ??
    getFlatEntryField(entry, "full_name") ??
    (entry?.purpose ? titleize(entry.purpose) : "Visitor")
  );
}

export function getVisitorPhone(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
) {
  return (
    presentString(entry?.visitor?.phone_number) ??
    getFlatEntryField(entry, "phone_number")
  );
}

export function getVisitorEmail(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
) {
  return presentString(entry?.visitor?.email) ?? getFlatEntryField(entry, "email");
}

export function getCompanionDetails(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
): NormalizedCompanion[] {
  return (
    entry?.companion_details
      ?.map((detail) => {
        const record = detail as Record<string, unknown>;
        const name =
          getStringFromRecord(record, "full_name") ??
          getStringFromRecord(record, "name") ??
          "";
        const phoneNumber =
          getStringFromRecord(record, "phone_number") ??
          getStringFromRecord(record, "phoneNumber") ??
          "";

        return { name, phoneNumber };
      })
      .filter((detail) => detail.name || detail.phoneNumber) ?? []
  );
}

export function getFlatLabel(entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry) {
  const flat = entry?.flat;

  if (!flat) {
    if (entry?.purpose === "staff") {
      return "Society staff";
    }
    return entry?.flat_id ? `Flat #${entry.flat_id}` : "Flat pending";
  }

  return `${flat.block ? `${flat.block}-` : ""}${flat.flat_number ?? flat.id ?? "-"}`;
}

export function getFlatLocationParts(entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry) {
  const flat = entry?.flat;

  if (!flat) {
    if (entry?.purpose === "staff") {
      return {
        number: "Society staff",
        wing: null as string | null,
      };
    }
    return {
      number: entry?.flat_id ? `#${entry.flat_id}` : "-",
      wing: null as string | null,
    };
  }

  const number = flat.flat_number ?? (flat.id ? `#${flat.id}` : "-");
  const wing = flat.block
    ? titleize(flat.block.replace(/-/g, " "))
    : flat.floor
      ? `Floor ${flat.floor}`
      : null;

  return { number, wing };
}

export function getFlatLocationLabel(entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry) {
  const { number, wing } = getFlatLocationParts(entry);

  if (wing) {
    return `${wing} · ${number}`;
  }

  return number;
}

export function formatRelativeTime(value?: string | null) {
  if (!value) {
    return "";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "";
  }

  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);

  if (seconds < 60) {
    return "Just now";
  }

  const minutes = Math.floor(seconds / 60);

  if (minutes < 60) {
    return `${minutes}m ago`;
  }

  const hours = Math.floor(minutes / 60);

  if (hours < 24) {
    return `${hours}h ago`;
  }

  const days = Math.floor(hours / 24);

  if (days === 1) {
    return "Yesterday";
  }

  if (days < 7) {
    return `${days}d ago`;
  }

  return formatDateTime(value);
}

export function formatWaitingDuration(value?: string | null) {
  if (!value) {
    return "Waiting";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "Waiting";
  }

  const minutes = Math.max(1, Math.floor((Date.now() - date.getTime()) / 60000));

  if (minutes === 1) {
    return "Waiting 1 minute";
  }

  return `Waiting ${minutes} minutes`;
}

export type ActivityDisplay = {
  flat: string;
  icon: { ios: string; android: string; web: string };
  time: string;
  tint: string;
  title: string;
};

export function getActivityDisplay(entry?: ModelsVisitorEntry): ActivityDisplay {
  const flat = getFlatLabel(entry);
  const time = formatRelativeTime(
    entry?.checked_in_at ?? entry?.checked_out_at ?? entry?.updated_at ?? entry?.created_at,
  );
  const name = getVisitorName(entry);
  const purpose = entry?.purpose;
  const status = entry?.status;

  if (status === "waiting_approval") {
    if (purpose === "delivery") {
      return {
        flat,
        icon: { ios: "shippingbox.fill", android: "local_shipping", web: "local_shipping" },
        time,
        tint: "#d97706",
        title: "Delivery waiting approval",
      };
    }

    if (purpose === "cab") {
      return {
        flat,
        icon: { ios: "car.fill", android: "local_taxi", web: "local_taxi" },
        time,
        tint: "#d97706",
        title: "Cab waiting approval",
      };
    }

    return {
      flat,
      icon: { ios: "clock.fill", android: "schedule", web: "schedule" },
      time,
      tint: "#d97706",
      title: `${name} waiting approval`,
    };
  }

  if (status === "checked_in") {
    if (purpose === "cab") {
      return {
        flat,
        icon: { ios: "car.fill", android: "local_taxi", web: "local_taxi" },
        time,
        tint: "#0f766e",
        title: "Cab arrived",
      };
    }

    if (purpose === "delivery") {
      return {
        flat,
        icon: { ios: "shippingbox.fill", android: "local_shipping", web: "local_shipping" },
        time,
        tint: "#0f766e",
        title: "Delivery checked in",
      };
    }

    return {
      flat,
      icon: { ios: "checkmark.circle.fill", android: "check_circle", web: "check_circle" },
      time,
      tint: "#0f766e",
      title: `${name} checked in`,
    };
  }

  if (status === "checked_out") {
    return {
      flat,
      icon: { ios: "arrow.right.circle.fill", android: "logout", web: "logout" },
      time,
      tint: "#64748b",
      title: `${name} checked out`,
    };
  }

  if (status === "approved") {
    return {
      flat,
      icon: { ios: "qrcode", android: "qr_code", web: "qr_code" },
      time,
      tint: "#059669",
      title: `${name} approved to enter`,
    };
  }

  return {
    flat,
    icon: { ios: "person.fill", android: "person", web: "person" },
    time,
    tint: "#64748b",
    title: name,
  };
}

export function formatTimeOfDay(value?: string | null) {
  if (!value) {
    return "";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "";
  }

  return new Intl.DateTimeFormat(undefined, {
    hour: "numeric",
    minute: "2-digit",
  }).format(date);
}

export function formatDateOnly(value?: string | null, timeZone?: string) {
  if (!value) {
    return "";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "";
  }

  return new Intl.DateTimeFormat(undefined, {
    day: "2-digit",
    month: "short",
    year: "numeric",
    timeZone,
  }).format(date);
}

export function formatDateTime(value?: string | null) {
  if (!value) {
    return "-";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat(undefined, {
    day: "2-digit",
    hour: "numeric",
    minute: "2-digit",
    month: "short",
  }).format(date);
}

export function formatActivityTimestamp(value?: string | null) {
  if (!value) {
    return "";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  const time = formatTimeOfDay(value);

  if (!time) {
    return formatDateTime(value);
  }

  if (isToday(value)) {
    return `Today, ${time}`;
  }

  const now = new Date();
  const yesterday = new Date(now);
  yesterday.setDate(yesterday.getDate() - 1);

  if (
    date.getFullYear() === yesterday.getFullYear() &&
    date.getMonth() === yesterday.getMonth() &&
    date.getDate() === yesterday.getDate()
  ) {
    return `Yesterday, ${time}`;
  }

  return formatDateTime(value);
}

export type VisitorStatusMeta = {
  bg: string;
  border: string;
  color: string;
  label: string;
};

export function getVisitorStatusMeta(status?: ModelsVisitorStatus): VisitorStatusMeta {
  switch (status) {
    case "waiting_approval":
      return {
        bg: "#fffbeb",
        border: "#fde68a",
        color: "#92400e",
        label: "Pending",
      };
    case "approved":
      return {
        bg: "#ecfdf5",
        border: "#bbf7d0",
        color: "#166534",
        label: "Approved",
      };
    case "checked_in":
      return {
        bg: "#ecfdf5",
        border: "#bbf7d0",
        color: "#166534",
        label: "Checked In",
      };
    case "checked_out":
      return {
        bg: "#f8fafc",
        border: "#e2e8f0",
        color: "#64748b",
        label: "Checked Out",
      };
    case "rejected":
      return {
        bg: "#fef2f2",
        border: "#fecaca",
        color: "#b91c1c",
        label: "Rejected",
      };
    case "expired":
      return {
        bg: "#fef2f2",
        border: "#fecaca",
        color: "#b91c1c",
        label: "Expired",
      };
    case "cancelled":
      return {
        bg: "#f8fafc",
        border: "#e2e8f0",
        color: "#64748b",
        label: "Cancelled",
      };
    case "auto_closed":
      return {
        bg: "#f8fafc",
        border: "#e2e8f0",
        color: "#64748b",
        label: "Auto Closed",
      };
    default:
      return {
        bg: "#f8fafc",
        border: "#e2e8f0",
        color: "#64748b",
        label: titleize(status),
      };
  }
}

export function statusTone(status?: ModelsVisitorStatus) {
  switch (status) {
    case "approved":
      return "border-emerald-200 bg-emerald-50 text-emerald-800";
    case "checked_in":
      return "border-teal-200 bg-teal-50 text-teal-800";
    case "checked_out":
      return "border-slate-200 bg-slate-100 text-slate-700";
    case "waiting_approval":
      return "border-amber-200 bg-amber-50 text-amber-800";
    case "rejected":
    case "cancelled":
    case "expired":
      return "border-rose-200 bg-rose-50 text-rose-700";
    default:
      return "border-stone-200 bg-stone-50 text-stone-700";
  }
}

export function isToday(value?: string | null) {
  if (!value) {
    return false;
  }

  const date = new Date(value);
  const now = new Date();

  return (
    date.getFullYear() === now.getFullYear() &&
    date.getMonth() === now.getMonth() &&
    date.getDate() === now.getDate()
  );
}

export function isExpectedTodayEntry(entry?: ModelsVisitorEntry) {
  if (!entry || entry.status !== "approved") {
    return false;
  }

  return isToday(entry.expected_at) || isToday(entry.created_at);
}

export function isTodayActivityEntry(entry?: ModelsVisitorEntry) {
  if (!entry) {
    return false;
  }

  return (
    isToday(entry.created_at) ||
    isToday(entry.checked_in_at) ||
    isToday(entry.checked_out_at) ||
    isToday(entry.expected_at)
  );
}

export type WaitingDurationTone = "success" | "warning" | "error";

export function getWaitingDuration(approvedAt?: string | null): {
  label: string;
  minutes: number;
  tone: WaitingDurationTone;
} {
  if (!approvedAt) {
    return { label: "Waiting", minutes: 0, tone: "success" };
  }

  const approved = new Date(approvedAt);
  const minutes = Math.max(0, Math.floor((Date.now() - approved.getTime()) / 60000));
  const label = minutes <= 1 ? "Waiting 1 min" : `Waiting ${minutes} min`;

  if (minutes < 5) {
    return { label, minutes, tone: "success" };
  }
  if (minutes < 15) {
    return { label, minutes, tone: "warning" };
  }
  return { label, minutes, tone: "error" };
}

export type VisitorTimelineBlock = {
  date?: string;
  label: string;
  time?: string;
};

export type VisitorDetailRow = {
  label: string;
  value: string;
};

export function getVisitorTimelineBlocks(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
): VisitorTimelineBlock[] {
  const blocks: VisitorTimelineBlock[] = [];

  if (entry?.checked_in_at) {
    blocks.push({
      label: "Check-in",
      date: formatDateOnly(entry.checked_in_at),
      time: formatTimeOfDay(entry.checked_in_at),
    });
  } else if (entry?.expected_at) {
    blocks.push({
      label: "Expected check-in",
      date: formatDateOnly(entry.expected_at),
      time: formatTimeOfDay(entry.expected_at),
    });
  }

  if (entry?.checked_out_at) {
    blocks.push({
      label: "Checkout",
      date: formatDateOnly(entry.checked_out_at),
      time: formatTimeOfDay(entry.checked_out_at),
    });
  } else if (entry?.expected_checkout_at) {
    blocks.push({
      label: "Expected checkout",
      date: formatDateOnly(entry.expected_checkout_at),
      time: formatTimeOfDay(entry.expected_checkout_at),
    });
  }

  return blocks;
}

export function getVisitorDetailRows(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
): VisitorDetailRow[] {
  if (!entry) {
    return [];
  }

  const rows: VisitorDetailRow[] = [];

  const phone = getVisitorPhone(entry);
  const email = getVisitorEmail(entry);

  if (phone) {
    rows.push({ label: "Mobile", value: phone });
  }
  if (email) {
    rows.push({ label: "Email", value: email });
  }
  if (entry.vehicle_number) {
    rows.push({ label: "Vehicle", value: entry.vehicle_number });
  }
  if (entry.companions_count && entry.companions_count > 0) {
    rows.push({ label: "Companions", value: String(entry.companions_count) });
  }
  if (entry.delivery_partner) {
    rows.push({ label: "Partner", value: entry.delivery_partner });
  }
  if (entry.service_provider) {
    rows.push({ label: "Provider", value: entry.service_provider });
  }

  return rows;
}

export function isVisitorCheckoutOverdue(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
): boolean {
  if (entry?.status !== "checked_in" || !entry.expected_checkout_at) {
    return false;
  }

  return new Date(entry.expected_checkout_at).getTime() < Date.now();
}

export function getVisitorStatusContextMessage(
  entry?: ModelsVisitorEntry | ModelsVisitorPendingEntry,
): string | null {
  if (!entry?.status) {
    return null;
  }

  switch (entry.status) {
    case "expired":
      if (entry.expected_checkout_at) {
        return `This visit expired after the expected checkout window (${formatDateTime(entry.expected_checkout_at)}).`;
      }
      return "This visit is no longer active at the gate.";
    case "checked_out":
      return entry.checked_out_at
        ? `Visitor checked out on ${formatDateTime(entry.checked_out_at)}.`
        : "This visit has been completed.";
    case "rejected":
      return "This visit was declined and cannot be used for entry.";
    case "cancelled":
      return "This visit was cancelled before check-in.";
    case "waiting_approval":
      return "Waiting for resident approval before the visitor can enter.";
    case "approved":
      return entry.expected_at
        ? `Approved for entry. Expected ${formatDateTime(entry.expected_at)}.`
        : "Approved and ready for gate check-in.";
    case "checked_in":
      if (isVisitorCheckoutOverdue(entry)) {
        return `Overdue checkout. Expected departure was ${formatDateTime(entry.expected_checkout_at!)}. Visitor remains checked in until guard checkout.`;
      }
      return entry.checked_in_at
        ? `Currently inside since ${formatDateTime(entry.checked_in_at)}.`
        : "Visitor is currently inside the society.";
    default:
      return null;
  }
}

export function waitingToneStyle(tone?: WaitingDurationTone) {
  switch (tone) {
    case "warning":
      return {
        backgroundColor: "#fffbeb",
        borderColor: "#fde68a",
        color: "#d97706",
      };
    case "error":
      return {
        backgroundColor: "#fef2f2",
        borderColor: "#fecaca",
        color: "#b91c1c",
      };
    default:
      return {
        backgroundColor: "#ecfdf5",
        borderColor: "#a7f3d0",
        color: "#059669",
      };
  }
}
