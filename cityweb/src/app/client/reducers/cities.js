import { ActionTypes } from '#app/actions'

const initialState = {
  cities: [],
  current: {display: null}
}

export default function cities(state = initialState, action) {
  switch (action.type) {
    case ActionTypes.SET_CITIES:
      return {
        ...state,
        cities: action.cities
      };

    case ActionTypes.SET_CURRENT_CITY:
      var city = $.grep(state.cities, function(city) {
        if(city.key == action.cityKey) {
          return city;
        }
      })[0];

      return {
        ...state,
        current: city
      };

    default:
      return state;
  }
}
