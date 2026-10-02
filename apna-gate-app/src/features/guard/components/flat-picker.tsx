import { useEffect, useState } from "react";
import {
  ActivityIndicator,
  FlatList,
  Modal,
  Platform,
  Pressable,
  StyleSheet,
  Text,
  TextInput,
  View,
  type TextStyle,
} from "react-native";
import { SymbolView } from "expo-symbols";
import {
  SafeAreaView,
  useSafeAreaInsets,
} from "react-native-safe-area-context";

import { AppStatusBar } from "@/components/layout/app-status-bar";
import { Row, Stack } from "@/components/layout";
import { dashboardActionToneStyles } from "@/components/dashboard";
import {
  flatFromResponse,
  type SelectedFlat,
} from "@/features/guard/hooks/use-guard-manual-entry";
import {
  useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextQuery,
  useGetV1SocietiesBySocietyIdFlatsQuery,
} from "@/lib/api/generated-api";
import { colors } from "@/theme/colors";
import { layout } from "@/theme/layout";
import { radius } from "@/theme/radius";
import { shadows } from "@/theme/shadows";
import { spacing } from "@/theme/spacing";
import { KeyboardAvoidingView } from "react-native";

const SEARCH_DEBOUNCE_MS = 400;
const SEARCH_LIMIT = 8;

const androidTouchableFocusProps =
  Platform.OS === "android" ? { focusable: false as const } : {};

const webNoOutline: TextStyle =
  Platform.OS === "web"
    ? ({ outlineStyle: "none" } as unknown as TextStyle)
    : {};

function FlatSearchField({
  label = "Search flat",
  value,
  onChangeText,
  placeholder,
  autoFocus,
}: {
  label?: string;
  value: string;
  onChangeText: (text: string) => void;
  placeholder?: string;
  autoFocus?: boolean;
}) {
  return (
    <Stack gap={6}>
      <Text style={styles.fieldLabel}>{label}</Text>
      <View style={styles.searchFieldWrapper}>
        <View style={styles.fieldInputIcon}>
          <SymbolView
            name={{ ios: "magnifyingglass", android: "search", web: "search" }}
            size={20}
            tintColor={colors.guard.textMuted}
          />
        </View>
        <TextInput
          autoCapitalize="none"
          autoCorrect={false}
          autoFocus={autoFocus}
          blurOnSubmit={false}
          cursorColor={colors.brand.orange}
          importantForAutofill="no"
          placeholder={placeholder}
          placeholderTextColor={colors.guard.textMuted}
          selectionColor={colors.accent.selection}
          underlineColorAndroid="transparent"
          value={value}
          onChangeText={onChangeText}
          style={[styles.searchInput, webNoOutline]}
        />
      </View>
    </Stack>
  );
}

function FlatSearchModal({
  visible,
  societyId,
  onClose,
  onSelect,
}: {
  visible: boolean;
  societyId: number;
  onClose: () => void;
  onSelect: (flat: SelectedFlat) => void;
}) {
  const insets = useSafeAreaInsets();
  const [query, setQuery] = useState("");
  const [debounced, setDebounced] = useState("");

  useEffect(() => {
    const t = setTimeout(() => setDebounced(query.trim()), SEARCH_DEBOUNCE_MS);
    return () => clearTimeout(t);
  }, [query]);

  useEffect(() => {
    if (!visible) {
      /* eslint-disable react-hooks/set-state-in-effect -- Reset search draft when the modal closes. */
      setQuery("");
      setDebounced("");
      /* eslint-enable react-hooks/set-state-in-effect */
    }
  }, [visible]);

  const canSearch = debounced.length > 0;
  const { data, isFetching, isLoading } =
    useGetV1SocietiesBySocietyIdFlatsQuery(
      {
        societyId,
        search: debounced,
        status: "occupied",
        isActive: true,
        limit: SEARCH_LIMIT,
      },
      { skip: !visible || !canSearch },
    );

  const flats = data?.data?.flats?.items ?? [];
  const total = data?.data?.flats?.total ?? 0;

  if (!visible) {
    return null;
  }

  return (
    <Modal animationType="slide" visible onRequestClose={onClose}>
      <SafeAreaView
        edges={["top", "left", "right"]}
        style={styles.modalScreen}
      >
        <AppStatusBar />
        <Row align="center" gap="sm" style={styles.modalHeader}>
          <Pressable
            style={styles.modalBackButton}
            onPress={onClose}
            {...androidTouchableFocusProps}
          >
            <SymbolView
              name={{
                ios: "chevron.left",
                android: "arrow_back",
                web: "arrow_back",
              }}
              size={20}
              tintColor={colors.guard.text}
            />
          </Pressable>
          <View style={styles.modalHeaderText}>
            <Text style={styles.modalTitle}>Find flat</Text>
            <Text style={styles.modalSubtitle}>Type to see matching flats</Text>
          </View>
        </Row>

        <View style={styles.modalSearchSection}>
          <FlatSearchField
            autoFocus={visible}
            placeholder="Flat no., resident, or wing"
            value={query}
            onChangeText={setQuery}
          />
        </View>

        {!canSearch ? (
          <View style={styles.modalEmptyState}>
            <Text style={styles.modalEmptyStateText}>
              Type to search flats{"\n"}e.g. G-02, left-wing
            </Text>
          </View>
        ) : isLoading ? (
          <ActivityIndicator
            color={colors.guard.teal}
            style={styles.modalLoading}
          />
        ) : flats.length === 0 ? (
          <Text style={styles.modalNoResults}>No flats found</Text>
        ) : (
          <KeyboardAvoidingView
            style={{ flex: 1 }}
            behavior={Platform.OS === "ios" ? "padding" : "height"}
          >
            <FlatList
              contentContainerStyle={[
                styles.modalListContent,
                {
                  paddingBottom:
                    spacing["3xl"] + Math.max(insets.bottom, 0),
                },
              ]}
              data={flats}
              keyExtractor={(item) => String(item.id)}
              keyboardShouldPersistTaps="always"
              ListFooterComponent={
                total > flats.length ? (
                  <Text style={styles.modalListFooter}>
                    {flats.length} of {total} — type more to refine
                  </Text>
                ) : isFetching ? (
                  <ActivityIndicator
                    color={colors.guard.teal}
                    style={styles.modalListLoading}
                  />
                ) : null
              }
              renderItem={({ item }) => (
                <Pressable
                  style={styles.flatListItem}
                  onPress={() => {
                    const flat = flatFromResponse(item);
                    if (flat) {
                      onSelect(flat);
                      onClose();
                    }
                  }}
                  {...androidTouchableFocusProps}
                >
                  <Text style={styles.flatListItemNumber}>
                    {item.flat_number ?? `#${item.id}`}
                  </Text>
                  <Text style={styles.flatListItemMeta}>
                    {[
                      item.block ? `Wing ${item.block}` : null,
                      item.floor ? `Floor ${item.floor}` : null,
                    ]
                      .filter(Boolean)
                      .join(" · ")}
                  </Text>
                </Pressable>
              )}
            />
          </KeyboardAvoidingView>
        )}
      </SafeAreaView>
    </Modal>
  );
}

function ResidentPreviewCard({
  flat,
  resident,
  loading,
  onPress,
}: {
  flat: SelectedFlat;
  resident?: string | null;
  loading?: boolean;
  onPress: () => void;
}) {
  const wingLabel = flat.block
    ? `Wing ${flat.block}`
    : flat.floor
      ? `Tower ${flat.floor}`
      : null;
  const residentLine = loading
    ? "Loading resident..."
    : [resident ?? "No primary resident", wingLabel]
        .filter(Boolean)
        .join(" · ");

  return (
    <Pressable
      style={styles.flatSearchCard}
      onPress={onPress}
      {...androidTouchableFocusProps}
    >
      <View
        style={[
          styles.flatSearchIcon,
          { backgroundColor: dashboardActionToneStyles.blue.backgroundColor },
        ]}
      >
        <SymbolView
          name={{ ios: "house.fill", android: "home", web: "home" }}
          size={24}
          tintColor={dashboardActionToneStyles.blue.iconColor}
        />
      </View>
      <Stack gap={2} style={styles.flatSearchCopy}>
        <Text style={styles.flatSearchTitle}>
          {flat.flat_number ?? `Flat ${flat.id}`}
        </Text>
        <Text numberOfLines={1} style={styles.flatSearchSubtitle}>
          {residentLine}
        </Text>
      </Stack>
      <Row align="center" gap="xs" style={styles.changeBadge}>
        <Text style={styles.changeBadgeText}>Change</Text>
        <SymbolView
          name={{
            ios: "chevron.right",
            android: "chevron_right",
            web: "chevron_right",
          }}
          size={10}
          tintColor={colors.brand.orange}
        />
      </Row>
    </Pressable>
  );
}

export type FlatPickerProps = {
  societyId: number;
  selected: SelectedFlat | null;
  error?: string;
  label?: string;
  onSelect: (flat: SelectedFlat) => void;
};

export function FlatPicker({
  societyId,
  selected,
  error,
  label,
  onSelect,
}: FlatPickerProps) {
  const [open, setOpen] = useState(false);
  const { data, isFetching } =
    useGetV1SocietiesBySocietyIdFlatsAndFlatIdVisitorContextQuery(
      { societyId, flatId: selected?.id ?? 0 },
      { skip: !selected?.id },
    );
  const resident = data?.data?.context?.primary_resident?.full_name;

  return (
    <Stack gap="md">
      {label ? <Text style={styles.fieldLabel}>{label}</Text> : null}
      {selected ? (
        <ResidentPreviewCard
          flat={selected}
          loading={isFetching}
          resident={resident}
          onPress={() => setOpen(true)}
        />
      ) : (
        <Pressable
          style={[styles.flatSearchCard, error && styles.flatSearchCardError]}
          onPress={() => setOpen(true)}
          {...androidTouchableFocusProps}
        >
          <View style={[styles.flatSearchIcon, styles.flatSearchIconOrange]}>
            <SymbolView
              name={{
                ios: "magnifyingglass",
                android: "search",
                web: "search",
              }}
              size={24}
              tintColor={colors.brand.orange}
            />
          </View>
          <Stack gap={2} style={styles.flatSearchCopy}>
            <Text style={styles.flatSearchTitle}>Search flat</Text>
            <Text style={styles.flatSearchSubtitle}>
              Search by flat no., resident name or wing
            </Text>
          </Stack>
          <SymbolView
            name={{
              ios: "chevron.right",
              android: "chevron_right",
              web: "chevron_right",
            }}
            size={16}
            tintColor={colors.guard.textMuted}
          />
        </Pressable>
      )}

      {error ? <Text style={styles.flatPickerError}>{error}</Text> : null}
      {open ? (
        <KeyboardAvoidingView
          style={{ flex: 1 }}
          behavior={Platform.OS === "ios" ? "padding" : "height"}
        >
          <FlatSearchModal
            societyId={societyId}
            visible
            onClose={() => setOpen(false)}
            onSelect={onSelect}
          />
        </KeyboardAvoidingView>
      ) : null}
    </Stack>
  );
}

const styles = StyleSheet.create({
  changeBadge: {
    backgroundColor: colors.brand.orangeSoft,
    borderRadius: 999,
    paddingHorizontal: spacing.md,
    paddingVertical: 6,
  },
  changeBadgeText: {
    color: colors.brand.orange,
    fontSize: 12,
    fontWeight: "600",
  },
  fieldInputIcon: {
    alignItems: "center",
    justifyContent: "center",
    width: 24,
  },
  fieldLabel: {
    color: colors.guard.textMuted,
    fontSize: 13,
    fontWeight: "500",
  },
  flatListItem: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingVertical: spacing.lg,
    ...shadows.card,
  },
  flatListItemMeta: {
    color: colors.text.muted,
    fontSize: 14,
  },
  flatListItemNumber: {
    color: colors.text.primary,
    flex: 1,
    fontSize: 17,
    fontWeight: "600",
  },
  flatPickerError: {
    color: colors.status.error,
    fontSize: 14,
    marginTop: spacing.sm,
  },
  flatSearchCard: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.md,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.lg,
    ...shadows.sm,
  },
  flatSearchCardError: {
    borderColor: "#fca5a5",
  },
  flatSearchCopy: {
    flex: 1,
    minWidth: 0,
  },
  flatSearchIcon: {
    alignItems: "center",
    borderRadius: radius.lg,
    height: 48,
    justifyContent: "center",
    width: 48,
  },
  flatSearchIconOrange: {
    backgroundColor: colors.brand.orangeSoft,
  },
  flatSearchSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 13,
    fontWeight: "500",
  },
  flatSearchTitle: {
    color: colors.guard.text,
    fontSize: 16,
    fontWeight: "700",
  },
  modalBackButton: {
    alignItems: "center",
    backgroundColor: colors.guard.screenBg,
    borderRadius: 999,
    height: 40,
    justifyContent: "center",
    width: 40,
  },
  modalEmptyState: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
    paddingHorizontal: layout.screenPaddingHorizontal,
  },
  modalEmptyStateText: {
    color: colors.guard.textMuted,
    fontSize: 15,
    lineHeight: 22,
    textAlign: "center",
  },
  modalHeader: {
    borderBottomColor: colors.guard.border,
    borderBottomWidth: StyleSheet.hairlineWidth,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingVertical: spacing.md,
  },
  modalHeaderText: {
    flex: 1,
  },
  modalListContent: {
    gap: spacing.sm,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.md,
  },
  modalListFooter: {
    color: colors.guard.textMuted,
    fontSize: 13,
    paddingVertical: spacing.md,
    textAlign: "center",
  },
  modalListLoading: {
    paddingVertical: spacing.lg,
  },
  modalLoading: {
    marginTop: spacing["3xl"],
  },
  modalNoResults: {
    color: colors.guard.textMuted,
    fontSize: 15,
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing["2xl"],
    textAlign: "center",
  },
  modalScreen: {
    backgroundColor: colors.guard.screenBg,
    flex: 1,
  },
  modalSearchSection: {
    paddingHorizontal: layout.screenPaddingHorizontal,
    paddingTop: spacing.lg,
  },
  modalSubtitle: {
    color: colors.guard.textMuted,
    fontSize: 13,
    marginTop: 2,
  },
  modalTitle: {
    color: colors.guard.text,
    fontSize: 18,
    fontWeight: "700",
  },
  searchFieldWrapper: {
    alignItems: "center",
    backgroundColor: colors.surface.card,
    borderColor: colors.guard.border,
    borderRadius: radius.lg,
    borderWidth: 1,
    flexDirection: "row",
    gap: spacing.sm,
    minHeight: layout.inputHeight,
    paddingHorizontal: spacing.md,
  },
  searchInput: {
    color: colors.guard.text,
    flex: 1,
    fontSize: 16,
    minHeight: 24,
    padding: 0,
  },
});
