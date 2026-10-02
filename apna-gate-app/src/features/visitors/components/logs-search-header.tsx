import { useEffect, useRef } from "react";
import { Pressable, Platform, StyleSheet, Text, TextInput, View } from "react-native";
import { SymbolView } from "expo-symbols";
import type { Href } from "expo-router";

import { ScreenBackHeader } from "@/components/layout/screen-back-header";
import { SegmentTabs } from "@/components/ui";
import { guardHomeRoute } from "@/features/guard/guard-routes";
import {
  getDateRangeLabel,
  type DateRangePreset,
  type LogsSegment,
} from "@/features/visitors/visitor-date-ranges";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";

const DEFAULT_SEGMENT_OPTIONS: { label: string; value: LogsSegment }[] = [
  { label: "Today", value: "today" },
  { label: "Expected", value: "expected" },
  { label: "Inside", value: "inside" },
  { label: "All", value: "all" },
];

type LogsSearchHeaderProps<T extends string = LogsSegment> = {
  activeFilterCount: number;
  datePreset: DateRangePreset;
  fallbackHomeRoute?: Href;
  filterSummary: string;
  hasExplicitFilters: boolean;
  isSearchActive: boolean;
  isSearchMode: boolean;
  onClearFilters: () => void;
  onDatePress: () => void;
  onExitSearch: () => void;
  onFilterPress: () => void;
  onOpenSearch: () => void;
  onSearchChange: (value: string) => void;
  onSegmentChange: (segment: T) => void;
  searchValue: string;
  segment: T;
  segmentOptions?: { label: string; value: T }[];
  title?: string;
};

export function LogsSearchHeader<T extends string = LogsSegment>({
  activeFilterCount,
  datePreset,
  fallbackHomeRoute,
  filterSummary,
  hasExplicitFilters,
  isSearchActive,
  isSearchMode,
  onClearFilters,
  onDatePress,
  onExitSearch,
  onFilterPress,
  onOpenSearch,
  onSearchChange,
  onSegmentChange,
  searchValue,
  segment,
  segmentOptions,
  title = "Visitor Logs",
}: LogsSearchHeaderProps<T>) {
  const searchRef = useRef<TextInput>(null);
  const options = segmentOptions ?? (DEFAULT_SEGMENT_OPTIONS as { label: string; value: T }[]);
  const hideDatePicker = segment === "expected" || segment === "today" || segment === "all";
  const showPresetControls = !isSearchMode && !hasExplicitFilters;

  useEffect(() => {
    if (isSearchMode) {
      const timer = setTimeout(() => searchRef.current?.focus(), 80);
      return () => clearTimeout(timer);
    }
  }, [isSearchMode]);

  return (
    <View style={styles.container}>
      <ScreenBackHeader
        fallbackHomeRoute={fallbackHomeRoute ?? guardHomeRoute()}
        title={isSearchMode ? "Search Visitors" : hasExplicitFilters ? "Filtered Results" : title}
        trailing={
          isSearchMode ? (
            <Pressable
              accessibilityRole="button"
              hitSlop={8}
              style={styles.headerAction}
              onPress={onExitSearch}
            >
              <Text style={styles.headerActionText}>Cancel</Text>
            </Pressable>
          ) : null
        }
      />

      <Pressable
        accessibilityRole="search"
        style={styles.searchRow}
        onPress={isSearchMode ? undefined : onOpenSearch}
      >
        <SymbolView
          name={{ ios: "magnifyingglass", android: "search", web: "search" }}
          size={18}
          tintColor={colors.brand.orange}
        />
        <TextInput
          ref={searchRef}
          autoCapitalize="none"
          autoCorrect={false}
          editable={isSearchMode}
          placeholder="Search visitor by name, phone, flat ..."
          placeholderTextColor={colors.guard.textMuted}
          selectionColor={colors.brand.orange}
          style={[styles.searchInput, Platform.OS === "web" ? styles.searchInputWeb : null]}
          value={searchValue}
          onChangeText={onSearchChange}
        />
        <Pressable
          accessibilityLabel="Open filters"
          accessibilityRole="button"
          hitSlop={8}
          style={styles.iconButton}
          onPress={onFilterPress}
        >
          <SymbolView
            name={{ ios: "line.3.horizontal.decrease", android: "filter_list", web: "filter_list" }}
            size={18}
            tintColor={colors.guard.text}
          />
          {hasExplicitFilters || activeFilterCount > 0 ? <View style={styles.filterDot} /> : null}
        </Pressable>
      </Pressable>

      {isSearchMode ? (
        <View style={styles.modeCopy}>
          <Text style={styles.searchHint}>
            {isSearchActive ? (
              <>
                Showing results for{" "}
                <Text style={styles.searchHintTerm}>&quot;{searchValue.trim()}&quot;</Text>
              </>
            ) : hasExplicitFilters ? (
              "Showing filtered visitor results."
            ) : (
              "Search across all visitor entries."
            )}
          </Text>
          {hasExplicitFilters && filterSummary ? (
            <Text style={styles.filterSummary}>{filterSummary}</Text>
          ) : null}
        </View>
      ) : hasExplicitFilters ? (
        <View style={styles.filteredPanel}>
          <View style={styles.filteredCopy}>
            <Text style={styles.filteredLabel}>Custom history query</Text>
            <Text style={styles.filteredSummary}>{filterSummary || "Filters active"}</Text>
          </View>
          <Pressable
            accessibilityRole="button"
            style={styles.clearFiltersButton}
            onPress={onClearFilters}
          >
            <Text style={styles.clearFiltersText}>Clear filters</Text>
          </Pressable>
        </View>
      ) : showPresetControls ? (
        <>
          <SegmentTabs
            compact
            options={options}
            value={segment}
            variant="underline"
            onChange={onSegmentChange}
          />

          {!hideDatePicker ? (
            <Pressable accessibilityRole="button" style={styles.dateRow} onPress={onDatePress}>
              <SymbolView
                name={{ ios: "calendar", android: "calendar_today", web: "calendar_today" }}
                size={14}
                tintColor={colors.brand.orange}
              />
              <Text style={styles.dateLabel}>{getDateRangeLabel(datePreset)}</Text>
              <SymbolView
                name={{ ios: "chevron.down", android: "expand_more", web: "expand_more" }}
                size={12}
                tintColor={colors.brand.orange}
              />
            </Pressable>
          ) : null}
        </>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  headerAction: {
    alignItems: "center",
    minHeight: 32,
    justifyContent: "center",
    paddingHorizontal: spacing.xs,
  },
  headerActionDot: {
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    height: 7,
    position: "absolute",
    right: 0,
    top: 2,
    width: 7,
  },
  headerActionText: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "700",
  },
  container: {
    gap: spacing.md,
    paddingHorizontal: layout.screenPaddingHorizontal,
  },
  dateLabel: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "600",
  },
  dateRow: {
    alignItems: "center",
    alignSelf: "flex-start",
    flexDirection: "row",
    gap: 6,
    paddingVertical: 2,
  },
  clearFiltersButton: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: "#F1E4DA",
    borderRadius: radius.lg,
    borderWidth: 1,
    justifyContent: "center",
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  clearFiltersText: {
    color: colors.brand.orange,
    fontSize: 13,
    fontWeight: "700",
  },
  filteredCopy: {
    flex: 1,
    gap: 3,
    minWidth: 0,
  },
  filteredLabel: {
    color: colors.guard.textMuted,
    fontSize: 11,
    fontWeight: "700",
    textTransform: "uppercase",
  },
  filteredPanel: {
    alignItems: "center",
    backgroundColor: colors.brand.orangeSoft,
    borderColor: "#fed7aa",
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  filteredSummary: {
    color: colors.guard.text,
    fontSize: 14,
    fontWeight: "700",
  },
  filterDot: {
    backgroundColor: colors.brand.orange,
    borderRadius: 999,
    height: 7,
    position: "absolute",
    right: 2,
    top: 2,
    width: 7,
  },
  iconButton: {
    alignItems: "center",
    height: 32,
    justifyContent: "center",
    width: 32,
  },
  searchHint: {
    color: colors.guard.textMuted,
    fontSize: 13,
  },
  filterSummary: {
    color: colors.guard.text,
    fontSize: 13,
    fontWeight: "700",
  },
  modeCopy: {
    gap: 4,
  },
  searchHintTerm: {
    color: colors.guard.text,
    fontWeight: "600",
  },
  searchInput: {
    backgroundColor: "transparent",
    borderWidth: 0,
    color: colors.guard.text,
    flex: 1,
    fontSize: 14,
    marginLeft: spacing.sm,
    minHeight: 40,
    paddingVertical: 0,
  },
  searchInputWeb: {
    outlineWidth: 0,
  },
  searchRow: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: "#F1E4DA",
    borderRadius: radius.xl,
    borderWidth: 1,
    flexDirection: "row",
    minHeight: 48,
    paddingHorizontal: spacing.lg,
    ...shadows.sm,
  },
});
