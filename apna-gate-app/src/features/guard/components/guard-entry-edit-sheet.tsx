import { useEffect, useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { selectIsPhotoUploading } from "@/redux/photoUploadsSlice";
import { startVisitorPhotoUpload } from "@/features/visitors/photos/visitor-photo-upload";
import { supportsGuardCompanions, usesCabOptionalCompanions } from "../guard-platform-policy";
import {
  PhotoPicker,
  photoFormData,
  releaseSelectedPhoto,
  type SelectedPhoto,
} from "@/features/visitors/photos/photo-picker";
import { VisitorPhoto } from "@/features/visitors/photos/visitor-photo";
import { AppToast, useToast } from "@/components/ui/toast";
import {
  logPhotoFailure,
  photoUploadError,
  PhotoOperationError,
} from "@/features/visitors/photos/photo-errors";
import { usePutVisitorPhotoMutation } from "@/lib/api/photo-api";
import {
  Modal,
  Keyboard,
  KeyboardAvoidingView,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from "react-native";
import { SymbolView } from "expo-symbols";
import {
  SafeAreaView,
} from "react-native-safe-area-context";

import { GuardBottomActions } from "@/features/guard/components/guard-bottom-actions";
import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Row, Stack } from "@/components/layout";
import { FlatPicker } from "@/features/guard/components/flat-picker";
import { getApiMessage } from "@/features/auth/api-error";
import { titleize } from "@/features/guard/guard-utils";
import {
  canEditVisitorFlat,
  selectedFlatFromEntry,
} from "@/features/guard/guard-entry-edit";
import {
  type CompanionDetail,
  type CompanionFieldErrors,
  emptyCompanion,
  hasCompanionDetailErrors,
  isValidPhone,
  resizeCompanions,
  serializeCompanionDetails,
  validateCompanionDetails,
} from "@/features/guard/guard-companions";
import {
  formatSelectedFlatLabel,
  type SelectedFlat,
} from "@/features/guard/hooks/use-guard-manual-entry";
import {
  getCompanionDetails,
  getVisitorEmail,
  getVisitorName,
  getVisitorPhone,
} from "@/features/visitors/visitor-utils";
import type { ModelsVisitorEntry } from "@/lib/api/generated-api";
import {
  type UpdateGuardVisitorEntryBody,
  usePatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdMutation,
} from "@/lib/api/guard-api-extensions";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type GuardEntryEditSheetProps = {
  entry?: ModelsVisitorEntry | null;
  onClose: () => void;
  onSaved?: (entry: ModelsVisitorEntry) => void;
  societyId: number;
  visible: boolean;
};

const VEHICLE_TYPES = ["cab", "auto", "car", "bike"] as const;

type EditFormValues = UpdateGuardVisitorEntryBody & {
  companions: CompanionDetail[];
  selectedFlat: SelectedFlat | null;
};

type EditFormErrors = {
  companionDetails?: CompanionFieldErrors[];
  form?: string;
  phoneNumber?: string;
};

function buildInitialValues(
  entry?: ModelsVisitorEntry | null,
  selectedFlat?: SelectedFlat | null,
): EditFormValues {
  const companionsCount = entry?.companions_count ?? 0;
  const companionDetails = getCompanionDetails(entry ?? undefined);

  return {
    full_name: entry ? getVisitorName(entry) : "",
    phone_number: getVisitorPhone(entry ?? undefined) ?? "",
    email: getVisitorEmail(entry ?? undefined) ?? "",
    vehicle_number: entry?.vehicle_number ?? "",
    vehicle_type: entry?.vehicle_type,
    notes: entry?.notes ?? "",
    companions: resizeCompanions(companionDetails, companionsCount),
    companions_count: companionsCount,
    flat_id: selectedFlat?.id,
    selectedFlat: selectedFlat ?? null,
  };
}

export function GuardEntryEditSheet({
  entry,
  onClose,
  onSaved,
  societyId,
  visible,
}: GuardEntryEditSheetProps) {
  const { showToast } = useToast();
  const dispatch = useAppDispatch();
  const backgroundPhotoUploading = useAppSelector((state) =>
    Platform.OS === "android" && selectIsPhotoUploading(state, { societyId, entryId: entry?.id ?? 0 }),
  );
  const allowCompanions = supportsGuardCompanions(Platform.OS, entry?.purpose);
  const [photo, setPhoto] = useState<SelectedPhoto>();
  const photoRef = useRef<SelectedPhoto | undefined>(undefined);
  const saveInFlight = useRef(false);
  const [photoBusy, setPhotoBusy] = useState(false);
  const [uploadPhoto, uploadState] = usePutVisitorPhotoMutation();
  const opened = useRef<number | undefined>(undefined);
  const initialValues = useRef(
    buildInitialValues(entry, selectedFlatFromEntry(entry)),
  );
  const [values, setValues] = useState(() =>
    buildInitialValues(entry, selectedFlatFromEntry(entry)),
  );
  const [errors, setErrors] = useState<EditFormErrors>({});
  const [patchEntry, patchState] =
    usePatchV1SocietiesBySocietyIdVisitorEntriesAndEntryIdMutation();
  const allowFlatEdit = canEditVisitorFlat(entry);
  const isSaving = patchState.isLoading || uploadState.isLoading || photoBusy;
  const close = () => {
    if (!isSaving) {
      releaseSelectedPhoto(photo);
      setPhoto(undefined);
      onClose();
    }
  };

  useEffect(() => {
    photoRef.current = photo;
  }, [photo]);

  useEffect(
    () => () => {
      releaseSelectedPhoto(photoRef.current);
    },
    [],
  );

  useEffect(() => {
    if (!visible) opened.current = undefined;
    if (visible && opened.current !== entry?.id) {
      opened.current = entry?.id;
      setValues(buildInitialValues(entry, selectedFlatFromEntry(entry)));
      initialValues.current = buildInitialValues(
        entry,
        selectedFlatFromEntry(entry),
      );
      setPhoto((current) => {
        releaseSelectedPhoto(current);
        return undefined;
      });
      setErrors({});
    }
  }, [entry, visible]);

  const setCompanionsCount = (count: number) => {
    const nextCount = Math.max(0, count);
    setValues((current) => ({
      ...current,
      companions: resizeCompanions(current.companions, nextCount),
      companions_count: nextCount,
    }));
    setErrors((current) => ({
      ...current,
      companionDetails: undefined,
      form: undefined,
    }));
  };

  const updateCompanion = (
    index: number,
    field: keyof CompanionDetail,
    value: string,
  ) => {
    setValues((current) => {
      const companions = resizeCompanions(
        current.companions,
        Math.max(current.companions.length, index + 1),
      );

      return {
        ...current,
        companions: companions.map((companion, companionIndex) =>
          companionIndex === index
            ? { ...companion, [field]: value }
            : companion,
        ),
      };
    });
    setErrors((current) => ({
      ...current,
      companionDetails: undefined,
      form: undefined,
    }));
  };

  const handleSave = async () => {
    if (!entry?.id || isSaving || (photo && backgroundPhotoUploading) || (Platform.OS === "android" && saveInFlight.current)) {
      return;
    }

    const fullName = values.full_name?.trim();
    if (!fullName) {
      setErrors({ form: "Visitor name is required." });
      return;
    }

    const phoneNumber = values.phone_number?.trim();
    if (
      phoneNumber !== initialValues.current.phone_number?.trim() &&
      phoneNumber &&
      !isValidPhone(phoneNumber)
    ) {
      setErrors({ phoneNumber: "Enter a valid 10-digit phone number." });
      return;
    }

    if (allowFlatEdit && !values.selectedFlat?.id) {
      setErrors({ form: "Select a visiting flat." });
      return;
    }

    const companionsCount =
      allowCompanions ? (values.companions_count ?? 0) : 0;

    if (allowCompanions && companionsCount > 0) {
      const companionDetailsErrors = validateCompanionDetails(
        values.companions,
        companionsCount,
      );

      if (hasCompanionDetailErrors(companionDetailsErrors)) {
        setErrors({
          companionDetails: companionDetailsErrors,
          form: "Complete each companion with a name or phone number.",
        });
        return;
      }
    }

    setErrors({});
    Keyboard.dismiss();

    const body: UpdateGuardVisitorEntryBody = {
      full_name: fullName,
      phone_number: phoneNumber ?? "",
      email: values.email?.trim() ?? "",
      vehicle_number: values.vehicle_number?.trim() ?? "",
      notes: values.notes?.trim() ?? "",
    };

    if (allowCompanions) {
      body.companion_details = serializeCompanionDetails(
        values.companions,
        companionsCount,
      ) ?? (usesCabOptionalCompanions(Platform.OS, entry.purpose) ? [] : undefined);
      body.companions_count = companionsCount;
    }

    if (entry.purpose === "cab" && values.vehicle_type) {
      body.vehicle_type = values.vehicle_type;
    }

    if (allowFlatEdit && values.selectedFlat?.id) {
      body.flat_id = values.selectedFlat.id;
    }

    const previous = initialValues.current;
    for (const key of Object.keys(
      body,
    ) as (keyof UpdateGuardVisitorEntryBody)[]) {
      const before =
        key === "companion_details"
          ? serializeCompanionDetails(
              previous.companions,
              previous.companions_count ?? 0,
            )
          : previous[key];
      const normalizedBefore =
        typeof before === "string" ? before.trim() : before;
      if (JSON.stringify(body[key]) === JSON.stringify(normalizedBefore))
        delete body[key];
    }

    let updated = entry;
    let entrySaved = false;
    saveInFlight.current = true;
    try {
      if (Object.keys(body).length > 0) {
        const response = await patchEntry({
          societyId,
          entryId: entry.id,
          body,
        }).unwrap();
        updated = response.data?.entry ?? entry;
        initialValues.current = values;
      }
      entrySaved = true;
      if (photo && Platform.OS === "android") {
        const transferredPhoto = photo;
        photoRef.current = undefined;
        setPhoto(undefined);
        startVisitorPhotoUpload(dispatch, {
          societyId,
          entryId: updated.id!,
          photo: transferredPhoto,
          photoReference: entry.visitor?.photo_url?.trim(),
          onFailure: () => showToast({
            title: "Details saved",
            message: "Photo couldn't be uploaded. You can add it from Edit Details.",
            variant: "warning",
          }),
        });
        onSaved?.(updated);
        onClose();
        return;
      }
      if (photo) {
        await uploadPhoto({
          societyId,
          entryId: updated.id!,
          photoReference: entry.visitor?.photo_url?.trim(),
          body: await photoFormData(photo),
        }).unwrap();
      }
      onSaved?.(updated);
      releaseSelectedPhoto(photo);
      photoRef.current = undefined;
      setPhoto(undefined);
      onClose();
    } catch (saveError) {
      releaseSelectedPhoto(photo);
      photoRef.current = undefined;
      setPhoto(undefined);
      const failure = photoUploadError(saveError);
      if (saveError instanceof PhotoOperationError) logPhotoFailure(failure);
      if (entrySaved) {
        onSaved?.(updated);
        onClose();
      }
      showToast({
        title: entrySaved ? "Details saved" : "Couldn't save changes",
        message: entrySaved
          ? "Photo couldn't be uploaded. Select it again to retry."
          : getApiMessage(saveError, "Please try again."),
        variant: entrySaved ? "warning" : "error",
      });
    } finally {
      saveInFlight.current = false;
    }
  };

  return (
    <Modal
      animationType="slide"
      presentationStyle="pageSheet"
      visible={visible}
      onRequestClose={close}
    >
      <SafeAreaView edges={["top", "left", "right"]} style={styles.screen}>
        <AppStatusBar />
        <KeyboardAvoidingView
          style={{ flex: 1 }}
          behavior={Platform.OS === "ios" ? "padding" : "height"}
        >
          <View style={styles.header}>
            <View style={styles.headerSide} />
            <View style={{ alignItems: "center" }}>
              <Text style={styles.headerTitle}>Edit visitor</Text>
              <Text style={styles.subtitle}>Update visitor information</Text>
            </View>
            <Pressable
              accessibilityLabel="Close"
              hitSlop={12}
              style={styles.headerSide}
              onPress={close}
            >
              <SymbolView
                name={{ ios: "xmark", android: "close", web: "close" }}
                size={18}
                tintColor={colors.text.secondary}
              />
            </Pressable>
          </View>

          <ScrollView
            style={{ flex: 1 }}
            contentContainerStyle={styles.content}
            keyboardShouldPersistTaps="handled"
          >
            <Stack gap="md">
              <PhotoPicker
                value={photo}
                onChange={setPhoto}
                disabled={isSaving || backgroundPhotoUploading}
                onBusyChange={setPhotoBusy}
                name={values.full_name || "Visitor"}
                purpose={titleize(entry?.purpose ?? "guest")}
                currentPhoto={
                  entry?.visitor?.photo_url ? (
                    <VisitorPhoto entry={entry} size={74} />
                  ) : undefined
                }
              />
              {backgroundPhotoUploading ? (
                <Text style={styles.subtitle}>Uploading photo… You can continue editing details.</Text>
              ) : null}
              <View
                pointerEvents={isSaving ? "none" : "auto"}
                style={{ gap: 16 }}
              >
                <View style={styles.contextCard}>
                  <Text style={styles.contextLabel}>Purpose</Text>
                  <Text style={styles.contextValue}>
                    {titleize(entry?.purpose ?? "guest")}
                  </Text>
                  {entry?.flat || entry?.flat_id ? (
                    <>
                      <Text
                        style={[styles.contextLabel, styles.contextLabelSpaced]}
                      >
                        Current flat
                      </Text>
                      <Text style={styles.contextValue}>
                        {formatSelectedFlatLabel(
                          selectedFlatFromEntry(entry),
                        ) || (entry.flat_id ? `Flat #${entry.flat_id}` : "—")}
                      </Text>
                    </>
                  ) : null}
                  {entry?.delivery_partner ? (
                    <>
                      <Text
                        style={[styles.contextLabel, styles.contextLabelSpaced]}
                      >
                        Delivery partner
                      </Text>
                      <Text style={styles.contextValue}>
                        {entry.delivery_partner}
                      </Text>
                    </>
                  ) : null}
                  {entry?.service_provider ? (
                    <>
                      <Text
                        style={[styles.contextLabel, styles.contextLabelSpaced]}
                      >
                        Service provider
                      </Text>
                      <Text style={styles.contextValue}>
                        {entry.service_provider}
                      </Text>
                    </>
                  ) : null}
                </View>

                {allowFlatEdit ? (
                  <FlatPicker
                    label="Visiting flat"
                    selected={values.selectedFlat ?? null}
                    societyId={societyId}
                    onSelect={(flat) =>
                      setValues((current) => ({
                        ...current,
                        selectedFlat: flat,
                        flat_id: flat.id,
                      }))
                    }
                  />
                ) : null}

                <View style={styles.fieldGroup}>
                  <Text style={styles.fieldLabel}>Full name</Text>
                  <TextInput
                    autoCapitalize="words"
                    style={styles.fieldInput}
                    value={values.full_name ?? ""}
                    onChangeText={(full_name) =>
                      setValues((current) => ({ ...current, full_name }))
                    }
                  />
                </View>
                <View style={styles.fieldGroup}>
                  <Text style={styles.fieldLabel}>Phone number</Text>
                  <TextInput
                    autoComplete="tel"
                    keyboardType="phone-pad"
                    style={[
                      styles.fieldInput,
                      errors.phoneNumber && styles.fieldInputError,
                    ]}
                    value={values.phone_number ?? ""}
                    onChangeText={(phone_number) => {
                      setValues((current) => ({ ...current, phone_number }));
                      setErrors((current) => ({
                        ...current,
                        form: undefined,
                        phoneNumber: undefined,
                      }));
                    }}
                  />
                  {errors.phoneNumber ? (
                    <Text style={styles.fieldError}>{errors.phoneNumber}</Text>
                  ) : null}
                </View>
                <View style={styles.fieldGroup}>
                  <Text style={styles.fieldLabel}>Email</Text>
                  <TextInput
                    autoCapitalize="none"
                    autoComplete="email"
                    keyboardType="email-address"
                    style={styles.fieldInput}
                    value={values.email ?? ""}
                    onChangeText={(email) =>
                      setValues((current) => ({ ...current, email }))
                    }
                  />
                </View>
                {allowCompanions ? (
                  <View style={styles.fieldGroup}>
                    <Text style={styles.fieldLabel}>Companion count</Text>
                    <TextInput
                      keyboardType="number-pad"
                      style={styles.fieldInput}
                      value={String(values.companions_count ?? 0)}
                      onChangeText={(text) =>
                        setCompanionsCount(Number(text.replace(/\D/g, "") || 0))
                      }
                    />
                  </View>
                ) : null}
                {allowCompanions &&
                (values.companions_count ?? 0) > 0 ? (
                  <Stack gap="md">
                    {Array.from(
                      { length: values.companions_count ?? 0 },
                      (_, index) => {
                        const companion =
                          values.companions[index] ?? emptyCompanion();
                        const companionErrors =
                          errors.companionDetails?.[index];

                        return (
                          <View
                            key={`edit-companion-${index}`}
                            style={styles.companionCard}
                          >
                            <Text style={styles.companionCardTitle}>
                              Companion {index + 1}
                            </Text>
                            <Text style={styles.companionCardHint}>
                              Name or phone is required
                            </Text>
                            <View style={styles.fieldGroup}>
                              <Text style={styles.fieldLabel}>
                                Companion name
                              </Text>
                              <TextInput
                                autoCapitalize="words"
                                placeholder="Optional if phone is provided"
                                placeholderTextColor={colors.text.placeholder}
                                style={[
                                  styles.fieldInput,
                                  companionErrors?.name &&
                                    styles.fieldInputError,
                                ]}
                                value={companion.name}
                                onChangeText={(value) =>
                                  updateCompanion(index, "name", value)
                                }
                              />
                              {companionErrors?.name ? (
                                <Text style={styles.fieldError}>
                                  {companionErrors.name}
                                </Text>
                              ) : null}
                            </View>
                            <View style={styles.fieldGroup}>
                              <Text style={styles.fieldLabel}>
                                Companion phone
                              </Text>
                              <TextInput
                                keyboardType="phone-pad"
                                placeholder="Optional if name is provided"
                                placeholderTextColor={colors.text.placeholder}
                                style={[
                                  styles.fieldInput,
                                  companionErrors?.phoneNumber &&
                                    styles.fieldInputError,
                                ]}
                                value={companion.phoneNumber}
                                onChangeText={(value) =>
                                  updateCompanion(index, "phoneNumber", value)
                                }
                              />
                              {companionErrors?.phoneNumber ? (
                                <Text style={styles.fieldError}>
                                  {companionErrors.phoneNumber}
                                </Text>
                              ) : null}
                            </View>
                          </View>
                        );
                      },
                    )}
                  </Stack>
                ) : null}
                <View style={styles.fieldGroup}>
                  <Text style={styles.fieldLabel}>Vehicle number</Text>
                  <TextInput
                    autoCapitalize="characters"
                    style={styles.fieldInput}
                    value={values.vehicle_number ?? ""}
                    onChangeText={(vehicle_number) =>
                      setValues((current) => ({ ...current, vehicle_number }))
                    }
                  />
                </View>
                {entry?.purpose === "cab" ? (
                  <View style={styles.fieldGroup}>
                    <Text style={styles.fieldLabel}>Vehicle type</Text>
                    <Row align="center" gap={8} style={styles.chipRow}>
                      {VEHICLE_TYPES.map((type) => {
                        const active = values.vehicle_type === type;
                        return (
                          <Pressable
                            key={type}
                            style={[styles.chip, active && styles.chipActive]}
                            onPress={() =>
                              setValues((current) => ({
                                ...current,
                                vehicle_type: type,
                              }))
                            }
                          >
                            <Text
                              style={[
                                styles.chipText,
                                active && styles.chipTextActive,
                              ]}
                            >
                              {titleize(type)}
                            </Text>
                          </Pressable>
                        );
                      })}
                    </Row>
                  </View>
                ) : null}
                <View style={styles.fieldGroup}>
                  <Text style={styles.fieldLabel}>Notes</Text>
                  <TextInput
                    multiline
                    numberOfLines={3}
                    style={[styles.fieldInput, styles.fieldInputMultiline]}
                    value={values.notes ?? ""}
                    onChangeText={(notes) =>
                      setValues((current) => ({ ...current, notes }))
                    }
                  />
                </View>
                {errors.form ? (
                  <Text style={styles.errorText}>{errors.form}</Text>
                ) : null}
              </View>
            </Stack>
          </ScrollView>

          <GuardBottomActions>
            <Pressable
              disabled={isSaving}
              style={[styles.saveButton, isSaving && styles.saveButtonDisabled]}
              onPress={() => void handleSave()}
            >
              <Text style={styles.saveButtonText}>
                {isSaving ? "Saving..." : "Save changes"}
              </Text>
            </Pressable>
          </GuardBottomActions>
        </KeyboardAvoidingView>
      </SafeAreaView>
      {visible ? <AppToast /> : null}
    </Modal>
  );
}

const styles = StyleSheet.create({
  chip: {
    backgroundColor: colors.dashboard.actionNeutralSoft,
    borderRadius: 999,
    paddingHorizontal: spacing.md,
    paddingVertical: 8,
  },
  chipActive: {
    backgroundColor: colors.brand.orangeSoft,
  },
  chipRow: {
    flexWrap: "wrap",
  },
  chipText: {
    color: colors.text.secondary,
    fontSize: 13,
    fontWeight: "600",
  },
  chipTextActive: {
    color: colors.brand.orange,
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
    color: colors.text.secondary,
    fontSize: 12,
  },
  companionCardTitle: {
    color: colors.text.primary,
    fontSize: 15,
    fontWeight: "700",
  },
  content: {
    paddingBottom: spacing.lg,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.lg,
  },
  subtitle: { color: "#7A879D", fontSize: 12, marginTop: 2 },
  contextCard: {
    backgroundColor: "#F7F8FA",
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    padding: spacing.md,
  },
  contextLabel: {
    color: colors.text.secondary,
    fontSize: 12,
    fontWeight: "600",
    textTransform: "uppercase",
  },
  contextLabelSpaced: {
    marginTop: spacing.sm,
  },
  contextValue: {
    color: colors.text.primary,
    fontSize: 15,
    fontWeight: "600",
    marginTop: 2,
  },
  errorText: {
    ...typography.bodySmall,
    color: colors.status.error,
  },
  fieldGroup: {
    gap: spacing.xs,
  },
  fieldInput: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: 10,
    borderWidth: 1,
    color: colors.text.primary,
    fontSize: 14,
    paddingHorizontal: spacing.md,
    paddingVertical: 11,
  },
  fieldInputError: {
    borderColor: colors.status.error,
  },
  fieldInputMultiline: {
    minHeight: 96,
    textAlignVertical: "top",
  },
  fieldLabel: {
    color: colors.text.secondary,
    fontSize: 13,
    fontWeight: "600",
  },
  fieldError: {
    ...typography.bodySmall,
    color: colors.status.error,
  },
  header: {
    alignItems: "center",
    borderBottomColor: colors.border.default,
    borderBottomWidth: StyleSheet.hairlineWidth,
    flexDirection: "row",
    justifyContent: "space-between",
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.md,
  },
  headerSide: {
    alignItems: "center",
    height: 32,
    justifyContent: "center",
    width: 32,
  },
  headerTitle: {
    ...typography.subtitle,
    color: colors.text.primary,
    fontWeight: "700",
  },
  saveButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: radius.xl,
    paddingVertical: spacing.md,
  },
  saveButtonDisabled: {
    opacity: 0.7,
  },
  saveButtonText: {
    textAlign: "center",
    ...typography.button,
    color: colors.text.inverse,
  },
  screen: {
    backgroundColor: "#FFFFFF",
    flex: 1,
  },
});
