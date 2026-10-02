import { useCallback, useMemo, useState } from "react";

import { getApiMessage } from "@/features/auth/api-error";
import { formatSelectedFlatLabel, type SelectedFlat } from "@/features/guard/hooks/use-guard-manual-entry";
import {
  type ModelsShortLinkResponse,
  type ModelsVisitorPurpose,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaffMutation,
} from "@/lib/api/generated-api";

type CreatedInviteState = {
  purpose: ModelsVisitorPurpose;
  expiresAt?: string;
  flat?: SelectedFlat | null;
  shortLink?: ModelsShortLinkResponse | null;
};

export function useGuardVisitorInvite(societyId: number) {
  const [selectedFlat, setSelectedFlat] = useState<SelectedFlat | null>(null);
  const [purpose, setPurpose] = useState<ModelsVisitorPurpose>("guest");
  const [createdInvite, setCreatedInvite] = useState<CreatedInviteState | null>(null);
  const [flatError, setFlatError] = useState<string>();

  const [createInvite, createInviteState] =
    usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesStaffMutation();

  const isFormValid = useMemo(
    () => Boolean(selectedFlat?.id) && Boolean(purpose),
    [purpose, selectedFlat?.id],
  );

  const resetForm = useCallback(() => {
    setSelectedFlat(null);
    setPurpose("guest");
    setFlatError(undefined);
  }, []);

  const clearCreatedInvite = useCallback(() => {
    setCreatedInvite(null);
    resetForm();
  }, [resetForm]);

  const submit = useCallback(async () => {
    if (!selectedFlat?.id) {
      setFlatError("Select a visiting flat");
      return { success: false as const, message: "Select a visiting flat." };
    }
    setFlatError(undefined);

    try {
      const response = await createInvite({
        flatId: selectedFlat.id,
        societyId,
        modelsCreateVisitorInviteRequest: { purpose },
      }).unwrap();

      const shortLink = response.data?.short_link ?? response.data?.link;
      if (!shortLink?.url) {
        return {
          success: false as const,
          message: response.message ?? "Invite created but share link was unavailable.",
        };
      }

      setCreatedInvite({
        purpose,
        expiresAt: response.data?.token?.expires_at,
        flat: selectedFlat,
        shortLink,
      });
      resetForm();

      return {
        success: true as const,
        message: response.message ?? "Visitor form link created",
      };
    } catch (error) {
      return {
        success: false as const,
        message: getApiMessage(error, "Please try again."),
      };
    }
  }, [createInvite, purpose, resetForm, selectedFlat, societyId]);

  return {
    clearCreatedInvite,
    createInviteState,
    createdInvite,
    flatError,
    flatLabel: formatSelectedFlatLabel(createdInvite?.flat ?? selectedFlat),
    isFormValid,
    purpose,
    selectedFlat,
    setPurpose,
    setSelectedFlat,
    submit,
  };
}
