import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router'
import geoevents from './geoevents';
import cities from './cities';

export default combineReducers({
  cities,
  geoevents,
  routing: routeReducer,
});
