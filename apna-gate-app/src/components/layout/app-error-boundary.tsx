import { Component, type ErrorInfo, type ReactNode } from "react";
import { Pressable, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type AppErrorBoundaryProps = {
  children: ReactNode;
  title?: string;
};

type AppErrorBoundaryState = {
  error: Error | null;
};

export class AppErrorBoundary extends Component<
  AppErrorBoundaryProps,
  AppErrorBoundaryState
> {
  state: AppErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error): AppErrorBoundaryState {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    if (__DEV__) {
      console.error("AppErrorBoundary caught:", error, info.componentStack);
    }
  }

  private handleRetry = () => {
    this.setState({ error: null });
  };

  render() {
    if (this.state.error) {
      return (
        <SafeAreaView style={styles.screen}>
          <View style={styles.content}>
            <Text style={styles.title}>{this.props.title ?? "Something went wrong"}</Text>
            <Text style={styles.message}>
              An unexpected error occurred. Please try again.
            </Text>
            <Pressable
              accessibilityRole="button"
              style={({ pressed }) => [styles.button, pressed ? styles.buttonPressed : null]}
              onPress={this.handleRetry}
            >
              <Text style={styles.buttonText}>Try again</Text>
            </Pressable>
          </View>
        </SafeAreaView>
      );
    }

    return this.props.children;
  }
}

const styles = StyleSheet.create({
  button: {
    backgroundColor: colors.guard.teal,
    borderRadius: 12,
    marginTop: spacing.lg,
    paddingHorizontal: spacing.xl,
    paddingVertical: spacing.md,
  },
  buttonPressed: {
    opacity: 0.88,
  },
  buttonText: {
    color: colors.text.inverse,
    fontSize: 15,
    fontWeight: "600",
    textAlign: "center",
  },
  content: {
    alignItems: "center",
    flex: 1,
    justifyContent: "center",
    paddingHorizontal: spacing.xl,
  },
  message: {
    color: colors.text.secondary,
    fontSize: 15,
    lineHeight: 22,
    marginTop: spacing.sm,
    textAlign: "center",
  },
  screen: {
    backgroundColor: colors.surface.screen,
    flex: 1,
  },
  title: {
    color: colors.text.primary,
    fontSize: 20,
    fontWeight: "700",
    textAlign: "center",
  },
});
