import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

type AppState = {
  selectedFlatId: number | null;
  selectedSocietyId: number | null;
};

const initialState: AppState = {
  selectedFlatId: null,
  selectedSocietyId: null,
};

const appSlice = createSlice({
  name: "app",
  initialState,
  reducers: {
    setSelectedFlat: (state, action: PayloadAction<number>) => {
      state.selectedFlatId = action.payload;
      state.selectedSocietyId = null;
    },
    setSelectedSociety: (state, action: PayloadAction<number>) => {
      state.selectedSocietyId = action.payload;
      state.selectedFlatId = null;
    },
    clearWorkspaceSelection: (state) => {
      state.selectedFlatId = null;
      state.selectedSocietyId = null;
    },
  },
});

export const { clearWorkspaceSelection, setSelectedFlat, setSelectedSociety } =
  appSlice.actions;
export default appSlice.reducer;
