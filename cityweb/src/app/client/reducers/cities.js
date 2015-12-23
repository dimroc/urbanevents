import { ActionTypes } from '#app/actions'
import { UPDATE_PATH } from 'redux-simple-router'

const initialState = {
  cities: [],
  across: []
}

export default function cities(state = initialState, action) {
  switch (action.type) {
    case ActionTypes.SET_ACROSS:
      return {
        ...state,
        across: action.cityGeoevents
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
