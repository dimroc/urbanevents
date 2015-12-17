/**
 * action types
 */

export const SET_CONFIG = 'SET_CONFIG';

/**
 * action creators
 */

export function setConfig(config) {
  return { type: SET_CONFIG, config };
}


/* CityWeb Action Types */
var keyMirror = require('keymirror')
var actionTypes = keyMirror({
  SET_GEOEVENTS: null,
  GET_CITIES: null
})

export const ActionTypes = actionTypes

/* CityWeb Action Creators */

export function setGeoevents(geoevents) {
  return { type: actionTypes.SET_GEOEVENTS, geoevents };
}

export function getCities() {
  return { type: actionTypes.GET_CITIES };
}
