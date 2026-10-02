import { useEffect, useRef, useState } from "react";
import { Modal, Pressable, ScrollView, StyleSheet, Text, View } from "react-native";
import { SymbolView } from "expo-symbols";
import { SafeAreaView } from "react-native-safe-area-context";

import { FilterChip, FilterChipCell, FilterChipGrid } from "@/components/ui/filter-chip";
import {
  DATE_RANGE_OPTIONS,
  getDateRangeLabel,
  type DateRangePreset,
} from "@/features/visitors/visitor-date-ranges";
import { titleize, visitorPurposes } from "@/features/visitors/visitor-utils";
import type {
  ModelsVisitorPurpose,
  ModelsVisitorStatus,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { theme } from "@/theme";

const G = theme.guard;
const accent = colors.brand.orange;
const accentSoft = colors.brand.orangeSoft;

type SheetFilters = {
  datePreset: DateRangePreset;
  purpose?: ModelsVisitorPurpose;
  status?: ModelsVisitorStatus;
};

type LogsFilterSheetProps = {
  datePreset: DateRangePreset;
  onApply: (filters: SheetFilters) => void;
  onClear: () => void;
  onClose: () => void;
  purpose?: ModelsVisitorPurpose;
  status?: ModelsVisitorStatus;
  visible: boolean;
};

const DEFAULT_FILTERS: SheetFilters = {
  datePreset: "today",
  purpose: undefined,
  status: undefined,
};

const STATUS_OPTIONS: { label: string; value?: ModelsVisitorStatus }[] = [
  { label: "All", value: undefined },
  { label: "Pending", value: "waiting_approval" },
  { label: "Approved", value: "approved" },
  { label: "Checked In", value: "checked_in" },
  { label: "Checked Out", value: "checked_out" },
  { label: "Rejected", value: "rejected" },
];

const PURPOSE_OPTIONS: { label: string; value?: ModelsVisitorPurpose }[] = [
  { label: "All", value: undefined },
  ...visitorPurposes.map((purpose) => ({
    label: titleize(purpose),
    value: purpose,
  })),
];

function SectionTitle({ title }: { title: string }) {
  return (
    <Text style={styles.sectionTitle}>{title}</Text>
  );
}

export function LogsFilterSheet({
  datePreset,
  onApply,
  onClear,
  onClose,
  purpose,
  status,
  visible,
}: LogsFilterSheetProps) {
  const [draft, setDraft] = useState<SheetFilters>({
    datePreset,
    purpose,
    status,
  });
  const [dateMenuOpen, setDateMenuOpen] = useState(false);
  const [dateSectionY, setDateSectionY] = useState(0);
  const scrollRef = useRef<ScrollView>(null);

  useEffect(() => {
    if (visible) {
      /* eslint-disable react-hooks/set-state-in-effect -- Re-seed filter draft values each time the sheet opens. */
      setDraft({ datePreset, purpose, status });
      setDateMenuOpen(false);
      /* eslint-enable react-hooks/set-state-in-effect */
    }
  }, [datePreset, purpose, status, visible]);

  useEffect(() => {
    if (!dateMenuOpen) {
      return;
    }

    const timer = setTimeout(() => {
      scrollRef.current?.scrollTo({
        animated: true,
        y: Math.max(dateSectionY - 12, 0),
      });
    }, 80);

    return () => clearTimeout(timer);
  }, [dateMenuOpen, dateSectionY]);

  const handleReset = () => {
    setDraft(DEFAULT_FILTERS);
    setDateMenuOpen(false);
    onClear();
  };

  return (
    <Modal animationType="slide" transparent visible={visible} onRequestClose={onClose}>
      <Pressable style={styles.backdrop} onPress={onClose}>
        <Pressable onPress={(event) => event.stopPropagation()}>
          <SafeAreaView edges={["bottom"]} style={styles.sheet}>
            <View style={styles.sheetBody}>
              <View style={styles.headerRow}>
                <View style={styles.headerCopy}>
                  <Text style={styles.headerTitle}>Refine visitor logs</Text>
                  <Text style={styles.headerSubtitle}>
                    Filter by status, purpose, and date range
                  </Text>
                </View>
                <Pressable accessibilityRole="button" hitSlop={8} onPress={onClose}>
                  <Text style={styles.closeText}>Close</Text>
                </Pressable>
              </View>

              <ScrollView
                ref={scrollRef}
                contentContainerStyle={styles.scrollContent}
                showsVerticalScrollIndicator={false}
                style={styles.scroll}
              >
                <View style={styles.section}>
                  <SectionTitle title="Status" />
                  <FilterChipGrid>
                    {STATUS_OPTIONS.map((option) => (
                      <FilterChipCell key={option.label}>
                        <FilterChip
                          label={option.label}
                          selected={draft.status === option.value}
                          onPress={() =>
                            setDraft((current) => ({ ...current, status: option.value }))
                          }
                        />
                      </FilterChipCell>
                    ))}
                  </FilterChipGrid>
                </View>

                <View style={styles.section}>
                  <SectionTitle title="Purpose" />
                  <FilterChipGrid>
                    {PURPOSE_OPTIONS.map((option) => (
                      <FilterChipCell key={option.label}>
                        <FilterChip
                          label={option.label}
                          selected={draft.purpose === option.value}
                          onPress={() =>
                            setDraft((current) => ({ ...current, purpose: option.value }))
                          }
                        />
                      </FilterChipCell>
                    ))}
                  </FilterChipGrid>
                </View>

                <View
                  style={styles.section}
                  onLayout={(event) => setDateSectionY(event.nativeEvent.layout.y)}
                >
                  <SectionTitle title="Date Range" />
                  <Pressable
                    accessibilityRole="button"
                    style={styles.dateSelector}
                    onPress={() => setDateMenuOpen((open) => !open)}
                  >
                    <SymbolView
                      name={{ ios: "calendar", android: "calendar_today", web: "calendar_today" }}
                      size={16}
                      tintColor={accent}
                    />
                    <Text style={styles.dateSelectorText}>{getDateRangeLabel(draft.datePreset)}</Text>
                    <SymbolView
                      name={{
                        ios: dateMenuOpen ? "chevron.up" : "chevron.down",
                        android: dateMenuOpen ? "expand_less" : "expand_more",
                        web: dateMenuOpen ? "expand_less" : "expand_more",
                      }}
                      size={16}
                      tintColor={G.textMuted}
                    />
                  </Pressable>

                  {dateMenuOpen ? (
                    <View style={styles.dateMenu}>
                      {DATE_RANGE_OPTIONS.map((option) => {
                        const active = draft.datePreset === option.value;

                        return (
                          <Pressable
                            key={option.value}
                            accessibilityRole="button"
                            style={[styles.dateOption, active ? styles.dateOptionActive : null]}
                            onPress={() => {
                              setDraft((current) => ({ ...current, datePreset: option.value }));
                              setDateMenuOpen(false);
                            }}
                          >
                            <Text
                              style={[
                                styles.dateOptionText,
                                active ? styles.dateOptionTextActive : null,
                              ]}
                            >
                              {active ? `✓ ${option.label}` : option.label}
                            </Text>
                          </Pressable>
                        );
                      })}
                    </View>
                  ) : null}
                </View>
              </ScrollView>

              <View style={styles.footer}>
                <Pressable
                  accessibilityRole="button"
                  style={({ pressed }) => [styles.footerButton, styles.resetButton, pressed && styles.buttonPressed]}
                  onPress={handleReset}
                >
                  <Text style={styles.resetButtonText}>Reset</Text>
                </Pressable>
                <Pressable
                  accessibilityRole="button"
                  style={({ pressed }) => [styles.footerButton, styles.applyButton, pressed && styles.buttonPressed]}
                  onPress={() => {
                    onApply(draft);
                    onClose();
                  }}
                >
                  <Text style={styles.applyButtonText}>Apply</Text>
                </Pressable>
              </View>
            </View>
          </SafeAreaView>
        </Pressable>
      </Pressable>
    </Modal>
  );
}

const styles = StyleSheet.create({
  applyButton: {
    backgroundColor: accent,
    borderColor: accent,
  },
  applyButtonText: {
    color: "#ffffff",
    fontSize: 15,
    fontWeight: "700",
  },
  backdrop: {
    backgroundColor: "rgba(15, 23, 42, 0.35)",
    flex: 1,
    justifyContent: "flex-end",
  },
  buttonPressed: {
    opacity: 0.88,
  },
  closeText: {
    color: accent,
    fontSize: 15,
    fontWeight: "600",
  },
  dateMenu: {
    backgroundColor: theme.surface.card,
    borderColor: theme.border.default,
    borderRadius: 14,
    borderWidth: 1,
    marginTop: 10,
    overflow: "hidden",
  },
  dateOption: {
    borderBottomColor: theme.border.default,
    borderBottomWidth: 1,
    paddingHorizontal: 14,
    paddingVertical: 13,
  },
  dateOptionActive: {
    backgroundColor: accentSoft,
  },
  dateOptionText: {
    color: G.textMuted,
    fontSize: 14,
    fontWeight: "500",
  },
  dateOptionTextActive: {
    color: accent,
    fontWeight: "700",
  },
  dateSelector: {
    alignItems: "center",
    backgroundColor: theme.surface.card,
    borderColor: theme.border.default,
    borderRadius: 12,
    borderWidth: 1.5,
    flexDirection: "row",
    gap: 10,
    minHeight: 48,
    paddingHorizontal: 14,
  },
  dateSelectorText: {
    color: G.text,
    flex: 1,
    fontSize: 15,
    fontWeight: "600",
  },
  footer: {
    borderTopColor: theme.border.default,
    borderTopWidth: 1,
    flexDirection: "row",
    gap: 12,
    paddingTop: 16,
  },
  footerButton: {
    alignItems: "center",
    borderRadius: 14,
    borderWidth: 1.5,
    flex: 1,
    justifyContent: "center",
    minHeight: 50,
  },
  headerCopy: {
    flex: 1,
    gap: 4,
    paddingRight: 12,
  },
  headerRow: {
    alignItems: "flex-start",
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 16,
  },
  headerSubtitle: {
    color: G.textMuted,
    fontSize: 13,
    lineHeight: 18,
  },
  headerTitle: {
    color: G.text,
    fontSize: 18,
    fontWeight: "700",
  },
  resetButton: {
    backgroundColor: theme.surface.muted,
    borderColor: theme.border.default,
  },
  resetButtonText: {
    color: G.text,
    fontSize: 15,
    fontWeight: "700",
  },
  scroll: {
    maxHeight: 460,
  },
  scrollContent: {
    gap: 22,
    paddingBottom: 132,
  },
  section: {
    gap: 12,
  },
  sectionTitle: {
    color: G.text,
    fontSize: 15,
    fontWeight: "700",
  },
  sheet: {
    backgroundColor: theme.surface.card,
  },
  sheetBody: {
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    paddingBottom: 20,
    paddingHorizontal: 20,
    paddingTop: 18,
  },
});
