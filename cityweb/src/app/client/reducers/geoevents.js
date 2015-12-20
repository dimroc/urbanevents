import { ActionTypes } from '#app/actions'

const initialState = []

export default function geoevents(state = initialState, action) {
  switch (action.type) {
    case ActionTypes.SET_GEOEVENTS:
      return action.geoevents
    default:
      return state
  }
}
