import { useResident } from "@/features/resident/resident-context";
import {
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesQuery,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdMembersQuery,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesMutation,
} from "@/lib/api/resident-api-extensions";

export function useResidentMembers() {
  const { canManageFlatMembers, flatId, societyId } = useResident();
  const shouldSkip = !societyId || !flatId;

  const membersQuery = useGetV1SocietiesBySocietyIdFlatsAndFlatIdMembersQuery(
    { societyId: societyId ?? 0, flatId: flatId ?? 0, limit: 50, offset: 0 },
    { skip: shouldSkip },
  );

  const invitesQuery = useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesQuery(
    { societyId: societyId ?? 0, flatId: flatId ?? 0 },
    { skip: shouldSkip || !canManageFlatMembers },
  );

  const [createInvite, createInviteState] =
    usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesMutation();

  const refetchAll = () => {
    if (!membersQuery.isUninitialized) {
      void membersQuery.refetch();
    }
    if (canManageFlatMembers && !invitesQuery.isUninitialized) {
      void invitesQuery.refetch();
    }
  };

  return {
    canManageFlatMembers,
    createInvite,
    createInviteState,
    invites: invitesQuery.data?.data?.invites ?? [],
    isCreating: createInviteState.isLoading,
    isLoading:
      shouldSkip ||
      (membersQuery.isLoading && !membersQuery.data) ||
      (canManageFlatMembers && invitesQuery.isLoading && !invitesQuery.data),
    isRefreshing: membersQuery.isFetching || invitesQuery.isFetching,
    members: membersQuery.data?.data?.residents ?? [],
    refetchAll,
  };
}
