import { ActionTypes } from '#app/actions'
import { UPDATE_PATH } from 'redux-simple-router'

const initialState = {
  cities: [],
  current: {display: null}
}

export default function cities(state = initialState, action) {
  switch (action.type) {
    case UPDATE_PATH:
      console.log("Handling update path", arguments)
      return state

    case ActionTypes.SET_ACROSS:
      let cityGeoevents = state.cities.map((city) => {
        city.geoevents = action.cityGeoevents[city.key] || []
        return city
      })

      return {
        ...state,
        cities: cityGeoevents
      }

    case ActionTypes.SET_CITIES:
      return {
        ...state,
        cities: action.cities
      };

    default:
      return state;
  }
}
