import React from 'react';
import PostDetail from "./postDetail";
import useCurrentPost from "../../../hooks/useCurrentPost";

const PostDetailContainer = (props) => {
    const postSlug = props.match.params.slug;
    const currentPost = useCurrentPost(postSlug);

    return <PostDetail
        currentPost={currentPost}
    />
};

export default PostDetailContainer;
