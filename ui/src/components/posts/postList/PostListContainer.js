import React, {useEffect, useState} from 'react';
import {connect} from "react-redux";

import PostList from "./PostList";
import {fetchPosts} from "../../../actions";
import {ArrayParam, useQueryParam} from "use-query-params";

const PostListContainer = ({fetchPosts, posts, loggedIn}) => {
    const [tagModalActive, setTagModalActive] = useState(false);
    const [tagFilters, onSetTagFilter]  = useQueryParam("tag_name", ArrayParam);

    useEffect(
        () => {
            (async () =>  await fetchPosts({filters: {tag_name: tagFilters}}))()
        },
        [fetchPosts, tagFilters]
    );

    return <PostList
        loggedIn={loggedIn}
        posts={posts}
        setTagModalActive={setTagModalActive}
        tagModalActive={tagModalActive}
        onSetTagFilter={onSetTagFilter}
        tagFilters={tagFilters}
    />
};

const mapStateToProps = (state) => {
    return {posts: state.posts, loggedIn: state.auth.token !== null}
};

export default connect(mapStateToProps, {fetchPosts})(PostListContainer);