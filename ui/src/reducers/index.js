import {combineReducers} from "redux";
import postReducer from "./postReducer";
import authReducer from "./authReducer";
import tagReducer from "./tagReducer";


export default combineReducers({
    posts: postReducer,
    auth: authReducer,
    tags: tagReducer,
});
