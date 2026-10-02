import assert from "node:assert/strict";
import { describe, it } from "node:test";

import {
  buildMemberInviteShareContent,
  buildVisitorInviteShareContent,
} from "./invite-share-content";

const shortLink = {
  code: "Ab3Xy9",
  url: "https://apnagate.org/link/Ab3Xy9",
};

describe("invite share builders", () => {
  it("builds visitor share content from the short link URL", () => {
    const result = buildVisitorInviteShareContent({
      expectedAt: "2026-08-30T13:30:00Z",
      expiresAt: "2026-08-30T16:30:00Z",
      flatLabel: "A-101",
      purpose: "guest",
      shortLink,
      societyName: "Green Valley Residency",
    });

    assert.equal(result.available, true);
    if (!result.available) return;

    assert.equal(result.content.title, "Visitor invite");
    assert.equal(result.content.url, shortLink.url);
    assert.match(result.content.message, /You're invited to visit Green Valley Residency/);
    assert.match(result.content.message, /Purpose: Guest/);
    assert.match(result.content.message, /Flat: A-101/);
    assert.match(result.content.message, /Expected:/);
    assert.match(result.content.message, /Complete your visitor details here:\nhttps:\/\/apnagate\.org\/link\/Ab3Xy9/);
    assert.doesNotMatch(result.content.message, /\/visit\//);
  });

  it("builds member share content from the short link URL", () => {
    const result = buildMemberInviteShareContent({
      expiresAt: "2026-08-31T18:29:00Z",
      flatLabel: "A-101",
      fullName: "Rahul",
      role: "family",
      shortLink,
      societyName: "Green Valley Residency",
    });

    assert.equal(result.available, true);
    if (!result.available) return;

    assert.equal(result.content.title, "Member invite");
    assert.equal(result.content.url, shortLink.url);
    assert.match(result.content.message, /Rahul, you've been invited to join Green Valley Residency/);
    assert.match(result.content.message, /Role: Family/);
    assert.match(result.content.message, /Accept your invitation:\nhttps:\/\/apnagate\.org\/link\/Ab3Xy9/);
    assert.doesNotMatch(result.content.message, /\/join\/flat\//);
  });

  it("returns unavailable when short_link.url is missing", () => {
    assert.deepEqual(
      buildVisitorInviteShareContent({
        purpose: "guest",
        shortLink: { code: "Ab3Xy9" },
      }),
      { available: false, reason: "short_link_unavailable" },
    );
  });

  it("omits visitor expected time without blank-line artifacts", () => {
    const result = buildVisitorInviteShareContent({
      purpose: "delivery",
      shortLink,
      societyName: "Green Valley Residency",
    });

    assert.equal(result.available, true);
    if (!result.available) return;

    assert.doesNotMatch(result.content.message, /^Expected:/m);
    assert.doesNotMatch(result.content.message, /\n\n\n/);
  });

  it("omits expiry sentence when expires_at is missing", () => {
    const result = buildMemberInviteShareContent({
      fullName: "Rahul",
      role: "tenant",
      shortLink,
      societyName: "Green Valley Residency",
    });

    assert.equal(result.available, true);
    if (!result.available) return;

    assert.doesNotMatch(result.content.message, /expires on/);
    assert.doesNotMatch(result.content.message, /\n\n\n/);
  });
});
