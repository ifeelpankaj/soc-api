import { useState } from "react";
import { Modal, Pressable, ScrollView, StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { SymbolView } from "expo-symbols";
import { Row } from "@/components/layout";
import { AppText } from "./app-text";
import { Button } from "./button";
import { FilterChip } from "./filter-chip";
import { colors } from "@/theme/colors";

const monthNames = Array.from({ length: 12 }, (_, index) =>
  new Date(2026, index, 1).toLocaleDateString("en-IN", { month: "long" }),
);

/** Month-only selection using the app's date-sheet and chip conventions. */
export function MonthPicker({ value, onChange }: {
  value: string;
  onChange: (value: string) => void;
}) {
  const [visible, setVisible] = useState(false);
  const [year, setYear] = useState(new Date().getFullYear());
  const selectedYear = Number(value.slice(0, 4));
  const selectedMonth = Number(value.slice(5, 7));
  return (
    <>
      <Row gap="sm">
        <Pressable accessibilityRole="button" accessibilityLabel="Select billing month"
          onPress={() => { setYear(selectedYear || new Date().getFullYear()); setVisible(true); }}
          style={styles.field}>
          <SymbolView name={{ ios: "calendar", android: "calendar_month", web: "calendar_month" }} size={18} tintColor={colors.text.muted} />
          <AppText style={{ flex: 1 }}>{value ? `${monthNames[selectedMonth - 1]} ${selectedYear}` : "Billing month"}</AppText>
          <AppText color="muted">⌄</AppText>
        </Pressable>
        {value ? <Button fullWidth={false} compact title="Clear" variant="ghost" onPress={() => onChange("")} /> : null}
      </Row>
      <Modal animationType="slide" transparent visible={visible} onRequestClose={() => setVisible(false)}>
        <Pressable style={styles.backdrop} onPress={() => setVisible(false)}>
          <Pressable style={styles.sheet} onPress={(event) => event.stopPropagation()}>
            <SafeAreaView edges={["bottom"]}>
              <ScrollView contentContainerStyle={styles.body}>
                <Row justify="space-between">
                  <AppText variant="title">Billing month</AppText>
                  <Button fullWidth={false} compact title="Close" variant="ghost" onPress={() => setVisible(false)} />
                </Row>
                <Row justify="space-between">
                  <Button fullWidth={false} compact title="‹" accessibilityLabel="Previous year" variant="secondary" disabled={year <= 1} onPress={() => setYear(year - 1)} />
                  <AppText style={{ fontWeight: "700" }}>{year}</AppText>
                  <Button fullWidth={false} compact title="›" accessibilityLabel="Next year" variant="secondary" disabled={year >= 9999} onPress={() => setYear(year + 1)} />
                </Row>
                <View style={styles.months}>
                  {monthNames.map((name, index) => (
                    <View key={name} style={styles.cell}>
                      <FilterChip label={name} selected={selectedYear === year && selectedMonth === index + 1}
                        onPress={() => { onChange(`${String(year).padStart(4, "0")}-${String(index + 1).padStart(2, "0")}`); setVisible(false); }} />
                    </View>
                  ))}
                </View>
              </ScrollView>
            </SafeAreaView>
          </Pressable>
        </Pressable>
      </Modal>
    </>
  );
}

const styles = StyleSheet.create({
  field: { flex: 1, minWidth: 0, minHeight: 48, flexDirection: "row", alignItems: "center", gap: 8, padding: 12, borderWidth: 1, borderColor: colors.border.default, borderRadius: 14, backgroundColor: colors.surface.card },
  backdrop: { flex: 1, justifyContent: "flex-end", backgroundColor: "rgba(15, 23, 42, 0.35)" },
  sheet: { backgroundColor: colors.surface.screen, borderTopLeftRadius: 24, borderTopRightRadius: 24, maxHeight: "85%" },
  body: { padding: 16, gap: 16 },
  months: { flexDirection: "row", flexWrap: "wrap", gap: 8 },
  cell: { width: "48%" },
});
