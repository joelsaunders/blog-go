import theBookOfJoel from "../apis/theBookOfJoel";
import {saveState} from "../localStorage";

export const FETCH_POSTS = 'FETCH_POSTS';
export const EDIT_POST = 'EDIT_POST';
export const CREATE_POST = 'CREATE_POST';
export const DELETE_POST = 'DELETE_POST';
export const SIGN_IN = 'SIGN_IN';
export const SIGN_OUT = 'SIGN_OUT';
export const FETCH_TAGS = 'FETCH_TAGS';
export const EDIT_TAG = 'EDIT_TAG';
export const DELETE_TAG = 'DELETE_TAG';
export const CREATE_TAG = 'CREATE_TAG';


export const fetchTags = () => async dispatch => {
    const response = await theBookOfJoel.get('api/v1/tags');
    dispatch({type: FETCH_TAGS, payload: response.data})
};

export const createTag = (data) => async (dispatch, getState) => {
    const token = getState().auth.token;

    try {
        const response = await theBookOfJoel.post(
            'api/v1/tags',
            {name: data},
            {headers: {Authorization: `Bearer ${token}`}}
        );
        dispatch({type: CREATE_TAG, payload: response.data});
    } catch (err) {
        if (err.response.status === 401) {
            dispatch(signOut())
        }
        console.log(err);
    }
};

export const editTag = (tagID, data) => async (dispatch, getState) => {
    const token = getState().auth.token;

    try {
        const response = await theBookOfJoel.patch(
            `api/v1/tags/${tagID}`,
            {id: tagID, name: data},
            {headers: {Authorization: `Bearer ${token}`}}
        );
        dispatch({type: EDIT_TAG, payload: response.data});
        dispatch(fetchPosts());
    } catch (err) {
        if (err.response.status === 401) {
            dispatch(signOut())
        }
        console.log(err);
    }
};

export const deleteTag = (tagID) => async (dispatch, getState) => {
    const token = getState().auth.token;

    try {
        await theBookOfJoel.delete(
            `/api/v1/tags/${tagID}`,
            {headers: {Authorization: `Bearer ${token}`}}
        );
        dispatch({type: DELETE_TAG, payload: {id: tagID}});
        dispatch(fetchPosts());
    } catch (err) {
        if (err.response.status === 401) {
            dispatch(signOut())
        }
        console.log(err);
    }
};

export const fetchPosts = (tagFilter) => async dispatch => {
    const baseUrlString = "api/v1/posts";
    const urlString = tagFilter ? `${baseUrlString}?tag_name=${tagFilter}`: baseUrlString;

    const response = await theBookOfJoel.get(urlString);
    dispatch({type: FETCH_POSTS, payload: response.data})
};

export const editPost = (postId, data) => async (dispatch, getState) => {
    const token = getState().auth.token;

    try {
        const response = await theBookOfJoel.patch(
            `/api/v1/posts/${postId}`,
            data,
            {headers: {Authorization: `Bearer ${token}`}}
        );
        dispatch({type: EDIT_POST, payload: {responseData: response.data, originalSlug: postId}})
    } catch (err) {
        if (err.response.status === 401) {
            dispatch(signOut())
        }
        console.log(err);
    }
};

export const createPost = (data) => async (dispatch, getState) => {
    const {token, user} = getState().auth;

    try {
        const response = await theBookOfJoel.post(
            `api/v1/posts`,
            {...data, ...{author: user, published: true}},
            {headers: {Authorization: `Bearer ${token}`}}
        );
        dispatch({type: CREATE_POST, payload: response.data})
    } catch (err) {
        if (err.response.status === 401) {
            dispatch(signOut())
        }
        console.log(err);
    }
};

export const deletePost = (postSlug) => async (dispatch, getState) => {
    const {token} = getState().auth;

    try {
        await theBookOfJoel.delete(
            `/api/v1/posts/${postSlug}`,
            {headers: {Authorization: `Bearer ${token}`}}
        );
        dispatch({type: DELETE_POST, payload: {slug: postSlug}})
    } catch (err) {
        if (err.response.status === 401) {
            dispatch(signOut())
        }
        console.log(err);
    }
};

export const signIn = (email, token) => {
    const payload = {user: email, token: token};
    saveState({auth: payload});
    return {type: SIGN_IN, payload: payload}
};

export const signOut = () => {
    const payload = {user: null, token: null};
    saveState({auth: payload});
    return {type: SIGN_OUT, payload: payload}
};
