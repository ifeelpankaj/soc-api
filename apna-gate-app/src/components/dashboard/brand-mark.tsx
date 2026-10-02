import { Image } from "expo-image";
import { StyleSheet, Text } from "react-native";

import { Row } from "@/components/layout";
import { colors } from "@/theme/colors";

type BrandMarkProps = {
  size?: "sm" | "md";
};

export function BrandMark({ size = "md" }: BrandMarkProps) {
  const iconSize = size === "sm" ? 32 : 36;
  const fontSize = size === "sm" ? 18 : 20;

  return (
    <Row align="center" gap="sm" justify="flex-start">
      <Image
        accessibilityLabel="Apna Gate logo"
        source={require("@/assets/images/public/logo.png")}
        style={{ height: iconSize, width: iconSize }}
      />
      <Text style={[styles.brandText, { fontSize }]}>Apna Gate</Text>
    </Row>
  );
}

const styles = StyleSheet.create({
  brandText: {
    color: colors.brand.navy,
    fontWeight: "800",
    letterSpacing: -0.5,
  },
});
