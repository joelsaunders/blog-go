import React from 'react';
import _ from 'lodash';

import PostItem from "./PostListItem";
import {Link} from "react-router-dom";
import TagModal from "../../tags/TagModal";

const renderCreatePostButton = (setTagModalActive) => {
    return <div className="float-right flex flex-row mt-2 pb-4 mr-2">
        <button onClick={setTagModalActive} className="mr-2">
            Manage Tags
        </button>
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

// Helper function to add/remove items from the filter array
const addRemoveTag = (tag, tagList) => {
    if (_.includes(tagList, tag)) {
        return _.filter(tagList, (item) => item !== tag)
    } else {
        return [...(tagList || []), tag]
    }
};

const PostList = ({posts, loggedIn, setTagModalActive, tagModalActive, onSetTagFilter, tagFilters}) => {
    return <div className="relative">
        {loggedIn? renderCreatePostButton(setTagModalActive): null}
        <div style={{width: "100%"}}>
            {
                Object.values(posts).map((post) => <PostItem
                    setTagFilter={(value) => onSetTagFilter(addRemoveTag(value, tagFilters))}
                    selectedTags={tagFilters}
                    key={post.slug}
                    post={post}
                />)
            }
        </div>
        {tagModalActive? <TagModal onDismiss={() => setTagModalActive(false)} />: null}
    </div>;
};

export default PostList;
