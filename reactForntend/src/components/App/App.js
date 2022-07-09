import './App.css';
import { Provider } from 'react-redux'
import store from '../../redux/Store'
import MyRouter from '../../router/MyRouter'
import 'antd/dist/antd.css';

function App() {
  if (!('indexedDB' in window)) {
    console.log('This browser doesn\'t support IndexedDB');
    return;
  }

  return (
    <Provider store={store}>
      <MyRouter />
    </Provider>
  );
}

export default App;
