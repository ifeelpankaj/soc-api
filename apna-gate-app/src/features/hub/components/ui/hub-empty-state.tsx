import { EmptyState } from "@/components/ui";

type HubEmptyStateProps = {
  title: string;
  message: string;
  actionLabel?: string;
  onAction?: () => void;
};

export function HubEmptyState({
  title,
  message,
  actionLabel,
  onAction,
}: HubEmptyStateProps) {
  return (
    <EmptyState
      actionLabel={actionLabel}
      message={message}
      title={title}
      onAction={onAction}
    />
  );
}
