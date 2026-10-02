import {
  Redirect,
  useFocusEffect,
  useNavigation,
  useRouter,
} from "expo-router";
import { usePreventRemove } from "expo-router/react-navigation";
import { useCallback, useEffect, useRef, useState } from "react";
import {
  BackHandler,
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
import { SafeAreaView } from "react-native-safe-area-context";
import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Button, ConfirmDialog, LoadingState, useToast } from "@/components/ui";
import { BottomActions } from "@/components/ui/bottom-actions";
import { useAuth } from "@/features/auth/use-auth";
import {
  passwordChangeError,
  validatePasswordChange,
  type PasswordErrors,
  type PasswordField,
  type PasswordValues,
} from "@/features/auth/password-change";
import { ProfileScreenHeader } from "@/features/profile/components/profile-screen-header";
import { usePostV1AuthChangePasswordMutation } from "@/lib/api/generated-api";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";

const EMPTY: PasswordValues = { current: "", next: "", confirm: "" };
const FIELDS: { key: PasswordField; label: string; placeholder: string }[] = [
  {
    key: "current",
    label: "Current password",
    placeholder: "Enter your current password",
  },
  { key: "next", label: "New password", placeholder: "Create a new password" },
  {
    key: "confirm",
    label: "Confirm new password",
    placeholder: "Enter your new password again",
  },
];

export default function ChangePasswordScreen() {
  const { status, homeRoute, signOutLocal } = useAuth();
  const router = useRouter();
  const navigation = useNavigation();
  const goBack = useBackAction(homeRoute ?? "/");
  const { showToast } = useToast();
  const [changePassword] = usePostV1AuthChangePasswordMutation();
  const [values, setValues] = useState<PasswordValues>(EMPTY);
  const [visible, setVisible] = useState({
    current: false,
    next: false,
    confirm: false,
  });
  const [errors, setErrors] = useState<PasswordErrors>({});
  const [failure, setFailure] = useState<ReturnType<
    typeof passwordChangeError
  > | null>(null);
  const [phase, setPhase] = useState<
    "editing" | "submitting" | "signingOut" | "complete"
  >("editing");
  const [discardOpen, setDiscardOpen] = useState(false);
  const [allowLeave, setAllowLeave] = useState(false);
  const [destination, setDestination] = useState<{
    href: "/login" | "/forget";
    changed: boolean;
  } | null>(null);
  const inputs = useRef<Partial<Record<PasswordField, TextInput | null>>>({});
  const pendingLeave = useRef<(() => void) | null>(null);
  const submitting = useRef(false);
  const navigated = useRef(false);
  const [focusedField, setFocusedField] = useState<PasswordField | null>(null);
  const busy = phase === "submitting" || phase === "signingOut";
  const dirty = Object.values(values).some(Boolean);

  const requestLeave = useCallback(
    (action: () => void) => {
      if (submitting.current) return;
      if (dirty) {
        pendingLeave.current = action;
        setDiscardOpen(true);
      } else action();
    },
    [dirty],
  );
  usePreventRemove(
    status === "authenticated" && !allowLeave && (dirty || busy),
    ({ data }) => {
      requestLeave(() => navigation.dispatch(data.action));
    },
  );
  useFocusEffect(
    useCallback(() => {
      const subscription = BackHandler.addEventListener(
        "hardwareBackPress",
        () => {
          requestLeave(goBack);
          return true;
        },
      );
      return () => subscription.remove();
    }, [goBack, requestLeave]),
  );
  useEffect(() => {
    if (allowLeave && pendingLeave.current) {
      const action = pendingLeave.current;
      pendingLeave.current = null;
      action();
    }
  }, [allowLeave]);

  useEffect(() => {
    if (phase === "editing" && failure?.field)
      inputs.current[failure.field]?.focus();
  }, [phase, failure]);
  useEffect(() => {
    if (!destination || navigated.current) return;
    navigated.current = true;
    router.dismissAll();
    router.replace(destination.href);
    if (destination.changed)
      showToast({
        title: "Password changed",
        message: "Sign in with your new password.",
        variant: "success",
        duration: 6000,
      });
  }, [destination, router, showToast]);

  const leaveForLogin = async (reset = false, changed = false) => {
    submitting.current = true;
    setPhase("signingOut");
    setValues(EMPTY);
    setAllowLeave(true);
    await signOutLocal();
    setPhase("complete");
    setDestination({ href: reset ? "/forget" : "/login", changed });
  };

  const submit = async () => {
    if (submitting.current) return;
    const validation = validatePasswordChange(values);
    setErrors(validation);
    setFailure(null);
    const first = FIELDS.find(({ key }) => validation[key]);
    if (first) {
      inputs.current[first.key]?.focus();
      showToast({
        title: "Check your passwords",
        message: validation[first.key],
        variant: "error",
        duration: 5000,
      });
      return;
    }
    submitting.current = true;
    setPhase("submitting");
    Keyboard.dismiss();
    try {
      await changePassword({
        modelsChangePasswordRequest: {
          current_password: values.current,
          new_password: values.next,
          confirm_password: values.confirm,
        },
      }).unwrap();
    } catch (reason) {
      const problem = passwordChangeError(reason);
      setFailure(problem);
      if (problem.field) {
        setErrors({ [problem.field]: problem.message });
        inputs.current[problem.field]?.focus();
      }
      showToast({
        title: problem.uncertain
          ? "Check your account"
          : "Password not changed",
        message: problem.message,
        variant: "error",
        duration: 6000,
      });
      submitting.current = false;
      setPhase("editing");
      return;
    }
    await leaveForLogin(false, true);
  };

  if (status === "loading") return <LoadingState message="Loading account" />;
  if (status !== "authenticated" && phase === "editing")
    return <Redirect href="/login" />;

  return (
    <SafeAreaView edges={["top", "left", "right"]} style={styles.screen}>
      <AppStatusBar />
      <View style={styles.header}>
        <ProfileScreenHeader
          title="Change password"
          fallbackHomeRoute={homeRoute ?? "/"}
          onBack={() => requestLeave(goBack)}
          backDisabled={busy}
        />
      </View>
      <KeyboardAvoidingView
        style={styles.screen}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        <ScrollView
          contentContainerStyle={styles.content}
          keyboardShouldPersistTaps="handled"
          keyboardDismissMode="on-drag"
        >
          <View style={styles.intro}>
            <Text style={styles.title}>Keep your account secure</Text>
            <Text style={styles.subtitle}>
              Choose a password you don’t use for other accounts.
            </Text>
          </View>
          <View style={styles.notice}>
            <Text style={styles.noticeTitle}>You’ll sign in again</Text>
            <Text style={styles.subtitle}>
              Changing your password signs you out on all devices.
            </Text>
          </View>
          <View style={styles.card}>
            {FIELDS.map(({ key, label, placeholder }, index) => (
              <View key={key} style={styles.field}>
                <Text style={styles.label}>{label}</Text>
                <View
                  style={[
                    styles.inputRow,
                    focusedField === key ? styles.inputFocused : undefined,
                    errors[key] ? styles.inputError : undefined,
                  ]}
                >
                  <TextInput
                    ref={(ref) => {
                      inputs.current[key] = ref;
                    }}
                    accessibilityLabel={label}
                    accessibilityHint={errors[key]}
                    onFocus={() => setFocusedField(key)}
                    onBlur={() => setFocusedField(null)}
                    value={values[key]}
                    onChangeText={(value) => {
                      setValues((previous) => ({ ...previous, [key]: value }));
                      setErrors((previous) => ({
                        ...previous,
                        [key]: undefined,
                        ...(key === "next" ? { confirm: undefined } : {}),
                      }));
                      setFailure(null);
                    }}
                    placeholder={placeholder}
                    placeholderTextColor={colors.text.muted}
                    style={styles.input}
                    secureTextEntry={!visible[key]}
                    autoCapitalize="none"
                    autoCorrect={false}
                    autoComplete={
                      key === "current" ? "current-password" : "new-password"
                    }
                    textContentType={
                      key === "current" ? "password" : "newPassword"
                    }
                    editable={!busy}
                    returnKeyType={index === 2 ? "done" : "next"}
                    submitBehavior={index === 2 ? "blurAndSubmit" : "submit"}
                    onSubmitEditing={() => {
                      if (index === 2) void submit();
                      else inputs.current[FIELDS[index + 1].key]?.focus();
                    }}
                  />
                  <Pressable
                    disabled={busy}
                    accessibilityRole="button"
                    accessibilityLabel={`${visible[key] ? "Hide" : "Show"} ${label.toLowerCase()}`}
                    onPress={() =>
                      setVisible((previous) => ({
                        ...previous,
                        [key]: !previous[key],
                      }))
                    }
                    style={styles.toggle}
                  >
                    <Text style={styles.toggleText}>
                      {visible[key] ? "Hide" : "Show"}
                    </Text>
                  </Pressable>
                </View>
                {key === "next" && (
                  <Text style={styles.helper}>
                    At least 8 characters, up to 72 bytes. Some characters use
                    more than one byte.
                  </Text>
                )}
                {errors[key] && (
                  <Text accessibilityRole="alert" style={styles.error}>
                    {errors[key]}
                  </Text>
                )}
              </View>
            ))}
          </View>
          {failure && !failure.field && (
            <Text accessibilityRole="alert" style={styles.error}>
              {failure.message}
            </Text>
          )}
          {(failure?.uncertain || failure?.reset) && (
            <View style={styles.recovery}>
              <Button
                title="Sign in again"
                variant="secondary"
                disabled={busy}
                onPress={() => void leaveForLogin()}
              />
              <Button
                title="Reset password"
                variant="ghost"
                disabled={busy}
                onPress={() => void leaveForLogin(true)}
              />
            </View>
          )}
        </ScrollView>
        <BottomActions>
          {busy && (
            <Text accessibilityLiveRegion="polite" style={styles.helper}>
              {phase === "submitting"
                ? "Updating your password…"
                : "Signing you out…"}
            </Text>
          )}
          <Button
            title="Change password"
            loading={busy}
            disabled={busy}
            onPress={() => void submit()}
          />
        </BottomActions>
      </KeyboardAvoidingView>
      <ConfirmDialog
        visible={discardOpen}
        title="Discard password changes?"
        message="Leave this screen and discard what you entered?"
        cancelLabel="Keep editing"
        confirmLabel="Discard changes"
        onCancel={() => {
          pendingLeave.current = null;
          setDiscardOpen(false);
        }}
        onConfirm={() => {
          setValues(EMPTY);
          setDiscardOpen(false);
          setAllowLeave(true);
        }}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: colors.surface.screen },
  header: { paddingHorizontal: 20, paddingTop: 8 },
  content: {
    padding: 20,
    gap: 20,
    width: "100%",
    maxWidth: 600,
    alignSelf: "center",
  },
  intro: { gap: 8 },
  title: {
    color: colors.text.primary,
    fontSize: 24,
    fontWeight: "700",
    lineHeight: 32,
  },
  subtitle: { color: colors.text.secondary, fontSize: 14, lineHeight: 21 },
  notice: {
    padding: 16,
    borderRadius: 16,
    backgroundColor: colors.brand.orangeSoft,
    gap: 4,
  },
  noticeTitle: { color: colors.text.primary, fontSize: 15, fontWeight: "600" },
  card: {
    backgroundColor: colors.surface.card,
    borderRadius: 20,
    padding: 18,
    gap: 24,
    borderWidth: 1,
    borderColor: colors.border.default,
  },
  field: { gap: 8 },
  label: { fontSize: 14, fontWeight: "600", color: colors.text.primary },
  inputRow: {
    flexDirection: "row",
    alignItems: "center",
    borderWidth: 1,
    borderColor: colors.border.input,
    borderRadius: 12,
    minHeight: 56,
  },
  inputError: { borderColor: colors.status.error },
  inputFocused: { borderColor: colors.brand.orange },
  input: {
    flex: 1,
    minWidth: 0,
    padding: 14,
    fontSize: 16,
    color: colors.text.primary,
  },
  toggle: {
    minWidth: 56,
    minHeight: 48,
    alignItems: "center",
    justifyContent: "center",
    paddingHorizontal: 10,
  },
  toggleText: { color: colors.brand.navy, fontSize: 14, fontWeight: "600" },
  helper: { color: colors.text.muted, fontSize: 12, lineHeight: 18 },
  error: { color: colors.status.error, fontSize: 14, lineHeight: 20 },
  recovery: { gap: 8 },
});
