import { File, Paths } from "expo-file-system";
import { isAvailableAsync, shareAsync } from "expo-sharing";

export async function openDocument(bytes: number[], filename: string) {
  if (!(await isAvailableAsync()))
    throw new Error("Document sharing is unavailable on this device.");
  const file = new File(Paths.cache, filename);
  file.write(new Uint8Array(bytes));
  // Android's chooser can resolve before the receiving app reads the file.
  // Leave the downloaded response in the OS-managed temporary cache, rather
  // than deleting the receiver's input immediately after choosing an app.
  await shareAsync(file.uri, {
    mimeType: "application/pdf",
    UTI: "com.adobe.pdf",
    dialogTitle: "View or save document",
  });
}
