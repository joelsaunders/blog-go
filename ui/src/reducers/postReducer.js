import _ from 'lodash';
import {CREATE_POST, EDIT_POST, FETCH_POSTS} from "../actions";

export default (state={}, action) => {
    switch (action.type) {
        case FETCH_POSTS:
            console.log(action.type);
            console.log(action.payload);
            return _.mapKeys(action.payload, 'slug');
        case EDIT_POST :
            return {...state, [action.payload.slug]: action.payload};
        case CREATE_POST:
            return {...state, [action.payload.slug]: action.payload};
        default:
            return state
    }
};
