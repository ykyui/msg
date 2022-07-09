import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import store from './Store'
import agent from "../agent"

const initialState = () => {
    return {}
}

export const longpoll = createAsyncThunk(
    'longpoll.longpoll',
    async (props, thunkAPI) => {
        let result
        try {
            result = await agent.Msg.longpollMsg()
        } catch (e) {
            await new Promise(resolve => setTimeout(resolve, 3000))
        }
        store.dispatch(longpoll())
        return result
    })

export const longpollSlice = createSlice({
    name: "longpoll",
    initialState: initialState,
    extraReducers: {
        [longpoll.fulfilled]: (state, { payload }) => {
            console.log(payload)
        }
    },
});


export default longpollSlice.reducer;