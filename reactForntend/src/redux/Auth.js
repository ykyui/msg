import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import agent from "../agent"

const initialState = () => {
    return { token: localStorage.getItem("token") ?? "" }
}

export const login = createAsyncThunk(
    'login.login',
    async (props, thunkAPI) => {
        const response = await agent.Auth.login(props.username, props.password)
        // store.dispatch(longpoll())
        return response.token
    })

export const authSlice = createSlice({
    name: "auth",
    initialState: initialState,
    extraReducers: {
        [login.fulfilled]: (state, { payload }) => {
            state.token = payload;
            console.log(payload)
            window.location.href = "/#"
        }
    },
    reducers: {
        logout: (state) => {
            localStorage.removeItem("token")
            state.token = ""
        }
    },
});

export const { logout } = authSlice.actions;

export const isLogin = (state) => state.auth.token != ""

export default authSlice.reducer;