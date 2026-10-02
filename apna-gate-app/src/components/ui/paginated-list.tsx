import { type ReactNode, type ComponentType } from "react";
import {
  ActivityIndicator,
  FlatList,
  RefreshControl,
  StyleSheet,
  View,
  type ListRenderItem,
  type ViewStyle,
} from "react-native";

import { Stack } from "@/components/layout/stack";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { spacing } from "@/theme/spacing";

import { EmptyState } from "./empty-state";

type PaginatedListProps<T> = {
  data: T[];
  keyExtractor: (item: T, index: number) => string;
  renderItem: ListRenderItem<T>;
  isLoading?: boolean;
  isRefreshing?: boolean;
  isLoadingMore?: boolean;
  hasMore?: boolean;
  onRefresh?: () => void;
  onLoadMore?: () => void;
  emptyComponent?: ReactNode;
  loadingComponent?: ReactNode;
  emptyTitle?: string;
  emptyMessage?: string;
  header?: ReactNode;
  footer?: ReactNode;
  contentContainerStyle?: ViewStyle;
  ItemSeparatorComponent?: ComponentType<unknown> | null;
};

export function PaginatedList<T>({
  data,
  keyExtractor,
  renderItem,
  isLoading,
  isRefreshing,
  isLoadingMore,
  hasMore,
  onRefresh,
  onLoadMore,
  emptyComponent,
  loadingComponent,
  emptyTitle,
  emptyMessage,
  header,
  footer,
  contentContainerStyle,
  ItemSeparatorComponent,
}: PaginatedListProps<T>) {
  const listFooter = (
    <View>
      {isLoadingMore ? (
        <View style={styles.loadingMore}>
          <ActivityIndicator color={colors.guard.teal} />
        </View>
      ) : null}
      {footer}
    </View>
  );

  return (
    <FlatList
      style={{ flex: 1 }}
      data={data}
      ItemSeparatorComponent={ItemSeparatorComponent}
      keyboardShouldPersistTaps="handled"
      keyExtractor={keyExtractor}
      renderItem={renderItem}
      ListHeaderComponent={
        header ? (
          <Stack gap="lg" style={styles.header}>
            {header}
          </Stack>
        ) : (
          <View style={styles.headerSpacer} />
        )
      }
      ListFooterComponent={listFooter}
      ListEmptyComponent={
        <View>
          {isLoading ? (
            loadingComponent ?? (
              <View style={styles.loading}>
                <ActivityIndicator color={colors.guard.teal} />
              </View>
            )
          ) : (
            emptyComponent ?? (
              <EmptyState
                title={emptyTitle ?? "Nothing here yet"}
                message={emptyMessage ?? ""}
                onAction={onRefresh}
                actionLabel="Refresh"
              />
            )
          )}
        </View>
      }
      contentContainerStyle={[
        styles.content,
        data.length === 0 && styles.contentEmpty,
        contentContainerStyle,
      ]}
      refreshControl={
        onRefresh ? (
          <RefreshControl
            refreshing={Boolean(isRefreshing)}
            tintColor={colors.guard.teal}
            onRefresh={onRefresh}
          />
        ) : undefined
      }
      onEndReached={() => {
        if (hasMore && !isLoading && !isLoadingMore) {
          onLoadMore?.();
        }
      }}
      onEndReachedThreshold={0.4}
    />
  );
}

const styles = StyleSheet.create({
  content: {
    gap: spacing.lg,
    paddingBottom: spacing["2xl"],
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing["2xl"],
  },
  contentEmpty: {
    flexGrow: 1,
  },
  header: {
    paddingBottom: spacing.lg,
  },
  headerSpacer: {
    height: spacing["2xl"],
  },
  loading: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
    paddingVertical: spacing["3xl"],
  },
  loadingMore: {
    alignItems: "center",
    paddingVertical: spacing.lg,
  },
});
