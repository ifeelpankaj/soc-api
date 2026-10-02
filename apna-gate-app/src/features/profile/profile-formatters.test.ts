import assert from "node:assert/strict";
import { describe, it } from "node:test";

import {
  formatDateOfBirthForInput,
  normalizeProfileDateInput,
} from "./profile-formatters";

describe("formatDateOfBirthForInput", () => {
  it("returns empty string for blank values", () => {
    assert.equal(formatDateOfBirthForInput(""), "");
    assert.equal(formatDateOfBirthForInput(null), "");
    assert.equal(formatDateOfBirthForInput(undefined), "");
  });

  it("extracts YYYY-MM-DD from RFC3339 timestamps", () => {
    assert.equal(formatDateOfBirthForInput("1990-01-15T00:00:00Z"), "1990-01-15");
  });

  it("keeps existing YYYY-MM-DD values", () => {
    assert.equal(formatDateOfBirthForInput("1990-01-15"), "1990-01-15");
  });
});

describe("normalizeProfileDateInput", () => {
  it("accepts RFC3339 timestamps", () => {
    assert.deepEqual(normalizeProfileDateInput("1990-01-15T00:00:00Z"), {
      ok: true,
      value: "1990-01-15",
    });
  });

  it("round-trips formatted input before save", () => {
    const formatted = formatDateOfBirthForInput("1990-01-15T00:00:00Z");
    assert.deepEqual(normalizeProfileDateInput(formatted), {
      ok: true,
      value: "1990-01-15",
    });
  });

  it("rejects unsupported formats", () => {
    assert.deepEqual(normalizeProfileDateInput("not-a-date"), {
      ok: false,
      error: "Use YYYY-MM-DD format for date of birth.",
    });
  });
});
