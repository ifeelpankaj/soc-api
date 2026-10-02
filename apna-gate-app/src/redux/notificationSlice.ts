import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

type NotificationState = {
  devicePushToken: string | null;
};

const initialState: NotificationState = {
  devicePushToken: null,
};

const notificationSlice = createSlice({
  name: "notifications",
  initialState,
  reducers: {
    setDevicePushToken: (state, action: PayloadAction<string | null>) => {
      state.devicePushToken = action.payload;
    },
  },
});

export const { setDevicePushToken } = notificationSlice.actions;
export default notificationSlice.reducer;
