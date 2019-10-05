import React, {useEffect} from 'react';
import {connect} from "react-redux";
import {fetchPosts} from "../../../actions";
import PostItem from "./PostListItem";
import {Link} from "react-router-dom";

const renderCreatePostButton = () => {
    return <div className="absolute right-0 mb-4 pb-4 mr-2">
        <Link to="/posts/create">
            <svg viewBox="0 0 24 24" width="24" height="24"
                 stroke="currentColor" strokeWidth="2" fill="none"
                 strokeLinecap="round" strokeLinejoin="round">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
        </Link>
    </div>
};

const PostList = ({fetchPosts, posts, loggedIn}) => {
    useEffect(
        () => {
            (async () =>  await fetchPosts())()
        },
        [fetchPosts]
    );

    return <div className="relative">
        <div style={{width: "100%"}}>
            {Object.values(posts).map((post) => <PostItem key={post.slug} post={post} />)}
        </div>
        {loggedIn? renderCreatePostButton(): null}
    </div>;
};

const mapStateToProps = (state) => {
    return {posts: state.posts, loggedIn: state.auth.token !== null}
};

export default connect(mapStateToProps, {fetchPosts})(PostList);
