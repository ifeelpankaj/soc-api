import { StyleSheet } from "react-native";

import { Row, Stack } from "@/components/layout";
import { AppText } from "@/components/ui/app-text";
import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

import { Badge } from "./badge";
import { Card } from "./card";

type ProfileHeroProps = {
  name: string;
  subtitle?: string;
  email?: string | null;
  badge?: string;
};

export function ProfileHero({ name, subtitle, email, badge }: ProfileHeroProps) {
  return (
    <Card style={styles.card}>
      <Row align="flex-start" gap="lg">
        <Stack gap="sm" style={styles.copy}>
          <AppText variant="eyebrow" color="#99f6e4">
            Profile
          </AppText>
          <AppText variant="titleLarge" color={colors.text.inverse}>
            {name}
          </AppText>
          {subtitle ? (
            <AppText variant="bodySmall" color="#cbd5e1">
              {subtitle}
            </AppText>
          ) : null}
          {email ? (
            <AppText variant="bodySmall" color="#94a3b8">
              {email}
            </AppText>
          ) : null}
        </Stack>
        {badge ? <Badge label={badge} tone="emerald" /> : null}
      </Row>
    </Card>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: "#020617",
    gap: spacing.lg,
  },
  copy: {
    flex: 1,
  },
});
