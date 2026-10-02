import { type TextStyle } from "react-native";

export const typography = {
  eyebrow: {
    fontSize: 12,
    fontWeight: "700",
    letterSpacing: 1.2,
    textTransform: "uppercase",
  },
  titleLarge: {
    fontSize: 30,
    fontWeight: "700",
    lineHeight: 36,
  },
  title: {
    fontSize: 24,
    fontWeight: "700",
    lineHeight: 32,
  },
  subtitle: {
    fontSize: 16,
    fontWeight: "400",
    lineHeight: 24,
  },
  body: {
    fontSize: 16,
    fontWeight: "400",
    lineHeight: 24,
  },
  bodySmall: {
    fontSize: 14,
    fontWeight: "400",
    lineHeight: 20,
  },
  caption: {
    fontSize: 12,
    fontWeight: "400",
    lineHeight: 16,
  },
  label: {
    fontSize: 12,
    fontWeight: "600",
    letterSpacing: 2.24,
    textTransform: "uppercase",
  },
  button: {
    fontSize: 16,
    fontWeight: "600",
  },
  buttonCompact: {
    fontSize: 14,
    fontWeight: "700",
  },
} as const satisfies Record<string, TextStyle>;
