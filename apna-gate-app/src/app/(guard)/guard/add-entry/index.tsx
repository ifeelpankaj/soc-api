import { useRouter } from "expo-router";

import { GuardSubScreen } from "@/features/guard/components/guard-sub-screen";
import { ManualEntryForm } from "@/features/guard/components/manual-entry/manual-entry-form";
import { guardEntryDetailRoute } from "@/features/guard/guard-routes";
import { useGuardDashboard } from "@/features/guard/hooks/use-guard-dashboard";
import { useGuardFeedback } from "@/features/guard/hooks/use-guard-feedback";
import { useGuardScreen } from "@/features/guard/hooks/use-guard-screen";

export default function GuardAddEntryScreen() {
  const router = useRouter();
  const feedback = useGuardFeedback();
  const { selectedMembership, selectedSocietyId } = useGuardScreen();
  const { visitorSettings } = useGuardDashboard();

  return (
    <ManualEntryForm
      renderLayout={(content, footer) => (
        <GuardSubScreen title="Add Entry" footer={footer}>{content}</GuardSubScreen>
      )}
      allowGuardEntry={visitorSettings?.allow_guard_entry !== false}
      societyId={selectedSocietyId ?? 0}
      societyName={
        selectedMembership
          ? `Society #${selectedMembership.society_id}`
          : undefined
      }
      onEntryCreated={({ entry }) => {
        if (!entry?.id) {
          return;
        }
        router.replace(guardEntryDetailRoute(entry.id));
      }}
      onError={(message) =>
        feedback.showActionResult(
          { success: false, message },
          {
            errorTitle: "Could not create entry",
            successTitle: "Entry created",
          },
        )
      }
      onSuccess={(message) =>
        feedback.showSuccess("Entry created", message)
      }
    />
  );
}
