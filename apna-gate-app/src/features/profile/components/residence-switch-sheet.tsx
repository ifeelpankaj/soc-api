import { Modal, Pressable, ScrollView, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { Button, StatusPill } from "@/components/ui";
import { titleize } from "@/features/guard/guard-utils";
import type { ModelsFlatResidentResponse } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { theme } from "@/theme";

type ResidenceSwitchSheetProps = {
  onClose: () => void;
  onSelect: (flatId: number) => void;
  residences: ModelsFlatResidentResponse[];
  selectedFlatId?: number | null;
  visible: boolean;
};

export function ResidenceSwitchSheet({
  onClose,
  onSelect,
  residences,
  selectedFlatId,
  visible,
}: ResidenceSwitchSheetProps) {
  return (
    <Modal animationType="slide" transparent visible={visible} onRequestClose={onClose}>
      <Pressable style={styles.backdrop} onPress={onClose}>
        <Pressable onPress={(event) => event.stopPropagation()}>
          <SafeAreaView edges={["bottom"]} style={styles.sheet}>
            <View style={styles.sheetBody}>
              <View style={styles.headerRow}>
                <View style={styles.headerCopy}>
                  <Text style={styles.headerTitle}>Switch flat</Text>
                  <Text style={styles.headerSubtitle}>
                    Choose which flat residence to use
                  </Text>
                </View>
                <Pressable accessibilityRole="button" hitSlop={8} onPress={onClose}>
                  <Text style={styles.closeText}>Close</Text>
                </Pressable>
              </View>

              <ScrollView
                contentContainerStyle={styles.scrollContent}
                showsVerticalScrollIndicator={false}
              >
                {residences.length === 0 ? (
                  <View style={styles.emptyCard}>
                    <Text style={styles.emptyTitle}>No active residences</Text>
                    <Text style={styles.emptyBody}>
                      Flat claims are approved by your society admin.
                    </Text>
                  </View>
                ) : (
                  residences.map((residence) => {
                    const isSelected = residence.flat_id === selectedFlatId;
                    const flatLabel = [
                      residence.block ? `Block ${residence.block}` : null,
                      residence.flat_number ? `Flat ${residence.flat_number}` : null,
                    ]
                      .filter(Boolean)
                      .join(" · ");

                    return (
                      <View
                        key={`residence-switch-${residence.id ?? residence.flat_id}`}
                        style={[
                          styles.residenceCard,
                          isSelected ? styles.residenceCardActive : null,
                        ]}
                      >
                        <View style={styles.residenceHeader}>
                          <View style={styles.residenceCopy}>
                            <Text style={styles.residenceTitle}>
                              {residence.society_name ?? "Society"}
                            </Text>
                            <Text style={styles.residenceSubtitle}>
                              {flatLabel || "Flat"} ·{" "}
                              {residence.is_primary ? "Primary" : titleize(residence.role ?? "resident")}
                            </Text>
                          </View>
                          <StatusPill status={residence.status} />
                        </View>
                        <Button
                          compact
                          disabled={isSelected}
                          title={isSelected ? "Current flat" : "Switch flat"}
                          variant={isSelected ? "ghost" : "primary"}
                          onPress={() => {
                            if (residence.flat_id) {
                              onSelect(residence.flat_id);
                              onClose();
                            }
                          }}
                        />
                      </View>
                    );
                  })
                )}
              </ScrollView>
            </View>
          </SafeAreaView>
        </Pressable>
      </Pressable>
    </Modal>
  );
}

const styles = StyleSheet.create({
  backdrop: {
    backgroundColor: "rgba(15, 23, 42, 0.42)",
    flex: 1,
    justifyContent: "flex-end",
  },
  closeText: {
    color: colors.brand.orange,
    fontSize: 14,
    fontWeight: "700",
  },
  emptyBody: {
    color: theme.text.secondary,
    fontSize: 14,
    lineHeight: 20,
    marginTop: 4,
  },
  emptyCard: {
    backgroundColor: theme.surface.card,
    borderColor: theme.border.default,
    borderRadius: 16,
    borderWidth: 1,
    padding: 16,
  },
  emptyTitle: {
    color: theme.text.primary,
    fontSize: 16,
    fontWeight: "700",
  },
  headerCopy: {
    flex: 1,
    gap: 4,
  },
  headerRow: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: 12,
    justifyContent: "space-between",
    marginBottom: 16,
  },
  headerSubtitle: {
    color: theme.text.secondary,
    fontSize: 14,
    lineHeight: 20,
  },
  headerTitle: {
    color: theme.text.primary,
    fontSize: 20,
    fontWeight: "700",
  },
  residenceCard: {
    backgroundColor: theme.surface.card,
    borderColor: theme.border.default,
    borderRadius: 16,
    borderWidth: 1,
    gap: 12,
    marginBottom: 12,
    padding: 16,
  },
  residenceCardActive: {
    backgroundColor: colors.brand.orangeSoft,
    borderColor: colors.brand.orange,
  },
  residenceCopy: {
    flex: 1,
    gap: 2,
  },
  residenceHeader: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: 12,
    justifyContent: "space-between",
  },
  residenceSubtitle: {
    color: theme.text.secondary,
    fontSize: 14,
    textTransform: "capitalize",
  },
  residenceTitle: {
    color: theme.text.primary,
    fontSize: 17,
    fontWeight: "700",
  },
  scrollContent: {
    paddingBottom: 8,
  },
  sheet: {
    backgroundColor: theme.surface.screen,
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    maxHeight: "82%",
  },
  sheetBody: {
    paddingHorizontal: 20,
    paddingTop: 20,
  },
});
