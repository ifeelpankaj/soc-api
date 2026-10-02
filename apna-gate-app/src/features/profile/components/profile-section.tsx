import type { ReactNode } from "react";
import { StyleSheet } from "react-native";

import { Stack } from "@/components/layout";
import { Card, SectionHeader } from "@/components/ui";

type ProfileSectionProps = {
  children: ReactNode;
  title: string;
};

export function ProfileSection({ children, title }: ProfileSectionProps) {
  return (
    <Stack gap={10}>
      <SectionHeader title={title} />
      <Card style={styles.card}>{children}</Card>
    </Stack>
  );
}

const styles = StyleSheet.create({
  card: {
    overflow: "hidden",
    padding: 0,
  },
});
