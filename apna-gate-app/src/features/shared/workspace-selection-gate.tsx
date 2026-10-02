import { ActivityIndicator, ScrollView, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Stack } from "@/components/layout";
import { Button, Card, EmptyState, LoadingState, ScreenHeader, StatusPill } from "@/components/ui";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

export type WorkspaceSelectionItem = {
  id: string;
  onSelect: () => void;
  primaryLabel: string;
  secondaryLabel: string;
  selectLabel: string;
  status?: string | null;
};

type WorkspaceSelectionGateProps = {
  emptyActionLabel?: string;
  emptyMessage: string;
  emptyTitle: string;
  eyebrow: string;
  isLoading?: boolean;
  items: WorkspaceSelectionItem[];
  loadingMessage?: string;
  onRefresh?: () => void;
  subtitle: string;
  title: string;
  useLoadingState?: boolean;
};

export function WorkspaceSelectionGate({
  emptyActionLabel = "Refresh",
  emptyMessage,
  emptyTitle,
  eyebrow,
  isLoading = false,
  items,
  loadingMessage = "Opening workspace",
  onRefresh,
  subtitle,
  title,
  useLoadingState = false,
}: WorkspaceSelectionGateProps) {
  if (isLoading) {
    if (useLoadingState) {
      return <LoadingState message={loadingMessage} />;
    }

    return (
      <SafeAreaView style={styles.screen}>
        <AppStatusBar />
        <View style={styles.inlineLoading}>
          <ActivityIndicator color={colors.guard.teal} size="small" />
        </View>
      </SafeAreaView>
    );
  }

  if (items.length === 0) {
    if (useLoadingState) {
      return (
        <SafeAreaView style={styles.screen}>
          <AppStatusBar />
          <ScrollView contentContainerStyle={styles.scrollContent}>
            <EmptyState
              actionLabel={emptyActionLabel}
              message={emptyMessage}
              title={emptyTitle}
              onAction={onRefresh}
            />
          </ScrollView>
        </SafeAreaView>
      );
    }

    return (
      <SafeAreaView style={styles.screen}>
        <AppStatusBar />
        <ScrollView contentContainerStyle={styles.scrollContent}>
          <EmptyState
            actionLabel={emptyActionLabel}
            message={emptyMessage}
            title={emptyTitle}
            onAction={onRefresh}
          />
        </ScrollView>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.screen}>
      <AppStatusBar />
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Stack gap="3xl">
          <ScreenHeader eyebrow={eyebrow} subtitle={subtitle} title={title} />

          <Stack gap="md">
            {items.map((item) => (
              <Card key={item.id} style={styles.itemCard}>
                <View style={styles.itemHeader}>
                  <View style={styles.itemCopy}>
                    <Text style={styles.primaryLabel}>{item.primaryLabel}</Text>
                    <Text style={styles.secondaryLabel}>{item.secondaryLabel}</Text>
                  </View>
                  {item.status ? <StatusPill status={item.status} /> : null}
                </View>
                <Button title={item.selectLabel} onPress={item.onSelect} />
              </Card>
            ))}
          </Stack>
        </Stack>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  inlineLoading: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
  },
  itemCard: {
    gap: spacing.lg,
  },
  itemCopy: {
    flex: 1,
  },
  itemHeader: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: spacing.lg,
    justifyContent: "space-between",
  },
  primaryLabel: {
    color: colors.text.primary,
    fontSize: 20,
    fontWeight: "700",
  },
  screen: {
    backgroundColor: colors.surface.screen,
    flex: 1,
  },
  scrollContent: {
    paddingHorizontal: spacing["2xl"],
    paddingVertical: spacing["3xl"],
  },
  secondaryLabel: {
    color: colors.text.secondary,
    fontSize: 16,
    marginTop: spacing.xs,
    textTransform: "capitalize",
  },
});
