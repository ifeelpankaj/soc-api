"""Inspect exported maintenance PDF fixtures. Requires pypdf (QA only).

Generate fixtures using MAINTENANCE_PDF_QA_DIR and the document Go tests.
Pass --render-with /path/to/pdftoppm to also render every page for visual review.
"""
import argparse
from pathlib import Path
import subprocess
from pypdf import PdfReader

parser = argparse.ArgumentParser()
parser.add_argument("directory", type=Path)
parser.add_argument("--render-with")
args = parser.parse_args()

required = {"invoice.pdf", "invoice-long.pdf", "invoice-large-amount.pdf",
            "receipt-admin.pdf", "receipt-resident.pdf", "receipt-reversed.pdf"}
paths = sorted(args.directory.glob("*.pdf"))
assert required <= {p.name for p in paths}, "Generate renderer fixtures first"
texts = {}
for path in paths:
    reader = PdfReader(path)
    text = "\n".join(page.extract_text() for page in reader.pages)
    texts[path.name] = " ".join(text.split())
    assert "SECRET FORMER RESIDENT" not in text + str(reader.metadata)
    for number, page in enumerate(reader.pages, 1):
        assert abs(float(page.mediabox.width) - 595.28) < 1
        assert abs(float(page.mediabox.height) - 841.89) < 1
        assert f"Page {number} of {len(reader.pages)}" in page.extract_text()
    if "invoice" in path.name:
        assert "Not a payment receipt." in text
        assert "UPI destination" not in text and "Bank reference" not in text
    if "resident" in path.name or "other" in path.name or "reversed" in path.name:
        for secret in ("001234567890", "000PRIVATEUTR123", "ADMIN PRIVATE REVERSAL NOTE",
                       "Verified by user ID", "Reversed by user ID", "Payer user ID"):
            assert secret not in text + str(reader.metadata), (path, secret)
    if "reversed" in path.name:
        assert "REVERSED" in text and "refund" in text
    if path.name.startswith("integration-"):
        assert "Original Society" in text and "CHANGED LIVE NAME" not in text
        assert "Changed Road" not in text and "CHANGED-" not in text
    if args.render_with:
        dest = args.directory / "rendered"
        dest.mkdir(exist_ok=True)
        subprocess.run([args.render_with, "-png", "-r", "90", str(path),
                        str(dest / path.stem)], check=True)
    print(f"PASS {path.name}: {len(reader.pages)} page(s)")

assert "INR 2,234.56" in texts["invoice.pdf"]
assert "INR 92,233,720,368,547,758.07" in texts["invoice-large-amount.pdf"]
assert texts["invoice-long.pdf"].count("Shared services:") == 90
assert "001234567890" in texts["receipt-admin.pdf"]
for name in ("integration-receipt-admin.pdf", "integration-receipt-payer.pdf"):
    if name in texts:
        assert "000PRIVATEUTR123" in texts[name]
        assert "old@bank" in texts[name] and "new@bank" not in texts[name]
print("All content, privacy, pagination and metadata assertions passed.")
