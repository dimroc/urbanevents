import { compose, applyMiddleware, createStore as reduxCreateStore} from 'redux';
import thunk from 'redux-thunk'
import { devTools, persistState } from 'redux-devtools';
import rootReducer from '#app/reducers/index.js';

const middleware = applyMiddleware(thunk);

let finalCreateStore = compose(
  middleware,
  devTools()
)(reduxCreateStore).bind(null, rootReducer);

export const createStore = finalCreateStore;
