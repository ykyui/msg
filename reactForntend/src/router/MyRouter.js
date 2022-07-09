import { BrowserRouter, Route, Routes } from "react-router-dom"
import Home from "../components/page/Home"
import Login from "../components/page/login/Login"
import { useSelector } from 'react-redux';
import { isLogin } from "../redux/Auth";

export default () => {
    const _isLogin = useSelector(isLogin);
    return (
        <BrowserRouter>
            <Routes>
                {_isLogin ? <Route path="*" element={<Home />} /> : <Route path="*" element={<Login />} />}
                {/* <Route path="/repos" element={Repos} />
<Route path="/about" element={About} />
<Route path="/sayhi" render={props => <SayHi name="joe" {...props} />} /> */}
            </Routes>
        </BrowserRouter>
    )
}