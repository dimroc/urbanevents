import { compose, createStore as reduxCreateStore} from 'redux';
import { devTools, persistState } from 'redux-devtools';
import rootReducer from './reducers';

let finalCreateStore;
if (process.env.NODE_ENV === 'production') {
  finalCreateStore = reduxCreateStore.bind(null, rootReducer);
} else {
  try {
    finalCreateStore = compose(
      devTools(),
      persistState(window.location.href.match(/[?&]debug_session=([^&]+)\b/))
    )(reduxCreateStore).bind(null, rootReducer);
    console.log('dev tools added');
  } catch (e) {
    finalCreateStore = reduxCreateStore.bind(null, rootReducer);
  }
}

export const createStore = finalCreateStore;
