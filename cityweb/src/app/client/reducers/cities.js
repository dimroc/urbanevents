import { ActionTypes } from '#app/actions'

const initialState = {}

export default function cities(state = initialState, action) {
  switch (action.type) {
    case ActionTypes.GET_CITIES:
      console.log(arguments);
      return {}
    default:
      return state
  }
}

export function getCityName(state, cityKey) {
  return "New York City"
}
