import { Platform, StyleSheet, Text, View } from "react-native";
import { VisitorPhotoPreview } from "@/features/visitors/photos/visitor-photo-viewer";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";
import { SymbolView, type SymbolViewProps } from "expo-symbols";

import {
  formatDateOnly,
  formatDateTime,
  formatTimeOfDay,
  getCompanionDetails,
  getFlatLocationLabel,
  getVisitorEmail,
  getVisitorName,
  getVisitorPhone,
  getVisitorStatusContextMessage,
  getVisitorStatusMeta,
  titleize,
} from "@/features/guard/guard-utils";
import type {
  ModelsVisitorEntry,
  ModelsVisitorPendingEntry,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

const CARD_BORDER = "rgba(16, 29, 54, 0.08)";
const HERO_BG = "#FFF8F2";

type GuardVisitorDetailViewProps = {
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry;
  isPhotoUploading?: boolean;
};
type SymbolName = Extract<NonNullable<SymbolViewProps["name"]>, object>;
type InfoField = {
  icon: SymbolName;
  label: string;
  value: string;
};

function SectionHeader({
  icon,
  title,
}: {
  icon: SymbolName;
  title: string;
}) {
  return (
    <View style={styles.sectionHeader}>
      <SymbolView name={icon} size={18} tintColor={colors.brand.orange} />
      <Text style={styles.sectionTitle}>{title}</Text>
    </View>
  );
}

function MetaTag({ icon, label }: { icon: SymbolName; label: string }) {
  return (
    <View style={styles.metaTag}>
      <SymbolView name={icon} size={12} tintColor={colors.guard.textMuted} />
      <Text numberOfLines={1} style={styles.metaTagText}>
        {label}
      </Text>
    </View>
  );
}

function TimeCard({
  accent,
  date,
  label,
  time,
}: {
  accent: "checkin" | "checkout";
  date?: string | null;
  label: string;
  time?: string | null;
}) {
  const isCheckIn = accent === "checkin";
  const iconColor = isCheckIn ? colors.status.success : colors.brand.orange;
  const iconBg = isCheckIn
    ? colors.status.successSoft
    : colors.brand.orangeSoft;

  return (
    <View style={styles.timeCard}>
      <View style={[styles.timeIconWrap, { backgroundColor: iconBg }]}>
        <SymbolView
          name={{
            ios: "calendar",
            android: "calendar_today",
            web: "calendar_today",
          }}
          size={16}
          tintColor={iconColor}
        />
      </View>
      <Text style={styles.timeLabel}>{label}</Text>
      <Text style={styles.timeDate}>{date || "-"}</Text>
      <View style={styles.timeRow}>
        <SymbolView
          name={{ ios: "clock", android: "schedule", web: "schedule" }}
          size={11}
          tintColor={colors.guard.textMuted}
        />
        <Text style={styles.timeValue}>{time || "-"}</Text>
      </View>
    </View>
  );
}

function InfoRow({
  icon,
  isLast,
  label,
  value,
}: {
  icon: SymbolName;
  isLast?: boolean;
  label: string;
  value: string;
}) {
  return (
    <>
      <View style={styles.infoRow}>
        <View style={styles.infoIconWrap}>
          <SymbolView name={icon} size={16} tintColor={colors.brand.orange} />
        </View>
        <View style={styles.infoCopy}>
          <Text style={styles.infoLabel}>{label}</Text>
          <Text style={styles.infoValue}>{value}</Text>
        </View>
      </View>
      {!isLast ? <View style={styles.rowDivider} /> : null}
    </>
  );
}

function getDetailStatusStyle(status?: ModelsVisitorEntry["status"]) {
  switch (status) {
    case "checked_in":
    case "checked_out":
    case "approved":
      return {
        backgroundColor: colors.status.successSoft,
        borderColor: "#bbf7d0",
        color: colors.status.success,
        showCheck: true,
      };
    case "waiting_approval":
      return {
        backgroundColor: colors.brand.orangeSoft,
        borderColor: "#FED7AA",
        color: colors.brand.orange,
        showCheck: false,
      };
    case "rejected":
    case "expired":
      return {
        backgroundColor: colors.status.errorSoft,
        borderColor: "#fecaca",
        color: colors.status.error,
        showCheck: false,
      };
    default:
      return null;
  }
}

function buildInfoFields(
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry,
): InfoField[] {
  const phone = getVisitorPhone(entry);
  const email = getVisitorEmail(entry);
  const fields: InfoField[] = [
    {
      icon: { ios: "person.fill", android: "person", web: "person" },
      label: "Visitor Name",
      value: getVisitorName(entry),
    },
  ];

  if (phone) {
    fields.push({
      icon: { ios: "phone.fill", android: "phone", web: "phone" },
      label: "Phone Number",
      value: phone,
    });
  }

  if (entry.purpose) {
    fields.push({
      icon: {
        ios: "person.text.rectangle.fill",
        android: "badge",
        web: "badge",
      },
      label: "Purpose",
      value: titleize(entry.purpose),
    });
  }

  if (email) {
    fields.push({
      icon: { ios: "envelope.fill", android: "email", web: "email" },
      label: "Email",
      value: email,
    });
  }

  if (entry.vehicle_number) {
    fields.push({
      icon: {
        ios: "car.fill",
        android: "directions_car",
        web: "directions_car",
      },
      label: "Vehicle",
      value: entry.vehicle_type
        ? `${entry.vehicle_number} (${titleize(entry.vehicle_type)})`
        : entry.vehicle_number,
    });
  }

  if (entry.companions_count && entry.companions_count > 0) {
    fields.push({
      icon: { ios: "person.2.fill", android: "group", web: "group" },
      label: "Companions",
      value: String(entry.companions_count),
    });
  }

  if (entry.delivery_partner) {
    fields.push({
      icon: {
        ios: "shippingbox.fill",
        android: "local_shipping",
        web: "local_shipping",
      },
      label: "Delivery Partner",
      value: entry.delivery_partner,
    });
  }

  if (entry.service_provider) {
    fields.push({
      icon: {
        ios: "wrench.and.screwdriver.fill",
        android: "build",
        web: "build",
      },
      label: "Service Provider",
      value: entry.service_provider,
    });
  }

  return fields;
}

function buildVisitFields(
  entry: ModelsVisitorEntry | ModelsVisitorPendingEntry,
): InfoField[] {
  const fields: InfoField[] = [];

  if (entry.expected_at) {
    fields.push({
      icon: { ios: "calendar.badge.clock", android: "event", web: "event" },
      label: "Expected Visit",
      value: formatDateTime(entry.expected_at),
    });
  }

  if (entry.approved_at) {
    fields.push({
      icon: {
        ios: "checkmark.seal.fill",
        android: "verified",
        web: "verified",
      },
      label: "Approved At",
      value: formatDateTime(entry.approved_at),
    });
  }

  if (entry.created_at) {
    fields.push({
      icon: { ios: "clock.fill", android: "schedule", web: "schedule" },
      label: "Created",
      value: formatDateTime(entry.created_at),
    });
  }

  return fields;
}

export function GuardVisitorDetailView({
  entry,
  isPhotoUploading = false,
}: GuardVisitorDetailViewProps) {
  const visitorName = getVisitorName(entry);
  const flatLocation = getFlatLocationLabel(entry);
  const purposeLabel = entry.purpose ? titleize(entry.purpose) : "Guest";
  const statusMeta = getVisitorStatusMeta(entry.status);
  const detailStatus = getDetailStatusStyle(entry.status);
  const contextMessage = getVisitorStatusContextMessage(entry);
  const infoFields = buildInfoFields(entry);
  const visitFields = buildVisitFields(entry);
  const companionDetails = getCompanionDetails(entry);

  const checkInDate = entry.checked_in_at
    ? formatDateOnly(entry.checked_in_at)
    : null;
  const checkInTime = entry.checked_in_at
    ? formatTimeOfDay(entry.checked_in_at)
    : null;
  const checkOutDate = entry.checked_out_at
    ? formatDateOnly(entry.checked_out_at)
    : entry.expected_checkout_at
      ? formatDateOnly(entry.expected_checkout_at)
      : null;
  const checkOutTime = entry.checked_out_at
    ? formatTimeOfDay(entry.checked_out_at)
    : entry.expected_checkout_at
      ? formatTimeOfDay(entry.expected_checkout_at)
      : null;
  const checkOutLabel = entry.checked_out_at ? "CHECK-OUT" : "EXPECTED OUT";

  const badgeBg = detailStatus?.backgroundColor ?? statusMeta.bg;
  const badgeBorder = detailStatus?.borderColor ?? statusMeta.border;
  const badgeColor = detailStatus?.color ?? statusMeta.color;

  return (
    <View style={styles.container}>
        <View style={styles.photoBanner}>
          {entry.visitor?.photo_url && !(Platform.OS === "android" && isPhotoUploading) ? (
            <VisitorPhotoPreview entry={entry} />
          ) : isPhotoUploading ? (
            <View style={styles.photoState}>
              <SymbolView
                name={{ ios: "icloud.and.arrow.up", android: "cloud_upload", web: "cloud_upload" }}
                size={30}
                tintColor={colors.brand.orange}
              />
              <Text style={styles.photoStateTitle}>Uploading photo…</Text>
              <Text style={styles.photoStateCopy}>You can continue using the app.</Text>
            </View>
          ) : (
            <View style={styles.photoState}>
              <SymbolView
                name={{ ios: "person.crop.rectangle", android: "no_photography", web: "no_photography" }}
                size={30}
                tintColor={colors.guard.textMuted}
              />
              <Text style={styles.photoStateTitle}>No visitor photo</Text>
            </View>
          )}
        </View>
        <View style={styles.heroCard}>
          <View style={styles.heroTop}>
            <VisitorPhoto entry={entry} />

            <View style={styles.heroCopy}>
              <Text numberOfLines={1} style={styles.visitorName}>
                {visitorName}
              </Text>
              <MetaTag
                icon={{ ios: "person.fill", android: "person", web: "person" }}
                label={purposeLabel}
              />
              <MetaTag
                icon={{
                  ios: "building.2.fill",
                  android: "apartment",
                  web: "apartment",
                }}
                label={flatLocation}
              />
            </View>

            <View
              style={[
                styles.statusBadge,
                { backgroundColor: badgeBg, borderColor: badgeBorder },
              ]}
            >
              {detailStatus?.showCheck ? (
                <SymbolView
                  name={{
                    ios: "checkmark.circle.fill",
                    android: "check_circle",
                    web: "check_circle",
                  }}
                  size={12}
                  tintColor={badgeColor}
                />
              ) : (
                <SymbolView
                  name={{ ios: "clock", android: "schedule", web: "schedule" }}
                  size={11}
                  tintColor={badgeColor}
                />
              )}
              <Text style={[styles.statusText, { color: badgeColor }]}>
                {statusMeta.label}
              </Text>
            </View>
          </View>

          <View style={styles.timeCardsRow}>
            <TimeCard
              accent="checkin"
              date={checkInDate}
              label="CHECK-IN"
              time={checkInTime}
            />
            <TimeCard
              accent="checkout"
              date={checkOutDate}
              label={checkOutLabel}
              time={checkOutTime}
            />
          </View>
        </View>

        {contextMessage ? (
          <View style={styles.contextBanner}>
            <SymbolView
              name={{ ios: "info.circle.fill", android: "info", web: "info" }}
              size={15}
              tintColor={colors.brand.orange}
            />
            <Text style={styles.contextText}>{contextMessage}</Text>
          </View>
        ) : null}

        <SectionHeader
          icon={{
            ios: "building.2.fill",
            android: "apartment",
            web: "apartment",
          }}
          title="Visiting Flat"
        />
        <View style={styles.flatCard}>
          <View style={styles.flatIconWrap}>
            <SymbolView
              name={{ ios: "house.fill", android: "home", web: "home" }}
              size={18}
              tintColor={colors.brand.orange}
            />
          </View>
          <View style={styles.flatCopy}>
            <Text style={styles.flatPrimary}>{flatLocation}</Text>
            <Text style={styles.flatSecondary}>Flat</Text>
          </View>
          <SymbolView
            name={{
              ios: "chevron.right",
              android: "chevron_right",
              web: "chevron_right",
            }}
            size={14}
            tintColor={colors.guard.textMuted}
          />
        </View>

        <SectionHeader
          icon={{ ios: "person.fill", android: "person", web: "person" }}
          title="Visitor Information"
        />
        <View style={styles.infoCard}>
          {infoFields.map((field, index) => (
            <InfoRow
              key={field.label}
              icon={field.icon}
              isLast={index === infoFields.length - 1}
              label={field.label}
              value={field.value}
            />
          ))}
        </View>

        {companionDetails.length > 0 ? (
          <>
            <SectionHeader
              icon={{ ios: "person.2.fill", android: "group", web: "group" }}
              title="Companion Details"
            />
            <View style={styles.infoCard}>
              {companionDetails.map((companion, index) => (
                <CompanionInfoRow
                  key={`${companion.name}-${companion.phoneNumber}-${index}`}
                  companion={companion}
                  index={index}
                  isLast={index === companionDetails.length - 1}
                />
              ))}
            </View>
          </>
        ) : null}

        {visitFields.length > 0 ? (
          <>
            <SectionHeader
              icon={{ ios: "clock.fill", android: "schedule", web: "schedule" }}
              title="Visit Details"
            />
            <View style={styles.infoCard}>
              {visitFields.map((field, index) => (
                <InfoRow
                  key={field.label}
                  icon={field.icon}
                  isLast={index === visitFields.length - 1}
                  label={field.label}
                  value={field.value}
                />
              ))}
            </View>
          </>
        ) : null}

        {entry.notes ? (
          <>
            <SectionHeader
              icon={{
                ios: "note.text",
                android: "sticky_note_2",
                web: "sticky_note_2",
              }}
              title="Notes"
            />
            <View style={styles.notesCard}>
              <Text style={styles.notesText}>{entry.notes}</Text>
            </View>
          </>
        ) : null}

        <View style={styles.footer}>
          <SymbolView
            name={{
              ios: "shield.checkered",
              android: "verified_user",
              web: "verified_user",
            }}
            size={16}
            tintColor={colors.brand.orange}
          />
          <Text style={styles.footerText}>
            Thank you for keeping our community safe.
          </Text>
        </View>
    </View>
  );
}

function CompanionInfoRow({
  companion,
  index,
  isLast,
}: {
  companion: { name: string; phoneNumber: string };
  index: number;
  isLast?: boolean;
}) {
  return (
    <InfoRow
      icon={{ ios: "person.2.fill", android: "group", web: "group" }}
      isLast={isLast}
      label={`Companion ${index + 1}`}
      value={[companion.name, companion.phoneNumber].filter(Boolean).join(" | ")}
    />
  );
}

const styles = StyleSheet.create({
  avatar: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    height: 56,
    justifyContent: "center",
    width: 56,
  },
  avatarText: {
    color: colors.text.inverse,
    fontSize: 24,
    fontWeight: "800",
  },
  container: {
    gap: spacing.lg,
    paddingBottom: spacing["3xl"],
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.sm,
  },
  contextBanner: {
    alignItems: "flex-start",
    backgroundColor: colors.brand.orangeSoft,
    borderColor: "rgba(255, 106, 26, 0.15)",
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  contextText: {
    color: colors.guard.textMuted,
    flex: 1,
    fontSize: 13,
    lineHeight: 18,
  },
  flatCard: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    ...shadows.sm,
  },
  flatCopy: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  flatIconWrap: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  flatPrimary: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "700",
  },
  flatSecondary: {
    color: colors.guard.textMuted,
    fontSize: 12,
    fontWeight: "500",
  },
  footer: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.lg,
    flexDirection: "row",
    gap: spacing.sm,
    justifyContent: "center",
    marginTop: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  footerText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "600",
  },
  heroCard: {
    backgroundColor: HERO_BG,
    borderColor: "rgba(255, 106, 26, 0.12)",
    borderRadius: radius.xl,
    borderWidth: 1,
    gap: spacing.md,
    padding: spacing.md,
    ...shadows.sm,
  },
  photoBanner: {
    height: 240,
    backgroundColor: HERO_BG,
    borderColor: "rgba(255, 106, 26, 0.12)",
    borderRadius: 16,
    borderWidth: 1,
    overflow: "hidden",
    width: "100%",
    ...shadows.sm,
  },
  photoState: {
    alignItems: "center",
    flex: 1,
    gap: spacing.xs,
    justifyContent: "center",
    padding: spacing.lg,
  },
  photoStateCopy: {
    color: colors.guard.textMuted,
    fontSize: 12,
    textAlign: "center",
  },
  photoStateTitle: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "700",
    textAlign: "center",
  },
  heroCopy: {
    flex: 1,
    gap: 6,
    minWidth: 0,
  },
  heroTop: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: spacing.md,
  },
  infoCard: {
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    overflow: "hidden",
    ...shadows.sm,
  },
  infoCopy: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  infoIconWrap: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    height: 36,
    justifyContent: "center",
    width: 36,
  },
  infoLabel: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "600",
  },
  infoRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  infoValue: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "700",
  },
  metaTag: {
    alignItems: "center",
    flexDirection: "row",
    gap: 5,
  },
  metaTagText: {
    color: colors.guard.textMuted,
    flexShrink: 1,
    fontSize: 12,
    fontWeight: "500",
  },
  notesCard: {
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    ...shadows.sm,
  },
  notesText: {
    color: colors.brand.navy,
    fontSize: 14,
    lineHeight: 20,
  },
  rowDivider: {
    backgroundColor: CARD_BORDER,
    height: 1,
    marginLeft: spacing.md + 36 + spacing.md,
  },
  sectionHeader: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
  },
  sectionTitle: {
    color: colors.brand.navy,
    fontSize: 16,
    fontWeight: "700",
  },
  statusBadge: {
    alignItems: "center",
    borderRadius: radius["2xl"],
    borderWidth: 1,
    flexDirection: "row",
    gap: 4,
    paddingHorizontal: 8,
    paddingVertical: 4,
  },
  statusText: {
    fontSize: 10,
    fontWeight: "700",
  },
  timeCard: {
    alignItems: "flex-start",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    flex: 1,
    gap: 4,
    minWidth: 0,
    padding: spacing.sm,
  },
  timeCardsRow: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  timeDate: {
    color: colors.brand.navy,
    fontSize: 13,
    fontWeight: "700",
  },
  timeIconWrap: {
    alignItems: "center",
    borderRadius: radius.sm,
    height: 28,
    justifyContent: "center",
    width: 28,
  },
  timeLabel: {
    color: colors.guard.textMuted,
    fontSize: 9,
    fontWeight: "700",
    letterSpacing: 0.4,
  },
  timeRow: {
    alignItems: "center",
    flexDirection: "row",
    gap: 4,
  },
  timeValue: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "600",
  },
  visitorName: {
    color: colors.brand.navy,
    fontSize: 20,
    fontWeight: "800",
  },
});
