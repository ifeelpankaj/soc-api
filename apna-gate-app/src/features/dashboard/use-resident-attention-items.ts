import { useRouter } from "expo-router";
import { useCallback, useMemo } from "react";

import type { AttentionItem } from "@/components/dashboard";
import { mapResidentAttentionItems } from "@/features/dashboard/map-resident-attention-items";
import { useResidentDashboard } from "@/features/resident/hooks/use-resident-dashboard";
import { residentHubPostRoute } from "@/features/hub/hub-routes";
import {
  residentVisitorsRoute,
  residentMaintenanceBillRoute,
  residentMaintenancePaymentRoute,
} from "@/features/resident/resident-routes";
import { useResident } from "@/features/resident/resident-context";
import { findHubChannel, useHubChannelsQuery } from "@/lib/api/hub-api";
import { useMaintenanceOutstandingQuery } from "@/lib/api/maintenance-api";

export function useResidentAttentionItems(): AttentionItem[] {
  const router = useRouter();
  const { societyId, flatId } = useResident();
  const dashboard = useResidentDashboard();
  const maintenanceQuery = useMaintenanceOutstandingQuery(
    { societyId: societyId ?? 0, flatId: flatId ?? 0 },
    { skip: !societyId || !flatId, refetchOnMountOrArgChange: true },
  );
  const channelsQuery = useHubChannelsQuery(societyId ?? 0, {
    skip: !societyId,
    refetchOnMountOrArgChange: true,
  });

  const onReviewVisitors = useCallback(
    () => router.push(residentVisitorsRoute()),
    [router],
  );

  const onOpenMaintenanceBill = useCallback(
    (billId: number) => {
      if (!societyId || !flatId) {
        return;
      }
      router.push(residentMaintenanceBillRoute(billId, societyId, flatId));
    },
    [flatId, router, societyId],
  );

  const onOpenMaintenancePay = useCallback(
    (billId: number) => {
      if (!societyId || !flatId) {
        return;
      }
      router.push(residentMaintenancePaymentRoute(billId, societyId, flatId));
    },
    [flatId, router, societyId],
  );

  const onOpenAnnouncement = useCallback(
    (postId: number, channelId: number) => {
      if (!societyId) {
        return;
      }
      router.push(residentHubPostRoute(societyId, postId, channelId, "home"));
    },
    [router, societyId],
  );

  const importantAnnouncement = useMemo(() => {
    const announcementChannel = findHubChannel(channelsQuery.data, "announcement");
    const post = announcementChannel?.important_announcement;
    if (!announcementChannel?.id || !post) {
      return null;
    }
    return { channelId: announcementChannel.id, post };
  }, [channelsQuery.data]);

  return useMemo(
    () =>
      mapResidentAttentionItems({
        pendingCount: dashboard.pendingCount,
        currentBill: maintenanceQuery.data?.current_bill,
        importantAnnouncement,
        onReviewVisitors,
        onOpenMaintenanceBill,
        onOpenMaintenancePay,
        onOpenAnnouncement,
      }),
    [
      dashboard.pendingCount,
      importantAnnouncement,
      maintenanceQuery.data?.current_bill,
      onOpenAnnouncement,
      onOpenMaintenanceBill,
      onOpenMaintenancePay,
      onReviewVisitors,
    ],
  );
}
