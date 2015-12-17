import { combineReducers } from 'redux';
import { routeReducer } from 'redux-simple-router'
import { SET_CONFIG } from '#app/actions';
import geoevents from './geoevents';
import cities from './cities';

function config(state = {}, action) {
  switch (action.type) {
    case SET_CONFIG:
      return action.config;
    default:
      return state;
  }
}

export default combineReducers(Object.assign({}, config, geoevents, {
                                            routing: routeReducer
}));
