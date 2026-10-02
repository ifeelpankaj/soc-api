import type { ReactNode } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { formatDateOnly, formatTimeOfDay, titleize } from "@/features/guard/guard-utils";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import {
  buildMemberInviteShareContent,
  copyInviteLink,
  type InviteShareContent,
  shareInviteNative,
  shareInviteOnTelegram,
  shareInviteOnWhatsApp,
} from "@/features/resident/invites/invite-share";
import type { ModelsShortLinkResponse } from "@/lib/api/generated-api";
import type { ModelsFlatMemberInviteRole } from "@/lib/api/resident-api-extensions";
import { colors } from "@/theme/colors";
import { androidCompactText } from "@/theme/platform-styles";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

const CARD_BORDER = "rgba(16, 29, 54, 0.08)";

type MemberInviteSuccessProps = {
  flatNumber?: string;
  invite: {
    expiresAt?: string;
    fullName: string;
    role: ModelsFlatMemberInviteRole;
    shortLink?: ModelsShortLinkResponse | null;
  };
  onDone: () => void;
  societyName?: string;
};

function DetailRow({
  children,
  icon,
  iconBg,
  iconColor,
  label,
  value,
}: {
  children?: ReactNode;
  icon: SymbolViewProps["name"];
  iconBg: string;
  iconColor: string;
  label: string;
  value?: string;
}) {
  return (
    <View style={styles.detailRow}>
      <View style={[styles.detailIconWrap, { backgroundColor: iconBg }]}>
        <SymbolView name={icon} size={16} tintColor={iconColor} />
      </View>
      <View style={styles.detailCopy}>
        <Text style={styles.detailLabel}>{label}</Text>
        {children ?? (
          <Text numberOfLines={2} style={styles.detailValue}>
            {value}
          </Text>
        )}
      </View>
    </View>
  );
}

function ShareOption({
  disabled,
  icon,
  label,
  onPress,
  primary,
}: {
  disabled?: boolean;
  icon: SymbolViewProps["name"];
  label: string;
  onPress: () => void;
  primary?: boolean;
}) {
  return (
    <Pressable
      accessibilityRole="button"
      disabled={disabled}
      style={({ pressed }) => [
        styles.shareOption,
        primary ? styles.shareOptionPrimary : styles.shareOptionSecondary,
        pressed && styles.shareOptionPressed,
        disabled && styles.disabledAction,
      ]}
      onPress={onPress}
    >
      <View style={styles.shareOptionLeft}>
        <SymbolView
          name={icon}
          size={18}
          tintColor={primary ? colors.text.inverse : colors.brand.navy}
        />
        <Text style={[styles.shareOptionText, primary && styles.shareOptionTextPrimary]}>
          {label}
        </Text>
      </View>
      <SymbolView
        name={{ ios: "chevron.right", android: "chevron_right", web: "chevron_right" }}
        size={14}
        tintColor={primary ? colors.text.inverse : colors.guard.textMuted}
      />
    </Pressable>
  );
}

export function MemberInviteSuccess({
  flatNumber,
  invite,
  onDone,
  societyName,
}: MemberInviteSuccessProps) {
  const feedback = useResidentFeedback();
  const shareResult = buildMemberInviteShareContent({
    expiresAt: invite.expiresAt,
    flatLabel: flatNumber,
    fullName: invite.fullName,
    role: invite.role,
    shortLink: invite.shortLink,
    societyName,
  });
  const joinUrl = shareResult.available ? shareResult.content.url : "Share link unavailable";
  const expiryDate = invite.expiresAt ? formatDateOnly(invite.expiresAt) : null;
  const expiryTime = invite.expiresAt ? formatTimeOfDay(invite.expiresAt) : null;

  const handleShare = async (
    shareFn: (content: InviteShareContent) => Promise<void>,
    errorLabel: string,
  ) => {
    if (!shareResult.available) {
      feedback.showInfo("Share link unavailable", "This invite does not have a short link yet.");
      return;
    }

    try {
      await shareFn(shareResult.content);
    } catch {
      feedback.showError("Share failed", errorLabel, "Please try again.");
    }
  };

  const handleCopy = async () => {
    if (!shareResult.available) {
      feedback.showInfo("Share link unavailable", "This invite does not have a short link yet.");
      return;
    }

    try {
      const copied = await copyInviteLink(shareResult.content);
      feedback.showSuccess(
        copied ? "Link copied" : "Share link",
        copied ? "Join link copied to clipboard." : "Use the share sheet to copy the join link.",
      );
    } catch {
      feedback.showError("Copy failed", "Unable to copy the join link.", "Please try again.");
    }
  };

  return (
    <View style={styles.content}>
      <View style={styles.heroRow}>
        <View style={styles.heroCopy}>
          <View style={styles.successBadge}>
            <SymbolView
              name={{ ios: "checkmark", android: "check", web: "check" }}
              size={28}
              tintColor={colors.text.inverse}
            />
          </View>
          <View style={styles.heroText}>
            <Text style={styles.pageTitle}>Invite ready!</Text>
            <Text style={styles.pageSubtitle}>
              Share this link with {invite.fullName} so they can accept and join your flat.
            </Text>
          </View>
        </View>

        <View style={styles.heroArt}>
          <SymbolView
            name={{ ios: "envelope.open.fill", android: "mail", web: "mail" }}
            size={42}
            tintColor={colors.brand.orange}
          />
          <View style={styles.sparkleOne}>
            <SymbolView
              name={{ ios: "sparkle", android: "auto_awesome", web: "auto_awesome" }}
              size={12}
              tintColor={colors.brand.orange}
            />
          </View>
          <View style={styles.sparkleTwo}>
            <SymbolView
              name={{ ios: "sparkle", android: "auto_awesome", web: "auto_awesome" }}
              size={10}
              tintColor={colors.brand.orange}
            />
          </View>
        </View>
      </View>

      <View style={styles.detailsCard}>
        <DetailRow
          icon={{ ios: "person.fill", android: "person", web: "person" }}
          iconBg={colors.brand.orangeSoft}
          iconColor={colors.brand.orange}
          label="Member"
          value={invite.fullName}
        />

        <View style={styles.rowDivider} />

        <DetailRow
          icon={{ ios: "checkmark.shield.fill", android: "verified_user", web: "verified_user" }}
          iconBg={colors.status.successSoft}
          iconColor={colors.status.success}
          label="Role"
          value={titleize(invite.role)}
        />

        <View style={styles.rowDivider} />

        <DetailRow
          icon={{ ios: "link", android: "link", web: "link" }}
          iconBg="#F5F3FF"
          iconColor="#7C3AED"
          label="Join link"
        >
          <View style={styles.linkBox}>
            <Text numberOfLines={2} selectable={shareResult.available} style={styles.linkText}>
              {joinUrl}
            </Text>
            <Pressable
              accessibilityLabel="Copy join link"
              accessibilityRole="button"
              disabled={!shareResult.available}
              hitSlop={8}
              style={({ pressed }) => [
                styles.copyButton,
                pressed && styles.copyButtonPressed,
                !shareResult.available && styles.disabledAction,
              ]}
              onPress={() => void handleCopy()}
            >
              <SymbolView
                name={{ ios: "doc.on.doc.fill", android: "content_copy", web: "content_copy" }}
                size={16}
                tintColor={colors.brand.orange}
              />
            </Pressable>
          </View>
        </DetailRow>

        {expiryDate || expiryTime ? (
          <>
            <View style={styles.cardDivider} />
            <View style={styles.expiryRow}>
              {expiryDate ? (
                <View style={styles.expiryItem}>
                  <SymbolView
                    name={{ ios: "calendar", android: "calendar_today", web: "calendar_today" }}
                    size={13}
                    tintColor={colors.guard.textMuted}
                  />
                  <Text style={styles.expiryText}>Expires on {expiryDate}</Text>
                </View>
              ) : null}
              {expiryTime ? (
                <View style={styles.expiryItem}>
                  <SymbolView
                    name={{ ios: "clock", android: "schedule", web: "schedule" }}
                    size={13}
                    tintColor={colors.guard.textMuted}
                  />
                  <Text style={styles.expiryText}>{expiryTime}</Text>
                </View>
              ) : null}
            </View>
          </>
        ) : null}
      </View>

      <View style={styles.shareSection}>
        <Text style={styles.shareHeading}>Share via</Text>
        <View style={styles.shareList}>
          <ShareOption
            icon={{ ios: "message.fill", android: "chat", web: "chat" }}
            label="Share on WhatsApp"
            disabled={!shareResult.available}
            primary
            onPress={() => void handleShare(shareInviteOnWhatsApp, "Unable to open WhatsApp.")}
          />
          <ShareOption
            icon={{ ios: "paperplane.fill", android: "send", web: "send" }}
            label="Share on Telegram"
            disabled={!shareResult.available}
            onPress={() => void handleShare(shareInviteOnTelegram, "Unable to open Telegram.")}
          />
          <ShareOption
            icon={{ ios: "link", android: "link", web: "link" }}
            label="Copy link"
            disabled={!shareResult.available}
            onPress={() => void handleCopy()}
          />
          <ShareOption
            icon={{ ios: "ellipsis", android: "more_horiz", web: "more_horiz" }}
            label="More options"
            disabled={!shareResult.available}
            onPress={() => void handleShare(shareInviteNative, "Unable to open the share sheet.")}
          />
        </View>
      </View>

      <Pressable
        accessibilityRole="button"
        style={({ pressed }) => [styles.doneButton, pressed && styles.doneButtonPressed]}
        onPress={onDone}
      >
        <Text style={styles.doneButtonText}>Done</Text>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  cardDivider: {
    backgroundColor: CARD_BORDER,
    height: 1,
    marginVertical: spacing.xs,
  },
  content: {
    gap: spacing.lg,
    paddingBottom: spacing["3xl"],
  },
  copyButton: {
    alignItems: "center",
    justifyContent: "center",
    padding: spacing.xs,
  },
  copyButtonPressed: {
    opacity: 0.7,
  },
  detailCopy: {
    flex: 1,
    gap: 0,
    justifyContent: "center",
    minWidth: 0,
  },
  detailIconWrap: {
    alignItems: "center",
    borderRadius: radius.full,
    height: 32,
    justifyContent: "center",
    width: 32,
  },
  detailLabel: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "600",
    lineHeight: 14,
    marginBottom: 1,
  },
  detailRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: 52,
    paddingVertical: spacing.xs,
  },
  detailValue: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "700",
    lineHeight: 20,
    ...androidCompactText,
  },
  detailsCard: {
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.xl,
    borderWidth: 1,
    gap: spacing.sm,
    padding: spacing.sm,
    ...shadows.sm,
  },
  doneButton: {
    alignItems: "center",
    paddingVertical: spacing.md,
  },
  doneButtonPressed: {
    opacity: 0.8,
  },
  doneButtonText: {
    color: colors.brand.orange,
    fontSize: 16,
    fontWeight: "700",
  },
  disabledAction: {
    opacity: 0.45,
  },
  expiryItem: {
    alignItems: "center",
    flexDirection: "row",
    flex: 1,
    gap: 6,
    minWidth: 0,
  },
  expiryRow: {
    flexDirection: "row",
    gap: spacing.sm,
    justifyContent: "space-between",
  },
  expiryText: {
    color: colors.guard.textMuted,
    flexShrink: 1,
    fontSize: 11,
    fontWeight: "500",
  },
  heroArt: {
    alignItems: "center",
    height: 72,
    justifyContent: "center",
    position: "relative",
    width: 72,
  },
  heroCopy: {
    flex: 1,
    flexDirection: "row",
    gap: spacing.md,
    minWidth: 0,
  },
  heroRow: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: spacing.sm,
  },
  heroText: {
    flex: 1,
    gap: 4,
    minWidth: 0,
    paddingTop: 4,
  },
  linkBox: {
    alignItems: "center",
    backgroundColor: colors.surface.secondary,
    borderColor: CARD_BORDER,
    borderRadius: radius.md,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.xs,
    minHeight: 36,
    paddingHorizontal: spacing.xs,
    paddingVertical: spacing.xs,
  },
  linkText: {
    color: colors.brand.navy,
    flex: 1,
    fontSize: 12,
    fontWeight: "500",
    lineHeight: 16,
    ...androidCompactText,
  },
  pageSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 13,
    fontWeight: "500",
    lineHeight: 18,
  },
  pageTitle: {
    color: colors.brand.navy,
    fontSize: 20,
    fontWeight: "800",
  },
  rowDivider: {
    backgroundColor: CARD_BORDER,
    height: 1,
  },
  shareHeading: {
    color: colors.brand.navy,
    fontSize: 14,
    fontWeight: "700",
  },
  shareList: {
    gap: spacing.sm,
  },
  shareOption: {
    alignItems: "center",
    borderRadius: radius.lg,
    flexDirection: "row",
    justifyContent: "space-between",
    minHeight: 48,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  shareOptionLeft: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
  },
  shareOptionPressed: {
    opacity: 0.92,
  },
  shareOptionPrimary: {
    backgroundColor: colors.brand.orange,
    ...shadows.sm,
  },
  shareOptionSecondary: {
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderWidth: 1,
  },
  shareOptionText: {
    color: colors.brand.navy,
    fontSize: 14,
    fontWeight: "600",
  },
  shareOptionTextPrimary: {
    color: colors.text.inverse,
  },
  shareSection: {
    gap: spacing.sm,
  },
  sparkleOne: {
    position: "absolute",
    right: 4,
    top: 8,
  },
  sparkleTwo: {
    bottom: 10,
    left: 0,
    position: "absolute",
  },
  successBadge: {
    alignItems: "center",
    backgroundColor: colors.status.success,
    borderRadius: 999,
    height: 52,
    justifyContent: "center",
    width: 52,
  },
});
