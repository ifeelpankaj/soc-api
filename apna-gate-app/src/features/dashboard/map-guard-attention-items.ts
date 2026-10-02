import type { AttentionItem } from "@/components/dashboard";

export type GuardAttentionInput = {
  pendingCount: number;
  readyCount: number;
  onReviewPending: () => void;
  onCheckIn: () => void;
};

export function mapGuardAttentionItems(input: GuardAttentionInput): AttentionItem[] {
  const items: AttentionItem[] = [];

  if (input.pendingCount > 0) {
    items.push({
      id: "guard-pending-approval",
      icon: { ios: "person.fill", android: "person", web: "person" },
      iconTone: "orange",
      title:
        input.pendingCount === 1
          ? "1 visitor waiting for approval"
          : `${input.pendingCount} visitors waiting for approval`,
      subtitle: "Requires your review",
      actionLabel: "Review →",
      onPress: input.onReviewPending,
    });
  }

  if (input.readyCount > 0) {
    items.push({
      id: "guard-ready-check-in",
      icon: {
        ios: "door.left.hand.open",
        android: "meeting_room",
        web: "meeting_room",
      },
      iconTone: "blue",
      title:
        input.readyCount === 1
          ? "1 visitor ready to check in"
          : `${input.readyCount} visitors ready to check in`,
      subtitle: "Expected visitor",
      actionLabel: "Check In →",
      onPress: input.onCheckIn,
    });
  }

  return items;
}
