import { colors } from "@/theme/colors";
import { Image } from "expo-image";
import { SymbolView } from "expo-symbols";
import { useRef, useState } from "react";
import {
  ActivityIndicator,
  Linking,
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
import { getApiMessage } from "@/features/auth/api-error";
import { useCompleteAuth } from "@/features/auth/use-complete-auth";
import { usePostV1AuthLoginMutation } from "@/lib/api/generated-api";
import { Redirect, useRouter } from "expo-router";

import { useAuth } from "@/features/auth/use-auth";
import { pendingNotificationRoute } from "@/features/web/notification-handoff";
import { useToast } from "@/components/ui";

const COUNTRY_CODE = "+91";

type LoginMethod = "email" | "phone";

export default function LoginScreen() {
  const [loginMethod, setLoginMethod] = useState<LoginMethod>("email");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [login, { isLoading: isLoggingIn }] = usePostV1AuthLoginMutation();
  const completeAuth = useCompleteAuth();
  const { homeRoute, status } = useAuth();
  const { showToast } = useToast();
  const router = useRouter();
  const submitting = useRef(false);
  const [completing, setCompleting] = useState(false);
  const isLoading = isLoggingIn || completing;

  if (status === "authenticated" && homeRoute) {
    return <Redirect href={pendingNotificationRoute() ?? homeRoute} />;
  }

  const handleForget = () => {
    router.push("/forget");
  };

  const handlePhoneChange = (text: string) => {
    // keep digits only, cap at 10 (standard Indian mobile number length)
    const digitsOnly = text.replace(/[^0-9]/g, "").slice(0, 10);
    setPhone(digitsOnly);
  };

  const handleLogin = async () => {
    if (submitting.current) return;
    if (!password) {
      showToast({
        title: "Missing details",
        message: "Enter your password.",
        variant: "error",
      });
      return;
    }

    if (loginMethod === "email" && !email.trim()) {
      showToast({
        title: "Missing details",
        message: "Enter your email.",
        variant: "error",
      });
      return;
    }

    if (loginMethod === "phone" && phone.trim().length < 10) {
      showToast({
        title: "Missing details",
        message: "Enter a valid 10-digit phone number.",
        variant: "error",
      });
      return;
    }

    submitting.current = true;
    setCompleting(true);
    try {
      const response = await login({
        modelsLoginRequest:
          loginMethod === "email"
            ? {
                email: email.trim().toLowerCase(),
                password,
              }
            : {
                phone_number: `${COUNTRY_CODE}${phone.trim()}`,
                password,
              },
      }).unwrap();

      await completeAuth(response.data ?? undefined);
    } catch (error) {
      showToast({
        title: "Sign in failed",
        message: getApiMessage(error, "Check your email/phone and password."),
        variant: "error",
      });
    } finally {
      submitting.current = false;
      setCompleting(false);
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
            <View style={styles.heroCopy}>
              <Text style={styles.heroTitle}>APNA GATE</Text>
              <Text style={styles.heroSubtitle}>
                Modern Security for Modern Societies.
              </Text>
            </View>
          </View>

          <View style={styles.sheet}>
            <View style={styles.header}>
              <Text style={styles.title}>Welcome Back to Apna Gate</Text>
              <Text style={styles.subtitle}>
                Continue securely to your community.
              </Text>
            </View>

            <View style={styles.divider} />

            {/* <View style={styles.notice}>
              <Text style={styles.noticeIcon}>{"\u24D8"}</Text>
              <View style={styles.noticeCopy}>
                <Text style={styles.noticeTitle}>New here?</Text>
                <Text style={styles.helper}>
                  Residents can login after their flat claim is approved.
                </Text>
              </View>
            </View> */}
            <Text style={styles.helper}>
              To claim your flat, visit{" "}
              <Text
                style={styles.link}
                onPress={() =>
                  Linking.openURL("https://web.apnagate.org/claim")
                }
              >
                apnagate.org/claim
              </Text>{" "}
              and enter your society code.
            </Text>

            <View style={styles.divider} />

            {/* Login method toggle */}
            <View style={styles.toggleRow}>
              <Pressable
                accessibilityRole="button"
                onPress={() => setLoginMethod("email")}
                style={[
                  styles.toggleButton,
                  loginMethod === "email" && styles.toggleButtonActive,
                ]}
              >
                <Text
                  style={[
                    styles.toggleText,
                    loginMethod === "email" && styles.toggleTextActive,
                  ]}
                >
                  Email
                </Text>
              </Pressable>
              <Pressable
                accessibilityRole="button"
                onPress={() => setLoginMethod("phone")}
                style={[
                  styles.toggleButton,
                  loginMethod === "phone" && styles.toggleButtonActive,
                ]}
              >
                <Text
                  style={[
                    styles.toggleText,
                    loginMethod === "phone" && styles.toggleTextActive,
                  ]}
                >
                  Phone
                </Text>
              </Pressable>
            </View>

            <View style={styles.form}>
              {loginMethod === "email" ? (
                <View style={styles.field}>
                  <Text style={styles.label}>Email</Text>
                  <TextInput
                    autoCapitalize="none"
                    autoComplete="email"
                    cursorColor="#ff6a1a"
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
                <View style={styles.field}>
                  <Text style={styles.label}>Phone Number</Text>
                  <View style={styles.phoneRow}>
                    <View style={styles.countryCode}>
                      <Text style={styles.countryCodeText}>{COUNTRY_CODE}</Text>
                    </View>
                    <TextInput
                      cursorColor="#ff6a1a"
                      keyboardType="number-pad"
                      maxLength={10}
                      onChangeText={handlePhoneChange}
                      placeholder="98765 43210"
                      placeholderTextColor="#aaa19a"
                      selectionColor="#ffc39a"
                      style={styles.phoneInput}
                      value={phone}
                    />
                  </View>
                </View>
              )}

              <View style={styles.field}>
                <Text style={styles.label}>Password</Text>
                <View style={styles.passwordRow}>
                  <TextInput
                    autoCapitalize="none"
                    cursorColor="#ff6a1a"
                    onChangeText={setPassword}
                    placeholder="Enter password"
                    placeholderTextColor="#aaa19a"
                    secureTextEntry={!showPassword}
                    selectionColor="#ffc39a"
                    style={styles.passwordInput}
                    value={password}
                  />
                  <Pressable
                    accessibilityLabel={
                      showPassword ? "Hide password" : "Show password"
                    }
                    accessibilityRole="button"
                    onPress={() => setShowPassword((current) => !current)}
                    style={styles.eyeButton}
                  >
                    <SymbolView
                      name={
                        showPassword
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

              <Pressable
                accessibilityRole="button"
                onPress={handleForget}
                style={styles.forgotLink}
              >
                <Text style={styles.forgotText}>Forgot Password?</Text>
              </Pressable>
            </View>

            <View style={styles.secureNote}>
              <Text style={styles.secureIcon}>{"\u{1F6E1}\uFE0F"}</Text>
              <View style={styles.secureCopy}>
                <Text style={styles.secureTitle}>Secure login</Text>
                <Text style={styles.helper}>
                  Your information is encrypted.
                </Text>
              </View>
            </View>

            <Pressable
              accessibilityRole="button"
              disabled={isLoading}
              onPress={handleLogin}
              style={[styles.button, isLoading && styles.buttonDisabled]}
            >
              {isLoading ? (
                <ActivityIndicator color="#ffffff" />
              ) : (
                <Text style={styles.buttonText}>Secure Sign In</Text>
              )}
            </Pressable>

            <View style={styles.divider} />

            <View style={styles.help}>
              <Text style={styles.helpTitle}>Need help?</Text>
              <Text style={styles.helper}>
                Contact your society administrator.
              </Text>
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
  countryCode: {
    alignItems: "center",
    borderRightColor: "#e4dcd6",
    borderRightWidth: 1,
    justifyContent: "center",
    paddingLeft: 15,
    paddingRight: 10,
  },
  countryCodeText: {
    color: colors.text.primary,
    fontSize: 15,
    fontWeight: "600",
  },
  divider: {
    backgroundColor: "#eee7e2",
    height: 1,
  },
  eyeButton: {
    alignItems: "center",
    flexDirection: "row",
    gap: 5,
    height: 52,
    justifyContent: "center",
    paddingHorizontal: 12,
  },
  eyeText: {
    color: colors.text.muted,
    fontSize: 13,
    fontWeight: "600",
  },
  field: {
    gap: 8,
  },
  forgotLink: {
    alignSelf: "flex-end",
    paddingVertical: 2,
  },
  forgotText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "600",
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
    gap: 5,
  },
  helpTitle: {
    color: colors.text.primary,
    fontSize: 14,
    fontWeight: "700",
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
  noPointerEvents: {
    pointerEvents: "none",
  },
  link: {
    color: colors.brand.orange,
    fontWeight: "600",
  },
  notice: {
    flexDirection: "row",
    gap: 12,
  },
  noticeCopy: {
    flex: 1,
    gap: 3,
  },
  noticeIcon: {
    color: colors.brand.orange,
    fontSize: 20,
    lineHeight: 24,
  },
  noticeTitle: {
    color: colors.text.primary,
    fontSize: 15,
    fontWeight: "700",
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
  phoneInput: {
    color: colors.text.primary,
    flex: 1,
    fontSize: 15,
    minHeight: 52,
    paddingHorizontal: 12,
  },
  phoneRow: {
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
  toggleButton: {
    alignItems: "center",
    borderRadius: 12,
    paddingVertical: 10,
    flex: 1,
  },
  toggleButtonActive: {
    backgroundColor: "#ffffff",
    boxShadow: "0 2px 6px rgba(0, 0, 0, 0.08)",
  },
  toggleRow: {
    backgroundColor: "#f5f1ed",
    borderRadius: 14,
    flexDirection: "row",
    gap: 4,
    padding: 4,
  },
  toggleText: {
    color: colors.text.muted,
    fontSize: 14,
    fontWeight: "600",
  },
  toggleTextActive: {
    color: colors.text.primary,
  },
});
