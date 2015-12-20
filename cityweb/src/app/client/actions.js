import { pushPath } from 'redux-simple-router'

/* CityWeb Action Types */
var keyMirror = require('keymirror')
var actionTypes = keyMirror({
  SET_GEOEVENTS: null,
  SET_CITIES: null,
  SET_CURRENT_CITY: null,
  GET_CITY: null
})

export const ActionTypes = actionTypes

/* CityWeb Action Creators */

export function getGeoeventsAsync(cityKey, q) {
  return (dispatch, getState) => {
    if (!q || q.length == 0) {
      return (dispatch) => { dispatch(setGeoevents([])) }
    }

    let url = "/api/v1/cities/" + cityKey + "/search?q=" + q;

    return fetch(url).then((result) => {
      return result.json();
    }).then(geoevents => {
      dispatch(setGeoevents(geoevents))
      dispatch(pushPath('/' + cityKey + '?q=' + q), getState());
    });
  }
}

export function setGeoevents(geoevents) {
  return { type: actionTypes.SET_GEOEVENTS, geoevents };
}

export function getCitiesAsync() {
  return (dispatch) => {
    return fetch('/api/v1/cities').then(function(result) {
      return result.json();
    }).then(cities => { dispatch(setCities(cities)) });
  };
}

export function setCities(cities) {
  return { type: actionTypes.SET_CITIES, cities };
}

export function setCurrentCity(cityKey) {
  return { type: actionTypes.SET_CURRENT_CITY, cityKey };
}

export function getCurrentCity() {
  return { type: actionTypes.GET_CURRENT_CITY };
}
