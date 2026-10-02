import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Image } from "expo-image";
import { Redirect, useLocalSearchParams, useRouter } from "expo-router";
import { LoadingState } from "@/components/ui/loading-state";
import { pendingNotificationRoute } from "@/features/web/notification-handoff";
import { validNotificationLink } from "@/features/web/notification-target";
import { useEffect, useRef, useState } from "react";
import {
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  useWindowDimensions,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { typography } from "@/theme/typography";
import { useAuth } from "@/features/auth/auth-provider";

const images = [
  require("../../assets/images/public/soc_img_one.png"),
  require("../../assets/images/public/soc_img_two.png"),
];

const features = [
  ["\u{1F6E1}\uFE0F", "Secure Entry", "Approve visitors in seconds."],
  ["\u{1F4F1}", "Digital Access", "QR-based entry for residents and guests."],
  [
    "\u{1F465}",
    "Community Management",
    "Manage residents, staff, and visitors from one place.",
  ],
] as const;

const HERO_BG = "#ffff";
const DIVIDER = "#eee7e2";
const FEATURE_DESC = "#81766f";
const FOOTER_MUTED = "#c9c1bb";

export default function Index() {
  const { width } = useWindowDimensions();
  const router = useRouter();
  const { homeRoute, status } = useAuth();
  const params = useLocalSearchParams<{ notification_id?: string; recipient?: string }>();
  const carouselRef = useRef<ScrollView>(null);
  const [activeSlide, setActiveSlide] = useState(0);
  useEffect(() => {
    const timer = setInterval(() => {
      setActiveSlide((current) => {
        const next = (current + 1) % images.length;
        carouselRef.current?.scrollTo({ x: next * width, animated: true });
        return next;
      });
    }, 3500);
    return () => clearInterval(timer);
  }, [width]);

  if (Platform.OS === "web") {
    if (status === "loading") return <LoadingState message="Loading Apna Gate" />;
    if (params.notification_id !== undefined || params.recipient !== undefined) {
      return <Redirect href={validNotificationLink(params.notification_id, params.recipient)
        ? { pathname: "/notification-open", params: { id: params.notification_id!, recipient: params.recipient! } }
        : "/notifications"} />;
    }
    const pending = status === "authenticated" ? pendingNotificationRoute() : null;
    if (pending) return <Redirect href={pending} />;
  }

  if (status === "authenticated" && homeRoute) {
    return <Redirect href={homeRoute} />;
  }

  return (
    <SafeAreaView style={styles.screen}>
      <AppStatusBar />
      <ScrollView bounces={false} contentContainerStyle={styles.scrollContent}>
        <View style={styles.hero}>
          <ScrollView
            ref={carouselRef}
            horizontal
            pagingEnabled
            bounces={false}
            showsHorizontalScrollIndicator={false}
            onMomentumScrollEnd={(event) =>
              setActiveSlide(
                Math.round(event.nativeEvent.contentOffset.x / width),
              )
            }
          >
            {images.map((source, index) => (
              <Image
                key={index}
                source={source}
                contentFit="cover"
                contentPosition="center"
                style={{ width, height: "100%" }}
              />
            ))}
          </ScrollView>

          <View style={styles.heroCopy}>
            <Text style={styles.heroTitle}>APNA GATE</Text>
            <Text style={styles.heroSubtitle}>
              Modern Security for Modern Societies.
            </Text>
          </View>

          <View style={styles.paginationWrap}>
            <View style={styles.pagination}>
              {images.map((_, index) => (
                <View
                  key={index}
                  style={[
                    styles.dot,
                    activeSlide === index
                      ? styles.dotActive
                      : styles.dotInactive,
                  ]}
                />
              ))}
            </View>
          </View>
        </View>

        <View style={styles.sheet}>
          <View style={styles.intro}>
            <Text style={styles.title}>Welcome to Apna Gate</Text>
            <Text style={styles.subtitle}>
              A smarter way to welcome visitors and keep your community secure.
            </Text>
          </View>

          <View style={styles.divider} />

          <View style={styles.features}>
            {features.map(([icon, title, description]) => (
              <View key={title} style={styles.featureRow}>
                <View style={styles.featureIconWrap}>
                  <Text style={styles.featureIcon}>{icon}</Text>
                </View>
                <View style={styles.featureCopy}>
                  <Text style={styles.featureTitle}>{title}</Text>
                  <Text style={styles.featureDescription}>{description}</Text>
                </View>
              </View>
            ))}
          </View>

          <View style={styles.divider} />

          <View style={styles.footer}>
            <Text style={styles.trustLabel}>
              Trusted by modern residential communities.
            </Text>
            <Pressable
              accessibilityRole="button"
              onPress={() => router.push("/login")}
              style={({ pressed }) => [
                styles.ctaButton,
                pressed && styles.ctaButtonPressed,
              ]}
            >
              <Text style={styles.ctaText}>Continue to your community.</Text>
            </Pressable>
            <View style={styles.legalRow}>
              <Text style={styles.legalText}>Terms</Text>
              <Text style={styles.legalSeparator}>{"\u2022"}</Text>
              <Text style={styles.legalText}>Privacy</Text>
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const heroTextShadow =
  Platform.OS === "web"
    ? { textShadow: "0 1px 6px rgba(0, 0, 0, 0.45)" as const }
    : {
        textShadowColor: "rgba(0, 0, 0, 0.45)",
        textShadowOffset: { width: 0, height: 1 },
        textShadowRadius: 6,
      };

const heroSubtitleTextShadow =
  Platform.OS === "web"
    ? { textShadow: "0 1px 5px rgba(0, 0, 0, 0.45)" as const }
    : {
        textShadowColor: "rgba(0, 0, 0, 0.45)",
        textShadowOffset: { width: 0, height: 1 },
        textShadowRadius: 5,
      };

const styles = StyleSheet.create({
  screen: {
    backgroundColor: HERO_BG,
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
  },
  noPointerEvents: {
    pointerEvents: "none",
  },
  hero: {
    backgroundColor: HERO_BG,
    height: 340,
    overflow: "hidden",
    position: "relative",
  },
  heroCopy: {
    bottom: 76,
    gap: spacing.sm,
    left: spacing["2xl"],
    position: "absolute",
    right: spacing["2xl"],
  },
  heroTitle: {
    color: colors.text.inverse,
    fontSize: 27,
    fontWeight: "500",
    lineHeight: 36,
    ...heroTextShadow,
  },
  heroSubtitle: {
    ...typography.label,
    color: "rgba(255, 255, 255, 0.9)",
    ...heroSubtitleTextShadow,
  },
  paginationWrap: {
    alignItems: "center",
    bottom: 42,
    left: 0,
    position: "absolute",
    right: 0,
  },
  pagination: {
    backgroundColor: "rgba(0, 0, 0, 0.35)",
    borderRadius: radius["2xl"],
    flexDirection: "row",
    gap: spacing.sm,
    justifyContent: "center",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  dot: {
    borderRadius: radius.sm,
    height: spacing.sm,
  },
  dotActive: {
    backgroundColor: colors.text.inverse,
    width: spacing["2xl"],
  },
  dotInactive: {
    backgroundColor: "rgba(255, 255, 255, 0.55)",
    width: spacing.sm,
  },
  sheet: {
    backgroundColor: colors.surface.card,
    borderTopLeftRadius: 28,
    borderTopRightRadius: 28,
    flexGrow: 1,
    gap: spacing["2xl"],
    marginTop: -20,
    paddingBottom: layout.screenPaddingBottom,
    paddingHorizontal: spacing["2xl"],
    paddingTop: spacing["3xl"],
    boxShadow: "0 -10px 24px rgba(0, 0, 0, 0.1)",
  },
  intro: {
    gap: spacing.sm,
  },
  title: {
    color: colors.text.secondaryDark,
    fontSize: 25,
    fontWeight: "500",
    lineHeight: 32,
  },
  subtitle: {
    color: colors.text.secondaryDark,
    fontSize: 15,
    lineHeight: 24,
  },
  divider: {
    backgroundColor: DIVIDER,
    height: 1,
  },
  features: {
    gap: spacing.xl,
  },
  featureRow: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: 14,
  },
  featureIconWrap: {
    alignItems: "center",
    backgroundColor: "#fff1e8",
    borderRadius: radius["2xl"],
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  featureIcon: {
    fontSize: 18,
  },
  featureCopy: {
    flex: 1,
    gap: spacing.xs,
  },
  featureTitle: {
    color: colors.text.secondaryDark,
    fontSize: 15,
    fontWeight: "600",
  },
  featureDescription: {
    color: FEATURE_DESC,
    fontSize: typography.bodySmall.fontSize,
    lineHeight: 20,
  },
  footer: {
    gap: spacing.lg,
  },
  trustLabel: {
    ...typography.label,
    color: colors.text.placeholder,
    letterSpacing: 0.8,
    textAlign: "center",
    textTransform: "uppercase",
  },
  ctaButton: {
    alignItems: "center",
    backgroundColor: colors.brand.primary,
    borderRadius: radius.lg,
    justifyContent: "center",
    minHeight: layout.buttonHeight,
    ...shadows.brand,
  },
  ctaButtonPressed: {
    opacity: 0.9,
  },
  ctaText: {
    ...typography.button,
    color: colors.text.inverse,
  },
  legalRow: {
    flexDirection: "row",
    gap: spacing.sm,
    justifyContent: "center",
  },
  legalText: {
    ...typography.caption,
    color: colors.text.placeholder,
  },
  legalSeparator: {
    ...typography.caption,
    color: FOOTER_MUTED,
  },
});
