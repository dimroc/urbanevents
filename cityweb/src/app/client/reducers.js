import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router'
import { SET_CONFIG } from './actions';
import geoevents from './reducers/geoevents';

function config(state = {}, action) {
  switch (action.type) {
    case SET_CONFIG:
      return action.config;
    default:
      return state;
  }
}

//export default combineReducers({config, geoevents});
export default combineReducers(Object.assign({}, config, geoevents, {
                                            routing: routeReducer
}));
