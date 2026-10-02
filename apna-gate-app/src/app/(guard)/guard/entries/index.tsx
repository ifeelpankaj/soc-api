import { GuardLogsScreen } from "@/features/guard/components/logs/guard-logs-screen";
import { KeyboardAvoidingView, Platform } from "react-native";

export default function GuardEntriesRoute() {
  return (
    <KeyboardAvoidingView
      style={{ flex: 1 }}
      behavior={Platform.OS === "ios" ? "padding" : "height"}
    >
      <GuardLogsScreen />
    </KeyboardAvoidingView>
  );
}
