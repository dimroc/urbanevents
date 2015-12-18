import { compose, applyMiddleware, createStore as reduxCreateStore} from 'redux';
import thunk from 'redux-thunk'
import { devTools, persistState } from 'redux-devtools';
import rootReducer from '#app/reducers/index.js';

let finalCreateStore;
const middleware = applyMiddleware(thunk);

if (process.env.NODE_ENV === 'production') {
  finalCreateStore = compose(middleware)(reduxCreateStore).bind(null, rootReducer);
} else {
  try {
    finalCreateStore = compose(
      middleware,
      devTools(),
      persistState(window.location.href.match(/[?&]debug_session=([^&]+)\b/))
    )(reduxCreateStore).bind(null, rootReducer);
    console.log('dev tools added');
  } catch (e) {
    finalCreateStore = compose(middleware)(reduxCreateStore).bind(null, rootReducer);
  }
}

export const createStore = finalCreateStore;
