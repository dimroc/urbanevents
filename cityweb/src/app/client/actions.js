/* CityWeb Action Types */
var keyMirror = require('keymirror')
var actionTypes = keyMirror({
  SET_GEOEVENTS: null,
  SET_CITIES: null,
  GET_CITY: null
})

export const ActionTypes = actionTypes

/* CityWeb Action Creators */

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
