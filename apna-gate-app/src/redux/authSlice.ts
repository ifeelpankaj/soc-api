import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

import type { ModelsUserResponse } from "@/lib/api/generated-api";

type AuthState = {
  user: ModelsUserResponse | null;
};

const initialState: AuthState = {
  user: null,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setCredentials: (
      state,
      action: PayloadAction<{ user: ModelsUserResponse | null }>,
    ) => {
      state.user = action.payload.user;
    },
    clearAuth: (state) => {
      state.user = null;
    },
  },
});

export const { clearAuth, setCredentials } = authSlice.actions;
export default authSlice.reducer;
