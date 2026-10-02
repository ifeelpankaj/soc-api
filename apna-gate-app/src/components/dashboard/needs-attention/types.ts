import type { SymbolViewProps } from "expo-symbols";

export type AttentionItem = {
  id: string;
  icon: SymbolViewProps["name"];
  iconTone?: "orange" | "blue" | "purple" | "neutral" | "teal";
  title: string;
  subtitle: string;
  badge?: string;
  actionLabel?: string;
  onPress: () => void;
};
