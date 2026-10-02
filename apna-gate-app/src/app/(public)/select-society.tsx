import { colors } from "@/theme/colors";
import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Image } from "expo-image";
import { useRouter } from "expo-router";
import {
  Platform,
  ActivityIndicator,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { saveWorkspace } from "@/features/auth/auth-storage";
import { guardHomeRoute } from "@/features/guard/guard-routes";
import {
  formatGlobalRole,
  formatMembershipRole,
  formatStatus,
  getAdminMembership,
  getMobileMemberships,
  getMobileResidences,
  requiresAdminPortal,
} from "@/features/auth/mobile-access";
import { useLogout } from "@/features/auth/use-logout";
import { getApiMessage } from "@/features/auth/api-error";
import type {
  ModelsFlatResidentResponse,
  ModelsSocietyMemberResponse,
} from "@/lib/api/generated-api";
import { useGetV1BootstrapQuery } from "@/lib/api/generated-api";
import { setSelectedFlat, setSelectedSociety } from "@/redux/appSlice";
import { useAppDispatch } from "@/redux/hooks";

function StatusBadge({ status }: { status?: string | null }) {
  const badgeStyle =
    status === "active"
      ? styles.statusBadgeActive
      : status === "pending"
        ? styles.statusBadgePending
        : styles.statusBadgeDefault;

  const textStyle =
    status === "active"
      ? styles.statusTextActive
      : status === "pending"
        ? styles.statusTextPending
        : styles.statusTextDefault;

  return (
    <View style={[styles.statusBadge, badgeStyle]}>
      <Text style={[styles.statusText, textStyle]}>{formatStatus(status)}</Text>
    </View>
  );
}

function WorkspaceCard({
  title,
  subtitle,
  status,
  actionLabel,
  disabled,
  onPress,
}: {
  title: string;
  subtitle: string;
  status?: string | null;
  actionLabel: string;
  disabled?: boolean;
  onPress: () => void;
}) {
  return (
    <View style={styles.card}>
      <View style={styles.cardHeader}>
        <View style={styles.cardCopy}>
          <Text style={styles.cardTitle}>{title}</Text>
          <Text style={styles.cardSubtitle}>{subtitle}</Text>
        </View>
        <StatusBadge status={status} />
      </View>
      <Pressable
        accessibilityRole="button"
        disabled={disabled}
        onPress={onPress}
        style={[styles.cardButton, disabled && styles.cardButtonDisabled]}
      >
        <Text style={styles.cardButtonText}>{actionLabel}</Text>
      </Pressable>
    </View>
  );
}

function ResidenceCard({
  residence,
  onPress,
}: {
  residence: ModelsFlatResidentResponse;
  onPress: () => void;
}) {
  const flatLabel = [
    residence.flat_number ? `Flat ${residence.flat_number}` : null,
    residence.block ? `Block ${residence.block}` : null,
  ]
    .filter(Boolean)
    .join(" · ");

  return (
    <WorkspaceCard
      actionLabel="Open resident home"
      disabled={residence.status !== "active"}
      onPress={onPress}
      status={residence.status}
      subtitle={flatLabel || "Linked residence"}
      title={residence.society_name ?? "Your society"}
    />
  );
}

function MembershipCard({
  membership,
  onPress,
}: {
  membership: ModelsSocietyMemberResponse;
  onPress: () => void;
}) {
  const isGuard = membership.role === "staff";

  return (
    <WorkspaceCard
      actionLabel={isGuard ? "Open guard desk" : "Open resident home"}
      disabled={membership.status !== "active"}
      onPress={onPress}
      status={membership.status}
      subtitle={`${formatMembershipRole(membership.role)} access`}
      title={`Society #${membership.society_id ?? "-"}`}
    />
  );
}

export default function SelectSocietyScreen() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const { signOut, isLoading: isSigningOut } = useLogout();
  const { data, isLoading, isFetching, isError, error, refetch } =
    useGetV1BootstrapQuery(undefined);

  const bootstrap = data?.data;
  const user = bootstrap?.user;
  const adminMembership = getAdminMembership(bootstrap?.memberships);
  const showAdminPortal =
    Platform.OS !== "web" && (requiresAdminPortal(user?.global_role) || !!adminMembership);
  const residences = getMobileResidences(bootstrap);
  const memberships = getMobileMemberships(bootstrap);
  const hasWorkspaces = residences.length > 0 || memberships.length > 0;

  const adminReason = requiresAdminPortal(user?.global_role)
    ? formatGlobalRole(user?.global_role)
    : formatMembershipRole(adminMembership?.role);

  const openResidence = (flatId: number) => {
    dispatch(setSelectedFlat(flatId));
    void saveWorkspace({ flatId });
    router.replace("/resident/dashboard");
  };

  const openMembership = (membership: ModelsSocietyMemberResponse) => {
    if (membership.role === "staff" && membership.society_id) {
      dispatch(setSelectedSociety(membership.society_id));
      void saveWorkspace({ societyId: membership.society_id });
      router.replace(guardHomeRoute());
      return;
    }

    router.replace("/resident/dashboard");
  };

  return (
    <SafeAreaView style={styles.screen}>
      <AppStatusBar />
      <ScrollView bounces={false} contentContainerStyle={styles.content}>
        <View style={styles.hero}>
          <Image
            source={require("../../../assets/images/public/soc_img_one.png")}
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
          {isLoading ? (
            <View style={styles.loadingBlock}>
              <ActivityIndicator color="#ff6a1a" size="large" />
              <Text style={styles.loadingText}>Checking your access...</Text>
            </View>
          ) : isError ? (
            <>
              <View style={styles.header}>
                <Text style={styles.title}>Unable to load access</Text>
                <Text style={styles.subtitle}>
                  {getApiMessage(error, "Check your connection and try again.")}
                </Text>
              </View>

              <View style={styles.divider} />

              <Pressable
                accessibilityRole="button"
                disabled={isFetching}
                onPress={refetch}
                style={[styles.button, isFetching && styles.buttonDisabled]}
              >
                {isFetching ? (
                  <ActivityIndicator color="#ffffff" />
                ) : (
                  <Text style={styles.buttonText}>Try again</Text>
                )}
              </Pressable>

              <Pressable
                accessibilityRole="button"
                disabled={isSigningOut}
                onPress={signOut}
                style={styles.refreshLink}
              >
                <Text style={styles.linkText}>Sign out</Text>
              </Pressable>
            </>
          ) : showAdminPortal ? (
            <>
              <View style={styles.header}>
                <Text style={styles.title}>Use the admin portal</Text>
                <Text style={styles.subtitle}>
                  This mobile app is for residents and guards only.
                </Text>
              </View>

              <View style={styles.divider} />

              <View style={styles.notice}>
                <Text style={styles.noticeIcon}>{"\u24D8"}</Text>
                <View style={styles.noticeCopy}>
                  <Text style={styles.noticeTitle}>{adminReason} account</Text>
                  <Text style={styles.helper}>
                    Please sign in through the web admin portal to manage
                    societies, settings, and staff.
                  </Text>
                </View>
              </View>

              <Pressable
                accessibilityRole="button"
                disabled={isSigningOut}
                onPress={signOut}
                style={[styles.button, isSigningOut && styles.buttonDisabled]}
              >
                {isSigningOut ? (
                  <ActivityIndicator color="#ffffff" />
                ) : (
                  <Text style={styles.buttonText}>Sign out</Text>
                )}
              </Pressable>
            </>
          ) : hasWorkspaces ? (
            <>
              <View style={styles.header}>
                <Text style={styles.title}>Choose your community</Text>
                <Text style={styles.subtitle}>
                  {Platform.OS === "web" ? "Select your resident workspace to continue." : "Select a residence or guard workspace to continue."}
                </Text>
              </View>

              <View style={styles.divider} />

              {residences.length > 0 ? (
                <View style={styles.section}>
                  <Text style={styles.sectionLabel}>Residences</Text>
                  {residences.map((residence) => (
                    <ResidenceCard
                      key={`residence-${residence.id ?? residence.flat_id}`}
                      onPress={() => {
                        if (residence.flat_id) {
                          openResidence(residence.flat_id);
                        }
                      }}
                      residence={residence}
                    />
                  ))}
                </View>
              ) : null}

              {memberships.length > 0 ? (
                <View style={styles.section}>
                  <Text style={styles.sectionLabel}>Workspaces</Text>
                  {memberships.map((membership) => (
                    <MembershipCard
                      key={`membership-${membership.id ?? membership.society_id}`}
                      membership={membership}
                      onPress={() => openMembership(membership)}
                    />
                  ))}
                </View>
              ) : null}

              <Pressable
                accessibilityRole="button"
                disabled={isFetching}
                onPress={refetch}
                style={styles.refreshLink}
              >
                <Text style={styles.linkText}>
                  {isFetching ? "Refreshing..." : "Refresh access"}
                </Text>
              </Pressable>
            </>
          ) : (
            <>
              <View style={styles.header}>
                <Text style={styles.title}>Access is pending</Text>
                <Text style={styles.subtitle}>
                  {Platform.OS === "web" ? "No active resident workspace is linked. Guard workspaces require the Android app." : "Your account is ready, but no active resident or guard access is linked yet."}
                </Text>
              </View>

              <View style={styles.divider} />

              <View style={styles.notice}>
                <Text style={styles.noticeIcon}>{"\u24D8"}</Text>
                <View style={styles.noticeCopy}>
                  <Text style={styles.noticeTitle}>Waiting for approval</Text>
                  <Text style={styles.helper}>
                    Ask your society administrator to approve your resident or
                    guard access, then refresh this screen.
                  </Text>
                </View>
              </View>

              <Pressable
                accessibilityRole="button"
                disabled={isFetching}
                onPress={refetch}
                style={[styles.button, isFetching && styles.buttonDisabled]}
              >
                {isFetching ? (
                  <ActivityIndicator color="#ffffff" />
                ) : (
                  <Text style={styles.buttonText}>Refresh access</Text>
                )}
              </Pressable>

              <Pressable
                accessibilityRole="button"
                disabled={isSigningOut}
                onPress={signOut}
                style={styles.refreshLink}
              >
                <Text style={styles.linkText}>Sign out</Text>
              </Pressable>
            </>
          )}
        </View>
      </ScrollView>
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
  card: {
    backgroundColor: "#fffaf6",
    borderColor: "#f1e4da",
    borderRadius: 18,
    borderWidth: 1,
    gap: 14,
    padding: 16,
  },
  cardButton: {
    alignItems: "center",
    backgroundColor: colors.brand.orange,
    borderRadius: 14,
    justifyContent: "center",
    minHeight: 46,
  },
  cardButtonDisabled: {
    opacity: 0.55,
  },
  cardButtonText: {
    color: "#ffffff",
    fontSize: 15,
    fontWeight: "700",
  },
  cardCopy: {
    flex: 1,
    gap: 4,
  },
  cardHeader: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: 12,
  },
  cardSubtitle: {
    color: colors.text.secondary,
    fontSize: 14,
    lineHeight: 20,
  },
  cardTitle: {
    color: colors.text.primary,
    fontSize: 18,
    fontWeight: "700",
  },
  content: {
    flexGrow: 1,
  },
  divider: {
    backgroundColor: "#eee7e2",
    height: 1,
  },
  header: {
    gap: 8,
  },
  helper: {
    color: colors.text.muted,
    fontSize: 13,
    lineHeight: 19,
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
  linkText: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "600",
  },
  loadingBlock: {
    alignItems: "center",
    gap: 14,
    paddingVertical: 48,
  },
  loadingText: {
    color: colors.text.secondary,
    fontSize: 15,
  },
  noPointerEvents: {
    pointerEvents: "none",
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
  refreshLink: {
    alignItems: "center",
    paddingVertical: 4,
  },
  screen: {
    backgroundColor: colors.surface.screen,
    flex: 1,
  },
  section: {
    gap: 12,
  },
  sectionLabel: {
    color: "#a89e97",
    fontSize: 12,
    fontWeight: "700",
    letterSpacing: 0.8,
    textTransform: "uppercase",
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
  statusBadge: {
    borderRadius: 999,
    paddingHorizontal: 10,
    paddingVertical: 5,
  },
  statusBadgeActive: {
    backgroundColor: "#ecfdf3",
  },
  statusBadgeDefault: {
    backgroundColor: "#f3f4f6",
  },
  statusBadgePending: {
    backgroundColor: "#fff7ed",
  },
  statusText: {
    fontSize: 11,
    fontWeight: "700",
  },
  statusTextActive: {
    color: "#15803d",
  },
  statusTextDefault: {
    color: "#4b5563",
  },
  statusTextPending: {
    color: "#c2410c",
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
