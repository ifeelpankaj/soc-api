import type { ReactNode } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { SymbolViewProps } from "expo-symbols";

import { Row, Stack } from "@/components/layout";
import { UserAvatar } from "@/components/ui";
import { BrandMark } from "@/components/dashboard/brand-mark";
import { colors } from "@/theme/colors";

export type DashboardHeaderAction = {
  accessibilityLabel: string;
  icon: SymbolViewProps["name"];
  onPress: () => void;
  notificationCount?: number;
  showBadge?: boolean;
};

export type DashboardStatusItem = {
  label: string;
  live?: boolean;
};

export type DashboardProfileAvatar = {
  imageUrl?: string | null;
  name?: string | null;
  onPress?: () => void;
  showOnlineDot?: boolean;
};

type DashboardHeaderProps = {
  actions?: DashboardHeaderAction[];
  greeting?: string;
  leading?: ReactNode;
  profileAvatar?: DashboardProfileAvatar;
  showBrand?: boolean;
  statusItems?: DashboardStatusItem[];
  title: string;
};

function NotificationBadge({ count }: { count: number }) {
  const label = count > 99 ? "99+" : String(count);

  return (
    <View style={styles.notificationBadge}>
      <Text style={styles.notificationBadgeText}>{label}</Text>
    </View>
  );
}

function OnlineDot() {
  return (
    <View style={styles.onlineDot}>
      <View style={styles.onlineDotInner} />
    </View>
  );
}

export function DashboardHeader({
  actions = [],
  greeting,
  leading,
  profileAvatar,
  showBrand = true,
  statusItems = [],
  title,
}: DashboardHeaderProps) {
  return (
    <Stack gap="lg">
      <Row align="center" justify="space-between">
        {showBrand ? <BrandMark size="sm" /> : leading}
        <Row align="center" gap="sm" justify="flex-start">
          {actions.map((action) => (
            <Pressable
              key={action.accessibilityLabel}
              accessibilityLabel={action.accessibilityLabel}
              accessibilityRole="button"
              hitSlop={8}
              style={({ pressed }) => [styles.iconButton, pressed && styles.iconButtonPressed]}
              onPress={action.onPress}
            >
              <SymbolView name={action.icon} size={22} tintColor={colors.guard.text} />
              {action.notificationCount && action.notificationCount > 0 ? (
                <NotificationBadge count={action.notificationCount} />
              ) : action.showBadge ? (
                <View style={styles.notificationDot} />
              ) : null}
            </Pressable>
          ))}
          {profileAvatar ? (
            <View style={styles.avatarWrap}>
              <UserAvatar
                imageUrl={profileAvatar.imageUrl}
                name={profileAvatar.name}
                size={40}
                onPress={profileAvatar.onPress}
              />
              {profileAvatar.showOnlineDot !== false ? <OnlineDot /> : null}
            </View>
          ) : null}
        </Row>
      </Row>

      <Stack gap="xs">
        {greeting ? <Text style={styles.greeting}>{greeting}</Text> : null}
        {title ? (
          <Row align="center" gap="sm" justify="flex-start">
            <Text numberOfLines={1} style={styles.location}>
              {title}
            </Text>
            {statusItems.length > 0 ? (
              <>
                <Text style={styles.separator}>•</Text>
                {statusItems.map((item) => (
                  <Row key={item.label} align="center" gap={6} justify="flex-start">
                    {item.live !== undefined ? (
                      <>
                        <View
                          style={[
                            styles.liveDot,
                            {
                              backgroundColor: item.live
                                ? colors.status.success
                                : colors.guard.textMuted,
                            },
                          ]}
                        />
                        <Text
                          style={[
                            styles.statusLabel,
                            {
                              color: item.live ? colors.status.success : colors.guard.textMuted,
                            },
                          ]}
                        >
                          {item.label}
                        </Text>
                      </>
                    ) : (
                      <Text style={styles.statusLabelMuted}>{item.label}</Text>
                    )}
                  </Row>
                ))}
              </>
            ) : null}
          </Row>
        ) : null}
      </Stack>
    </Stack>
  );
}

const styles = StyleSheet.create({
  avatarWrap: {
    position: "relative",
  },
  greeting: {
    color: colors.brand.navy,
    fontSize: 26,
    fontWeight: "700",
    letterSpacing: -0.5,
  },
  iconButton: {
    alignItems: "center",
    height: 40,
    justifyContent: "center",
    position: "relative",
    width: 40,
  },
  iconButtonPressed: {
    opacity: 0.65,
  },
  liveDot: {
    borderRadius: 999,
    height: 8,
    width: 8,
  },
  location: {
    color: colors.guard.textMuted,
    fontSize: 14,
    fontWeight: "500",
    maxWidth: "60%",
  },
  notificationBadge: {
    alignItems: "center",
    backgroundColor: colors.status.error,
    borderColor: colors.surface.card,
    borderRadius: 999,
    borderWidth: 2,
    justifyContent: "center",
    minHeight: 18,
    minWidth: 18,
    paddingHorizontal: 4,
    position: "absolute",
    right: 2,
    top: 2,
  },
  notificationBadgeText: {
    color: colors.text.inverse,
    fontSize: 10,
    fontWeight: "700",
    lineHeight: 12,
  },
  notificationDot: {
    backgroundColor: colors.status.error,
    borderColor: colors.surface.card,
    borderRadius: 999,
    borderWidth: 2,
    height: 10,
    position: "absolute",
    right: 6,
    top: 6,
    width: 10,
  },
  onlineDot: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderRadius: 999,
    bottom: 0,
    height: 14,
    justifyContent: "center",
    position: "absolute",
    right: 0,
    width: 14,
  },
  onlineDotInner: {
    backgroundColor: colors.status.success,
    borderRadius: 999,
    height: 10,
    width: 10,
  },
  separator: {
    color: colors.guard.textMuted,
    fontSize: 14,
  },
  statusLabel: {
    fontSize: 13,
    fontWeight: "600",
  },
  statusLabelMuted: {
    color: colors.guard.textMuted,
    fontSize: 14,
    fontWeight: "500",
  },
});
