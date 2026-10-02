import { useEffect, useState } from "react";
import {
  KeyboardAvoidingView,
  Modal,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from "react-native";
import { SymbolView } from "expo-symbols";
import { SafeAreaView } from "react-native-safe-area-context";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { BottomActions } from "@/components/ui/bottom-actions";
import { Stack } from "@/components/layout";
import { getApiMessage } from "@/features/auth/api-error";
import {
  formatDateOfBirthForInput,
  normalizeProfileDateInput,
  PROFILE_GENDER_OPTIONS,
  type ProfileGenderValue,
} from "@/features/profile/profile-formatters";
import type { ModelsUserResponse } from "@/lib/api/generated-api";
import {
  type UpdateProfileBody,
  usePatchV1AuthProfileMutation,
} from "@/lib/api/auth-api-extensions";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";

type ProfileEditSheetProps = {
  onClose: () => void;
  onSaved?: () => void;
  user?: ModelsUserResponse | null;
  visible: boolean;
};

function buildInitialValues(
  user?: ModelsUserResponse | null,
): UpdateProfileBody {
  return {
    first_name: user?.first_name ?? "",
    last_name: user?.last_name ?? "",
    phone_number: user?.phone_number ?? "",
    date_of_birth: formatDateOfBirthForInput(user?.date_of_birth),
    gender: user?.gender ?? "",
  };
}

export function ProfileEditSheet({
  onClose,
  onSaved,
  user,
  visible,
}: ProfileEditSheetProps) {
  const [values, setValues] = useState<UpdateProfileBody>(() =>
    buildInitialValues(user),
  );
  const [error, setError] = useState<string | null>(null);
  const [dobError, setDobError] = useState<string | null>(null);
  const [patchProfile, patchState] = usePatchV1AuthProfileMutation();

  useEffect(() => {
    if (visible) {
      /* eslint-disable react-hooks/set-state-in-effect -- Re-seed the edit draft with the current profile on open. */
      setValues(buildInitialValues(user));
      setError(null);
      setDobError(null);
      /* eslint-enable react-hooks/set-state-in-effect */
    }
  }, [user, visible]);

  const handleSave = async () => {
    const firstName = values.first_name?.trim();
    if (!firstName) {
      setError("First name is required.");
      return;
    }

    const phoneNumber = values.phone_number?.trim();
    if (phoneNumber && phoneNumber.length < 8) {
      setError("Enter a valid phone number.");
      return;
    }

    const normalizedDob = normalizeProfileDateInput(values.date_of_birth);
    if (!normalizedDob.ok) {
      setDobError(normalizedDob.error);
      setError(null);
      return;
    }

    const gender = values.gender?.trim();
    if (
      gender &&
      !PROFILE_GENDER_OPTIONS.some((option) => option.value === gender)
    ) {
      setError("Select a valid gender option.");
      return;
    }

    setError(null);
    setDobError(null);

    try {
      await patchProfile({
        first_name: firstName,
        last_name: values.last_name?.trim() || undefined,
        phone_number: phoneNumber || undefined,
        date_of_birth: normalizedDob.value || undefined,
        gender: (gender as ProfileGenderValue | undefined) || undefined,
      }).unwrap();
      onSaved?.();
      onClose();
    } catch (saveError) {
      setError(getApiMessage(saveError, "Could not update profile."));
    }
  };

  return (
    <Modal
      animationType="slide"
      presentationStyle="pageSheet"
      visible={visible}
      onRequestClose={onClose}
    >
      <SafeAreaView edges={["top", "left", "right"]} style={styles.screen}>
        <AppStatusBar />
        <KeyboardAvoidingView
          behavior={Platform.OS === "ios" ? "padding" : "height"}
          style={styles.screen}
        >
          <View style={styles.header}>
            <View style={styles.headerSide} />
            <Text style={styles.headerTitle}>Edit Profile</Text>
            <Pressable
              accessibilityLabel="Close"
              hitSlop={12}
              style={styles.headerSide}
              onPress={onClose}
            >
              <SymbolView
                name={{ ios: "xmark", android: "close", web: "close" }}
                size={18}
                tintColor={colors.text.secondary}
              />
            </Pressable>
          </View>

          <ScrollView
            contentContainerStyle={styles.content}
            keyboardShouldPersistTaps="handled"
          >
            <Stack gap="md">
              <View style={styles.fieldGroup}>
                <Text style={styles.fieldLabel}>First name</Text>
                <TextInput
                  autoCapitalize="words"
                  style={styles.fieldInput}
                  value={values.first_name ?? ""}
                  onChangeText={(first_name) =>
                    setValues((current) => ({ ...current, first_name }))
                  }
                />
              </View>
              <View style={styles.fieldGroup}>
                <Text style={styles.fieldLabel}>Last name</Text>
                <TextInput
                  autoCapitalize="words"
                  style={styles.fieldInput}
                  value={values.last_name ?? ""}
                  onChangeText={(last_name) =>
                    setValues((current) => ({ ...current, last_name }))
                  }
                />
              </View>
              <View style={styles.fieldGroup}>
                <Text style={styles.fieldLabel}>Phone number</Text>
                <TextInput
                  autoComplete="tel"
                  keyboardType="phone-pad"
                  style={styles.fieldInput}
                  value={values.phone_number ?? ""}
                  onChangeText={(phone_number) =>
                    setValues((current) => ({ ...current, phone_number }))
                  }
                />
              </View>
              <View style={styles.fieldGroup}>
                <Text style={styles.fieldLabel}>Gender</Text>
                <View style={styles.genderRow}>
                  {PROFILE_GENDER_OPTIONS.map((option) => {
                    const selected = values.gender === option.value;

                    return (
                      <Pressable
                        key={option.value}
                        accessibilityRole="button"
                        accessibilityState={{ selected }}
                        style={[
                          styles.genderChip,
                          selected && styles.genderChipSelected,
                        ]}
                        onPress={() =>
                          setValues((current) => ({
                            ...current,
                            gender: selected ? "" : option.value,
                          }))
                        }
                      >
                        <Text
                          style={[
                            styles.genderChipText,
                            selected && styles.genderChipTextSelected,
                          ]}
                        >
                          {option.label}
                        </Text>
                      </Pressable>
                    );
                  })}
                </View>
              </View>
              <View style={styles.fieldGroup}>
                <Text style={styles.fieldLabel}>Date of birth</Text>
                <Text style={styles.fieldHint}>
                  Use YYYY-MM-DD (for example, 1990-01-15)
                </Text>
                <TextInput
                  autoComplete="birthdate-full"
                  placeholder="1990-01-15"
                  style={[
                    styles.fieldInput,
                    dobError ? styles.fieldInputError : null,
                  ]}
                  value={values.date_of_birth ?? ""}
                  onChangeText={(date_of_birth) => {
                    setDobError(null);
                    setValues((current) => ({ ...current, date_of_birth }));
                  }}
                />
                {dobError ? (
                  <Text style={styles.errorText}>{dobError}</Text>
                ) : null}
              </View>
              {error ? <Text style={styles.errorText}>{error}</Text> : null}
            </Stack>
          </ScrollView>

          <BottomActions>
            <Pressable
              disabled={patchState.isLoading}
              style={[
                styles.saveButton,
                patchState.isLoading && styles.saveButtonDisabled,
              ]}
              onPress={() => void handleSave()}
            >
              <Text style={styles.saveButtonText}>
                {patchState.isLoading ? "Saving..." : "Save profile"}
              </Text>
            </Pressable>
          </BottomActions>
        </KeyboardAvoidingView>
      </SafeAreaView>
    </Modal>
  );
}

const styles = StyleSheet.create({
  content: {
    paddingBottom: spacing["3xl"],
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.lg,
  },
  errorText: {
    ...typography.bodySmall,
    color: colors.status.error,
  },
  fieldGroup: {
    gap: spacing.xs,
  },
  fieldHint: {
    ...typography.bodySmall,
    color: colors.text.secondary,
  },
  fieldInput: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius.lg,
    borderWidth: 1,
    color: colors.text.primary,
    fontSize: 16,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  fieldInputError: {
    borderColor: colors.status.error,
  },
  fieldLabel: {
    color: colors.text.secondary,
    fontSize: 13,
    fontWeight: "600",
  },
  footer: {
    borderTopColor: colors.border.default,
    borderTopWidth: StyleSheet.hairlineWidth,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.lg,
  },
  genderChip: {
    backgroundColor: colors.surface.card,
    borderColor: colors.border.default,
    borderRadius: radius["2xl"],
    borderWidth: 1,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  genderChipSelected: {
    backgroundColor: colors.brand.orange,
    borderColor: colors.brand.orange,
  },
  genderChipText: {
    ...typography.bodySmall,
    color: colors.text.secondary,
    fontWeight: "600",
  },
  genderChipTextSelected: {
    color: colors.text.inverse,
  },
  genderRow: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: spacing.sm,
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
    ...typography.button,
    color: colors.text.inverse,
  },
  screen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
});
