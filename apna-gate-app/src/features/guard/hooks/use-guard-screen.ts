import { useGuardSociety } from "@/features/guard/guard-context";
import { titleize } from "@/features/guard/guard-utils";
import {
  useGetV1SocietiesBySocietyIdQuery,
} from "@/lib/api/generated-api";

export function useGuardScreen() {
  const {
    isLoading,
    memberships,
    refetch,
    requiresSelection,
    selectedMembership,
    selectedSocietyId,
    user,
  } = useGuardSociety();

  const societyQuery = useGetV1SocietiesBySocietyIdQuery(
    { societyId: selectedSocietyId ?? 0 },
    { skip: !selectedSocietyId },
  );

  const societyName =
    societyQuery.data?.data?.society?.name ??
    (selectedSocietyId ? `Society #${selectedSocietyId}` : undefined);
  const societyCode = societyQuery.data?.data?.society?.society_code;

  const isReady =
    !isLoading &&
    memberships.length > 0 &&
    !requiresSelection &&
    !!selectedSocietyId;

  return {
    isLoading,
    isReady,
    memberships,
    refetch,
    requiresSelection,
    roleLabel: titleize(selectedMembership?.role ?? "staff"),
    selectedMembership,
    selectedSocietyId,
    societyCode,
    societyName,
    user,
  };
}
