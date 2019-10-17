import _ from 'lodash';
import {CREATE_POST, EDIT_POST, FETCH_POSTS} from "../actions";

export default (state={}, action) => {
    switch (action.type) {
        case FETCH_POSTS:
            return _.mapKeys(action.payload, 'slug');
        case EDIT_POST :
            //handle the case where the slug may have changed
            const newState = _.mapKeys(state, (value, key) => {
                if (key === action.payload.originalSlug) {
                    return action.payload.responseData.slug
                }
                return key
            });
            return {...newState, [action.payload.responseData.slug]: action.payload.responseData};
        case CREATE_POST:
            return {...state, [action.payload.slug]: action.payload};
        default:
            return state
    }
};
