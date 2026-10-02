import { Modal, Pressable, ScrollView, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { Button, StatusPill } from "@/components/ui";
import { titleize } from "@/features/guard/guard-utils";
import type { ModelsSocietyMemberResponse } from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { theme } from "@/theme";

type SocietySwitchSheetProps = {
  memberships: ModelsSocietyMemberResponse[];
  onClose: () => void;
  onSelect: (societyId: number) => void;
  selectedSocietyId?: number | null;
  visible: boolean;
};

export function SocietySwitchSheet({
  memberships,
  onClose,
  onSelect,
  selectedSocietyId,
  visible,
}: SocietySwitchSheetProps) {
  return (
    <Modal animationType="slide" transparent visible={visible} onRequestClose={onClose}>
      <Pressable style={styles.backdrop} onPress={onClose}>
        <Pressable onPress={(event) => event.stopPropagation()}>
          <SafeAreaView edges={["bottom"]} style={styles.sheet}>
            <View style={styles.sheetBody}>
              <View style={styles.headerRow}>
                <View style={styles.headerCopy}>
                  <Text style={styles.headerTitle}>Switch society</Text>
                  <Text style={styles.headerSubtitle}>
                    Choose which society workspace to use
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
                {memberships.length === 0 ? (
                  <View style={styles.emptyCard}>
                    <Text style={styles.emptyTitle}>No active staff access</Text>
                    <Text style={styles.emptyBody}>
                      Ask your society admin to link your guard account.
                    </Text>
                  </View>
                ) : (
                  memberships.map((membership) => {
                    const isSelected = membership.society_id === selectedSocietyId;

                    return (
                      <View
                        key={`society-switch-${membership.id ?? membership.society_id}`}
                        style={[styles.membershipCard, isSelected ? styles.membershipCardActive : null]}
                      >
                        <View style={styles.membershipHeader}>
                          <View style={styles.membershipCopy}>
                            <Text style={styles.membershipTitle}>
                              Society #{membership.society_id}
                            </Text>
                            <Text style={styles.membershipSubtitle}>
                              {titleize(membership.role ?? "staff")} access
                            </Text>
                          </View>
                          <StatusPill status={membership.status} />
                        </View>
                        <Button
                          compact
                          disabled={isSelected}
                          title={isSelected ? "Current society" : "Switch society"}
                          variant={isSelected ? "ghost" : "primary"}
                          onPress={() => {
                            if (membership.society_id) {
                              onSelect(membership.society_id);
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
  membershipCard: {
    backgroundColor: theme.surface.card,
    borderColor: theme.border.default,
    borderRadius: 16,
    borderWidth: 1,
    gap: 12,
    marginBottom: 12,
    padding: 16,
  },
  membershipCardActive: {
    backgroundColor: colors.brand.orangeSoft,
    borderColor: colors.brand.orange,
  },
  membershipCopy: {
    flex: 1,
    gap: 2,
  },
  membershipHeader: {
    alignItems: "flex-start",
    flexDirection: "row",
    gap: 12,
    justifyContent: "space-between",
  },
  membershipSubtitle: {
    color: theme.text.secondary,
    fontSize: 14,
    textTransform: "capitalize",
  },
  membershipTitle: {
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
