import React from 'react';
import {editPost} from "../../../actions";
import {connect} from "react-redux";
import FormikPostForm from "../../forms/postform/PostForm";
import useCurrentPost from "../../../hooks/useCurrentPost";
import PostEditComponent from "./PostEditComponent";


const PostEditContainer = (props) => {
    const postSlug = props.match.params.slug;
    const currentPost = useCurrentPost(postSlug);


    if (!currentPost) {
        return <div className="ui very padded loading segment" />
    }

    return <PostEditComponent currentPost={currentPost}>
        <FormikPostForm
            submitAction={(payload) => props.editPost(postSlug, payload)}
            post={currentPost}
            match={props.match}
        />
    </PostEditComponent>;
};

export default connect(null, {editPost})(PostEditContainer);
