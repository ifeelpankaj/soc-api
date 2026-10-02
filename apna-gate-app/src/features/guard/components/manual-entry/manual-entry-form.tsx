import { type ReactNode, useState } from "react";
import { useToast } from "@/components/ui/toast";
import {
  PhotoPicker,
  releaseSelectedPhoto,
} from "@/features/visitors/photos/photo-picker";
import { startVisitorPhotoUpload } from "@/features/visitors/photos/visitor-photo-upload";
import { useAppDispatch } from "@/redux/hooks";
import {
  ActivityIndicator,
  Keyboard,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
  type TextInputProps,
  type TextStyle,
} from "react-native";
import { SymbolView } from "expo-symbols";
import { GuardBottomActions } from "@/features/guard/components/guard-bottom-actions";

import { Row, Stack } from "@/components/layout";
import {
  dashboardActionToneStyles,
  type DashboardActionTone,
} from "@/components/dashboard";
import { FlatPicker } from "@/features/guard/components/flat-picker";
import {
  DELIVERY_PARTNERS,
  DELIVERY_PARTNER_OTHER_LABEL,
  titleize,
  visitorPurposes,
} from "@/features/guard/guard-utils";
import {
  formatSelectedFlatLabel,
  type SelectedFlat,
  useGuardManualEntry,
} from "@/features/guard/hooks/use-guard-manual-entry";
import { usesCabOptionalCompanions } from "@/features/guard/guard-platform-policy";
import { useGuardVisitorInvite } from "@/features/guard/hooks/useGuardVisitorInvite";
import {
  buildVisitorInviteShareContent,
  copyInviteLink,
  shareInviteNative,
  shareInviteOnWhatsApp,
} from "@/features/resident/invites/invite-share";
import {
  type ModelsShortLinkResponse,
  type ModelsVisitorEntry,
  type ModelsVisitorPurpose,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { androidCompactText } from "@/theme/platform-styles";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

type EntryMode = "full_entry" | "form_link";

/** Keeps Android from tabbing focus through touchable chips while typing in inputs. */
const androidTouchableFocusProps =
  Platform.OS === "android" ? { focusable: false as const } : {};

/** Android autofill and keyboard resize can cycle focus across inputs when the keyboard opens. */
const androidTextInputProps =
  Platform.OS === "android"
    ? ({
        autoComplete: "off",
        importantForAutofill: "no",
        showSoftInputOnFocus: true,
      } as const)
    : {};

/** Removes the default black browser outline on web when an input is focused. */
const webNoOutline: TextStyle =
  Platform.OS === "web"
    ? ({ outlineStyle: "none" } as unknown as TextStyle)
    : {};

const PURPOSE_META = {
  guest: {
    ios: "person.fill",
    android: "person",
    web: "person",
    label: "Guest",
  },
  delivery: {
    ios: "shippingbox.fill",
    android: "local_shipping",
    web: "local_shipping",
    label: "Delivery",
  },
  cab: {
    ios: "car.fill",
    android: "local_taxi",
    web: "local_taxi",
    label: "Cab",
  },
  service: {
    ios: "wrench.fill",
    android: "build",
    web: "build",
    label: "Service",
  },
  maintenance: {
    ios: "hammer.fill",
    android: "handyman",
    web: "handyman",
    label: "Maint.",
  },
  staff: {
    ios: "person.badge.shield.checkmark.fill",
    android: "badge",
    web: "badge",
    label: "Staff",
  },
  other: {
    ios: "ellipsis.circle.fill",
    android: "more_horiz",
    web: "more_horiz",
    label: "Other",
  },
} as const;

const PURPOSE_TONES: Record<ModelsVisitorPurpose, DashboardActionTone> = {
  guest: "blue",
  delivery: "orange",
  cab: "purple",
  service: "neutral",
  maintenance: "orange",
  staff: "blue",
  other: "neutral",
};

const FIELD_ICONS = {
  person: { ios: "person.fill", android: "person", web: "person" },
  phone: { ios: "phone.fill", android: "phone", web: "phone" },
  email: { ios: "envelope.fill", android: "email", web: "email" },
  note: { ios: "note.text", android: "note", web: "note" },
  car: { ios: "car.fill", android: "directions_car", web: "directions_car" },
  tag: { ios: "number", android: "tag", web: "tag" },
  building: { ios: "building.2.fill", android: "business", web: "business" },
} as const;

function SectionHeader({ accent, title }: { accent?: boolean; title: string }) {
  if (accent) {
    return (
      <Row align="center" gap="sm">
        <View style={styles.sectionAccent} />
        <Text style={styles.sectionTitle}>{title}</Text>
      </Row>
    );
  }

  return <Text style={styles.sectionTitle}>{title}</Text>;
}

function Field({
  label,
  error,
  icon,
  multiline,
  style,
  onBlur,
  onFocus,
  ...props
}: TextInputProps & {
  error?: string;
  icon?: (typeof FIELD_ICONS)[keyof typeof FIELD_ICONS];
  label: string;
}) {
  const [focused, setFocused] = useState(false);

  const sharedInputProps = {
    blurOnSubmit: false,
    cursorColor: colors.brand.orange,
    multiline,
    placeholderTextColor: colors.guard.textMuted,
    selectionColor: colors.accent.selection,
    textAlignVertical: multiline ? ("top" as const) : undefined,
    underlineColorAndroid: "transparent" as const,
    onBlur: (event: Parameters<NonNullable<TextInputProps["onBlur"]>>[0]) => {
      setFocused(false);
      onBlur?.(event);
    },
    onFocus: (event: Parameters<NonNullable<TextInputProps["onFocus"]>>[0]) => {
      setFocused(true);
      onFocus?.(event);
    },
    ...androidTextInputProps,
    ...props,
  };

  if (Platform.OS === "android") {
    return (
      <View focusable={false} style={styles.fieldGroup}>
        <View
          style={[
            styles.fieldCardAndroid,
            multiline && styles.fieldCardAndroidMultiline,
            focused && styles.fieldCardAndroidFocused,
            error ? styles.fieldCardAndroidError : null,
          ]}
        >
          {icon ? (
            <View style={styles.fieldIconWrapAndroid}>
              <SymbolView
                name={{ ios: icon.ios, android: icon.android, web: icon.web }}
                size={18}
                tintColor={colors.brand.orange}
              />
            </View>
          ) : null}
          <View style={styles.fieldCopyAndroid}>
            <Text style={styles.fieldLabelAndroid}>{label}</Text>
            <TextInput
              {...sharedInputProps}
              style={[
                styles.fieldInputAndroid,
                multiline && styles.fieldInputAndroidMultiline,
                webNoOutline,
                style,
              ]}
            />
          </View>
        </View>
        {error ? <Text style={styles.fieldError}>{error}</Text> : null}
      </View>
    );
  }

  return (
    <View focusable={false} style={styles.fieldGroup}>
      <Text style={styles.fieldLabel}>{label}</Text>
      <View
        style={[
          styles.fieldInputShell,
          multiline && styles.fieldInputShellMultiline,
          error ? styles.fieldInputShellError : null,
        ]}
      >
        {icon ? (
          <View pointerEvents="none" style={styles.fieldInputIcon}>
            <SymbolView
              name={{ ios: icon.ios, android: icon.android, web: icon.web }}
              size={18}
              tintColor={colors.brand.orange}
            />
          </View>
        ) : null}
        <TextInput
          {...sharedInputProps}
          style={[
            styles.fieldInputShellInput,
            multiline && styles.fieldInputMultiline,
            webNoOutline,
            style,
          ]}
        />
      </View>
      {error ? <Text style={styles.fieldError}>{error}</Text> : null}
    </View>
  );
}

function PurposePicker({
  value,
  onChange,
}: {
  value: ModelsVisitorPurpose;
  onChange: (p: ModelsVisitorPurpose) => void;
}) {
  return (
    <Stack gap="md">
      <SectionHeader title="Select purpose" />
      <View style={styles.purposeGrid}>
        {visitorPurposes.map((p) => {
          const meta = PURPOSE_META[p];
          const tone = PURPOSE_TONES[p];
          const toneStyle = dashboardActionToneStyles[tone];
          const active = value === p;

          return (
            <Pressable
              key={p}
              accessibilityRole="button"
              accessibilityState={{ selected: active }}
              {...androidTouchableFocusProps}
              style={({ pressed }) => [
                styles.purposeTile,
                active && styles.purposeTileActive,
                {
                  opacity: pressed ? 0.85 : 1,
                  transform: [{ scale: pressed ? 0.97 : 1 }],
                },
              ]}
              onPress={() => onChange(p)}
            >
              <View
                style={[
                  styles.purposeIconWrap,
                  { backgroundColor: toneStyle.backgroundColor },
                ]}
              >
                <SymbolView
                  name={{ ios: meta.ios, android: meta.android, web: meta.web }}
                  size={24}
                  tintColor={toneStyle.iconColor}
                />
              </View>
              <Text
                numberOfLines={1}
                style={[
                  styles.purposeTileLabel,
                  active && styles.purposeTileLabelActive,
                  {
                    color: active
                      ? toneStyle.iconColor
                      : colors.guard.textMuted,
                  },
                ]}
              >
                {meta.label}
              </Text>
            </Pressable>
          );
        })}
      </View>
    </Stack>
  );
}

function EntryLinkModeButton({
  active,
  onPress,
}: {
  active: boolean;
  onPress: () => void;
}) {
  return (
    <Pressable
      accessibilityLabel="Send visitor link"
      accessibilityRole="button"
      accessibilityState={{ selected: active }}
      hitSlop={8}
      {...androidTouchableFocusProps}
      style={({ pressed }) => [
        styles.linkModeButton,
        active && styles.linkModeButtonActive,
        pressed && styles.linkModeButtonPressed,
      ]}
      onPress={onPress}
    >
      <SymbolView
        name={{ ios: "link", android: "link", web: "link" }}
        size={20}
        tintColor={active ? colors.brand.orange : colors.text.placeholder}
      />
    </Pressable>
  );
}

function InviteLinkSuccessCard({
  invite,
  onClear,
  societyName,
}: {
  invite: {
    purpose: ModelsVisitorPurpose;
    expiresAt?: string;
    flat?: SelectedFlat | null;
    shortLink?: ModelsShortLinkResponse | null;
  };
  onClear: () => void;
  societyName?: string | null;
}) {
  const shareResult = buildVisitorInviteShareContent({
    expiresAt: invite.expiresAt,
    flatLabel: formatSelectedFlatLabel(invite.flat),
    purpose: invite.purpose,
    shortLink: invite.shortLink,
    societyName,
  });
  const formUrl = shareResult.available
    ? shareResult.content.url
    : "Share link unavailable";

  return (
    <Stack gap="md" style={styles.successCard}>
      <Text style={styles.successCardTitle}>Link ready</Text>
      <Text style={styles.successCardSubtitle}>
        {titleize(invite.purpose)} · {formatSelectedFlatLabel(invite.flat)}
      </Text>
      <Text selectable style={styles.successCardUrl}>
        {formUrl}
      </Text>
      {invite.expiresAt ? (
        <Text style={styles.successCardExpiry}>
          Expires {new Date(invite.expiresAt).toLocaleString()}
        </Text>
      ) : null}
      <Stack gap="sm">
        <Pressable
          disabled={!shareResult.available}
          style={[
            styles.primaryButton,
            !shareResult.available && styles.disabledAction,
          ]}
          onPress={() =>
            shareResult.available && shareInviteOnWhatsApp(shareResult.content)
          }
        >
          <Text style={styles.primaryButtonText}>Share on WhatsApp</Text>
        </Pressable>
        <Row gap="sm">
          <Pressable
            disabled={!shareResult.available}
            style={[
              styles.secondaryButton,
              !shareResult.available && styles.disabledAction,
            ]}
            onPress={() =>
              shareResult.available && copyInviteLink(shareResult.content)
            }
          >
            <Text style={styles.secondaryButtonText}>Copy link</Text>
          </Pressable>
          <Pressable
            disabled={!shareResult.available}
            style={[
              styles.secondaryButton,
              !shareResult.available && styles.disabledAction,
            ]}
            onPress={() =>
              shareResult.available && shareInviteNative(shareResult.content)
            }
          >
            <Text style={styles.secondaryButtonText}>Share</Text>
          </Pressable>
        </Row>
        <Pressable style={styles.secondaryButton} onPress={onClear}>
          <Text style={styles.secondaryButtonText}>New link</Text>
        </Pressable>
      </Stack>
    </Stack>
  );
}

type ManualEntryFormProps = {
  renderLayout: (content: ReactNode, footer: ReactNode) => ReactNode;
  allowGuardEntry?: boolean;
  societyId: number;
  societyName?: string | null;
  onEntryCreated?: (result: {
    entry?: ModelsVisitorEntry;
    qrToken?: string;
  }) => void;
  onSuccess: (message: string) => void;
  onError: (message: string) => void;
};


function CompanionFields({ form }: { form: ReturnType<typeof useGuardManualEntry> }) {
  return (
    <>
      <Field
        key="companions-count"
        error={form.errors.companionsCount}
        icon={FIELD_ICONS.tag}
        keyboardType="number-pad"
        label="Companion count"
        placeholder="0"
        value={String(form.companionsCount)}
        onChangeText={(value) =>
          form.setCompanionsCount(
            Number(value.replace(/\D/g, "") || 0),
          )
        }
      />
      {form.companionsCount > 0 ? (
        <Stack gap="md">
          {Array.from({ length: form.companionsCount }, (_, index) => {
            const companion = form.companions[index] ?? {
              name: "",
              phoneNumber: "",
            };
            const companionErrors =
              form.errors.companionDetails?.[index];

            return (
              <View
                key={`companion-${index}`}
                style={styles.companionCard}
              >
                <Text style={styles.companionCardTitle}>
                  Companion {index + 1}
                </Text>
                <Text style={styles.companionCardHint}>
                  Name or phone is required
                </Text>
                <Field
                  autoCapitalize="words"
                  error={companionErrors?.name}
                  icon={FIELD_ICONS.person}
                  label="Companion name"
                  placeholder="Optional if phone is provided"
                  value={companion.name}
                  onChangeText={(value) =>
                    form.updateCompanion(index, "name", value)
                  }
                />
                <Field
                  error={companionErrors?.phoneNumber}
                  icon={FIELD_ICONS.phone}
                  keyboardType="phone-pad"
                  label="Companion phone"
                  placeholder="Optional if name is provided"
                  value={companion.phoneNumber}
                  onChangeText={(value) =>
                    form.updateCompanion(index, "phoneNumber", value)
                  }
                />
              </View>
            );
          })}
        </Stack>
      ) : null}
    </>
  );
}


export function ManualEntryForm({
  renderLayout,
  allowGuardEntry = true,
  societyId,
  societyName,
  onEntryCreated,
  onSuccess,
  onError,
}: ManualEntryFormProps) {
  const form = useGuardManualEntry(societyId);
  const dispatch = useAppDispatch();
  const { showToast } = useToast();
  const inviteForm = useGuardVisitorInvite(societyId);
  const [entryMode, setEntryMode] = useState<EntryMode>("full_entry");
  const [showExtra, setShowExtra] = useState(false);

  const handleSubmit = async () => {
    Keyboard.dismiss();
    if (entryMode === "form_link") {
      const result = await inviteForm.submit();
      if (result.success) onSuccess(result.message);
      else onError(result.message);
      return;
    }

    const result = await form.submit();
    if (result.success) {
      if (result.photo && result.entry?.id) {
        startVisitorPhotoUpload(dispatch, {
          societyId,
          entryId: result.entry.id,
          photo: result.photo,
          onFailure: () =>
            showToast({
              title: "Visitor entry created",
              message:
                "Photo couldn't be uploaded. You can add it from Edit Details.",
              variant: "warning",
            }),
        });
      } else if (result.photo) {
        releaseSelectedPhoto(result.photo);
      }
      onEntryCreated?.({ entry: result.entry, qrToken: result.qrToken });
      onSuccess(result.message);
    } else onError(result.message);
  };

  const isFormLinkMode = entryMode === "form_link";
  const canSubmit = isFormLinkMode
    ? inviteForm.isFormValid && !inviteForm.createInviteState.isLoading
    : allowGuardEntry &&
      form.isFormValid &&
      !form.isSaving &&
      !form.photoBusy;
  const isSubmitting = isFormLinkMode
    ? inviteForm.createInviteState.isLoading
    : form.isSaving;
  const selectedFlat = isFormLinkMode
    ? inviteForm.selectedFlat
    : form.selectedFlat;
  const flatError = isFormLinkMode ? inviteForm.flatError : form.errors.flat;
  const purpose = isFormLinkMode ? inviteForm.purpose : form.purpose;
  const setSelectedFlat = isFormLinkMode
    ? inviteForm.setSelectedFlat
    : form.setSelectedFlat;
  const setPurpose = isFormLinkMode ? inviteForm.setPurpose : form.setPurpose;

  const content = (
    <ScrollView
      pointerEvents={isSubmitting || form.photoBusy ? "none" : "auto"}
      contentContainerStyle={styles.formScrollContent}
      keyboardShouldPersistTaps="handled"
      nestedScrollEnabled={Platform.OS === "android"}
      removeClippedSubviews={Platform.OS === "android" ? false : undefined}
      showsVerticalScrollIndicator={false}
      style={styles.formScroll}
    >
      <Row align="flex-start" justify="space-between">
        <Stack gap="xs" style={styles.formTitleBlock}>
          <Row align="baseline" gap={4} justify="flex-start">
            <Text style={styles.formTitleAdd}>Add</Text>
            <Text style={styles.formTitleVisitor}>Visitor</Text>
          </Row>
          <Text style={styles.formSubtitle}>
            {isFormLinkMode
              ? "Visitor fills details online and gets a gate QR code."
              : societyName
                ? `Gate entry · ${societyName}`
                : "Quick gate entry"}
          </Text>
        </Stack>
        <EntryLinkModeButton
          active={isFormLinkMode}
          onPress={() =>
            !form.isSaving &&
            !form.photoBusy &&
            setEntryMode((mode) =>
              mode === "form_link" ? "full_entry" : "form_link",
            )
          }
        />
      </Row>

      {!allowGuardEntry ? (
        <View style={styles.permissionBanner}>
          <Text style={styles.permissionBannerTitle}>Guard entry disabled</Text>
          <Text style={styles.permissionBannerBody}>
            Your admin has not allowed guard entry for this society.
          </Text>
        </View>
      ) : null}

      <View style={{ gap: 20 }}>
        {purpose !== "staff" && (isFormLinkMode || purpose !== "service") ? (
          <FlatPicker
            error={flatError}
            selected={selectedFlat}
            societyId={societyId}
            onSelect={setSelectedFlat}
          />
        ) : null}

        <PurposePicker value={purpose} onChange={setPurpose} />
      </View>

      {!isFormLinkMode ? (
        <>
          <SectionHeader accent title="Visitor details" />
          <PhotoPicker
            value={form.photo}
            onChange={form.setPhoto}
            disabled={form.isSaving}
            onBusyChange={form.setPhotoBusy}
          />
          <View style={{ gap: 24 }}>
            <Stack gap="lg">
              {form.purpose !== "delivery" ? (
                <Field
                  key="visitor-name"
                  autoCapitalize="words"
                  error={form.errors.fullName}
                  icon={FIELD_ICONS.person}
                  label={
                    form.purpose === "cab" ? "Driver name" : "Visitor name"
                  }
                  placeholder="Full name"
                  value={form.fullName}
                  onChangeText={form.setFullName}
                />
              ) : null}
              <Field
                key="phone-number"
                error={form.errors.phoneNumber}
                icon={FIELD_ICONS.phone}
                keyboardType="phone-pad"
                label="Phone number"
                placeholder={
                  form.purpose === "cab" ? "Optional" : "10-digit mobile"
                }
                value={form.phoneNumber}
                onChangeText={form.setPhoneNumber}
              />
              {form.purpose === "guest" ? <CompanionFields form={form} /> : null}
              {form.purpose === "delivery" ? (
                <>
                  <Text style={styles.fieldLabel}>Delivery partner</Text>
                  <View style={styles.partnerChipGrid}>
                    {DELIVERY_PARTNERS.map((partner) => {
                      const active =
                        !form.deliveryPartnerIsOther &&
                        form.deliveryPartner === partner;
                      return (
                        <Pressable
                          key={partner}
                          style={[
                            styles.partnerChip,
                            active && styles.partnerChipActive,
                          ]}
                          onPress={() => form.selectDeliveryPartner(partner)}
                          {...androidTouchableFocusProps}
                        >
                          <Text
                            style={[
                              styles.partnerChipText,
                              active && styles.partnerChipTextActive,
                            ]}
                          >
                            {partner}
                          </Text>
                        </Pressable>
                      );
                    })}
                    <Pressable
                      style={[
                        styles.partnerChip,
                        form.deliveryPartnerIsOther && styles.partnerChipActive,
                      ]}
                      onPress={form.selectCustomDeliveryPartner}
                      {...androidTouchableFocusProps}
                    >
                      <Text
                        style={[
                          styles.partnerChipText,
                          form.deliveryPartnerIsOther &&
                            styles.partnerChipTextActive,
                        ]}
                      >
                        {DELIVERY_PARTNER_OTHER_LABEL}
                      </Text>
                    </Pressable>
                  </View>
                  {form.deliveryPartnerIsOther ? (
                    <Field
                      autoCapitalize="words"
                      error={form.errors.deliveryPartner}
                      icon={FIELD_ICONS.building}
                      label="Delivery from"
                      placeholder="Company or service name"
                      value={form.deliveryPartner}
                      onChangeText={form.setDeliveryPartner}
                    />
                  ) : null}
                  {!form.deliveryPartnerIsOther &&
                  form.errors.deliveryPartner ? (
                    <Text style={styles.flatPickerError}>
                      {form.errors.deliveryPartner}
                    </Text>
                  ) : null}
                </>
              ) : null}
              {form.purpose === "cab" ? (
                <>
                  <Field
                    autoCapitalize="characters"
                    error={form.errors.vehicleNumber}
                    icon={FIELD_ICONS.car}
                    label="Vehicle number"
                    placeholder="Required"
                    value={form.vehicleNumber}
                    onChangeText={form.setVehicleNumber}
                  />
                  <Text style={styles.fieldLabel}>Vehicle type</Text>
                  <Row align="center" gap={8} style={styles.vehicleTypeRow}>
                    {(["cab", "auto", "car", "bike"] as const).map((type) => {
                      const active = form.vehicleType === type;
                      return (
                        <Pressable
                          key={type}
                          style={[
                            styles.partnerChip,
                            active && styles.partnerChipActive,
                          ]}
                          onPress={() => form.setVehicleType(type)}
                          {...androidTouchableFocusProps}
                        >
                          <Text
                            style={[
                              styles.partnerChipText,
                              active && styles.partnerChipTextActive,
                            ]}
                          >
                            {titleize(type)}
                          </Text>
                        </Pressable>
                      );
                    })}
                  </Row>
                  {form.errors.vehicleType ? (
                    <Text style={styles.flatPickerError}>
                      {form.errors.vehicleType}
                    </Text>
                  ) : null}
                </>
              ) : null}
              {form.purpose === "service" || form.purpose === "maintenance" ? (
                <Field
                  error={form.errors.serviceProvider}
                  icon={FIELD_ICONS.building}
                  label={
                    form.purpose === "maintenance"
                      ? "Vendor / company"
                      : "Service provider"
                  }
                  placeholder={
                    form.purpose === "service"
                      ? "e.g. Newspaper delivery, garbage collection"
                      : "Required"
                  }
                  value={form.serviceProvider}
                  onChangeText={form.setServiceProvider}
                />
              ) : null}
            </Stack>

            <Pressable
              style={styles.extraToggle}
              onPress={() => setShowExtra((v) => !v)}
              {...androidTouchableFocusProps}
            >
              <View style={styles.extraToggleIcon}>
                <SymbolView
                  name={{
                    ios: showExtra ? "minus" : "plus",
                    android: showExtra ? "remove" : "add",
                    web: showExtra ? "remove" : "add",
                  }}
                  size={14}
                  tintColor={colors.brand.orange}
                />
              </View>
              <Text style={styles.extraToggleText}>
                {showExtra ? "Hide details" : "Additional details (optional)"}
              </Text>
              {form.optionalFieldsCount > 0 ? (
                <View style={styles.optionalCountBadge}>
                  <Text style={styles.optionalCountText}>
                    {form.optionalFieldsCount}
                  </Text>
                </View>
              ) : null}
              <SymbolView
                name={{
                  ios: showExtra ? "chevron.up" : "chevron.down",
                  android: showExtra ? "expand_less" : "expand_more",
                  web: showExtra ? "expand_less" : "expand_more",
                }}
                size={16}
                tintColor={colors.guard.textMuted}
              />
            </Pressable>

            {showExtra ? (
              <Stack gap="lg">
                <Field
                  autoCapitalize="none"
                  icon={FIELD_ICONS.email}
                  keyboardType="email-address"
                  label="Email"
                  placeholder="Optional"
                  value={form.email}
                  onChangeText={form.setEmail}
                />
                {usesCabOptionalCompanions(Platform.OS, form.purpose) ? (
                  <CompanionFields form={form} />
                ) : (
                  <Field
                    autoCapitalize="characters"
                    icon={FIELD_ICONS.car}
                    label="Vehicle number"
                    placeholder="Optional"
                    value={form.vehicleNumber}
                    onChangeText={form.setVehicleNumber}
                  />
                )}
                <Field
                  icon={FIELD_ICONS.note}
                  label="Notes"
                  multiline
                  placeholder="Optional"
                  value={form.notes}
                  onChangeText={form.setNotes}
                />
              </Stack>
            ) : null}

          </View>
        </>
      ) : (
        <>
          {inviteForm.createdInvite ? (
            <InviteLinkSuccessCard
              invite={inviteForm.createdInvite}
              onClear={inviteForm.clearCreatedInvite}
              societyName={societyName}
            />
          ) : null}
        </>
      )}
    </ScrollView>
  );
  const footer = (
    <GuardBottomActions>
      <Pressable
        disabled={!canSubmit}
        style={[
          styles.submitButton,
          canSubmit
            ? styles.submitButtonEnabled
            : styles.submitButtonDisabled,
        ]}
        onPress={handleSubmit}
        {...androidTouchableFocusProps}
      >
        {isSubmitting ? (
          <ActivityIndicator color={colors.text.inverse} />
        ) : (
          <Row align="center" gap="sm" justify="center">
            <SymbolView
              name={{
                ios: "checkmark.shield.fill",
                android: "verified_user",
                web: "verified_user",
              }}
              size={20}
              tintColor={
                canSubmit ? colors.text.inverse : colors.guard.textMuted
              }
            />
            <Text
              style={[
                styles.submitButtonText,
                !canSubmit && styles.submitButtonTextDisabled,
              ]}
            >
              {isFormLinkMode
                ? "Create link"
                : allowGuardEntry
                  ? "Create Entry"
                  : "Guard entry disabled"}
            </Text>
            <SymbolView
              name={{
                ios: "chevron.right",
                android: "chevron_right",
                web: "chevron_right",
              }}
              size={16}
              tintColor={
                canSubmit ? colors.text.inverse : colors.guard.textMuted
              }
            />
          </Row>
        )}
      </Pressable>
    </GuardBottomActions>
  );
  return renderLayout(content, footer);
}

const styles = StyleSheet.create({
  changeBadge: {
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 999,
    paddingHorizontal: spacing.md,
    paddingVertical: 6,
  },
  changeBadgeText: {
    color: colors.brand.orange,
    fontSize: 12,
    fontWeight: "600",
  },
  flatSearchCard: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
    ...shadows.sm,
  },
  flatSearchCardError: {
    borderColor: "#fca5a5",
  },
  flatSearchCopy: {
    flex: 1,
    minWidth: 0,
  },
  flatSearchIcon: {
    alignItems: "center",
    borderRadius: radius.lg,
    height: 48,
    justifyContent: "center",
    width: 48,
  },
  flatSearchIconOrange: {
    backgroundColor: colors.brand.orangeSoft,
  },
  flatSearchSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 13,
    fontWeight: "500",
  },
  flatSearchTitle: {
    color: colors.guard.text,
    fontSize: 16,
    fontWeight: "700",
  },
  entryCreatedMeta: {
    color: colors.text.muted,
    fontSize: 14,
  },
  extraToggle: {
    alignItems: "center",
    flexDirection: "row",
    gap: spacing.sm,
    paddingVertical: spacing.sm,
  },
  extraToggleIcon: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 999,
    height: 24,
    justifyContent: "center",
    width: 24,
  },
  extraToggleText: {
    color: colors.guard.textMuted,
    flex: 1,
    fontSize: 14,
    fontWeight: "500",
  },
  companionCard: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    gap: spacing.sm,
    padding: spacing.md,
  },
  companionCardHint: {
    color: colors.guard.textMuted,
    fontSize: 12,
    marginBottom: spacing.xs,
  },
  companionCardTitle: {
    color: colors.brand.navy,
    fontSize: 14,
    fontWeight: "700",
  },
  fieldError: {
    color: colors.status.error,
    fontSize: 14,
  },
  fieldGroup: {
    gap: spacing.xs,
  },
  fieldCardAndroid: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: layout.inputHeight,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  fieldCardAndroidError: {
    borderColor: "#fca5a5",
  },
  fieldCardAndroidFocused: {
    borderColor: colors.brand.orange,
  },
  fieldCardAndroidMultiline: {
    alignItems: "flex-start",
    minHeight: 112,
    paddingVertical: spacing.md,
  },
  fieldCopyAndroid: {
    flex: 1,
    gap: 0,
    justifyContent: "center",
    minWidth: 0,
  },
  fieldIconWrapAndroid: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.md,
    height: 32,
    justifyContent: "center",
    width: 32,
  },
  fieldInputAndroid: {
    color: colors.guard.text,
    flex: 1,
    fontSize: 15,
    fontWeight: "600",
    minHeight: 20,
    padding: 0,
    paddingVertical: 0,
    ...androidCompactText,
  },
  fieldInputAndroidMultiline: {
    minHeight: 72,
    paddingTop: spacing.xs,
    textAlignVertical: "top",
  },
  fieldLabelAndroid: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "600",
    letterSpacing: 0.4,
    textTransform: "uppercase",
  },
  fieldInputIcon: {
    alignItems: "center",
    justifyContent: "center",
    width: 24,
  },
  fieldInputMultiline: {
    minHeight: 72,
    paddingTop: Platform.OS === "android" ? spacing.xs : 0,
    textAlignVertical: "top",
  },
  fieldInputShell: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: layout.inputHeight,
    paddingHorizontal: spacing.md,
    paddingVertical: Platform.OS === "android" ? spacing.sm : spacing.md,
  },
  fieldInputShellError: {
    borderColor: "#fca5a5",
  },
  fieldInputShellInput: {
    color: colors.guard.text,
    flex: 1,
    fontSize: 16,
    lineHeight: 22,
    minHeight: Platform.OS === "android" ? 40 : 24,
    minWidth: 0,
    padding: 0,
    ...androidCompactText,
  },
  fieldInputShellMultiline: {
    alignItems: "flex-start",
    minHeight: 112,
    paddingVertical: spacing.md,
  },
  partnerChip: {
    backgroundColor: colors.dashboard.actionNeutralSoft,
    borderRadius: 999,
    paddingHorizontal: spacing.md,
    paddingVertical: 8,
  },
  partnerChipGrid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 8,
  },
  partnerChipActive: {
    backgroundColor: colors.brand.orangeSoft,
  },
  partnerChipText: {
    color: colors.text.secondary,
    fontSize: 13,
    fontWeight: "600",
  },
  partnerChipTextActive: {
    color: colors.brand.orange,
  },
  vehicleTypeRow: {
    flexWrap: "wrap",
  },
  fieldLabel: {
    color: colors.guard.textMuted,
    fontSize: 13,
    fontWeight: "500",
  },
  flatListItem: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingVertical: spacing.lg,
    ...shadows.card,
  },
  flatListItemMeta: {
    color: colors.text.muted,
    fontSize: 14,
  },
  flatListItemNumber: {
    color: colors.text.primary,
    flex: 1,
    fontSize: 17,
    fontWeight: "600",
  },
  flatPickerError: {
    color: colors.status.error,
    fontSize: 14,
    marginTop: spacing.sm,
  },
  flexButton: {
    flex: 1,
  },
  formScroll: {
    flex: 1,
  },
  formScrollContent: {
    paddingBottom: spacing.lg,
    gap: spacing["2xl"],
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.md,
  },
  formSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 14,
    marginTop: spacing.xs,
  },
  formTitleAdd: {
    color: colors.brand.navy,
    fontSize: 28,
    fontWeight: "700",
    letterSpacing: -0.5,
  },
  formTitleBlock: {
    flex: 1,
    paddingRight: spacing.md,
  },
  formTitleVisitor: {
    color: colors.brand.orange,
    fontSize: 28,
    fontWeight: "700",
    letterSpacing: -0.5,
  },
  linkModeButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    height: 44,
    justifyContent: "center",
    width: 44,
    ...shadows.sm,
  },
  linkModeButtonActive: {
    backgroundColor: colors.brand.orangeSoft,
    borderColor: colors.brand.orange,
  },
  linkModeButtonPressed: {
    opacity: 0.7,
  },
  modalBackButton: {
    alignItems: "center",
    backgroundColor: colors.guard.screenBg,
    borderRadius: 999,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  modalEmptyState: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
    paddingHorizontal: spacing["3xl"],
  },
  modalEmptyStateText: {
    color: colors.text.muted,
    fontSize: 16,
    textAlign: "center",
  },
  modalHeader: {
    backgroundColor: colors.surface.card,
    borderBottomColor: colors.guard.border,
    borderBottomWidth: 1,
    paddingHorizontal: spacing.lg,
    paddingVertical: 14,
  },
  modalHeaderText: {
    flex: 1,
  },
  modalListContent: {
    gap: spacing.sm,
    paddingBottom: spacing["2xl"],
    paddingHorizontal: spacing.lg,
  },
  modalListFooter: {
    color: colors.text.placeholder,
    fontSize: 13,
    paddingVertical: spacing.md,
    textAlign: "center",
  },
  modalListLoading: {
    marginVertical: spacing.sm,
  },
  modalLoading: {
    marginTop: 40,
  },
  modalNoResults: {
    color: colors.text.muted,
    marginTop: 40,
    textAlign: "center",
  },
  permissionBanner: {
    backgroundColor: colors.status.warningSoft,
    borderColor: "#fde68a",
    borderRadius: radius.lg,
    borderWidth: 1,
    gap: 4,
    padding: spacing.lg,
  },
  permissionBannerBody: {
    color: "#92400e",
    fontSize: 14,
    lineHeight: 20,
  },
  permissionBannerTitle: {
    color: "#78350f",
    fontSize: 15,
    fontWeight: "700",
  },
  modalScreen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
  modalSearchSection: {
    backgroundColor: colors.guard.screenBg,
    borderBottomColor: colors.guard.border,
    borderBottomWidth: 1,
    paddingBottom: spacing.lg,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.md,
  },
  modalSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 13,
  },
  modalTitle: {
    color: colors.guard.text,
    fontSize: 18,
    fontWeight: "700",
    letterSpacing: -0.5,
  },
  optionalCountBadge: {
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 999,
    paddingHorizontal: spacing.sm,
    paddingVertical: 2,
  },
  optionalCountText: {
    color: colors.brand.orange,
    fontSize: 11,
    fontWeight: "700",
  },
  purposeGrid: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.md,
    justifyContent: "flex-start",
  },
  purposeIconWrap: {
    alignItems: "center",
    borderRadius: radius.lg,
    height: 44,
    justifyContent: "center",
    width: 44,
  },
  purposeTile: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    gap: spacing.sm,
    minWidth: 72,
    paddingVertical: spacing.md,
    width: "22%",
    ...shadows.sm,
  },
  purposeTileActive: {
    borderColor: colors.brand.orange,
    borderWidth: 2,
  },
  purposeTileLabel: {
    color: colors.guard.textMuted,
    fontSize: 12,
    fontWeight: "500",
    textAlign: "center",
  },
  purposeTileLabelActive: {
    color: colors.guard.text,
    fontWeight: "700",
  },
  primaryButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: radius.lg,
    paddingVertical: spacing.md,
  },
  primaryButtonText: {
    color: colors.text.inverse,
    fontWeight: "600",
  },
  sectionAccent: {
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    height: 18,
    width: 4,
  },
  sectionTitle: {
    color: colors.brand.navy,
    fontSize: 17,
    fontWeight: "700",
    letterSpacing: -0.2,
  },
  searchFieldWrapper: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: layout.inputHeight,
    paddingHorizontal: spacing.md,
  },
  searchInput: {
    alignSelf: "stretch",
    backgroundColor: "transparent",
    borderWidth: 0,
    color: colors.guard.text,
    flex: 1,
    fontSize: 16,
    lineHeight: 22,
    minHeight: Platform.OS === "android" ? 44 : 40,
    paddingVertical: 0,
    ...androidCompactText,
  },
  secondaryButton: {
    alignItems: "center",
    backgroundColor: colors.surface.screen,
    borderRadius: radius.lg,
    paddingVertical: spacing.md,
  },
  secondaryButtonText: {
    color: colors.text.secondary,
    fontWeight: "600",
  },
  submitButton: {
    alignItems: "center",
    borderRadius: radius.xl,
    minHeight: layout.buttonHeight,
    paddingVertical: spacing.md,
    justifyContent: "center",
  },
  submitButtonDisabled: {
    backgroundColor: colors.border.default,
  },
  disabledAction: {
    opacity: 0.45,
  },
  submitButtonEnabled: {
    backgroundColor: colors.brand.orange,
    ...shadows.cta,
  },
  submitButtonText: {
    flexShrink: 1,
    textAlign: "center",
    color: colors.text.inverse,
    fontSize: 17,
    fontWeight: "700",
    letterSpacing: 0.5,
  },
  submitButtonTextDisabled: {
    color: colors.guard.textMuted,
  },
  successCard: {
    backgroundColor: colors.surface.card,
    borderColor: colors.accent.selection,
    borderRadius: 20,
    borderWidth: 1,
    padding: layout.screenPaddingHorizontal,
    ...shadows.hero,
  },
  successCardExpiry: {
    color: colors.text.muted,
    fontSize: 13,
  },
  successCardSubtitle: {
    color: colors.text.primary,
    fontSize: 16,
    fontWeight: "500",
  },
  successCardTitle: {
    color: colors.guard.teal,
    fontSize: 15,
    fontWeight: "600",
  },
  successHint: {
    color: colors.text.secondary,
    fontSize: 13,
    lineHeight: 18,
  },
  successCardUrl: {
    backgroundColor: colors.surface.screen,
    borderRadius: radius.lg,
    color: colors.text.secondary,
    fontSize: 13,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
});
