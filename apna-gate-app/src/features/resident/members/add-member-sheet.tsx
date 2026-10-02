import { type ReactNode, useRef, useState } from "react";
import {
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
  type TextInputProps,
} from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { BottomActions } from "@/components/ui/bottom-actions";
import { getApiMessage } from "@/features/auth/api-error";
import { useResidentFeedback } from "@/features/resident/hooks/use-resident-feedback";
import { MemberInviteSuccess } from "@/features/resident/members/member-invite-success";
import { useResident } from "@/features/resident/resident-context";
import {
  type ModelsFlatMemberInviteRole,
  usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesMutation,
} from "@/lib/api/resident-api-extensions";
import type { ModelsShortLinkResponse } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { androidCompactText } from "@/theme/platform-styles";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

const CARD_BORDER = "rgba(16, 29, 54, 0.08)";

const ROLE_OPTIONS: {
  description: string;
  icon: SymbolViewProps["name"];
  label: string;
  value: ModelsFlatMemberInviteRole;
}[] = [
  {
    description: "Family member",
    icon: { ios: "person.2.fill", android: "group", web: "group" },
    label: "Family",
    value: "family",
  },
  {
    description: "Tenant / Occupant",
    icon: { ios: "person.fill", android: "person", web: "person" },
    label: "Tenant",
    value: "tenant",
  },
];

const ROLE_INFO: Record<ModelsFlatMemberInviteRole, string> = {
  family:
    "Family members can access resident features and view flat related information.",
  tenant: "Tenants can access resident features for the flat they occupy.",
};

type CreatedInvite = {
  expiresAt?: string;
  fullName: string;
  role: ModelsFlatMemberInviteRole;
  shortLink?: ModelsShortLinkResponse | null;
};

type AddMemberSheetProps = {
  renderLayout: (content: ReactNode, footer?: ReactNode) => ReactNode;
  onCreated: (inviteId?: number) => void;
  onSuccessViewChange?: (visible: boolean) => void;
};

function FormField({
  icon,
  label,
  ...props
}: TextInputProps & {
  icon: SymbolViewProps["name"];
  label: string;
}) {
  const [focused, setFocused] = useState(false);

  return (
    <View style={[styles.fieldCard, focused && styles.fieldCardFocused]}>
      <View style={styles.fieldIconWrap}>
        <SymbolView name={icon} size={18} tintColor={colors.brand.orange} />
      </View>
      <View style={styles.fieldCopy}>
        <Text style={styles.fieldLabel}>{label}</Text>
        <TextInput
          cursorColor={colors.brand.orange}
          placeholderTextColor={colors.guard.textMuted}
          selectionColor={colors.brand.orangeSoft}
          underlineColorAndroid="transparent"
          onBlur={(event) => {
            setFocused(false);
            props.onBlur?.(event);
          }}
          onFocus={(event) => {
            setFocused(true);
            props.onFocus?.(event);
          }}
          style={[
            styles.fieldInput,
            Platform.OS === "web" ? styles.fieldInputWeb : null,
          ]}
          {...props}
        />
      </View>
    </View>
  );
}

function RoleOptionCard({
  description,
  icon,
  label,
  selected,
  onSelect,
}: {
  description: string;
  icon: SymbolViewProps["name"];
  label: string;
  selected: boolean;
  onSelect: () => void;
}) {
  return (
    <Pressable
      accessibilityRole="button"
      accessibilityState={{ selected }}
      style={({ pressed }) => [
        styles.roleCard,
        selected && styles.roleCardSelected,
        pressed && styles.roleCardPressed,
      ]}
      onPress={onSelect}
    >
      <View
        style={[styles.roleIconWrap, selected && styles.roleIconWrapSelected]}
      >
        <SymbolView
          name={icon}
          size={18}
          tintColor={selected ? colors.brand.orange : colors.guard.textMuted}
        />
      </View>
      <View style={styles.roleCopy}>
        <Text style={[styles.roleLabel, selected && styles.roleLabelSelected]}>
          {label}
        </Text>
        <Text style={styles.roleDescription}>{description}</Text>
      </View>
      {selected ? (
        <View style={styles.roleCheck}>
          <SymbolView
            name={{ ios: "checkmark", android: "check", web: "check" }}
            size={12}
            tintColor={colors.text.inverse}
          />
        </View>
      ) : null}
    </Pressable>
  );
}

export function AddMemberSheet({
  onCreated,
  onSuccessViewChange,
  renderLayout,
}: AddMemberSheetProps) {
  const { flatId, selectedResidence, societyId } = useResident();
  const feedback = useResidentFeedback();
  const submitting = useRef(false);
  const [fullName, setFullName] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<ModelsFlatMemberInviteRole>("family");
  const [createdInvite, setCreatedInvite] = useState<CreatedInvite | null>(
    null,
  );
  const [createInvite, createInviteState] =
    usePostV1SocietiesBySocietyIdFlatsAndFlatIdMemberInvitesMutation();

  if (createdInvite) {
    return renderLayout(
      <ScrollView
        contentContainerStyle={styles.content}
        showsVerticalScrollIndicator={false}
      >
        <MemberInviteSuccess
          invite={createdInvite}
          flatNumber={selectedResidence?.flat_number ?? undefined}
          societyName={selectedResidence?.society_name ?? undefined}
          onDone={() => {
            setCreatedInvite(null);
            setFullName("");
            setPhone("");
            setEmail("");
            onSuccessViewChange?.(false);
            onCreated();
          }}
        />
      </ScrollView>,
    );
  }

  const handleCreate = async () => {
    if (!societyId || !flatId || submitting.current) {
      return;
    }

    const trimmedName = fullName.trim();
    const trimmedPhone = phone.trim();
    const trimmedEmail = email.trim();

    if (!trimmedName) {
      feedback.showInfo("Missing name", "Enter the member's full name.");
      return;
    }

    if (!trimmedPhone && !trimmedEmail) {
      feedback.showInfo(
        "Contact required",
        "Enter a phone number or email for the invite.",
      );
      return;
    }

    submitting.current = true;
    try {
      const response = await createInvite({
        societyId,
        flatId,
        modelsCreateFlatMemberInviteRequest: {
          full_name: trimmedName,
          role,
          phone: trimmedPhone || undefined,
          email: trimmedEmail || undefined,
        },
      }).unwrap();

      const inviteId = response.data?.invite?.id;
      if (inviteId) {
        feedback.showSuccess(
          "Invite ready",
          "Share the member invite from the detail page.",
        );
        onCreated(inviteId);
        return;
      }

      const shortLink = response.data?.short_link ?? response.data?.link;
      if (!shortLink?.url) {
        feedback.showSuccess("Invite created", "Share link unavailable.");
        onCreated();
        return;
      }

      setCreatedInvite({
        fullName: trimmedName,
        role,
        shortLink,
        expiresAt: response.data?.token?.expires_at,
      });
      onSuccessViewChange?.(true);
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

  const content = (
    <ScrollView
      contentContainerStyle={styles.content}
      keyboardShouldPersistTaps="handled"
      showsVerticalScrollIndicator={false}
    >
      <Text style={styles.pageSubtitle}>
        Send a join link to a family member or tenant for flat{" "}
        {selectedResidence?.flat_number ?? "—"}.
      </Text>

      <View style={styles.formCard}>
        <FormField
          autoCapitalize="words"
          autoCorrect={false}
          icon={{ ios: "person.fill", android: "person", web: "person" }}
          label="Full name"
          placeholder="Enter full name"
          value={fullName}
          onChangeText={setFullName}
        />
        <FormField
          autoComplete="tel"
          icon={{ ios: "phone.fill", android: "phone", web: "phone" }}
          keyboardType="phone-pad"
          label="Phone number"
          placeholder="Enter phone number"
          value={phone}
          onChangeText={setPhone}
        />
        <FormField
          autoCapitalize="none"
          autoComplete="email"
          icon={{ ios: "envelope.fill", android: "email", web: "email" }}
          keyboardType="email-address"
          label="Email address"
          placeholder="Enter email address"
          value={email}
          onChangeText={setEmail}
        />

        <View style={styles.roleSection}>
          <Text style={styles.roleHeading}>Role</Text>
          <View style={styles.roleRow}>
            {ROLE_OPTIONS.map((option) => (
              <RoleOptionCard
                key={option.value}
                description={option.description}
                icon={option.icon}
                label={option.label}
                selected={role === option.value}
                onSelect={() => setRole(option.value)}
              />
            ))}
          </View>
        </View>
      </View>

      <View style={styles.infoBanner}>
        <SymbolView
          name={{
            ios: "lock.shield.fill",
            android: "verified_user",
            web: "verified_user",
          }}
          size={18}
          tintColor={colors.brand.orange}
        />
        <Text style={styles.infoText}>{ROLE_INFO[role]}</Text>
      </View>
    </ScrollView>
  );
  return renderLayout(
    content,
    <BottomActions>
      <Pressable
        accessibilityRole="button"
        disabled={createInviteState.isLoading}
        style={({ pressed }) => [
          styles.createButton,
          pressed && !createInviteState.isLoading && styles.createButtonPressed,
          createInviteState.isLoading && styles.createButtonDisabled,
        ]}
        onPress={() => void handleCreate()}
      >
        {createInviteState.isLoading ? (
          <Text style={styles.createButtonText}>Creating...</Text>
        ) : (
          <>
            <SymbolView
              name={{ ios: "link", android: "link", web: "link" }}
              size={18}
              tintColor={colors.text.inverse}
            />
            <Text style={styles.createButtonText}>Create link</Text>
          </>
        )}
      </Pressable>
    </BottomActions>,
  );
}

export function AddMemberFlatBadge({
  flatNumber,
}: {
  flatNumber?: string | null;
}) {
  if (!flatNumber) {
    return null;
  }

  return (
    <View style={styles.flatBadge}>
      <View style={styles.flatBadgeIconWrap}>
        <SymbolView
          name={{
            ios: "building.2.fill",
            android: "apartment",
            web: "apartment",
          }}
          size={16}
          tintColor={colors.brand.orange}
        />
      </View>
      <Text style={styles.flatBadgeText}>{flatNumber}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: spacing.lg,
    paddingBottom: layout.screenPaddingBottom,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.xs,
  },
  createButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    flexDirection: "row",
    gap: spacing.sm,
    justifyContent: "center",
    minHeight: 52,
    paddingHorizontal: spacing.xl,
    ...shadows.sm,
  },
  createButtonDisabled: {
    opacity: 0.7,
  },
  createButtonPressed: {
    opacity: 0.92,
  },
  createButtonText: {
    color: colors.text.inverse,
    fontSize: 16,
    fontWeight: "700",
  },
  fieldCard: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: 52,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
  },
  fieldCardFocused: {
    borderColor: colors.brand.orange,
  },
  fieldCopy: {
    flex: 1,
    gap: 0,
    justifyContent: "center",
    minWidth: 0,
  },
  fieldIconWrap: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    height: 32,
    justifyContent: "center",
    width: 32,
  },
  fieldInput: {
    color: colors.brand.navy,
    flex: 1,
    fontSize: 15,
    fontWeight: "600",
    minHeight: 20,
    padding: 0,
    paddingVertical: 0,
    ...androidCompactText,
  },
  fieldInputWeb: {
    outlineWidth: 0,
  },
  fieldLabel: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "600",
    lineHeight: 14,
    marginBottom: 1,
  },
  flatBadge: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    gap: 4,
    minWidth: 52,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    ...shadows.sm,
  },
  flatBadgeIconWrap: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.full,
    height: 28,
    justifyContent: "center",
    width: 28,
  },
  flatBadgeText: {
    color: colors.brand.navy,
    fontSize: 12,
    fontWeight: "800",
  },
  formCard: {
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.xl,
    borderWidth: 1,
    gap: spacing.sm,
    padding: spacing.sm,
    ...shadows.sm,
  },
  infoBanner: {
    alignItems: "flex-start",
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: radius.lg,
    flexDirection: "row",
    gap: spacing.sm,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  infoText: {
    color: colors.brand.navy,
    flex: 1,
    fontSize: 13,
    lineHeight: 19,
  },
  pageSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 14,
    fontWeight: "500",
    lineHeight: 20,
  },
  roleCard: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: CARD_BORDER,
    borderRadius: radius.lg,
    borderWidth: 1,
    flex: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minWidth: 0,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.sm,
  },
  roleCardPressed: {
    opacity: 0.92,
  },
  roleCardSelected: {
    backgroundColor: colors.brand.orangeSoft,
    borderColor: colors.brand.orange,
  },
  roleCheck: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    height: 20,
    justifyContent: "center",
    width: 20,
  },
  roleCopy: {
    flex: 1,
    gap: 2,
    minWidth: 0,
  },
  roleDescription: {
    color: colors.guard.textMuted,
    fontSize: 10,
    fontWeight: "500",
  },
  roleHeading: {
    color: colors.brand.navy,
    fontSize: 15,
    fontWeight: "700",
  },
  roleIconWrap: {
    alignItems: "center",
    backgroundColor: colors.surface.secondary,
    borderRadius: radius.full,
    height: 32,
    justifyContent: "center",
    width: 32,
  },
  roleIconWrapSelected: {
    backgroundColor: colors.surface.card,
  },
  roleLabel: {
    color: colors.brand.navy,
    fontSize: 13,
    fontWeight: "700",
  },
  roleLabelSelected: {
    color: colors.brand.orange,
  },
  roleRow: {
    flexDirection: "row",
    gap: spacing.sm,
  },
  roleSection: {
    gap: spacing.sm,
    marginTop: spacing.xs,
  },
});
