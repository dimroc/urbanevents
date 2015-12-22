import { ActionTypes } from '#app/actions'

const initialState = {
  q: null,
  geoevents: []
}

export default function geoevents(state = initialState, action) {
  switch (action.type) {
    case ActionTypes.SET_GEOEVENTS:
      return {
        q: action.q,
        geoevents: action.geoevents
      }
    default:
      return state
  }
}
