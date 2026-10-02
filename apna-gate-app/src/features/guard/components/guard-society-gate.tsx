import { useMemo } from "react";

import { useGuardSociety } from "@/features/guard/guard-context";
import { WorkspaceSelectionGate } from "@/features/shared/workspace-selection-gate";

export function GuardSocietyGate() {
  const { isLoading, memberships, refetch, selectSociety } = useGuardSociety();

  const items = useMemo(
    () =>
      memberships.map((membership) => ({
        id: `guard-society-${membership.id ?? membership.society_id}`,
        onSelect: () => {
          if (membership.society_id) {
            selectSociety(membership.society_id);
          }
        },
        primaryLabel: `Society #${membership.society_id}`,
        secondaryLabel: `${membership.role ?? "staff"} access`,
        selectLabel: "Operate this gate",
        status: membership.status,
      })),
    [memberships, selectSociety],
  );

  if (!isLoading && memberships.length === 0) {
    return (
      <WorkspaceSelectionGate
        emptyMessage="Your account is signed in, but active guard access is not linked yet."
        emptyTitle="No guard access"
        eyebrow="Guard workspace"
        isLoading={false}
        items={[]}
        subtitle="Choose the society gate you are operating before continuing."
        title="Select society"
        onRefresh={refetch}
      />
    );
  }

  return (
    <WorkspaceSelectionGate
      emptyMessage="Your account is signed in, but active guard access is not linked yet."
      emptyTitle="No guard access"
      eyebrow="Guard workspace"
      isLoading={isLoading}
      items={items}
      loadingMessage="Opening guard workspace"
      subtitle="Choose the society gate you are operating before continuing."
      title="Select society"
      onRefresh={refetch}
    />
  );
}
