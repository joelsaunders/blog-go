import React from 'react';
import FormikPostForm from "../../forms/postform/PostForm";
import {createPost} from "../../../actions";
import {connect} from "react-redux";
import PostCreateComponent from "./PostCreateComponent";

const PostCreateContainer = (props) => {

    return <PostCreateComponent>
        <FormikPostForm
            submitAction={props.createPost}
            match={props.match}
        />
    </PostCreateComponent>;
};

export default connect(null, {createPost})(PostCreateContainer);
