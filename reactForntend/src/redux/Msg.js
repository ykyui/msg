import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import agent from "../agent"
import store from "./Store";
import { msgContent } from '../idb/idb';


const initialState = () => {
    return { contacts: [] }
}

export const initMsg = createAsyncThunk(
    'initMsg.init',
    async (props, thunkAPI) => {
        const msg = await msgContent()
        return msg
    })

export const msgSlice = createSlice({
    name: "msg",
    initialState: initialState,
    extraReducers: {
        [initMsg.fulfilled]: (state, { payload }) => {
            console.log("init msg ",payload)
        }
    },
    reducers: {

    },
});

export const { } = msgSlice.actions;

export default msgSlice.reducer;