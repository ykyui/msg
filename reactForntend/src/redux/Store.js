import { configureStore } from "@reduxjs/toolkit";
import authReducer, { login, logout } from "./Auth";
import msgReducer from "./Msg";
import longpollReducer from "./Longpoll";
import agent from '../agent';
import { longpoll } from "./Longpoll";
import { initMsg } from "./Msg";

const localStorageMiddleware = (store) => (next) => (action) => {
    switch (action.type) {
        case login.fulfilled.type:
            // window.localStorage.setItem('jwt', action.payload);
            agent.setToken(action.payload);
            store.dispatch(initMsg())
            store.dispatch(longpoll())
            break;

        case logout.type:
            window.localStorage.removeItem('jwt');
            agent.setToken(undefined);
            break;
    }

    return next(action);
};

const store = configureStore({
    reducer: {
        auth: authReducer,
        longpoll: longpollReducer,
        msg:msgReducer
    },
    middleware: (getDefaultMiddleware) => [
        ...getDefaultMiddleware(),
        localStorageMiddleware
    ]
})

export default store