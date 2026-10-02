import type { AttentionItem } from "@/components/dashboard";
import type { HubContent } from "@/lib/api/hub-api";
import type { MaintenanceBill } from "@/lib/api/maintenance-api";
import {
  formatBillingMonth,
  formatMoney,
  getMaintenanceBillPresentation,
} from "@/features/resident/maintenance/maintenance-presentation";
import { formatDateOnly } from "@/features/guard/guard-utils";

export type ResidentAttentionInput = {
  pendingCount: number;
  currentBill?: MaintenanceBill | null;
  importantAnnouncement?: {
    channelId: number;
    post: HubContent;
  } | null;
  onReviewVisitors: () => void;
  onOpenMaintenanceBill: (billId: number) => void;
  onOpenMaintenancePay: (billId: number) => void;
  onOpenAnnouncement: (postId: number, channelId: number) => void;
};

function maintenanceAttentionItem(
  bill: MaintenanceBill,
  input: ResidentAttentionInput,
): AttentionItem | null {
  const presentation = getMaintenanceBillPresentation({ bill });
  if (presentation.state === "PENDING_REVIEW") {
    return null;
  }

  const amountLine = `${formatMoney(bill.total_paise)} · ${formatBillingMonth(bill.billing_month)}`;
  const dueLine = bill.due_message
    ? bill.due_message
    : bill.due_date
      ? `Due ${formatDateOnly(bill.due_date, bill.timezone)}`
      : undefined;

  if (presentation.state === "REJECTED") {
    return {
      id: `maintenance-rejected-${bill.id}`,
      icon: { ios: "indianrupeesign.circle", android: "currency_rupee", web: "currency_rupee" },
      iconTone: "orange",
      title: "Maintenance payment rejected",
      subtitle: amountLine,
      actionLabel: "Review →",
      onPress: () => input.onOpenMaintenanceBill(bill.id),
    };
  }

  if (presentation.state === "OVERDUE") {
    return {
      id: `maintenance-overdue-${bill.id}`,
      icon: { ios: "indianrupeesign.circle", android: "currency_rupee", web: "currency_rupee" },
      iconTone: "orange",
      title: "Maintenance overdue",
      subtitle: amountLine,
      actionLabel: "Pay Now →",
      onPress: () => input.onOpenMaintenancePay(bill.id),
    };
  }

  if (presentation.state === "UNPAID") {
    return {
      id: `maintenance-due-${bill.id}`,
      icon: { ios: "indianrupeesign.circle", android: "currency_rupee", web: "currency_rupee" },
      iconTone: "teal",
      title: "Maintenance due soon",
      subtitle: dueLine ? `${amountLine.split(" · ")[0]} · ${dueLine}` : amountLine,
      actionLabel: "Pay Now →",
      onPress: () => input.onOpenMaintenancePay(bill.id),
    };
  }

  return null;
}

function announcementAttentionItem(
  input: ResidentAttentionInput,
): AttentionItem | null {
  const payload = input.importantAnnouncement;
  if (!payload?.post?.is_important) {
    return null;
  }

  const postId = payload.post.post_id ?? payload.post.id;
  const title = payload.post.title?.trim() || "Important announcement";
  const schedule =
    payload.post.body.trim().split("\n")[0]?.slice(0, 80) ||
    "Official society update";

  return {
    id: `announcement-important-${postId}`,
    icon: { ios: "megaphone.fill", android: "campaign", web: "campaign" },
    iconTone: "purple",
    title,
    subtitle: schedule,
    actionLabel: "View →",
    onPress: () => input.onOpenAnnouncement(postId, payload.channelId),
  };
}

export function mapResidentAttentionItems(input: ResidentAttentionInput): AttentionItem[] {
  const items: AttentionItem[] = [];

  if (input.pendingCount > 0) {
    items.push({
      id: "resident-pending-visitors",
      icon: { ios: "person.fill", android: "person", web: "person" },
      iconTone: "orange",
      title:
        input.pendingCount === 1
          ? "1 visitor awaiting approval"
          : `${input.pendingCount} visitors awaiting approval`,
      subtitle: "Waiting for your action",
      actionLabel: "Review →",
      onPress: input.onReviewVisitors,
    });
  }

  if (input.currentBill) {
    const maintenanceItem = maintenanceAttentionItem(input.currentBill, input);
    if (maintenanceItem) {
      items.push(maintenanceItem);
    }
  }

  const announcementItem = announcementAttentionItem(input);
  if (announcementItem) {
    items.push(announcementItem);
  }

  return items;
}
