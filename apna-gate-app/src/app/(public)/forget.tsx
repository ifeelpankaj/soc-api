import { colors } from "@/theme/colors";
import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Image } from "expo-image";
import { SymbolView } from "expo-symbols";
import { useRouter } from "expo-router";
import { useRef, useState } from "react";
import {
  ActivityIndicator,
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

import { useToast } from "@/components/ui";
import { getApiMessage } from "@/features/auth/api-error";
import {
  usePostV1AuthForgotPasswordMutation,
  usePostV1AuthResetPasswordMutation,
} from "@/lib/api/generated-api";

type Step = "email" | "reset";

const isValidEmail = (value: string) => /\S+@\S+\.\S+/.test(value);

export default function ForgetScreen() {
  const router = useRouter();
  const submitting = useRef(false);
  const { showToast } = useToast();

  const [step, setStep] = useState<Step>("email");
  const [email, setEmail] = useState("");
  const [otp, setOtp] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const [forgotPassword, { isLoading }] = usePostV1AuthForgotPasswordMutation();
  const [resetPassword, { isLoading: isResetLoading }] =
    usePostV1AuthResetPasswordMutation();

  const normalizedEmail = email.trim().toLowerCase();
  const isBusy = isLoading || isResetLoading;

  const handleSendCode = async () => {
    if (submitting.current) return;
    if (!isValidEmail(normalizedEmail)) {
      showToast({
        title: "Invalid email",
        message: "Enter a valid email address.",
        variant: "error",
      });
      return;
    }

    submitting.current = true;
    try {
      await forgotPassword({
        modelsForgotPasswordRequest: { email: normalizedEmail },
      }).unwrap();
      setStep("reset");
      showToast({
        title: "Check your email",
        message: "If an account exists, a reset code has been sent.",
        variant: "success",
      });
    } catch (error) {
      showToast({
        title: "Request failed",
        message: getApiMessage(error, "Unable to send reset code."),
        variant: "error",
      });
    } finally {
      submitting.current = false;
    }
  };

  const handleResetPassword = async () => {
    if (submitting.current) return;
    if (!otp.trim()) {
      showToast({
        title: "Missing code",
        message: "Enter the code sent to your email.",
        variant: "error",
      });
      return;
    }
    if (newPassword.length < 8) {
      showToast({
        title: "Weak password",
        message: "Password must be at least 8 characters.",
        variant: "error",
      });
      return;
    }
    if (newPassword !== confirmPassword) {
      showToast({
        title: "Password mismatch",
        message: "Passwords do not match.",
        variant: "error",
      });
      return;
    }

    submitting.current = true;
    try {
      await resetPassword({
        modelsResetPasswordRequest: {
          email: normalizedEmail,
          otp: otp.trim(),
          new_password: newPassword,
          confirm_password: confirmPassword,
        },
      }).unwrap();

      showToast({
        title: "Password updated",
        message: "Sign in with your new password.",
        variant: "success",
      });
      router.replace("/login");
    } catch (error) {
      showToast({
        title: "Reset failed",
        message: getApiMessage(error, "Unable to reset password."),
        variant: "error",
      });
    } finally {
      submitting.current = false;
    }
  };

  return (
    <SafeAreaView style={styles.screen}>
      <AppStatusBar />
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        style={styles.screen}
      >
        <ScrollView
          bounces={false}
          contentContainerStyle={styles.content}
          keyboardShouldPersistTaps="handled"
        >
          <View style={styles.hero}>
            <Image
              source={require("../../../assets/images/public/soc_img_two.png")}
              contentFit="cover"
              contentPosition="center"
              style={StyleSheet.absoluteFill}
            />
            {/* <LinearGradient
              colors={["rgba(23,17,15,0.05)", "rgba(23,17,15,0.78)"]}
              style={[StyleSheet.absoluteFill, styles.noPointerEvents]}
            /> */}
            <View style={styles.heroCopy}>
              <Text style={styles.heroTitle}>APNA GATE</Text>
              <Text style={styles.heroSubtitle}>
                Modern Security for Modern Societies.
              </Text>
            </View>
          </View>

          <View style={styles.sheet}>
            <View style={styles.header}>
              <Text style={styles.title}>
                {step === "email" ? "Forgot Password" : "Reset Password"}
              </Text>
              <Text style={styles.subtitle}>
                {step === "email"
                  ? "Enter your email and we will send a reset code."
                  : `Enter the code sent to ${normalizedEmail}.`}
              </Text>
            </View>

            <View style={styles.divider} />

            <View style={styles.form}>
              {step === "email" ? (
                <View style={styles.field}>
                  <Text style={styles.label}>Email</Text>
                  <TextInput
                    autoCapitalize="none"
                    autoComplete="email"
                    cursorColor="#ff6a1a"
                    editable={!isBusy}
                    keyboardType="email-address"
                    onChangeText={setEmail}
                    placeholder="name@example.com"
                    placeholderTextColor="#aaa19a"
                    selectionColor="#ffc39a"
                    style={styles.input}
                    value={email}
                  />
                </View>
              ) : (
                <>
                  <View style={styles.field}>
                    <Text style={styles.label}>Reset Code</Text>
                    <TextInput
                      autoCapitalize="none"
                      cursorColor="#ff6a1a"
                      editable={!isBusy}
                      keyboardType="number-pad"
                      onChangeText={setOtp}
                      placeholder="Enter 6-digit code"
                      placeholderTextColor="#aaa19a"
                      selectionColor="#ffc39a"
                      style={styles.input}
                      value={otp}
                    />
                  </View>

                  <View style={styles.field}>
                    <Text style={styles.label}>New Password</Text>
                    <View style={styles.passwordRow}>
                      <TextInput
                        autoCapitalize="none"
                        cursorColor="#ff6a1a"
                        editable={!isBusy}
                        onChangeText={setNewPassword}
                        placeholder="At least 8 characters"
                        placeholderTextColor="#aaa19a"
                        secureTextEntry={!showNewPassword}
                        selectionColor="#ffc39a"
                        style={styles.passwordInput}
                        value={newPassword}
                      />
                      <Pressable
                        accessibilityLabel={
                          showNewPassword ? "Hide password" : "Show password"
                        }
                        accessibilityRole="button"
                        onPress={() =>
                          setShowNewPassword((current) => !current)
                        }
                        style={styles.eyeButton}
                      >
                        <SymbolView
                          name={
                            showNewPassword
                              ? {
                                  ios: "eye.slash",
                                  android: "visibility_off",
                                  web: "visibility_off",
                                }
                              : {
                                  ios: "eye",
                                  android: "visibility",
                                  web: "visibility",
                                }
                          }
                          size={18}
                          tintColor="#81766f"
                        />
                      </Pressable>
                    </View>
                  </View>

                  <View style={styles.field}>
                    <Text style={styles.label}>Confirm Password</Text>
                    <View style={styles.passwordRow}>
                      <TextInput
                        autoCapitalize="none"
                        cursorColor="#ff6a1a"
                        editable={!isBusy}
                        onChangeText={setConfirmPassword}
                        placeholder="Re-enter new password"
                        placeholderTextColor="#aaa19a"
                        secureTextEntry={!showConfirmPassword}
                        selectionColor="#ffc39a"
                        style={styles.passwordInput}
                        value={confirmPassword}
                      />
                      <Pressable
                        accessibilityLabel={
                          showConfirmPassword
                            ? "Hide password"
                            : "Show password"
                        }
                        accessibilityRole="button"
                        onPress={() =>
                          setShowConfirmPassword((current) => !current)
                        }
                        style={styles.eyeButton}
                      >
                        <SymbolView
                          name={
                            showConfirmPassword
                              ? {
                                  ios: "eye.slash",
                                  android: "visibility_off",
                                  web: "visibility_off",
                                }
                              : {
                                  ios: "eye",
                                  android: "visibility",
                                  web: "visibility",
                                }
                          }
                          size={18}
                          tintColor="#81766f"
                        />
                      </Pressable>
                    </View>
                  </View>
                </>
              )}
            </View>

            <View style={styles.secureNote}>
              <Text style={styles.secureIcon}>{"\u{1F6E1}\uFE0F"}</Text>
              <View style={styles.secureCopy}>
                <Text style={styles.secureTitle}>Secure reset</Text>
                <Text style={styles.helper}>
                  Your new password is encrypted before storage.
                </Text>
              </View>
            </View>

            <Pressable
              accessibilityRole="button"
              disabled={isBusy}
              onPress={step === "email" ? handleSendCode : handleResetPassword}
              style={[styles.button, isBusy && styles.buttonDisabled]}
            >
              {isBusy ? (
                <ActivityIndicator color="#ffffff" />
              ) : (
                <Text style={styles.buttonText}>
                  {step === "email" ? "Send Reset Code" : "Update Password"}
                </Text>
              )}
            </Pressable>

            <View style={styles.divider} />

            <View style={styles.help}>
              {step === "reset" ? (
                <Pressable
                  accessibilityRole="button"
                  disabled={isBusy}
                  onPress={() => setStep("email")}
                >
                  <Text style={styles.linkText}>Use a different email</Text>
                </Pressable>
              ) : null}
              <Pressable
                accessibilityRole="button"
                disabled={isBusy}
                onPress={() => router.replace("/login")}
              >
                <Text style={styles.linkText}>Back to sign in</Text>
              </Pressable>
            </View>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  button: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: 16,
    boxShadow: "0 10px 20px rgba(255, 106, 26, 0.18)",
    justifyContent: "center",
    minHeight: 54,
  },
  buttonDisabled: {
    opacity: 0.75,
  },
  buttonText: {
    color: "#ffffff",
    fontSize: 16,
    fontWeight: "700",
  },
  content: {
    flexGrow: 1,
  },
  divider: {
    backgroundColor: "#eee7e2",
    height: 1,
  },
  eyeButton: {
    alignItems: "center",
    height: 52,
    justifyContent: "center",
    paddingHorizontal: 12,
  },
  field: {
    gap: 8,
  },
  form: {
    gap: 14,
  },
  header: {
    gap: 8,
  },
  helper: {
    color: colors.text.muted,
    fontSize: 13,
    lineHeight: 19,
  },
  help: {
    alignItems: "center",
    gap: 10,
  },
  hero: {
    height: 260,
    overflow: "hidden",
  },
  heroCopy: {
    alignItems: "center",
    bottom: 34,
    gap: 10,
    left: 24,
    position: "absolute",
    right: 24,
  },
  heroSubtitle: {
    color: "rgba(255,255,255,0.9)",
    fontSize: 12,
    fontWeight: "700",
    letterSpacing: 0.6,
    textAlign: "center",
    textTransform: "uppercase",
  },
  heroTitle: {
    color: "#ffffff",
    fontSize: 30,
    fontWeight: "700",
    letterSpacing: 0.8,
    textAlign: "center",
  },
  input: {
    backgroundColor: "#ffffff",
    borderColor: colors.border.input,
    borderRadius: 14,
    borderWidth: 1,
    color: colors.text.primary,
    fontSize: 15,
    minHeight: 52,
    paddingHorizontal: 15,
  },
  label: {
    color: colors.text.primary,
    fontSize: 14,
    fontWeight: "600",
  },
  linkText: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "600",
  },
  noPointerEvents: {
    pointerEvents: "none",
  },
  passwordInput: {
    color: colors.text.primary,
    flex: 1,
    fontSize: 15,
    minHeight: 52,
    paddingLeft: 15,
  },
  passwordRow: {
    alignItems: "center",
    backgroundColor: "#ffffff",
    borderColor: colors.border.input,
    borderRadius: 14,
    borderWidth: 1,
    flexDirection: "row",
  },
  screen: {
    backgroundColor: colors.surface.screen,
    flex: 1,
  },
  secureCopy: {
    flex: 1,
    gap: 3,
  },
  secureIcon: {
    fontSize: 20,
    lineHeight: 24,
  },
  secureNote: {
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 14,
    flexDirection: "row",
    gap: 12,
    padding: 14,
  },
  secureTitle: {
    color: colors.text.primary,
    fontSize: 15,
    fontWeight: "700",
  },
  sheet: {
    backgroundColor: "#ffffff",
    borderTopLeftRadius: 28,
    borderTopRightRadius: 28,
    boxShadow: "0 -10px 24px rgba(0, 0, 0, 0.1)",
    flexGrow: 1,
    gap: 22,
    marginTop: -24,
    paddingBottom: 28,
    paddingHorizontal: 24,
    paddingTop: 30,
  },
  subtitle: {
    color: colors.text.secondary,
    fontSize: 15,
    lineHeight: 22,
  },
  title: {
    color: colors.text.primary,
    fontSize: 28,
    fontWeight: "600",
    lineHeight: 34,
  },
});
