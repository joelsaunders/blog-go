import _ from 'lodash';
import {FETCH_TAGS} from "../actions";

export default (state={}, action) => {
    switch (action.type) {
        case FETCH_TAGS:
            return _.mapKeys(action.payload, 'id');
        default:
            return state
    }
};
