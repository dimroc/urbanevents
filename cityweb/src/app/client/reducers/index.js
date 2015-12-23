import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router'
import cities from './cities';

export default combineReducers({
  cities,
  routing: routeReducer,
});
