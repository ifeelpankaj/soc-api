import { useMemo } from "react";

import { useResident } from "@/features/resident/resident-context";
import { WorkspaceSelectionGate } from "@/features/shared/workspace-selection-gate";

export function ResidentSocietyGate() {
  const { isLoading, refetch, residences, selectResidence } = useResident();

  const items = useMemo(
    () =>
      residences.map((residence) => ({
        id: `residence-${residence.id ?? residence.flat_id}`,
        onSelect: () => {
          if (residence.flat_id) {
            selectResidence(residence.flat_id);
          }
        },
        primaryLabel: residence.society_name ?? "Your society",
        secondaryLabel: `Flat ${residence.flat_number ?? "-"}${residence.block ? ` · Block ${residence.block}` : ""}`,
        selectLabel: "Open this flat",
        status: residence.status,
      })),
    [residences, selectResidence],
  );

  return (
    <WorkspaceSelectionGate
      emptyMessage="Your account is signed in, but a flat has not been linked yet."
      emptyTitle="No active residence"
      eyebrow="Resident workspace"
      isLoading={isLoading}
      items={items}
      loadingMessage="Opening resident workspace"
      subtitle="Choose the flat you want to manage before continuing."
      title="Select residence"
      useLoadingState
      onRefresh={refetch}
    />
  );
}
