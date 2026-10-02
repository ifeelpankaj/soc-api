import assert from "node:assert/strict";
import test from "node:test";
import {
  API_TIMEOUT_MS,
  isAuthenticationFailure,
  isReadRequest,
} from "./request-policy";
import { getApiMessage } from "./api-error";

test("only reads are eligible for automatic replay after authentication recovery", () => {
  assert.equal(isReadRequest("/visitors"), true);
  assert.equal(isReadRequest({ method: "get" }), true);
  for (const method of ["POST", "PATCH", "PUT", "DELETE"])
    assert.equal(isReadRequest({ method }), false);
});
test("temporary failures never invalidate a saved session", () => {
  for (const status of ["FETCH_ERROR", "TIMEOUT_ERROR", 429, 500, 503])
    assert.equal(isAuthenticationFailure({ status }), false);
  assert.equal(isAuthenticationFailure({ status: 401 }), true);
  assert.equal(isAuthenticationFailure({ status: 403 }), true);
  assert.equal(API_TIMEOUT_MS, 30_000);
});
test("network and timeout messages offer recovery without leaking transport details", () => {
  assert.match(
    getApiMessage({ status: "FETCH_ERROR", error: "private host" }, "Failed"),
    /connection/,
  );
  assert.match(
    getApiMessage({ status: "TIMEOUT_ERROR" }, "Failed"),
    /try again/,
  );
});
