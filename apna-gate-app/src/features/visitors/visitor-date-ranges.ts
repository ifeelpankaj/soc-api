import type { ModelsVisitorPurpose, ModelsVisitorStatus } from "@/lib/api/generated-api";

export type LogsSegment = "today" | "expected" | "inside" | "all";
export type DateRangePreset =
  | "today"
  | "yesterday"
  | "last_7_days"
  | "this_month"
  | "all";

const IST_OFFSET_MS = 5.5 * 60 * 60 * 1000;
const IST_DAY_MS = 24 * 60 * 60 * 1000;

function getTodayPartsIST(now = new Date()) {
  const formatter = new Intl.DateTimeFormat("en-CA", {
    timeZone: "Asia/Kolkata",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
  const [year, month, day] = formatter.format(now).split("-").map(Number);
  return { year, month, day };
}

function istDayStartUtc(year: number, month: number, day: number) {
  return new Date(Date.UTC(year, month - 1, day, 0, 0, 0, 0) - IST_OFFSET_MS);
}

function addIstDays(year: number, month: number, day: number, delta: number) {
  const start = istDayStartUtc(year, month, day);
  return new Date(start.getTime() + delta * IST_DAY_MS);
}

function startOfDay(date: Date) {
  const next = new Date(date);
  next.setHours(0, 0, 0, 0);
  return next;
}

function endOfDay(date: Date) {
  const next = new Date(date);
  next.setHours(23, 59, 59, 999);
  return next;
}

export type HalfOpenRange = {
  eventFrom: string;
  eventTo: string;
};

export type VisitorEntryEvent = "activity" | "checked_in" | "checked_out" | "created" | "expected";

/** Half-open [from, to) range aligned to Asia/Kolkata day boundaries. */
export function getHalfOpenRangeIST(preset: DateRangePreset): HalfOpenRange | undefined {
  if (preset === "all") {
    return undefined;
  }

  const { year, month, day } = getTodayPartsIST();
  const todayStart = istDayStartUtc(year, month, day);
  const tomorrowStart = addIstDays(year, month, day, 1);

  switch (preset) {
    case "today":
      return {
        eventFrom: todayStart.toISOString(),
        eventTo: tomorrowStart.toISOString(),
      };
    case "yesterday": {
      const yesterdayStart = addIstDays(year, month, day, -1);
      return {
        eventFrom: yesterdayStart.toISOString(),
        eventTo: todayStart.toISOString(),
      };
    }
    case "last_7_days": {
      const weekStart = addIstDays(year, month, day, -6);
      return {
        eventFrom: weekStart.toISOString(),
        eventTo: tomorrowStart.toISOString(),
      };
    }
    case "this_month": {
      const monthStart = istDayStartUtc(year, month, 1);
      return {
        eventFrom: monthStart.toISOString(),
        eventTo: tomorrowStart.toISOString(),
      };
    }
    default:
      return undefined;
  }
}

export type VisitorEntryQueryFilters = {
  event?: VisitorEntryEvent;
  eventFrom?: string;
  eventTo?: string;
  purpose?: ModelsVisitorPurpose;
  status?: ModelsVisitorStatus;
};

export function buildVisitorEntryQueryFilters(input: {
  datePreset: DateRangePreset;
  isSearchActive: boolean;
  purpose?: ModelsVisitorPurpose;
  segment: LogsSegment;
  sheetStatus?: ModelsVisitorStatus;
}): VisitorEntryQueryFilters {
  const { datePreset, isSearchActive, purpose, segment, sheetStatus } = input;

  if (isSearchActive) {
    return {};
  }

  const range = getHalfOpenRangeIST(datePreset);
  const todayRange = getHalfOpenRangeIST("today");

  if (sheetStatus === "checked_out") {
    const effectiveRange = range ?? todayRange;
    return {
      purpose,
      status: "checked_out",
      event: effectiveRange ? "checked_out" : undefined,
      eventFrom: effectiveRange?.eventFrom,
      eventTo: effectiveRange?.eventTo,
    };
  }

  if (sheetStatus === "checked_in") {
    if (!range) {
      return { purpose };
    }

    return {
      purpose,
      event: "checked_in",
      eventFrom: range.eventFrom,
      eventTo: range.eventTo,
    };
  }

  if (sheetStatus === "approved") {
    const effectiveRange = range ?? todayRange;
    return {
      purpose,
      status: "approved",
      event: effectiveRange ? "expected" : undefined,
      eventFrom: effectiveRange?.eventFrom,
      eventTo: effectiveRange?.eventTo,
    };
  }

  if (sheetStatus === "waiting_approval") {
    const effectiveRange = range ?? todayRange;
    return {
      purpose,
      status: "waiting_approval",
      event: effectiveRange ? "created" : undefined,
      eventFrom: effectiveRange?.eventFrom,
      eventTo: effectiveRange?.eventTo,
    };
  }

  if (sheetStatus === "rejected") {
    if (!range) {
      return { purpose, status: "rejected" };
    }

    return {
      purpose,
      status: "rejected",
      event: "created",
      eventFrom: range.eventFrom,
      eventTo: range.eventTo,
    };
  }

  switch (segment) {
    case "today":
      return {
        purpose,
        event: "activity",
        eventFrom: todayRange?.eventFrom,
        eventTo: todayRange?.eventTo,
      };
    case "expected":
      return {
        purpose,
        status: "approved",
        event: "expected",
        eventFrom: todayRange?.eventFrom,
        eventTo: todayRange?.eventTo,
      };
    case "inside":
      return { purpose, status: "checked_in" };
    case "all":
    default:
      return { purpose };
  }
}

export function buildExplicitVisitorEntryQueryFilters(input: {
  applyDateRange?: boolean;
  datePreset: DateRangePreset;
  purpose?: ModelsVisitorPurpose;
  status?: ModelsVisitorStatus;
}): VisitorEntryQueryFilters {
  const { applyDateRange = true, datePreset, purpose, status } = input;
  const range = applyDateRange ? getHalfOpenRangeIST(datePreset) : undefined;

  if (status === "checked_out") {
    return {
      purpose,
      status,
      event: range ? "checked_out" : undefined,
      eventFrom: range?.eventFrom,
      eventTo: range?.eventTo,
    };
  }

  if (status === "checked_in") {
    return {
      purpose,
      event: range ? "checked_in" : undefined,
      eventFrom: range?.eventFrom,
      eventTo: range?.eventTo,
      status: range ? undefined : status,
    };
  }

  if (status === "approved") {
    return {
      purpose,
      status,
      event: range ? "expected" : undefined,
      eventFrom: range?.eventFrom,
      eventTo: range?.eventTo,
    };
  }

  if (status === "waiting_approval") {
    return {
      purpose,
      status,
      event: range ? "created" : undefined,
      eventFrom: range?.eventFrom,
      eventTo: range?.eventTo,
    };
  }

  if (status === "rejected") {
    return {
      purpose,
      status,
      event: range ? "created" : undefined,
      eventFrom: range?.eventFrom,
      eventTo: range?.eventTo,
    };
  }

  return {
    purpose,
    event: range ? "activity" : undefined,
    eventFrom: range?.eventFrom,
    eventTo: range?.eventTo,
  };
}

export function getVisitorFilterSummary(input: {
  datePreset: DateRangePreset;
  purpose?: ModelsVisitorPurpose;
  status?: ModelsVisitorStatus;
}) {
  const parts = [
    input.datePreset !== "today" ? getDateRangeLabel(input.datePreset) : null,
    input.purpose ? titleizeFilterValue(input.purpose) : null,
    input.status ? titleizeFilterValue(input.status) : null,
  ].filter(Boolean);

  return parts.length > 0 ? parts.join(" • ") : "";
}

function titleizeFilterValue(value: string) {
  return value
    .replace(/_/g, " ")
    .split(" ")
    .filter(Boolean)
    .map((part) => `${part.charAt(0).toUpperCase()}${part.slice(1)}`)
    .join(" ");
}

export function getDateRange(preset: DateRangePreset) {
  const now = new Date();

  switch (preset) {
    case "today":
      return {
        createdFrom: startOfDay(now).toISOString(),
        createdTo: endOfDay(now).toISOString(),
      };
    case "yesterday": {
      const day = new Date(now);
      day.setDate(day.getDate() - 1);
      return {
        createdFrom: startOfDay(day).toISOString(),
        createdTo: endOfDay(day).toISOString(),
      };
    }
    case "last_7_days": {
      const start = new Date(now);
      start.setDate(start.getDate() - 6);
      return {
        createdFrom: startOfDay(start).toISOString(),
        createdTo: endOfDay(now).toISOString(),
      };
    }
    case "this_month": {
      const start = new Date(now.getFullYear(), now.getMonth(), 1);
      return {
        createdFrom: startOfDay(start).toISOString(),
        createdTo: endOfDay(now).toISOString(),
      };
    }
    default:
      return { createdFrom: undefined, createdTo: undefined };
  }
}

export function getDateRangeLabel(preset: DateRangePreset) {
  switch (preset) {
    case "today":
      return "Today";
    case "yesterday":
      return "Yesterday";
    case "last_7_days":
      return "Last 7 Days";
    case "this_month":
      return "This Month";
    default:
      return "All Time";
  }
}

export function getLogsSegmentSummaryTitle(segment: LogsSegment) {
  switch (segment) {
    case "today":
      return "Entries with activity today";
    case "expected":
      return "Expected guest arrivals";
    case "inside":
      return "Visitors currently inside";
    case "all":
    default:
      return "All visitor entries";
  }
}

export const DATE_RANGE_OPTIONS: { label: string; value: DateRangePreset }[] = [
  { label: "Today", value: "today" },
  { label: "Yesterday", value: "yesterday" },
  { label: "Last 7 Days", value: "last_7_days" },
  { label: "This Month", value: "this_month" },
  { label: "All Time", value: "all" },
];

export function getTodayRange() {
  return getHalfOpenRangeIST("today");
}
