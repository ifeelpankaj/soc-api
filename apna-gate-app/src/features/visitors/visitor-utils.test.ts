import assert from "node:assert/strict";
import { describe, it } from "node:test";

import type { ModelsVisitorEntry } from "@/lib/api/generated-api";

import {
  getCompanionDetails,
  getVisitorEmail,
  getVisitorName,
  getVisitorPhone,
  residentInvitePurposes,
  visitorPurposes,
} from "./visitor-utils";

type VisitorEntryTestShape = ModelsVisitorEntry & {
  email?: string;
  full_name?: string;
  phone_number?: string;
};

it("resident invite purposes exclude Service and Staff without changing guard purposes", () => {
  assert.deepEqual(residentInvitePurposes, ["guest", "delivery", "cab", "maintenance", "other"]);
  assert.ok(visitorPurposes.includes("service"));
  assert.ok(visitorPurposes.includes("staff"));
});

describe("visitor field resolvers", () => {
  it("prefers nested visitor values", () => {
    const entry: VisitorEntryTestShape = {
      email: "flat@example.com",
      full_name: "Flat Name",
      phone_number: "9999999999",
      visitor: {
        email: "nested@example.com",
        full_name: "Nested Name",
        phone_number: "8888888888",
      },
    };

    assert.equal(getVisitorName(entry), "Nested Name");
    assert.equal(getVisitorPhone(entry), "8888888888");
    assert.equal(getVisitorEmail(entry), "nested@example.com");
  });

  it("falls back to flattened values", () => {
    const entry: VisitorEntryTestShape = {
      email: "flat@example.com",
      full_name: "Flat Name",
      phone_number: "9999999999",
      visitor: {},
    };

    assert.equal(getVisitorName(entry), "Flat Name");
    assert.equal(getVisitorPhone(entry), "9999999999");
    assert.equal(getVisitorEmail(entry), "flat@example.com");
  });

  it("treats empty nested strings as missing", () => {
    const entry: VisitorEntryTestShape = {
      email: "flat@example.com",
      full_name: "Flat Name",
      phone_number: "9999999999",
      visitor: {
        email: " ",
        full_name: "",
        phone_number: "",
      },
    };

    assert.equal(getVisitorName(entry), "Flat Name");
    assert.equal(getVisitorPhone(entry), "9999999999");
    assert.equal(getVisitorEmail(entry), "flat@example.com");
  });
});

describe("companion resolvers", () => {
  it("normalizes supported companion API shapes", () => {
    const entry: ModelsVisitorEntry = {
      companion_details: [
        { full_name: "First Companion", phone_number: "7777777777" },
        { name: "Second Companion", phoneNumber: "6666666666" },
      ],
    };

    assert.deepEqual(getCompanionDetails(entry), [
      { name: "First Companion", phoneNumber: "7777777777" },
      { name: "Second Companion", phoneNumber: "6666666666" },
    ]);
  });

  it("filters completely empty companion objects", () => {
    const entry: ModelsVisitorEntry = {
      companion_details: [
        {},
        { full_name: "  ", phone_number: "" },
        { name: "Only Name" },
        { phoneNumber: "5555555555" },
      ],
    };

    assert.deepEqual(getCompanionDetails(entry), [
      { name: "Only Name", phoneNumber: "" },
      { name: "", phoneNumber: "5555555555" },
    ]);
  });
});
