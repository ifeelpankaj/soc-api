import { useCallback, useMemo, useState } from "react";

import { getApiMessage, isSubscriptionError } from "@/features/auth/api-error";
import { useGuardFeedback } from "@/features/guard/hooks/use-guard-feedback";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";
import type { GuardDeskVisitorSettings } from "@/lib/api/guard-api-extensions";
import {
  useGetV1SocietiesBySocietyIdGuardDeskBootstrapQuery,
  usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutMutation,
} from "@/lib/api/generated-api";

export function useGuardDashboard() {
  const { selectedSocietyId, societyName } = useGuardScreen();
  const { showError, showSuccess } = useGuardFeedback();
  const [checkoutEntryId, setCheckoutEntryId] = useState<number | null>(null);

  const shouldSkip = !selectedSocietyId;
  const queryOpts = {
    skip: shouldSkip,
    refetchOnFocus: true,
    refetchOnReconnect: true,
  };

  const bootstrapQuery = useGetV1SocietiesBySocietyIdGuardDeskBootstrapQuery(
    { societyId: selectedSocietyId ?? 0 },
    queryOpts,
  );
  const [checkOut] =
    usePostV1SocietiesBySocietyIdVisitorEntriesAndEntryIdCheckOutMutation();

  const desk = bootstrapQuery.data?.data?.desk;

  const refetchAll = useCallback(() => {
    void bootstrapQuery.refetch();
  }, [bootstrapQuery]);

  const checkOutEntry = useCallback(
    async (entryId?: number) => {
      if (!entryId || !selectedSocietyId) {
        return;
      }

      setCheckoutEntryId(entryId);

      try {
        await checkOut({ societyId: selectedSocietyId, entryId }).unwrap();
        showSuccess("Checked out", "Visitor has left the society.");
        void bootstrapQuery.refetch();
      } catch (error) {
        showError("Checkout failed", error, "Please try again.");
      } finally {
        setCheckoutEntryId(null);
      }
    },
    [bootstrapQuery, checkOut, selectedSocietyId, showError, showSuccess],
  );

  const stats = desk?.stats;
  const pendingEntries = desk?.pending_preview ?? [];
  const visitorSettings = (desk as { visitor_settings?: GuardDeskVisitorSettings } | undefined)
    ?.visitor_settings;

  const resolvedSocietyName = useMemo(
    () => desk?.society?.name ?? societyName ?? `Society #${selectedSocietyId}`,
    [desk?.society?.name, selectedSocietyId, societyName],
  );

  const isInitialLoading =
    shouldSkip || (bootstrapQuery.isLoading && !bootstrapQuery.data);

  const bootstrapError = bootstrapQuery.isError ? bootstrapQuery.error : null;
  const isSubscriptionBlocked = isSubscriptionError(bootstrapError);
  const hasError = bootstrapQuery.isError && !isSubscriptionBlocked;

  return {
    checkOutEntry,
    checkoutEntryId,
    errorMessage: hasError
      ? getApiMessage(bootstrapError, "Unable to load guard desk data.")
      : null,
    expectedGuestsCount: desk?.expected_guests_count ?? 0,
    hasError,
    isSubscriptionBlocked,
    isInitialLoading,
    isRefreshing: bootstrapQuery.isFetching && !bootstrapQuery.isLoading,
    pendingEntries,
    refetchAll,
    societyName: resolvedSocietyName,
    stats,
    visitorSettings,
  };
}
