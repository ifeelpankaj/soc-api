import type { ReactElement, ReactNode } from "react";
import { useCallback, useState } from "react";
import { useFocusEffect } from "expo-router";
import {
  BackHandler,
  KeyboardAvoidingView,
  Platform,
  StyleSheet,
  View,
  type StyleProp,
  type ViewStyle,
} from "react-native";
import {
  SafeAreaView,
  useSafeAreaInsets,
} from "react-native-safe-area-context";
import type { Href } from "expo-router";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Stack } from "@/components/layout";
import { EmptyState, PaginatedList } from "@/components/ui";
import { LogEntryDivider } from "@/features/visitors/components/log-entry-row";
import { VisitorEntriesSkeleton } from "@/features/visitors/components/visitor-entries-skeleton";
import { LogsDateSheet } from "@/features/visitors/components/logs-date-sheet";
import { LogsFilterSheet } from "@/features/visitors/components/logs-filter-sheet";
import { LogsSearchHeader } from "@/features/visitors/components/logs-search-header";
import { guardHomeRoute } from "@/features/guard/guard-routes";
import type { DateRangePreset } from "@/features/visitors/visitor-date-ranges";
import type {
  ModelsVisitorEntry,
  ModelsVisitorPurpose,
  ModelsVisitorStatus,
} from "@/lib/api/generated-api";
import { useBackAction } from "@/lib/navigation/use-back-action";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type VisitorEntriesListControls<TSegment extends string> = {
  activeFilterCount: number;
  applySheetFilters: (filters: {
    datePreset: DateRangePreset;
    purpose?: ModelsVisitorPurpose;
    status?: ModelsVisitorStatus;
  }) => void;
  clearSheetFilters: () => void;
  datePreset: DateRangePreset;
  hasMore: boolean;
  hasExplicitFilters: boolean;
  hasSearchQuery: boolean;
  isLoading: boolean;
  isLoadingMore: boolean;
  isRefreshing: boolean;
  isSearchActive: boolean;
  items: ModelsVisitorEntry[];
  loadMore: () => void;
  purpose?: ModelsVisitorPurpose;
  refresh: () => void;
  searchInput: string;
  filterSummary: string;
  segment: TSegment;
  selectSegment: (segment: TSegment) => void;
  setDatePreset: (preset: DateRangePreset) => void;
  setSearchInput: (value: string) => void;
  sheetStatus?: ModelsVisitorStatus;
};

type VisitorEntriesListScreenProps<TSegment extends string> = {
  backgroundColor?: string;
  controls: VisitorEntriesListControls<TSegment>;
  emptyMessage: string;
  emptyTitle: string;
  fallbackHomeRoute?: Href;
  footer?: ReactNode;
  headerExtra?: ReactNode;
  keyPrefix: string;
  listContentStyle?: StyleProp<ViewStyle>;
  renderRow: (entry: ModelsVisitorEntry) => ReactElement;
  searchTitle?: string;
  segmentOptions?: { label: string; value: TSegment }[];
  sheetFooter?: ReactNode;
};

export function VisitorEntriesListScreen<TSegment extends string>({
  backgroundColor = colors.surface.screen,
  controls,
  emptyMessage,
  emptyTitle,
  fallbackHomeRoute,
  footer,
  headerExtra,
  keyPrefix,
  listContentStyle,
  renderRow,
  searchTitle,
  segmentOptions,
  sheetFooter,
}: VisitorEntriesListScreenProps<TSegment>) {
  const [filterOpen, setFilterOpen] = useState(false);
  const [dateOpen, setDateOpen] = useState(false);
  const [searchMode, setSearchMode] = useState(false);
  const insets = useSafeAreaInsets();
  const listContentPaddingBottom = spacing["2xl"] + Math.max(insets.bottom, 0);
  const showSearchEmpty = searchMode && !controls.hasSearchQuery && !controls.hasExplicitFilters;
  const handleBack = useBackAction(fallbackHomeRoute ?? guardHomeRoute());

  const exitSearch = useCallback(() => {
    controls.setSearchInput("");
    setSearchMode(false);
  }, [controls]);

  const handleHardwareBack = useCallback(() => {
    if (filterOpen) {
      setFilterOpen(false);
      return;
    }

    if (dateOpen) {
      setDateOpen(false);
      return;
    }

    if (searchMode) {
      exitSearch();
      return;
    }

    handleBack();
  }, [dateOpen, exitSearch, filterOpen, handleBack, searchMode]);

  useFocusEffect(
    useCallback(() => {
      if (Platform.OS !== "android") {
        return undefined;
      }

      const subscription = BackHandler.addEventListener("hardwareBackPress", () => {
        handleHardwareBack();
        return true;
      });

      return () => subscription.remove();
    }, [handleHardwareBack]),
  );

  return (
    <SafeAreaView
      edges={["top", "left", "right", "bottom"]}
      style={[styles.screen, { backgroundColor }]}
    >
      <AppStatusBar />
      <KeyboardAvoidingView
        behavior={Platform.OS === "ios" ? "padding" : "height"}
        style={styles.screen}
      >
        <PaginatedList<ModelsVisitorEntry>
          ItemSeparatorComponent={LogEntryDivider}
          contentContainerStyle={StyleSheet.flatten([
            styles.listContent,
            { paddingBottom: listContentPaddingBottom },
            listContentStyle,
          ])}
          data={showSearchEmpty ? [] : controls.items}
          emptyComponent={
            showSearchEmpty ? (
              <EmptyState
                message="Find visitors by name, phone number, vehicle number, or flat."
                title="Search visitor logs"
              />
            ) : undefined
          }
          emptyMessage={emptyMessage}
          emptyTitle={emptyTitle}
          footer={footer ?? <View style={styles.footerSpacer} />}
          hasMore={showSearchEmpty ? false : controls.hasMore}
          header={
            <Stack gap="md" style={styles.header}>
              <LogsSearchHeader<TSegment>
                activeFilterCount={controls.activeFilterCount}
                datePreset={controls.datePreset}
                fallbackHomeRoute={fallbackHomeRoute}
                filterSummary={controls.filterSummary}
                hasExplicitFilters={controls.hasExplicitFilters}
                isSearchActive={controls.isSearchActive}
                isSearchMode={searchMode}
                searchValue={controls.searchInput}
                segment={controls.segment}
                segmentOptions={segmentOptions}
                title={searchTitle}
                onClearFilters={controls.clearSheetFilters}
                onDatePress={() => setDateOpen(true)}
                onExitSearch={exitSearch}
                onFilterPress={() => setFilterOpen(true)}
                onOpenSearch={() => setSearchMode(true)}
                onSearchChange={controls.setSearchInput}
                onSegmentChange={controls.selectSegment}
              />
              {headerExtra}
            </Stack>
          }
          isLoading={showSearchEmpty ? false : controls.isLoading}
          loadingComponent={<VisitorEntriesSkeleton />}
          isLoadingMore={showSearchEmpty ? false : controls.isLoadingMore}
          isRefreshing={controls.isRefreshing}
          keyExtractor={(item) => `${keyPrefix}-${item.id}`}
          renderItem={({ item }) => renderRow(item)}
          onLoadMore={controls.loadMore}
          onRefresh={controls.refresh}
        />
      </KeyboardAvoidingView>

      <LogsFilterSheet
        datePreset={controls.datePreset}
        purpose={controls.purpose}
        status={controls.sheetStatus}
        visible={filterOpen}
        onApply={controls.applySheetFilters}
        onClear={controls.clearSheetFilters}
        onClose={() => setFilterOpen(false)}
      />

      <LogsDateSheet
        selected={controls.datePreset}
        visible={dateOpen}
        onApply={controls.setDatePreset}
        onClose={() => setDateOpen(false)}
      />

      {sheetFooter}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  footerSpacer: {
    height: spacing.sm,
  },
  header: {
    paddingBottom: spacing.sm,
  },
  listContent: {
    gap: spacing.xs,
    paddingBottom: spacing["2xl"],
    paddingHorizontal: 0,
    paddingTop: spacing.md,
  },
  screen: {
    flex: 1,
  },
});

export type { VisitorEntriesListControls };
