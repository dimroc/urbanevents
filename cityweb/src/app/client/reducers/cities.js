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
      var city = null;
      var index = state.cities.forEach((cityEntry) => {
        if(cityEntry.key === action.cityKey) {
          city = cityEntry
          return;
        }
      });

      console.log("## Selected city", city.display);
      return {
        ...state,
        current: city
      };

    default:
      return state;
  }
}
