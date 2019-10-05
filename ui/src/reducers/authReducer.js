import {SIGN_IN, SIGN_OUT} from "../actions";

export default (state={user: null, token: null}, action) => {
    switch (action.type) {
        case SIGN_IN:
            return {
                ...state,
                token: action.payload.token,
                user: action.payload.user
            };
        case SIGN_OUT:
            return {
                ...state,
                token: null,
                user: null
            };
        default:
            return state
    }
}
