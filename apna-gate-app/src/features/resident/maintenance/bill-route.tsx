import { useLocalSearchParams } from "expo-router";
import { useResident } from "@/features/resident/resident-context";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import { residentMaintenanceRoute } from "@/features/resident/resident-routes";
import { EmptyState } from "@/components/ui";
import type { BillArgs } from "@/lib/api/maintenance-api";
import type { ComponentType } from "react";

export function BillRoute({
  screen: Screen,
}: {
  screen: ComponentType<BillArgs>;
}) {
  const params = useLocalSearchParams<{
    billId: string;
    societyId?: string;
    flatId?: string;
  }>();
  const { societyId, flatId } = useResident();
  const billId = Number(params.billId);
  const valid =
    Number.isSafeInteger(billId) &&
    billId > 0 &&
    societyId &&
    flatId &&
    (!params.societyId || Number(params.societyId) === societyId) &&
    (!params.flatId || Number(params.flatId) === flatId);
  if (!valid)
    return (
      <ResidentSubScreen
        title="Maintenance"
        fallbackHomeRoute={residentMaintenanceRoute()}
      >
        <EmptyState
          title="Bill unavailable"
          message="Open a bill from Maintenance for your selected residence."
        />
      </ResidentSubScreen>
    );
  return (
    <Screen
      key={`${societyId}:${flatId}:${billId}`}
      societyId={societyId}
      flatId={flatId}
      billId={billId}
    />
  );
}
