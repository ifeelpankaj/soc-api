import { useEffect, useRef, useState } from "react";
import { Button, AppText } from "@/components/ui";
import { Stack } from "@/components/layout";
import { useReadDocumentMutation } from "@/lib/api/document-api";
import { getFriendlyApiMessage } from "@/features/auth/api-error";
import { openDocument } from "./open-document";

export function DocumentAction({
  path,
  filename,
  title,
}: {
  path: string;
  filename: string;
  title: string;
}) {
  const [read, { reset }] = useReadDocumentMutation();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string>();
  const active = useRef(false);
  const inFlight = useRef(false);
  useEffect(() => {
    active.current = true;
    return () => {
      active.current = false;
    };
  }, []);
  const open = async () => {
    if (inFlight.current) return;
    inFlight.current = true;
    setBusy(true);
    setError(undefined);
    try {
      const bytes = await read({ path }).unwrap();
      if (active.current) await openDocument(bytes, filename);
    } catch (e) {
      if (active.current)
        setError(
          getFriendlyApiMessage(
            e,
            "Could not open the document. Please retry.",
          ),
        );
    } finally {
      reset();
      inFlight.current = false;
      if (active.current) setBusy(false);
    }
  };
  return (
    <Stack gap="sm">
      <Button
        title={title}
        variant="secondary"
        loading={busy}
        onPress={() => void open()}
      />
      {error ? (
        <AppText color="error" variant="bodySmall">
          {error}
        </AppText>
      ) : null}
    </Stack>
  );
}
