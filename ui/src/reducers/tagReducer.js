import _ from 'lodash';
import {FETCH_TAGS, EDIT_TAG, CREATE_TAG, DELETE_TAG} from "../actions";

export default (state={}, action) => {
    switch (action.type) {
        case FETCH_TAGS:
            return _.mapKeys(action.payload, 'id');
        case EDIT_TAG:
            return {...state, [action.payload.id]: action.payload};
        case CREATE_TAG:
            return {...state, [action.payload.id]: action.payload};
        case DELETE_TAG:
            const {[action.payload.id]: __, ...withoutID} = state;
            return withoutID;
        default:
            return state
    }
};
