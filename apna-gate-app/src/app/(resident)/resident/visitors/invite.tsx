import { useRouter } from "expo-router";
import { useRef, useState } from "react";
import { SymbolView } from "expo-symbols";
import { Pressable, ScrollView, StyleSheet, Text, View } from "react-native";

import { Button, Card } from "@/components/ui";
import { BottomActions } from "@/components/ui/bottom-actions";
import { invitePurposePresentation } from "@/features/resident/invites/invite-purpose-presentation";
import { getApiMessage } from "@/features/auth/api-error";
import {
  titleize,
  residentInvitePurposes,
} from "@/features/visitors/visitor-utils";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import {
  buildVisitorInviteShareContent,
  copyInviteLink,
  type InviteShareContent,
  shareInviteNative,
  shareInviteOnTelegram,
  shareInviteOnWhatsApp,
} from "@/features/resident/invites/invite-share";
import { useResident } from "@/features/resident/resident-context";
import {
  residentDashboardRoute,
  residentInviteDetailRoute,
} from "@/features/resident/resident-routes";
import {
  type ModelsShortLinkResponse,
  type ModelsVisitorPurpose,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesMutation,
} from "@/lib/api/generated-api";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type CreatedInvite = {
  purpose: ModelsVisitorPurpose;
  shortLink?: ModelsShortLinkResponse | null;
  expiresAt?: string;
  flatLabel?: string;
  societyName?: string | null;
};

export default function ResidentInviteScreen() {
  const { societyId, flatId } = useResident();
  return <ResidentInviteForm key={`${societyId}:${flatId}`} />;
}

function ResidentInviteForm() {
  const router = useRouter();
  const handleBack = useBackAction(residentDashboardRoute());
  const feedback = useResidentFeedback();
  const { flatId, canManageFlatVisitors, selectedResidence, societyId } =
    useResident();
  const [purpose, setPurpose] = useState<ModelsVisitorPurpose>("guest");
  const [createdInvite, setCreatedInvite] = useState<CreatedInvite | null>(
    null,
  );
  const submitting = useRef(false);
  const [createInvite, createInviteState] =
    usePostV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesMutation();

  if (!canManageFlatVisitors) {
    return (
      <ResidentSubScreen title="Visitor Invite">
        <ScrollView contentContainerStyle={styles.scrollContent}>
          <Card style={styles.restrictedCard}>
            <Text style={styles.restrictedTitle}>Visitor access required</Text>
            <Text style={styles.restrictedBody}>
              Active flat residents with visitor management access can create
              pre-approved visitor invites.
            </Text>
          </Card>
        </ScrollView>
      </ResidentSubScreen>
    );
  }

  const handleCreateInvite = async () => {
    if (!societyId || !flatId || !canManageFlatVisitors || submitting.current) {
      return;
    }
    submitting.current = true;
    try {
      const response = await createInvite({
        societyId,
        flatId,
        modelsCreateVisitorInviteRequest: { purpose },
      }).unwrap();

      const inviteId = response.data?.invite?.id;
      if (inviteId) {
        feedback.showSuccess(
          "Invite ready",
          "Share the invite link from the detail page.",
        );
        router.replace(residentInviteDetailRoute("visitor", inviteId));
        return;
      }

      const shortLink = response.data?.short_link ?? response.data?.link;
      if (!shortLink?.url) {
        feedback.showSuccess("Invite created", "Share link unavailable.");
        handleBack();
        return;
      }

      setCreatedInvite({
        purpose,
        flatLabel: selectedResidence?.flat_number ?? undefined,
        shortLink,
        societyName: selectedResidence?.society_name,
        expiresAt: response.data?.token?.expires_at,
      });
      feedback.showSuccess("Invite ready", "Share the link with your visitor.");
    } catch (error) {
      feedback.showError(
        "Invite failed",
        error,
        getApiMessage(error, "Please try again."),
      );
    } finally {
      submitting.current = false;
    }
  };

  const shareResult = createdInvite
    ? buildVisitorInviteShareContent(createdInvite)
    : { available: false as const, reason: "short_link_unavailable" as const };

  const handleShare = async (
    shareFn: (content: InviteShareContent) => Promise<void>,
    errorLabel: string,
  ) => {
    if (!shareResult.available) {
      feedback.showInfo(
        "Share link unavailable",
        "This invite does not have a short link yet.",
      );
      return;
    }

    try {
      await shareFn(shareResult.content);
    } catch {
      feedback.showError("Share failed", errorLabel, "Please try again.");
    }
  };

  const handleCopyLink = async () => {
    if (!shareResult.available) {
      feedback.showInfo(
        "Share link unavailable",
        "This invite does not have a short link yet.",
      );
      return;
    }

    try {
      const copied = await copyInviteLink(shareResult.content);
      feedback.showSuccess(
        copied ? "Link copied" : "Share link",
        copied
          ? "Visitor form link copied to clipboard."
          : "Use the share sheet to copy the visitor form link.",
      );
    } catch {
      feedback.showError(
        "Copy failed",
        "Unable to copy the visitor form link.",
        "Please try again.",
      );
    }
  };

  if (createdInvite) {
    const formUrl = shareResult.available
      ? shareResult.content.url
      : "Share link unavailable";

    return (
      <ResidentSubScreen title="Visitor Invite">
        <ScrollView contentContainerStyle={styles.scrollContent}>
          <View style={styles.content}>
            <View style={styles.intro}>
              <Text style={styles.pageTitle}>Invite ready</Text>
              <Text style={styles.pageSubtitle}>
                Share this form link with your visitor so they can complete
                entry details on web.
              </Text>
            </View>

            <Card style={styles.detailsCard}>
              <View style={styles.fieldGroup}>
                <Text style={styles.fieldLabel}>Purpose</Text>
                <Text style={styles.fieldValue}>
                  {titleize(createdInvite.purpose)}
                </Text>
              </View>

              <View style={styles.fieldGroup}>
                <Text style={styles.fieldLabel}>Form link</Text>
                <Text selectable style={styles.codeBlock}>
                  {formUrl}
                </Text>
              </View>

              {createdInvite.expiresAt ? (
                <Text style={styles.expiresText}>
                  Expires {new Date(createdInvite.expiresAt).toLocaleString()}
                </Text>
              ) : null}
            </Card>

            <Button
              disabled={!shareResult.available}
              title="Share on WhatsApp"
              onPress={() =>
                void handleShare(
                  shareInviteOnWhatsApp,
                  "Unable to open WhatsApp.",
                )
              }
            />
            <Button
              disabled={!shareResult.available}
              title="Share on Telegram"
              variant="secondary"
              onPress={() =>
                void handleShare(
                  shareInviteOnTelegram,
                  "Unable to open Telegram.",
                )
              }
            />
            <Button
              disabled={!shareResult.available}
              title="More options"
              variant="secondary"
              onPress={() =>
                void handleShare(
                  shareInviteNative,
                  "Unable to open the share sheet.",
                )
              }
            />
            <Button
              disabled={!shareResult.available}
              title="Copy link"
              variant="secondary"
              onPress={() => void handleCopyLink()}
            />
            <Button title="Done" variant="secondary" onPress={handleBack} />
          </View>
        </ScrollView>
      </ResidentSubScreen>
    );
  }

  return (
    <ResidentSubScreen
      title="Invite a visitor"
      footer={
        <BottomActions>
          <Button
            title="Create invite"
            disabled={!societyId || !flatId}
            loading={createInviteState.isLoading}
            onPress={handleCreateInvite}
          />
        </BottomActions>
      }
    >
      <ScrollView
        style={{ flex: 1 }}
        contentContainerStyle={styles.scrollContent}
      >
        <View style={styles.content}>
          <View style={styles.intro}>
            <View style={styles.residenceBadge}>
              <SymbolView
                name={{ ios: "house.fill", android: "home", web: "home" }}
                size={20}
                tintColor={colors.brand.orange}
              />
              <View style={{ flex: 1 }}>
                <Text style={styles.residenceTitle}>
                  {selectedResidence?.flat_number
                    ? `Flat ${selectedResidence.flat_number}`
                    : "Your home"}
                </Text>
                {selectedResidence?.society_name ? (
                  <Text style={styles.pageSubtitle}>
                    {selectedResidence.society_name}
                  </Text>
                ) : null}
              </View>
            </View>
            <Text style={styles.pageTitle}>Who are you welcoming?</Text>
            <Text style={styles.pageSubtitle}>
              Choose a purpose. We&apos;ll create a link for your visitor.
            </Text>
          </View>

          <View style={styles.purposeCard}>
            <View style={styles.purposeOptions}>
              {residentInvitePurposes.map((option) => {
                const active = option === purpose;
                const presentation = invitePurposePresentation[option];

                return (
                  <Pressable
                    key={option}
                    accessibilityRole="radio"
                    accessibilityState={{
                      checked: active,
                      disabled: createInviteState.isLoading,
                    }}
                    disabled={createInviteState.isLoading}
                    onPress={() => setPurpose(option)}
                    style={[
                      styles.purposeChip,
                      active && styles.purposeChipActive,
                    ]}
                  >
                    <View
                      style={[
                        styles.purposeIcon,
                        active && styles.purposeIconActive,
                      ]}
                    >
                      <SymbolView
                        name={presentation.icon}
                        size={23}
                        tintColor={
                          active ? colors.brand.orange : colors.text.secondary
                        }
                      />
                    </View>
                    <View style={styles.purposeCopy}>
                      <Text
                        style={[
                          styles.purposeChipText,
                          active && styles.purposeChipTextActive,
                        ]}
                      >
                        {titleize(option)}
                      </Text>
                      <Text style={styles.pageSubtitle}>
                        {presentation.description}
                      </Text>
                    </View>
                    <SymbolView
                      name={{
                        ios: active ? "checkmark.circle.fill" : "circle",
                        android: active
                          ? "check_circle"
                          : "radio_button_unchecked",
                        web: active ? "check_circle" : "radio_button_unchecked",
                      }}
                      size={22}
                      tintColor={
                        active ? colors.brand.orange : colors.border.default
                      }
                    />
                  </Pressable>
                );
              })}
            </View>
          </View>
          <View style={styles.guide}>
            <Text style={styles.guideTitle}>An easier welcome</Text>
            <Text style={styles.pageSubtitle}>
              Create your link, share it with your visitor, and let them fill in
              their details before arriving.
            </Text>
          </View>
        </View>
      </ScrollView>
    </ResidentSubScreen>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: spacing["2xl"],
  },
  scrollContent: {
    paddingTop: spacing.lg,
    paddingBottom: layout.screenPaddingBottom,
    paddingHorizontal: layout.screenPaddingHorizontal,
  },
  restrictedCard: {
    gap: spacing.sm,
  },
  restrictedTitle: {
    ...typography.body,
    color: colors.text.primary,
    fontWeight: "700",
  },
  restrictedBody: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  intro: {
    gap: spacing.xs,
  },
  pageTitle: {
    ...typography.title,
    color: colors.text.primary,
  },
  pageSubtitle: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  detailsCard: {
    gap: spacing.lg,
  },
  purposeCard: {
    gap: spacing.lg,
  },
  fieldGroup: {
    gap: spacing.xs,
  },
  fieldLabel: {
    ...typography.eyebrow,
    color: colors.text.muted,
  },
  fieldValue: {
    ...typography.subtitle,
    color: colors.text.primary,
    fontWeight: "600",
    textTransform: "capitalize",
  },
  codeBlock: {
    ...typography.bodySmall,
    backgroundColor: colors.surface.screen,
    borderRadius: radius.md,
    color: colors.text.primary,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  expiresText: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  purposeOptions: {
    gap: spacing.md,
  },
  purposeChip: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.md,
    minHeight: 80,
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
  },
  purposeChipActive: {
    backgroundColor: colors.brand.orangeSoft,
    borderColor: colors.brand.orange,
  },
  purposeChipText: {
    ...typography.body,
    color: colors.text.secondary,
    fontWeight: "600",
    textTransform: "capitalize",
  },
  purposeChipTextActive: {
    color: colors.brand.navy,
  },
  residenceBadge: {
    flexDirection: "row",
    alignItems: "center",
    gap: spacing.md,
    marginBottom: spacing.lg,
  },
  residenceTitle: {
    ...typography.bodySmall,
    fontWeight: "600",
    color: colors.text.primary,
  },
  purposeIcon: {
    width: 44,
    height: 44,
    borderRadius: radius.md,
    backgroundColor: colors.surface.muted,
    alignItems: "center",
    justifyContent: "center",
  },
  purposeIconActive: { backgroundColor: colors.surface.card },
  purposeCopy: { flex: 1, gap: spacing.xs },
  guide: {
    gap: spacing.sm,
    paddingHorizontal: spacing.sm,
    paddingBottom: spacing.sm,
  },
  guideTitle: {
    ...typography.bodySmall,
    fontWeight: "600",
    color: colors.text.primary,
  },
});
