import { createStore, combineReducers, applyMiddleware, compose } from 'redux';
import userReducer from './reducers/userReducer';
import { composeWithDevTools } from 'remote-redux-devtools'
import thunk from 'redux-thunk';

const initialState = {};

const middleware = [thunk];

const reducers = combineReducers({
  user: userReducer,
});

const composeEnhancers = composeWithDevTools({
  realtime: true,
  name: 'Your Instance Name',
  hostname: 'localhost',
  port: 3000 // the port your remotedev server is running at
})

const enhancer = composeEnhancers(applyMiddleware(...middleware));

const store = createStore(reducers, initialState, enhancer);

export default store;
