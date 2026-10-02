import type { SymbolViewProps } from "expo-symbols";

export const invitePurposePresentation = {
  guest: {
    description: "Friends, family and personal visits",
    icon: { ios: "person.2", android: "group", web: "group" },
  },
  delivery: {
    description: "Food, parcels and doorstep deliveries",
    icon: { ios: "shippingbox", android: "inventory_2", web: "inventory_2" },
  },
  cab: {
    description: "A pickup or drop-off at your home",
    icon: { ios: "car", android: "local_taxi", web: "local_taxi" },
  },
  maintenance: {
    description: "Repairs and maintenance for your flat",
    icon: {
      ios: "wrench.and.screwdriver",
      android: "handyman",
      web: "handyman",
    },
  },
  other: {
    description: "Another reason to visit",
    icon: { ios: "ellipsis.circle", android: "more_horiz", web: "more_horiz" },
  },
} satisfies Record<
  string,
  { description: string; icon: SymbolViewProps["name"] }
>;
