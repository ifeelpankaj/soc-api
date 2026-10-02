import { useFocusEffect, useLocalSearchParams } from "expo-router";
import { SymbolView, type SymbolViewProps } from "expo-symbols";
import { useCallback, useEffect } from "react";
import {
  AppState,
  Pressable,
  RefreshControl,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import Svg, { Path } from "react-native-svg";

import { LoadingState } from "@/components/ui/loading-state";
import { StatusPill } from "@/components/ui/status-pill";
import { ResidentSubScreen } from "@/features/resident/components/resident-sub-screen";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import { useResident } from "@/features/resident/resident-context";
import {
  buildMemberInviteShareContent,
  buildVisitorInviteShareContent,
  copyInviteLink,
  type InviteShareResult,
  shareInviteBySms,
  shareInviteNative,
  shareInviteOnWhatsApp,
} from "@/features/resident/invites/invite-share";
import { useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdQuery } from "@/lib/api/resident-api-extensions";
import {
  type ModelsFlatMemberInviteResponse,
  type ModelsVisitorInviteHistoryItem,
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteIdQuery,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";

type DetailType = "member" | "visitor";
type AppSymbolName = SymbolViewProps["name"];

const INVITE_SYMBOLS = {
  arrowClockwise: {
    ios: "arrow.clockwise",
    android: "refresh",
    web: "refresh",
  },
  calendar: {
    ios: "calendar",
    android: "calendar_today",
    web: "calendar_today",
  },
  checkmarkCircle: {
    ios: "checkmark.circle",
    android: "check_circle",
    web: "check_circle",
  },
  checkmarkShield: {
    ios: "checkmark.shield",
    android: "verified_user",
    web: "verified_user",
  },
  clock: { ios: "clock", android: "schedule", web: "schedule" },
  copy: { ios: "doc.on.doc", android: "content_copy", web: "content_copy" },
  link: { ios: "link", android: "link", web: "link" },
  listClipboard: {
    ios: "list.clipboard",
    android: "format_list_bulleted",
    web: "format_list_bulleted",
  },
  message: { ios: "message.fill", android: "sms", web: "sms" },
  person: { ios: "person", android: "person", web: "person" },
  person2: { ios: "person.2", android: "groups", web: "groups" },
  personCircle: {
    ios: "person.crop.circle",
    android: "account_circle",
    web: "account_circle",
  },
  phone: { ios: "phone", android: "phone", web: "phone" },
  phoneBubble: { ios: "phone.bubble.left.fill", android: "chat", web: "chat" },
  share: { ios: "square.and.arrow.up", android: "share", web: "share" },
  signOut: {
    ios: "rectangle.portrait.and.arrow.right",
    android: "logout",
    web: "logout",
  },
} satisfies Record<string, AppSymbolName>;

export function ResidentInviteDetailScreen({ type }: { type: DetailType }) {
  const params = useLocalSearchParams<{ inviteId?: string }>();
  const inviteId = Number(params.inviteId);
  const { flatId, selectedResidence, societyId } = useResident();
  const feedback = useResidentFeedback();
  const canQuery = Boolean(
    societyId && flatId && Number.isFinite(inviteId) && inviteId > 0,
  );

  const visitorQuery =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorInvitesInviteIdQuery(
      { societyId: societyId ?? 0, flatId: flatId ?? 0, inviteId },
      { skip: !canQuery || type !== "visitor" },
    );
  const memberQuery =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesInviteIdQuery(
      { societyId: societyId ?? 0, flatId: flatId ?? 0, inviteId },
      { skip: !canQuery || type !== "member" },
    );

  const isVisitorDetail = type === "visitor";
  const invite = isVisitorDetail
    ? visitorQuery.data?.data?.invite
    : memberQuery.data?.data?.invite;
  const isFetching = isVisitorDetail
    ? visitorQuery.isFetching
    : memberQuery.isFetching;
  const isLoading = isVisitorDetail
    ? visitorQuery.isLoading
    : memberQuery.isLoading;
  const isNotFound = isVisitorDetail
    ? isNotFoundError(visitorQuery.error)
    : isNotFoundError(memberQuery.error);
  const visitorRefetch = visitorQuery.refetch;
  const visitorIsUninitialized = visitorQuery.isUninitialized;
  const memberRefetch = memberQuery.refetch;
  const memberIsUninitialized = memberQuery.isUninitialized;

  const refetchCurrentDetail = useCallback(() => {
    if (isNotFound) {
      return;
    }
    if (type === "visitor") {
      if (!visitorIsUninitialized) {
        void visitorRefetch();
      }
      return;
    }
    if (!memberIsUninitialized) {
      void memberRefetch();
    }
  }, [
    isNotFound,
    memberIsUninitialized,
    memberRefetch,
    type,
    visitorIsUninitialized,
    visitorRefetch,
  ]);

  useFocusEffect(refetchCurrentDetail);

  useEffect(() => {
    const subscription = AppState.addEventListener("change", (state) => {
      if (state === "active") {
        refetchCurrentDetail();
      }
    });
    return () => subscription.remove();
  }, [refetchCurrentDetail]);

  const shareResult = buildShareModel(type, invite, selectedResidence);

  const share = async (method: "copy" | "native" | "sms" | "whatsapp") => {
    if (!shareResult.available) {
      feedback.showInfo(
        "Share link unavailable",
        "This invite does not have a short link yet.",
      );
      return;
    }

    const { content } = shareResult;
    try {
      if (method === "sms") {
        await shareInviteBySms(content);
      } else if (method === "whatsapp") {
        await shareInviteOnWhatsApp(content);
      } else if (method === "copy") {
        const copied = await copyInviteLink(content);
        feedback.showSuccess(
          copied ? "Link copied" : "Share link",
          copied
            ? "Invite link copied to clipboard."
            : "Use the share sheet to copy the invite link.",
        );
      } else {
        await shareInviteNative(content);
      }
    } catch (error) {
      feedback.showError(
        "Share failed",
        error,
        "Unable to share this invite right now.",
      );
    }
  };

  return (
    <ResidentSubScreen
      title={type === "visitor" ? "Visitor invite" : "Member invite"}
    >
      {isLoading && !invite ? (
        <LoadingState message="Loading invite" />
      ) : (
        <ScrollView
          contentContainerStyle={styles.content}
          refreshControl={
            <RefreshControl
              refreshing={isFetching}
              onRefresh={refetchCurrentDetail}
            />
          }
          showsVerticalScrollIndicator={false}
        >
          {isNotFound ? (
            <EmptyInvite message="Invite not found or no longer available." />
          ) : type === "visitor" ? (
            <VisitorDetail
              invite={invite as ModelsVisitorInviteHistoryItem | undefined}
            />
          ) : (
            <MemberDetail
              invite={invite as ModelsFlatMemberInviteResponse | undefined}
            />
          )}

          <SharePanel
            disabled={!shareResult.available}
            inviteType={type}
            url={
              shareResult.available
                ? shareResult.content.url
                : "Share link unavailable"
            }
            onCopy={() => void share("copy")}
            onNative={() => void share("native")}
            onSms={() => void share("sms")}
            onWhatsApp={() => void share("whatsapp")}
          />

          <RefreshStatusButton onPress={refetchCurrentDetail} />
        </ScrollView>
      )}
    </ResidentSubScreen>
  );
}

function VisitorDetail({
  invite,
}: {
  invite?: ModelsVisitorInviteHistoryItem;
}) {
  if (!invite) {
    return <EmptyInvite message="Unable to load this guest invite." />;
  }
  const name = invite.visitor?.full_name || "Visitor";
  let inviteTitle = `${name} is invited for ${invite.purpose}`;
  if (invite.purpose === "guest") {
    inviteTitle = `${name} is invited as a guest`;
  }
  return (
    <DetailCard
      icon={INVITE_SYMBOLS.personCircle}
      status={invite.status}
      title={inviteTitle || `Inviting for ${invite.purpose ?? "You"}`}
      items={[
        {
          icon: INVITE_SYMBOLS.person,
          label: "Purpose",
          value: titleize(invite.purpose ?? "guest"),
        },
        {
          icon: INVITE_SYMBOLS.calendar,
          label: "Created on",
          value: formatDate(invite.created_at),
        },
        {
          icon: INVITE_SYMBOLS.clock,
          label: "Expires on",
          value: formatDate(invite.expires_at),
        },
        {
          icon: INVITE_SYMBOLS.phone,
          label: "Visitor contact",
          value:
            invite.visitor?.phone_number ||
            invite.visitor?.email ||
            "Not submitted yet",
        },
        {
          icon: INVITE_SYMBOLS.listClipboard,
          label: "Entry status",
          value: invite.entry?.status
            ? titleize(invite.entry.status)
            : "Not submitted yet",
        },
        ...(invite.entry?.checked_in_at
          ? [
              {
                icon: INVITE_SYMBOLS.checkmarkCircle,
                label: "Checked in",
                value: formatDate(invite.entry.checked_in_at),
              },
            ]
          : []),
        ...(invite.entry?.checked_out_at
          ? [
              {
                icon: INVITE_SYMBOLS.signOut,
                label: "Checked out",
                value: formatDate(invite.entry.checked_out_at),
              },
            ]
          : []),
      ]}
    />
  );
}

function MemberDetail({ invite }: { invite?: ModelsFlatMemberInviteResponse }) {
  if (!invite) {
    return <EmptyInvite message="Unable to load this member invite." />;
  }
  return (
    <DetailCard
      icon={INVITE_SYMBOLS.person2}
      status={invite.status}
      title={invite.full_name || "Member invite"}
      items={[
        {
          icon: INVITE_SYMBOLS.person2,
          label: "Role",
          value: titleize(invite.role ?? "member"),
        },
        {
          icon: INVITE_SYMBOLS.phone,
          label: "Contact",
          value: invite.phone || invite.email || "No contact",
        },
        {
          icon: INVITE_SYMBOLS.calendar,
          label: "Created on",
          value: formatDate(invite.created_at),
        },
        {
          icon: INVITE_SYMBOLS.clock,
          label: "Expires on",
          value: formatDate(invite.expires_at),
        },
      ]}
    />
  );
}

type DetailItem = {
  icon: AppSymbolName;
  label: string;
  value?: string;
};

function DetailCard({
  icon,
  items,
  status,
  title,
}: {
  icon: AppSymbolName;
  items: DetailItem[];
  status?: string;
  title: string;
}) {
  const rows = chunkItems(items, 2);

  return (
    <View style={styles.card}>
      <View style={styles.header}>
        <View style={styles.headerIcon}>
          <SymbolView name={icon} size={30} tintColor={colors.brand.orange} />
        </View>
        <View style={styles.headerText}>
          <Text style={styles.title}>{title}</Text>
          {status ? <StatusPill status={status} /> : null}
        </View>
      </View>

      <View style={styles.detailGrid}>
        {rows.map((row, index) => (
          <View
            key={row.map((item) => item.label).join(":")}
            style={[styles.detailRow, index > 0 && styles.detailRowDivider]}
          >
            {row.map((item, itemIndex) => (
              <InfoItem
                key={item.label}
                icon={item.icon}
                label={item.label}
                showSideDivider={itemIndex === 0 && row.length > 1}
                value={item.value}
              />
            ))}
          </View>
        ))}
      </View>
    </View>
  );
}

function InfoItem({
  icon,
  label,
  showSideDivider,
  value,
}: {
  icon: AppSymbolName;
  label: string;
  showSideDivider?: boolean;
  value?: string;
}) {
  return (
    <View style={[styles.infoItem, showSideDivider && styles.infoItemDivider]}>
      <View style={styles.infoIcon}>
        <SymbolView name={icon} size={21} tintColor={colors.brand.orange} />
      </View>
      <View style={styles.infoText}>
        <Text style={styles.infoLabel}>{label}</Text>
        <Text numberOfLines={2} style={styles.infoValue}>
          {value || "-"}
        </Text>
      </View>
    </View>
  );
}

function SharePanel({
  disabled,
  inviteType,
  onCopy,
  onNative,
  onSms,
  onWhatsApp,
  url,
}: {
  disabled?: boolean;
  inviteType: DetailType;
  onCopy: () => void;
  onNative: () => void;
  onSms: () => void;
  onWhatsApp: () => void;
  url: string;
}) {
  return (
    <View style={styles.sharePanel}>
      <View style={styles.shareHeader}>
        <View style={styles.shareHeaderText}>
          <Text style={styles.sectionTitle}>Share invite</Text>
          <Text style={styles.sectionSubtitle}>
            Send this invite link to your{" "}
            {inviteType === "visitor" ? "guest" : "member"}
          </Text>
        </View>
        <View style={styles.secureBadge}>
          <SymbolView
            name={INVITE_SYMBOLS.checkmarkShield}
            size={18}
            tintColor={colors.accent.primary}
          />
          <Text style={styles.secureBadgeText}>Secure link</Text>
        </View>
      </View>

      <View style={[styles.linkRow, disabled && styles.shareButtonDisabled]}>
        <SymbolView
          name={INVITE_SYMBOLS.link}
          size={20}
          tintColor={colors.text.secondary}
        />
        <Text numberOfLines={1} style={styles.linkText}>
          {url}
        </Text>
        <View style={styles.linkDivider} />
        <Pressable
          accessibilityRole="button"
          disabled={disabled}
          style={styles.copyInlineButton}
          onPress={onCopy}
        >
          <SymbolView
            name={INVITE_SYMBOLS.copy}
            size={19}
            tintColor={colors.brand.orange}
          />
          <Text style={styles.copyInlineText}>Copy</Text>
        </Pressable>
      </View>

      <View style={styles.shareGrid}>
        <ShareAction
          disabled={disabled}
          label="WhatsApp"
          tone="whatsapp"
          onPress={onWhatsApp}
        />
        <ShareAction
          disabled={disabled}
          icon={INVITE_SYMBOLS.message}
          label="SMS"
          tone="sms"
          onPress={onSms}
        />

        <ShareAction
          disabled={disabled}
          icon={INVITE_SYMBOLS.share}
          label="Share"
          tone="share"
          onPress={onNative}
        />
      </View>

      <View style={styles.shareNote}>
        <View style={styles.shareNoteIcon}>
          <SymbolView
            name={INVITE_SYMBOLS.person2}
            size={22}
            tintColor={colors.brand.orange}
          />
        </View>
        <View style={styles.shareNoteText}>
          <Text style={styles.shareNoteTitle}>
            {inviteType === "visitor" ? "Visitors" : "Members"} can fill their
            details using this link.
          </Text>
          <Text style={styles.shareNoteSubtitle}>
            Share only with the person you are inviting.
          </Text>
        </View>
      </View>
    </View>
  );
}

function ShareAction({
  disabled,
  icon,
  label,
  tone,
  onPress,
}: {
  disabled?: boolean;
  icon?: AppSymbolName;
  label: string;
  tone: "copy" | "share" | "sms" | "whatsapp";
  onPress: () => void;
}) {
  const iconColor = getShareToneColor(tone);

  return (
    <Pressable
      accessibilityRole="button"
      disabled={disabled}
      style={({ pressed }) => [
        styles.shareButton,
        pressed && styles.shareButtonPressed,
        disabled && styles.shareButtonDisabled,
      ]}
      onPress={onPress}
    >
      {tone === "whatsapp" ? (
        <WhatsAppIcon color={iconColor} size={28} />
      ) : icon ? (
        <SymbolView name={icon} size={27} tintColor={iconColor} />
      ) : null}
      <Text numberOfLines={2} style={styles.shareButtonText}>
        {label}
      </Text>
    </Pressable>
  );
}

function WhatsAppIcon({ color, size }: { color: string; size: number }) {
  return (
    <Svg height={size} viewBox="0 0 32 32" width={size}>
      <Path
        d="M16 3.8c-6.62 0-12 5.2-12 11.6 0 2.3.7 4.43 1.9 6.23L4.65 28l6.65-1.55A12.5 12.5 0 0 0 16 27c6.62 0 12-5.2 12-11.6S22.62 3.8 16 3.8Z"
        fill="none"
        stroke={color}
        strokeLinejoin="round"
        strokeWidth={2.5}
      />
      <Path
        d="M11.05 10.55c.38-.73.78-.78 1.25-.78h.88c.28.02.58.08.78.58.28.68.95 2.35 1.02 2.52.1.23.13.5-.06.8-.18.27-.35.43-.58.68-.2.22-.42.47-.17.88.25.43 1.08 1.72 2.35 2.78 1.62 1.35 2.92 1.78 3.35 1.98.42.2.68.17.95-.1.3-.33 1.15-1.28 1.47-1.72.3-.43.6-.35 1.03-.2.42.13 2.62 1.18 3.08 1.4.45.22.75.33.85.52.1.2.1 1.15-.25 2.25-.35 1.08-2.05 2.07-2.87 2.15-.75.08-1.7.12-2.75-.17-.63-.18-1.43-.45-2.47-.87-4.35-1.82-7.18-6.05-7.4-6.33-.23-.3-1.77-2.28-1.77-4.35 0-2.05 1.12-3.07 1.35-3.52Z"
        fill={color}
        transform="scale(.72) translate(3.7 3.8)"
      />
    </Svg>
  );
}

function RefreshStatusButton({ onPress }: { onPress: () => void }) {
  return (
    <Pressable
      accessibilityRole="button"
      style={({ pressed }) => [
        styles.refreshButton,
        pressed && styles.shareButtonPressed,
      ]}
      onPress={onPress}
    >
      <SymbolView
        name={INVITE_SYMBOLS.arrowClockwise}
        size={21}
        tintColor={colors.brand.orange}
      />
      <Text style={styles.refreshButtonText}>Refresh status</Text>
    </Pressable>
  );
}

function EmptyInvite({ message }: { message: string }) {
  return (
    <View style={styles.card}>
      <Text style={styles.title}>{message}</Text>
    </View>
  );
}

function buildShareModel(
  type: DetailType,
  invite?: ModelsVisitorInviteHistoryItem | ModelsFlatMemberInviteResponse,
  residence?: ReturnType<typeof useResident>["selectedResidence"],
): InviteShareResult {
  if (!invite) {
    return { available: false, reason: "short_link_unavailable" };
  }
  const isVisitor = type === "visitor";
  const visitor = invite as ModelsVisitorInviteHistoryItem;
  const member = invite as ModelsFlatMemberInviteResponse;

  if (isVisitor) {
    return buildVisitorInviteShareContent({
      expiresAt: visitor.expires_at,
      flatLabel: residence?.flat_number,
      purpose: visitor.purpose,
      shortLink: visitor.short_link,
      societyName: residence?.society_name,
    });
  }

  return buildMemberInviteShareContent({
    expiresAt: member.expires_at,
    flatLabel: residence?.flat_number,
    fullName: member.full_name,
    role: member.role,
    shortLink: member.short_link,
    societyName: residence?.society_name,
  });
}

function formatDate(value?: string) {
  return value ? new Date(value).toLocaleString() : "-";
}

function titleize(value: string) {
  return value
    .replace(/[_-]+/g, " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

function chunkItems<T>(items: T[], size: number) {
  const rows: T[][] = [];
  for (let index = 0; index < items.length; index += size) {
    rows.push(items.slice(index, index + size));
  }
  return rows;
}

function getShareToneColor(tone: "copy" | "share" | "sms" | "whatsapp") {
  switch (tone) {
    case "whatsapp":
      return "#22c55e";
    case "sms":
      return "#2f80ed";
    case "copy":
      return "#7c3aed";
    case "share":
      return colors.brand.orange;
  }
}

function isNotFoundError(error: unknown) {
  return Boolean(
    error &&
    typeof error === "object" &&
    "status" in error &&
    error.status === 404,
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    gap: spacing.xl,
    padding: spacing.xl,
  },
  content: {
    gap: spacing.lg,
    padding: spacing.lg,
    paddingBottom: spacing["3xl"],
  },
  copyInlineButton: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: 40,
    paddingLeft: spacing.md,
  },
  copyInlineText: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "900",
  },
  detailGrid: {
    borderTopColor: colors.border.default,
    borderTopWidth: StyleSheet.hairlineWidth,
  },
  detailRow: {
    flexDirection: "row",
    minHeight: 88,
  },
  detailRowDivider: {
    borderTopColor: colors.border.default,
    borderTopWidth: StyleSheet.hairlineWidth,
  },
  header: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.lg,
  },
  headerIcon: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    height: 72,
    justifyContent: "center",
    width: 72,
  },
  headerText: {
    flex: 1,
    gap: spacing.sm,
  },
  infoIcon: {
    alignItems: "center",
    backgroundColor: colors.surface.secondary,
    borderRadius: radius.lg,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  infoItem: {
    alignItems: "center",
    flex: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingVertical: spacing.lg,
  },
  infoItemDivider: {
    borderRightColor: colors.border.default,
    borderRightWidth: StyleSheet.hairlineWidth,
    marginRight: spacing.lg,
    paddingRight: spacing.lg,
  },
  infoLabel: {
    color: colors.text.muted,
    fontSize: 13,
    fontWeight: "800",
    letterSpacing: 0,
  },
  infoText: {
    flex: 1,
    gap: spacing.xs,
    minWidth: 0,
  },
  infoValue: {
    color: colors.text.primary,
    fontSize: 16,
    fontWeight: "900",
    lineHeight: 21,
  },
  linkDivider: {
    alignSelf: "stretch",
    backgroundColor: colors.border.default,
    marginLeft: spacing.sm,
    width: StyleSheet.hairlineWidth,
  },
  linkRow: {
    alignItems: "center",
    backgroundColor: colors.surface.input,
    borderRadius: radius.lg,
    flexDirection: "row",
    gap: spacing.md,
    minHeight: 62,
    paddingHorizontal: spacing.md,
  },
  linkText: {
    color: colors.text.primary,
    flex: 1,
    fontSize: 14,
    fontWeight: "900",
    minWidth: 0,
  },
  refreshButton: {
    alignItems: "center",
    backgroundColor: colors.surface.input,
    borderColor: colors.border.input,
    borderRadius: radius.xl,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    justifyContent: "center",
    minHeight: 58,
    paddingHorizontal: spacing.xl,
  },
  refreshButtonText: {
    color: colors.brand.orange,
    fontSize: 17,
    fontWeight: "900",
  },
  sectionTitle: {
    color: colors.brand.navy,
    fontSize: 20,
    fontWeight: "900",
  },
  sectionSubtitle: {
    color: colors.text.muted,
    fontSize: 14,
    fontWeight: "700",
    lineHeight: 20,
  },
  secureBadge: {
    alignItems: "center",
    backgroundColor: colors.accent.primarySoft,
    borderRadius: radius.lg,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  secureBadgeText: {
    color: colors.accent.primary,
    fontSize: 13,
    fontWeight: "900",
  },
  shareButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.border.input,
    borderRadius: radius.md,
    borderWidth: 1,
    flex: 1,
    gap: spacing.sm,
    justifyContent: "center",
    minHeight: 78,
    minWidth: 0,
    paddingHorizontal: 4,
    paddingVertical: spacing.sm,
  },
  shareButtonPressed: {
    opacity: 0.82,
  },
  shareButtonDisabled: {
    opacity: 0.45,
  },
  shareButtonText: {
    color: colors.text.primary,
    fontSize: 11,
    fontWeight: "800",
    lineHeight: 14,
    textAlign: "center",
  },
  shareGrid: {
    flexDirection: "row",
    gap: 7,
  },
  sharePanel: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    gap: spacing.lg,
    padding: spacing.xl,
  },
  shareHeader: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: spacing.md,
    justifyContent: "space-between",
  },
  shareHeaderText: {
    flex: 1,
    gap: spacing.sm,
    minWidth: 0,
  },
  shareNote: {
    alignItems: "center",
    backgroundColor: colors.surface.input,
    borderRadius: radius.lg,
    flexDirection: "row",
    gap: spacing.md,
    padding: spacing.lg,
  },
  shareNoteIcon: {
    width: 32,
  },
  shareNoteSubtitle: {
    color: colors.text.muted,
    fontSize: 13,
    fontWeight: "700",
    lineHeight: 19,
  },
  shareNoteText: {
    flex: 1,
    gap: spacing.xs,
    minWidth: 0,
  },
  shareNoteTitle: {
    color: colors.text.primary,
    fontSize: 14,
    fontWeight: "900",
    lineHeight: 20,
  },
  title: {
    color: colors.brand.navy,
    fontSize: 24,
    fontWeight: "900",
  },
});
